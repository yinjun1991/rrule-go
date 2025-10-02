package rrule

import (
    "testing"
    "time"
)

func TestNewRRule_ValidateBounds_NegativeCases(t *testing.T) {
    base := ROption{Freq: DAILY, Dtstart: mustUTC(2024, 1, 1, 0, 0, 0)}

    cases := []ROption{
        // bysecond out of range
        merge(base, ROption{Bysecond: []int{60}}),
        // byminute out of range
        merge(base, ROption{Byminute: []int{60}}),
        // byhour out of range
        merge(base, ROption{Byhour: []int{24}}),
        // bymonthday out of range
        merge(base, ROption{Bymonthday: []int{0}}),
        merge(base, ROption{Bymonthday: []int{32}}),
        merge(base, ROption{Bymonthday: []int{-32}}),
        // byyearday out of range
        merge(base, ROption{Byyearday: []int{0}}),
        merge(base, ROption{Byyearday: []int{367}}),
        merge(base, ROption{Byyearday: []int{-367}}),
        // byweekno out of range
        merge(base, ROption{Byweekno: []int{0}}),
        merge(base, ROption{Byweekno: []int{54}}),
        merge(base, ROption{Byweekno: []int{-54}}),
        // bymonth out of range
        merge(base, ROption{Bymonth: []int{0}}),
        merge(base, ROption{Bymonth: []int{13}}),
        // bysetpos out of range
        merge(base, ROption{Bysetpos: []int{0}}),
        merge(base, ROption{Bysetpos: []int{367}}),
        merge(base, ROption{Bysetpos: []int{-367}}),
    }

    for i, opt := range cases {
        if _, err := NewRRule(opt); err == nil {
            t.Errorf("case %d: expected error for out of bounds option, got nil", i)
        }
    }
}

func TestNewRRule_ValidateBounds_ByDayN_OutOfRange(t *testing.T) {
    base := ROption{Freq: MONTHLY, Dtstart: mustUTC(2024, 1, 1, 0, 0, 0)}
    opt := merge(base, ROption{Byweekday: []Weekday{MO.Nth(54)}})
    if _, err := NewRRule(opt); err == nil {
        t.Fatal("expected error for BYDAY N=54, got nil")
    }
}

func TestNewRRule_ValidateBounds_IntervalNegative(t *testing.T) {
    opt := ROption{Freq: DAILY, Dtstart: mustUTC(2024, 1, 1, 0, 0, 0), Interval: -1}
    if _, err := NewRRule(opt); err == nil {
        t.Fatal("expected error for Interval < 0, got nil")
    }
}

func TestNewRRule_CountNegativeBecomesUnlimited(t *testing.T) {
    // Count < 0 becomes 0 (unlimited) in buildRRule
    r, err := NewRRule(ROption{Freq: DAILY, Dtstart: mustUTC(2024, 1, 1, 0, 0, 0), Count: -5})
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if r.count != 0 {
        t.Errorf("expected count to normalize to 0 (unlimited), got %d", r.count)
    }
}

// helpers
func mustUTC(y int, m int, d int, hh int, mm int, ss int) (t time.Time) {
    return time.Date(y, time.Month(m), d, hh, mm, ss, 0, time.UTC)
}

func merge(a, b ROption) ROption {
    out := a
    if len(b.Bysecond) > 0 { out.Bysecond = b.Bysecond }
    if len(b.Byminute) > 0 { out.Byminute = b.Byminute }
    if len(b.Byhour) > 0 { out.Byhour = b.Byhour }
    if len(b.Bymonthday) > 0 { out.Bymonthday = b.Bymonthday }
    if len(b.Byyearday) > 0 { out.Byyearday = b.Byyearday }
    if len(b.Byweekno) > 0 { out.Byweekno = b.Byweekno }
    if len(b.Bymonth) > 0 { out.Bymonth = b.Bymonth }
    if len(b.Bysetpos) > 0 { out.Bysetpos = b.Bysetpos }
    if len(b.Byweekday) > 0 { out.Byweekday = b.Byweekday }
    if b.Interval != 0 { out.Interval = b.Interval }
    return out
}
