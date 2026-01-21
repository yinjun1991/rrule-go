package rrule

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRecurrenceSetHelper_GetCount(t *testing.T) {
	utc := time.UTC
	t.Run("Count present", func(t *testing.T) {
		rr := []string{"RRULE:FREQ=DAILY;COUNT=7"}
		helper, err := ParseRecurrenceSet(rr, time.Date(2025, 1, 1, 0, 0, 0, 0, utc), false)
		assert.NoError(t, err)
		assert.Equal(t, 7, helper.GetCount())
	})

	t.Run("No count returns 0", func(t *testing.T) {
		rr := []string{"RRULE:FREQ=WEEKLY;BYDAY=MO,WE"}
		helper, err := ParseRecurrenceSet(rr, time.Date(2025, 1, 6, 0, 0, 0, 0, utc), false)
		assert.NoError(t, err)
		assert.Equal(t, 0, helper.GetCount())
	})

	t.Run("All-day with count", func(t *testing.T) {
		rr := []string{"RRULE:FREQ=DAILY;COUNT=2"}
		helper, err := ParseRecurrenceSet(rr, time.Date(2025, 2, 1, 0, 0, 0, 0, utc), true)
		assert.NoError(t, err)
		assert.Equal(t, 2, helper.GetCount())
	})
}

func TestParseRecurrenceSet_DTStartMismatch(t *testing.T) {
	t.Run("Timed event mismatch", func(t *testing.T) {
		rr := []string{"RRULE:FREQ=MONTHLY;BYMONTHDAY=21"}
		helper, err := ParseRecurrenceSet(rr, time.Date(2025, time.October, 22, 9, 0, 0, 0, time.UTC), false)
		assert.Nil(t, err)
		err = helper.ValidateDTStartAlignment(false)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not align with recurrence set")
	})

	t.Run("All-day event mismatch", func(t *testing.T) {
		rr := []string{"RRULE:FREQ=WEEKLY;BYDAY=MO"}
		helper, err := ParseRecurrenceSet(rr, time.Date(2025, time.October, 22, 0, 0, 0, 0, time.UTC), true)
		assert.Nil(t, err)
		err = helper.ValidateDTStartAlignment(false)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not align with recurrence set")
	})
}

func TestParseRecurrenceSet_DTStartIrregularFirstInstance(t *testing.T) {
	dtstart := time.Date(2024, time.January, 1, 9, 0, 0, 0, time.UTC)
	rr := []string{
		"RRULE:FREQ=MONTHLY;BYDAY=FR;BYSETPOS=1",
		"RDATE:20240101T090000Z",
	}

	helper, err := ParseRecurrenceSet(rr, dtstart, false)
	assert.NoError(t, err)
	assert.NotNil(t, helper)

	first := helper.set.After(dtstart, true)
	assert.False(t, first.IsZero())
	assert.Equal(t, dtstart, first)
}

