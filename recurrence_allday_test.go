package rrule

import (
	"strings"
	"testing"
	"time"
)

// TestAllDayDaily tests daily recurrence for all-day events.
func TestAllDayDaily(t *testing.T) {
	r, _ := newRecurrence(ROption{
		Freq:    DAILY,
		Count:   3,
		AllDay:  true,
		Dtstart: time.Date(2023, 1, 1, 14, 30, 0, 0, time.UTC),
	})
	want := []time.Time{
		time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
	}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

// TestAllDayWeekly tests weekly recurrence for all-day events.
func TestAllDayWeekly(t *testing.T) {
	r, _ := newRecurrence(ROption{
		Freq:    WEEKLY,
		Count:   3,
		AllDay:  true,
		Dtstart: time.Date(2023, 1, 1, 9, 15, 30, 0, time.UTC), // Sunday
	})
	want := []time.Time{
		time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 8, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
	}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

// TestAllDayMonthly tests monthly recurrence for all-day events.
func TestAllDayMonthly(t *testing.T) {
	r, _ := newRecurrence(ROption{
		Freq:    MONTHLY,
		Count:   3,
		AllDay:  true,
		Dtstart: time.Date(2023, 1, 15, 16, 45, 0, 0, time.UTC),
	})
	want := []time.Time{
		time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 2, 15, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 3, 15, 0, 0, 0, 0, time.UTC),
	}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

// TestAllDayYearly tests yearly recurrence for all-day events.
func TestAllDayYearly(t *testing.T) {
	r, _ := newRecurrence(ROption{
		Freq:    YEARLY,
		Count:   3,
		AllDay:  true,
		Dtstart: time.Date(2023, 6, 15, 23, 59, 59, 0, time.UTC),
	})
	want := []time.Time{
		time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC),
	}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

// TestAllDayWithUntil tests UNTIL boundary handling for all-day events.
func TestAllDayWithUntil(t *testing.T) {
	r, _ := newRecurrence(ROption{
		Freq:    DAILY,
		AllDay:  true,
		Dtstart: time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
		Until:   time.Date(2023, 1, 3, 23, 59, 59, 0, time.UTC),
	})
	want := []time.Time{
		time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
	}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

// TestAllDayWithUntilMidnight tests the midnight UNTIL boundary for all-day events.
func TestAllDayWithUntilMidnight(t *testing.T) {
	r, _ := newRecurrence(ROption{
		Freq:    DAILY,
		AllDay:  true,
		Dtstart: time.Date(2023, 1, 1, 15, 0, 0, 0, time.UTC),
		Until:   time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
	})
	want := []time.Time{
		time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
	}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

// TestAllDayWithCount tests COUNT handling for all-day events.
func TestAllDayWithCount(t *testing.T) {
	r, _ := newRecurrence(ROption{
		Freq:    WEEKLY,
		Count:   5,
		AllDay:  true,
		Dtstart: time.Date(2023, 1, 1, 8, 45, 30, 0, time.UTC),
	})
	want := []time.Time{
		time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 8, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 22, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 29, 0, 0, 0, 0, time.UTC),
	}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

// TestAllDayDTStartProcessing tests DTStart normalization for all-day events.
func TestAllDayDTStartProcessing(t *testing.T) {
	// DTStart values at different times should normalize to 00:00:00.
	testCases := []struct {
		name    string
		dtstart time.Time
	}{
		{"Morning", time.Date(2023, 1, 1, 8, 30, 15, 0, time.UTC)},
		{"Noon", time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)},
		{"Evening", time.Date(2023, 1, 1, 18, 45, 30, 0, time.UTC)},
		{"Late Night", time.Date(2023, 1, 1, 23, 59, 59, 0, time.UTC)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := newRecurrence(ROption{
				Freq:    DAILY,
				Count:   2,
				AllDay:  true,
				Dtstart: tc.dtstart,
			})
			want := []time.Time{
				time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			}
			value := r.All()
			if !timesEqual(value, want) {
				t.Errorf("get %v, want %v", value, want)
			}
		})
	}
}

// TestAllDaySetWithRDate tests Set RDate for all-day events.
func TestAllDaySetWithRDate(t *testing.T) {
	set := New(ROption{
		Freq:    WEEKLY,
		Count:   2,
		AllDay:  true,
		Dtstart: time.Date(2023, 1, 1, 14, 30, 0, 0, time.UTC),
	})
	if set == nil {
		t.Fatal("failed to create recurrence")
	}

	// Add extra dates.
	set.RDate(time.Date(2023, 1, 20, 16, 45, 0, 0, time.UTC))
	set.RDate(time.Date(2023, 1, 25, 9, 15, 30, 0, time.UTC))

	want := []time.Time{
		time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 8, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 20, 0, 0, 0, 0, time.UTC), // RDate normalized to floating midnight
		time.Date(2023, 1, 25, 0, 0, 0, 0, time.UTC), // RDate normalized to floating midnight
	}
	value := set.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

// TestAllDaySetWithExDate tests Set ExDate for all-day events.
func TestAllDaySetWithExDate(t *testing.T) {
	set := New(ROption{
		Freq:    DAILY,
		Count:   5,
		AllDay:  true,
		Dtstart: time.Date(2023, 1, 1, 11, 20, 0, 0, time.UTC),
	})
	if set == nil {
		t.Fatal("failed to create recurrence")
	}

	// Exclude specific dates (use midnight to match all-day events).
	set.ExDate(time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC))
	set.ExDate(time.Date(2023, 1, 4, 0, 0, 0, 0, time.UTC))

	want := []time.Time{
		time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 5, 0, 0, 0, 0, time.UTC),
	}
	value := set.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

// TestAllDaySetComplex tests complex Set combinations for all-day events.
func TestAllDaySetComplex(t *testing.T) {
	set := New(ROption{
		Freq:    WEEKLY,
		Count:   3,
		AllDay:  true,
		Dtstart: time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
	})
	if set == nil {
		t.Fatal("failed to create recurrence")
	}

	// Add an extra date.
	set.RDate(time.Date(2023, 1, 10, 12, 45, 0, 0, time.UTC))

	// Exclude a date (use midnight to match all-day events).
	set.ExDate(time.Date(2023, 1, 8, 0, 0, 0, 0, time.UTC))

	want := []time.Time{
		time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC), // RDate normalized to floating midnight
		time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
	}
	value := set.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

// TestAllDayTimezoneHandling tests timezone handling for all-day events.
// Per RFC 5545, all-day events use floating time and represent 00:00:00 in any timezone.
func TestAllDayTimezoneHandling(t *testing.T) {
	testCases := []struct {
		name string
		tz   *time.Location
	}{
		{"UTC", time.UTC},
		{"EST", time.FixedZone("EST", -5*3600)}, // UTC-5
		{"JST", time.FixedZone("JST", 9*3600)},  // UTC+9
		{"CET", time.FixedZone("CET", 1*3600)},  // UTC+1
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := newRecurrence(ROption{
				Freq:    DAILY,
				Count:   2,
				AllDay:  true,
				Dtstart: time.Date(2023, 1, 1, 15, 30, 0, 0, tc.tz),
			})

			// Per RFC 5545, all-day events should convert to floating time (no timezone binding).
			// In Go we represent floating time with UTC since it is timezone-agnostic.
			// This ensures users in any timezone see the same day's 00:00:00.
			want := []time.Time{
				time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			}
			value := r.All()
			if !timesEqual(value, want) {
				t.Errorf("get %v, want %v", value, want)
			}
		})
	}
}

