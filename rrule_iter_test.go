package rrule

import (
	"testing"
	"time"
)

// TestIteratorFrequencyAdvanced 测试高级频率场景：间隔、跨年、跨月等
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

// TestIteratorByRulesCombinations 测试BY*规则组合
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
				Bymonth: []int{6}, // 只在6月
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
				Bymonth: []int{3, 9}, // 3月和9月
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
				Byweekday: []Weekday{FR},                                // 只在周五
			},
			expected: []time.Time{
				time.Date(2020, 1, 10, 10, 0, 0, 0, time.UTC), // 第一个周五
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
				Byweekday: []Weekday{MO, WE, FR},                        // 周一、周三、周五
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
				Bymonthday: []int{15}, // 每月15日
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
				Bymonthday: []int{-1}, // 每月最后一天
			},
			expected: []time.Time{
				time.Date(2020, 1, 31, 10, 0, 0, 0, time.UTC),
				time.Date(2020, 2, 29, 10, 0, 0, 0, time.UTC), // 闰年2月
				time.Date(2020, 3, 31, 10, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "BYHOUR_multiple",
			opt: ROption{
				Freq:    DAILY,
				Count:   4,
				Dtstart: time.Date(2020, 1, 1, 8, 0, 0, 0, time.UTC),
				Byhour:  []int{9, 15}, // 9点和15点
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
				Byminute: []int{15, 45}, // 15分和45分
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
				Bysecond: []int{10, 50}, // 10秒和50秒
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

// TestIteratorLeapYearAndBoundaries 测试闰年和边界情况
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
				Dtstart: time.Date(2020, 2, 29, 12, 0, 0, 0, time.UTC), // 闰年2月29日
			},
			expected: []time.Time{
				time.Date(2020, 2, 29, 12, 0, 0, 0, time.UTC),
				time.Date(2024, 2, 29, 12, 0, 0, 0, time.UTC), // 下一个闰年
				time.Date(2028, 2, 29, 12, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "month_boundary_31_to_30_days",
			opt: ROption{
				Freq:    MONTHLY,
				Count:   3,
				Dtstart: time.Date(2020, 1, 31, 15, 0, 0, 0, time.UTC), // 1月31日
			},
			expected: []time.Time{
				time.Date(2020, 1, 31, 15, 0, 0, 0, time.UTC),
				time.Date(2020, 3, 31, 15, 0, 0, 0, time.UTC), // 跳过2月（只有28/29天）
				time.Date(2020, 5, 31, 15, 0, 0, 0, time.UTC), // 跳过4月（只有30天）
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

// TestIteratorAllDayVsTimedEvents 测试全天事件与定时事件的差异
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
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), // 全天事件时间被重置为00:00:00
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
				time.Date(2020, 1, 1, 14, 30, 45, 0, time.UTC), // 保持原始时间
				time.Date(2020, 1, 2, 14, 30, 45, 0, time.UTC),
				time.Date(2020, 1, 3, 14, 30, 45, 0, time.UTC),
			},
		},
		{
			name: "allday_with_until",
			opt: ROption{
				Freq:    DAILY,
				Dtstart: time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				Until:   time.Date(2020, 1, 3, 15, 30, 0, 0, time.UTC), // Until时间也会被重置
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
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), // 全天事件小时固定为0
				time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "allday_with_byhour",
			opt: ROption{
				Freq:    DAILY,
				Count:   2,
				Dtstart: time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				Byhour:  []int{9, 15}, // 全天事件应该忽略BYHOUR
				AllDay:  true,
			},
			expected: []time.Time{
				time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 15, 0, 0, 0, time.UTC),
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

// TestIteratorTimezoneHandling 测试多时区处理
func TestIteratorTimezoneHandling(t *testing.T) {
	// 创建不同时区
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
				time.Date(2020, 1, 2, 0, 0, 0, 0, east), // 跨日边界
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
				AllDay:  true, // 全天事件应该与时区无关
			},
			expected: []time.Time{
				time.Date(2020, 1, 1, 0, 0, 0, 0, utc), // 转换为UTC的00:00
				time.Date(2020, 1, 2, 0, 0, 0, 0, utc),
			},
		},
		{
			name: "mixed_timezone_until",
			opt: ROption{
				Freq:    DAILY,
				Dtstart: time.Date(2020, 1, 1, 10, 0, 0, 0, east),
				Until:   time.Date(2020, 1, 2, 23, 0, 0, 0, utc), // UTC时间，对应东八区次日07:00
			},
			expected: []time.Time{
				time.Date(2020, 1, 1, 10, 0, 0, 0, east),
				time.Date(2020, 1, 2, 10, 0, 0, 0, east),
				// 2020-01-03 10:00 东八区 = 2020-01-03 02:00 UTC，晚于Until时间
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

// TestIteratorErrorAndNegativeCases 测试错误和负面情况
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
				Bymonth: []int{0}, // 无效月份
			},
			expectError: true,
		},
		{
			name: "invalid_bymonth_negative",
			opt: ROption{
				Freq:    YEARLY,
				Count:   1,
				Dtstart: time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				Bymonth: []int{-1}, // 无效月份
			},
			expectError: true,
		},
		{
			name: "invalid_byhour_negative",
			opt: ROption{
				Freq:    DAILY,
				Count:   1,
				Dtstart: time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				Byhour:  []int{-1}, // 无效小时
			},
			expectError: true,
		},
		{
			name: "invalid_byminute_negative",
			opt: ROption{
				Freq:     HOURLY,
				Count:    1,
				Dtstart:  time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				Byminute: []int{-1}, // 无效分钟
			},
			expectError: true,
		},
		{
			name: "invalid_bysecond_negative",
			opt: ROption{
				Freq:     MINUTELY,
				Count:    1,
				Dtstart:  time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				Bysecond: []int{-1}, // 无效秒数
			},
			expectError: true,
		},
		{
			name: "until_before_dtstart",
			opt: ROption{
				Freq:    DAILY,
				Dtstart: time.Date(2020, 1, 10, 10, 0, 0, 0, time.UTC),
				Until:   time.Date(2020, 1, 5, 10, 0, 0, 0, time.UTC), // Until在Dtstart之前
			},
			expectError: false,
			expected:    []time.Time{}, // 应该返回空结果
		},
		{
			name: "zero_interval",
			opt: ROption{
				Freq:     DAILY,
				Interval: 0, // 会被自动修正为1
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
				Interval: -1, // 无效间隔
				Count:    1,
				Dtstart:  time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewRRule(tt.opt)
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

// TestIteratorPerformanceAndMemory 测试性能和内存使用
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

			// 测试性能
			rrule, err := NewRRule(tt.opt)
			if err != nil {
				t.Fatalf("Failed to create RRule: %v", err)
			}

			results := rrule.All()
			elapsed := time.Since(start)

			// 验证结果数量
			if len(results) != tt.count {
				t.Errorf("Expected %d results, got %d", tt.count, len(results))
			}

			// 验证性能
			if elapsed > tt.maxTime {
				t.Errorf("Performance test failed: took %v, expected < %v", elapsed, tt.maxTime)
			}

			t.Logf("Generated %d results in %v", len(results), elapsed)
		})
	}
}

