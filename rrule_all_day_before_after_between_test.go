package rrule

import (
    "testing"
    "time"
)

func TestAllDay_BeforeAfterBetween(t *testing.T) {
    dt := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
    r, err := NewRRule(ROption{Freq: DAILY, Count: 3, AllDay: true, Dtstart: dt})
    if err != nil { t.Fatal(err) }

    // Occurrences: 2023-01-01, 2023-01-02, 2023-01-03 (midnight UTC)
    d1 := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
    d2 := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)
    d3 := time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC)

    // After at start
    if got := r.After(d1, false); !got.Equal(d2) {
        t.Errorf("After exclude start: got %v want %v", got, d2)
    }
    if got := r.After(d1, true); !got.Equal(d1) {
        t.Errorf("After include start: got %v want %v", got, d1)
    }

    // Before at end
    if got := r.Before(d3, false); !got.Equal(d2) {
        t.Errorf("Before exclude end: got %v want %v", got, d2)
    }
    if got := r.Before(d3, true); !got.Equal(d3) {
        t.Errorf("Before include end: got %v want %v", got, d3)
    }

    // Between window
    // (d1, d3) exclusive should give only d2
    got := r.Between(d1, d3, false)
    if len(got) != 1 || !got[0].Equal(d2) {
        t.Errorf("Between exclusive got %v want [%v]", got, d2)
    }
    // [d1, d3] inclusive should give d1,d2,d3
    got = r.Between(d1, d3, true)
    want := []time.Time{d1, d2, d3}
    if !timesEqual(got, want) {
        t.Errorf("Between inclusive got %v want %v", got, want)
    }
}

