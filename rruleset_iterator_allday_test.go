package rrule

import (
	"testing"
	"time"
)

func TestSetIterator_AllDay_BasicFunctionality(t *testing.T) {
	// Test basic all-day event iteration with RRule
	dtstart := time.Date(2025, 1, 1, 10, 30, 0, 0, time.UTC)
	rrule, err := NewRRule(ROption{
		Freq:    DAILY,
		Count:   1,
		Dtstart: dtstart,
	})
	if err != nil {
		t.Fatal(err)
	}

	set := &Set{
		dtstart: dtstart,
		rrule:   rrule,
		allDay:  true,
	}

	iter := set.Iterator()
	result, ok := iter()
	if !ok {
		t.Fatal("Expected iterator to return a value")
	}

	// Should return normalized floating time (00:00:00 UTC)
	expected := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	// Should not have more values
	_, ok = iter()
	if ok {
		t.Error("Expected iterator to be exhausted")
	}
}

func TestSetIterator_AllDay_RDateNormalization(t *testing.T) {
	// Test rdate normalization for all-day events
	dtstart := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	rdates := []time.Time{
		time.Date(2025, 1, 2, 14, 30, 0, 0, time.UTC),  // Should normalize to 00:00:00
		time.Date(2025, 1, 3, 9, 15, 30, 0, time.UTC),  // Should normalize to 00:00:00
		time.Date(2025, 1, 4, 23, 59, 59, 0, time.UTC), // Should normalize to 00:00:00
	}

	set := &Set{
		dtstart: dtstart,
		rdate:   rdates,
		allDay:  true,
	}

	iter := set.Iterator()
	expected := []time.Time{
		time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC), // Only rdates, not dtstart
		time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 4, 0, 0, 0, 0, time.UTC),
	}

	for i, exp := range expected {
		result, ok := iter()
		if !ok {
			t.Fatalf("Expected iterator to return value at index %d", i)
		}
		if !result.Equal(exp) {
			t.Errorf("At index %d: expected %v, got %v", i, exp, result)
		}
	}

	// Should be exhausted
	_, ok := iter()
	if ok {
		t.Error("Expected iterator to be exhausted")
	}
}

func TestSetIterator_AllDay_ExDateFiltering(t *testing.T) {
	// Test exdate filtering for all-day events
	dtstart := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	rdates := []time.Time{
		time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 4, 0, 0, 0, 0, time.UTC),
	}
	exdates := []time.Time{
		time.Date(2025, 1, 2, 15, 30, 0, 0, time.UTC), // Should exclude Jan 2 (normalized)
		time.Date(2025, 1, 4, 8, 45, 0, 0, time.UTC),  // Should exclude Jan 4 (normalized)
	}

	set := &Set{
		dtstart: dtstart,
		rdate:   rdates,
		exdate:  exdates,
		allDay:  true,
	}

	iter := set.Iterator()
	expected := []time.Time{
		time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC), // Only Jan 3 (not excluded), no dtstart
	}

	for i, exp := range expected {
		result, ok := iter()
		if !ok {
			t.Fatalf("Expected iterator to return value at index %d", i)
		}
		if !result.Equal(exp) {
			t.Errorf("At index %d: expected %v, got %v", i, exp, result)
		}
	}

	// Should be exhausted
	_, ok := iter()
	if ok {
		t.Error("Expected iterator to be exhausted")
	}
}

func TestSetIterator_AllDay_TimeZoneIndependence(t *testing.T) {
	// Test that all-day events are timezone-independent
	locations := []*time.Location{
		time.UTC,
		time.FixedZone("EST", -5*3600), // UTC-5
		time.FixedZone("JST", 9*3600),  // UTC+9
		time.FixedZone("PST", -8*3600), // UTC-8
	}

	for _, loc := range locations {
		t.Run("Location_"+loc.String(), func(t *testing.T) {
			// Create times in different timezones
			dtstart := time.Date(2025, 1, 1, 14, 30, 0, 0, loc)
			rdates := []time.Time{
				time.Date(2025, 1, 2, 9, 15, 0, 0, loc),
				time.Date(2025, 1, 3, 18, 45, 0, 0, loc),
			}

			set := &Set{
				dtstart: dtstart,
				rdate:   rdates,
				allDay:  true,
			}

			iter := set.Iterator()
			expected := []time.Time{
				time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC), // Normalized floating time
				time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC), // Normalized floating time
			}

			for i, exp := range expected {
				result, ok := iter()
				if !ok {
					t.Fatalf("Expected iterator to return value at index %d", i)
				}
				if !result.Equal(exp) {
					t.Errorf("At index %d: expected %v, got %v", i, exp, result)
				}
			}
		})
	}
}

