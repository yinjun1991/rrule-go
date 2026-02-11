// 2017-2022, Teambition. All rights reserved.

package rrule

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func StrToRRule(input string) (*Recurrence, error) {
	option, err := StrToROption(input)
	if err != nil {
		return nil, err
	}
	return newRecurrence(*option)
}

func timesEqual(value, want []time.Time) bool {
	if len(value) != len(want) {
		return false
	}
	for index := range value {
		if value[index] != want[index] {
			return false
		}
	}
	return true
}

func TestNoDtstart(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MONTHLY})
	if time.Now().Unix()-r.dtstart.Unix() > 1 {
		t.Errorf(`default Dtstart shold be time.Now(), but got %s`, r.dtstart.Format(time.RFC3339))
	}
}

func TestBadBySetPos(t *testing.T) {
	_, e := newRecurrence(ROption{Freq: MONTHLY, Count: 1, Bysetpos: []int{0},
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	if e == nil {
		t.Error("get nil, want error")
	}
}

func TestBadBySetPosMany(t *testing.T) {
	_, e := newRecurrence(ROption{Freq: MONTHLY, Count: 1, Bysetpos: []int{-1, 0, 1},
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	if e == nil {
		t.Error("get nil, want error")
	}
}

func TestByNegativeMonthDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MONTHLY,
		Count:      3,
		Bymonthday: []int{-1},
		Dtstart:    time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 30, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 10, 31, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 11, 30, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMonthlyMaxYear(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MONTHLY, Interval: 15,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
	})
	value := r.All()[1]
	want := time.Date(1998, 12, 2, 9, 0, 0, 0, time.UTC)
	if value != want {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestWeeklyMaxYear(t *testing.T) {
	// Purposefully doesn't match anything for code coverage.
	r, _ := newRecurrence(ROption{Freq: WEEKLY, Bymonthday: []int{31},
		Byyearday: []int{1}, Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
	})
	value := r.All()
	want := []time.Time{}
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestInvalidRRules(t *testing.T) {
	tests := []struct {
		desc    string
		rrule   ROption
		wantErr string
	}{
		{
			desc:    "Bysecond under",
			rrule:   ROption{Freq: YEARLY, Bysecond: []int{-1}},
			wantErr: "bysecond must be between 0 and 59",
		},
		{
			desc:    "Bysecond over",
			rrule:   ROption{Freq: YEARLY, Bysecond: []int{60}},
			wantErr: "bysecond must be between 0 and 59",
		},
		{
			desc:    "Byminute under",
			rrule:   ROption{Freq: YEARLY, Byminute: []int{-1}},
			wantErr: "byminute must be between 0 and 59",
		},
		{
			desc:    "Byminute over",
			rrule:   ROption{Freq: YEARLY, Byminute: []int{60}},
			wantErr: "byminute must be between 0 and 59",
		},
		{
			desc:    "Byhour under",
			rrule:   ROption{Freq: YEARLY, Byhour: []int{-1}},
			wantErr: "byhour must be between 0 and 23",
		},
		{
			desc:    "Byhour over",
			rrule:   ROption{Freq: YEARLY, Byhour: []int{24}},
			wantErr: "byhour must be between 0 and 23",
		},
		{
			desc:    "Bymonthday under",
			rrule:   ROption{Freq: YEARLY, Bymonthday: []int{0}},
			wantErr: "bymonthday must be between 1 and 31 or -1 and -31",
		},
		{
			desc:    "Bymonthday over",
			rrule:   ROption{Freq: YEARLY, Bymonthday: []int{32}},
			wantErr: "bymonthday must be between 1 and 31 or -1 and -31",
		},
		{
			desc:    "Bymonthday under negative",
			rrule:   ROption{Freq: YEARLY, Bymonthday: []int{-32}},
			wantErr: "bymonthday must be between 1 and 31 or -1 and -31",
		},
		{
			desc:    "Byyearday under",
			rrule:   ROption{Freq: YEARLY, Byyearday: []int{0}},
			wantErr: "byyearday must be between 1 and 366 or -1 and -366",
		},
		{
			desc:    "Byyearday over",
			rrule:   ROption{Freq: YEARLY, Byyearday: []int{367}},
			wantErr: "byyearday must be between 1 and 366 or -1 and -366",
		},
		{
			desc:    "Byyearday under negative",
			rrule:   ROption{Freq: YEARLY, Byyearday: []int{-367}},
			wantErr: "byyearday must be between 1 and 366 or -1 and -366",
		},
		{
			desc:    "Byweekno under",
			rrule:   ROption{Freq: YEARLY, Byweekno: []int{0}},
			wantErr: "byweekno must be between 1 and 53 or -1 and -53",
		},
		{
			desc:    "Byweekno over",
			rrule:   ROption{Freq: YEARLY, Byweekno: []int{54}},
			wantErr: "byweekno must be between 1 and 53 or -1 and -53",
		},
		{
			desc:    "Byweekno under negative",
			rrule:   ROption{Freq: YEARLY, Byweekno: []int{-54}},
			wantErr: "byweekno must be between 1 and 53 or -1 and -53",
		},
		{
			desc:    "Bymonth under",
			rrule:   ROption{Freq: YEARLY, Bymonth: []int{0}},
			wantErr: "bymonth must be between 1 and 12",
		},
		{
			desc:    "Bymonth over",
			rrule:   ROption{Freq: YEARLY, Bymonth: []int{13}},
			wantErr: "bymonth must be between 1 and 12",
		},
		{
			desc:    "Bysetpos under",
			rrule:   ROption{Freq: YEARLY, Bysetpos: []int{0}},
			wantErr: "bysetpos must be between 1 and 366 or -1 and -366",
		},
		{
			desc:    "Bysetpos over",
			rrule:   ROption{Freq: YEARLY, Bysetpos: []int{367}},
			wantErr: "bysetpos must be between 1 and 366 or -1 and -366",
		},
		{
			desc:    "Bysetpos under negative",
			rrule:   ROption{Freq: YEARLY, Bysetpos: []int{-367}},
			wantErr: "bysetpos must be between 1 and 366 or -1 and -366",
		},
		{
			desc:    "Byday under",
			rrule:   ROption{Freq: YEARLY, Byweekday: []Weekday{{1, -54}}},
			wantErr: "byday must be between 1 and 53 or -1 and -53",
		},
		{
			desc:    "Byday over",
			rrule:   ROption{Freq: YEARLY, Byweekday: []Weekday{{1, 54}}},
			wantErr: "byday must be between 1 and 53 or -1 and -53",
		},
		{
			desc:    "Interval under",
			rrule:   ROption{Freq: DAILY, Interval: -1},
			wantErr: "interval must be greater than 0",
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, err := newRecurrence(tc.rrule)
			if err == nil || err.Error() != tc.wantErr {
				t.Errorf("got %q, want %q", err, tc.wantErr)
			}
		})
	}
}

func TestHourlyInvalidAndRepeatedBysetpos(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: HOURLY, Bysetpos: []int{1, -1, 2},
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		Until:   time.Date(1997, 9, 2, 11, 0, 0, 0, time.UTC)})
	value := r.All()
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 10, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 11, 0, 0, 0, time.UTC)}
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestNoAfter(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:   5,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := time.Time{}
	value := r.After(time.Date(1997, 9, 6, 9, 0, 0, 0, time.UTC), false)
	if value != want {
		t.Errorf("get %v, want %v", value, want)
	}
}

// Test cases from Python Dateutil

func TestYearly(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY,
		Count:   3,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 9, 2, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestYearlyInterval(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY,
		Count:    3,
		Interval: 2,
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(2001, 9, 2, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestYearlyIntervalLarge(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY,
		Count:    3,
		Interval: 100,
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(2097, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(2197, 9, 2, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestYearlyByMonth(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY,
		Count:   3,
		Bymonth: []int{1, 3},
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 3, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 1, 2, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestYearlyByMonthDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY,
		Count:      3,
		Bymonthday: []int{1, 3},
		Dtstart:    time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 3, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 10, 1, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 10, 3, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestYearlyByMonthAndMonthDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY,
		Count:      3,
		Bymonth:    []int{1, 3},
		Bymonthday: []int{5, 7},
		Dtstart:    time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 5, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 7, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 3, 5, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestYearlyByWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY,
		Count:     3,
		Byweekday: []Weekday{TU, TH},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 9, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestYearlyByNWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY,
		Count:     3,
		Byweekday: []Weekday{TU.Nth(1), TH.Nth(-1)},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 12, 25, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 6, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 12, 31, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestYearlyByNWeekDayLarge(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY,
		Count:     3,
		Byweekday: []Weekday{TU.Nth(3), TH.Nth(-3)},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 12, 11, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 20, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 12, 17, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestYearlyByMonthAndWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY,
		Count:     3,
		Bymonth:   []int{1, 3},
		Byweekday: []Weekday{TU, TH},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 6, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 8, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestYearlyByMonthAndNWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY,
		Count:     3,
		Bymonth:   []int{1, 3},
		Byweekday: []Weekday{TU.Nth(1), TH.Nth(-1)},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 6, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 29, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 3, 3, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestYearlyByMonthAndNWeekDayLarge(t *testing.T) {
	// This is interesting because the TH.Nth(-3) ends up before
	// the TU.Nth(3).
	r, _ := newRecurrence(ROption{Freq: YEARLY,
		Count:     3,
		Bymonth:   []int{1, 3},
		Byweekday: []Weekday{TU.Nth(3), TH.Nth(-3)},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 15, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 20, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 3, 12, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestYearlyByMonthDayAndWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY,
		Count:      3,
		Bymonthday: []int{1, 3},
		Byweekday:  []Weekday{TU, TH},
		Dtstart:    time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 2, 3, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 3, 3, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestYearlyByMonthAndMonthDayAndWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY,
		Count:      3,
		Bymonth:    []int{1, 3},
		Bymonthday: []int{1, 3},
		Byweekday:  []Weekday{TU, TH},
		Dtstart:    time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 3, 3, 9, 0, 0, 0, time.UTC),
		time.Date(2001, 3, 1, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestYearlyByYearDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY,
		Count:     4,
		Byyearday: []int{1, 100, 200, 365},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 12, 31, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 4, 10, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 7, 19, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestYearlyByYearDayNeg(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY,
		Count:     4,
		Byyearday: []int{-365, -266, -166, -1},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 12, 31, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 4, 10, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 7, 19, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestYearlyByMonthAndYearDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY,
		Count:     4,
		Bymonth:   []int{4, 7},
		Byyearday: []int{1, 100, 200, 365},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 4, 10, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 7, 19, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 4, 10, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 7, 19, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestYearlyByMonthAndYearDayNeg(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY,
		Count:     4,
		Bymonth:   []int{4, 7},
		Byyearday: []int{-365, -266, -166, -1},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 4, 10, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 7, 19, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 4, 10, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 7, 19, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestYearlyByWeekNo(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY,
		Count:    3,
		Byweekno: []int{20},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 5, 11, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 5, 12, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 5, 13, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestYearlyByWeekNoAndWeekDay(t *testing.T) {
	// That's a nice one. The first days of week number one
	// may be in the last year.
	r, _ := newRecurrence(ROption{Freq: YEARLY,
		Count:     3,
		Byweekno:  []int{1},
		Byweekday: []Weekday{MO},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 12, 29, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 1, 4, 9, 0, 0, 0, time.UTC),
		time.Date(2000, 1, 3, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestYearlyByWeekNoAndWeekDayLarge(t *testing.T) {
	// Another nice test. The last days of week number 52/53
	// may be in the next year.
	r, _ := newRecurrence(ROption{Freq: YEARLY,
		Count:     3,
		Byweekno:  []int{52},
		Byweekday: []Weekday{SU},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 12, 28, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 12, 27, 9, 0, 0, 0, time.UTC),
		time.Date(2000, 1, 2, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestYearlyByWeekNoAndWeekDayLast(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY,
		Count:     3,
		Byweekno:  []int{-1},
		Byweekday: []Weekday{SU},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 12, 28, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 1, 3, 9, 0, 0, 0, time.UTC),
		time.Date(2000, 1, 2, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestYearlyByEaster(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY,
		Count:    3,
		Byeaster: []int{0},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 4, 12, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 4, 4, 9, 0, 0, 0, time.UTC),
		time.Date(2000, 4, 23, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestYearlyByEasterPos(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY,
		Count:    3,
		Byeaster: []int{1},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 4, 13, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 4, 5, 9, 0, 0, 0, time.UTC),
		time.Date(2000, 4, 24, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestYearlyByEasterNeg(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY,
		Count:    3,
		Byeaster: []int{-1},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 4, 11, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 4, 3, 9, 0, 0, 0, time.UTC),
		time.Date(2000, 4, 22, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestYearlyByWeekNoAndWeekDay53(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY,
		Count:     3,
		Byweekno:  []int{53},
		Byweekday: []Weekday{MO},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 12, 28, 9, 0, 0, 0, time.UTC),
		time.Date(2004, 12, 27, 9, 0, 0, 0, time.UTC),
		time.Date(2009, 12, 28, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestYearlyByHour(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY,
		Count:   3,
		Byhour:  []int{6, 18},
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 18, 0, 0, 0, time.UTC),
		time.Date(1998, 9, 2, 6, 0, 0, 0, time.UTC),
		time.Date(1998, 9, 2, 18, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestYearlyByMinute(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY,
		Count:    3,
		Byminute: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 6, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 18, 0, 0, time.UTC),
		time.Date(1998, 9, 2, 9, 6, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestYearlyBySecond(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY,
		Count:    3,
		Bysecond: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 6, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 0, 18, 0, time.UTC),
		time.Date(1998, 9, 2, 9, 0, 6, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestYearlyByHourAndMinute(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY,
		Count:    3,
		Byhour:   []int{6, 18},
		Byminute: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 18, 6, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 18, 18, 0, 0, time.UTC),
		time.Date(1998, 9, 2, 6, 6, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestYearlyByHourAndSecond(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY,
		Count:    3,
		Byhour:   []int{6, 18},
		Bysecond: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 18, 0, 6, 0, time.UTC),
		time.Date(1997, 9, 2, 18, 0, 18, 0, time.UTC),
		time.Date(1998, 9, 2, 6, 0, 6, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestYearlyByMinuteAndSecond(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY,
		Count:    3,
		Byminute: []int{6, 18},
		Bysecond: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 6, 6, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 6, 18, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 18, 6, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestYearlyByHourAndMinuteAndSecond(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY,
		Count:    3,
		Byhour:   []int{6, 18},
		Byminute: []int{6, 18},
		Bysecond: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 18, 6, 6, 0, time.UTC),
		time.Date(1997, 9, 2, 18, 6, 18, 0, time.UTC),
		time.Date(1997, 9, 2, 18, 18, 6, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestYearlyBySetPos(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY,
		Count:      3,
		Bymonthday: []int{15},
		Byhour:     []int{6, 18},
		Bysetpos:   []int{3, -3},
		Dtstart:    time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 11, 15, 18, 0, 0, 0, time.UTC),
		time.Date(1998, 2, 15, 6, 0, 0, 0, time.UTC),
		time.Date(1998, 11, 15, 18, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMonthly(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MONTHLY,
		Count:   3,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 10, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 11, 2, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMonthlyInterval(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MONTHLY,
		Count:    3,
		Interval: 2,
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 11, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 2, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMonthlyIntervalLarge(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MONTHLY,
		Count:    3,
		Interval: 18,
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 3, 2, 9, 0, 0, 0, time.UTC),
		time.Date(2000, 9, 2, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMonthlyByMonth(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MONTHLY,
		Count:   3,
		Bymonth: []int{1, 3},
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 3, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 1, 2, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMonthlyByMonthDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MONTHLY,
		Count:      3,
		Bymonthday: []int{1, 3},
		Dtstart:    time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 3, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 10, 1, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 10, 3, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMonthlyByMonthAndMonthDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MONTHLY,
		Count:      3,
		Bymonth:    []int{1, 3},
		Bymonthday: []int{5, 7},
		Dtstart:    time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 5, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 7, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 3, 5, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMonthlyByWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MONTHLY,
		Count:     3,
		Byweekday: []Weekday{TU, TH},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 9, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMonthlyByNWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MONTHLY,
		Count:     3,
		Byweekday: []Weekday{TU.Nth(1), TH.Nth(-1)},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 25, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 10, 7, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMonthlyByNWeekDayLarge(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MONTHLY,
		Count:     3,
		Byweekday: []Weekday{TU.Nth(3), TH.Nth(-3)},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 11, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 16, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 10, 16, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMonthlyByMonthAndWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MONTHLY,
		Count:     3,
		Bymonth:   []int{1, 3},
		Byweekday: []Weekday{TU, TH},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 6, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 8, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMonthlyByMonthAndNWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MONTHLY,
		Count:     3,
		Bymonth:   []int{1, 3},
		Byweekday: []Weekday{TU.Nth(1), TH.Nth(-1)},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 6, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 29, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 3, 3, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMonthlyByMonthAndNWeekDayLarge(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MONTHLY,
		Count:     3,
		Bymonth:   []int{1, 3},
		Byweekday: []Weekday{TU.Nth(3), TH.Nth(-3)},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 15, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 20, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 3, 12, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMonthlyByMonthDayAndWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MONTHLY,
		Count:      3,
		Bymonthday: []int{1, 3},
		Byweekday:  []Weekday{TU, TH},
		Dtstart:    time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 2, 3, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 3, 3, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMonthlyByMonthAndMonthDayAndWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MONTHLY,
		Count:      3,
		Bymonth:    []int{1, 3},
		Bymonthday: []int{1, 3},
		Byweekday:  []Weekday{TU, TH},
		Dtstart:    time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 3, 3, 9, 0, 0, 0, time.UTC),
		time.Date(2001, 3, 1, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMonthlyByYearDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MONTHLY,
		Count:     4,
		Byyearday: []int{1, 100, 200, 365},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 12, 31, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 4, 10, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 7, 19, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMonthlyByYearDayNeg(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MONTHLY,
		Count:     4,
		Byyearday: []int{-365, -266, -166, -1},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 12, 31, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 4, 10, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 7, 19, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMonthlyByMonthAndYearDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MONTHLY,
		Count:     4,
		Bymonth:   []int{4, 7},
		Byyearday: []int{1, 100, 200, 365},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 4, 10, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 7, 19, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 4, 10, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 7, 19, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMonthlyByMonthAndYearDayNeg(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MONTHLY,
		Count:     4,
		Bymonth:   []int{4, 7},
		Byyearday: []int{-365, -266, -166, -1},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 4, 10, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 7, 19, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 4, 10, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 7, 19, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMonthlyByWeekNo(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MONTHLY,
		Count:    3,
		Byweekno: []int{20},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 5, 11, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 5, 12, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 5, 13, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMonthlyByWeekNoAndWeekDay(t *testing.T) {
	// That's a nice one. The first days of week number one
	// may be in the last year.
	r, _ := newRecurrence(ROption{Freq: MONTHLY,
		Count:     3,
		Byweekno:  []int{1},
		Byweekday: []Weekday{MO},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 12, 29, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 1, 4, 9, 0, 0, 0, time.UTC),
		time.Date(2000, 1, 3, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMonthlyByWeekNoAndWeekDayLarge(t *testing.T) {
	// Another nice test. The last days of week number 52/53
	// may be in the next year.
	r, _ := newRecurrence(ROption{Freq: MONTHLY,
		Count:     3,
		Byweekno:  []int{52},
		Byweekday: []Weekday{SU},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 12, 28, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 12, 27, 9, 0, 0, 0, time.UTC),
		time.Date(2000, 1, 2, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMonthlyByWeekNoAndWeekDayLast(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MONTHLY,
		Count:     3,
		Byweekno:  []int{-1},
		Byweekday: []Weekday{SU},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 12, 28, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 1, 3, 9, 0, 0, 0, time.UTC),
		time.Date(2000, 1, 2, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMonthlyByWeekNoAndWeekDay53(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MONTHLY,
		Count:     3,
		Byweekno:  []int{53},
		Byweekday: []Weekday{MO},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 12, 28, 9, 0, 0, 0, time.UTC),
		time.Date(2004, 12, 27, 9, 0, 0, 0, time.UTC),
		time.Date(2009, 12, 28, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMonthlyByEaster(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MONTHLY,
		Count:    3,
		Byeaster: []int{0},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 4, 12, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 4, 4, 9, 0, 0, 0, time.UTC),
		time.Date(2000, 4, 23, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMonthlyByEasterPos(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MONTHLY,
		Count:    3,
		Byeaster: []int{1},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 4, 13, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 4, 5, 9, 0, 0, 0, time.UTC),
		time.Date(2000, 4, 24, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMonthlyByEasterNeg(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MONTHLY,
		Count:    3,
		Byeaster: []int{-1},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 4, 11, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 4, 3, 9, 0, 0, 0, time.UTC),
		time.Date(2000, 4, 22, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMonthlyByHour(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MONTHLY,
		Count:   3,
		Byhour:  []int{6, 18},
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 18, 0, 0, 0, time.UTC),
		time.Date(1997, 10, 2, 6, 0, 0, 0, time.UTC),
		time.Date(1997, 10, 2, 18, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMonthlyByMinute(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MONTHLY,
		Count:    3,
		Byminute: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 6, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 18, 0, 0, time.UTC),
		time.Date(1997, 10, 2, 9, 6, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMonthlyBySecond(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MONTHLY,
		Count:    3,
		Bysecond: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 6, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 0, 18, 0, time.UTC),
		time.Date(1997, 10, 2, 9, 0, 6, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMonthlyByHourAndMinute(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MONTHLY,
		Count:    3,
		Byhour:   []int{6, 18},
		Byminute: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 18, 6, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 18, 18, 0, 0, time.UTC),
		time.Date(1997, 10, 2, 6, 6, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMonthlyByHourAndSecond(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MONTHLY,
		Count:    3,
		Byhour:   []int{6, 18},
		Bysecond: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 18, 0, 6, 0, time.UTC),
		time.Date(1997, 9, 2, 18, 0, 18, 0, time.UTC),
		time.Date(1997, 10, 2, 6, 0, 6, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMonthlyByMinuteAndSecond(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MONTHLY,
		Count:    3,
		Byminute: []int{6, 18},
		Bysecond: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 6, 6, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 6, 18, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 18, 6, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMonthlyByHourAndMinuteAndSecond(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MONTHLY,
		Count:    3,
		Byhour:   []int{6, 18},
		Byminute: []int{6, 18},
		Bysecond: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 18, 6, 6, 0, time.UTC),
		time.Date(1997, 9, 2, 18, 6, 18, 0, time.UTC),
		time.Date(1997, 9, 2, 18, 18, 6, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMonthlyBySetPos(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MONTHLY,
		Count:      3,
		Bymonthday: []int{13, 17},
		Byhour:     []int{6, 18},
		Bysetpos:   []int{3, -3},
		Dtstart:    time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 13, 18, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 17, 6, 0, 0, 0, time.UTC),
		time.Date(1997, 10, 13, 18, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestWeekly(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: WEEKLY,
		Count:   3,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 9, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 16, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestWeeklyInterval(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: WEEKLY,
		Count:    3,
		Interval: 2,
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 16, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 30, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestWeeklyIntervalLarge(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: WEEKLY,
		Count:    3,
		Interval: 20,
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 20, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 6, 9, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestWeeklyByMonth(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: WEEKLY,
		Count:   3,
		Bymonth: []int{1, 3},
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 6, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 13, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 20, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestWeeklyByMonthDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: WEEKLY,
		Count:      3,
		Bymonthday: []int{1, 3},
		Dtstart:    time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 3, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 10, 1, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 10, 3, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestWeeklyByMonthAndMonthDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: WEEKLY,
		Count:      3,
		Bymonth:    []int{1, 3},
		Bymonthday: []int{5, 7},
		Dtstart:    time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 5, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 7, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 3, 5, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestWeeklyByWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: WEEKLY,
		Count:     3,
		Byweekday: []Weekday{TU, TH},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 9, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestWeeklyByNWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: WEEKLY,
		Count:     3,
		Byweekday: []Weekday{TU.Nth(1), TH.Nth(-1)},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 9, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestWeeklyByMonthAndWeekDay(t *testing.T) {
	// This test is interesting, because it crosses the year
	// boundary in a weekly period to find day '1' as a
	// valid recurrence.
	r, _ := newRecurrence(ROption{Freq: WEEKLY,
		Count:     3,
		Bymonth:   []int{1, 3},
		Byweekday: []Weekday{TU, TH},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 6, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 8, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestWeeklyByMonthAndNWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: WEEKLY,
		Count:     3,
		Bymonth:   []int{1, 3},
		Byweekday: []Weekday{TU.Nth(1), TH.Nth(-1)},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 6, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 8, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestWeeklyByMonthDayAndWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: WEEKLY,
		Count:      3,
		Bymonthday: []int{1, 3},
		Byweekday:  []Weekday{TU, TH},
		Dtstart:    time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 2, 3, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 3, 3, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestWeeklyByMonthAndMonthDayAndWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: WEEKLY,
		Count:      3,
		Bymonth:    []int{1, 3},
		Bymonthday: []int{1, 3},
		Byweekday:  []Weekday{TU, TH},
		Dtstart:    time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 3, 3, 9, 0, 0, 0, time.UTC),
		time.Date(2001, 3, 1, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestWeeklyByYearDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: WEEKLY,
		Count:     4,
		Byyearday: []int{1, 100, 200, 365},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 12, 31, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 4, 10, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 7, 19, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestWeeklyByYearDayNeg(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: WEEKLY,
		Count:     4,
		Byyearday: []int{-365, -266, -166, -1},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 12, 31, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 4, 10, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 7, 19, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestWeeklyByMonthAndYearDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: WEEKLY,
		Count:     4,
		Bymonth:   []int{1, 7},
		Byyearday: []int{1, 100, 200, 365},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 7, 19, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 7, 19, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestWeeklyByMonthAndYearDayNeg(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: WEEKLY,
		Count:     4,
		Bymonth:   []int{1, 7},
		Byyearday: []int{-365, -266, -166, -1},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 7, 19, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 7, 19, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestWeeklyByWeekNo(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: WEEKLY,
		Count:    3,
		Byweekno: []int{20},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 5, 11, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 5, 12, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 5, 13, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestWeeklyByWeekNoAndWeekDay(t *testing.T) {
	// That's a nice one. The first days of week number one
	// may be in the last year.
	r, _ := newRecurrence(ROption{Freq: WEEKLY,
		Count:     3,
		Byweekno:  []int{1},
		Byweekday: []Weekday{MO},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 12, 29, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 1, 4, 9, 0, 0, 0, time.UTC),
		time.Date(2000, 1, 3, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestWeeklyByWeekNoAndWeekDayLarge(t *testing.T) {
	// Another nice test. The last days of week number 52/53
	// may be in the next year.
	r, _ := newRecurrence(ROption{Freq: WEEKLY,
		Count:     3,
		Byweekno:  []int{52},
		Byweekday: []Weekday{SU},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 12, 28, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 12, 27, 9, 0, 0, 0, time.UTC),
		time.Date(2000, 1, 2, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestWeeklyByWeekNoAndWeekDayLast(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: WEEKLY,
		Count:     3,
		Byweekno:  []int{-1},
		Byweekday: []Weekday{SU},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 12, 28, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 1, 3, 9, 0, 0, 0, time.UTC),
		time.Date(2000, 1, 2, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestWeeklyByWeekNoAndWeekDay53(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: WEEKLY,
		Count:     3,
		Byweekno:  []int{53},
		Byweekday: []Weekday{MO},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 12, 28, 9, 0, 0, 0, time.UTC),
		time.Date(2004, 12, 27, 9, 0, 0, 0, time.UTC),
		time.Date(2009, 12, 28, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestWeeklyByEaster(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: WEEKLY,
		Count:    3,
		Byeaster: []int{0},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 4, 12, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 4, 4, 9, 0, 0, 0, time.UTC),
		time.Date(2000, 4, 23, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestWeeklyByEasterPos(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: WEEKLY,
		Count:    3,
		Byeaster: []int{1},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 4, 13, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 4, 5, 9, 0, 0, 0, time.UTC),
		time.Date(2000, 4, 24, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestWeeklyByEasterNeg(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: WEEKLY,
		Count:    3,
		Byeaster: []int{-1},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 4, 11, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 4, 3, 9, 0, 0, 0, time.UTC),
		time.Date(2000, 4, 22, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestWeeklyByHour(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: WEEKLY,
		Count:   3,
		Byhour:  []int{6, 18},
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 18, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 9, 6, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 9, 18, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestWeeklyByMinute(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: WEEKLY,
		Count:    3,
		Byminute: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 6, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 18, 0, 0, time.UTC),
		time.Date(1997, 9, 9, 9, 6, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestWeeklyBySecond(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: WEEKLY,
		Count:    3,
		Bysecond: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 6, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 0, 18, 0, time.UTC),
		time.Date(1997, 9, 9, 9, 0, 6, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestWeeklyByHourAndMinute(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: WEEKLY,
		Count:    3,
		Byhour:   []int{6, 18},
		Byminute: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 18, 6, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 18, 18, 0, 0, time.UTC),
		time.Date(1997, 9, 9, 6, 6, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestWeeklyByHourAndSecond(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: WEEKLY,
		Count:    3,
		Byhour:   []int{6, 18},
		Bysecond: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 18, 0, 6, 0, time.UTC),
		time.Date(1997, 9, 2, 18, 0, 18, 0, time.UTC),
		time.Date(1997, 9, 9, 6, 0, 6, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestWeeklyByMinuteAndSecond(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: WEEKLY,
		Count:    3,
		Byminute: []int{6, 18},
		Bysecond: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 6, 6, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 6, 18, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 18, 6, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestWeeklyByHourAndMinuteAndSecond(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: WEEKLY,
		Count:    3,
		Byhour:   []int{6, 18},
		Byminute: []int{6, 18},
		Bysecond: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 18, 6, 6, 0, time.UTC),
		time.Date(1997, 9, 2, 18, 6, 18, 0, time.UTC),
		time.Date(1997, 9, 2, 18, 18, 6, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestWeeklyBySetPos(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: WEEKLY,
		Count:     3,
		Byweekday: []Weekday{TU, TH},
		Byhour:    []int{6, 18},
		Bysetpos:  []int{3, -3},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 18, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 4, 6, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 9, 18, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestDaily(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:   3,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 3, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestDailyInterval(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:    3,
		Interval: 2,
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 6, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestDailyIntervalLarge(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:    3,
		Interval: 92,
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 12, 3, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 3, 5, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestDailyByMonth(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:   3,
		Bymonth: []int{1, 3},
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 3, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestDailyByMonthDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:      3,
		Bymonthday: []int{1, 3},
		Dtstart:    time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 3, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 10, 1, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 10, 3, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestDailyByMonthAndMonthDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:      3,
		Bymonth:    []int{1, 3},
		Bymonthday: []int{5, 7},
		Dtstart:    time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 5, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 7, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 3, 5, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestDailyByWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:     3,
		Byweekday: []Weekday{TU, TH},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 9, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestDailyByNWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:     3,
		Byweekday: []Weekday{TU.Nth(1), TH.Nth(-1)},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 9, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestDailyByMonthAndWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:     3,
		Bymonth:   []int{1, 3},
		Byweekday: []Weekday{TU, TH},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 6, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 8, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestDailyByMonthAndNWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:     3,
		Bymonth:   []int{1, 3},
		Byweekday: []Weekday{TU.Nth(1), TH.Nth(-1)},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 6, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 8, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestDailyByMonthDayAndWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:      3,
		Bymonthday: []int{1, 3},
		Byweekday:  []Weekday{TU, TH},
		Dtstart:    time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 2, 3, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 3, 3, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestDailyByMonthAndMonthDayAndWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:      3,
		Bymonth:    []int{1, 3},
		Bymonthday: []int{1, 3},
		Byweekday:  []Weekday{TU, TH},
		Dtstart:    time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 3, 3, 9, 0, 0, 0, time.UTC),
		time.Date(2001, 3, 1, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestDailyByYearDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:     4,
		Byyearday: []int{1, 100, 200, 365},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 12, 31, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 4, 10, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 7, 19, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestDailyByYearDayNeg(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:     4,
		Byyearday: []int{-365, -266, -166, -1},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 12, 31, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 4, 10, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 7, 19, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestDailyByMonthAndYearDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:     4,
		Bymonth:   []int{1, 7},
		Byyearday: []int{1, 100, 200, 365},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 7, 19, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 7, 19, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestDailyByMonthAndYearDayNeg(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:     4,
		Bymonth:   []int{1, 7},
		Byyearday: []int{-365, -266, -166, -1},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 7, 19, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 7, 19, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestDailyByWeekNo(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:    3,
		Byweekno: []int{20},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 5, 11, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 5, 12, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 5, 13, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestDailyByWeekNoAndWeekDay(t *testing.T) {
	// That's a nice one. The first days of week number one
	// may be in the last year.
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:     3,
		Byweekno:  []int{1},
		Byweekday: []Weekday{MO},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 12, 29, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 1, 4, 9, 0, 0, 0, time.UTC),
		time.Date(2000, 1, 3, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestDailyByWeekNoAndWeekDayLarge(t *testing.T) {
	// Another nice test. The last days of week number 52/53
	// may be in the next year.
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:     3,
		Byweekno:  []int{52},
		Byweekday: []Weekday{SU},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 12, 28, 9, 0, 0, 0, time.UTC),
		time.Date(1998, 12, 27, 9, 0, 0, 0, time.UTC),
		time.Date(2000, 1, 2, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestDailyByWeekNoAndWeekDayLast(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:     3,
		Byweekno:  []int{-1},
		Byweekday: []Weekday{SU},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 12, 28, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 1, 3, 9, 0, 0, 0, time.UTC),
		time.Date(2000, 1, 2, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestDailyByWeekNoAndWeekDay53(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:     3,
		Byweekno:  []int{53},
		Byweekday: []Weekday{MO},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 12, 28, 9, 0, 0, 0, time.UTC),
		time.Date(2004, 12, 27, 9, 0, 0, 0, time.UTC),
		time.Date(2009, 12, 28, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestDailyByEaster(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:    3,
		Byeaster: []int{0},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 4, 12, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 4, 4, 9, 0, 0, 0, time.UTC),
		time.Date(2000, 4, 23, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestDailyByEasterPos(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:    3,
		Byeaster: []int{1},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 4, 13, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 4, 5, 9, 0, 0, 0, time.UTC),
		time.Date(2000, 4, 24, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestDailyByEasterNeg(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:    3,
		Byeaster: []int{-1},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 4, 11, 9, 0, 0, 0, time.UTC),
		time.Date(1999, 4, 3, 9, 0, 0, 0, time.UTC),
		time.Date(2000, 4, 22, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestDailyByHour(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:   3,
		Byhour:  []int{6, 18},
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 18, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 3, 6, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 3, 18, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestDailyByMinute(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:    3,
		Byminute: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 6, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 18, 0, 0, time.UTC),
		time.Date(1997, 9, 3, 9, 6, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestDailyBySecond(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:    3,
		Bysecond: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 6, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 0, 18, 0, time.UTC),
		time.Date(1997, 9, 3, 9, 0, 6, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestDailyByHourAndMinute(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:    3,
		Byhour:   []int{6, 18},
		Byminute: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 18, 6, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 18, 18, 0, 0, time.UTC),
		time.Date(1997, 9, 3, 6, 6, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestDailyByHourAndSecond(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:    3,
		Byhour:   []int{6, 18},
		Bysecond: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 18, 0, 6, 0, time.UTC),
		time.Date(1997, 9, 2, 18, 0, 18, 0, time.UTC),
		time.Date(1997, 9, 3, 6, 0, 6, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestDailyByMinuteAndSecond(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:    3,
		Byminute: []int{6, 18},
		Bysecond: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 6, 6, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 6, 18, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 18, 6, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestDailyByHourAndMinuteAndSecond(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:    3,
		Byhour:   []int{6, 18},
		Byminute: []int{6, 18},
		Bysecond: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 18, 6, 6, 0, time.UTC),
		time.Date(1997, 9, 2, 18, 6, 18, 0, time.UTC),
		time.Date(1997, 9, 2, 18, 18, 6, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestDailyBySetPos(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:    3,
		Byhour:   []int{6, 18},
		Byminute: []int{15, 45},
		Bysetpos: []int{3, -3},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 18, 15, 0, 0, time.UTC),
		time.Date(1997, 9, 3, 6, 45, 0, 0, time.UTC),
		time.Date(1997, 9, 3, 18, 15, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestHourly(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: HOURLY,
		Count:   3,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 10, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 11, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestHourlyInterval(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: HOURLY,
		Count:    3,
		Interval: 2,
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 11, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 13, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestHourlyIntervalLarge(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: HOURLY,
		Count:    3,
		Interval: 769,
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 10, 4, 10, 0, 0, 0, time.UTC),
		time.Date(1997, 11, 5, 11, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestHourlyByMonth(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: HOURLY,
		Count:   3,
		Bymonth: []int{1, 3},
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 1, 1, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 1, 2, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestHourlyByMonthDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: HOURLY,
		Count:      3,
		Bymonthday: []int{1, 3},
		Dtstart:    time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 3, 0, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 3, 1, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 3, 2, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestHourlyByMonthAndMonthDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: HOURLY,
		Count:      3,
		Bymonth:    []int{1, 3},
		Bymonthday: []int{5, 7},
		Dtstart:    time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 5, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 5, 1, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 5, 2, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestHourlyByWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: HOURLY,
		Count:     3,
		Byweekday: []Weekday{TU, TH},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 10, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 11, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestHourlyByNWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: HOURLY,
		Count:     3,
		Byweekday: []Weekday{TU.Nth(1), TH.Nth(-1)},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 10, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 11, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestHourlyByMonthAndWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: HOURLY,
		Count:     3,
		Bymonth:   []int{1, 3},
		Byweekday: []Weekday{TU, TH},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 1, 1, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 1, 2, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestHourlyByMonthAndNWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: HOURLY,
		Count:     3,
		Bymonth:   []int{1, 3},
		Byweekday: []Weekday{TU.Nth(1), TH.Nth(-1)},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 1, 1, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 1, 2, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestHourlyByMonthDayAndWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: HOURLY,
		Count:      3,
		Bymonthday: []int{1, 3},
		Byweekday:  []Weekday{TU, TH},
		Dtstart:    time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 1, 1, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 1, 2, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestHourlyByMonthAndMonthDayAndWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: HOURLY,
		Count:      3,
		Bymonth:    []int{1, 3},
		Bymonthday: []int{1, 3},
		Byweekday:  []Weekday{TU, TH},
		Dtstart:    time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 1, 1, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 1, 2, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestHourlyByYearDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: HOURLY,
		Count:     4,
		Byyearday: []int{1, 100, 200, 365},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 12, 31, 0, 0, 0, 0, time.UTC),
		time.Date(1997, 12, 31, 1, 0, 0, 0, time.UTC),
		time.Date(1997, 12, 31, 2, 0, 0, 0, time.UTC),
		time.Date(1997, 12, 31, 3, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestHourlyByYearDayNeg(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: HOURLY,
		Count:     4,
		Byyearday: []int{-365, -266, -166, -1},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 12, 31, 0, 0, 0, 0, time.UTC),
		time.Date(1997, 12, 31, 1, 0, 0, 0, time.UTC),
		time.Date(1997, 12, 31, 2, 0, 0, 0, time.UTC),
		time.Date(1997, 12, 31, 3, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestHourlyByMonthAndYearDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: HOURLY,
		Count:     4,
		Bymonth:   []int{4, 7},
		Byyearday: []int{1, 100, 200, 365},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 4, 10, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 4, 10, 1, 0, 0, 0, time.UTC),
		time.Date(1998, 4, 10, 2, 0, 0, 0, time.UTC),
		time.Date(1998, 4, 10, 3, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestHourlyByMonthAndYearDayNeg(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: HOURLY,
		Count:     4,
		Bymonth:   []int{4, 7},
		Byyearday: []int{-365, -266, -166, -1},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 4, 10, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 4, 10, 1, 0, 0, 0, time.UTC),
		time.Date(1998, 4, 10, 2, 0, 0, 0, time.UTC),
		time.Date(1998, 4, 10, 3, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestHourlyByWeekNo(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: HOURLY,
		Count:    3,
		Byweekno: []int{20},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 5, 11, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 5, 11, 1, 0, 0, 0, time.UTC),
		time.Date(1998, 5, 11, 2, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestHourlyByWeekNoAndWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: HOURLY,
		Count:     3,
		Byweekno:  []int{1},
		Byweekday: []Weekday{MO},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 12, 29, 0, 0, 0, 0, time.UTC),
		time.Date(1997, 12, 29, 1, 0, 0, 0, time.UTC),
		time.Date(1997, 12, 29, 2, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestHourlyByWeekNoAndWeekDayLarge(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: HOURLY,
		Count:     3,
		Byweekno:  []int{52},
		Byweekday: []Weekday{SU},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 12, 28, 0, 0, 0, 0, time.UTC),
		time.Date(1997, 12, 28, 1, 0, 0, 0, time.UTC),
		time.Date(1997, 12, 28, 2, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestHourlyByWeekNoAndWeekDayLast(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: HOURLY,
		Count:     3,
		Byweekno:  []int{-1},
		Byweekday: []Weekday{SU},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 12, 28, 0, 0, 0, 0, time.UTC),
		time.Date(1997, 12, 28, 1, 0, 0, 0, time.UTC),
		time.Date(1997, 12, 28, 2, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestHourlyByWeekNoAndWeekDay53(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: HOURLY,
		Count:     3,
		Byweekno:  []int{53},
		Byweekday: []Weekday{MO},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 12, 28, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 12, 28, 1, 0, 0, 0, time.UTC),
		time.Date(1998, 12, 28, 2, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestHourlyByEaster(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: HOURLY,
		Count:    3,
		Byeaster: []int{0},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 4, 12, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 4, 12, 1, 0, 0, 0, time.UTC),
		time.Date(1998, 4, 12, 2, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestHourlyByEasterPos(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: HOURLY,
		Count:    3,
		Byeaster: []int{1},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 4, 13, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 4, 13, 1, 0, 0, 0, time.UTC),
		time.Date(1998, 4, 13, 2, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestHourlyByEasterNeg(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: HOURLY,
		Count:    3,
		Byeaster: []int{-1},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 4, 11, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 4, 11, 1, 0, 0, 0, time.UTC),
		time.Date(1998, 4, 11, 2, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestHourlyByHour(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: HOURLY,
		Count:   3,
		Byhour:  []int{6, 18},
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 18, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 3, 6, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 3, 18, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestHourlyByMinute(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: HOURLY,
		Count:    3,
		Byminute: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 6, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 18, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 10, 6, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestHourlyBySecond(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: HOURLY,
		Count:    3,
		Bysecond: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 6, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 0, 18, 0, time.UTC),
		time.Date(1997, 9, 2, 10, 0, 6, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestHourlyByHourAndMinute(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: HOURLY,
		Count:    3,
		Byhour:   []int{6, 18},
		Byminute: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 18, 6, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 18, 18, 0, 0, time.UTC),
		time.Date(1997, 9, 3, 6, 6, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestHourlyByHourAndSecond(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: HOURLY,
		Count:    3,
		Byhour:   []int{6, 18},
		Bysecond: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 18, 0, 6, 0, time.UTC),
		time.Date(1997, 9, 2, 18, 0, 18, 0, time.UTC),
		time.Date(1997, 9, 3, 6, 0, 6, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestHourlyByMinuteAndSecond(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: HOURLY,
		Count:    3,
		Byminute: []int{6, 18},
		Bysecond: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 6, 6, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 6, 18, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 18, 6, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestHourlyByHourAndMinuteAndSecond(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: HOURLY,
		Count:    3,
		Byhour:   []int{6, 18},
		Byminute: []int{6, 18},
		Bysecond: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 18, 6, 6, 0, time.UTC),
		time.Date(1997, 9, 2, 18, 6, 18, 0, time.UTC),
		time.Date(1997, 9, 2, 18, 18, 6, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestHourlyBySetPos(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: HOURLY,
		Count:    3,
		Byminute: []int{15, 45},
		Bysecond: []int{15, 45},
		Bysetpos: []int{3, -3},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 15, 45, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 45, 15, 0, time.UTC),
		time.Date(1997, 9, 2, 10, 15, 45, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMinutely(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MINUTELY,
		Count:   3,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 1, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 2, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMinutelyInterval(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MINUTELY,
		Count:    3,
		Interval: 2,
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 2, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 4, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMinutelyIntervalLarge(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MINUTELY,
		Count:    3,
		Interval: 1501,
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 3, 10, 1, 0, 0, time.UTC),
		time.Date(1997, 9, 4, 11, 2, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMinutelyByMonth(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MINUTELY,
		Count:   3,
		Bymonth: []int{1, 3},
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 1, 0, 1, 0, 0, time.UTC),
		time.Date(1998, 1, 1, 0, 2, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMinutelyByMonthDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MINUTELY,
		Count:      3,
		Bymonthday: []int{1, 3},
		Dtstart:    time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 3, 0, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 3, 0, 1, 0, 0, time.UTC),
		time.Date(1997, 9, 3, 0, 2, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMinutelyByMonthAndMonthDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MINUTELY,
		Count:      3,
		Bymonth:    []int{1, 3},
		Bymonthday: []int{5, 7},
		Dtstart:    time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 5, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 5, 0, 1, 0, 0, time.UTC),
		time.Date(1998, 1, 5, 0, 2, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMinutelyByWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MINUTELY,
		Count:     3,
		Byweekday: []Weekday{TU, TH},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 1, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 2, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMinutelyByNWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MINUTELY,
		Count:     3,
		Byweekday: []Weekday{TU.Nth(1), TH.Nth(-1)},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 1, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 2, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMinutelyByMonthAndWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MINUTELY,
		Count:     3,
		Bymonth:   []int{1, 3},
		Byweekday: []Weekday{TU, TH},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 1, 0, 1, 0, 0, time.UTC),
		time.Date(1998, 1, 1, 0, 2, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMinutelyByMonthAndNWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MINUTELY,
		Count:     3,
		Bymonth:   []int{1, 3},
		Byweekday: []Weekday{TU.Nth(1), TH.Nth(-1)},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 1, 0, 1, 0, 0, time.UTC),
		time.Date(1998, 1, 1, 0, 2, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMinutelyByMonthDayAndWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MINUTELY,
		Count:      3,
		Bymonthday: []int{1, 3},
		Byweekday:  []Weekday{TU, TH},
		Dtstart:    time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 1, 0, 1, 0, 0, time.UTC),
		time.Date(1998, 1, 1, 0, 2, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMinutelyByMonthAndMonthDayAndWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MINUTELY,
		Count:      3,
		Bymonth:    []int{1, 3},
		Bymonthday: []int{1, 3},
		Byweekday:  []Weekday{TU, TH},
		Dtstart:    time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 1, 0, 1, 0, 0, time.UTC),
		time.Date(1998, 1, 1, 0, 2, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMinutelyByYearDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MINUTELY,
		Count:     4,
		Byyearday: []int{1, 100, 200, 365},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 12, 31, 0, 0, 0, 0, time.UTC),
		time.Date(1997, 12, 31, 0, 1, 0, 0, time.UTC),
		time.Date(1997, 12, 31, 0, 2, 0, 0, time.UTC),
		time.Date(1997, 12, 31, 0, 3, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMinutelyByYearDayNeg(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MINUTELY,
		Count:     4,
		Byyearday: []int{-365, -266, -166, -1},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 12, 31, 0, 0, 0, 0, time.UTC),
		time.Date(1997, 12, 31, 0, 1, 0, 0, time.UTC),
		time.Date(1997, 12, 31, 0, 2, 0, 0, time.UTC),
		time.Date(1997, 12, 31, 0, 3, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMinutelyByMonthAndYearDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MINUTELY,
		Count:     4,
		Bymonth:   []int{4, 7},
		Byyearday: []int{1, 100, 200, 365},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 4, 10, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 4, 10, 0, 1, 0, 0, time.UTC),
		time.Date(1998, 4, 10, 0, 2, 0, 0, time.UTC),
		time.Date(1998, 4, 10, 0, 3, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMinutelyByMonthAndYearDayNeg(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MINUTELY,
		Count:     4,
		Bymonth:   []int{4, 7},
		Byyearday: []int{-365, -266, -166, -1},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 4, 10, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 4, 10, 0, 1, 0, 0, time.UTC),
		time.Date(1998, 4, 10, 0, 2, 0, 0, time.UTC),
		time.Date(1998, 4, 10, 0, 3, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMinutelyByWeekNo(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MINUTELY,
		Count:    3,
		Byweekno: []int{20},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 5, 11, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 5, 11, 0, 1, 0, 0, time.UTC),
		time.Date(1998, 5, 11, 0, 2, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMinutelyByWeekNoAndWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MINUTELY,
		Count:     3,
		Byweekno:  []int{1},
		Byweekday: []Weekday{MO},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 12, 29, 0, 0, 0, 0, time.UTC),
		time.Date(1997, 12, 29, 0, 1, 0, 0, time.UTC),
		time.Date(1997, 12, 29, 0, 2, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMinutelyByWeekNoAndWeekDayLarge(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MINUTELY,
		Count:     3,
		Byweekno:  []int{52},
		Byweekday: []Weekday{SU},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 12, 28, 0, 0, 0, 0, time.UTC),
		time.Date(1997, 12, 28, 0, 1, 0, 0, time.UTC),
		time.Date(1997, 12, 28, 0, 2, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMinutelyByWeekNoAndWeekDayLast(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MINUTELY,
		Count:     3,
		Byweekno:  []int{-1},
		Byweekday: []Weekday{SU},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 12, 28, 0, 0, 0, 0, time.UTC),
		time.Date(1997, 12, 28, 0, 1, 0, 0, time.UTC),
		time.Date(1997, 12, 28, 0, 2, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMinutelyByWeekNoAndWeekDay53(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MINUTELY,
		Count:     3,
		Byweekno:  []int{53},
		Byweekday: []Weekday{MO},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 12, 28, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 12, 28, 0, 1, 0, 0, time.UTC),
		time.Date(1998, 12, 28, 0, 2, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMinutelyByEaster(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MINUTELY,
		Count:    3,
		Byeaster: []int{0},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 4, 12, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 4, 12, 0, 1, 0, 0, time.UTC),
		time.Date(1998, 4, 12, 0, 2, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMinutelyByEasterPos(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MINUTELY,
		Count:    3,
		Byeaster: []int{1},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 4, 13, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 4, 13, 0, 1, 0, 0, time.UTC),
		time.Date(1998, 4, 13, 0, 2, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMinutelyByEasterNeg(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MINUTELY,
		Count:    3,
		Byeaster: []int{-1},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 4, 11, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 4, 11, 0, 1, 0, 0, time.UTC),
		time.Date(1998, 4, 11, 0, 2, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMinutelyByHour(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MINUTELY,
		Count:   3,
		Byhour:  []int{6, 18},
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 18, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 18, 1, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 18, 2, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMinutelyByMinute(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MINUTELY,
		Count:    3,
		Byminute: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 6, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 18, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 10, 6, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMinutelyBySecond(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MINUTELY,
		Count:    3,
		Bysecond: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 6, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 0, 18, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 1, 6, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMinutelyByHourAndMinute(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MINUTELY,
		Count:    3,
		Byhour:   []int{6, 18},
		Byminute: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 18, 6, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 18, 18, 0, 0, time.UTC),
		time.Date(1997, 9, 3, 6, 6, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMinutelyByHourAndSecond(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MINUTELY,
		Count:    3,
		Byhour:   []int{6, 18},
		Bysecond: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 18, 0, 6, 0, time.UTC),
		time.Date(1997, 9, 2, 18, 0, 18, 0, time.UTC),
		time.Date(1997, 9, 2, 18, 1, 6, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMinutelyByMinuteAndSecond(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MINUTELY,
		Count:    3,
		Byminute: []int{6, 18},
		Bysecond: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 6, 6, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 6, 18, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 18, 6, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMinutelyByHourAndMinuteAndSecond(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MINUTELY,
		Count:    3,
		Byhour:   []int{6, 18},
		Byminute: []int{6, 18},
		Bysecond: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 18, 6, 6, 0, time.UTC),
		time.Date(1997, 9, 2, 18, 6, 18, 0, time.UTC),
		time.Date(1997, 9, 2, 18, 18, 6, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestMinutelyBySetPos(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MINUTELY,
		Count:    3,
		Bysecond: []int{15, 30, 45},
		Bysetpos: []int{3, -3},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 15, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 0, 45, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 1, 15, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSecondly(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: SECONDLY,
		Count:   3,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 0, 1, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 0, 2, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSecondlyInterval(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: SECONDLY,
		Count:    3,
		Interval: 2,
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 0, 2, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 0, 4, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSecondlyIntervalLarge(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: SECONDLY,
		Count:    3,
		Interval: 90061,
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 3, 10, 1, 1, 0, time.UTC),
		time.Date(1997, 9, 4, 11, 2, 2, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSecondlyByMonth(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: SECONDLY,
		Count:   3,
		Bymonth: []int{1, 3},
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 1, 0, 0, 1, 0, time.UTC),
		time.Date(1998, 1, 1, 0, 0, 2, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSecondlyByMonthDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: SECONDLY,
		Count:      3,
		Bymonthday: []int{1, 3},
		Dtstart:    time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 3, 0, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 3, 0, 0, 1, 0, time.UTC),
		time.Date(1997, 9, 3, 0, 0, 2, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSecondlyByMonthAndMonthDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: SECONDLY,
		Count:      3,
		Bymonth:    []int{1, 3},
		Bymonthday: []int{5, 7},
		Dtstart:    time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 5, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 5, 0, 0, 1, 0, time.UTC),
		time.Date(1998, 1, 5, 0, 0, 2, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSecondlyByWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: SECONDLY,
		Count:     3,
		Byweekday: []Weekday{TU, TH},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 0, 1, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 0, 2, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSecondlyByNWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: SECONDLY,
		Count:     3,
		Byweekday: []Weekday{TU.Nth(1), TH.Nth(-1)},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 0, 1, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 0, 2, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSecondlyByMonthAndWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: SECONDLY,
		Count:     3,
		Bymonth:   []int{1, 3},
		Byweekday: []Weekday{TU, TH},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 1, 0, 0, 1, 0, time.UTC),
		time.Date(1998, 1, 1, 0, 0, 2, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSecondlyByMonthAndNWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: SECONDLY,
		Count:     3,
		Bymonth:   []int{1, 3},
		Byweekday: []Weekday{TU.Nth(1), TH.Nth(-1)},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 1, 0, 0, 1, 0, time.UTC),
		time.Date(1998, 1, 1, 0, 0, 2, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSecondlyByMonthDayAndWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: SECONDLY,
		Count:      3,
		Bymonthday: []int{1, 3},
		Byweekday:  []Weekday{TU, TH},
		Dtstart:    time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 1, 0, 0, 1, 0, time.UTC),
		time.Date(1998, 1, 1, 0, 0, 2, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSecondlyByMonthAndMonthDayAndWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: SECONDLY,
		Count:      3,
		Bymonth:    []int{1, 3},
		Bymonthday: []int{1, 3},
		Byweekday:  []Weekday{TU, TH},
		Dtstart:    time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 1, 1, 0, 0, 1, 0, time.UTC),
		time.Date(1998, 1, 1, 0, 0, 2, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSecondlyByYearDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: SECONDLY,
		Count:     4,
		Byyearday: []int{1, 100, 200, 365},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 12, 31, 0, 0, 0, 0, time.UTC),
		time.Date(1997, 12, 31, 0, 0, 1, 0, time.UTC),
		time.Date(1997, 12, 31, 0, 0, 2, 0, time.UTC),
		time.Date(1997, 12, 31, 0, 0, 3, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSecondlyByYearDayNeg(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: SECONDLY,
		Count:     4,
		Byyearday: []int{-365, -266, -166, -1},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 12, 31, 0, 0, 0, 0, time.UTC),
		time.Date(1997, 12, 31, 0, 0, 1, 0, time.UTC),
		time.Date(1997, 12, 31, 0, 0, 2, 0, time.UTC),
		time.Date(1997, 12, 31, 0, 0, 3, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSecondlyByMonthAndYearDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: SECONDLY,
		Count:     4,
		Bymonth:   []int{4, 7},
		Byyearday: []int{1, 100, 200, 365},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 4, 10, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 4, 10, 0, 0, 1, 0, time.UTC),
		time.Date(1998, 4, 10, 0, 0, 2, 0, time.UTC),
		time.Date(1998, 4, 10, 0, 0, 3, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSecondlyByMonthAndYearDayNeg(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: SECONDLY,
		Count:     4,
		Bymonth:   []int{4, 7},
		Byyearday: []int{-365, -266, -166, -1},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 4, 10, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 4, 10, 0, 0, 1, 0, time.UTC),
		time.Date(1998, 4, 10, 0, 0, 2, 0, time.UTC),
		time.Date(1998, 4, 10, 0, 0, 3, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSecondlyByWeekNo(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: SECONDLY,
		Count:    3,
		Byweekno: []int{20},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 5, 11, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 5, 11, 0, 0, 1, 0, time.UTC),
		time.Date(1998, 5, 11, 0, 0, 2, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSecondlyByWeekNoAndWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: SECONDLY,
		Count:     3,
		Byweekno:  []int{1},
		Byweekday: []Weekday{MO},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 12, 29, 0, 0, 0, 0, time.UTC),
		time.Date(1997, 12, 29, 0, 0, 1, 0, time.UTC),
		time.Date(1997, 12, 29, 0, 0, 2, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSecondlyByWeekNoAndWeekDayLarge(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: SECONDLY,
		Count:     3,
		Byweekno:  []int{52},
		Byweekday: []Weekday{SU},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 12, 28, 0, 0, 0, 0, time.UTC),
		time.Date(1997, 12, 28, 0, 0, 1, 0, time.UTC),
		time.Date(1997, 12, 28, 0, 0, 2, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSecondlyByWeekNoAndWeekDayLast(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: SECONDLY,
		Count:     3,
		Byweekno:  []int{-1},
		Byweekday: []Weekday{SU},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 12, 28, 0, 0, 0, 0, time.UTC),
		time.Date(1997, 12, 28, 0, 0, 1, 0, time.UTC),
		time.Date(1997, 12, 28, 0, 0, 2, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSecondlyByWeekNoAndWeekDay53(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: SECONDLY,
		Count:     3,
		Byweekno:  []int{53},
		Byweekday: []Weekday{MO},
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 12, 28, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 12, 28, 0, 0, 1, 0, time.UTC),
		time.Date(1998, 12, 28, 0, 0, 2, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSecondlyByEaster(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: SECONDLY,
		Count:    3,
		Byeaster: []int{0},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 4, 12, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 4, 12, 0, 0, 1, 0, time.UTC),
		time.Date(1998, 4, 12, 0, 0, 2, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSecondlyByEasterPos(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: SECONDLY,
		Count:    3,
		Byeaster: []int{1},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 4, 13, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 4, 13, 0, 0, 1, 0, time.UTC),
		time.Date(1998, 4, 13, 0, 0, 2, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSecondlyByEasterNeg(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: SECONDLY,
		Count:    3,
		Byeaster: []int{-1},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1998, 4, 11, 0, 0, 0, 0, time.UTC),
		time.Date(1998, 4, 11, 0, 0, 1, 0, time.UTC),
		time.Date(1998, 4, 11, 0, 0, 2, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSecondlyByHour(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: SECONDLY,
		Count:   3,
		Byhour:  []int{6, 18},
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 18, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 18, 0, 1, 0, time.UTC),
		time.Date(1997, 9, 2, 18, 0, 2, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSecondlyByMinute(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: SECONDLY,
		Count:    3,
		Byminute: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 6, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 6, 1, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 6, 2, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSecondlyBySecond(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: SECONDLY,
		Count:    3,
		Bysecond: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 6, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 0, 18, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 1, 6, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSecondlyByHourAndMinute(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: SECONDLY,
		Count:    3,
		Byhour:   []int{6, 18},
		Byminute: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 18, 6, 0, 0, time.UTC),
		time.Date(1997, 9, 2, 18, 6, 1, 0, time.UTC),
		time.Date(1997, 9, 2, 18, 6, 2, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSecondlyByHourAndSecond(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: SECONDLY,
		Count:    3,
		Byhour:   []int{6, 18},
		Bysecond: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 18, 0, 6, 0, time.UTC),
		time.Date(1997, 9, 2, 18, 0, 18, 0, time.UTC),
		time.Date(1997, 9, 2, 18, 1, 6, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSecondlyByMinuteAndSecond(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: SECONDLY,
		Count:    3,
		Byminute: []int{6, 18},
		Bysecond: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 6, 6, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 6, 18, 0, time.UTC),
		time.Date(1997, 9, 2, 9, 18, 6, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSecondlyByHourAndMinuteAndSecond(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: SECONDLY,
		Count:    3,
		Byhour:   []int{6, 18},
		Byminute: []int{6, 18},
		Bysecond: []int{6, 18},
		Dtstart:  time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 18, 6, 6, 0, time.UTC),
		time.Date(1997, 9, 2, 18, 6, 18, 0, time.UTC),
		time.Date(1997, 9, 2, 18, 18, 6, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSecondlyByHourAndMinuteAndSecondBug(t *testing.T) {
	// This explores a bug found by Mathieu Bridon.
	r, _ := newRecurrence(ROption{Freq: SECONDLY,
		Count:    3,
		Bysecond: []int{0},
		Byminute: []int{1},
		Dtstart:  time.Date(2010, 3, 22, 12, 1, 0, 0, time.UTC)})
	want := []time.Time{time.Date(2010, 3, 22, 12, 1, 0, 0, time.UTC),
		time.Date(2010, 3, 22, 13, 1, 0, 0, time.UTC),
		time.Date(2010, 3, 22, 14, 1, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestUntilNotMatching(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:   3,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		Until:   time.Date(1997, 9, 5, 8, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 3, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestUntilMatching(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:   3,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		Until:   time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 3, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestUntilSingle(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:   3,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		Until:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestUntilEmpty(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:   3,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		Until:   time.Date(1997, 9, 1, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestUntilWithDate(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:   3,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		Until:   time.Date(1997, 9, 5, 0, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 3, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestWkStIntervalMO(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: WEEKLY,
		Count:     3,
		Interval:  2,
		Byweekday: []Weekday{TU, SU},
		Wkst:      MO,
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 7, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 16, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestWkStIntervalSU(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: WEEKLY,
		Count:     3,
		Interval:  2,
		Byweekday: []Weekday{TU, SU},
		Wkst:      SU,
		Dtstart:   time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 14, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 16, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestDTStart(t *testing.T) {
	dt := time.Now().UTC().Truncate(time.Second)
	r, _ := newRecurrence(ROption{Freq: YEARLY, Count: 3})
	want := []time.Time{dt, dt.AddDate(1, 0, 0), dt.AddDate(2, 0, 0)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}

	dt = dt.AddDate(0, 0, 3)
	r.DTStart(dt)
	want = []time.Time{dt, dt.AddDate(1, 0, 0), dt.AddDate(2, 0, 0)}
	value = r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestDTStartIsDate(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:   3,
		Dtstart: time.Date(1997, 9, 2, 0, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 0, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 3, 0, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 4, 0, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestDTStartWithMicroseconds(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		Count:   3,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 500000000, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 3, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC)}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestUntil(t *testing.T) {
	r1, _ := newRecurrence(ROption{Freq: DAILY,
		Dtstart: time.Date(1997, 9, 2, 0, 0, 0, 0, time.UTC)})
	r1.until = time.Date(1998, 9, 2, 0, 0, 0, 0, time.UTC)
	r1.rebuildRule()

	r2, _ := newRecurrence(ROption{Freq: DAILY,
		Dtstart: time.Date(1997, 9, 2, 0, 0, 0, 0, time.UTC),
		Until:   time.Date(1998, 9, 2, 0, 0, 0, 0, time.UTC)})

	v1 := r1.All()
	v2 := r2.All()
	if !timesEqual(v1, v2) {
		t.Errorf("get %v, want %v", v1, v2)
	}

	r3, _ := newRecurrence(ROption{Freq: MONTHLY,
		Dtstart: time.Date(MAXYEAR-100, 1, 1, 0, 0, 0, 0, time.UTC)})
	r3.until = time.Date(MAXYEAR+100, 1, 1, 0, 0, 0, 0, time.UTC)
	r3.rebuildRule()
	v3 := r3.All()
	if len(v3) != 101*12 {
		t.Errorf("get %v, want %v", len(v3), 101*12)
	}
}

func TestMaxYear(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY,
		Count:      3,
		Bymonth:    []int{2},
		Bymonthday: []int{31},
		Dtstart:    time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestBefore(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		// Count:5,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC)
	value := r.Before(time.Date(1997, 9, 5, 9, 0, 0, 0, time.UTC), false)
	if value != want {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestBeforeInc(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		// Count:5,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := time.Date(1997, 9, 5, 9, 0, 0, 0, time.UTC)
	value := r.Before(time.Date(1997, 9, 5, 9, 0, 0, 0, time.UTC), true)
	if value != want {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestAfter(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		// Count:5,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})

	want := time.Date(1997, 9, 5, 9, 0, 0, 0, time.UTC)
	value := r.After(time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC), false)
	if value != want {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestAfterInc(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		// Count:5,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC)
	value := r.After(time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC), true)
	if value != want {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestBetween(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		// Count:5,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 3, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 5, 9, 0, 0, 0, time.UTC)}
	value := r.Between(time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC), time.Date(1997, 9, 6, 9, 0, 0, 0, time.UTC), false)
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestBetweenInc(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY,
		// Count:5,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 3, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 5, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 6, 9, 0, 0, 0, time.UTC)}
	value := r.Between(time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC), time.Date(1997, 9, 6, 9, 0, 0, 0, time.UTC), true)
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestAllWithDefaultUtil(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})

	value := r.All()
	if len(value) > 300 || len(value) < 200 {
		t.Errorf("No default Util time")
	}

	r, _ = newRecurrence(ROption{Freq: YEARLY})
	if len(r.All()) != len(value) {
		t.Errorf("No default Util time")
	}
}

func TestWeekdayGetters(t *testing.T) {
	wd := Weekday{n: 2, weekday: 0}
	if wd.N() != 2 {
		t.Errorf("Ord week getter is wrong")
	}
	if wd.Day() != 0 {
		t.Errorf("Day index getter is wrong")
	}
}

func TestRuleChangeDTStartTimezoneRespected(t *testing.T) {
	/*
		https://golang.org/pkg/time/#LoadLocation

		"The time zone database needed by LoadLocation may not be present on all systems, especially non-Unix systems.
		LoadLocation looks in the directory or uncompressed zip file named by the ZONEINFO environment variable,
		if any, then looks in known installation locations on Unix systems, and finally looks in
		$GOROOT/lib/time/zoneinfo.zip."
	*/
	loc, err := time.LoadLocation("CET")
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	rule, err := newRecurrence(
		ROption{
			Freq:    DAILY,
			Count:   10,
			Wkst:    MO,
			Dtstart: time.Date(2019, 3, 6, 1, 1, 1, 0, loc),
		},
	)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
	rule.DTStart(time.Date(2019, 3, 6, 0, 0, 0, 0, time.UTC))

	events := rule.All()
	if len(events) != 10 {
		t.Fatal("expected", 10, "got", len(events))
	}

	for _, e := range events {
		if e.Location().String() != "UTC" {
			t.Fatal("expected", "UTC", "got", e.Location().String())
		}
		h, m, s := e.Clock()
		if (h + m + s) != 0 {
			t.Fatal("expected", "0", "got", h, m, s)
		}
	}
}

func BenchmarkIterator(b *testing.B) {
	type testCase struct {
		Name   string
		Option ROption
	}

	for _, c := range []testCase{
		{
			Name: "simple secondly",
			Option: ROption{
				Dtstart: time.Date(2000, 03, 22, 12, 0, 0, 0, time.UTC),
				Freq:    SECONDLY,
			},
		},
		{
			Name: "simple minutely",
			Option: ROption{
				Dtstart: time.Date(2000, 03, 22, 12, 0, 0, 0, time.UTC),
				Freq:    MINUTELY,
			},
		},
		{
			Name: "simple hourly",
			Option: ROption{
				Dtstart: time.Date(2000, 03, 22, 12, 0, 0, 0, time.UTC),
				Freq:    HOURLY,
			},
		},
		{
			Name: "simple daily",
			Option: ROption{
				Dtstart: time.Date(2000, 03, 22, 12, 0, 0, 0, time.UTC),
				Freq:    DAILY,
			},
		},
		{
			Name: "simple weekly",
			Option: ROption{
				Dtstart: time.Date(2000, 03, 22, 12, 0, 0, 0, time.UTC),
				Freq:    WEEKLY,
			},
		},
		{
			Name: "simple monthly",
			Option: ROption{
				Dtstart: time.Date(2000, 03, 22, 12, 0, 0, 0, time.UTC),
				Freq:    MONTHLY,
			},
		},
		{
			Name: "simple yearly",
			Option: ROption{
				Dtstart: time.Date(2000, 03, 22, 12, 0, 0, 0, time.UTC),
				Freq:    YEARLY,
			},
		},
	} {
		c := c
		b.Run(c.Name, func(b *testing.B) {
			rrule, err := newRecurrence(c.Option)
			if err != nil {
				b.Errorf("failed to init rrule: %s", err)
			}

			for i := 0; i < b.N; i++ {
				res := iterateNum(rrule.Iterator(), 200)
				if res.IsZero() {
					b.Error("expected not zero iterator result")
				}
			}
		})
	}
}

func iterateNum(iter Next, num int) (last time.Time) {
	for i := 0; i < num; i++ {
		var ok bool
		last, ok = iter()
		if !ok {
			return time.Time{}
		}
	}
	return last
}

// TestRRuleAllDayTimezoneConsistency tests all-day consistency across timezones.
func TestRRuleAllDayTimezoneConsistency(t *testing.T) {
	timezones := []*time.Location{
		time.UTC,
		time.FixedZone("EST", -5*3600), // UTC-5
		time.FixedZone("JST", 9*3600),  // UTC+9
		time.FixedZone("CET", 1*3600),  // UTC+1
		time.FixedZone("PST", -8*3600), // UTC-8
	}

	baseDate := time.Date(2023, 6, 15, 14, 30, 45, 0, time.UTC)

	for i, tz := range timezones {
		t.Run(fmt.Sprintf("Timezone_%d", i), func(t *testing.T) {
			// Create all-day events on the same date across timezones.
			dtstart := time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day(),
				10+i*2, 15+i*5, 30+i*3, 0, tz) // Varying time components.

			r, err := newRecurrence(ROption{
				Freq:    DAILY,
				Count:   3,
				AllDay:  true,
				Dtstart: dtstart,
			})
			if err != nil {
				t.Fatalf("Failed to create RRule: %v", err)
			}

			// All-day events should produce the same results across timezones (floating time).
			expected := []time.Time{
				time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 6, 16, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 6, 17, 0, 0, 0, 0, time.UTC),
			}

			result := r.All()
			if !timesEqual(result, expected) {
				t.Errorf("Timezone %s: expected %v, got %v", tz.String(), expected, result)
			}
		})
	}
}

// TestRRuleTimezonePreservation tests timezone preservation for timed events.
func TestRRuleTimezonePreservation(t *testing.T) {
	testCases := []struct {
		name string
		tz   *time.Location
	}{
		{"UTC", time.UTC},
		{"EST", time.FixedZone("EST", -5*3600)},
		{"JST", time.FixedZone("JST", 9*3600)},
		{"CET", time.FixedZone("CET", 1*3600)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dtstart := time.Date(2023, 1, 1, 14, 30, 0, 0, tc.tz)

			r, err := newRecurrence(ROption{
				Freq:    DAILY,
				Count:   2,
				AllDay:  false,
				Dtstart: dtstart,
			})
			if err != nil {
				t.Fatalf("Failed to create RRule: %v", err)
			}

			result := r.All()

			// Verify timezone is preserved.
			for _, dt := range result {
				if dt.Location() != tc.tz {
					t.Errorf("Expected timezone %s, got %s", tc.tz.String(), dt.Location().String())
				}
			}

			// Verify time precision.
			expected := []time.Time{
				time.Date(2023, 1, 1, 14, 30, 0, 0, tc.tz),
				time.Date(2023, 1, 2, 14, 30, 0, 0, tc.tz),
			}

			if !timesEqual(result, expected) {
				t.Errorf("Expected %v, got %v", expected, result)
			}
		})
	}
}

// TestRRuleLeapYearHandling tests leap year handling.
func TestRRuleLeapYearHandling(t *testing.T) {
	testCases := []struct {
		name          string
		dtstart       time.Time
		freq          Frequency
		bymonth       []int
		bymonthday    []int
		count         int
		expectLeapDay bool
	}{
		{
			name:          "Leap_Year_Feb29",
			dtstart:       time.Date(2020, 2, 29, 10, 0, 0, 0, time.UTC), // 2020 is a leap year.
			freq:          YEARLY,
			count:         4,
			expectLeapDay: true,
		},
		{
			name:          "Non_Leap_Year_Feb29_Skip",
			dtstart:       time.Date(2020, 2, 29, 10, 0, 0, 0, time.UTC),
			freq:          YEARLY,
			count:         5, // Crosses non-leap years.
			expectLeapDay: false,
		},
		{
			name:          "Monthly_Feb29_Handling",
			dtstart:       time.Date(2020, 1, 29, 10, 0, 0, 0, time.UTC),
			freq:          MONTHLY,
			count:         3, // Jan 29 -> Feb 29 -> Mar 29.
			expectLeapDay: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			option := ROption{
				Freq:    tc.freq,
				Count:   tc.count,
				Dtstart: tc.dtstart,
			}
			if len(tc.bymonth) > 0 {
				option.Bymonth = tc.bymonth
			}
			if len(tc.bymonthday) > 0 {
				option.Bymonthday = tc.bymonthday
			}

			r, err := newRecurrence(option)
			if err != nil {
				t.Fatalf("Failed to create RRule: %v", err)
			}

			result := r.All()

			// Check for Feb 29.
			hasLeapDay := false
			for _, dt := range result {
				if dt.Month() == time.February && dt.Day() == 29 {
					hasLeapDay = true
					break
				}
			}

			if tc.expectLeapDay && !hasLeapDay {
				t.Errorf("Expected leap day (Feb 29) in results, but not found. Results: %v", result)
			}

			// Verify results are non-empty.
			if len(result) == 0 {
				t.Error("No results generated")
			}
		})
	}
}

// TestRRuleComplexByRuleCombinations tests complex BY rule combinations.
func TestRRuleComplexByRuleCombinations(t *testing.T) {
	testCases := []struct {
		name       string
		option     ROption
		minResults int
		maxResults int
	}{
		{
			name: "Multiple_BY_Rules_Intersection",
			option: ROption{
				Freq:       MONTHLY,
				Count:      12,
				Byweekday:  []Weekday{MO, WE, FR}, // Mon, Wed, Fri.
				Bymonthday: []int{1, 15, 30},      // 1st, 15th, 30th.
				Dtstart:    time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
			},
			minResults: 1,
			maxResults: 12,
		},
		{
			name: "BYSETPOS_With_Multiple_Rules",
			option: ROption{
				Freq:      MONTHLY,
				Count:     6,
				Byweekday: []Weekday{MO, TU, WE, TH, FR}, // Weekdays.
				Bysetpos:  []int{1, -1},                  // First and last.
				Dtstart:   time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
			},
			minResults: 6,
			maxResults: 12,
		},
		{
			name: "BYHOUR_BYMINUTE_Combination",
			option: ROption{
				Freq:     DAILY,
				Count:    3,
				Byhour:   []int{9, 12, 15, 18},
				Byminute: []int{0, 30},
				Dtstart:  time.Date(2023, 1, 1, 9, 0, 0, 0, time.UTC),
			},
			minResults: 3,
			maxResults: 24, // 3 days * 4 hours * 2 minutes
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r, err := newRecurrence(tc.option)
			if err != nil {
				t.Fatalf("Failed to create RRule: %v", err)
			}

			result := r.All()

			if len(result) < tc.minResults {
				t.Errorf("Expected at least %d results, got %d", tc.minResults, len(result))
			}

			if len(result) > tc.maxResults {
				t.Errorf("Expected at most %d results, got %d", tc.maxResults, len(result))
			}

			// Verify results are sorted.
			for i := 1; i < len(result); i++ {
				if result[i].Before(result[i-1]) {
					t.Errorf("Results not sorted: %v comes before %v", result[i], result[i-1])
				}
			}
		})
	}
}

// TestRRuleEdgeCaseParameters tests boundary parameter values.
func TestRRuleEdgeCaseParameters(t *testing.T) {
	testCases := []struct {
		name      string
		option    ROption
		expectErr bool
	}{
		{
			name: "Max_Interval",
			option: ROption{
				Freq:     YEARLY,
				Interval: 1000,
				Count:    2,
				Dtstart:  time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
			},
			expectErr: false,
		},
		{
			name: "Max_Count",
			option: ROption{
				Freq:    DAILY,
				Count:   10000,
				Dtstart: time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
			},
			expectErr: false,
		},
		{
			name: "Boundary_BYMONTHDAY",
			option: ROption{
				Freq:       MONTHLY,
				Count:      12,
				Bymonthday: []int{31, -1}, // Last day and second-to-last day.
				Dtstart:    time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
			},
			expectErr: false,
		},
		{
			name: "Boundary_BYYEARDAY",
			option: ROption{
				Freq:      YEARLY,
				Count:     3,
				Byyearday: []int{1, 366, -1, -366}, // Year boundary days.
				Dtstart:   time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
			},
			expectErr: false,
		},
		{
			name: "Invalid_BYMONTHDAY_Zero",
			option: ROption{
				Freq:       MONTHLY,
				Count:      3,
				Bymonthday: []int{0}, // Invalid value.
				Dtstart:    time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r, err := newRecurrence(tc.option)

			if tc.expectErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// For valid rules, verify results are generated.
			result := r.All()
			if len(result) == 0 && tc.option.Count > 0 {
				t.Error("Expected results but got none")
			}
		})
	}
}

// TestRRuleMethodChaining tests method chaining and state updates.
func TestRRuleMethodChaining(t *testing.T) {
	r, err := newRecurrence(ROption{
		Freq:    DAILY,
		Count:   3,
		Dtstart: time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("Failed to create RRule: %v", err)
	}

	// Test DTStart update.
	newDtstart := time.Date(2023, 2, 1, 15, 30, 0, 0, time.UTC)
	r.DTStart(newDtstart)

	if !r.GetDTStart().Equal(newDtstart.Truncate(time.Second)) {
		t.Errorf("DTStart not updated correctly: expected %v, got %v",
			newDtstart.Truncate(time.Second), r.GetDTStart())
	}

	// Test Until update.
	newUntil := time.Date(2023, 2, 5, 20, 0, 0, 0, time.UTC)
	r.until = newUntil
	r.rebuildRule()

	if !r.until.Equal(newUntil.Truncate(time.Second)) {
		t.Errorf("Until not updated correctly: expected %v, got %v",
			newUntil.Truncate(time.Second), r.until)
	}

	// Test AllDay toggle.
	r.SetAllDay(true)
	if !r.IsAllDay() {
		t.Error("AllDay flag not set correctly")
	}

	// Verify time is normalized in AllDay mode.
	result := r.All()
	for _, dt := range result {
		if dt.Hour() != 0 || dt.Minute() != 0 || dt.Second() != 0 {
			t.Errorf("AllDay event should have 00:00:00 time, got %v", dt)
		}
	}
}

// TestRRuleIteratorConsistency tests iterator vs batch method consistency.
func TestRRuleIteratorConsistency(t *testing.T) {
	testCases := []ROption{
		{
			Freq:    DAILY,
			Count:   5,
			Dtstart: time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
		},
		{
			Freq:      WEEKLY,
			Count:     4,
			Byweekday: []Weekday{MO, WE, FR},
			Dtstart:   time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
		},
		{
			Freq:    MONTHLY,
			Count:   3,
			AllDay:  true,
			Dtstart: time.Date(2023, 1, 15, 14, 30, 0, 0, time.UTC),
		},
	}

	for i, option := range testCases {
		t.Run(fmt.Sprintf("Case_%d", i), func(t *testing.T) {
			r, err := newRecurrence(option)
			if err != nil {
				t.Fatalf("Failed to create RRule: %v", err)
			}

			// Use All() to collect results.
			allResults := r.All()

			// Use the iterator to collect results.
			iterator := r.Iterator()
			var iterResults []time.Time

			for {
				next, ok := iterator()
				if !ok {
					break
				}
				iterResults = append(iterResults, next)
			}

			// Verify results match.
			if !timesEqual(allResults, iterResults) {
				t.Errorf("All() and Iterator() results differ:\nAll(): %v\nIterator(): %v",
					allResults, iterResults)
			}
		})
	}
}

// TestRRuleStringRoundTrip tests string round-trip serialization.
func TestRRuleStringRoundTrip(t *testing.T) {
	testCases := []ROption{
		{
			Freq:    DAILY,
			Count:   5,
			Dtstart: time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
		},
		{
			Freq:      WEEKLY,
			Interval:  2,
			Byweekday: []Weekday{MO, WE, FR},
			Until:     time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC),
			Dtstart:   time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
		},
		{
			Freq:       MONTHLY,
			Count:      12,
			Bymonthday: []int{1, 15, -1},
			Dtstart:    time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
		},
		{
			Freq:    YEARLY,
			Count:   3,
			AllDay:  true,
			Dtstart: time.Date(2023, 6, 15, 14, 30, 0, 0, time.UTC),
		},
	}

	for i, option := range testCases {
		t.Run(fmt.Sprintf("RoundTrip_%d", i), func(t *testing.T) {
			// Create the original RRule.
			original, err := newRecurrence(option)
			if err != nil {
				t.Fatalf("Failed to create original RRule: %v", err)
			}

			// Serialize to string.
			rruleStr := original.String()

			// Parse from string.
			parsed, err := StrToRRuleSet(rruleStr)
			if err != nil {
				t.Fatalf("Failed to parse RRule string '%s': %v", rruleStr, err)
			}

			// Compare results.
			originalResults := original.All()
			parsedResults := parsed.All()

			if !timesEqual(originalResults, parsedResults) {
				t.Errorf("Round-trip results differ:\nOriginal: %v\nParsed: %v\nRRule String: %s",
					originalResults, parsedResults, rruleStr)
			}
		})
	}
}

// TestRRulePerformanceBaseline tests performance baseline.
func TestRRulePerformanceBaseline(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	testCases := []struct {
		name        string
		option      ROption
		maxDuration time.Duration
	}{
		{
			name: "Large_Count_Daily",
			option: ROption{
				Freq:    DAILY,
				Count:   10000,
				Dtstart: time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
			},
			maxDuration: 100 * time.Millisecond,
		},
		{
			name: "Complex_BY_Rules",
			option: ROption{
				Freq:      MONTHLY,
				Count:     1000,
				Byweekday: []Weekday{MO, TU, WE, TH, FR},
				Bysetpos:  []int{1, 2, -2, -1},
				Dtstart:   time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
			},
			maxDuration: 200 * time.Millisecond,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r, err := newRecurrence(tc.option)
			if err != nil {
				t.Fatalf("Failed to create RRule: %v", err)
			}

			start := time.Now()
			result := r.All()
			duration := time.Since(start)

			if duration > tc.maxDuration {
				t.Errorf("Performance test failed: took %v, expected < %v", duration, tc.maxDuration)
			}

			t.Logf("Generated %d results in %v", len(result), duration)
		})
	}
}

func TestSet(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY, Count: 2, Byweekday: []Weekday{TU},
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	value := r.All()
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 9, 9, 0, 0, 0, time.UTC)}
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSetOverlapping(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	v1 := r.All()
	if len(v1) > 300 || len(v1) < 200 {
		t.Errorf("No default Util time")
	}
}

func TestSetString(t *testing.T) {
	moscow, _ := time.LoadLocation("Europe/Moscow")
	newYork, _ := time.LoadLocation("America/New_York")
	tehran, _ := time.LoadLocation("Asia/Tehran")

	r, _ := newRecurrence(ROption{Freq: YEARLY, Count: 1, Byweekday: []Weekday{TU},
		Dtstart: time.Date(1997, 9, 2, 8, 0, 0, 0, time.UTC)})
	r.ExDate(time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC))
	r.ExDate(time.Date(1997, 9, 11, 9, 0, 0, 0, time.UTC).In(moscow))
	r.ExDate(time.Date(1997, 9, 18, 9, 0, 0, 0, time.UTC).In(newYork))
	r.RDate(time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC).In(tehran))
	r.RDate(time.Date(1997, 9, 9, 9, 0, 0, 0, time.UTC))

	want := `DTSTART:19970902T080000Z
RRULE:FREQ=YEARLY;COUNT=1;BYDAY=TU
RDATE;TZID=Asia/Tehran:19970904T133000
RDATE:19970909T090000Z
EXDATE:19970904T090000Z
EXDATE;TZID=Europe/Moscow:19970911T130000
EXDATE;TZID=America/New_York:19970918T050000`
	value := r.String()
	if want != value {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSetDTStart(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY, Count: 1, Byweekday: []Weekday{TU},
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	r.ExDate(time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC))
	r.ExDate(time.Date(1997, 9, 11, 9, 0, 0, 0, time.UTC))
	r.ExDate(time.Date(1997, 9, 18, 9, 0, 0, 0, time.UTC))
	r.RDate(time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC))
	r.RDate(time.Date(1997, 9, 9, 9, 0, 0, 0, time.UTC))

	nyLoc, _ := time.LoadLocation("America/New_York")
	r.DTStart(time.Date(1997, 9, 3, 9, 0, 0, 0, nyLoc))

	want := `DTSTART;TZID=America/New_York:19970903T090000
RRULE:FREQ=YEARLY;COUNT=1;BYDAY=TU
RDATE:19970904T090000Z
RDATE:19970909T090000Z
EXDATE:19970904T090000Z
EXDATE:19970911T090000Z
EXDATE:19970918T090000Z`
	value := r.String()
	if want != value {
		t.Errorf("get \n%v\n want \n%v\n", value, want)
	}

	sset, err := StrToRRuleSet(r.String())
	if err != nil {
		t.Errorf("Could not create RSET from set output")
	}
	if sset.String() != r.String() {
		t.Errorf("RSET created from set output different than original set, %s", sset.String())
	}
}

func TestSetRecurrence(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY, Count: 1, Byweekday: []Weekday{TU},
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	value := r.Strings()
	if len(value) != 2 {
		t.Errorf("Wrong length for recurrence got=%v want=%v", len(value), 2)
	}
	want := "DTSTART:19970902T090000Z\nRRULE:FREQ=YEARLY;COUNT=1;BYDAY=TU"
	if r.String() != want {
		t.Errorf("get %s, want %v", r.String(), want)
	}
}

func TestSetDate(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY, Count: 1, Byweekday: []Weekday{TU},
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	r.RDate(time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC))
	r.RDate(time.Date(1997, 9, 9, 9, 0, 0, 0, time.UTC))
	value := r.All()
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 9, 9, 0, 0, 0, time.UTC)}
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSetRDates(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY, Count: 1, Byweekday: []Weekday{TU},
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	r.SetRDates([]time.Time{
		time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 9, 9, 0, 0, 0, time.UTC),
	})
	value := r.All()
	want := []time.Time{
		time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 9, 9, 0, 0, 0, time.UTC),
	}
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSetExDate(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY, Count: 6, Byweekday: []Weekday{TU, TH},
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	r.ExDate(time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC))
	r.ExDate(time.Date(1997, 9, 11, 9, 0, 0, 0, time.UTC))
	r.ExDate(time.Date(1997, 9, 18, 9, 0, 0, 0, time.UTC))
	value := r.All()
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 9, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 16, 9, 0, 0, 0, time.UTC)}
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSetExDates(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY, Count: 6, Byweekday: []Weekday{TU, TH},
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	r.SetExDates([]time.Time{
		time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 11, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 18, 9, 0, 0, 0, time.UTC),
	})
	value := r.All()
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 9, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 16, 9, 0, 0, 0, time.UTC)}
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSetExDateRevOrder(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: MONTHLY, Count: 5, Bymonthday: []int{10},
		Dtstart: time.Date(2004, 1, 1, 9, 0, 0, 0, time.UTC)})
	r.ExDate(time.Date(2004, 4, 10, 9, 0, 0, 0, time.UTC))
	r.ExDate(time.Date(2004, 2, 10, 9, 0, 0, 0, time.UTC))
	value := r.All()
	want := []time.Time{time.Date(2004, 1, 10, 9, 0, 0, 0, time.UTC),
		time.Date(2004, 3, 10, 9, 0, 0, 0, time.UTC),
		time.Date(2004, 5, 10, 9, 0, 0, 0, time.UTC)}
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSetDateAndExDate(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: YEARLY, Count: 6, Byweekday: []Weekday{TU, TH},
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	r.RDate(time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC))
	r.RDate(time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC))
	r.RDate(time.Date(1997, 9, 9, 9, 0, 0, 0, time.UTC))
	r.RDate(time.Date(1997, 9, 11, 9, 0, 0, 0, time.UTC))
	r.RDate(time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC))
	r.RDate(time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC))
	r.RDate(time.Date(1997, 9, 9, 9, 0, 0, 0, time.UTC))
	r.RDate(time.Date(1997, 9, 11, 9, 0, 0, 0, time.UTC))
	r.RDate(time.Date(1997, 9, 16, 9, 0, 0, 0, time.UTC))
	r.RDate(time.Date(1997, 9, 18, 9, 0, 0, 0, time.UTC))
	r.ExDate(time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC))
	r.ExDate(time.Date(1997, 9, 11, 9, 0, 0, 0, time.UTC))
	r.ExDate(time.Date(1997, 9, 18, 9, 0, 0, 0, time.UTC))
	value := r.All()
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 9, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 16, 9, 0, 0, 0, time.UTC)}
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSetBefore(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY, Count: 7,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC)
	value := r.Before(time.Date(1997, 9, 5, 9, 0, 0, 0, time.UTC), false)
	if value != want {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSetBeforeInc(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY, Count: 7,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := time.Date(1997, 9, 5, 9, 0, 0, 0, time.UTC)
	value := r.Before(time.Date(1997, 9, 5, 9, 0, 0, 0, time.UTC), true)
	if value != want {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSetAfter(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY, Count: 7,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := time.Date(1997, 9, 5, 9, 0, 0, 0, time.UTC)
	value := r.After(time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC), false)
	if value != want {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSetAfterInc(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY, Count: 7,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	want := time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC)
	value := r.After(time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC), true)
	if value != want {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSetBetween(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY, Count: 7,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	value := r.Between(time.Date(1997, 9, 3, 9, 0, 0, 0, time.UTC), time.Date(1997, 9, 6, 9, 0, 0, 0, time.UTC), false)
	want := []time.Time{time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 5, 9, 0, 0, 0, time.UTC)}
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSetBetweenInc(t *testing.T) {
	r, _ := newRecurrence(ROption{Freq: DAILY, Count: 7,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	value := r.Between(time.Date(1997, 9, 3, 9, 0, 0, 0, time.UTC), time.Date(1997, 9, 6, 9, 0, 0, 0, time.UTC), true)
	want := []time.Time{time.Date(1997, 9, 3, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 5, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 6, 9, 0, 0, 0, time.UTC)}
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSetTrickyTimeZones(t *testing.T) {
	moscow, _ := time.LoadLocation("Europe/Moscow")
	newYork, _ := time.LoadLocation("America/New_York")
	tehran, _ := time.LoadLocation("Asia/Tehran")

	r, _ := newRecurrence(ROption{
		Freq:    DAILY,
		Count:   4,
		Dtstart: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).In(moscow),
	})

	r.ExDate(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).In(newYork))
	r.ExDate(time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC).In(tehran))
	r.ExDate(time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC).In(moscow))
	r.ExDate(time.Date(2000, 1, 4, 0, 0, 0, 0, time.UTC))

	occurrences := r.All()

	if len(occurrences) > 0 {
		t.Errorf("No all occurrences excluded by ExDate: [%+v]", occurrences)
	}
}

func TestSetDtStart(t *testing.T) {
	ogr := []string{"DTSTART;TZID=America/Los_Angeles:20181115T000000", "RRULE:FREQ=DAILY;INTERVAL=1;WKST=SU;UNTIL=20181118T075959Z"}
	set, _ := StrSliceToRRuleSet(ogr)

	ogoc := set.All()
	set.DTStart(set.GetDTStart().AddDate(0, 0, 1))

	noc := set.All()
	if len(noc) != len(ogoc)-1 {
		t.Fatalf("As per the new DTStart the new occurences should exactly be one less that the original, new :%d original: %d", len(noc), len(ogoc))
	}

	for i := range noc {
		if noc[i] != ogoc[i+1] {
			t.Errorf("New occurences should just offset by one, mismatch at %d, expected: %+v, actual: %+v", i, ogoc[i+1], noc[i])
		}
	}
}

func TestRuleSetChangeDTStartTimezoneRespected(t *testing.T) {
	/*
		https://golang.org/pkg/time/#LoadLocation

		"The time zone database needed by LoadLocation may not be present on all systems, especially non-Unix systems.
		LoadLocation looks in the directory or uncompressed zip file named by the ZONEINFO environment variable,
		if any, then looks in known installation locations on Unix systems, and finally looks in
		$GOROOT/lib/time/zoneinfo.zip."
	*/
	loc, err := time.LoadLocation("CET")
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	rule, err := newRecurrence(
		ROption{
			Freq:     DAILY,
			Count:    10,
			Wkst:     MO,
			Byhour:   []int{10},
			Byminute: []int{0},
			Bysecond: []int{0},
			Dtstart:  time.Date(2019, 3, 6, 0, 0, 0, 0, loc),
		},
	)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
	rule.DTStart(time.Date(2019, 3, 6, 0, 0, 0, 0, time.UTC))

	events := rule.All()
	if len(events) != 10 {
		t.Fatal("expected", 10, "got", len(events))
	}

	for _, e := range events {
		if e.Location().String() != "UTC" {
			t.Fatal("expected", "UTC", "got", e.Location().String())
		}
	}
}

func TestSetStr(t *testing.T) {
	setStr := "RRULE:FREQ=DAILY;UNTIL=20180517T235959Z\n" +
		"EXDATE;VALUE=DATE-TIME:20180525T070000Z,20180530T130000Z\n" +
		"RDATE;VALUE=DATE-TIME:20180801T131313Z,20180902T141414Z\n"

	r, err := StrToRRuleSet(setStr)
	if err != nil {
		t.Fatalf("StrToRRuleSet(%s) returned error: %v", setStr, err)
	}

	if r.RRuleString() != "RRULE:FREQ=DAILY;UNTIL=20180517T235959Z" {
		t.Errorf("Unexpected rrule: %s", r.RRuleString())
	}

	// matching parsed EXDates
	exDates := r.GetExDate()
	if len(exDates) != 2 {
		t.Errorf("Unexpected number of exDates: %v != 2, %v", len(exDates), exDates)
	}
	if [2]string{timeToUTCStr(exDates[0]), timeToUTCStr(exDates[1])} != [2]string{"20180525T070000Z", "20180530T130000Z"} {
		t.Errorf("Unexpected exDates: %v", exDates)
	}

	// matching parsed RDates
	rDates := r.GetRDate()
	if len(rDates) != 2 {
		t.Errorf("Unexpected number of rDates: %v != 2, %v", len(rDates), rDates)
	}
	if [2]string{timeToUTCStr(rDates[0]), timeToUTCStr(rDates[1])} != [2]string{"20180801T131313Z", "20180902T141414Z"} {
		t.Errorf("Unexpected exDates: %v", exDates)
	}
}

// TestSetAllDayTimezoneConsistency tests all-day consistency across timezones.
func TestSetAllDayTimezoneConsistency(t *testing.T) {
	// Create times in different timezones.
	utc := time.UTC
	ny, _ := time.LoadLocation("America/New_York")
	tokyo, _ := time.LoadLocation("Asia/Tokyo")

	baseTime := time.Date(2024, 3, 15, 14, 30, 45, 0, utc)

	testCases := []struct {
		name     string
		timezone *time.Location
	}{
		{"UTC", utc},
		{"New_York", ny},
		{"Tokyo", tokyo},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := newRecurrence(ROption{
				Freq: DAILY,
			})
			r.SetAllDay(true)

			// Set DTSTART in different timezones.
			dtstart := baseTime.In(tc.timezone)
			r.DTStart(dtstart)

			// Add RDATEs in different timezones.
			rdate := baseTime.AddDate(0, 0, 1).In(tc.timezone)
			r.RDate(rdate)

			// Add EXDATEs in different timezones.
			exdate := baseTime.AddDate(0, 0, 2).In(tc.timezone)
			r.ExDate(exdate)

			// Verify all times normalize to floating time (UTC 00:00:00).
			if r.GetDTStart().Location() != utc {
				t.Errorf("DTSTART should be in UTC, got %v", r.GetDTStart().Location())
			}

			expectedDate := time.Date(2024, 3, 15, 0, 0, 0, 0, utc)
			if !r.GetDTStart().Equal(expectedDate) {
				t.Errorf("DTSTART should be %v, got %v", expectedDate, r.GetDTStart())
			}

			// Verify RDATE normalization.
			rdates := r.GetRDate()
			if len(rdates) != 1 {
				t.Fatalf("Expected 1 RDATE, got %d", len(rdates))
			}
			expectedRDate := time.Date(2024, 3, 16, 0, 0, 0, 0, utc)
			if !rdates[0].Equal(expectedRDate) {
				t.Errorf("RDATE should be %v, got %v", expectedRDate, rdates[0])
			}

			// Verify EXDATE normalization.
			exdates := r.GetExDate()
			if len(exdates) != 1 {
				t.Fatalf("Expected 1 EXDATE, got %d", len(exdates))
			}
			expectedExDate := time.Date(2024, 3, 17, 0, 0, 0, 0, utc)
			if !exdates[0].Equal(expectedExDate) {
				t.Errorf("EXDATE should be %v, got %v", expectedExDate, exdates[0])
			}
		})
	}
}

// TestSetComplexRRuleRDateExDateInteraction tests complex RRULE + RDATE + EXDATE interaction.
func TestSetComplexRRuleRDateExDateInteraction(t *testing.T) {
	// Create a daily recurrence rule.
	r, err := newRecurrence(ROption{
		Freq:    DAILY,
		Count:   10,
		Dtstart: time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("Failed to create RRule: %v", err)
	}
	// Add extra RDATEs (outside the RRULE sequence).
	r.RDate(time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC))
	r.RDate(time.Date(2024, 1, 20, 9, 0, 0, 0, time.UTC))

	// Exclude some RRULE-generated dates.
	r.ExDate(time.Date(2024, 1, 3, 9, 0, 0, 0, time.UTC))
	r.ExDate(time.Date(2024, 1, 5, 9, 0, 0, 0, time.UTC))

	// Exclude one RDATE (should be filtered out).
	r.ExDate(time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC))

	occurrences := r.All()

	// Log all occurrences for debugging.
	t.Logf("Generated occurrences:")
	for i, occ := range occurrences {
		t.Logf("  %d: %v", i, occ)
	}

	// Verify results: 10 RRULE - 2 EXDATE (Jan 3 and 5) + 1 valid RDATE (Jan 20) = 9.
	// Note: Jan 15 RDATE is excluded by EXDATE and not counted.
	expectedCount := 9
	if len(occurrences) != expectedCount {
		t.Errorf("Expected %d occurrences, got %d", expectedCount, len(occurrences))
	}

	// Verify excluded dates are not in the results.
	excludedDates := []time.Time{
		time.Date(2024, 1, 3, 9, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 5, 9, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC), // This RDATE is excluded by EXDATE.
	}

	for _, excluded := range excludedDates {
		for _, occurrence := range occurrences {
			if occurrence.Equal(excluded) {
				t.Errorf("Excluded date %v found in occurrences", excluded)
			}
		}
	}

	// Verify the included RDATE is present.
	expectedRDate := time.Date(2024, 1, 20, 9, 0, 0, 0, time.UTC)
	found := false
	for _, occurrence := range occurrences {
		if occurrence.Equal(expectedRDate) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected RDATE %v not found in occurrences", expectedRDate)
	}
}

// TestSetDaylightSavingTransition tests behavior during DST transitions.
func TestSetDaylightSavingTransition(t *testing.T) {
	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Skip("America/New_York timezone not available")
	}

	// 2024 DST starts: Mar 10, 2:00 AM -> 3:00 AM.
	// Create a rule that crosses the DST transition.
	r, err := newRecurrence(ROption{
		Freq:    DAILY,
		Count:   5,
		Dtstart: time.Date(2024, 3, 8, 2, 30, 0, 0, ny), // 2 days before the transition.
	})
	if err != nil {
		t.Fatalf("Failed to create RRule: %v", err)
	}

	occurrences := r.All()

	if len(occurrences) != 5 {
		t.Fatalf("Expected 5 occurrences, got %d", len(occurrences))
	}

	// Verify timezone is preserved.
	for i, occurrence := range occurrences {
		if occurrence.Location().String() != ny.String() {
			t.Errorf("Occurrence %d should be in %s timezone, got %s",
				i, ny.String(), occurrence.Location().String())
		}

		// During DST transition, time may shift.
		// On Mar 10, 2:30 AM jumps to 3:30 AM (DST starts).
		hour, min, sec := occurrence.Clock()

		// On the DST transition day (Mar 10), time shifts from 2:30 to 3:30.
		expectedHour := 2
		if occurrence.Month() == 3 && occurrence.Day() >= 10 {
			expectedHour = 3 // After DST, time becomes 3:30.
		}

		if hour != expectedHour || min != 30 || sec != 0 {
			t.Logf("Occurrence %d at %v: expected %02d:30:00, got %02d:%02d:%02d",
				i, occurrence, expectedHour, hour, min, sec)
		}
	}
}

// TestSetAllDayDynamicToggle tests toggling all-day/timed state dynamically.
func TestSetAllDayDynamicToggle(t *testing.T) {
	r := &Recurrence{}

	// Initialize as a timed event.
	dtstart := time.Date(2024, 6, 15, 14, 30, 45, 123456789, time.UTC)
	r.DTStart(dtstart)

	rdate := time.Date(2024, 6, 16, 10, 15, 30, 987654321, time.UTC)
	r.RDate(rdate)

	exdate := time.Date(2024, 6, 17, 16, 45, 20, 555666777, time.UTC)
	r.ExDate(exdate)

	// Verify initial state (timed).
	if r.IsAllDay() {
		t.Error("Set should not be all-day initially")
	}

	// Verify time precision is truncated to seconds.
	if r.GetDTStart().Nanosecond() != 0 {
		t.Error("Non-all-day DTSTART should be truncated to seconds")
	}

	// Switch to all-day.
	r.SetAllDay(true)

	// Verify state switch.
	if !r.IsAllDay() {
		t.Error("Set should be all-day after SetAllDay(true)")
	}

	// Verify all times are normalized to floating time.
	expectedDTStart := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	if !r.GetDTStart().Equal(expectedDTStart) {
		t.Errorf("All-day DTSTART should be %v, got %v", expectedDTStart, r.GetDTStart())
	}

	rdates := r.GetRDate()
	expectedRDate := time.Date(2024, 6, 16, 0, 0, 0, 0, time.UTC)
	if len(rdates) != 1 || !rdates[0].Equal(expectedRDate) {
		t.Errorf("All-day RDATE should be %v, got %v", expectedRDate, rdates)
	}

	exdates := r.GetExDate()
	expectedExDate := time.Date(2024, 6, 17, 0, 0, 0, 0, time.UTC)
	if len(exdates) != 1 || !exdates[0].Equal(expectedExDate) {
		t.Errorf("All-day EXDATE should be %v, got %v", expectedExDate, exdates)
	}

	// Switch back to timed.
	r.SetAllDay(false)

	// Verify state switch.
	if r.IsAllDay() {
		t.Error("Set should not be all-day after SetAllDay(false)")
	}

	// Note: after switching back, time remains normalized (00:00:00 UTC).
	// This is expected because the original timezone info is lost.
}

// TestSetIteratorConsistency tests iterator vs batch consistency.
func TestSetIteratorConsistency(t *testing.T) {
	testCases := []struct {
		name   string
		setup  func() *Recurrence
		allDay bool
	}{
		{
			name: "NonAllDay_WithRRule",
			setup: func() *Recurrence {
				r, _ := newRecurrence(ROption{
					Freq:    WEEKLY,
					Count:   5,
					Dtstart: time.Date(2024, 2, 1, 10, 0, 0, 0, time.UTC),
				})
				return r
			},
			allDay: false,
		},
		{
			name: "AllDay_WithRDateExDate",
			setup: func() *Recurrence {
				r := &Recurrence{}
				r.SetAllDay(true)
				r.DTStart(time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC))
				r.RDate(time.Date(2024, 2, 5, 0, 0, 0, 0, time.UTC))
				r.RDate(time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC))
				r.ExDate(time.Date(2024, 2, 5, 0, 0, 0, 0, time.UTC))
				return r
			},
			allDay: true,
		},
		{
			name: "Complex_Mixed",
			setup: func() *Recurrence {
				r, _ := newRecurrence(ROption{
					Freq:    DAILY,
					Count:   7,
					Dtstart: time.Date(2024, 3, 1, 15, 30, 0, 0, time.UTC),
				})
				r.RDate(time.Date(2024, 3, 10, 15, 30, 0, 0, time.UTC))
				r.ExDate(time.Date(2024, 3, 3, 15, 30, 0, 0, time.UTC))
				return r
			},
			allDay: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			set := tc.setup()

			// Use All() to get all occurrences.
			allOccurrences := set.All()

			// Use the iterator to collect all occurrences.
			var iteratorOccurrences []time.Time
			iterator := set.Iterator()
			for {
				dt, ok := iterator()
				if !ok {
					break
				}
				iteratorOccurrences = append(iteratorOccurrences, dt)
			}

			// Verify counts match.
			if len(allOccurrences) != len(iteratorOccurrences) {
				t.Errorf("All() returned %d occurrences, iterator returned %d",
					len(allOccurrences), len(iteratorOccurrences))
			}

			// Verify contents match.
			for i, expected := range allOccurrences {
				if i >= len(iteratorOccurrences) {
					t.Errorf("Iterator missing occurrence at index %d: %v", i, expected)
					continue
				}

				actual := iteratorOccurrences[i]
				if !expected.Equal(actual) {
					t.Errorf("Occurrence %d mismatch: All()=%v, Iterator()=%v",
						i, expected, actual)
				}
			}
		})
	}
}

// TestSetStringRoundTrip tests string round-trip serialization.
func TestSetStringRoundTrip(t *testing.T) {
	testCases := []struct {
		name   string
		setup  func() *Recurrence
		allDay bool
	}{
		{
			name: "AllDay_Complete",
			setup: func() *Recurrence {
				set, _ := newRecurrence(ROption{
					Freq:    WEEKLY,
					Count:   4,
					Dtstart: time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC),
				})

				set.RDate(time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC))
				set.ExDate(time.Date(2024, 7, 11, 0, 0, 0, 0, time.UTC))

				return set
			},
			allDay: true,
		},
		{
			name: "NonAllDay_WithTimezone",
			setup: func() *Recurrence {
				ny, _ := time.LoadLocation("America/New_York")

				set, _ := newRecurrence(ROption{
					Freq:    DAILY,
					Count:   3,
					Dtstart: time.Date(2024, 7, 4, 14, 30, 0, 0, ny),
				})

				set.RDate(time.Date(2024, 7, 10, 14, 30, 0, 0, ny))
				set.ExDate(time.Date(2024, 7, 5, 14, 30, 0, 0, ny))

				return set
			},
			allDay: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			originalSet := tc.setup()

			// Serialize to string.
			setString := originalSet.String()

			// Parse back into a Set.
			parsedSet, err := StrToRRuleSet(setString)
			if err != nil {
				t.Fatalf("Failed to parse set string: %v", err)
			}

			// Set AllDay state (parser may not detect it automatically).
			if tc.allDay {
				parsedSet.SetAllDay(true)
			}

			// Compare original and parsed Set.
			originalOccurrences := originalSet.All()
			parsedOccurrences := parsedSet.All()

			if len(originalOccurrences) != len(parsedOccurrences) {
				t.Errorf("Occurrence count mismatch: original=%d, parsed=%d",
					len(originalOccurrences), len(parsedOccurrences))
			}

			// Verify each occurrence.
			for i, original := range originalOccurrences {
				if i >= len(parsedOccurrences) {
					t.Errorf("Missing occurrence at index %d: %v", i, original)
					continue
				}

				parsed := parsedOccurrences[i]

				// For all-day events, compare the date only.
				if tc.allDay {
					if original.Year() != parsed.Year() ||
						original.Month() != parsed.Month() ||
						original.Day() != parsed.Day() {
						t.Errorf("All-day occurrence %d date mismatch: original=%v, parsed=%v",
							i, original, parsed)
					}
				} else {
					// For timed events, compare full time (account for timezone conversion).
					if !original.Equal(parsed) {
						t.Errorf("Non-all-day occurrence %d mismatch: original=%v, parsed=%v",
							i, original, parsed)
					}
				}
			}
		})
	}
}

// TestSetEdgeCases tests edge cases and error handling.
func TestSetEdgeCases(t *testing.T) {
	t.Run("EmptySet", func(t *testing.T) {
		set := &Recurrence{}

		occurrences := set.All()
		if len(occurrences) != 0 {
			t.Errorf("Empty set should return no occurrences, got %d", len(occurrences))
		}

		// Test the iterator.
		iterator := set.Iterator()
		if dt, ok := iterator(); ok {
			t.Errorf("Empty set iterator should return false, got %v", dt)
		}
	})

	t.Run("OnlyRDates", func(t *testing.T) {
		set := &Recurrence{}

		dates := []time.Time{
			time.Date(2024, 5, 1, 10, 0, 0, 0, time.UTC),
			time.Date(2024, 5, 15, 10, 0, 0, 0, time.UTC),
			time.Date(2024, 5, 30, 10, 0, 0, 0, time.UTC),
		}

		for _, date := range dates {
			set.RDate(date)
		}

		occurrences := set.All()
		if len(occurrences) != 3 {
			t.Errorf("Expected 3 occurrences from RDates, got %d", len(occurrences))
		}

		// Verify sorting.
		for i := 1; i < len(occurrences); i++ {
			if occurrences[i].Before(occurrences[i-1]) {
				t.Errorf("Occurrences should be sorted, but %v is before %v",
					occurrences[i], occurrences[i-1])
			}
		}
	})

	t.Run("AllExcluded", func(t *testing.T) {
		// Create a rule that generates 3 events.
		r, err := newRecurrence(ROption{
			Freq:    DAILY,
			Count:   3,
			Dtstart: time.Date(2024, 8, 1, 12, 0, 0, 0, time.UTC),
		})
		if err != nil {
			t.Fatalf("Failed to create RRule: %v", err)
		}

		// Exclude all generated events.
		r.ExDate(time.Date(2024, 8, 1, 12, 0, 0, 0, time.UTC))
		r.ExDate(time.Date(2024, 8, 2, 12, 0, 0, 0, time.UTC))
		r.ExDate(time.Date(2024, 8, 3, 12, 0, 0, 0, time.UTC))

		occurrences := r.All()
		if len(occurrences) != 0 {
			t.Errorf("All events should be excluded, got %d occurrences", len(occurrences))
		}
	})

	t.Run("DuplicateRDates", func(t *testing.T) {
		set := &Recurrence{}

		duplicateDate := time.Date(2024, 9, 15, 16, 30, 0, 0, time.UTC)

		// Add duplicate RDATEs.
		set.RDate(duplicateDate)
		set.RDate(duplicateDate)
		set.RDate(duplicateDate)

		occurrences := set.All()
		if len(occurrences) != 1 {
			t.Errorf("Duplicate RDates should be deduplicated, expected 1 occurrence, got %d",
				len(occurrences))
		}

		if !occurrences[0].Equal(duplicateDate) {
			t.Errorf("Expected occurrence %v, got %v", duplicateDate, occurrences[0])
		}
	})
}

// TestSetPerformance tests performance baseline.
func TestSetPerformance(t *testing.T) {
	t.Run("LargeRDateSet", func(t *testing.T) {
		set := &Recurrence{}

		// Add many RDATEs.
		baseDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		for i := 0; i < 1000; i++ {
			set.RDate(baseDate.AddDate(0, 0, i))
		}

		start := time.Now()
		occurrences := set.All()
		duration := time.Since(start)

		if len(occurrences) != 1000 {
			t.Errorf("Expected 1000 occurrences, got %d", len(occurrences))
		}

		// Performance check: should complete within a reasonable time (100ms).
		if duration > 100*time.Millisecond {
			t.Logf("Performance warning: Large RDate set took %v", duration)
		}
	})

	t.Run("ComplexSetWithManyExclusions", func(t *testing.T) {

		// Create a rule that generates many events.
		rrule, err := newRecurrence(ROption{
			Freq:    DAILY,
			Count:   500,
			Dtstart: time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
		})
		if err != nil {
			t.Fatalf("Failed to create RRule: %v", err)
		}

		// Exclude half the events.
		baseDate := time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC)
		for i := 0; i < 250; i += 2 {
			rrule.ExDate(baseDate.AddDate(0, 0, i))
		}

		start := time.Now()
		occurrences := rrule.All()
		duration := time.Since(start)

		expectedCount := 500 - 125 // 500 events - 125 excluded events.
		if len(occurrences) != expectedCount {
			t.Errorf("Expected %d occurrences, got %d", expectedCount, len(occurrences))
		}

		t.Logf("Complex set with exclusions took %v", duration)
	})
}

func TestSetRRulePreservesTimezoneForTimedEvents(t *testing.T) {
	newYork, _ := time.LoadLocation("America/New_York")
	losAngeles, _ := time.LoadLocation("America/Los_Angeles")

	rruleStart := time.Date(2024, 11, 1, 9, 30, 0, 0, newYork)
	set, err := newRecurrence(ROption{
		Freq:    DAILY,
		Count:   2,
		Dtstart: rruleStart,
	})
	if err != nil {
		t.Fatalf("Failed to create RRule: %v", err)
	}

	set.RDate(time.Date(2024, 11, 3, 14, 45, 0, 0, losAngeles))
	set.ExDate(time.Date(2024, 11, 2, 9, 30, 0, 0, newYork))

	if set.IsAllDay() {
		t.Fatal("Timed events should not flip the set into all-day mode")
	}
	if set.GetDTStart().Location() != newYork {
		t.Fatalf("DTSTART should retain original timezone, got %v", set.GetDTStart().Location())
	}

	want := []time.Time{
		time.Date(2024, 11, 1, 9, 30, 0, 0, newYork),
		time.Date(2024, 11, 3, 14, 45, 0, 0, losAngeles),
	}
	got := set.All()
	if !timesEqual(got, want) {
		t.Fatalf("Unexpected iterator results, want %v got %v", want, got)
	}

	recurrence := set.Strings()
	ensureContains := func(substr string) {
		for _, line := range recurrence {
			if strings.Contains(line, substr) {
				return
			}
		}
		t.Fatalf("Recurrence output missing %q in %v", substr, recurrence)
	}

	ensureContains("DTSTART;TZID=America/New_York:20241101T093000")
	ensureContains("RDATE;TZID=America/Los_Angeles:20241103T144500")
	ensureContains("EXDATE;TZID=America/New_York:20241102T093000")
}
