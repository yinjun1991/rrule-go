package rrule

import (
    "strings"
    "testing"
    "time"
)

// All-day RRULE output should omit BYHOUR/BYMINUTE/BYSECOND, even if provided.
func TestROption_RRuleString_AllDay_OmitsTimeParts(t *testing.T) {
    opt := ROption{
        Freq:      DAILY,
        AllDay:    true,
        Dtstart:   time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
        Byhour:    []int{9, 15},
        Byminute:  []int{30},
        Bysecond:  []int{0},
    }
    got := opt.RRuleString()
    if strings.Contains(got, "BYHOUR=") || strings.Contains(got, "BYMINUTE=") || strings.Contains(got, "BYSECOND=") {
        t.Fatalf("all-day RRULE must not include time parts, got: %q", got)
    }
}

// Set serialization should also omit time parts for all-day rules.
func TestSet_String_AllDay_OmitsTimePartsInRRULE(t *testing.T) {
    lines := []string{
        "DTSTART;VALUE=DATE:20230101",
        "RRULE:FREQ=DAILY;BYHOUR=9;BYMINUTE=30;BYSECOND=0",
    }
    set, err := StrSliceToRRuleSet(lines)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    out := set.String(true)
    for _, line := range strings.Split(out, "\n") {
        if strings.HasPrefix(line, "RRULE:") {
            if strings.Contains(line, "BYHOUR=") || strings.Contains(line, "BYMINUTE=") || strings.Contains(line, "BYSECOND=") {
                t.Fatalf("all-day Set RRULE must not include time parts, got: %q", line)
            }
        }
    }
}

