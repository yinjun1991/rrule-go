// 2017-2022, Teambition. All rights reserved.

package rrule

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// MAXYEAR
const (
	MAXYEAR = 9999
)

const (
	// DateTimeFormat is date-time format used in iCalendar (RFC 5545)
	DateTimeFormat = "20060102T150405Z"
	// LocalDateTimeFormat is a date-time format without Z prefix
	LocalDateTimeFormat = "20060102T150405"
	// DateFormat is date format used in iCalendar (RFC 5545)
	DateFormat = "20060102"
)

func timeToStr(time time.Time) string {
	return time.UTC().Format(DateTimeFormat)
}

func strToTimeInLoc(str string, loc *time.Location) (time.Time, error) {
	if len(str) == len(DateFormat) {
		return time.ParseInLocation(DateFormat, str, loc)
	}
	if len(str) == len(LocalDateTimeFormat) {
		return time.ParseInLocation(LocalDateTimeFormat, str, loc)
	}
	// date-time format carries zone info
	return time.Parse(DateTimeFormat, str)
}

func appendIntsOption(options []string, key string, value []int) []string {
	if len(value) == 0 {
		return options
	}
	valueStr := make([]string, len(value))
	for i, v := range value {
		valueStr[i] = strconv.Itoa(v)
	}
	return append(options, fmt.Sprintf("%s=%s", key, strings.Join(valueStr, ",")))
}

func strToInts(value string) ([]int, error) {
	contents := strings.Split(value, ",")
	result := make([]int, len(contents))
	var e error
	for i, s := range contents {
		result[i], e = strconv.Atoi(s)
		if e != nil {
			return nil, e
		}
	}
	return result, nil
}

// https://tools.ietf.org/html/rfc5545#section-3.3.5
// DTSTART:19970714T133000                       ; Local time
// DTSTART:19970714T173000Z                      ; UTC time
// DTSTART;TZID=America/New_York:19970714T133000 ; Local time and time zone reference
func timeToRFCDatetimeStr(time time.Time) string {
	if time.Location().String() != "UTC" {
		return fmt.Sprintf(";TZID=%s:%s", time.Location().String(), time.Format(LocalDateTimeFormat))
	}
	return fmt.Sprintf(":%s", time.Format(DateTimeFormat))
}

// StrToDates is intended to parse RDATE and EXDATE properties supporting only
// VALUE=DATE-TIME (DATE and PERIOD are not supported).
// Accepts string with format: "VALUE=DATE-TIME;[TZID=...]:{time},{time},...,{time}"
// or simply "{time},{time},...{time}" and parses it to array of dates
// In case no time zone specified in str, when all dates are parsed in UTC
func StrToDates(str string) (ts []time.Time, err error) {
	return StrToDatesInLoc(str, time.UTC)
}

// StrToDatesInLoc same as StrToDates but it consideres default location to parse dates in
// in case no location specified with TZID parameter
func StrToDatesInLoc(str string, defaultLoc *time.Location) (ts []time.Time, err error) {
	tmp := strings.Split(str, ":")
	if len(tmp) > 2 {
		return nil, fmt.Errorf("bad format")
	}
	loc := defaultLoc
	if len(tmp) == 2 {
		params := strings.Split(tmp[0], ";")
		for _, param := range params {
			if strings.HasPrefix(param, "TZID=") {
				loc, err = parseTZID(param)
			} else if param != "VALUE=DATE-TIME" && param != "VALUE=DATE" {
				err = fmt.Errorf("unsupported: %v", param)
			}
			if err != nil {
				return nil, fmt.Errorf("bad dates param: %s", err.Error())
			}
		}
		tmp = tmp[1:]
	}
	for _, datestr := range strings.Split(tmp[0], ",") {
		t, err := strToTimeInLoc(datestr, loc)
		if err != nil {
			return nil, fmt.Errorf("strToTime failed: %v", err)
		}
		ts = append(ts, t)
	}
	return
}

// processRRuleName processes the name of an RRule off a multi-line RRule set
func processRRuleName(line string) (string, error) {
	line = strings.ToUpper(strings.TrimSpace(line))
	if line == "" {
		return "", fmt.Errorf("bad format %v", line)
	}

	nameLen := strings.IndexAny(line, ";:")
	if nameLen <= 0 {
		return "", fmt.Errorf("bad format %v", line)
	}

	name := line[:nameLen]
	if strings.IndexAny(name, "=") > 0 {
		return "", fmt.Errorf("bad format %v", line)
	}

	return name, nil
}

