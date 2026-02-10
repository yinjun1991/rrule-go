// 2017-2022, Teambition. All rights reserved.

package rrule

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

// Recurrence allows more complex recurrence setups, mixing multiple rules, dates, exclusion rules, and exclusion dates
type Recurrence struct {
	Options                 ROption
	freq                    Frequency
	dtstart                 time.Time
	interval                int
	wkst                    int
	count                   int
	until                   time.Time
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
	rdate                   []time.Time
	exdate                  []time.Time
	allDay                  bool
	hasRule                 bool
}

func New(option ROption) *Recurrence {
	rec, err := newRecurrence(option)
	if err != nil {
		return nil
	}
	return rec
}

func Parse(lines []string) *Recurrence {
	set, err := StrSliceToRRuleSet(lines)
	if err != nil {
		return nil
	}
	return set
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

	r.freq = arg.Freq

	if arg.Interval < 1 {
		arg.Interval = 1
	}
	r.interval = arg.Interval

	if arg.Count < 0 {
		arg.Count = 0
	}
	r.count = arg.Count

	if arg.Dtstart.IsZero() {
		arg.Dtstart = time.Now().UTC()
	}

	if arg.AllDay {
		year, month, day := arg.Dtstart.Date()
		arg.Dtstart = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	} else {
		arg.Dtstart = arg.Dtstart.Truncate(time.Second)
	}
	r.dtstart = arg.Dtstart

	if arg.Until.IsZero() {
		r.until = r.dtstart.Add(time.Duration(1<<63 - 1))
	} else {
		if arg.AllDay {
			year, month, day := arg.Until.Date()
			arg.Until = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		} else {
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
				r.bysecond = []int{0}
			} else {
				r.bysecond = []int{r.dtstart.Second()}
			}
		}
	} else {
		r.bysecond = arg.Bysecond
	}

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

