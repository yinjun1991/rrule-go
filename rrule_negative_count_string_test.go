package rrule

import (
    "testing"
    "time"
)

// Negative COUNT should not appear in the generated RRULE string.
func TestROption_RRuleString_OmitsNegativeCount(t *testing.T) {
    opt := ROption{Freq: DAILY, Count: -5}
    got := opt.RRuleString()
    want := "FREQ=DAILY"
    if got != want {
        t.Fatalf("expected %q, got %q", want, got)
    }
}

// When parsing an RRULE with COUNT=-1, String() should omit COUNT.
func TestStrToRRule_String_OmitsNegativeCount_NoDTSTART(t *testing.T) {
    r, err := StrToRRule("RRULE:FREQ=DAILY;COUNT=-1")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    got := r.String()
    want := "FREQ=DAILY"
    if got != want {
        t.Fatalf("expected %q, got %q", want, got)
    }
}

// Same as above but with DTSTART present. COUNT should be omitted from output.
func TestStrToRRule_String_OmitsNegativeCount_WithDTSTART(t *testing.T) {
    dt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
    in := "DTSTART:" + dt.Format(DateTimeFormat) + "\nRRULE:FREQ=DAILY;COUNT=-1"

    r, err := StrToRRule(in)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    got := r.String()
    want := "DTSTART:" + dt.Format(DateTimeFormat) + "\nRRULE:FREQ=DAILY"
    if got != want {
        t.Fatalf("expected %q, got %q", want, got)
    }
}