func TestSetIterator_AllDay_WithRRule(t *testing.T) {
	// Test all-day events with rrule integration
	dtstart := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	rrule, err := NewRRule(ROption{
		Freq:    DAILY,
		Count:   3,
		Dtstart: dtstart,
	})
	if err != nil {
		t.Fatal(err)
	}

	set := &Set{
		dtstart: dtstart,
		rrule:   rrule,
		allDay:  true,
	}

	iter := set.Iterator()
	expected := []time.Time{
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC),
	}

	for i, exp := range expected {
		result, ok := iter()
		if !ok {
			t.Fatalf("Expected iterator to return value at index %d", i)
		}
		if !result.Equal(exp) {
			t.Errorf("At index %d: expected %v, got %v", i, exp, result)
		}
	}

	// Should be exhausted
	_, ok := iter()
	if ok {
		t.Error("Expected iterator to be exhausted")
	}
}

func TestSetIterator_AllDay_DuplicateElimination(t *testing.T) {
	// Test duplicate elimination for all-day events
	dtstart := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	rdates := []time.Time{
		time.Date(2025, 1, 1, 8, 30, 0, 0, time.UTC),  // Same date as dtstart, different time
		time.Date(2025, 1, 2, 14, 15, 0, 0, time.UTC), // Unique date
		time.Date(2025, 1, 2, 20, 45, 0, 0, time.UTC), // Same date as above, different time
	}

	set := &Set{
		dtstart: dtstart,
		rdate:   rdates,
		allDay:  true,
	}

	iter := set.Iterator()
	expected := []time.Time{
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), // Only one occurrence for Jan 1
		time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC), // Only one occurrence for Jan 2
	}

	for i, exp := range expected {
		result, ok := iter()
		if !ok {
			t.Fatalf("Expected iterator to return value at index %d", i)
		}
		if !result.Equal(exp) {
			t.Errorf("At index %d: expected %v, got %v", i, exp, result)
		}
	}

	// Should be exhausted
	_, ok := iter()
	if ok {
		t.Error("Expected iterator to be exhausted")
	}
}

func TestSetIterator_NonAllDay_Unchanged(t *testing.T) {
	// Test that non-all-day events are not affected by the changes
	dtstart := time.Date(2025, 1, 1, 10, 30, 0, 0, time.UTC)
	rdates := []time.Time{
		time.Date(2025, 1, 2, 14, 15, 0, 0, time.UTC),
		time.Date(2025, 1, 3, 9, 45, 0, 0, time.UTC),
	}

	set := &Set{
		dtstart: dtstart,
		rdate:   rdates,
		allDay:  false, // Non-all-day event
	}

	iter := set.Iterator()
	expected := []time.Time{
		time.Date(2025, 1, 2, 14, 15, 0, 0, time.UTC), // Only rdates, no dtstart
		time.Date(2025, 1, 3, 9, 45, 0, 0, time.UTC),
	}

	for i, exp := range expected {
		result, ok := iter()
		if !ok {
			t.Fatalf("Expected iterator to return value at index %d", i)
		}
		if !result.Equal(exp) {
			t.Errorf("At index %d: expected %v, got %v", i, exp, result)
		}
	}

	// Should be exhausted
	_, ok := iter()
	if ok {
		t.Error("Expected iterator to be exhausted")
	}
}

func TestSetIterator_AllDay_EdgeCases(t *testing.T) {
	// Test edge cases for all-day events
	t.Run("EmptyRDates", func(t *testing.T) {
		dtstart := time.Date(2025, 1, 1, 15, 30, 0, 0, time.UTC)
		rrule, err := NewRRule(ROption{
			Freq:    DAILY,
			Count:   1,
			Dtstart: dtstart,
		})
		if err != nil {
			t.Fatal(err)
		}

		set := &Set{
			dtstart: dtstart,
			rrule:   rrule,
			rdate:   []time.Time{}, // Empty rdates
			allDay:  true,
		}

		iter := set.Iterator()
		result, ok := iter()
		if !ok {
			t.Fatal("Expected iterator to return dtstart")
		}

		expected := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		if !result.Equal(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("EmptyExDates", func(t *testing.T) {
		dtstart := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		rrule, err := NewRRule(ROption{
			Freq:    DAILY,
			Count:   1,
			Dtstart: dtstart,
		})
		if err != nil {
			t.Fatal(err)
		}

		set := &Set{
			dtstart: dtstart,
			rrule:   rrule,
			exdate:  []time.Time{}, // Empty exdates
			allDay:  true,
		}

		iter := set.Iterator()
		result, ok := iter()
		if !ok {
			t.Fatal("Expected iterator to return dtstart")
		}

		expected := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		if !result.Equal(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("LeapYear", func(t *testing.T) {
		// Test leap year handling
		dtstart := time.Date(2024, 2, 29, 12, 0, 0, 0, time.UTC) // Leap year
		rrule, err := NewRRule(ROption{
			Freq:    DAILY,
			Count:   1,
			Dtstart: dtstart,
		})
		if err != nil {
			t.Fatal(err)
		}

		set := &Set{
			dtstart: dtstart,
			rrule:   rrule,
			allDay:  true,
		}

		iter := set.Iterator()
		result, ok := iter()
		if !ok {
			t.Fatal("Expected iterator to return value")
		}

		expected := time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC)
		if !result.Equal(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
}
