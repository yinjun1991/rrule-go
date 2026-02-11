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
	str := rruleStringFromOption(option)
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

func rruleStringFromOption(option *ROption) string {
	result := []string{fmt.Sprintf("FREQ=%v", option.Freq)}
	if option.Interval != 0 {
		result = append(result, fmt.Sprintf("INTERVAL=%v", option.Interval))
	}
	if option.Wkst != MO {
		result = append(result, fmt.Sprintf("WKST=%v", option.Wkst))
	}
	if option.Count > 0 {
		result = append(result, fmt.Sprintf("COUNT=%v", option.Count))
	}
	if !option.Until.IsZero() {
		if option.AllDay {
			loc := time.UTC
			if !option.Dtstart.IsZero() {
				loc = option.Dtstart.Location()
			}
			until := option.Until.In(loc)
			result = append(result, fmt.Sprintf("UNTIL=%v", until.Format(DateFormat)))
		} else {
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

func validateBounds(arg ROption) error {
	bounds := []struct {
		field     []int
		param     string
		bound     []int
		plusMinus bool
	}{
		{arg.Bysecond, "bysecond", []int{0, 59}, false},
		{arg.Byminute, "byminute", []int{0, 59}, false},
		{arg.Byhour, "byhour", []int{0, 23}, false},
		{arg.Bymonthday, "bymonthday", []int{1, 31}, true},
		{arg.Byyearday, "byyearday", []int{1, 366}, true},
		{arg.Byweekno, "byweekno", []int{1, 53}, true},
		{arg.Bymonth, "bymonth", []int{1, 12}, false},
		{arg.Bysetpos, "bysetpos", []int{1, 366}, true},
	}

	checkBounds := func(param string, value int, bounds []int, plusMinus bool) error {
		if !(value >= bounds[0] && value <= bounds[1]) && (!plusMinus || !(value <= -bounds[0] && value >= -bounds[1])) {
			plusMinusBounds := ""
			if plusMinus {
				plusMinusBounds = fmt.Sprintf(" or %d and %d", -bounds[0], -bounds[1])
			}
			return fmt.Errorf("%s must be between %d and %d%s", param, bounds[0], bounds[1], plusMinusBounds)
		}
		return nil
	}

	for _, b := range bounds {
		for _, value := range b.field {
			if err := checkBounds(b.param, value, b.bound, b.plusMinus); err != nil {
				return err
			}
		}
	}

	for _, w := range arg.Byweekday {
		if w.n > 53 || w.n < -53 {
			return errors.New("byday must be between 1 and 53 or -1 and -53")
		}
	}

	if arg.Interval < 0 {
		return errors.New("interval must be greater than 0")
	}

	return nil
}
