package rrule

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ROption offers options to construct a RRule instance.
// For performance, it is strongly recommended providing explicit ROption.Dtstart.
// If Dtstart is zero, it defaults to time.Now() in the effective location (Dtstart's location, otherwise Location or UTC).
type ROption struct {
	Freq       Frequency
	Dtstart    time.Time
	Location   *time.Location // Default timezone for local date-times; falls back to UTC.
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

	dtstart := option.Dtstart
	if option.Location != nil {
		dtstart = dtstart.In(option.Location)
	}

	// Handle AllDay events: use DATE format as per RFC 5545
	if option.AllDay {
		// All-day events should use VALUE=DATE format
		dateStr := dtstart.Format(DateFormat)
		return fmt.Sprintf("DTSTART;VALUE=DATE:%s\nRRULE:%s", dateStr, str)
	}

	// Non-all-day events use DATE-TIME format
	return fmt.Sprintf("DTSTART%s\nRRULE:%s", timeToRFCDatetimeStr(dtstart), str)
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
			loc := option.Location
			if loc == nil {
				loc = option.Dtstart.Location()
			}
			until := option.Until.In(loc)
			// For all-day events, use DATE format (no time part) as per RFC 5545
			result = append(result, fmt.Sprintf("UNTIL=%v", until.Format(DateFormat)))
		} else {
			// For date-time events, UNTIL is represented in UTC
			result = append(result, fmt.Sprintf("UNTIL=%v", timeToUTCStr(option.Until)))
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
	defaultLoc := time.UTC
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

	result := ROption{Location: defaultLoc}
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

		result.Dtstart, err = StrToDtStart(dtstartValue, defaultLoc)
		if err != nil {
			return nil, fmt.Errorf("StrToDtStart failed: %s", err)
		}
		if !result.Dtstart.IsZero() {
			result.Location = result.Dtstart.Location()
			defaultLoc = result.Location
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
		var err error
		switch key {
		case "FREQ":
			result.Freq, err = StrToFreq(value)
			freqSet = true
		case "DTSTART":
			result.Dtstart, err = strToTimeInLoc(value, defaultLoc)
		case "INTERVAL":
			result.Interval, err = strconv.Atoi(value)
		case "WKST":
			result.Wkst, err = strToWeekday(value)
		case "COUNT":
			result.Count, err = strconv.Atoi(value)
		case "UNTIL":
			result.Until, err = strToTimeInLoc(value, defaultLoc)
		case "BYSETPOS":
			result.Bysetpos, err = strToInts(value)
		case "BYMONTH":
			result.Bymonth, err = strToInts(value)
		case "BYMONTHDAY":
			result.Bymonthday, err = strToInts(value)
		case "BYYEARDAY":
			result.Byyearday, err = strToInts(value)
		case "BYWEEKNO":
			result.Byweekno, err = strToInts(value)
		case "BYDAY":
			result.Byweekday, err = strToWeekdays(value)
		case "BYHOUR":
			result.Byhour, err = strToInts(value)
		case "BYMINUTE":
			result.Byminute, err = strToInts(value)
		case "BYSECOND":
			result.Bysecond, err = strToInts(value)
		case "BYEASTER":
			result.Byeaster, err = strToInts(value)
		default:
			return nil, errors.New("unknown RRULE property: " + key)
		}
		if err != nil {
			return nil, err
		}
	}

	if !result.Dtstart.IsZero() {
		result.Location = result.Dtstart.Location()
	} else if result.Location == nil {
		result.Location = time.UTC
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

func (option *ROption) GetLocation() *time.Location {
	if option.Location != nil {
		return option.Location
	}

	if option.Dtstart.IsZero() {
		return time.UTC
	}

	return option.Dtstart.Location()
}
