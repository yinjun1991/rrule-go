package rrule

import (
	"testing"
	"time"
)

// TestAllDayDaily 测试全天事件的每日循环
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

// TestAllDayWeekly 测试全天事件的每周循环
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

// TestAllDayMonthly 测试全天事件的每月循环
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

// TestAllDayYearly 测试全天事件的每年循环
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

// TestAllDayWithUntil 测试全天事件的Until边界处理
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

// TestAllDayWithUntilMidnight 测试全天事件Until为午夜的边界情况
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

// TestAllDayWithCount 测试全天事件的Count计数逻辑
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

// TestAllDayDTStartProcessing 测试全天事件DTStart的时间处理
func TestAllDayDTStartProcessing(t *testing.T) {
	// 测试不同时间的DTStart都应该被规范化为00:00:00
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

// TestAllDaySetWithRDate 测试全天事件的Set RDate功能
func TestAllDaySetWithRDate(t *testing.T) {
	set := Set{}

	// 添加基础规则
	r, _ := NewRRule(ROption{
		Freq:    WEEKLY,
		Count:   2,
		AllDay:  true,
		Dtstart: time.Date(2023, 1, 1, 14, 30, 0, 0, time.UTC),
	})
	set.RRule(r)

	// 添加额外日期
	set.RDate(time.Date(2023, 1, 20, 16, 45, 0, 0, time.UTC))
	set.RDate(time.Date(2023, 1, 25, 9, 15, 30, 0, time.UTC))

	want := []time.Time{
		time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 8, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 20, 16, 45, 0, 0, time.UTC), // RDate保持原时间
		time.Date(2023, 1, 25, 9, 15, 30, 0, time.UTC), // RDate保持原时间
	}
	value := set.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

// TestAllDaySetWithExDate 测试全天事件的Set ExDate功能
func TestAllDaySetWithExDate(t *testing.T) {
	set := Set{}

	// 添加基础规则
	r, _ := NewRRule(ROption{
		Freq:    DAILY,
		Count:   5,
		AllDay:  true,
		Dtstart: time.Date(2023, 1, 1, 11, 20, 0, 0, time.UTC),
	})
	set.RRule(r)

	// 排除特定日期（需要使用午夜时间来匹配全天事件）
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

// TestAllDaySetComplex 测试全天事件的复杂Set组合
func TestAllDaySetComplex(t *testing.T) {
	set := Set{}

	// 添加基础规则
	r, _ := NewRRule(ROption{
		Freq:    WEEKLY,
		Count:   3,
		AllDay:  true,
		Dtstart: time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
	})
	set.RRule(r)

	// 添加额外日期
	set.RDate(time.Date(2023, 1, 10, 12, 45, 0, 0, time.UTC))

	// 排除日期（需要使用午夜时间来匹配全天事件）
	set.ExDate(time.Date(2023, 1, 8, 0, 0, 0, 0, time.UTC))

	want := []time.Time{
		time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 10, 12, 45, 0, 0, time.UTC), // RDate保持原时间
		time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
	}
	value := set.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

// TestAllDayTimezoneHandling 测试全天事件的时区处理
// 根据RFC 5545，全天事件应该使用浮动时间，在所有时区都表示同一天的00:00:00
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

			// 根据RFC 5545浮动时间规范，全天事件应该转换为浮动时间（无时区绑定）
			// 在Go中，我们用UTC表示浮动时间，因为它不依赖本地时区
			// 这样确保在任何时区的用户看到的都是同一天的00:00:00
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

// TestAllDayLeapYear 测试全天事件的闰年处理
func TestAllDayLeapYear(t *testing.T) {
	r, _ := NewRRule(ROption{
		Freq:    YEARLY,
		Count:   4,
		AllDay:  true,
		Dtstart: time.Date(2020, 2, 29, 12, 0, 0, 0, time.UTC), // 闰年2月29日
	})
	want := []time.Time{
		time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC), // 下一个闰年
		time.Date(2028, 2, 29, 0, 0, 0, 0, time.UTC),
		time.Date(2032, 2, 29, 0, 0, 0, 0, time.UTC),
	}
	value := r.All()
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

// TestAllDayYearBoundary 测试全天事件的跨年边界
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

// TestAllDayMonthBoundary 测试全天事件的跨月边界
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

// TestAllDayIterator 测试全天事件的Iterator方法
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

// TestAllDayWithByWeekDay 测试全天事件结合ByWeekDay规则
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

// TestAllDayWithByMonthDay 测试全天事件结合ByMonthDay规则
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

// TestAllDayConsistencyWithNonAllDay 测试全天事件与非全天事件的一致性
func TestAllDayConsistencyWithNonAllDay(t *testing.T) {
	// 全天事件
	allDayRule, _ := NewRRule(ROption{
		Freq:    DAILY,
		Count:   3,
		AllDay:  true,
		Dtstart: time.Date(2023, 1, 1, 14, 30, 0, 0, time.UTC),
	})

	// 非全天事件（相同日期，午夜时间）
	nonAllDayRule, _ := NewRRule(ROption{
		Freq:    DAILY,
		Count:   3,
		AllDay:  false,
		Dtstart: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	})

	allDayResults := allDayRule.All()
	nonAllDayResults := nonAllDayRule.All()

	// 全天事件的结果应该与午夜非全天事件的结果相同
	if !timesEqual(allDayResults, nonAllDayResults) {
		t.Errorf("AllDay results %v should match non-AllDay midnight results %v",
			allDayResults, nonAllDayResults)
	}
}
