package rrule

import (
    "testing"
)

// Confirms current behavior: DTSTART not first is ignored (non-fatal) per current parser.
func TestSetParse_DTSTARTNotFirst_Ignored(t *testing.T) {
    lines := []string{
        "RRULE:FREQ=DAILY;COUNT=1",
        "DTSTART:20250101T000000Z",
    }
    set, err := StrSliceToRRuleSet(lines)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if !set.GetDTStart().IsZero() {
        t.Errorf("expected dtstart to remain zero when DTSTART is not first, got %v", set.GetDTStart())
    }
}