// TestIteratorMemoryReuse 测试内存重用机制
func TestIteratorMemoryReuse(t *testing.T) {
	opt := ROption{
		Freq:    DAILY,
		Dtstart: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		Count:   100,
	}

	rrule, err := NewRRule(opt)
	if err != nil {
		t.Fatalf("Failed to create RRule: %v", err)
	}

	// 多次调用All()方法，验证内存重用
	var results [][]time.Time
	for i := 0; i < 5; i++ {
		result := rrule.All()
		results = append(results, result)

		// 验证结果一致性
		if len(result) != 100 {
			t.Errorf("Iteration %d: expected 100 results, got %d", i, len(result))
		}
	}

	// 验证所有结果相同
	for i := 1; i < len(results); i++ {
		if !timesEqual(results[0], results[i]) {
			t.Errorf("Results differ between iterations 0 and %d", i)
		}
	}
}

// TestIteratorConcurrency 测试并发安全性
func TestIteratorConcurrency(t *testing.T) {
	opt := ROption{
		Freq:    DAILY,
		Dtstart: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		Count:   1000,
	}

	rrule, err := NewRRule(opt)
	if err != nil {
		t.Fatalf("Failed to create RRule: %v", err)
	}

	// 并发调用All()方法
	const numGoroutines = 10
	results := make([][]time.Time, numGoroutines)
	done := make(chan int, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			results[index] = rrule.All()
			done <- index
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// 验证所有结果相同
	for i := 1; i < numGoroutines; i++ {
		if !timesEqual(results[0], results[i]) {
			t.Errorf("Concurrent results differ between goroutine 0 and %d", i)
		}
	}
}
