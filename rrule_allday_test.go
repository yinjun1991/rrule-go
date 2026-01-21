package rrule

import (
	"testing"
	"time"
)

// TestAllDayDaily tests daily recurrence for all-day events.
func TestAllDayDaily(t *testing.T) {
	r, _ := NewRRule(ROption{
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
	r, _ := NewRRule(ROption{
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
	r, _ := NewRRule(ROption{
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
	r, _ := NewRRule(ROption{
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
	r, _ := NewRRule(ROption{
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
	r, _ := NewRRule(ROption{
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
	r, _ := NewRRule(ROption{
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
			r, _ := NewRRule(ROption{
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
	set := Set{}

	// Add base rule.
	r, _ := NewRRule(ROption{
		Freq:    WEEKLY,
		Count:   2,
		AllDay:  true,
		Dtstart: time.Date(2023, 1, 1, 14, 30, 0, 0, time.UTC),
	})
	set.RRule(r)

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
	set := Set{}

	// Add base rule.
	r, _ := NewRRule(ROption{
		Freq:    DAILY,
		Count:   5,
		AllDay:  true,
		Dtstart: time.Date(2023, 1, 1, 11, 20, 0, 0, time.UTC),
	})
	set.RRule(r)

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
	set := Set{}

	// Add base rule.
	r, _ := NewRRule(ROption{
		Freq:    WEEKLY,
		Count:   3,
		AllDay:  true,
		Dtstart: time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
	})
	set.RRule(r)

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
			r, _ := NewRRule(ROption{
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
	r, _ := NewRRule(ROption{
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
	r, _ := NewRRule(ROption{
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
	r, _ := NewRRule(ROption{
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
	r, _ := NewRRule(ROption{
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
	r, _ := NewRRule(ROption{
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
	r, _ := NewRRule(ROption{
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
	allDayRule, _ := NewRRule(ROption{
		Freq:    DAILY,
		Count:   3,
		AllDay:  true,
		Dtstart: time.Date(2023, 1, 1, 14, 30, 0, 0, time.UTC),
	})

	// Timed event at midnight on the same date.
	nonAllDayRule, _ := NewRRule(ROption{
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
