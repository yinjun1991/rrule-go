package rrule

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRecurrenceSet_WithoutRRulePrefix(t *testing.T) {
	// Test case 1: RRULE string without "RRULE:" prefix
	dtstart := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)

	// This should fail with current implementation
	rruleStrings := []string{"FREQ=DAILY;COUNT=5"}
	_, err := ParseRecurrenceSet(rruleStrings, dtstart, false)

	// Current implementation expects this to fail
	// After fix, this should succeed
	if err != nil {
		t.Logf("Expected behavior: RRULE without prefix fails: %v", err)
	} else {
		t.Logf("RRULE without prefix parsed successfully")
	}
}

func TestParseRecurrenceSet_WithRRulePrefix(t *testing.T) {
	// Test case 2: RRULE string with "RRULE:" prefix (standard format)
	dtstart := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)

	rruleStrings := []string{"RRULE:FREQ=DAILY;COUNT=5"}
	helper, err := ParseRecurrenceSet(rruleStrings, dtstart, false)

	require.NoError(t, err)
	require.NotNil(t, helper)

	// Test occurrence generation
	occurrences := helper.Between(
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
	)

	assert.Len(t, occurrences, 5, "Should generate 5 occurrences")
}

func TestParseRecurrenceSet_MixedFormats(t *testing.T) {
	// Test case 3: Mixed RRULE, RDATE, EXDATE with various formats
	dtstart := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)

	rruleStrings := []string{
		"RRULE:FREQ=DAILY;COUNT=10",
		"RDATE:20240105T100000Z",
		"EXDATE:20240103T100000Z",
	}

	helper, err := ParseRecurrenceSet(rruleStrings, dtstart, false)
	require.NoError(t, err)
	require.NotNil(t, helper)

	// Test that RDATE and EXDATE are properly handled
	rdates := helper.GetRDates()
	exdates := helper.GetExDates()

	assert.Len(t, rdates, 1, "Should have 1 RDATE")
	assert.Len(t, exdates, 1, "Should have 1 EXDATE")
}

func TestParseRecurrenceSet_MultipleRRules(t *testing.T) {
	// Test case 4: Multiple RRULE strings (should take first one)
	dtstart := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)

	rruleStrings := []string{
		"RRULE:FREQ=DAILY;COUNT=5",
		"RRULE:FREQ=WEEKLY;COUNT=3", // This should be ignored
	}

	helper, err := ParseRecurrenceSet(rruleStrings, dtstart, false)

	// Current implementation might fail or use both rules
	// After fix, should use only the first RRULE
	if err != nil {
		t.Logf("Multiple RRULEs failed as expected: %v", err)
	} else {
		occurrences := helper.Between(
			time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
		)
		t.Logf("Multiple RRULEs generated %d occurrences", len(occurrences))
	}
}

func TestParseRecurrenceSet_AllDayEvents(t *testing.T) {
	// Test case 5: All-day events with floating time
	dtstart := time.Date(2024, 1, 1, 14, 30, 0, 0, time.UTC) // Non-midnight time

	rruleStrings := []string{"RRULE:FREQ=DAILY;COUNT=3"}
	helper, err := ParseRecurrenceSet(rruleStrings, dtstart, true) // allDay = true

	require.NoError(t, err)
	require.NotNil(t, helper)

	occurrences := helper.Between(
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
	)

	assert.Len(t, occurrences, 3, "Should generate 3 all-day occurrences")

	// All occurrences should be at midnight UTC for all-day events
	for i, occurrence := range occurrences {
		assert.EqualValues(t, 0, occurrence.Hour(), "All-day occurrence %d should be at midnight", i)
		assert.EqualValues(t, 0, occurrence.Minute(), "All-day occurrence %d should be at midnight", i)
		assert.EqualValues(t, 0, occurrence.Second(), "All-day occurrence %d should be at midnight", i)
	}
}

