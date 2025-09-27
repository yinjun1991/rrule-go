// 2017-2022, Teambition. All rights reserved.

package rrule

import (
	"errors"
	"fmt"
	"sort"
	"time"
)

// RRule offers a small, complete, and very fast, implementation of the recurrence rules
// documented in the iCalendar RFC, including support for caching of results.
type RRule struct {
	Options                 ROption
	freq                    Frequency
	dtstart                 time.Time
	interval                int
	wkst                    int
	count                   int       // 初始值为 0 表示不限制 count
	until                   time.Time // zero time 表示不限制 until
	bysetpos                []int
	bymonth                 []int
	bymonthday, bynmonthday []int
	byyearday               []int
	byweekno                []int
	byweekday               []int
	bynweekday              []Weekday
	byhour                  []int
	byminute                []int
	bysecond                []int
	byeaster                []int
	timeset                 []time.Time
	len                     int
}

// NewRRule construct a new RRule instance
func NewRRule(arg ROption) (*RRule, error) {
	if err := validateBounds(arg); err != nil {
		return nil, err
	}
	r := buildRRule(arg)
	return &r, nil
}

func buildRRule(arg ROption) RRule {
	r := RRule{}
	r.Options = arg
	// FREQ default to YEARLY
	r.freq = arg.Freq

	// INTERVAL default to 1
	if arg.Interval < 1 {
		arg.Interval = 1
	}
	r.interval = arg.Interval

	if arg.Count < 0 {
		arg.Count = 0
	}
	r.count = arg.Count

	// DTSTART default to now
	if arg.Dtstart.IsZero() {
		arg.Dtstart = time.Now().UTC()
	}

	// Handle AllDay events: convert to floating time (UTC) as per RFC 5545
	if arg.AllDay {
		// All-day events should use floating time (no timezone binding)
		// In Go, we represent floating time as UTC to ensure consistency
		year, month, day := arg.Dtstart.Date()
		arg.Dtstart = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	} else {
		// Non all-day events: truncate to second precision
		arg.Dtstart = arg.Dtstart.Truncate(time.Second)
	}
	r.dtstart = arg.Dtstart

	// UNTIL
	if arg.Until.IsZero() {
		// add largest representable duration (approximately 290 years).
		r.until = r.dtstart.Add(time.Duration(1<<63 - 1))
	} else {
		// Handle AllDay events: convert to floating time (UTC) as per RFC 5545
		if arg.AllDay {
			// All-day events should use floating time (no timezone binding)
			// In Go, we represent floating time as UTC to ensure consistency
			year, month, day := arg.Until.Date()
			arg.Until = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		} else {
			// Non all-day events: truncate to second precision
			arg.Until = arg.Until.Truncate(time.Second)
		}
		r.until = arg.Until
	}

	r.wkst = arg.Wkst.weekday
	r.bysetpos = arg.Bysetpos

	if len(arg.Byweekno) == 0 &&
		len(arg.Byyearday) == 0 &&
		len(arg.Bymonthday) == 0 &&
		len(arg.Byweekday) == 0 &&
		len(arg.Byeaster) == 0 {
		if r.freq == YEARLY {
			if len(arg.Bymonth) == 0 {
				arg.Bymonth = []int{int(r.dtstart.Month())}
			}
			arg.Bymonthday = []int{r.dtstart.Day()}
		} else if r.freq == MONTHLY {
			arg.Bymonthday = []int{r.dtstart.Day()}
		} else if r.freq == WEEKLY {
			arg.Byweekday = []Weekday{{weekday: toPyWeekday(r.dtstart.Weekday())}}
		}
	}
	r.bymonth = arg.Bymonth
	r.byyearday = arg.Byyearday
	r.byeaster = arg.Byeaster
	for _, mday := range arg.Bymonthday {
		if mday > 0 {
			r.bymonthday = append(r.bymonthday, mday)
		} else if mday < 0 {
			r.bynmonthday = append(r.bynmonthday, mday)
		}
	}
	r.byweekno = arg.Byweekno
	for _, wday := range arg.Byweekday {
		if wday.n == 0 || r.freq > MONTHLY {
			r.byweekday = append(r.byweekday, wday.weekday)
		} else {
			r.bynweekday = append(r.bynweekday, wday)
		}
	}
	if len(arg.Byhour) == 0 {
		if r.freq < HOURLY {
			if arg.AllDay {
				// All-day events should have hour set to 0
				r.byhour = []int{0}
			} else {
				r.byhour = []int{r.dtstart.Hour()}
			}
		}
	} else {
		r.byhour = arg.Byhour
	}
	if len(arg.Byminute) == 0 {
		if r.freq < MINUTELY {
			if arg.AllDay {
				// All-day events should have minute set to 0
				r.byminute = []int{0}
			} else {
				r.byminute = []int{r.dtstart.Minute()}
			}
		}
	} else {
		r.byminute = arg.Byminute
	}
	if len(arg.Bysecond) == 0 {
		if r.freq < SECONDLY {
			if arg.AllDay {
				// All-day events should have second set to 0
				r.bysecond = []int{0}
			} else {
				r.bysecond = []int{r.dtstart.Second()}
			}
		}
	} else {
		r.bysecond = arg.Bysecond
	}

	// Reset the timeset value
	r.timeset = nil

	if r.freq < HOURLY {
		r.timeset = make([]time.Time, 0, len(r.byhour)*len(r.byminute)*len(r.bysecond))
		for _, hour := range r.byhour {
			for _, minute := range r.byminute {
				for _, second := range r.bysecond {
					r.timeset = append(r.timeset, time.Date(1, 1, 1, hour, minute, second, 0, r.dtstart.Location()))
				}
			}
		}
		sort.Sort(timeSlice(r.timeset))
	}

	return r
}

