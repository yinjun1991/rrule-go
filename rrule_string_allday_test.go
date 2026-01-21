package rrule

import (
	"strings"
	"testing"
	"time"
)

// TestRRuleStringAllDayUntil tests RRuleString() handling of UNTIL for all-day events.
func TestRRuleStringAllDayUntil(t *testing.T) {
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
			// All-day event with an UNTIL parameter.
			option := ROption{
				Freq:    DAILY,
				AllDay:  true,
				Dtstart: time.Date(2023, 1, 1, 14, 30, 0, 0, tc.tz),
				Until:   time.Date(2023, 1, 5, 16, 45, 30, 0, tc.tz),
			}

			output := option.RRuleString()
			t.Logf("Timezone %s RRuleString output: %s", tc.name, output)

			// Verify the UNTIL parameter is handled correctly.
			if !strings.Contains(output, "UNTIL=") {
				t.Errorf("Expected UNTIL parameter in output for timezone %s, got: %s", tc.name, output)
			}

			// Verify UNTIL uses DATE format (no time part).
			if !strings.Contains(output, "UNTIL=20230105") {
				t.Errorf("Expected UNTIL=20230105 in output for timezone %s, got: %s", tc.name, output)
			}

			// Verify UNTIL has no time part (no "T" time).
			if strings.Contains(output, "UNTIL=20230105T") {
				t.Errorf("UNTIL should use DATE format (no time part) for all-day events in timezone %s, got: %s", tc.name, output)
			}
		})
	}
}

// TestRRuleStringAllDayConsistency tests RRuleString consistency for all-day events.
func TestRRuleStringAllDayConsistency(t *testing.T) {
	// Create all-day events on the same date across timezones.
	options := []ROption{
		{
			Freq:    WEEKLY,
			Count:   3,
			AllDay:  true,
			Dtstart: time.Date(2023, 6, 15, 8, 30, 0, 0, time.UTC),
		},
		{
			Freq:    WEEKLY,
			Count:   3,
			AllDay:  true,
			Dtstart: time.Date(2023, 6, 15, 14, 45, 0, 0, time.FixedZone("EST", -5*3600)),
		},
		{
			Freq:    WEEKLY,
			Count:   3,
			AllDay:  true,
			Dtstart: time.Date(2023, 6, 15, 22, 15, 0, 0, time.FixedZone("JST", 9*3600)),
		},
	}

	var outputs []string
	for i, option := range options {
		output := option.RRuleString()
		outputs = append(outputs, output)
		t.Logf("Option %d RRuleString output: %s", i+1, output)
	}

	// Verify outputs match since RRuleString omits DTSTART and includes only RRULE.
	expected := "FREQ=WEEKLY;COUNT=3"
	for i, output := range outputs {
		if output != expected {
			t.Errorf("Option %d: expected %s, got %s", i+1, expected, output)
		}
	}
}

// TestRRuleStringAllDayWithUntilTimezone tests UNTIL handling across timezones for all-day events.
func TestRRuleStringAllDayWithUntilTimezone(t *testing.T) {
	// Create all-day events on the same date across timezones with UNTIL.
	options := []ROption{
		{
			Freq:    DAILY,
			AllDay:  true,
			Dtstart: time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
			Until:   time.Date(2023, 1, 3, 15, 30, 0, 0, time.UTC),
		},
		{
			Freq:    DAILY,
			AllDay:  true,
			Dtstart: time.Date(2023, 1, 1, 16, 30, 0, 0, time.FixedZone("EST", -5*3600)),
			Until:   time.Date(2023, 1, 3, 20, 45, 0, 0, time.FixedZone("EST", -5*3600)),
		},
	}

	var outputs []string
	for i, option := range options {
		output := option.RRuleString()
		outputs = append(outputs, output)
		t.Logf("Option %d RRuleString output: %s", i+1, output)
	}

	// Verify UNTIL uses DATE format.
	for i, output := range outputs {
		if !strings.Contains(output, "UNTIL=20230103") {
			t.Errorf("Option %d: expected UNTIL=20230103 in output, got: %s", i+1, output)
		}
		// Verify UNTIL has no time part (no "T" time).
		if strings.Contains(output, "UNTIL=20230103T") {
			t.Errorf("Option %d: UNTIL should use DATE format (no time part) for all-day events, got: %s", i+1, output)
		}
	}
}