func TestParseRecurrenceSet_DTStartMissingFromSet(t *testing.T) {
	dtstart := time.Date(2024, time.January, 1, 9, 0, 0, 0, time.UTC)
	rr := []string{
		"RRULE:FREQ=MONTHLY;BYDAY=FR;BYSETPOS=1",
		"RDATE:20231225T090000Z",
	}

	helper, err := ParseRecurrenceSet(rr, dtstart, false)
	assert.Nil(t, err)
	err = helper.ValidateDTStartAlignment(false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not align with recurrence set")
	assert.Contains(t, err.Error(), "2024-01-05")
}

func TestParseRecurrenceSetIgnoreEXDATE_AllowsExcludedDTStart(t *testing.T) {
	loc, err := time.LoadLocation("America/Los_Angeles")
	assert.NoError(t, err)

	dtstart := time.Date(2025, time.September, 1, 9, 30, 0, 0, loc)
	rr := []string{
		"RRULE:FREQ=DAILY;COUNT=5",
		"EXDATE;TZID=America/Los_Angeles:20250901T093000",
	}

	helper, parseErr := ParseRecurrenceSet(rr, dtstart, false)
	assert.NoError(t, parseErr)
	assert.NotNil(t, helper)

	next, nextErr := helper.Next(dtstart)
	assert.NoError(t, nextErr)
	assert.Equal(t, time.Date(2025, time.September, 2, 9, 30, 0, 0, loc), next.In(loc))
}

func TestParseRecurrenceSet_EXDATEWithTZID_ExcludedInstance(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	assert.NoError(t, err)

	dtstart := time.Date(2025, time.December, 20, 21, 0, 0, 0, loc)
	excludedUTC := time.Date(2025, time.December, 20, 13, 0, 0, 0, time.UTC)
	rr := []string{
		"RRULE:FREQ=WEEKLY;BYDAY=SA",
		"EXDATE;TZID=Asia/Shanghai:20251220T210000",
	}

	helper, parseErr := ParseRecurrenceSet(rr, dtstart, false)
	assert.NoError(t, parseErr)
	assert.NotNil(t, helper)

	occurrences, nextErr := helper.NextN(dtstart.Add(-time.Second), 2)
	assert.NoError(t, nextErr)
	assert.Len(t, occurrences, 2)
	var occurrences2 []time.Time
	for _, occ := range occurrences {
		occurrences2 = append(occurrences2, occ.UTC())
	}
	assert.NotContains(t, occurrences2, dtstart.UTC())
	assert.NotContains(t, occurrences2, excludedUTC.UTC())
	assert.Equal(t, time.Date(2025, time.December, 27, 21, 0, 0, 0, loc).UTC(), occurrences2[0])
}

func TestParseRecurrenceSet_BYDAY1SA_Expansion(t *testing.T) {
	dtstart := time.Date(2025, time.March, 1, 10, 0, 0, 0, time.UTC)
	rr := []string{"RRULE:FREQ=MONTHLY;BYDAY=1SA;COUNT=3"}

	helper, err := ParseRecurrenceSet(rr, dtstart, false)
	assert.NoError(t, err)
	assert.NotNil(t, helper)

	occs, err := helper.NextN(dtstart.Add(-time.Second), 3)
	assert.NoError(t, err)
	assert.Len(t, occs, 3)
	assert.Equal(t, dtstart, occs[0])
	assert.Equal(t, time.Date(2025, time.April, 5, 10, 0, 0, 0, time.UTC), occs[1])
	assert.Equal(t, time.Date(2025, time.May, 3, 10, 0, 0, 0, time.UTC), occs[2])
}

func TestParseRecurrenceSet_BYSECOND_Expansion(t *testing.T) {
	dtstart := time.Date(2025, time.March, 1, 10, 0, 5, 0, time.UTC)
	rr := []string{"RRULE:FREQ=MINUTELY;BYSECOND=5,20;COUNT=4"}

	helper, err := ParseRecurrenceSet(rr, dtstart, false)
	assert.NoError(t, err)
	assert.NotNil(t, helper)

	occs, err := helper.NextN(dtstart.Add(-time.Second), 4)
	assert.NoError(t, err)
	assert.Len(t, occs, 4)
	assert.Equal(t, time.Date(2025, time.March, 1, 10, 0, 5, 0, time.UTC), occs[0])
	assert.Equal(t, time.Date(2025, time.March, 1, 10, 0, 20, 0, time.UTC), occs[1])
	assert.Equal(t, time.Date(2025, time.March, 1, 10, 1, 5, 0, time.UTC), occs[2])
	assert.Equal(t, time.Date(2025, time.March, 1, 10, 1, 20, 0, time.UTC), occs[3])
}

func TestParseRecurrenceSet_BYMINUTE_Expansion(t *testing.T) {
	dtstart := time.Date(2025, time.March, 1, 10, 15, 0, 0, time.UTC)
	rr := []string{"RRULE:FREQ=HOURLY;BYMINUTE=15,45;COUNT=4"}

	helper, err := ParseRecurrenceSet(rr, dtstart, false)
	assert.NoError(t, err)
	assert.NotNil(t, helper)

	occs, err := helper.NextN(dtstart.Add(-time.Second), 4)
	assert.NoError(t, err)
	assert.Len(t, occs, 4)
	assert.Equal(t, time.Date(2025, time.March, 1, 10, 15, 0, 0, time.UTC), occs[0])
	assert.Equal(t, time.Date(2025, time.March, 1, 10, 45, 0, 0, time.UTC), occs[1])
	assert.Equal(t, time.Date(2025, time.March, 1, 11, 15, 0, 0, time.UTC), occs[2])
	assert.Equal(t, time.Date(2025, time.March, 1, 11, 45, 0, 0, time.UTC), occs[3])
}

func TestParseRecurrenceSet_BYHOUR_Expansion(t *testing.T) {
	dtstart := time.Date(2025, time.March, 1, 9, 0, 0, 0, time.UTC)
	rr := []string{"RRULE:FREQ=DAILY;BYHOUR=9,17;COUNT=4"}

	helper, err := ParseRecurrenceSet(rr, dtstart, false)
	assert.NoError(t, err)
	assert.NotNil(t, helper)

	occs, err := helper.NextN(dtstart.Add(-time.Second), 4)
	assert.NoError(t, err)
	assert.Len(t, occs, 4)
	assert.Equal(t, time.Date(2025, time.March, 1, 9, 0, 0, 0, time.UTC), occs[0])
	assert.Equal(t, time.Date(2025, time.March, 1, 17, 0, 0, 0, time.UTC), occs[1])
	assert.Equal(t, time.Date(2025, time.March, 2, 9, 0, 0, 0, time.UTC), occs[2])
	assert.Equal(t, time.Date(2025, time.March, 2, 17, 0, 0, 0, time.UTC), occs[3])
}

func TestParseRecurrenceSet_BYMONTHDAY_Expansion(t *testing.T) {
	dtstart := time.Date(2025, time.March, 15, 10, 0, 0, 0, time.UTC)
	rr := []string{"RRULE:FREQ=MONTHLY;BYMONTHDAY=15;COUNT=3"}

	helper, err := ParseRecurrenceSet(rr, dtstart, false)
	assert.NoError(t, err)
	assert.NotNil(t, helper)

	occs, err := helper.NextN(dtstart.Add(-time.Second), 3)
	assert.NoError(t, err)
	assert.Len(t, occs, 3)
	assert.Equal(t, time.Date(2025, time.March, 15, 10, 0, 0, 0, time.UTC), occs[0])
	assert.Equal(t, time.Date(2025, time.April, 15, 10, 0, 0, 0, time.UTC), occs[1])
	assert.Equal(t, time.Date(2025, time.May, 15, 10, 0, 0, 0, time.UTC), occs[2])
}

func TestParseRecurrenceSet_BYYEARDAY_Expansion(t *testing.T) {
	dtstart := time.Date(2025, time.January, 1, 8, 0, 0, 0, time.UTC)
	rr := []string{"RRULE:FREQ=YEARLY;BYYEARDAY=1;COUNT=3"}

	helper, err := ParseRecurrenceSet(rr, dtstart, false)
	assert.NoError(t, err)
	assert.NotNil(t, helper)

	occs, err := helper.NextN(dtstart.Add(-time.Second), 3)
	assert.NoError(t, err)
	assert.Len(t, occs, 3)
	assert.Equal(t, time.Date(2025, time.January, 1, 8, 0, 0, 0, time.UTC), occs[0])
	assert.Equal(t, time.Date(2026, time.January, 1, 8, 0, 0, 0, time.UTC), occs[1])
	assert.Equal(t, time.Date(2027, time.January, 1, 8, 0, 0, 0, time.UTC), occs[2])
}

func TestParseRecurrenceSet_BYWEEKNO_WKST_Expansion(t *testing.T) {
	dtstart := time.Date(2025, time.January, 5, 9, 0, 0, 0, time.UTC) // Sunday
	rr := []string{"RRULE:FREQ=YEARLY;BYWEEKNO=2;BYDAY=SU;WKST=SU;COUNT=3"}

	helper, err := ParseRecurrenceSet(rr, dtstart, false)
	assert.NoError(t, err)
	assert.NotNil(t, helper)

	occs, err := helper.NextN(dtstart.Add(-time.Second), 3)
	assert.NoError(t, err)
	assert.Len(t, occs, 3)
	assert.Equal(t, time.Date(2025, time.January, 5, 9, 0, 0, 0, time.UTC), occs[0])
	assert.Equal(t, time.Date(2026, time.January, 11, 9, 0, 0, 0, time.UTC), occs[1])
	assert.Equal(t, time.Date(2027, time.January, 10, 9, 0, 0, 0, time.UTC), occs[2])
}

func TestParseRecurrenceSet_BYMONTH_Expansion(t *testing.T) {
	dtstart := time.Date(2025, time.February, 10, 9, 0, 0, 0, time.UTC)
	rr := []string{"RRULE:FREQ=YEARLY;BYMONTH=2,6;BYMONTHDAY=10;COUNT=4"}

	helper, err := ParseRecurrenceSet(rr, dtstart, false)
	assert.NoError(t, err)
	assert.NotNil(t, helper)

	occs, err := helper.NextN(dtstart.Add(-time.Second), 4)
	assert.NoError(t, err)
	assert.Len(t, occs, 4)
	assert.Equal(t, time.Date(2025, time.February, 10, 9, 0, 0, 0, time.UTC), occs[0])
	assert.Equal(t, time.Date(2025, time.June, 10, 9, 0, 0, 0, time.UTC), occs[1])
	assert.Equal(t, time.Date(2026, time.February, 10, 9, 0, 0, 0, time.UTC), occs[2])
	assert.Equal(t, time.Date(2026, time.June, 10, 9, 0, 0, 0, time.UTC), occs[3])
}

func TestParseRecurrenceSet_BYDAY_Negative_LastSaturday(t *testing.T) {
	dtstart := time.Date(2025, time.March, 29, 10, 0, 0, 0, time.UTC)
	rr := []string{"RRULE:FREQ=MONTHLY;BYDAY=-1SA;COUNT=3"}

	helper, err := ParseRecurrenceSet(rr, dtstart, false)
	assert.NoError(t, err)
	assert.NotNil(t, helper)

	occs, err := helper.NextN(dtstart.Add(-time.Second), 3)
	assert.NoError(t, err)
	assert.Len(t, occs, 3)
	assert.Equal(t, time.Date(2025, time.March, 29, 10, 0, 0, 0, time.UTC), occs[0])
	assert.Equal(t, time.Date(2025, time.April, 26, 10, 0, 0, 0, time.UTC), occs[1])
	assert.Equal(t, time.Date(2025, time.May, 31, 10, 0, 0, 0, time.UTC), occs[2])
}

func TestParseRecurrenceSet_BYDAY_BYSETPOS_Negative_LastSaturday(t *testing.T) {
	dtstart := time.Date(2025, time.March, 29, 10, 0, 0, 0, time.UTC)
	rr := []string{"RRULE:FREQ=MONTHLY;BYDAY=SA;BYSETPOS=-1;COUNT=3"}

	helper, err := ParseRecurrenceSet(rr, dtstart, false)
	assert.NoError(t, err)
	assert.NotNil(t, helper)

	occs, err := helper.NextN(dtstart.Add(-time.Second), 3)
	assert.NoError(t, err)
	assert.Len(t, occs, 3)
	assert.Equal(t, time.Date(2025, time.March, 29, 10, 0, 0, 0, time.UTC), occs[0])
	assert.Equal(t, time.Date(2025, time.April, 26, 10, 0, 0, 0, time.UTC), occs[1])
	assert.Equal(t, time.Date(2025, time.May, 31, 10, 0, 0, 0, time.UTC), occs[2])
}

func TestParseRecurrenceSet_BYMONTHDAY_Negative_LastDay(t *testing.T) {
	dtstart := time.Date(2025, time.March, 31, 8, 30, 0, 0, time.UTC)
	rr := []string{"RRULE:FREQ=MONTHLY;BYMONTHDAY=-1;COUNT=3"}

	helper, err := ParseRecurrenceSet(rr, dtstart, false)
	assert.NoError(t, err)
	assert.NotNil(t, helper)

	occs, err := helper.NextN(dtstart.Add(-time.Second), 3)
	assert.NoError(t, err)
	assert.Len(t, occs, 3)
	assert.Equal(t, time.Date(2025, time.March, 31, 8, 30, 0, 0, time.UTC), occs[0])
	assert.Equal(t, time.Date(2025, time.April, 30, 8, 30, 0, 0, time.UTC), occs[1])
	assert.Equal(t, time.Date(2025, time.May, 31, 8, 30, 0, 0, time.UTC), occs[2])
}

func TestParseRecurrenceSet_BYMONTHDAY_Negative_LastDay_AllDay(t *testing.T) {
	dtstart := time.Date(2025, time.March, 31, 0, 0, 0, 0, time.UTC)
	rr := []string{"RRULE:FREQ=MONTHLY;BYMONTHDAY=-1;COUNT=3"}

	helper, err := ParseRecurrenceSet(rr, dtstart, true)
	assert.NoError(t, err)
	assert.NotNil(t, helper)

	occs, err := helper.NextN(dtstart.Add(-time.Second), 3)
	assert.NoError(t, err)
	assert.Len(t, occs, 3)
	assert.Equal(t, time.Date(2025, time.March, 31, 0, 0, 0, 0, time.UTC), occs[0])
	assert.Equal(t, time.Date(2025, time.April, 30, 0, 0, 0, 0, time.UTC), occs[1])
	assert.Equal(t, time.Date(2025, time.May, 31, 0, 0, 0, 0, time.UTC), occs[2])
}

func TestParseRecurrenceSet_BYYEARDAY_Negative_LastDayOfYear(t *testing.T) {
	dtstart := time.Date(2025, time.December, 31, 10, 0, 0, 0, time.UTC)
	rr := []string{"RRULE:FREQ=YEARLY;BYYEARDAY=-1;COUNT=3"}

	helper, err := ParseRecurrenceSet(rr, dtstart, false)
	assert.NoError(t, err)
	assert.NotNil(t, helper)

	occs, err := helper.NextN(dtstart.Add(-time.Second), 3)
	assert.NoError(t, err)
	assert.Len(t, occs, 3)
	assert.Equal(t, time.Date(2025, time.December, 31, 10, 0, 0, 0, time.UTC), occs[0])
	assert.Equal(t, time.Date(2026, time.December, 31, 10, 0, 0, 0, time.UTC), occs[1])
	assert.Equal(t, time.Date(2027, time.December, 31, 10, 0, 0, 0, time.UTC), occs[2])
}

func TestParseRecurrenceSet_WithEXDATE_SkipsExcludedOccurrences(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	assert.NoError(t, err)

	// DTSTART on Monday 2024-01-01 09:00 (which is a Monday)
	dtstart := time.Date(2024, time.January, 1, 9, 0, 0, 0, loc)

	// Test case: Weekly recurring on Mon, Wed, Fri with COUNT=10, excluding Wed 2024-01-03
	rr := []string{
		"FREQ=WEEKLY;BYDAY=MO,WE,FR;COUNT=10",
		"EXDATE;TZID=America/New_York:20240103T090000",
	}

	helper, parseErr := ParseRecurrenceSet(rr, dtstart, false)
	assert.NoError(t, parseErr)
	assert.NotNil(t, helper)

	// Get first 5 occurrences after dtstart (not including dtstart)
	occs, nextErr := helper.NextN(dtstart, 5)
	assert.NoError(t, nextErr)
	assert.Len(t, occs, 5)

	// Expected occurrences (skipping 2024-01-03 which is excluded):
	// Mon 2024-01-01 09:00 (dtstart - included in NextN when called with dtstart)
	// Wed 2024-01-03 09:00 - EXCLUDED by EXDATE
	// Fri 2024-01-05 09:00
	// Mon 2024-01-08 09:00
	// Wed 2024-01-10 09:00
	// Fri 2024-01-12 09:00
	// Mon 2024-01-15 09:00
	assert.Equal(t, time.Date(2024, time.January, 5, 9, 0, 0, 0, loc), occs[0])  // Fri
	assert.Equal(t, time.Date(2024, time.January, 8, 9, 0, 0, 0, loc), occs[1])  // Mon
	assert.Equal(t, time.Date(2024, time.January, 10, 9, 0, 0, 0, loc), occs[2]) // Wed
	assert.Equal(t, time.Date(2024, time.January, 12, 9, 0, 0, 0, loc), occs[3]) // Fri
	assert.Equal(t, time.Date(2024, time.January, 15, 9, 0, 0, 0, loc), occs[4]) // Mon

	// Verify that EXDATE is properly stored
	exdates := helper.GetExDates()
	assert.Len(t, exdates, 1)
	expectedExdate := time.Date(2024, time.January, 3, 9, 0, 0, 0, loc)
	assert.Equal(t, expectedExdate, exdates[0])

	// Verify the RRule can be serialized back correctly
	rruleStrings := helper.ToRRuleStrings()
	assert.Contains(t, rruleStrings[0], "FREQ=WEEKLY")
	assert.Contains(t, rruleStrings[0], "BYDAY=MO,WE,FR")
	assert.Contains(t, rruleStrings[0], "COUNT=10")

	// Check if EXDATE is preserved in output
	var hasExdate bool
	for _, line := range rruleStrings {
		if strings.Contains(line, "EXDATE") {
			hasExdate = true
			break
		}
	}
	assert.True(t, hasExdate, "EXDATE should be preserved in ToRRuleStrings output")
}

func TestParseRecurrenceSet_WithEXDATE_RFC5545Compliance(t *testing.T) {
	t.Run("EXDATE with UTC format", func(t *testing.T) {
		dtstart := time.Date(2024, time.January, 1, 9, 0, 0, 0, time.UTC)
		rr := []string{
			"FREQ=DAILY;COUNT=5",
			"EXDATE:20240103T090000Z",
		}

		helper, err := ParseRecurrenceSet(rr, dtstart, false)
		assert.NoError(t, err)
		assert.NotNil(t, helper)

		// COUNT=5 generates: Jan 1 (dtstart), Jan 2, Jan 3 (excluded), Jan 4, Jan 5
		// NextN(dtstart, n) returns occurrences AFTER dtstart, so: Jan 2, Jan 4, Jan 5 (3 occurrences)
		occs, err := helper.NextN(dtstart, 5)
		assert.NoError(t, err)
		assert.Len(t, occs, 3)

		// Should skip 2024-01-03
		assert.Equal(t, time.Date(2024, time.January, 2, 9, 0, 0, 0, time.UTC), occs[0])
		assert.Equal(t, time.Date(2024, time.January, 4, 9, 0, 0, 0, time.UTC), occs[1])
		assert.Equal(t, time.Date(2024, time.January, 5, 9, 0, 0, 0, time.UTC), occs[2])
	})

	t.Run("EXDATE with TZID format", func(t *testing.T) {
		loc, err := time.LoadLocation("America/Los_Angeles")
		assert.NoError(t, err)

		dtstart := time.Date(2024, time.January, 1, 14, 30, 0, 0, loc)
		rr := []string{
			"FREQ=DAILY;COUNT=5",
			"EXDATE;TZID=America/Los_Angeles:20240103T143000",
		}

		helper, parseErr := ParseRecurrenceSet(rr, dtstart, false)
		assert.NoError(t, parseErr)
		assert.NotNil(t, helper)

		// COUNT=5 generates: Jan 1 (dtstart), Jan 2, Jan 3 (excluded), Jan 4, Jan 5
		// NextN(dtstart, n) returns occurrences AFTER dtstart, so: Jan 2, Jan 4, Jan 5 (3 occurrences)
		occs, nextErr := helper.NextN(dtstart, 5)
		assert.NoError(t, nextErr)
		assert.Len(t, occs, 3)

		// Should skip 2024-01-03
		assert.Equal(t, time.Date(2024, time.January, 2, 14, 30, 0, 0, loc), occs[0])
		assert.Equal(t, time.Date(2024, time.January, 4, 14, 30, 0, 0, loc), occs[1])
		assert.Equal(t, time.Date(2024, time.January, 5, 14, 30, 0, 0, loc), occs[2])
	})

	t.Run("Multiple EXDATEs", func(t *testing.T) {
		dtstart := time.Date(2024, time.January, 1, 9, 0, 0, 0, time.UTC)
		rr := []string{
			"FREQ=DAILY;COUNT=10",
			"EXDATE:20240103T090000Z",
			"EXDATE:20240105T090000Z",
		}

		helper, err := ParseRecurrenceSet(rr, dtstart, false)
		assert.NoError(t, err)
		assert.NotNil(t, helper)

		exdates := helper.GetExDates()
		assert.Len(t, exdates, 2)

		occs, err := helper.NextN(dtstart, 5)
		assert.NoError(t, err)

		// Should skip 2024-01-03 and 2024-01-05
		assert.Equal(t, time.Date(2024, time.January, 2, 9, 0, 0, 0, time.UTC), occs[0])
		assert.Equal(t, time.Date(2024, time.January, 4, 9, 0, 0, 0, time.UTC), occs[1])
		assert.Equal(t, time.Date(2024, time.January, 6, 9, 0, 0, 0, time.UTC), occs[2])
		assert.Equal(t, time.Date(2024, time.January, 7, 9, 0, 0, 0, time.UTC), occs[3])
		assert.Equal(t, time.Date(2024, time.January, 8, 9, 0, 0, 0, time.UTC), occs[4])
	})
}
