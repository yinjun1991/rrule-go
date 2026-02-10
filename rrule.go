// 2017-2022, Teambition. All rights reserved.

package rrule

import (
	"errors"
	"fmt"
	"sort"
	"time"
)

func New(option ROption) *Recurrence {
	rec, err := newRecurrence(option)
	if err != nil {
		return nil
	}
	return rec
}

func newRecurrence(option ROption) (*Recurrence, error) {
	rec := &Recurrence{}
	if err := rec.setRuleOptions(option); err != nil {
		return nil, err
	}
	return rec, nil
}

func (r *Recurrence) normalizeAllDayTimes() {
	if !r.dtstart.IsZero() {
		year, month, day := r.dtstart.Date()
		r.dtstart = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	}
	for i, rdate := range r.rdate {
		year, month, day := rdate.Date()
		r.rdate[i] = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	}
	for i, exdate := range r.exdate {
		year, month, day := exdate.Date()
		r.exdate[i] = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	}
}

func (r *Recurrence) setRuleOptions(option ROption) error {
	if option.AllDay && !r.allDay {
		r.allDay = true
		r.normalizeAllDayTimes()
	} else if r.allDay && !option.AllDay {
		option.AllDay = true
	}

	if option.Dtstart.IsZero() && !r.dtstart.IsZero() {
		option.Dtstart = r.dtstart
	}

	if err := r.applyRule(option); err != nil {
		return err
	}

	if !r.dtstart.IsZero() {
		r.Options.Dtstart = r.dtstart
	}

	r.hasRule = true
	return nil
}

func (r *Recurrence) rebuildRule() {
	_ = r.applyRule(r.Options)
	r.hasRule = true
}

func (r *Recurrence) applyRule(arg ROption) error {
	if err := validateBounds(arg); err != nil {
		return err
	}
	r.Options = arg
	r.allDay = arg.AllDay

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

	return nil
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

func (r *Recurrence) ruleIterator() Next {
	if !r.hasRule {
		return func() (time.Time, bool) {
			return time.Time{}, false
		}
	}
	iterator := rIterator{}
	iterator.year, iterator.month, iterator.day = r.dtstart.Date()

	// Handle AllDay events: use 00:00:00 for all-day events
	if r.allDay {
		iterator.hour, iterator.minute, iterator.second = 0, 0, 0
	} else {
		iterator.hour, iterator.minute, iterator.second = r.dtstart.Clock()
	}
	iterator.weekday = toPyWeekday(r.dtstart.Weekday())

	iterator.ii = iterInfo{recurrence: r}
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
