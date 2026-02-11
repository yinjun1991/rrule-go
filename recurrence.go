// 2017-2022, Teambition. All rights reserved.

package rrule

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Recurrence allows more complex recurrence setups, mixing multiple rules, dates, exclusion rules, and exclusion dates
type Recurrence struct {
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
	intervalExplicit        bool
	bymonthExplicit         bool
	bymonthdayExplicit      bool
	byweekdayExplicit       bool
	byhourExplicit          bool
	byminuteExplicit        bool
	bysecondExplicit        bool
	hasRule                 bool
}

// New builds a recurrence from ROption.
// Returns an error when option validation fails; otherwise returns a Recurrence
// even if the option contains no rule parts.
func New(option ROption) (*Recurrence, error) {
	rec := &Recurrence{}
	if err := rec.setRuleOptions(option); err != nil {
		return nil, err
	}
	return rec, nil
}

// NewWithDTStart builds a recurrence using the provided dtstart and allDay flags.
// Any DTSTART lines in the input are ignored; the dtstart argument is authoritative.
// Returns an error when the input lines are malformed; returns an empty Recurrence
// when lines are empty or normalize to no usable rules.
func NewWithDTStart(dtstart time.Time, allDay bool, lines ...string) (*Recurrence, error) {
	rec := &Recurrence{}
	rec.SetAllDay(allDay)
	rec.DTStart(dtstart)

	if len(lines) == 0 {
		return rec, nil
	}

	normalized, err := NormalizeRecurrenceRuleset(lines)
	if err != nil {
		return nil, err
	}
	if len(normalized) == 0 {
		return rec, nil
	}

	filtered := make([]string, 0, len(normalized))
	for _, line := range normalized {
		if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(line)), "DTSTART") {
			continue
		}
		filtered = append(filtered, line)
	}
	if len(filtered) == 0 {
		return rec, nil
	}

	defaultLoc := time.UTC
	if !rec.GetDTStart().IsZero() {
		defaultLoc = rec.GetDTStart().Location()
	}

	var dtstartLineForRRULE string
	if !rec.GetDTStart().IsZero() {
		if rec.allDay {
			dtstartLineForRRULE = fmt.Sprintf("DTSTART;VALUE=DATE:%s", rec.GetDTStart().Format(DateFormat))
		} else {
			dtstartLineForRRULE = fmt.Sprintf("DTSTART%s", timeToRFCDatetimeStr(rec.GetDTStart()))
		}
	}

	for _, line := range filtered {
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
			rOpt, err := parseROptionFromString(rruleInput)
			if err != nil {
				return nil, fmt.Errorf("parseROption failed: %v", err)
			}
			err = rec.setRuleOptions(*rOpt)
			if err != nil {
				return nil, fmt.Errorf("NewRRule failed: %v", err)
			}
		case "RDATE", "EXDATE":
			if !rec.allDay && containsValueDateParam(rule) {
				rec.SetAllDay(true)
			}

			ts, err := StrToDatesInLoc(rule, defaultLoc)
			if err != nil {
				return nil, fmt.Errorf("strToDates failed: %v", err)
			}
			for _, t := range ts {
				if name == "RDATE" {
					rec.RDate(t)
				} else {
					rec.ExDate(t)
				}
			}
		}
	}

	return rec, nil
}

