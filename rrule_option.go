package rrule

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ROption offers options to construct a RRule instance.
// For performance, it is strongly recommended providing explicit ROption.Dtstart, which defaults to `time.Now().UTC().Truncate(time.Second)`.
type ROption struct {
	Freq       Frequency
	Dtstart    time.Time
	Interval   int
	Wkst       Weekday
	Count      int
	Until      time.Time
	Bysetpos   []int
	Bymonth    []int
	Bymonthday []int
	Byyearday  []int
	Byweekno   []int
	Byweekday  []Weekday
	Byhour     []int
	Byminute   []int
	Bysecond   []int
	Byeaster   []int
	AllDay     bool
}

// String returns RRULE string with DTSTART if exists. e.g.
//
//	DTSTART;TZID=America/New_York:19970105T083000
//	RRULE:FREQ=YEARLY;INTERVAL=2;BYMONTH=1;BYDAY=SU;BYHOUR=8,9;BYMINUTE=30
func (option *ROption) String() string {
	str := option.RRuleString()
	if option.Dtstart.IsZero() {
		return str
	}

	// Handle AllDay events: use DATE format as per RFC 5545
	if option.AllDay {
		// All-day events should use VALUE=DATE format
		year, month, day := option.Dtstart.Date()
		dateStr := fmt.Sprintf("%04d%02d%02d", year, int(month), day)
		return fmt.Sprintf("DTSTART;VALUE=DATE:%s\nRRULE:%s", dateStr, str)
	}

	// Non-all-day events use DATE-TIME format
	return fmt.Sprintf("DTSTART%s\nRRULE:%s", timeToRFCDatetimeStr(option.Dtstart), str)
}

// RRuleString returns RRULE string exclude DTSTART
func (option *ROption) RRuleString() string {
	result := []string{fmt.Sprintf("FREQ=%v", option.Freq)}
	if option.Interval != 0 {
		result = append(result, fmt.Sprintf("INTERVAL=%v", option.Interval))
	}
	if option.Wkst != MO {
		result = append(result, fmt.Sprintf("WKST=%v", option.Wkst))
	}
	// Only include COUNT when it is a positive integer.
	// Negative or zero COUNT means "unlimited" and should not appear in RRULE output.
	if option.Count > 0 {
		result = append(result, fmt.Sprintf("COUNT=%v", option.Count))
	}
	if !option.Until.IsZero() {
		// RFC 5545: UNTIL value type must match DTSTART value type
		// For all-day events (floating time), UNTIL should also use floating time
		if option.AllDay {
			// For all-day events, use DATE format (no time part) as per RFC 5545
			// Convert to date at 00:00:00 in UTC to represent floating time
			floatingTime := time.Date(option.Until.Year(), option.Until.Month(), option.Until.Day(), 0, 0, 0, 0, time.UTC)
			result = append(result, fmt.Sprintf("UNTIL=%v", floatingTime.Format(DateFormat)))
		} else {
			result = append(result, fmt.Sprintf("UNTIL=%v", timeToStr(option.Until)))
		}
	}
	result = appendIntsOption(result, "BYSETPOS", option.Bysetpos)
	result = appendIntsOption(result, "BYMONTH", option.Bymonth)
	result = appendIntsOption(result, "BYMONTHDAY", option.Bymonthday)
	result = appendIntsOption(result, "BYYEARDAY", option.Byyearday)
	result = appendIntsOption(result, "BYWEEKNO", option.Byweekno)
	if len(option.Byweekday) != 0 {
		valueStr := make([]string, len(option.Byweekday))
		for i, wday := range option.Byweekday {
			valueStr[i] = wday.String()
		}
		result = append(result, fmt.Sprintf("BYDAY=%s", strings.Join(valueStr, ",")))
	}
	// For all-day events, time-of-day components are ignored by engines,
	// so we omit BYHOUR/BYMINUTE/BYSECOND in the RRULE output for better interoperability.
	if !option.AllDay {
		result = appendIntsOption(result, "BYHOUR", option.Byhour)
		result = appendIntsOption(result, "BYMINUTE", option.Byminute)
		result = appendIntsOption(result, "BYSECOND", option.Bysecond)
	}
	result = appendIntsOption(result, "BYEASTER", option.Byeaster)
	return strings.Join(result, ";")
}