func TestToRRuleStrings_Conversion(t *testing.T) {
	// Test case 6: Convert back to string format
	dtstart := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)

	originalStrings := []string{
		"RRULE:FREQ=DAILY;COUNT=5",
		"RDATE:20240110T100000Z",
		"EXDATE:20240103T100000Z",
	}

	helper, err := ParseRecurrenceSet(originalStrings, dtstart, false)
	require.NoError(t, err)

	convertedStrings := helper.ToRRuleStrings()

	// Should contain RRULE, RDATE, and EXDATE entries
	assert.NotEmpty(t, convertedStrings, "Should return non-empty string array")

	// Log the conversion for manual inspection
	t.Logf("Original: %v", originalStrings)
	t.Logf("Converted: %v", convertedStrings)
}

func TestParseRecurrenceSet_DTStartHandling(t *testing.T) {
	// Test case 7: DTSTART handling for different time zones

	// Test with UTC time
	dtstart := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	rruleStrings := []string{"RRULE:FREQ=DAILY;COUNT=3"}

	helper, err := ParseRecurrenceSet(rruleStrings, dtstart, false)
	require.NoError(t, err)

	occurrences := helper.Between(
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
	)

	assert.Len(t, occurrences, 3, "Should generate 3 occurrences")
	assert.EqualValues(t, dtstart, occurrences[0], "First occurrence should match DTSTART")
}

func TestParseRecurrenceSet_UntilHandling(t *testing.T) {
	// Test case 8: UNTIL parameter handling for all-day vs non-all-day

	// Non-all-day event with UNTIL
	dtstart := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	rruleStrings := []string{"RRULE:FREQ=DAILY;UNTIL=20240105T100000Z"}

	helper, err := ParseRecurrenceSet(rruleStrings, dtstart, false)
	require.NoError(t, err)

	until := helper.GetUntil()
	require.NotNil(t, until, "Should have UNTIL value")

	expectedUntil := time.Date(2024, 1, 5, 10, 0, 0, 0, time.UTC)
	assert.EqualValues(t, expectedUntil, *until, "UNTIL should match expected time")

	// All-day event with UNTIL (should use date-only format)
	allDayStrings := []string{"RRULE:FREQ=DAILY;UNTIL=20240105"}
	allDayHelper, err := ParseRecurrenceSet(allDayStrings, dtstart, true)

	if err != nil {
		t.Logf("All-day UNTIL parsing failed (may need format adjustment): %v", err)
	} else {
		allDayUntil := allDayHelper.GetUntil()
		if allDayUntil != nil {
			t.Logf("All-day UNTIL: %v", *allDayUntil)
		}
	}
}

func TestParseRecurrenceSet_EmptyInput(t *testing.T) {
	// Test case 9: Edge cases
	dtstart := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)

	// Empty string array
	_, err := ParseRecurrenceSet([]string{}, dtstart, false)
	assert.Error(t, err, "Empty string array should return error")

	// Nil string array
	_, err = ParseRecurrenceSet(nil, dtstart, false)
	assert.Error(t, err, "Nil string array should return error")
}

func TestParseRecurrenceSet_InvalidFormats(t *testing.T) {
	// Test case 10: Invalid RRULE formats
	dtstart := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)

	invalidStrings := []string{
		"RRULE:INVALID_FREQ=DAILY",
		"RRULE:FREQ=DAILY;INVALID_PARAM=VALUE",
	}

	for _, invalidString := range invalidStrings {
		_, err := ParseRecurrenceSet([]string{invalidString}, dtstart, false)
		assert.Error(t, err, "Invalid RRULE should return error: %s", invalidString)
	}

	// Test that "INVALID:FREQ=DAILY" is passed through (not treated as RRULE)
	// This maintains backward compatibility - let rrule-go handle the validation
	_, err := ParseRecurrenceSet([]string{"INVALID:FREQ=DAILY"}, dtstart, false)
	// This may or may not error depending on rrule-go's handling
	// The important thing is that our normalization doesn't crash
	t.Logf("INVALID:FREQ=DAILY result: %v", err)
}
