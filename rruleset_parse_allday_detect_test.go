package rrule

import (
    "strings"
    "testing"
    "time"
)

// Ensures DTSTART;VALUE=DATE auto-detects all-day in StrToRRuleSet and propagates to RRULE
func TestStrToRRuleSet_AllDayDetectionAndPropagation(t *testing.T) {
    setStr := strings.Join([]string{
        "DTSTART;VALUE=DATE:20230901",
        "RRULE:FREQ=DAILY;COUNT=3;UNTIL=20230903",
        "RDATE;VALUE=DATE:20230902",
        "EXDATE;VALUE=DATE:20230901",
    }, "\n")

    set, err := StrToRRuleSet(setStr)
    if err != nil {
        t.Fatalf("StrToRRuleSet failed: %v", err)
    }

    // All-day should be auto-detected
    if !set.IsAllDay() {
        t.Fatal("expected set to be all-day after parsing DTSTART;VALUE=DATE")
    }

    // RRULE serialization should use DATE UNTIL and DTSTART should use VALUE=DATE
    out := set.String(true)
    if !strings.Contains(out, "DTSTART;VALUE=DATE:20230901") {
        t.Errorf("expected DTSTART;VALUE=DATE in output, got: %s", out)
    }
    if !strings.Contains(out, "RRULE:FREQ=DAILY;COUNT=3;UNTIL=20230903") {
        t.Errorf("expected UNTIL as DATE in RRULE, got: %s", out)
    }
    if strings.Contains(out, "UNTIL=20230903T") {
        t.Errorf("UNTIL must not include time part for all-day, got: %s", out)
    }

    // Verify occurrences: DTSTART=2023-09-01, COUNT=3 -> 1st, 2nd, 3rd; EXDATE removes 1st; RDATE adds 2nd (duplicate)
    want := []time.Time{
        time.Date(2023, 9, 2, 0, 0, 0, 0, time.UTC),
        time.Date(2023, 9, 3, 0, 0, 0, 0, time.UTC),
    }
    got := set.All()
    if !timesEqual(got, want) {
        t.Errorf("occurrences mismatch, got %v, want %v", got, want)
    }
}