// Parse builds a recurrence from lines.
// Returns an error when lines are malformed; returns an empty Recurrence when
// lines are empty or normalize to no usable rules.
func Parse(lines ...string) (*Recurrence, error) {
	if len(lines) == 0 {
		return &Recurrence{}, nil
	}

	normalized, err := NormalizeRecurrenceRuleset(lines)
	if err != nil {
		return nil, err
	}
	if len(normalized) == 0 {
		return &Recurrence{}, nil
	}
	lines = normalized

	defaultLoc := time.UTC
	set := Recurrence{}
	var dtstartLineForRRULE string

	firstName, err := processRRuleName(lines[0])
	if err != nil {
		return nil, err
	}
	if firstName == "DTSTART" {
		dtstartField := lines[0][len(firstName)+1:]
		if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(dtstartField)), "VALUE=DATE:") {
			set.SetAllDay(true)
		}

		dt, err := StrToDtStart(dtstartField, defaultLoc)
		if err != nil {
			return nil, fmt.Errorf("StrToDtStart failed: %v", err)
		}
		defaultLoc = dt.Location()
		set.DTStart(dt)
		if !set.GetDTStart().IsZero() {
			if set.allDay {
				dtstartLineForRRULE = fmt.Sprintf("DTSTART;VALUE=DATE:%s", set.GetDTStart().Format(DateFormat))
			} else {
				dtstartLineForRRULE = fmt.Sprintf("DTSTART%s", timeToRFCDatetimeStr(set.GetDTStart()))
			}
		}
		lines = lines[1:]
	}

	for _, line := range lines {
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
			rOpt, err := parseROptionFromString(rruleInput)
			if err != nil {
				return nil, fmt.Errorf("parseROption failed: %v", err)
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

// NormalizeRecurrenceRuleset cleans and normalizes recurrence lines.
// It trims whitespace, removes empty entries, normalizes RRULE prefixing,
// and keeps only the first RRULE when multiple are present.
func NormalizeRecurrenceRuleset(ruleset []string) ([]string, error) {
	if len(ruleset) == 0 {
		return nil, nil
	}

	normalized := make([]string, 0, len(ruleset))
	foundRRule := false

	for _, rule := range ruleset {
		rule = strings.TrimSpace(rule)
		if rule == "" {
			continue
		}

		upperRule := strings.ToUpper(rule)
		if strings.HasPrefix(upperRule, "DTSTART") {
			normalized = append(normalized, rule)
			continue
		}

		normalizedRule, err := normalizeRecurrenceLine(rule)
		if err != nil {
			return nil, fmt.Errorf("invalid recurrence string '%s': %w", rule, err)
		}

		if strings.HasPrefix(strings.ToUpper(normalizedRule), "RRULE:") {
			if foundRRule {
				continue
			}
			foundRRule = true
		}

		normalized = append(normalized, normalizedRule)
	}

	if len(normalized) == 0 {
		return nil, nil
	}

	return normalized, nil
}

// normalizeRecurrenceLine normalizes a single RRULE/RDATE/EXDATE line.
func normalizeRecurrenceLine(rule string) (string, error) {
	upperRule := strings.ToUpper(rule)

	if strings.HasPrefix(upperRule, "DTSTART") {
		return rule, nil
	}
	if strings.HasPrefix(upperRule, "RRULE:") {
		content := rule[len("RRULE:"):]
		if err := validateRRuleProperties(content); err != nil {
			return "", err
		}
		return rule, nil
	}
	if strings.HasPrefix(upperRule, "RDATE:") || strings.HasPrefix(upperRule, "RDATE;") ||
		strings.HasPrefix(upperRule, "EXDATE:") || strings.HasPrefix(upperRule, "EXDATE;") {
		return rule, nil
	}
	if isRRuleProperties(rule) {
		if err := validateRRuleProperties(rule); err != nil {
			return "", err
		}
		return "RRULE:" + rule, nil
	}
	return "", fmt.Errorf("unrecognized rule format")
}

// validateRRuleProperties validates the RRULE properties without the RRULE prefix.
func validateRRuleProperties(content string) error {
	content = strings.TrimSpace(content)
	if content == "" {
		return fmt.Errorf("empty rrule content")
	}

	upperContent := strings.ToUpper(content)
	if !strings.Contains(upperContent, "FREQ=") {
		return fmt.Errorf("rrule must contain FREQ parameter")
	}

	return nil
}

// isRRuleProperties reports whether a string looks like RRULE properties.
func isRRuleProperties(content string) bool {
	upperContent := strings.ToUpper(strings.TrimSpace(content))
	return strings.Contains(upperContent, "FREQ=") &&
		!strings.HasPrefix(upperContent, "RRULE:") &&
		!strings.HasPrefix(upperContent, "RDATE:") &&
		!strings.HasPrefix(upperContent, "EXDATE:") &&
		!strings.HasPrefix(upperContent, "DTSTART")
}

func ParseRRuleString(rfcString string) (*Recurrence, error) {
	option, err := parseROptionFromString(rfcString)
	if err != nil {
		return nil, err
	}
	return New(*option)
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
	}

	if option.Dtstart.IsZero() && !r.dtstart.IsZero() {
		option.Dtstart = r.dtstart
	}

	if err := r.applyRule(option); err != nil {
		return err
	}

	r.hasRule = true
	return nil
}

func (r *Recurrence) rebuildRule() {
	_ = r.applyRule(r.ruleOptionFromState())
	r.hasRule = true
}

func (r *Recurrence) applyRule(option ROption) error {
	if err := validateBounds(option); err != nil {
		return err
	}
	if option.AllDay {
		option.Byhour = nil
		option.Byminute = nil
		option.Bysecond = nil
		if option.Freq >= HOURLY {
			option.Freq = DAILY
		}
	}
	r.allDay = option.AllDay
	r.intervalExplicit = option.Interval > 0
	r.bymonthExplicit = len(option.Bymonth) != 0
	r.bymonthdayExplicit = len(option.Bymonthday) != 0
	r.byweekdayExplicit = len(option.Byweekday) != 0
	r.byhourExplicit = len(option.Byhour) != 0
	r.byminuteExplicit = len(option.Byminute) != 0
	r.bysecondExplicit = len(option.Bysecond) != 0

	r.freq = option.Freq

	if option.Interval < 1 {
		option.Interval = 1
	}
	r.interval = option.Interval

	if option.Count < 0 {
		option.Count = 0
	}
	r.count = option.Count

	if option.Dtstart.IsZero() {
		option.Dtstart = time.Now().UTC()
	}

	if option.AllDay {
		year, month, day := option.Dtstart.Date()
		option.Dtstart = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	} else {
		option.Dtstart = option.Dtstart.Truncate(time.Second)
	}
	r.dtstart = option.Dtstart

	if option.Until.IsZero() {
		r.until = r.dtstart.Add(time.Duration(1<<63 - 1))
	} else {
		if option.AllDay {
			year, month, day := option.Until.Date()
			option.Until = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		} else {
			option.Until = option.Until.Truncate(time.Second)
		}
		r.until = option.Until
	}

	if option.RDate != nil {
		r.SetRDates(option.RDate)
	}
	if option.EXDate != nil {
		r.SetExDates(option.EXDate)
	}

	r.wkst = option.Wkst.weekday
	r.bymonthday = nil
	r.bynmonthday = nil
	r.byweekday = nil
	r.bynweekday = nil
	r.bysetpos = option.Bysetpos

	if len(option.Byweekno) == 0 &&
		len(option.Bymonthday) == 0 &&
		len(option.Byyearday) == 0 &&
		len(option.Byweekday) == 0 &&
		len(option.Byeaster) == 0 {
		if r.freq == YEARLY {
			if len(option.Bymonth) == 0 {
				option.Bymonth = []int{int(r.dtstart.Month())}
			}
			option.Bymonthday = []int{r.dtstart.Day()}
		} else if r.freq == MONTHLY {
			option.Bymonthday = []int{r.dtstart.Day()}
		} else if r.freq == WEEKLY {
			option.Byweekday = []Weekday{{weekday: toPyWeekday(r.dtstart.Weekday())}}
		}
	}
	r.bymonth = option.Bymonth
	r.byyearday = option.Byyearday
	r.byeaster = option.Byeaster
	for _, mday := range option.Bymonthday {
		if mday > 0 {
			r.bymonthday = append(r.bymonthday, mday)
		} else if mday < 0 {
			r.bynmonthday = append(r.bynmonthday, mday)
		}
	}
	r.byweekno = option.Byweekno
	for _, wday := range option.Byweekday {
		if wday.n == 0 || r.freq > MONTHLY {
			r.byweekday = append(r.byweekday, wday.weekday)
		} else {
			r.bynweekday = append(r.bynweekday, wday)
		}
	}
	if len(option.Byhour) == 0 {
		if r.freq < HOURLY {
			if option.AllDay {
				r.byhour = []int{0}
			} else {
				r.byhour = []int{r.dtstart.Hour()}
			}
		}
	} else {
		r.byhour = option.Byhour
	}
	if len(option.Byminute) == 0 {
		if r.freq < MINUTELY {
			if option.AllDay {
				r.byminute = []int{0}
			} else {
				r.byminute = []int{r.dtstart.Minute()}
			}
		}
	} else {
		r.byminute = option.Byminute
	}
	if len(option.Bysecond) == 0 {
		if r.freq < SECONDLY {
			if option.AllDay {
				r.bysecond = []int{0}
			} else {
				r.bysecond = []int{r.dtstart.Second()}
			}
		}
	} else {
		r.bysecond = option.Bysecond
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

func (r *Recurrence) ruleOptionFromState() ROption {
	option := ROption{
		Freq:      r.freq,
		Dtstart:   r.dtstart,
		Interval:  r.interval,
		Wkst:      Weekday{weekday: r.wkst},
		Count:     r.count,
		AllDay:    r.allDay,
		Bysetpos:  cloneIntSlice(r.bysetpos),
		Bymonth:   cloneIntSlice(r.bymonth),
		Byyearday: cloneIntSlice(r.byyearday),
		Byweekno:  cloneIntSlice(r.byweekno),
		Byhour:    cloneIntSlice(r.byhour),
		Byminute:  cloneIntSlice(r.byminute),
		Bysecond:  cloneIntSlice(r.bysecond),
		Byeaster:  cloneIntSlice(r.byeaster),
	}

	if !r.intervalExplicit && r.interval == 1 {
		option.Interval = 0
	}

	if !r.until.IsZero() {
		maxUntil := r.dtstart.Add(time.Duration(1<<63 - 1))
		if !r.until.Equal(maxUntil) {
			option.Until = r.until
		}
	}

	if !r.bymonthExplicit {
		option.Bymonth = nil
	}

	if r.bymonthdayExplicit {
		byMonthDay := make([]int, 0, len(r.bymonthday)+len(r.bynmonthday))
		byMonthDay = append(byMonthDay, r.bymonthday...)
		byMonthDay = append(byMonthDay, r.bynmonthday...)
		option.Bymonthday = byMonthDay
	} else {
		option.Bymonthday = nil
	}

	if r.byweekdayExplicit {
		byWeekday := make([]Weekday, 0, len(r.byweekday)+len(r.bynweekday))
		for _, wday := range r.byweekday {
			byWeekday = append(byWeekday, Weekday{weekday: wday})
		}
		byWeekday = append(byWeekday, r.bynweekday...)
		option.Byweekday = byWeekday
	} else {
		option.Byweekday = nil
	}

	if !r.byhourExplicit {
		option.Byhour = nil
	}
	if !r.byminuteExplicit {
		option.Byminute = nil
	}
	if !r.bysecondExplicit {
		option.Bysecond = nil
	}

	return option
}

func cloneIntSlice(values []int) []int {
	if len(values) == 0 {
		return nil
	}
	out := make([]int, len(values))
	copy(out, values)
	return out
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

// Strings returns a slice of all the recurrence rules for a set
func (set *Recurrence) Strings() []string {
	var res []string

	str := set.DTStartString()
	if str != "" {
		res = append(res, str)
	}

	if set.hasRule {
		str = set.RRuleString()
		if str != "" {
			res = append(res, str)
		}
	}

	str = set.RDateString()
	if str != "" {
		res = append(res, str)
	}

	str = set.EXDateString()
	if str != "" {
		res = append(res, str)
	}

	return res
}

// String returns the full RFC 5545 recurrence text, one property per line.
// Example:
// DTSTART:20240101T090000Z
// RRULE:FREQ=DAILY;COUNT=2
// EXDATE:20240102T090000Z
func (set *Recurrence) String() string {
	res := set.Strings()
	return strings.Join(res, "\n")
}

// DTStartString returns DTSTART serialized as a single line.
// Example: DTSTART;VALUE=DATE:20240101
// Example: DTSTART:20240101T090000Z
// Example: DTSTART;TZID=Asia/Shanghai:20240101T090000
func (set *Recurrence) DTStartString() string {
	if set.dtstart.IsZero() {
		return ""
	}
	// No colon, DTSTART may have TZID, which would require a semicolon after DTSTART
	// RFC 5545: For all-day events, use VALUE=DATE format
	if set.allDay {
		// All-day events should use VALUE=DATE format as per RFC 5545
		return fmt.Sprintf("DTSTART;VALUE=DATE:%s", set.dtstart.Format(DateFormat))
	}

	return fmt.Sprintf("DTSTART%s", timeToRFCDatetimeStr(set.dtstart))
}

// RRuleString returns RRULE serialized as a single line without DTSTART.
// Example: RRULE:FREQ=DAILY;COUNT=5
func (set *Recurrence) RRuleString() string {
	return fmt.Sprintf("RRULE:%s", set.rrulePropertiesString())
}

// EXDateString returns EXDATE lines serialized as a single string.
// When DTSTART is available, EXDATE is normalized to the DTSTART timezone to match
// RFC 5545 requirements; UTC values use a trailing Z without TZID.
// When DTSTART is missing, EXDATE entries are grouped by timezone.
// Example: EXDATE;VALUE=DATE:20240110,20240112
// Example: EXDATE:20240110T090000Z,20240112T090000Z
// Example: EXDATE;TZID=Asia/Shanghai:20240110T090000,20240112T090000
func (set *Recurrence) EXDateString() string {
	if len(set.exdate) == 0 {
		return ""
	}
	if set.allDay {
		values := make([]string, 0, len(set.exdate))
		for _, item := range set.exdate {
			values = append(values, item.Format(DateFormat))
		}
		return fmt.Sprintf("EXDATE;VALUE=DATE:%s", strings.Join(values, ","))
	}

	valuesByTZID := make(map[string][]string)
	var tzidOrder []string
	for _, item := range set.exdate {
		tzid := item.Location().String()
		if _, ok := valuesByTZID[tzid]; !ok {
			tzidOrder = append(tzidOrder, tzid)
		}
		if tzid == "UTC" {
			valuesByTZID[tzid] = append(valuesByTZID[tzid], item.Format(DateTimeFormat))
		} else {
			valuesByTZID[tzid] = append(valuesByTZID[tzid], item.Format(LocalDateTimeFormat))
		}
	}

	lines := make([]string, 0, len(valuesByTZID))
	for _, tzid := range tzidOrder {
		values := strings.Join(valuesByTZID[tzid], ",")
		if tzid == "UTC" {
			lines = append(lines, fmt.Sprintf("EXDATE:%s", values))
		} else {
			lines = append(lines, fmt.Sprintf("EXDATE;TZID=%s:%s", tzid, values))
		}
	}
	return strings.Join(lines, "\n")
}

// RDateString returns RDATE lines serialized as a single string.
// When DTSTART is available, RDATE is normalized to the DTSTART timezone to match
// RFC 5545 requirements; UTC values use a trailing Z without TZID.
// When DTSTART is missing, RDATE entries are grouped by timezone.
// Example: RDATE;VALUE=DATE:20240301,20240303
// Example: RDATE:20240301T090000Z,20240305T090000Z
// Example: RDATE;TZID=Asia/Shanghai:20240301T090000,20240305T090000
func (set *Recurrence) RDateString() string {
	if len(set.rdate) == 0 {
		return ""
	}
	if set.allDay {
		values := make([]string, 0, len(set.rdate))
		for _, item := range set.rdate {
			values = append(values, item.Format(DateFormat))
		}
		return fmt.Sprintf("RDATE;VALUE=DATE:%s", strings.Join(values, ","))
	}

	valuesByTZID := make(map[string][]string)
	var tzidOrder []string
	for _, item := range set.rdate {
		tzid := item.Location().String()
		if _, ok := valuesByTZID[tzid]; !ok {
			tzidOrder = append(tzidOrder, tzid)
		}
		if tzid == "UTC" {
			valuesByTZID[tzid] = append(valuesByTZID[tzid], item.Format(DateTimeFormat))
		} else {
			valuesByTZID[tzid] = append(valuesByTZID[tzid], item.Format(LocalDateTimeFormat))
		}
	}

	lines := make([]string, 0, len(valuesByTZID))
	for _, tzid := range tzidOrder {
		values := strings.Join(valuesByTZID[tzid], ",")
		if tzid == "UTC" {
			lines = append(lines, fmt.Sprintf("RDATE:%s", values))
		} else {
			lines = append(lines, fmt.Sprintf("RDATE;TZID=%s:%s", tzid, values))
		}
	}
	return strings.Join(lines, "\n")
}

// DTStart sets dtstart property for set.
// It will be truncated to second precision.
func (set *Recurrence) DTStart(dtstart time.Time) {
	prevDtstart := set.dtstart
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

	if set.hasRule {
		if !set.until.IsZero() && !prevDtstart.IsZero() {
			maxUntil := prevDtstart.Add(time.Duration(1<<63 - 1))
			if set.until.Equal(maxUntil) {
				set.until = set.dtstart.Add(time.Duration(1<<63 - 1))
			}
		}
		set.rebuildRule()
	}
}

// GetDTStart gets DTSTART for set
func (set *Recurrence) GetDTStart() time.Time {
	return set.dtstart
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
	}

	if set.hasRule {
		set.rebuildRule()
	}
}

// IsAllDay returns whether the set is configured for all-day events.
func (set *Recurrence) IsAllDay() bool {
	return set.allDay
}

// rrulePropertiesString returns the RRULE value without the "RRULE:" prefix.
// Example: FREQ=DAILY;COUNT=5
// Example: FREQ=WEEKLY;INTERVAL=2;BYDAY=MO,WE,FR
func (set *Recurrence) rrulePropertiesString() string {
	result := []string{fmt.Sprintf("FREQ=%v", set.freq)}
	if set.intervalExplicit && set.interval != 1 {
		result = append(result, fmt.Sprintf("INTERVAL=%v", set.interval))
	}
	if set.wkst != MO.weekday {
		result = append(result, fmt.Sprintf("WKST=%v", Weekday{weekday: set.wkst}))
	}
	if set.count > 0 {
		result = append(result, fmt.Sprintf("COUNT=%v", set.count))
	}
	if !set.until.IsZero() {
		maxUntil := set.dtstart.Add(time.Duration(1<<63 - 1))
		if !set.until.Equal(maxUntil) {
			if set.allDay {
				until := set.until.In(set.dtstart.Location())
				result = append(result, fmt.Sprintf("UNTIL=%v", until.Format(DateFormat)))
			} else {
				result = append(result, fmt.Sprintf("UNTIL=%v", timeToUTCStr(set.until)))
			}
		}
	}
	result = appendIntsOption(result, "BYSETPOS", set.bysetpos)
	if set.bymonthExplicit {
		result = appendIntsOption(result, "BYMONTH", set.bymonth)
	}
	if set.bymonthdayExplicit {
		byMonthDay := make([]int, 0, len(set.bymonthday)+len(set.bynmonthday))
		byMonthDay = append(byMonthDay, set.bymonthday...)
		byMonthDay = append(byMonthDay, set.bynmonthday...)
		result = appendIntsOption(result, "BYMONTHDAY", byMonthDay)
	}
	result = appendIntsOption(result, "BYYEARDAY", set.byyearday)
	result = appendIntsOption(result, "BYWEEKNO", set.byweekno)
	if set.byweekdayExplicit {
		byWeekday := make([]Weekday, 0, len(set.byweekday)+len(set.bynweekday))
		for _, wday := range set.byweekday {
			byWeekday = append(byWeekday, Weekday{weekday: wday})
		}
		byWeekday = append(byWeekday, set.bynweekday...)
		valueStr := make([]string, len(byWeekday))
		for i, wday := range byWeekday {
			valueStr[i] = wday.String()
		}
		result = append(result, fmt.Sprintf("BYDAY=%s", strings.Join(valueStr, ",")))
	}
	if !set.allDay {
		if set.byhourExplicit {
			result = appendIntsOption(result, "BYHOUR", set.byhour)
		}
		if set.byminuteExplicit {
			result = appendIntsOption(result, "BYMINUTE", set.byminute)
		}
		if set.bysecondExplicit {
			result = appendIntsOption(result, "BYSECOND", set.bysecond)
		}
	}
	result = appendIntsOption(result, "BYEASTER", set.byeaster)
	return strings.Join(result, ";")
}

func parseROptionFromString(rfcString string) (*ROption, error) {
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
		return nil, errors.New("RRULE property FREQ is required")
	}
	return &result, nil
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
	return Parse(ss...)
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