func (r *Recurrence) ruleIterator() Next {
	if !r.hasRule {
		return func() (time.Time, bool) {
			return time.Time{}, false
		}
	}
	iterator := rIterator{}
	iterator.year, iterator.month, iterator.day = r.dtstart.Date()

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

// Recurrence returns a slice of all the recurrence rules for a set
func (set *Recurrence) Recurrence(includeDTSTART bool) []string {
	var res []string

	if !set.dtstart.IsZero() && includeDTSTART {
		// No colon, DTSTART may have TZID, which would require a semicolon after DTSTART
		// RFC 5545: For all-day events, use VALUE=DATE format
		if set.allDay {
			// All-day events should use VALUE=DATE format as per RFC 5545
			year, month, day := set.dtstart.Date()
			dateStr := fmt.Sprintf("%04d%02d%02d", year, int(month), day)
			res = append(res, fmt.Sprintf("DTSTART;VALUE=DATE:%s", dateStr))
		} else {
			res = append(res, fmt.Sprintf("DTSTART%s", timeToRFCDatetimeStr(set.dtstart)))
		}
	}

	if set.hasRule {
		res = append(res, fmt.Sprintf("RRULE:%s", set.Options.RRuleString()))
	}

	for _, item := range set.rdate {
		// RFC 5545: RDATE values should match DTSTART value type
		if set.allDay {
			// All-day events: use VALUE=DATE format as per RFC 5545
			year, month, day := item.Date()
			dateStr := fmt.Sprintf("%04d%02d%02d", year, int(month), day)
			res = append(res, fmt.Sprintf("RDATE;VALUE=DATE:%s", dateStr))
		} else {
			res = append(res, fmt.Sprintf("RDATE%s", timeToRFCDatetimeStr(item)))
		}
	}

	for _, item := range set.exdate {
		// RFC 5545: EXDATE values should match DTSTART value type
		if set.allDay {
			// All-day events: use VALUE=DATE format as per RFC 5545
			year, month, day := item.Date()
			dateStr := fmt.Sprintf("%04d%02d%02d", year, int(month), day)
			res = append(res, fmt.Sprintf("EXDATE;VALUE=DATE:%s", dateStr))
		} else {
			res = append(res, fmt.Sprintf("EXDATE%s", timeToRFCDatetimeStr(item)))
		}
	}
	return res
}

func (set *Recurrence) String(includeDTSTART bool) string {
	res := set.Recurrence(includeDTSTART)
	return strings.Join(res, "\n")
}

// DTStart sets dtstart property for set.
// It will be truncated to second precision.
func (set *Recurrence) DTStart(dtstart time.Time) {
	// Handle AllDay events: convert to floating time (UTC) as per RFC 5545
	if set.allDay {
		// All-day events should use floating time (no timezone binding)
		// In Go, we represent floating time as UTC to ensure consistency
		year, month, day := dtstart.Date()
		set.dtstart = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	} else {
		// Non all-day events: truncate to second precision
		set.dtstart = dtstart.Truncate(time.Second)
	}

	set.Options.Dtstart = set.dtstart
	if set.hasRule {
		set.rebuildRule()
	}
}

// GetDTStart gets DTSTART for set
func (set *Recurrence) GetDTStart() time.Time {
	return set.dtstart
}

func (set *Recurrence) GetOptions() *ROption {
	if !set.hasRule {
		return nil
	}
	return &set.Options
}

// RDate include the given datetime instance in the recurrence set generation.
// It will be truncated to second precision.
func (set *Recurrence) RDate(rdate time.Time) {
	// Handle AllDay events: convert to floating time (UTC) as per RFC 5545
	if set.allDay {
		// All-day events should use floating time (no timezone binding)
		// In Go, we represent floating time as UTC to ensure consistency
		year, month, day := rdate.Date()
		set.rdate = append(set.rdate, time.Date(year, month, day, 0, 0, 0, 0, time.UTC))
	} else {
		// Non all-day events: truncate to second precision
		set.rdate = append(set.rdate, rdate.Truncate(time.Second))
	}
}

// SetRDates sets explicitly added dates (rdates) in the set.
// It will be truncated to second precision.
func (set *Recurrence) SetRDates(rdates []time.Time) {
	set.rdate = make([]time.Time, 0, len(rdates))
	for _, rdate := range rdates {
		// Handle AllDay events: convert to floating time (UTC) as per RFC 5545
		if set.allDay {
			// All-day events should use floating time (no timezone binding)
			// In Go, we represent floating time as UTC to ensure consistency
			year, month, day := rdate.Date()
			set.rdate = append(set.rdate, time.Date(year, month, day, 0, 0, 0, 0, time.UTC))
		} else {
			// Non all-day events: truncate to second precision
			set.rdate = append(set.rdate, rdate.Truncate(time.Second))
		}
	}
}

// GetRDate returns explicitly added dates (rdates) in the set
func (set *Recurrence) GetRDate() []time.Time {
	return set.rdate
}

// ExDate include the given datetime instance in the recurrence set exclusion list.
// Dates included that way will not be generated,
// even if some inclusive rrule or rdate matches them.
// It will be truncated to second precision.
func (set *Recurrence) ExDate(exdate time.Time) {
	// Handle AllDay events: convert to floating time (UTC) as per RFC 5545
	if set.allDay {
		// All-day events should use floating time (no timezone binding)
		// In Go, we represent floating time as UTC to ensure consistency
		year, month, day := exdate.Date()
		set.exdate = append(set.exdate, time.Date(year, month, day, 0, 0, 0, 0, time.UTC))
	} else {
		// Non all-day events: truncate to second precision
		set.exdate = append(set.exdate, exdate.Truncate(time.Second))
	}
}

// SetExDates sets explicitly excluded dates (exdates) in the set.
// It will be truncated to second precision.
func (set *Recurrence) SetExDates(exdates []time.Time) {
	set.exdate = make([]time.Time, 0, len(exdates))
	for _, exdate := range exdates {
		// Handle AllDay events: convert to floating time (UTC) as per RFC 5545
		if set.allDay {
			// All-day events should use floating time (no timezone binding)
			// In Go, we represent floating time as UTC to ensure consistency
			year, month, day := exdate.Date()
			set.exdate = append(set.exdate, time.Date(year, month, day, 0, 0, 0, 0, time.UTC))
		} else {
			// Non all-day events: truncate to second precision
			set.exdate = append(set.exdate, exdate.Truncate(time.Second))
		}
	}
}

// GetExDate returns explicitly excluded dates (exdates) in the set.
func (set *Recurrence) GetExDate() []time.Time {
	return set.exdate
}

// SetAllDay sets the all-day flag for the set.
// When set to true, all time values (dtstart, rdate, exdate) will be normalized to floating time.
func (set *Recurrence) SetAllDay(allDay bool) {
	set.allDay = allDay

	// If switching to all-day, normalize existing times
	if allDay {
		// Normalize dtstart
		if !set.dtstart.IsZero() {
			year, month, day := set.dtstart.Date()
			set.dtstart = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		}

		// Normalize rdate
		for i, rdate := range set.rdate {
			year, month, day := rdate.Date()
			set.rdate[i] = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		}

		// Normalize exdate
		for i, exdate := range set.exdate {
			year, month, day := exdate.Date()
			set.exdate[i] = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		}

		if !set.dtstart.IsZero() {
			set.Options.Dtstart = set.dtstart
		}
	}

	if set.hasRule {
		set.Options.AllDay = allDay
		set.rebuildRule()
	}
}

// IsAllDay returns whether the set is configured for all-day events.
func (set *Recurrence) IsAllDay() bool {
	return set.allDay
}

type genItem struct {
	dt  time.Time
	gen Next
}

type genItemSlice []genItem

func (s genItemSlice) Len() int           { return len(s) }
func (s genItemSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s genItemSlice) Less(i, j int) bool { return s[i].dt.Before(s[j].dt) }

func addGenList(genList *[]genItem, next Next) {
	dt, ok := next()
	if ok {
		*genList = append(*genList, genItem{dt, next})
	}
}

// Iterator returns an iterator for Recurrence
func (set *Recurrence) Iterator() (next func() (time.Time, bool)) {
	rlist := []genItem{}
	exlist := []genItem{}

	// Normalize rdate times for all-day events
	rdates := set.rdate
	if set.allDay && len(set.rdate) > 0 {
		rdates = make([]time.Time, len(set.rdate))
		for i, t := range set.rdate {
			// Convert to floating time (00:00:00 UTC represents floating time)
			rdates[i] = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
		}
	}

	sort.Sort(timeSlice(rdates))
	addGenList(&rlist, timeSliceIterator(rdates))
	if set.hasRule {
		addGenList(&rlist, set.ruleIterator())
	}
	sort.Sort(genItemSlice(rlist))

	// Normalize exdate times for all-day events
	exdates := set.exdate
	if set.allDay && len(set.exdate) > 0 {
		exdates = make([]time.Time, len(set.exdate))
		for i, t := range set.exdate {
			// Convert to floating time (00:00:00 UTC represents floating time)
			exdates[i] = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
		}
	}

	sort.Sort(timeSlice(exdates))
	addGenList(&exlist, timeSliceIterator(exdates))
	sort.Sort(genItemSlice(exlist))

	lastdt := time.Time{}
	return func() (time.Time, bool) {
		for len(rlist) != 0 {
			dt := rlist[0].dt
			var ok bool
			rlist[0].dt, ok = rlist[0].gen()
			if !ok {
				rlist = rlist[1:]
			}
			sort.Sort(genItemSlice(rlist))

			// Normalize dt for all-day events to ensure consistent comparison
			normalizedDt := dt
			if set.allDay {
				normalizedDt = time.Date(dt.Year(), dt.Month(), dt.Day(), 0, 0, 0, 0, time.UTC)
			}

			// Use normalized time for comparison
			normalizedLastdt := lastdt
			if set.allDay && !lastdt.IsZero() {
				normalizedLastdt = time.Date(lastdt.Year(), lastdt.Month(), lastdt.Day(), 0, 0, 0, 0, time.UTC)
			}

			if lastdt.IsZero() || !normalizedLastdt.Equal(normalizedDt) {
				for len(exlist) != 0 && exlist[0].dt.Before(normalizedDt) {
					exlist[0].dt, ok = exlist[0].gen()
					if !ok {
						exlist = exlist[1:]
					}
					sort.Sort(genItemSlice(exlist))
				}
				lastdt = normalizedDt
				if len(exlist) == 0 || !normalizedDt.Equal(exlist[0].dt) {
					return normalizedDt, true
				}
			}
		}
		return time.Time{}, false
	}
}

// All returns all occurrences of the Recurrence.
// It is only supported second precision.
func (set *Recurrence) All() []time.Time {
	return all(set.Iterator())
}

// Between returns all the occurrences of the rrule between after and before.
// The inc keyword defines what happens if after and/or before are themselves occurrences.
// With inc == True, they will be included in the list, if they are found in the recurrence set.
// It is only supported second precision.
func (set *Recurrence) Between(after, before time.Time, inc bool) []time.Time {
	return between(set.Iterator(), after, before, inc)
}

// Before Returns the last recurrence before the given datetime instance,
// or time.Time's zero value if no recurrence match.
// The inc keyword defines what happens if dt is an occurrence.
// With inc == True, if dt itself is an occurrence, it will be returned.
// It is only supported second precision.
func (set *Recurrence) Before(dt time.Time, inc bool) time.Time {
	return before(set.Iterator(), dt, inc)
}

// After returns the first recurrence after the given datetime instance,
// or time.Time's zero value if no recurrence match.
// The inc keyword defines what happens if dt is an occurrence.
// With inc == True, if dt itself is an occurrence, it will be returned.
// It is only supported second precision.
func (set *Recurrence) After(dt time.Time, inc bool) time.Time {
	return after(set.Iterator(), dt, inc)
}

// StrToRRuleSet converts string to RRuleSet
func StrToRRuleSet(s string) (*Recurrence, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, errors.New("empty string")
	}
	ss := strings.Split(s, "\n")
	return StrSliceToRRuleSet(ss)
}

