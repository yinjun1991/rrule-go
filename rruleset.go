// 2017-2022, Teambition. All rights reserved.

package rrule

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

// Set allows more complex recurrence setups, mixing multiple rules, dates, exclusion rules, and exclusion dates
type Set struct {
	dtstart time.Time
	rrule   *RRule
	rdate   []time.Time
	exdate  []time.Time
	allDay  bool // RFC 5545: All-day events use floating time (no timezone binding)
}

// Recurrence returns a slice of all the recurrence rules for a set
func (set *Set) Recurrence(includeDTSTART bool) []string {
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

	if set.rrule != nil {
		res = append(res, fmt.Sprintf("RRULE:%s", set.rrule.Options.RRuleString()))
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

func (set *Set) String(includeDTSTART bool) string {
	res := set.Recurrence(includeDTSTART)
	return strings.Join(res, "\n")
}

// DTStart sets dtstart property for set.
// It will be truncated to second precision.
func (set *Set) DTStart(dtstart time.Time) {
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

	if set.rrule != nil {
		set.rrule.DTStart(set.dtstart)
	}
}

// GetDTStart gets DTSTART for set
func (set *Set) GetDTStart() time.Time {
	return set.dtstart
}

// RRule set the RRULE for set.
// There is the only one RRULE in the set as https://tools.ietf.org/html/rfc5545#appendix-A.1
func (set *Set) RRule(rrule *RRule) {
	if !rrule.Options.Dtstart.IsZero() {
		set.dtstart = rrule.dtstart
	} else if !set.dtstart.IsZero() {
		rrule.DTStart(set.dtstart)
	}
	set.rrule = rrule
}

// GetRRule returns the rrules in the set
func (set *Set) GetRRule() *RRule {
	return set.rrule
}

// RDate include the given datetime instance in the recurrence set generation.
// It will be truncated to second precision.
func (set *Set) RDate(rdate time.Time) {
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
func (set *Set) SetRDates(rdates []time.Time) {
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
func (set *Set) GetRDate() []time.Time {
	return set.rdate
}

// ExDate include the given datetime instance in the recurrence set exclusion list.
// Dates included that way will not be generated,
// even if some inclusive rrule or rdate matches them.
// It will be truncated to second precision.
func (set *Set) ExDate(exdate time.Time) {
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
func (set *Set) SetExDates(exdates []time.Time) {
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
func (set *Set) GetExDate() []time.Time {
	return set.exdate
}

// SetAllDay sets the all-day flag for the set.
// When set to true, all time values (dtstart, rdate, exdate) will be normalized to floating time.
func (set *Set) SetAllDay(allDay bool) {
	set.allDay = allDay

	// Only call SetAllDay on rrule if it exists
	if set.rrule != nil {
		set.rrule.SetAllDay(allDay)
	}

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

		// Update rrule if exists
		if set.rrule != nil {
			set.rrule.DTStart(set.dtstart)
		}
	}
}

// IsAllDay returns whether the set is configured for all-day events.
func (set *Set) IsAllDay() bool {
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

// Iterator returns an iterator for rrule.Set
func (set *Set) Iterator() (next func() (time.Time, bool)) {
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
	if set.rrule != nil {
		addGenList(&rlist, set.rrule.Iterator())
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

// All returns all occurrences of the rrule.Set.
// It is only supported second precision.
func (set *Set) All() []time.Time {
	return all(set.Iterator())
}

// Between returns all the occurrences of the rrule between after and before.
// The inc keyword defines what happens if after and/or before are themselves occurrences.
// With inc == True, they will be included in the list, if they are found in the recurrence set.
// It is only supported second precision.
func (set *Set) Between(after, before time.Time, inc bool) []time.Time {
	return between(set.Iterator(), after, before, inc)
}

// Before Returns the last recurrence before the given datetime instance,
// or time.Time's zero value if no recurrence match.
// The inc keyword defines what happens if dt is an occurrence.
// With inc == True, if dt itself is an occurrence, it will be returned.
// It is only supported second precision.
func (set *Set) Before(dt time.Time, inc bool) time.Time {
	return before(set.Iterator(), dt, inc)
}

// After returns the first recurrence after the given datetime instance,
// or time.Time's zero value if no recurrence match.
// The inc keyword defines what happens if dt is an occurrence.
// With inc == True, if dt itself is an occurrence, it will be returned.
// It is only supported second precision.
func (set *Set) After(dt time.Time, inc bool) time.Time {
	return after(set.Iterator(), dt, inc)
}

// StrToRRuleSet converts string to RRuleSet
func StrToRRuleSet(s string) (*Set, error) {
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
func StrSliceToRRuleSet(ss []string) (*Set, error) {
	return StrSliceToRRuleSetInLoc(ss, time.UTC)
}

// StrSliceToRRuleSetInLoc is same as StrSliceToRRuleSet, but by default parses local times
// in specified default location
func StrSliceToRRuleSetInLoc(ss []string, defaultLoc *time.Location) (*Set, error) {
	if len(ss) == 0 {
		return &Set{}, nil
	}

	set := Set{}

	// According to RFC DTSTART is always the first line.
	firstName, err := processRRuleName(ss[0])
	if err != nil {
		return nil, err
	}
	if firstName == "DTSTART" {
		dt, err := StrToDtStart(ss[0][len(firstName)+1:], defaultLoc)
		if err != nil {
			return nil, fmt.Errorf("StrToDtStart failed: %v", err)
		}
		// default location should be taken from DTSTART property to correctly
		// parse local times met in RDATE,EXDATE and other rules
		defaultLoc = dt.Location()
		set.DTStart(dt)
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
			rOpt, err := StrToROptionInLocation(rule, defaultLoc)
			if err != nil {
				return nil, fmt.Errorf("StrToROption failed: %v", err)
			}
			r, err := NewRRule(*rOpt)
			if err != nil {
				return nil, fmt.Errorf("NewRRule failed: %v", r)
			}

			set.RRule(r)
		case "RDATE", "EXDATE":
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
