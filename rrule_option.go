package rrule

import (
	"errors"
	"fmt"
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
	RDate      []time.Time
	EXDate     []time.Time
	AllDay     bool
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