// validateBounds checks the RRule's options are within the boundaries defined
// in RRFC 5545. This is useful to ensure that the RRule can even have any times,
// as going outside these bounds trivially will never have any dates. This can catch
// obvious user error.
func validateBounds(arg ROption) error {
	bounds := []struct {
		field     []int
		param     string
		bound     []int
		plusMinus bool // If the bound also applies for -x to -y.
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

	// Days can optionally specify weeks, like BYDAY=+2MO for the 2nd Monday
	// of the month/year.
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

// Iterator return an iterator for RRule
func (r *RRule) Iterator() Next {
	iterator := rIterator{}
	iterator.year, iterator.month, iterator.day = r.dtstart.Date()

	// Handle AllDay events: use 00:00:00 for all-day events
	if r.Options.AllDay {
		iterator.hour, iterator.minute, iterator.second = 0, 0, 0
	} else {
		iterator.hour, iterator.minute, iterator.second = r.dtstart.Clock()
	}
	iterator.weekday = toPyWeekday(r.dtstart.Weekday())

	iterator.ii = iterInfo{rrule: r}
	iterator.ii.rebuild(iterator.year, iterator.month)

	if r.freq < HOURLY {
		iterator.timeset = r.timeset
	} else {
		if r.freq >= HOURLY && len(r.byhour) != 0 && !contains(r.byhour, iterator.hour) ||
			r.freq >= MINUTELY && len(r.byminute) != 0 && !contains(r.byminute, iterator.minute) ||
			r.freq >= SECONDLY && len(r.bysecond) != 0 && !contains(r.bysecond, iterator.second) {
			iterator.timeset = nil
		} else {
			iterator.ii.fillTimeSet(&iterator.timeset, r.freq, iterator.hour, iterator.minute, iterator.second)
		}
	}
	iterator.count = r.count
	return iterator.next
}

// All returns all occurrences of the RRule.
// It is only supported second precision.
func (r *RRule) All() []time.Time {
	return all(r.Iterator())
}

// Between returns all the occurrences of the RRule between after and before.
// The inc keyword defines what happens if after and/or before are themselves occurrences.
// With inc == True, they will be included in the list, if they are found in the recurrence set.
// It is only supported second precision.
func (r *RRule) Between(after, before time.Time, inc bool) []time.Time {
	return between(r.Iterator(), after, before, inc)
}

// Before returns the last recurrence before the given datetime instance,
// or time.Time's zero value if no recurrence match.
// The inc keyword defines what happens if dt is an occurrence.
// With inc == True, if dt itself is an occurrence, it will be returned.
// It is only supported second precision.
func (r *RRule) Before(dt time.Time, inc bool) time.Time {
	return before(r.Iterator(), dt, inc)
}

// After returns the first recurrence after the given datetime instance,
// or time.Time's zero value if no recurrence match.
// The inc keyword defines what happens if dt is an occurrence.
// With inc == True, if dt itself is an occurrence, it will be returned.
// It is only supported second precision.
func (r *RRule) After(dt time.Time, inc bool) time.Time {
	return after(r.Iterator(), dt, inc)
}

// DTStart set a new DTSTART for the rule and recalculates the timeset if needed.
// It will be truncated to second precision.
// Default to `time.Now().UTC().Truncate(time.Second)`.
func (r *RRule) DTStart(dt time.Time) {
	// Handle AllDay events: convert to floating time (UTC) as per RFC 5545
	if r.Options.AllDay {
		// All-day events should use floating time (no timezone binding)
		// In Go, we represent floating time as UTC to ensure consistency
		year, month, day := dt.Date()
		r.Options.Dtstart = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	} else {
		// Non all-day events: truncate to second precision
		r.Options.Dtstart = dt.Truncate(time.Second)
	}
	*r = buildRRule(r.Options)
}

// GetDTStart gets DTSTART time for rrule
func (r *RRule) GetDTStart() time.Time {
	return r.dtstart
}

// Until set a new UNTIL for the rule and recalculates the timeset if needed.
// It will be truncated to second precision.
// Default to `Dtstart.Add(time.Duration(1<<63 - 1))`, approximately 290 years.
func (r *RRule) Until(ut time.Time) {
	// Handle AllDay events: convert to floating time (UTC) as per RFC 5545
	if r.Options.AllDay {
		// All-day events should use floating time (no timezone binding)
		// In Go, we represent floating time as UTC to ensure consistency
		year, month, day := ut.Date()
		r.Options.Until = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	} else {
		// Non all-day events: truncate to second precision
		r.Options.Until = ut.Truncate(time.Second)
	}
	*r = buildRRule(r.Options)
}

// GetUntil gets UNTIL time for rrule
func (r *RRule) GetUntil() time.Time {
	return r.until
}

// IsAllDay returns whether the set is configured for all-day events.
func (r *RRule) IsAllDay() bool {
	return r.Options.AllDay
}

func (r *RRule) SetAllDay(allDay bool) {
	r.Options.AllDay = allDay

	// If switching to all-day, normalize existing time values
	if allDay {
		// Normalize Dtstart to floating time (00:00:00 UTC)
		if !r.Options.Dtstart.IsZero() {
			year, month, day := r.Options.Dtstart.Date()
			r.Options.Dtstart = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		}

		// Normalize Until to floating time (00:00:00 UTC)
		if !r.Options.Until.IsZero() {
			year, month, day := r.Options.Until.Date()
			r.Options.Until = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		}
	}

	// Rebuild the RRule with updated options
	*r = buildRRule(r.Options)
}

func (r *RRule) String() string {
	return r.Options.String()
}

// StrToRRule converts string to RRule
func StrToRRule(rfcString string) (*RRule, error) {
	option, e := StrToROption(rfcString)
	if e != nil {
		return nil, e
	}
	return NewRRule(*option)
}
