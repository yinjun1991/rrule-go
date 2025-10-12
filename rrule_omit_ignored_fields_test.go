package rrule

import (
    "testing"
    "time"
)

// UNTIL zero-value should not appear in RRULE output.
func TestROption_RRuleString_OmitsZeroUntil(t *testing.T) {
    opt := ROption{Freq: DAILY, Until: time.Time{}}
    got := opt.RRuleString()
    want := "FREQ=DAILY"
    if got != want {
        t.Fatalf("expected %q, got %q", want, got)
    }
}

// INTERVAL=0 should be treated as default and omitted from RRULE output on round-trip.
func TestStrToRRule_String_OmitsIntervalZero(t *testing.T) {
    r, err := StrToRRule("RRULE:FREQ=DAILY;INTERVAL=0")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    got := r.String()
    want := "FREQ=DAILY"
    if got != want {
        t.Fatalf("expected %q, got %q", want, got)
    }
}

// WKST=MO (default) should be omitted from RRULE output on round-trip.
func TestStrToRRule_String_OmitsDefaultWkst(t *testing.T) {
    r, err := StrToRRule("RRULE:FREQ=WEEKLY;WKST=MO")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    got := r.String()
    want := "FREQ=WEEKLY"
    if got != want {
        t.Fatalf("expected %q, got %q", want, got)
    }
}

// When COUNT is invalid (negative) but UNTIL is valid, omit COUNT and keep UNTIL.
func TestStrToRRule_String_OmitsNegativeCount_KeepsUntil(t *testing.T) {
    until := time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)
    in := "RRULE:FREQ=DAILY;COUNT=-1;UNTIL=" + until.Format(DateTimeFormat)
    r, err := StrToRRule(in)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    got := r.String()
    want := "FREQ=DAILY;UNTIL=" + until.Format(DateTimeFormat)
    if got != want {
        t.Fatalf("expected %q, got %q", want, got)
    }
}

