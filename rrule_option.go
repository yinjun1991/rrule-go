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
// If Dtstart is zero, it defaults to time.Now().UTC().
// AllDay uses floating DATE semantics (VALUE=DATE); non-all-day does not support floating DATE-TIME.
type ROption struct {
	Freq       Frequency
	Dtstart    time.Time // Caller must set the timezone on Dtstart first; if AllDay is true, recurrence starts from Dtstart's local date (e.g., 2024-06-01T23:00:00+02:00 starts on 2024-06-01).
	Interval   int
	Wkst       Weekday
	Count      int
	Until      time.Time // For all-day, only the date fields of Until are used (time-of-day is ignored); for non-all-day, Until must be UTC and preserves time-of-day.
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
			loc := time.UTC
			if !option.Dtstart.IsZero() {
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

	result := ROption{}
	var dtstartIsDate bool
	var dtstartHasTZID bool
	var dtstartIsUTC bool
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
		dtstartIsDate, dtstartHasTZID, dtstartIsUTC = detectDtstartKind(dtstartValue)
		// Check if this is an all-day event (VALUE=DATE)
		if dtstartIsDate {
			result.AllDay = true
		}

		result.Dtstart, err = StrToDtStart(dtstartValue, defaultLoc)
		if err != nil {
			return nil, fmt.Errorf("StrToDtStart failed: %s", err)
		}
		if !result.Dtstart.IsZero() {
			defaultLoc = result.Dtstart.Location()
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
			if dtstartIsDate {
				if len(value) != len(DateFormat) || strings.Contains(value, "T") {
					return nil, fmt.Errorf("UNTIL must be DATE when DTSTART is DATE")
				}
			} else if dtstartHasTZID || dtstartIsUTC {
				if !strings.HasSuffix(strings.ToUpper(value), "Z") {
					return nil, fmt.Errorf("UNTIL must be UTC when DTSTART uses TZID or UTC")
				}
			}
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

	if !freqSet {
		// Per RFC 5545, FREQ is mandatory and supposed to be the first
		// parameter. We'll just confirm it exists because we do not
		// have a meaningful default nor a way to confirm if we parsed
		// a value from the options this returns.
		return nil, errors.New("RRULE property FREQ is required")
	}
	return &result, nil
}

func detectDtstartKind(dtstartValue string) (bool, bool, bool) {
	upper := strings.ToUpper(dtstartValue)
	if strings.HasPrefix(upper, "VALUE=DATE:") {
		return true, false, false
	}
	hasTZID := strings.HasPrefix(upper, "TZID=") || strings.Contains(upper, ";TZID=")
	timePart := upper
	if idx := strings.LastIndex(upper, ":"); idx >= 0 {
		timePart = upper[idx+1:]
	}
	isUTC := strings.HasSuffix(timePart, "Z")
	return false, hasTZID, isUTC
}