// TestAllDayLeapYear tests leap year handling for all-day events.
func TestAllDayLeapYear(t *testing.T) {
	r, _ := newRecurrence(ROption{
		Freq:    YEARLY,
		Count:   4,
		AllDay:  true,
		Dtstart: time.Date(2020, 2, 29, 12, 0, 0, 0, time.UTC), // Leap day.
	})
	want := []time.Time{
		time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC), // Next leap year.
		time.Date(2028, 2, 29, 0, 0, 0, 0, time.UTC),
		time.Date(2032, 2, 29, 0, 0, 0, 0, time.UTC),
	}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

// TestAllDayYearBoundary tests year boundary handling for all-day events.
func TestAllDayYearBoundary(t *testing.T) {
	r, _ := newRecurrence(ROption{
		Freq:    DAILY,
		Count:   5,
		AllDay:  true,
		Dtstart: time.Date(2022, 12, 30, 18, 45, 0, 0, time.UTC),
	})
	want := []time.Time{
		time.Date(2022, 12, 30, 0, 0, 0, 0, time.UTC),
		time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
	}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

// TestAllDayMonthBoundary tests month boundary handling for all-day events.
func TestAllDayMonthBoundary(t *testing.T) {
	r, _ := newRecurrence(ROption{
		Freq:    DAILY,
		Count:   4,
		AllDay:  true,
		Dtstart: time.Date(2023, 1, 30, 20, 15, 0, 0, time.UTC),
	})
	want := []time.Time{
		time.Date(2023, 1, 30, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 2, 2, 0, 0, 0, 0, time.UTC),
	}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

// TestAllDayIterator tests the iterator for all-day events.
func TestAllDayIterator(t *testing.T) {
	r, _ := newRecurrence(ROption{
		Freq:    DAILY,
		Count:   3,
		AllDay:  true,
		Dtstart: time.Date(2023, 1, 1, 13, 25, 45, 0, time.UTC),
	})

	iter := r.Iterator()
	var results []time.Time

	for {
		dt, ok := iter()
		if !ok {
			break
		}
		results = append(results, dt)
	}

	want := []time.Time{
		time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
	}

	if !timesEqual(results, want) {
		t.Errorf("get %v, want %v", results, want)
	}
}

// TestAllDayWithByWeekDay tests all-day events with ByWeekday rules.
func TestAllDayWithByWeekDay(t *testing.T) {
	r, _ := newRecurrence(ROption{
		Freq:      WEEKLY,
		Count:     4,
		AllDay:    true,
		Byweekday: []Weekday{MO, WE, FR},
		Dtstart:   time.Date(2023, 1, 2, 11, 30, 0, 0, time.UTC), // Monday
	})
	want := []time.Time{
		time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC), // Monday
		time.Date(2023, 1, 4, 0, 0, 0, 0, time.UTC), // Wednesday
		time.Date(2023, 1, 6, 0, 0, 0, 0, time.UTC), // Friday
		time.Date(2023, 1, 9, 0, 0, 0, 0, time.UTC), // Next Monday
	}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

// TestAllDayWithByMonthDay tests all-day events with ByMonthday rules.
func TestAllDayWithByMonthDay(t *testing.T) {
	r, _ := newRecurrence(ROption{
		Freq:       MONTHLY,
		Count:      3,
		AllDay:     true,
		Bymonthday: []int{1, 15},
		Dtstart:    time.Date(2023, 1, 1, 17, 45, 0, 0, time.UTC),
	})
	want := []time.Time{
		time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
	}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

// TestAllDayConsistencyWithNonAllDay tests consistency between all-day and timed events.
func TestAllDayConsistencyWithNonAllDay(t *testing.T) {
	// All-day event.
	allDayRule, _ := newRecurrence(ROption{
		Freq:    DAILY,
		Count:   3,
		AllDay:  true,
		Dtstart: time.Date(2023, 1, 1, 14, 30, 0, 0, time.UTC),
	})

	// Timed event at midnight on the same date.
	nonAllDayRule, _ := newRecurrence(ROption{
		Freq:    DAILY,
		Count:   3,
		AllDay:  false,
		Dtstart: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	})

	allDayResults := allDayRule.All()
	nonAllDayResults := nonAllDayRule.All()

	// All-day results should match timed midnight results.
	if !timesEqual(allDayResults, nonAllDayResults) {
		t.Errorf("AllDay results %v should match non-AllDay midnight results %v",
			allDayResults, nonAllDayResults)
	}
}

func TestAllDay_BeforeAfterBetween(t *testing.T) {
	dt := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	r, err := newRecurrence(ROption{Freq: DAILY, Count: 3, AllDay: true, Dtstart: dt})
	if err != nil {
		t.Fatal(err)
	}

	// Occurrences: 2023-01-01, 2023-01-02, 2023-01-03 (midnight UTC)
	d1 := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	d2 := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)
	d3 := time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC)

	// After at start
	if got := r.After(d1, false); !got.Equal(d2) {
		t.Errorf("After exclude start: got %v want %v", got, d2)
	}
	if got := r.After(d1, true); !got.Equal(d1) {
		t.Errorf("After include start: got %v want %v", got, d1)
	}

	// Before at end
	if got := r.Before(d3, false); !got.Equal(d2) {
		t.Errorf("Before exclude end: got %v want %v", got, d2)
	}
	if got := r.Before(d3, true); !got.Equal(d3) {
		t.Errorf("Before include end: got %v want %v", got, d3)
	}

	// Between window
	// (d1, d3) exclusive should give only d2
	got := r.Between(d1, d3, false)
	if len(got) != 1 || !got[0].Equal(d2) {
		t.Errorf("Between exclusive got %v want [%v]", got, d2)
	}
	// [d1, d3] inclusive should give d1,d2,d3
	got = r.Between(d1, d3, true)
	want := []time.Time{d1, d2, d3}
	if !timesEqual(got, want) {
		t.Errorf("Between inclusive got %v want %v", got, want)
	}
}

func TestRecurrenceIterator_AllDay_BasicFunctionality(t *testing.T) {
	// Test basic all-day event iteration with RRule
	dtstart := time.Date(2025, 1, 1, 10, 30, 0, 0, time.UTC)
	set := New(ROption{
		Freq:    DAILY,
		Count:   1,
		AllDay:  true,
		Dtstart: dtstart,
	})
	if set == nil {
		t.Fatal("failed to create recurrence")
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

func TestRecurrenceIterator_AllDay_RDateNormalization(t *testing.T) {
	// Test rdate normalization for all-day events
	rdates := []time.Time{
		time.Date(2025, 1, 2, 14, 30, 0, 0, time.UTC),  // Should normalize to 00:00:00
		time.Date(2025, 1, 3, 9, 15, 30, 0, time.UTC),  // Should normalize to 00:00:00
		time.Date(2025, 1, 4, 23, 59, 59, 0, time.UTC), // Should normalize to 00:00:00
	}

	set := &Recurrence{}
	set.SetAllDay(true)
	set.SetRDates(rdates)

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

func TestRecurrenceIterator_AllDay_ExDateFiltering(t *testing.T) {
	// Test exdate filtering for all-day events
	rdates := []time.Time{
		time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 4, 0, 0, 0, 0, time.UTC),
	}
	exdates := []time.Time{
		time.Date(2025, 1, 2, 15, 30, 0, 0, time.UTC), // Should exclude Jan 2 (normalized)
		time.Date(2025, 1, 4, 8, 45, 0, 0, time.UTC),  // Should exclude Jan 4 (normalized)
	}

	set := &Recurrence{}
	set.SetAllDay(true)
	set.SetRDates(rdates)
	set.SetExDates(exdates)

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

func TestRecurrenceIterator_AllDay_TimeZoneIndependence(t *testing.T) {
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
			rdates := []time.Time{
				time.Date(2025, 1, 2, 9, 15, 0, 0, loc),
				time.Date(2025, 1, 3, 18, 45, 0, 0, loc),
			}

			set := &Recurrence{}
			set.SetAllDay(true)
			set.SetRDates(rdates)

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

func TestRecurrenceIterator_AllDay_WithRRule(t *testing.T) {
	// Test all-day events with rrule integration
	dtstart := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	set := New(ROption{
		Freq:    DAILY,
		Count:   3,
		AllDay:  true,
		Dtstart: dtstart,
	})
	if set == nil {
		t.Fatal("failed to create recurrence")
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

func TestRecurrenceIterator_AllDay_DuplicateElimination(t *testing.T) {
	// Test duplicate elimination for all-day events
	rdates := []time.Time{
		time.Date(2025, 1, 1, 8, 30, 0, 0, time.UTC),  // Same date as dtstart, different time
		time.Date(2025, 1, 2, 14, 15, 0, 0, time.UTC), // Unique date
		time.Date(2025, 1, 2, 20, 45, 0, 0, time.UTC), // Same date as above, different time
	}

	set := &Recurrence{}
	set.SetAllDay(true)
	set.SetRDates(rdates)

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

func TestRecurrenceIterator_NonAllDay_Unchanged(t *testing.T) {
	// Test that non-all-day events are not affected by the changes
	rdates := []time.Time{
		time.Date(2025, 1, 2, 14, 15, 0, 0, time.UTC),
		time.Date(2025, 1, 3, 9, 45, 0, 0, time.UTC),
	}

	set := &Recurrence{}
	set.SetAllDay(false)
	set.SetRDates(rdates)

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

func TestRecurrenceIterator_AllDay_EdgeCases(t *testing.T) {
	// Test edge cases for all-day events
	t.Run("EmptyRDates", func(t *testing.T) {
		dtstart := time.Date(2025, 1, 1, 15, 30, 0, 0, time.UTC)
		set, err := newRecurrence(ROption{
			Freq:    DAILY,
			Count:   1,
			AllDay:  true,
			Dtstart: dtstart,
		})
		if err != nil {
			t.Fatal(err)
		}

		set.SetRDates([]time.Time{})

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
		set, err := newRecurrence(ROption{
			Freq:    DAILY,
			Count:   1,
			AllDay:  true,
			Dtstart: dtstart,
		})
		if err != nil {
			t.Fatal(err)
		}

		set.SetExDates([]time.Time{})

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
		set, err := newRecurrence(ROption{
			Freq:    DAILY,
			Count:   1,
			AllDay:  true,
			Dtstart: dtstart,
		})
		if err != nil {
			t.Fatal(err)
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

func TestRecurrenceAllDay_BasicFunctionality(t *testing.T) {
	set := &Recurrence{}

	// Test initial state
	if set.IsAllDay() {
		t.Error("New Set should not be all-day by default")
	}

	// Test setting all-day flag
	set.SetAllDay(true)
	if !set.IsAllDay() {
		t.Error("Set should be all-day after SetAllDay(true)")
	}

	// Test unsetting all-day flag
	set.SetAllDay(false)
	if set.IsAllDay() {
		t.Error("Set should not be all-day after SetAllDay(false)")
	}
}

func TestRecurrenceAllDay_DTStartNormalization(t *testing.T) {
	set := &Recurrence{}
	set.SetAllDay(true)

	// Test DTStart with timezone - should be normalized to UTC midnight
	loc, _ := time.LoadLocation("America/New_York")
	dtstart := time.Date(2024, 1, 15, 14, 30, 45, 0, loc)

	set.DTStart(dtstart)

	expected := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	if !set.GetDTStart().Equal(expected) {
		t.Errorf("Expected DTStart %v, got %v", expected, set.GetDTStart())
	}
}

func TestRecurrenceAllDay_RDateNormalization(t *testing.T) {
	set := &Recurrence{}
	set.SetAllDay(true)

	// Test RDate with timezone - should be normalized to UTC midnight
	loc, _ := time.LoadLocation("Europe/London")
	rdate1 := time.Date(2024, 2, 10, 9, 15, 30, 0, loc)
	rdate2 := time.Date(2024, 2, 12, 18, 45, 0, 0, loc)

	set.RDate(rdate1)
	set.RDate(rdate2)

	rdates := set.GetRDate()
	if len(rdates) != 2 {
		t.Errorf("Expected 2 RDates, got %d", len(rdates))
	}

	expected1 := time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC)
	expected2 := time.Date(2024, 2, 12, 0, 0, 0, 0, time.UTC)

	if !rdates[0].Equal(expected1) {
		t.Errorf("Expected first RDate %v, got %v", expected1, rdates[0])
	}
	if !rdates[1].Equal(expected2) {
		t.Errorf("Expected second RDate %v, got %v", expected2, rdates[1])
	}
}

func TestRecurrenceAllDay_SetRDatesNormalization(t *testing.T) {
	set := &Recurrence{}
	set.SetAllDay(true)

	// Test SetRDates with multiple timezones
	loc1, _ := time.LoadLocation("Asia/Tokyo")
	loc2, _ := time.LoadLocation("America/Los_Angeles")

	rdates := []time.Time{
		time.Date(2024, 3, 5, 10, 30, 0, 0, loc1),
		time.Date(2024, 3, 7, 22, 15, 45, 0, loc2),
		time.Date(2024, 3, 9, 6, 0, 0, 0, time.UTC),
	}

	set.SetRDates(rdates)

	result := set.GetRDate()
	if len(result) != 3 {
		t.Errorf("Expected 3 RDates, got %d", len(result))
	}

	expected := []time.Time{
		time.Date(2024, 3, 5, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 3, 7, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 3, 9, 0, 0, 0, 0, time.UTC),
	}

	for i, exp := range expected {
		if !result[i].Equal(exp) {
			t.Errorf("Expected RDate[%d] %v, got %v", i, exp, result[i])
		}
	}
}

func TestRecurrenceAllDay_ExDateNormalization(t *testing.T) {
	set := &Recurrence{}
	set.SetAllDay(true)

	// Test ExDate with timezone - should be normalized to UTC midnight
	loc, _ := time.LoadLocation("Australia/Sydney")
	exdate1 := time.Date(2024, 4, 20, 11, 45, 30, 0, loc)
	exdate2 := time.Date(2024, 4, 22, 16, 30, 15, 0, loc)

	set.ExDate(exdate1)
	set.ExDate(exdate2)

	exdates := set.GetExDate()
	if len(exdates) != 2 {
		t.Errorf("Expected 2 ExDates, got %d", len(exdates))
	}

	expected1 := time.Date(2024, 4, 20, 0, 0, 0, 0, time.UTC)
	expected2 := time.Date(2024, 4, 22, 0, 0, 0, 0, time.UTC)

	if !exdates[0].Equal(expected1) {
		t.Errorf("Expected first ExDate %v, got %v", expected1, exdates[0])
	}
	if !exdates[1].Equal(expected2) {
		t.Errorf("Expected second ExDate %v, got %v", expected2, exdates[1])
	}
}

func TestRecurrenceAllDay_SetExDatesNormalization(t *testing.T) {
	set := &Recurrence{}
	set.SetAllDay(true)

	// Test SetExDates with multiple timezones
	loc1, _ := time.LoadLocation("Europe/Paris")
	loc2, _ := time.LoadLocation("America/Chicago")

	exdates := []time.Time{
		time.Date(2024, 5, 10, 8, 30, 0, 0, loc1),
		time.Date(2024, 5, 12, 20, 15, 45, 0, loc2),
		time.Date(2024, 5, 14, 12, 0, 0, 0, time.UTC),
	}

	set.SetExDates(exdates)

	result := set.GetExDate()
	if len(result) != 3 {
		t.Errorf("Expected 3 ExDates, got %d", len(result))
	}

	expected := []time.Time{
		time.Date(2024, 5, 10, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 5, 12, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 5, 14, 0, 0, 0, 0, time.UTC),
	}

	for i, exp := range expected {
		if !result[i].Equal(exp) {
			t.Errorf("Expected ExDate[%d] %v, got %v", i, exp, result[i])
		}
	}
}

func TestRecurrenceAllDay_ExistingTimesNormalization(t *testing.T) {
	set := &Recurrence{}

	// Set up non-all-day times first
	loc, _ := time.LoadLocation("America/New_York")
	dtstart := time.Date(2024, 6, 15, 14, 30, 45, 0, loc)
	rdate := time.Date(2024, 6, 17, 10, 15, 30, 0, loc)
	exdate := time.Date(2024, 6, 19, 16, 45, 0, 0, loc)

	set.DTStart(dtstart)
	set.RDate(rdate)
	set.ExDate(exdate)

	// Verify non-all-day times are preserved with timezone
	if set.GetDTStart().Location() == time.UTC {
		t.Error("DTStart should preserve timezone before SetAllDay")
	}

	// Switch to all-day - should normalize existing times
	set.SetAllDay(true)

	// Verify all times are normalized to UTC midnight
	expectedDTStart := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	if !set.GetDTStart().Equal(expectedDTStart) {
		t.Errorf("Expected normalized DTStart %v, got %v", expectedDTStart, set.GetDTStart())
	}

	rdates := set.GetRDate()
	expectedRDate := time.Date(2024, 6, 17, 0, 0, 0, 0, time.UTC)
	if len(rdates) != 1 || !rdates[0].Equal(expectedRDate) {
		t.Errorf("Expected normalized RDate %v, got %v", expectedRDate, rdates)
	}

	exdates := set.GetExDate()
	expectedExDate := time.Date(2024, 6, 19, 0, 0, 0, 0, time.UTC)
	if len(exdates) != 1 || !exdates[0].Equal(expectedExDate) {
		t.Errorf("Expected normalized ExDate %v, got %v", expectedExDate, exdates)
	}
}

func TestRecurrenceAllDay_RecurrenceSerialization(t *testing.T) {
	set := &Recurrence{}
	set.SetAllDay(true)

	// Set up all-day event data
	dtstart := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)
	rdate := time.Date(2024, 7, 5, 0, 0, 0, 0, time.UTC)
	exdate := time.Date(2024, 7, 10, 0, 0, 0, 0, time.UTC)

	set.DTStart(dtstart)
	set.RDate(rdate)
	set.ExDate(exdate)

	// Test Recurrence() serialization
	recurrence := set.Strings()

	// Verify DTSTART format (should use VALUE=DATE as per RFC 5545)
	dtStartFound := false
	for _, line := range recurrence {
		if strings.HasPrefix(line, "DTSTART;VALUE=DATE:") {
			dtStartFound = true
			expected := "DTSTART;VALUE=DATE:20240701"
			if line != expected {
				t.Errorf("Expected DTSTART %s, got %s", expected, line)
			}
			break
		}
	}
	if !dtStartFound {
		t.Error("DTSTART not found in recurrence")
	}

	// Verify RDATE format (should use VALUE=DATE as per RFC 5545)
	rDateFound := false
	for _, line := range recurrence {
		if strings.HasPrefix(line, "RDATE;VALUE=DATE:") {
			rDateFound = true
			expected := "RDATE;VALUE=DATE:20240705"
			if line != expected {
				t.Errorf("Expected RDATE %s, got %s", expected, line)
			}
			break
		}
	}
	if !rDateFound {
		t.Error("RDATE not found in recurrence")
	}

	// Verify EXDATE format (should use VALUE=DATE as per RFC 5545)
	exDateFound := false
	for _, line := range recurrence {
		if strings.HasPrefix(line, "EXDATE;VALUE=DATE:") {
			exDateFound = true
			expected := "EXDATE;VALUE=DATE:20240710"
			if line != expected {
				t.Errorf("Expected EXDATE %s, got %s", expected, line)
			}
			break
		}
	}
	if !exDateFound {
		t.Error("EXDATE not found in recurrence")
	}
}

func TestRecurrenceAllDay_StringSerialization(t *testing.T) {
	set := &Recurrence{}
	set.SetAllDay(true)

	// Set up all-day event data
	dtstart := time.Date(2024, 8, 15, 0, 0, 0, 0, time.UTC)
	set.DTStart(dtstart)

	// Test String() serialization
	str := set.String()

	// Should contain VALUE=DATE format as per RFC 5545
	if !strings.Contains(str, "DTSTART;VALUE=DATE:20240815") {
		t.Errorf("String() should contain VALUE=DATE format, got: %s", str)
	}

	// Should not contain time part for all-day events
	if strings.Contains(str, "20240815T000000") {
		t.Errorf("String() should not contain time part for all-day events, got: %s", str)
	}
}

func TestRecurrenceAllDay_NonAllDayPreservesTimezone(t *testing.T) {
	set := &Recurrence{}
	// Keep as non-all-day (default)

	// Set up times with timezone
	loc, _ := time.LoadLocation("Europe/Berlin")
	dtstart := time.Date(2024, 9, 20, 14, 30, 45, 0, loc)
	rdate := time.Date(2024, 9, 22, 10, 15, 30, 0, loc)
	exdate := time.Date(2024, 9, 24, 16, 45, 0, 0, loc)

	set.DTStart(dtstart)
	set.RDate(rdate)
	set.ExDate(exdate)

	// Verify times are truncated to seconds but preserve timezone info
	expectedDTStart := dtstart.Truncate(time.Second)
	if !set.GetDTStart().Equal(expectedDTStart) {
		t.Errorf("Expected DTStart %v, got %v", expectedDTStart, set.GetDTStart())
	}

	rdates := set.GetRDate()
	expectedRDate := rdate.Truncate(time.Second)
	if len(rdates) != 1 || !rdates[0].Equal(expectedRDate) {
		t.Errorf("Expected RDate %v, got %v", expectedRDate, rdates)
	}

	exdates := set.GetExDate()
	expectedExDate := exdate.Truncate(time.Second)
	if len(exdates) != 1 || !exdates[0].Equal(expectedExDate) {
		t.Errorf("Expected ExDate %v, got %v", expectedExDate, exdates)
	}
}

func TestRecurrenceAllDay_EdgeCases(t *testing.T) {
	set := &Recurrence{}
	set.SetAllDay(true)

	// Test with zero time
	zeroTime := time.Time{}
	set.DTStart(zeroTime)

	// Zero time should remain zero (not normalized)
	if !set.GetDTStart().IsZero() {
		t.Error("Zero time should remain zero")
	}

	// Test leap year date
	leapDate := time.Date(2024, 2, 29, 15, 30, 0, 0, time.UTC)
	set.DTStart(leapDate)

	expected := time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC)
	if !set.GetDTStart().Equal(expected) {
		t.Errorf("Expected leap year date %v, got %v", expected, set.GetDTStart())
	}

	// Test year boundary
	yearBoundary := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
	set.RDate(yearBoundary)

	rdates := set.GetRDate()
	expectedYearBoundary := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	if len(rdates) != 1 || !rdates[0].Equal(expectedYearBoundary) {
		t.Errorf("Expected year boundary date %v, got %v", expectedYearBoundary, rdates)
	}
}

// TestAllDaySetStringWithRDate tests all-day Set string serialization with RDATE.
func TestAllDaySetStringWithRDate(t *testing.T) {
	testCases := []struct {
		name     string
		dtstart  time.Time
		rdates   []time.Time
		expected []string
	}{
		{
			name:    "Single RDATE",
			dtstart: time.Date(2023, 5, 1, 9, 30, 0, 0, time.UTC),
			rdates: []time.Time{
				time.Date(2023, 5, 5, 14, 15, 0, 0, time.UTC),
			},
			expected: []string{
				"DTSTART;VALUE=DATE:20230501",
				"RDATE;VALUE=DATE:20230505",
			},
		},
		{
			name:    "Multiple RDATEs with different timezones",
			dtstart: time.Date(2023, 6, 10, 8, 0, 0, 0, time.FixedZone("EST", -5*3600)),
			rdates: []time.Time{
				time.Date(2023, 6, 15, 16, 30, 0, 0, time.FixedZone("JST", 9*3600)),
				time.Date(2023, 6, 20, 22, 45, 0, 0, time.FixedZone("CET", 1*3600)),
				time.Date(2023, 6, 25, 11, 0, 0, 0, time.UTC),
			},
			expected: []string{
				"DTSTART;VALUE=DATE:20230610",
				"RDATE;VALUE=DATE:20230615",
				"RDATE;VALUE=DATE:20230620",
				"RDATE;VALUE=DATE:20230625",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			set := &Recurrence{}
			set.SetAllDay(true)
			set.DTStart(tc.dtstart)

			for _, rdate := range tc.rdates {
				set.RDate(rdate)
			}

			output := set.String()
			t.Logf("RDATE test %s output: %s", tc.name, output)

			// Verify all expected strings are present.
			for _, expected := range tc.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected %s in output, got: %s", expected, output)
				}
			}

			// Verify RDATE uses VALUE=DATE format and has no time part.
			lines := strings.Split(output, "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "RDATE") {
					if !strings.Contains(line, "VALUE=DATE") {
						t.Errorf("RDATE should use VALUE=DATE format for all-day events, got: %s", line)
					}
					// Check for a time component (T followed by digits).
					if strings.Contains(line, "T") && !strings.Contains(line, "VALUE=DATE") {
						t.Errorf("RDATE should not contain time part for all-day events, got: %s", line)
					}
				}
			}
		})
	}
}

// TestAllDaySetStringWithExDate tests all-day Set string serialization with EXDATE.
func TestAllDaySetStringWithExDate(t *testing.T) {
	testCases := []struct {
		name     string
		dtstart  time.Time
		exdates  []time.Time
		expected []string
	}{
		{
			name:    "Single EXDATE",
			dtstart: time.Date(2023, 7, 1, 10, 0, 0, 0, time.UTC),
			exdates: []time.Time{
				time.Date(2023, 7, 4, 15, 30, 0, 0, time.UTC),
			},
			expected: []string{
				"DTSTART;VALUE=DATE:20230701",
				"EXDATE;VALUE=DATE:20230704",
			},
		},
		{
			name:    "Multiple EXDATEs with different timezones",
			dtstart: time.Date(2023, 8, 1, 12, 0, 0, 0, time.FixedZone("PST", -8*3600)),
			exdates: []time.Time{
				time.Date(2023, 8, 5, 6, 0, 0, 0, time.FixedZone("EST", -5*3600)),
				time.Date(2023, 8, 10, 18, 30, 0, 0, time.FixedZone("JST", 9*3600)),
				time.Date(2023, 8, 15, 23, 59, 59, 0, time.UTC),
			},
			expected: []string{
				"DTSTART;VALUE=DATE:20230801",
				"EXDATE;VALUE=DATE:20230805",
				"EXDATE;VALUE=DATE:20230810",
				"EXDATE;VALUE=DATE:20230815",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			set := &Recurrence{}
			set.SetAllDay(true)
			set.DTStart(tc.dtstart)

			for _, exdate := range tc.exdates {
				set.ExDate(exdate)
			}

			output := set.String()
			t.Logf("EXDATE test %s output: %s", tc.name, output)

			// Verify all expected strings are present.
			for _, expected := range tc.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected %s in output, got: %s", expected, output)
				}
			}

			// Verify EXDATE uses VALUE=DATE format and has no time part.
			lines := strings.Split(output, "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "EXDATE") {
					if !strings.Contains(line, "VALUE=DATE") {
						t.Errorf("EXDATE should use VALUE=DATE format for all-day events, got: %s", line)
					}
					// Check for a time component (T followed by digits).
					if strings.Contains(line, "T") && !strings.Contains(line, "VALUE=DATE") {
						t.Errorf("EXDATE should not contain time part for all-day events, got: %s", line)
					}
				}
			}
		})
	}
}

// TestAllDaySetStringComplex tests complex all-day Set scenarios (RRULE + RDATE + EXDATE + UNTIL).
func TestAllDaySetStringComplex(t *testing.T) {
	dtstart := time.Date(2023, 9, 1, 14, 30, 0, 0, time.FixedZone("EST", -5*3600))
	until := time.Date(2023, 9, 30, 23, 59, 59, 0, time.UTC)
	set, err := newRecurrence(ROption{
		Freq:    WEEKLY,
		AllDay:  true,
		Dtstart: dtstart,
		Until:   until,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Add RDATE.
	set.RDate(time.Date(2023, 9, 15, 16, 0, 0, 0, time.FixedZone("JST", 9*3600)))
	set.RDate(time.Date(2023, 9, 25, 8, 30, 0, 0, time.UTC))

	// Add EXDATE.
	set.ExDate(time.Date(2023, 9, 8, 12, 0, 0, 0, time.UTC))
	set.ExDate(time.Date(2023, 9, 22, 20, 15, 0, 0, time.FixedZone("CET", 1*3600)))

	output := set.String()
	t.Logf("Complex all-day set output: %s", output)

	expectedStrings := []string{
		"DTSTART;VALUE=DATE:20230901",
		"RRULE:FREQ=WEEKLY;UNTIL=20230930",
		"RDATE;VALUE=DATE:20230915",
		"RDATE;VALUE=DATE:20230925",
		"EXDATE;VALUE=DATE:20230908",
		"EXDATE;VALUE=DATE:20230922",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected %s in output, got: %s", expected, output)
		}
	}

	// Verify all date-related fields use DATE format without time parts.
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if (strings.HasPrefix(line, "DTSTART") ||
			strings.HasPrefix(line, "RDATE") ||
			strings.HasPrefix(line, "EXDATE")) &&
			strings.Contains(line, "T") &&
			!strings.Contains(line, "VALUE=DATE") {
			t.Errorf("Date field should use VALUE=DATE format for all-day events, got: %s", line)
		}
		// UNTIL in RRULE should use DATE format (no "T").
		// Check whether UNTIL contains a time component (T followed by digits).
		if strings.Contains(line, "UNTIL=") {
			// Extract the UNTIL value.
			parts := strings.Split(line, "UNTIL=")
			if len(parts) > 1 {
				untilValue := strings.Split(parts[1], ";")[0]  // Handle possible trailing params.
				untilValue = strings.Split(untilValue, " ")[0] // Handle possible whitespace.
				// Check for a time component (T followed by digits).
				if strings.Contains(untilValue, "T") {
					t.Errorf("UNTIL should use DATE format (no time part) for all-day events, got: %s", untilValue)
				}
			}
		}
	}
}

func TestStrSliceToRRuleSetDetectsAllDayFromRDate(t *testing.T) {
	lines := []string{
		"RDATE;VALUE=DATE:20240301,20240303",
	}

	set, err := StrSliceToRRuleSetInLoc(lines, time.UTC)
	if err != nil {
		t.Fatalf("StrSliceToRRuleSetInLoc failed: %v", err)
	}
	if !set.IsAllDay() {
		t.Fatal("Set should be marked all-day when RDATE uses VALUE=DATE")
	}

	rdates := set.GetRDate()
	if len(rdates) != 2 {
		t.Fatalf("Expected 2 RDATEs, got %d", len(rdates))
	}
	want := []time.Time{
		time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 3, 3, 0, 0, 0, 0, time.UTC),
	}
	if !timesEqual(rdates, want) {
		t.Fatalf("RDATEs were not normalized to floating midnight: %v", rdates)
	}

	recurrence := set.RRuleString()
	if !strings.Contains(recurrence, "RDATE;VALUE=DATE:20240301") ||
		!strings.Contains(recurrence, "RDATE;VALUE=DATE:20240303") {
		t.Fatalf("Expected VALUE=DATE serialization, got %q", recurrence)
	}
}

func TestStrSliceToRRuleSetDetectsAllDayFromExDate(t *testing.T) {
	lines := []string{
		"EXDATE;VALUE=DATE:20250110,20250112",
	}

	set, err := StrSliceToRRuleSetInLoc(lines, time.UTC)
	if err != nil {
		t.Fatalf("StrSliceToRRuleSetInLoc failed: %v", err)
	}
	if !set.IsAllDay() {
		t.Fatal("Set should be marked all-day when EXDATE uses VALUE=DATE")
	}

	exdates := set.GetExDate()
	if len(exdates) != 2 {
		t.Fatalf("Expected 2 EXDATEs, got %d", len(exdates))
	}
	want := []time.Time{
		time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 12, 0, 0, 0, 0, time.UTC),
	}
	if !timesEqual(exdates, want) {
		t.Fatalf("EXDATEs were not normalized to floating midnight: %v", exdates)
	}

	recurrence := set.RRuleString()
	if !strings.Contains(recurrence, "EXDATE;VALUE=DATE:20250110") ||
		!strings.Contains(recurrence, "EXDATE;VALUE=DATE:20250112") {
		t.Fatalf("Expected VALUE=DATE serialization, got %q", recurrence)
	}
}
