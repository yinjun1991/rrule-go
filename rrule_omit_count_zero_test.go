package rrule

import "testing"

func TestStrToRRule_String_OmitsCountZero(t *testing.T) {
    r, err := StrToRRule("RRULE:FREQ=DAILY;COUNT=0")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if got, want := r.String(), "FREQ=DAILY"; got != want {
        t.Fatalf("expected %q, got %q", want, got)
    }
}

func TestSet_String_OmitsCountZero(t *testing.T) {
    lines := []string{
        "DTSTART:20240101T000000Z",
        "RRULE:FREQ=DAILY;COUNT=0",
    }
    set, err := StrSliceToRRuleSet(lines)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    out := set.String(true)
    if out == "" {
        t.Fatal("unexpected empty set output")
    }
    if idx := indexOfRRULE(out); idx == -1 {
        t.Fatalf("RRULE line not found in: %q", out)
    } else {
        rline := rrLineAt(out, idx)
        if rline != "RRULE:FREQ=DAILY" {
            t.Fatalf("expected RRULE to omit COUNT=0, got: %q", rline)
        }
    }
}

// small helpers for the test file
func indexOfRRULE(s string) int {
    i := 0
    for _, line := range splitLines(s) {
        if len(line) >= 6 && line[:6] == "RRULE:" {
            return i
        }
        i++
    }
    return -1
}

func rrLineAt(s string, idx int) string {
    i := 0
    for _, line := range splitLines(s) {
        if i == idx {
            return line
        }
        i++
    }
    return ""
}

func splitLines(s string) []string {
    var lines []string
    start := 0
    for i := 0; i < len(s); i++ {
        if s[i] == '\n' {
            lines = append(lines, s[start:i])
            start = i + 1
        }
    }
    lines = append(lines, s[start:])
    return lines
}