// StrSliceToRRuleSet converts given str slice to RRuleSet
// In case there is a time met in any rule without specified time zone, when
// it is parsed in UTC (see StrSliceToRRuleSetInLoc)
func StrSliceToRRuleSet(ss []string) (*Recurrence, error) {
	return StrSliceToRRuleSetInLoc(ss, time.UTC)
}

// StrSliceToRRuleSetInLoc is same as StrSliceToRRuleSet, but by default parses local times
// in specified default location
func StrSliceToRRuleSetInLoc(ss []string, defaultLoc *time.Location) (*Recurrence, error) {
	if len(ss) == 0 {
		return &Recurrence{}, nil
	}

	set := Recurrence{}
	var dtstartLineForRRULE string

	// According to RFC DTSTART is always the first line.
	firstName, err := processRRuleName(ss[0])
	if err != nil {
		return nil, err
	}
	if firstName == "DTSTART" {
		// Detect all-day (VALUE=DATE) before parsing to ensure normalization
		dtstartField := ss[0][len(firstName)+1:]
		if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(dtstartField)), "VALUE=DATE:") {
			// Set as all-day before applying DTStart so normalization uses date semantics
			set.SetAllDay(true)
		}

		dt, err := StrToDtStart(dtstartField, defaultLoc)
		if err != nil {
			return nil, fmt.Errorf("StrToDtStart failed: %v", err)
		}
		// default location should be taken from DTSTART property to correctly
		// parse local times met in RDATE,EXDATE and other rules
		defaultLoc = dt.Location()
		set.DTStart(dt)
		if !set.GetDTStart().IsZero() {
			if set.allDay {
				dtstartLineForRRULE = fmt.Sprintf("DTSTART;VALUE=DATE:%s", set.GetDTStart().Format(DateFormat))
			} else {
				dtstartLineForRRULE = fmt.Sprintf("DTSTART%s", timeToRFCDatetimeStr(set.GetDTStart()))
			}
		}
		// We've processed the first one
		ss = ss[1:]
	}

	for _, line := range ss {
		name, err := processRRuleName(line)
		if err != nil {
			return nil, err
		}
		rule := line[len(name)+1:]

		switch name {
		case "RRULE":
			rruleInput := line
			if dtstartLineForRRULE != "" {
				rruleInput = dtstartLineForRRULE + "\n" + line
			}
			rOpt, err := StrToROption(rruleInput)
			if err != nil {
				return nil, fmt.Errorf("StrToROption failed: %v", err)
			}
			err = set.setRuleOptions(*rOpt)
			if err != nil {
				return nil, fmt.Errorf("NewRRule failed: %v", err)
			}
		case "RDATE", "EXDATE":
			if !set.allDay && containsValueDateParam(rule) {
				set.SetAllDay(true)
			}

			ts, err := StrToDatesInLoc(rule, defaultLoc)
			if err != nil {
				return nil, fmt.Errorf("strToDates failed: %v", err)
			}
			for _, t := range ts {
				if name == "RDATE" {
					set.RDate(t)
				} else {
					set.ExDate(t)
				}
			}
		}
	}

	return &set, nil
}

func containsValueDateParam(rule string) bool {
	upper := strings.ToUpper(rule)
	paramSection := upper
	if idx := strings.Index(paramSection, ":"); idx != -1 {
		paramSection = paramSection[:idx]
	}
	for _, part := range strings.Split(paramSection, ";") {
		if strings.TrimSpace(part) == "VALUE=DATE" {
			return true
		}
	}
	return false
}

func normalizeRecurrenceStrings(ruleset []string) ([]string, error) {
	if len(ruleset) == 0 {
		return nil, fmt.Errorf("empty input strings")
	}

	normalized, err := NormalizeRecurrenceRuleset(ruleset)
	if err != nil {
		return nil, err
	}

	var result []string
	var foundRRule bool

	for _, str := range normalized {
		str = strings.TrimSpace(str)
		if str == "" {
			continue
		}

		if strings.HasPrefix(strings.ToUpper(str), "RRULE:") {
			if !foundRRule {
				result = append(result, str)
				foundRRule = true
			}
			continue
		}

		result = append(result, str)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no valid recurrence strings found")
	}

	return result, nil
}
