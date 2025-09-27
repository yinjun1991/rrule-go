package rrule

import (
	"testing"
	"time"
)

// TestIteratorBasicFrequencies 测试基本频率的迭代器功能
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
			r, err := NewRRule(tt.opt)
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

// TestIteratorAllDayEvents 测试全天事件的迭代器
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
			r, err := NewRRule(tt.opt)
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

// TestIteratorTimezones 测试多时区处理
func TestIteratorTimezones(t *testing.T) {
	// 加载时区
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
			r, err := NewRRule(tt.opt)
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

// TestIteratorByRules 测试 BY* 规则组合
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
			r, err := NewRRule(tt.opt)
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

// TestIteratorBoundaryConditions 测试边界条件
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
			r, err := NewRRule(tt.opt)
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

// TestIteratorInterval 测试间隔设置
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
			r, err := NewRRule(tt.opt)
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

// TestIteratorUntil 测试 UNTIL 限制
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
			r, err := NewRRule(tt.opt)
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

// TestIteratorComplexByRules 测试复杂的 BY* 规则组合
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
			r, err := NewRRule(tt.opt)
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

// TestIteratorEdgeCases 测试边缘情况
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
			name: "Zero_count",
			opt: ROption{
				Freq:    DAILY,
				Count:   0,
				Dtstart: time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
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
			r, err := NewRRule(tt.opt)
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

// TestIteratorPerformance 测试迭代器性能和内存使用
func TestIteratorPerformance(t *testing.T) {
	// 测试大量迭代的性能
	r, err := NewRRule(ROption{
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

	// 验证结果的连续性
	for i := 1; i < len(result); i++ {
		expected := result[i-1].AddDate(0, 0, 1)
		if !result[i].Equal(expected) {
			t.Errorf("Result %d: expected %v, got %v", i, expected, result[i])
			break
		}
	}
}

// TestIteratorNext 测试迭代器的 next() 方法
func TestIteratorNext(t *testing.T) {
	r, err := NewRRule(ROption{
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

	// 应该没有更多值
	_, ok := iter()
	if ok {
		t.Error("Iterator should have ended")
	}
}

// TestIteratorMaxYear 测试最大年份限制
func TestIteratorMaxYear(t *testing.T) {
	r, err := NewRRule(ROption{
		Freq:    YEARLY,
		Count:   10,
		Dtstart: time.Date(9995, 1, 1, 10, 0, 0, 0, time.UTC), // Close to MAXYEAR
	})
	if err != nil {
		t.Fatalf("Failed to create RRule: %v", err)
	}

	result := r.All()
	// 应该在达到 MAXYEAR 时停止
	for _, dt := range result {
		if dt.Year() > MAXYEAR {
			t.Errorf("Result year %d exceeds MAXYEAR %d", dt.Year(), MAXYEAR)
		}
	}
}