// StrToDtStart accepts string with format: "(TZID={timezone}:)?{time}" or "VALUE=DATE:{date}" and parses it to a date
// may be used to parse DTSTART rules, without the DTSTART; part.
func StrToDtStart(str string, defaultLoc *time.Location) (time.Time, error) {
	// Handle VALUE=DATE parameter for all-day events
	if strings.HasPrefix(str, "VALUE=DATE:") {
		dateStr := str[len("VALUE=DATE:"):]
		// Parse DATE format (YYYYMMDD) for all-day events
		return strToTimeInLoc(dateStr, time.UTC) // All-day events use floating time (UTC)
	}

	tmp := strings.Split(str, ":")
	if len(tmp) > 2 || len(tmp) == 0 {
		return time.Time{}, fmt.Errorf("bad format")
	}

	if len(tmp) == 2 {
		// tzid
		loc, err := parseTZID(tmp[0])
		if err != nil {
			return time.Time{}, err
		}
		return strToTimeInLoc(tmp[1], loc)
	}
	// no tzid, len == 1
	return strToTimeInLoc(tmp[0], defaultLoc)
}

func parseTZID(s string) (*time.Location, error) {
	if !strings.HasPrefix(s, "TZID=") || len(s) == len("TZID=") {
		return nil, fmt.Errorf("bad TZID parameter format")
	}
	return time.LoadLocation(s[len("TZID="):])
}

// Python: MO-SU: 0 - 6
// Golang: SU-SAT 0 - 6
func toPyWeekday(from time.Weekday) int {
	return []int{6, 0, 1, 2, 3, 4, 5}[from]
}

// year -> 1 if leap year, else 0."
func isLeap(year int) int {
	if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
		return 1
	}
	return 0
}

// daysIn returns the number of days in a month for a given year.
func daysIn(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// mod in Python
func pymod(a, b int) int {
	r := a % b
	// If r and b differ in sign, add b to wrap the result to the correct sign.
	if r*b < 0 {
		r += b
	}
	return r
}

// divmod in Python
func divmod(a, b int) (div, mod int) {
	return int(math.Floor(float64(a) / float64(b))), pymod(a, b)
}

func contains(list []int, elem int) bool {
	for _, t := range list {
		if t == elem {
			return true
		}
	}
	return false
}

func timeContains(list []time.Time, elem time.Time) bool {
	for _, t := range list {
		if t.Equal(elem) {
			return true
		}
	}
	return false
}

func repeat(value, count int) []int {
	result := []int{}
	for i := 0; i < count; i++ {
		result = append(result, value)
	}
	return result
}

func concat(slices ...[]int) []int {
	result := []int{}
	for _, item := range slices {
		result = append(result, item...)
	}
	return result
}

func rang(start, end int) []int {
	result := []int{}
	for i := start; i < end; i++ {
		result = append(result, i)
	}
	return result
}

func pySubscript(slice []int, index int) (int, error) {
	if index < 0 {
		index += len(slice)
	}
	if index < 0 || index >= len(slice) {
		return 0, errors.New("index error")
	}
	return slice[index], nil
}

func timeSliceIterator(s []time.Time) func() (time.Time, bool) {
	index := 0
	return func() (time.Time, bool) {
		if index >= len(s) {
			return time.Time{}, false
		}
		result := s[index]
		index++
		return result, true
	}
}

func easter(year int) time.Time {
	g := year % 19
	c := year / 100
	h := (c - c/4 - (8*c+13)/25 + 19*g + 15) % 30
	i := h - (h/28)*(1-(h/28)*(29/(h+1))*((21-g)/11))
	j := (year + year/4 + i + 2 - c + c/4) % 7
	p := i - j
	d := 1 + (p+27+(p+6)/40)%31
	m := 3 + (p+26)/30
	return time.Date(year, time.Month(m), d, 0, 0, 0, 0, time.UTC)
}

func all(next Next) []time.Time {
	result := []time.Time{}
	for {
		v, ok := next()
		if !ok {
			return result
		}
		result = append(result, v)
	}
}

func between(next Next, after, before time.Time, inc bool) []time.Time {
	result := []time.Time{}
	for {
		v, ok := next()
		if !ok || inc && v.After(before) || !inc && !v.Before(before) {
			return result
		}
		if inc && !v.Before(after) || !inc && v.After(after) {
			result = append(result, v)
		}
	}
}

func before(next Next, dt time.Time, inc bool) time.Time {
	result := time.Time{}
	for {
		v, ok := next()
		if !ok || inc && v.After(dt) || !inc && !v.Before(dt) {
			return result
		}
		result = v
	}
}

func after(next Next, dt time.Time, inc bool) time.Time {
	for {
		v, ok := next()
		if !ok {
			return time.Time{}
		}
		if inc && !v.Before(dt) || !inc && v.After(dt) {
			return v
		}
	}
}

type optInt struct {
	Int     int
	Defined bool
}

func prepareTimeSet(set *[]time.Time, length int) {
	if len(*set) < length {
		*set = make([]time.Time, 0, length)
		return
	}

	*set = (*set)[:0]
}
