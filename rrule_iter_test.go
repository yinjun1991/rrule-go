package rrule

import (
	"testing"
	"time"
)

// TestIteratorFrequencyAdvanced tests advanced frequency scenarios: intervals, year/month boundaries, etc.
func TestIteratorFrequencyAdvanced(t *testing.T) {
	tests := []struct {
		name     string
		opt      ROption
		expected []time.Time
	}{
		{
			name: "YEARLY_with_interval",
			opt: ROption{
				Freq:     YEARLY,
				Interval: 2,
				Count:    3,
				Dtstart:  time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
			},
			expected: []time.Time{
				time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC),
				time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "MONTHLY_cross_year",
			opt: ROption{
				Freq:    MONTHLY,
				Count:   3,
				Dtstart: time.Date(2020, 11, 15, 14, 30, 0, 0, time.UTC),
			},
			expected: []time.Time{
				time.Date(2020, 11, 15, 14, 30, 0, 0, time.UTC),
				time.Date(2020, 12, 15, 14, 30, 0, 0, time.UTC),
				time.Date(2021, 1, 15, 14, 30, 0, 0, time.UTC),
			},
		},
		{
			name: "WEEKLY_cross_year",
			opt: ROption{
				Freq:    WEEKLY,
				Count:   3,
				Dtstart: time.Date(2020, 12, 28, 9, 0, 0, 0, time.UTC), // Monday
			},
			expected: []time.Time{
				time.Date(2020, 12, 28, 9, 0, 0, 0, time.UTC),
				time.Date(2021, 1, 4, 9, 0, 0, 0, time.UTC),
				time.Date(2021, 1, 11, 9, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "DAILY_cross_month",
			opt: ROption{
				Freq:    DAILY,
				Count:   3,
				Dtstart: time.Date(2020, 1, 30, 8, 0, 0, 0, time.UTC),
			},
			expected: []time.Time{
				time.Date(2020, 1, 30, 8, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 31, 8, 0, 0, 0, time.UTC),
				time.Date(2020, 2, 1, 8, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "HOURLY_cross_day",
			opt: ROption{
				Freq:    HOURLY,
				Count:   3,
				Dtstart: time.Date(2020, 1, 1, 23, 30, 0, 0, time.UTC),
			},
			expected: []time.Time{
				time.Date(2020, 1, 1, 23, 30, 0, 0, time.UTC),
				time.Date(2020, 1, 2, 0, 30, 0, 0, time.UTC),
				time.Date(2020, 1, 2, 1, 30, 0, 0, time.UTC),
			},
		},
		{
			name: "MINUTELY_cross_hour",
			opt: ROption{
				Freq:    MINUTELY,
				Count:   3,
				Dtstart: time.Date(2020, 1, 1, 10, 59, 15, 0, time.UTC),
			},
			expected: []time.Time{
				time.Date(2020, 1, 1, 10, 59, 15, 0, time.UTC),
				time.Date(2020, 1, 1, 11, 0, 15, 0, time.UTC),
				time.Date(2020, 1, 1, 11, 1, 15, 0, time.UTC),
			},
		},
		{
			name: "SECONDLY_cross_minute",
			opt: ROption{
				Freq:    SECONDLY,
				Count:   3,
				Dtstart: time.Date(2020, 1, 1, 10, 30, 59, 0, time.UTC),
			},
			expected: []time.Time{
				time.Date(2020, 1, 1, 10, 30, 59, 0, time.UTC),
				time.Date(2020, 1, 1, 10, 31, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 10, 31, 1, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := newRecurrence(tt.opt)
			if err != nil {
				t.Fatalf("Failed to create RRule: %v", err)
			}
			result := r.All()
			if !timesEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestIteratorBasicFrequencies tests iterator behavior for basic frequencies.
func TestIteratorBasicFrequencies(t *testing.T) {
	tests := []struct {
		name     string
		opt      ROption
		expected []time.Time
	}{
		{
			name: "YEARLY_basic",
			opt: ROption{
				Freq:    YEARLY,
				Count:   3,
				Dtstart: time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
			},
			expected: []time.Time{
				time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				time.Date(2021, 1, 1, 10, 0, 0, 0, time.UTC),
				time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "MONTHLY_basic",
			opt: ROption{
				Freq:    MONTHLY,
				Count:   3,
				Dtstart: time.Date(2020, 1, 15, 14, 30, 0, 0, time.UTC),
			},
			expected: []time.Time{
				time.Date(2020, 1, 15, 14, 30, 0, 0, time.UTC),
				time.Date(2020, 2, 15, 14, 30, 0, 0, time.UTC),
				time.Date(2020, 3, 15, 14, 30, 0, 0, time.UTC),
			},
		},
		{
			name: "WEEKLY_basic",
			opt: ROption{
				Freq:    WEEKLY,
				Count:   3,
				Dtstart: time.Date(2020, 1, 6, 9, 0, 0, 0, time.UTC), // Monday
			},
			expected: []time.Time{
				time.Date(2020, 1, 6, 9, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 13, 9, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 20, 9, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "DAILY_basic",
			opt: ROption{
				Freq:    DAILY,
				Count:   3,
				Dtstart: time.Date(2020, 1, 1, 8, 0, 0, 0, time.UTC),
			},
			expected: []time.Time{
				time.Date(2020, 1, 1, 8, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 2, 8, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 3, 8, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "HOURLY_basic",
			opt: ROption{
				Freq:    HOURLY,
				Count:   3,
				Dtstart: time.Date(2020, 1, 1, 10, 30, 0, 0, time.UTC),
			},
			expected: []time.Time{
				time.Date(2020, 1, 1, 10, 30, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 11, 30, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 12, 30, 0, 0, time.UTC),
			},
		},
		{
			name: "MINUTELY_basic",
			opt: ROption{
				Freq:    MINUTELY,
				Count:   3,
				Dtstart: time.Date(2020, 1, 1, 10, 30, 15, 0, time.UTC),
			},
			expected: []time.Time{
				time.Date(2020, 1, 1, 10, 30, 15, 0, time.UTC),
				time.Date(2020, 1, 1, 10, 31, 15, 0, time.UTC),
				time.Date(2020, 1, 1, 10, 32, 15, 0, time.UTC),
			},
		},
		{
			name: "SECONDLY_basic",
			opt: ROption{
				Freq:    SECONDLY,
				Count:   3,
				Dtstart: time.Date(2020, 1, 1, 10, 30, 15, 0, time.UTC),
			},
			expected: []time.Time{
				time.Date(2020, 1, 1, 10, 30, 15, 0, time.UTC),
				time.Date(2020, 1, 1, 10, 30, 16, 0, time.UTC),
				time.Date(2020, 1, 1, 10, 30, 17, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := newRecurrence(tt.opt)
			if err != nil {
				t.Fatalf("Failed to create RRule: %v", err)
			}
			result := r.All()
			if !timesEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestIteratorByRulesCombinations tests BY* rule combinations.
func TestIteratorByRulesCombinations(t *testing.T) {
	tests := []struct {
		name     string
		opt      ROption
		expected []time.Time
	}{
		{
			name: "BYMONTH_single",
			opt: ROption{
				Freq:    YEARLY,
				Count:   3,
				Dtstart: time.Date(2020, 1, 15, 10, 0, 0, 0, time.UTC),
				Bymonth: []int{6}, // Only in June.
			},
			expected: []time.Time{
				time.Date(2020, 6, 15, 10, 0, 0, 0, time.UTC),
				time.Date(2021, 6, 15, 10, 0, 0, 0, time.UTC),
				time.Date(2022, 6, 15, 10, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "BYMONTH_multiple",
			opt: ROption{
				Freq:    YEARLY,
				Count:   4,
				Dtstart: time.Date(2020, 1, 15, 10, 0, 0, 0, time.UTC),
				Bymonth: []int{3, 9}, // March and September.
			},
			expected: []time.Time{
				time.Date(2020, 3, 15, 10, 0, 0, 0, time.UTC),
				time.Date(2020, 9, 15, 10, 0, 0, 0, time.UTC),
				time.Date(2021, 3, 15, 10, 0, 0, 0, time.UTC),
				time.Date(2021, 9, 15, 10, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "BYWEEKDAY_single",
			opt: ROption{
				Freq:      WEEKLY,
				Count:     3,
				Dtstart:   time.Date(2020, 1, 6, 10, 0, 0, 0, time.UTC), // Monday
				Byweekday: []Weekday{FR},                                // Only Fridays.
			},
			expected: []time.Time{
				time.Date(2020, 1, 10, 10, 0, 0, 0, time.UTC), // First Friday.
				time.Date(2020, 1, 17, 10, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 24, 10, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "BYWEEKDAY_multiple",
			opt: ROption{
				Freq:      WEEKLY,
				Count:     4,
				Dtstart:   time.Date(2020, 1, 6, 10, 0, 0, 0, time.UTC), // Monday
				Byweekday: []Weekday{MO, WE, FR},                        // Mon, Wed, Fri.
			},
			expected: []time.Time{
				time.Date(2020, 1, 6, 10, 0, 0, 0, time.UTC),  // Monday
				time.Date(2020, 1, 8, 10, 0, 0, 0, time.UTC),  // Wednesday
				time.Date(2020, 1, 10, 10, 0, 0, 0, time.UTC), // Friday
				time.Date(2020, 1, 13, 10, 0, 0, 0, time.UTC), // Next Monday
			},
		},
		{
			name: "BYMONTHDAY_positive",
			opt: ROption{
				Freq:       MONTHLY,
				Count:      3,
				Dtstart:    time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				Bymonthday: []int{15}, // 15th of each month.
			},
			expected: []time.Time{
				time.Date(2020, 1, 15, 10, 0, 0, 0, time.UTC),
				time.Date(2020, 2, 15, 10, 0, 0, 0, time.UTC),
				time.Date(2020, 3, 15, 10, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "BYMONTHDAY_negative",
			opt: ROption{
				Freq:       MONTHLY,
				Count:      3,
				Dtstart:    time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				Bymonthday: []int{-1}, // Last day of each month.
			},
			expected: []time.Time{
				time.Date(2020, 1, 31, 10, 0, 0, 0, time.UTC),
				time.Date(2020, 2, 29, 10, 0, 0, 0, time.UTC), // Leap-year February.
				time.Date(2020, 3, 31, 10, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "BYHOUR_multiple",
			opt: ROption{
				Freq:    DAILY,
				Count:   4,
				Dtstart: time.Date(2020, 1, 1, 8, 0, 0, 0, time.UTC),
				Byhour:  []int{9, 15}, // 9:00 and 15:00.
			},
			expected: []time.Time{
				time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 15, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 2, 9, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 2, 15, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "BYMINUTE_multiple",
			opt: ROption{
				Freq:     HOURLY,
				Count:    4,
				Dtstart:  time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				Byminute: []int{15, 45}, // 15 and 45 minutes.
			},
			expected: []time.Time{
				time.Date(2020, 1, 1, 10, 15, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 10, 45, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 11, 15, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 11, 45, 0, 0, time.UTC),
			},
		},
		{
			name: "BYSECOND_multiple",
			opt: ROption{
				Freq:     MINUTELY,
				Count:    4,
				Dtstart:  time.Date(2020, 1, 1, 10, 30, 0, 0, time.UTC),
				Bysecond: []int{10, 50}, // 10 and 50 seconds.
			},
			expected: []time.Time{
				time.Date(2020, 1, 1, 10, 30, 10, 0, time.UTC),
				time.Date(2020, 1, 1, 10, 30, 50, 0, time.UTC),
				time.Date(2020, 1, 1, 10, 31, 10, 0, time.UTC),
				time.Date(2020, 1, 1, 10, 31, 50, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := newRecurrence(tt.opt)
			if err != nil {
				t.Fatalf("Failed to create RRule: %v", err)
			}
			result := r.All()
			if !timesEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestIteratorAllDayEvents tests the iterator for all-day events.
func TestIteratorAllDayEvents(t *testing.T) {
	tests := []struct {
		name     string
		opt      ROption
		expected []time.Time
	}{
		{
			name: "AllDay_DAILY",
			opt: ROption{
				Freq:    DAILY,
				Count:   3,
				AllDay:  true,
				Dtstart: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: []time.Time{
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "AllDay_WEEKLY",
			opt: ROption{
				Freq:    WEEKLY,
				Count:   3,
				AllDay:  true,
				Dtstart: time.Date(2020, 1, 6, 0, 0, 0, 0, time.UTC), // Monday
			},
			expected: []time.Time{
				time.Date(2020, 1, 6, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 13, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 20, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "AllDay_MONTHLY",
			opt: ROption{
				Freq:    MONTHLY,
				Count:   3,
				AllDay:  true,
				Dtstart: time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC),
			},
			expected: []time.Time{
				time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 2, 15, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 3, 15, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := newRecurrence(tt.opt)
			if err != nil {
				t.Fatalf("Failed to create RRule: %v", err)
			}
			result := r.All()
			if !timesEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestIteratorTimezones tests multi-timezone handling.
func TestIteratorTimezones(t *testing.T) {
	// Load timezones.
	nyc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Skip("Skipping timezone test: America/New_York not available")
	}
	tokyo, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Skip("Skipping timezone test: Asia/Tokyo not available")
	}

	tests := []struct {
		name     string
		opt      ROption
		expected []time.Time
	}{
		{
			name: "NYC_timezone",
			opt: ROption{
				Freq:    DAILY,
				Count:   2,
				Dtstart: time.Date(2020, 1, 1, 10, 0, 0, 0, nyc),
			},
			expected: []time.Time{
				time.Date(2020, 1, 1, 10, 0, 0, 0, nyc),
				time.Date(2020, 1, 2, 10, 0, 0, 0, nyc),
			},
		},
		{
			name: "Tokyo_timezone",
			opt: ROption{
				Freq:    WEEKLY,
				Count:   2,
				Dtstart: time.Date(2020, 1, 6, 15, 30, 0, 0, tokyo),
			},
			expected: []time.Time{
				time.Date(2020, 1, 6, 15, 30, 0, 0, tokyo),
				time.Date(2020, 1, 13, 15, 30, 0, 0, tokyo),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := newRecurrence(tt.opt)
			if err != nil {
				t.Fatalf("Failed to create RRule: %v", err)
			}
			result := r.All()
			if !timesEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestIteratorByRules tests BY* rule combinations.
func TestIteratorByRules(t *testing.T) {
	tests := []struct {
		name     string
		opt      ROption
		expected []time.Time
	}{
		{
			name: "BYMONTH_yearly",
			opt: ROption{
				Freq:    YEARLY,
				Count:   3,
				Bymonth: []int{3, 6, 9},
				Dtstart: time.Date(2020, 1, 15, 10, 0, 0, 0, time.UTC),
			},
			expected: []time.Time{
				time.Date(2020, 3, 15, 10, 0, 0, 0, time.UTC),
				time.Date(2020, 6, 15, 10, 0, 0, 0, time.UTC),
				time.Date(2020, 9, 15, 10, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "BYWEEKDAY_weekly",
			opt: ROption{
				Freq:      WEEKLY,
				Count:     4,
				Byweekday: []Weekday{MO, WE, FR},
				Dtstart:   time.Date(2020, 1, 6, 9, 0, 0, 0, time.UTC), // Monday
			},
			expected: []time.Time{
				time.Date(2020, 1, 6, 9, 0, 0, 0, time.UTC),  // MO
				time.Date(2020, 1, 8, 9, 0, 0, 0, time.UTC),  // WE
				time.Date(2020, 1, 10, 9, 0, 0, 0, time.UTC), // FR
				time.Date(2020, 1, 13, 9, 0, 0, 0, time.UTC), // MO next week
			},
		},
		{
			name: "BYMONTHDAY_monthly",
			opt: ROption{
				Freq:       MONTHLY,
				Count:      4,
				Bymonthday: []int{1, 15},
				Dtstart:    time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
			},
			expected: []time.Time{
				time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 15, 10, 0, 0, 0, time.UTC),
				time.Date(2020, 2, 1, 10, 0, 0, 0, time.UTC),
				time.Date(2020, 2, 15, 10, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "BYHOUR_daily",
			opt: ROption{
				Freq:    DAILY,
				Count:   4,
				Byhour:  []int{9, 15},
				Dtstart: time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
			},
			expected: []time.Time{
				time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 15, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 2, 9, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 2, 15, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := newRecurrence(tt.opt)
			if err != nil {
				t.Fatalf("Failed to create RRule: %v", err)
			}
			result := r.All()
			if !timesEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestIteratorBoundaryConditions tests boundary conditions.
func TestIteratorBoundaryConditions(t *testing.T) {
	tests := []struct {
		name     string
		opt      ROption
		expected []time.Time
	}{
		{
			name: "LeapYear_Feb29",
			opt: ROption{
				Freq:       YEARLY,
				Count:      2,
				Bymonth:    []int{2},
				Bymonthday: []int{29},
				Dtstart:    time.Date(2020, 2, 29, 10, 0, 0, 0, time.UTC), // 2020 is leap year
			},
			expected: []time.Time{
				time.Date(2020, 2, 29, 10, 0, 0, 0, time.UTC),
				time.Date(2024, 2, 29, 10, 0, 0, 0, time.UTC), // Next leap year
			},
		},
		{
			name: "YearEnd_crossover",
			opt: ROption{
				Freq:    DAILY,
				Count:   3,
				Dtstart: time.Date(2020, 12, 30, 23, 59, 59, 0, time.UTC),
			},
			expected: []time.Time{
				time.Date(2020, 12, 30, 23, 59, 59, 0, time.UTC),
				time.Date(2020, 12, 31, 23, 59, 59, 0, time.UTC),
				time.Date(2021, 1, 1, 23, 59, 59, 0, time.UTC),
			},
		},
		{
			name: "MonthEnd_February",
			opt: ROption{
				Freq:       MONTHLY,
				Count:      3,
				Bymonthday: []int{28},
				Dtstart:    time.Date(2021, 1, 28, 12, 0, 0, 0, time.UTC), // 2021 is not leap year
			},
			expected: []time.Time{
				time.Date(2021, 1, 28, 12, 0, 0, 0, time.UTC),
				time.Date(2021, 2, 28, 12, 0, 0, 0, time.UTC),
				time.Date(2021, 3, 28, 12, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := newRecurrence(tt.opt)
			if err != nil {
				t.Fatalf("Failed to create RRule: %v", err)
			}
			result := r.All()
			if !timesEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestIteratorInterval tests interval settings.
func TestIteratorInterval(t *testing.T) {
	tests := []struct {
		name     string
		opt      ROption
		expected []time.Time
	}{
		{
			name: "DAILY_interval_2",
			opt: ROption{
				Freq:     DAILY,
				Interval: 2,
				Count:    3,
				Dtstart:  time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
			},
			expected: []time.Time{
				time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 3, 10, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 5, 10, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "WEEKLY_interval_3",
			opt: ROption{
				Freq:     WEEKLY,
				Interval: 3,
				Count:    3,
				Dtstart:  time.Date(2020, 1, 6, 9, 0, 0, 0, time.UTC), // Monday
			},
			expected: []time.Time{
				time.Date(2020, 1, 6, 9, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 27, 9, 0, 0, 0, time.UTC), // 3 weeks later
				time.Date(2020, 2, 17, 9, 0, 0, 0, time.UTC), // 6 weeks later
			},
		},
		{
			name: "MONTHLY_interval_2",
			opt: ROption{
				Freq:     MONTHLY,
				Interval: 2,
				Count:    3,
				Dtstart:  time.Date(2020, 1, 15, 14, 0, 0, 0, time.UTC),
			},
			expected: []time.Time{
				time.Date(2020, 1, 15, 14, 0, 0, 0, time.UTC),
				time.Date(2020, 3, 15, 14, 0, 0, 0, time.UTC),
				time.Date(2020, 5, 15, 14, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "HOURLY_interval_6",
			opt: ROption{
				Freq:     HOURLY,
				Interval: 6,
				Count:    4,
				Dtstart:  time.Date(2020, 1, 1, 6, 0, 0, 0, time.UTC),
			},
			expected: []time.Time{
				time.Date(2020, 1, 1, 6, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 18, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := newRecurrence(tt.opt)
			if err != nil {
				t.Fatalf("Failed to create RRule: %v", err)
			}
			result := r.All()
			if !timesEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestIteratorUntil tests UNTIL limits.
func TestIteratorUntil(t *testing.T) {
	tests := []struct {
		name     string
		opt      ROption
		expected []time.Time
	}{
		{
			name: "DAILY_until",
			opt: ROption{
				Freq:    DAILY,
				Until:   time.Date(2020, 1, 5, 10, 0, 0, 0, time.UTC),
				Dtstart: time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
			},
			expected: []time.Time{
				time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 2, 10, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 3, 10, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 4, 10, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 5, 10, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "WEEKLY_until",
			opt: ROption{
				Freq:    WEEKLY,
				Until:   time.Date(2020, 1, 20, 9, 0, 0, 0, time.UTC),
				Dtstart: time.Date(2020, 1, 6, 9, 0, 0, 0, time.UTC),
			},
			expected: []time.Time{
				time.Date(2020, 1, 6, 9, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 13, 9, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 20, 9, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := newRecurrence(tt.opt)
			if err != nil {
				t.Fatalf("Failed to create RRule: %v", err)
			}
			result := r.All()
			if !timesEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestIteratorComplexByRules tests complex BY* rule combinations.
func TestIteratorComplexByRules(t *testing.T) {
	tests := []struct {
		name     string
		opt      ROption
		expected []time.Time
	}{
		{
			name: "BYSETPOS_first_and_last",
			opt: ROption{
				Freq:      MONTHLY,
				Count:     4,
				Byweekday: []Weekday{MO, TU, WE, TH, FR}, // Weekdays
				Bysetpos:  []int{1, -1},                  // First and last
				Dtstart:   time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
			},
			expected: []time.Time{
				time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),  // First weekday of Jan
				time.Date(2020, 1, 31, 9, 0, 0, 0, time.UTC), // Last weekday of Jan
				time.Date(2020, 2, 3, 9, 0, 0, 0, time.UTC),  // First weekday of Feb
				time.Date(2020, 2, 28, 9, 0, 0, 0, time.UTC), // Last weekday of Feb
			},
		},
		{
			name: "BYMONTH_and_BYWEEKDAY",
			opt: ROption{
				Freq:      YEARLY,
				Count:     4,
				Bymonth:   []int{3, 6, 9, 12},
				Byweekday: []Weekday{FR},
				Dtstart:   time.Date(2020, 1, 1, 15, 0, 0, 0, time.UTC),
			},
			expected: []time.Time{
				time.Date(2020, 3, 6, 15, 0, 0, 0, time.UTC),
				time.Date(2020, 3, 13, 15, 0, 0, 0, time.UTC),
				time.Date(2020, 3, 20, 15, 0, 0, 0, time.UTC),
				time.Date(2020, 3, 27, 15, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := newRecurrence(tt.opt)
			if err != nil {
				t.Fatalf("Failed to create RRule: %v", err)
			}
			result := r.All()
			if !timesEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestIteratorLeapYearAndBoundaries tests leap years and boundaries.
func TestIteratorLeapYearAndBoundaries(t *testing.T) {
	tests := []struct {
		name     string
		opt      ROption
		expected []time.Time
	}{
		{
			name: "leap_year_feb_29",
			opt: ROption{
				Freq:    YEARLY,
				Count:   3,
				Dtstart: time.Date(2020, 2, 29, 12, 0, 0, 0, time.UTC), // Leap day.
			},
			expected: []time.Time{
				time.Date(2020, 2, 29, 12, 0, 0, 0, time.UTC),
				time.Date(2024, 2, 29, 12, 0, 0, 0, time.UTC), // Next leap year.
				time.Date(2028, 2, 29, 12, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "month_boundary_31_to_30_days",
			opt: ROption{
				Freq:    MONTHLY,
				Count:   3,
				Dtstart: time.Date(2020, 1, 31, 15, 0, 0, 0, time.UTC), // Jan 31.
			},
			expected: []time.Time{
				time.Date(2020, 1, 31, 15, 0, 0, 0, time.UTC),
				time.Date(2020, 3, 31, 15, 0, 0, 0, time.UTC), // Skip February (28/29 days).
				time.Date(2020, 5, 31, 15, 0, 0, 0, time.UTC), // Skip April (30 days).
			},
		},
		{
			name: "year_boundary_december_to_january",
			opt: ROption{
				Freq:    MONTHLY,
				Count:   3,
				Dtstart: time.Date(2020, 12, 15, 10, 30, 0, 0, time.UTC),
			},
			expected: []time.Time{
				time.Date(2020, 12, 15, 10, 30, 0, 0, time.UTC),
				time.Date(2021, 1, 15, 10, 30, 0, 0, time.UTC),
				time.Date(2021, 2, 15, 10, 30, 0, 0, time.UTC),
			},
		},
		{
			name: "week_boundary_sunday_to_monday",
			opt: ROption{
				Freq:    WEEKLY,
				Count:   3,
				Dtstart: time.Date(2020, 12, 27, 16, 0, 0, 0, time.UTC), // Sunday
			},
			expected: []time.Time{
				time.Date(2020, 12, 27, 16, 0, 0, 0, time.UTC),
				time.Date(2021, 1, 3, 16, 0, 0, 0, time.UTC),
				time.Date(2021, 1, 10, 16, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "day_boundary_23_59_to_00_00",
			opt: ROption{
				Freq:    DAILY,
				Count:   3,
				Dtstart: time.Date(2020, 12, 31, 23, 59, 0, 0, time.UTC),
			},
			expected: []time.Time{
				time.Date(2020, 12, 31, 23, 59, 0, 0, time.UTC),
				time.Date(2021, 1, 1, 23, 59, 0, 0, time.UTC),
				time.Date(2021, 1, 2, 23, 59, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := newRecurrence(tt.opt)
			if err != nil {
				t.Fatalf("Failed to create RRule: %v", err)
			}
			result := r.All()
			if !timesEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestIteratorEdgeCases tests edge cases.
func TestIteratorEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		opt      ROption
		expected []time.Time
	}{
		{
			name: "Empty_result_impossible_date",
			opt: ROption{
				Freq:       YEARLY,
				Count:      3,
				Bymonth:    []int{2},
				Bymonthday: []int{30}, // February 30th doesn't exist
				Dtstart:    time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
			},
			expected: []time.Time{}, // Should be empty
		},
		{
			name: "Cross_year_weekly",
			opt: ROption{
				Freq:    WEEKLY,
				Count:   3,
				Dtstart: time.Date(2020, 12, 28, 10, 0, 0, 0, time.UTC), // Monday
			},
			expected: []time.Time{
				time.Date(2020, 12, 28, 10, 0, 0, 0, time.UTC),
				time.Date(2021, 1, 4, 10, 0, 0, 0, time.UTC),
				time.Date(2021, 1, 11, 10, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := newRecurrence(tt.opt)
			if err != nil {
				t.Fatalf("Failed to create RRule: %v", err)
			}
			result := r.All()
			if !timesEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestIteratorAllDayVsTimedEvents tests differences between all-day and timed events.
func TestIteratorAllDayVsTimedEvents(t *testing.T) {
	tests := []struct {
		name     string
		opt      ROption
		expected []time.Time
	}{
		{
			name: "allday_daily_event",
			opt: ROption{
				Freq:    DAILY,
				Count:   3,
				Dtstart: time.Date(2020, 1, 1, 14, 30, 45, 0, time.UTC),
				AllDay:  true,
			},
			expected: []time.Time{
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), // All-day events normalize to 00:00:00.
				time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "timed_daily_event",
			opt: ROption{
				Freq:    DAILY,
				Count:   3,
				Dtstart: time.Date(2020, 1, 1, 14, 30, 45, 0, time.UTC),
				AllDay:  false,
			},
			expected: []time.Time{
				time.Date(2020, 1, 1, 14, 30, 45, 0, time.UTC), // Keep original time.
				time.Date(2020, 1, 2, 14, 30, 45, 0, time.UTC),
				time.Date(2020, 1, 3, 14, 30, 45, 0, time.UTC),
			},
		},
		{
			name: "allday_with_until",
			opt: ROption{
				Freq:    DAILY,
				Dtstart: time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				Until:   time.Date(2020, 1, 3, 15, 30, 0, 0, time.UTC), // UNTIL time is normalized too.
				AllDay:  true,
			},
			expected: []time.Time{
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "allday_hourly_becomes_daily",
			opt: ROption{
				Freq:    HOURLY,
				Count:   3,
				Dtstart: time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				AllDay:  true,
			},
			expected: []time.Time{
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "allday_with_byhour",
			opt: ROption{
				Freq:    DAILY,
				Count:   2,
				Dtstart: time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				Byhour:  []int{9, 15}, // All-day events should ignore BYHOUR.
				AllDay:  true,
			},
			expected: []time.Time{
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := newRecurrence(tt.opt)
			if err != nil {
				t.Fatalf("Failed to create RRule: %v", err)
			}
			result := r.All()
			if !timesEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestIteratorTimezoneHandling tests multi-timezone handling.
func TestIteratorTimezoneHandling(t *testing.T) {
	// Create different timezones.
	utc := time.UTC
	east, _ := time.LoadLocation("Asia/Shanghai")    // UTC+8
	west, _ := time.LoadLocation("America/New_York") // UTC-5/-4

	tests := []struct {
		name     string
		opt      ROption
		expected []time.Time
	}{
		{
			name: "utc_daily",
			opt: ROption{
				Freq:    DAILY,
				Count:   3,
				Dtstart: time.Date(2020, 1, 1, 10, 0, 0, 0, utc),
			},
			expected: []time.Time{
				time.Date(2020, 1, 1, 10, 0, 0, 0, utc),
				time.Date(2020, 1, 2, 10, 0, 0, 0, utc),
				time.Date(2020, 1, 3, 10, 0, 0, 0, utc),
			},
		},
		{
			name: "positive_offset_daily",
			opt: ROption{
				Freq:    DAILY,
				Count:   3,
				Dtstart: time.Date(2020, 1, 1, 10, 0, 0, 0, east),
			},
			expected: []time.Time{
				time.Date(2020, 1, 1, 10, 0, 0, 0, east),
				time.Date(2020, 1, 2, 10, 0, 0, 0, east),
				time.Date(2020, 1, 3, 10, 0, 0, 0, east),
			},
		},
		{
			name: "negative_offset_daily",
			opt: ROption{
				Freq:    DAILY,
				Count:   3,
				Dtstart: time.Date(2020, 1, 1, 10, 0, 0, 0, west),
			},
			expected: []time.Time{
				time.Date(2020, 1, 1, 10, 0, 0, 0, west),
				time.Date(2020, 1, 2, 10, 0, 0, 0, west),
				time.Date(2020, 1, 3, 10, 0, 0, 0, west),
			},
		},
		{
			name: "timezone_boundary_hourly",
			opt: ROption{
				Freq:    HOURLY,
				Count:   5,
				Dtstart: time.Date(2020, 1, 1, 22, 0, 0, 0, east), // 22:00 in Shanghai
			},
			expected: []time.Time{
				time.Date(2020, 1, 1, 22, 0, 0, 0, east),
				time.Date(2020, 1, 1, 23, 0, 0, 0, east),
				time.Date(2020, 1, 2, 0, 0, 0, 0, east), // Crosses day boundary.
				time.Date(2020, 1, 2, 1, 0, 0, 0, east),
				time.Date(2020, 1, 2, 2, 0, 0, 0, east),
			},
		},
		{
			name: "allday_timezone_independence",
			opt: ROption{
				Freq:    DAILY,
				Count:   2,
				Dtstart: time.Date(2020, 1, 1, 15, 30, 0, 0, east),
				AllDay:  true, // All-day events should be timezone-independent.
			},
			expected: []time.Time{
				time.Date(2020, 1, 1, 0, 0, 0, 0, utc), // Converted to UTC 00:00.
				time.Date(2020, 1, 2, 0, 0, 0, 0, utc),
			},
		},
		{
			name: "mixed_timezone_until",
			opt: ROption{
				Freq:    DAILY,
				Dtstart: time.Date(2020, 1, 1, 10, 0, 0, 0, east),
				Until:   time.Date(2020, 1, 2, 23, 0, 0, 0, utc), // UTC time; next day 07:00 in UTC+8.
			},
			expected: []time.Time{
				time.Date(2020, 1, 1, 10, 0, 0, 0, east),
				time.Date(2020, 1, 2, 10, 0, 0, 0, east),
				// 2020-01-03 10:00 UTC+8 = 2020-01-03 02:00 UTC, after the Until time.
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := newRecurrence(tt.opt)
			if err != nil {
				t.Fatalf("Failed to create RRule: %v", err)
			}
			result := r.All()
			if !timesEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestIteratorPerformance tests iterator performance and memory usage.
func TestIteratorPerformance(t *testing.T) {
	// Test performance with many iterations.
	r, err := newRecurrence(ROption{
		Freq:    DAILY,
		Count:   1000,
		Dtstart: time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("Failed to create RRule: %v", err)
	}

	result := r.All()
	if len(result) != 1000 {
		t.Errorf("Expected 1000 results, got %d", len(result))
	}

	// Verify results are continuous.
	for i := 1; i < len(result); i++ {
		expected := result[i-1].AddDate(0, 0, 1)
		if !result[i].Equal(expected) {
			t.Errorf("Result %d: expected %v, got %v", i, expected, result[i])
			break
		}
	}
}

// TestIteratorNext tests the iterator's next() method.
func TestIteratorNext(t *testing.T) {
	r, err := newRecurrence(ROption{
		Freq:    DAILY,
		Count:   3,
		Dtstart: time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("Failed to create RRule: %v", err)
	}

	iter := r.Iterator()
	expected := []time.Time{
		time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
		time.Date(2020, 1, 2, 10, 0, 0, 0, time.UTC),
		time.Date(2020, 1, 3, 10, 0, 0, 0, time.UTC),
	}

	for i, exp := range expected {
		val, ok := iter()
		if !ok {
			t.Fatalf("Iterator ended prematurely at index %d", i)
		}
		if !val.Equal(exp) {
			t.Errorf("Index %d: expected %v, got %v", i, exp, val)
		}
	}

	// There should be no more values.
	_, ok := iter()
	if ok {
		t.Error("Iterator should have ended")
	}
}

// TestIteratorMaxYear tests the max year limit.
func TestIteratorMaxYear(t *testing.T) {
	r, err := newRecurrence(ROption{
		Freq:    YEARLY,
		Count:   10,
		Dtstart: time.Date(9995, 1, 1, 10, 0, 0, 0, time.UTC), // Close to MAXYEAR
	})
	if err != nil {
		t.Fatalf("Failed to create RRule: %v", err)
	}

	result := r.All()
	// Should stop at MAXYEAR.
	for _, dt := range result {
		if dt.Year() > MAXYEAR {
			t.Errorf("Result year %d exceeds MAXYEAR %d", dt.Year(), MAXYEAR)
		}
	}
}

// TestIteratorErrorAndNegativeCases tests error and negative cases.
func TestIteratorErrorAndNegativeCases(t *testing.T) {
	tests := []struct {
		name        string
		opt         ROption
		expectError bool
		expected    []time.Time
	}{
		{
			name: "invalid_bymonth_zero",
			opt: ROption{
				Freq:    YEARLY,
				Count:   1,
				Dtstart: time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				Bymonth: []int{0}, // Invalid month.
			},
			expectError: true,
		},
		{
			name: "invalid_bymonth_negative",
			opt: ROption{
				Freq:    YEARLY,
				Count:   1,
				Dtstart: time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				Bymonth: []int{-1}, // Invalid month.
			},
			expectError: true,
		},
		{
			name: "invalid_byhour_negative",
			opt: ROption{
				Freq:    DAILY,
				Count:   1,
				Dtstart: time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				Byhour:  []int{-1}, // Invalid hour.
			},
			expectError: true,
		},
		{
			name: "invalid_byminute_negative",
			opt: ROption{
				Freq:     HOURLY,
				Count:    1,
				Dtstart:  time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				Byminute: []int{-1}, // Invalid minute.
			},
			expectError: true,
		},
		{
			name: "invalid_bysecond_negative",
			opt: ROption{
				Freq:     MINUTELY,
				Count:    1,
				Dtstart:  time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				Bysecond: []int{-1}, // Invalid second.
			},
			expectError: true,
		},
		{
			name: "until_before_dtstart",
			opt: ROption{
				Freq:    DAILY,
				Dtstart: time.Date(2020, 1, 10, 10, 0, 0, 0, time.UTC),
				Until:   time.Date(2020, 1, 5, 10, 0, 0, 0, time.UTC), // Until is before Dtstart.
			},
			expectError: false,
			expected:    []time.Time{}, // Should return empty results.
		},
		{
			name: "zero_interval",
			opt: ROption{
				Freq:     DAILY,
				Interval: 0, // Will be auto-corrected to 1.
				Count:    1,
				Dtstart:  time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
			},
			expectError: false,
			expected:    []time.Time{time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC)},
		},
		{
			name: "negative_interval",
			opt: ROption{
				Freq:     DAILY,
				Interval: -1, // Invalid interval.
				Count:    1,
				Dtstart:  time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := newRecurrence(tt.opt)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			result := r.All()
			if !timesEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestIteratorPerformanceAndMemory tests performance and memory usage.
func TestIteratorPerformanceAndMemory(t *testing.T) {
	tests := []struct {
		name    string
		opt     ROption
		count   int
		maxTime time.Duration
	}{
		{
			name: "large_count_daily",
			opt: ROption{
				Freq:    DAILY,
				Dtstart: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				Count:   10000,
			},
			count:   10000,
			maxTime: 100 * time.Millisecond,
		},
		{
			name: "large_count_hourly",
			opt: ROption{
				Freq:    HOURLY,
				Dtstart: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				Count:   5000,
			},
			count:   5000,
			maxTime: 200 * time.Millisecond,
		},
		{
			name: "complex_byrules_performance",
			opt: ROption{
				Freq:      MONTHLY,
				Dtstart:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				Bymonth:   []int{1, 3, 5, 7, 9, 11},
				Byweekday: []Weekday{MO, WE, FR},
				Count:     1000,
			},
			count:   1000,
			maxTime: 500 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()

			// Measure performance.
			rrule, err := newRecurrence(tt.opt)
			if err != nil {
				t.Fatalf("Failed to create RRule: %v", err)
			}

			results := rrule.All()
			elapsed := time.Since(start)

			// Verify result count.
			if len(results) != tt.count {
				t.Errorf("Expected %d results, got %d", tt.count, len(results))
			}

			// Verify performance.
			if elapsed > tt.maxTime {
				t.Errorf("Performance test failed: took %v, expected < %v", elapsed, tt.maxTime)
			}

			t.Logf("Generated %d results in %v", len(results), elapsed)
		})
	}
}

// TestIteratorMemoryReuse tests memory reuse.
func TestIteratorMemoryReuse(t *testing.T) {
	opt := ROption{
		Freq:    DAILY,
		Dtstart: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		Count:   100,
	}

	rrule, err := newRecurrence(opt)
	if err != nil {
		t.Fatalf("Failed to create RRule: %v", err)
	}

	// Call All() multiple times to verify memory reuse.
	var results [][]time.Time
	for i := 0; i < 5; i++ {
		result := rrule.All()
		results = append(results, result)

		// Verify result consistency.
		if len(result) != 100 {
			t.Errorf("Iteration %d: expected 100 results, got %d", i, len(result))
		}
	}

	// Verify all results match.
	for i := 1; i < len(results); i++ {
		if !timesEqual(results[0], results[i]) {
			t.Errorf("Results differ between iterations 0 and %d", i)
		}
	}
}

// TestIteratorConcurrency tests concurrency safety.
func TestIteratorConcurrency(t *testing.T) {
	opt := ROption{
		Freq:    DAILY,
		Dtstart: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		Count:   1000,
	}

	// Call All() concurrently.
	const numGoroutines = 10
	results := make([][]time.Time, numGoroutines)
	done := make(chan int, numGoroutines)
	errs := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			rrule, err := newRecurrence(opt)
			if err != nil {
				errs <- err
				done <- index
				return
			}
			results[index] = rrule.All()
			done <- index
		}(i)
	}

	// Wait for all goroutines to finish.
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatalf("Failed to create RRule: %v", err)
		}
	}

	// Verify all results match.
	for i := 1; i < numGoroutines; i++ {
		if !timesEqual(results[0], results[i]) {
			t.Errorf("Concurrent results differ between goroutine 0 and %d", i)
		}
	}
}