// StrToROption converts string to ROption.
func StrToROption(rfcString string) (*ROption, error) {
	return StrToROptionInLocation(rfcString, time.UTC)
}

// StrToROptionInLocation is same as StrToROption but in case local
// time is supplied as date-time/date field (ex. UNTIL), it is parsed
// as a time in a given location (time zone)
func StrToROptionInLocation(rfcString string, loc *time.Location) (*ROption, error) {
	rfcString = strings.TrimSpace(rfcString)
	strs := strings.Split(rfcString, "\n")
	var rruleStr, dtstartStr string
	switch len(strs) {
	case 1:
		rruleStr = strs[0]
	case 2:
		dtstartStr = strs[0]
		rruleStr = strs[1]
	default:
		return nil, errors.New("invalid RRULE string")
	}

	result := ROption{}
	freqSet := false

	if dtstartStr != "" {
		firstName, err := processRRuleName(dtstartStr)
		if err != nil {
			return nil, fmt.Errorf("expect DTSTART but: %s", err)
		}
		if firstName != "DTSTART" {
			return nil, fmt.Errorf("expect DTSTART but: %s", firstName)
		}

		dtstartValue := dtstartStr[len(firstName)+1:]
		// Check if this is an all-day event (VALUE=DATE)
		if strings.HasPrefix(dtstartValue, "VALUE=DATE:") {
			result.AllDay = true
		}

		result.Dtstart, err = StrToDtStart(dtstartValue, loc)
		if err != nil {
			return nil, fmt.Errorf("StrToDtStart failed: %s", err)
		}
	}

	rruleStr = strings.TrimPrefix(rruleStr, "RRULE:")
	for _, attr := range strings.Split(rruleStr, ";") {
		keyValue := strings.Split(attr, "=")
		if len(keyValue) != 2 {
			return nil, errors.New("wrong format")
		}
		key, value := keyValue[0], keyValue[1]
		if len(value) == 0 {
			return nil, errors.New(key + " option has no value")
		}
		var e error
		switch key {
		case "FREQ":
			result.Freq, e = StrToFreq(value)
			freqSet = true
		case "DTSTART":
			result.Dtstart, e = strToTimeInLoc(value, loc)
		case "INTERVAL":
			result.Interval, e = strconv.Atoi(value)
		case "WKST":
			result.Wkst, e = strToWeekday(value)
		case "COUNT":
			result.Count, e = strconv.Atoi(value)
		case "UNTIL":
			result.Until, e = strToTimeInLoc(value, loc)
		case "BYSETPOS":
			result.Bysetpos, e = strToInts(value)
		case "BYMONTH":
			result.Bymonth, e = strToInts(value)
		case "BYMONTHDAY":
			result.Bymonthday, e = strToInts(value)
		case "BYYEARDAY":
			result.Byyearday, e = strToInts(value)
		case "BYWEEKNO":
			result.Byweekno, e = strToInts(value)
		case "BYDAY":
			result.Byweekday, e = strToWeekdays(value)
		case "BYHOUR":
			result.Byhour, e = strToInts(value)
		case "BYMINUTE":
			result.Byminute, e = strToInts(value)
		case "BYSECOND":
			result.Bysecond, e = strToInts(value)
		case "BYEASTER":
			result.Byeaster, e = strToInts(value)
		default:
			return nil, errors.New("unknown RRULE property: " + key)
		}
		if e != nil {
			return nil, e
		}
	}
	if !freqSet {
		// Per RFC 5545, FREQ is mandatory and supposed to be the first
		// parameter. We'll just confirm it exists because we do not
		// have a meaningful default nor a way to confirm if we parsed
		// a value from the options this returns.
		return nil, errors.New("RRULE property FREQ is required")
	}
	return &result, nil
}
