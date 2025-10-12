package rrule

import (
    "strings"
    "testing"
)

// Helper: extract RRULE line from set.String(true)
func getRRULELine(s string) string {
    for _, line := range strings.Split(s, "\n") {
        if strings.HasPrefix(line, "RRULE:") {
            return line
        }
    }
    return ""
}

func TestSet_String_OmitsIgnoredParams_AllDay(t *testing.T) {
    lines := []string{
        "DTSTART;VALUE=DATE:20240101",
        "RRULE:FREQ=DAILY;COUNT=-1;INTERVAL=0;WKST=MO;BYHOUR=9;BYMINUTE=30;BYSECOND=0",
    }
    set, err := StrSliceToRRuleSet(lines)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    out := set.String(true)
    rline := getRRULELine(out)
    if rline == "" {
        t.Fatalf("RRULE line not found in: %q", out)
    }
    // All-day: drop COUNT<=0, INTERVAL=0, default WKST=MO, and time parts
    if want := "RRULE:FREQ=DAILY"; rline != want {
        t.Fatalf("expected %q, got %q", want, rline)
    }
}

func TestSet_String_OmitsIgnoredParams_Timed(t *testing.T) {
    lines := []string{
        "DTSTART:20240101T100000Z",
        "RRULE:FREQ=DAILY;COUNT=-1;INTERVAL=0;WKST=MO;BYHOUR=9;BYMINUTE=30;BYSECOND=0",
    }
    set, err := StrSliceToRRuleSet(lines)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    out := set.String(true)
    rline := getRRULELine(out)
    if rline == "" {
        t.Fatalf("RRULE line not found in: %q", out)
    }
    // Timed: drop COUNT<=0, INTERVAL=0, default WKST=MO, keep time parts
    if want := "RRULE:FREQ=DAILY;BYHOUR=9;BYMINUTE=30;BYSECOND=0"; rline != want {
        t.Fatalf("expected %q, got %q", want, rline)
    }
}

func TestSet_String_OmitsNegativeCount_KeepsUntil(t *testing.T) {
    lines := []string{
        "DTSTART:20240101T000000Z",
        "RRULE:FREQ=DAILY;COUNT=-1;UNTIL=20240103T000000Z",
    }
    set, err := StrSliceToRRuleSet(lines)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    rline := getRRULELine(set.String(true))
    if want := "RRULE:FREQ=DAILY;UNTIL=20240103T000000Z"; rline != want {
        t.Fatalf("expected %q, got %q", want, rline)
    }
}

