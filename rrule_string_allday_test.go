package rrule

import (
	"strings"
	"testing"
	"time"
)

// TestRRuleStringAllDayUntil 测试RRuleString()方法对全天事件UNTIL参数的处理
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
			// 全天事件，带UNTIL参数
			option := ROption{
				Freq:    DAILY,
				AllDay:  true,
				Dtstart: time.Date(2023, 1, 1, 14, 30, 0, 0, tc.tz),
				Until:   time.Date(2023, 1, 5, 16, 45, 30, 0, tc.tz),
			}

			output := option.RRuleString()
			t.Logf("Timezone %s RRuleString output: %s", tc.name, output)

			// 验证UNTIL参数是否正确处理
			if !strings.Contains(output, "UNTIL=") {
				t.Errorf("Expected UNTIL parameter in output for timezone %s, got: %s", tc.name, output)
			}

			// 验证UNTIL使用DATE格式（不包含时间部分）
			if !strings.Contains(output, "UNTIL=20230105") {
				t.Errorf("Expected UNTIL=20230105 in output for timezone %s, got: %s", tc.name, output)
			}

			// 验证UNTIL不包含时间部分（T后面跟数字）
			if strings.Contains(output, "UNTIL=20230105T") {
				t.Errorf("UNTIL should use DATE format (no time part) for all-day events in timezone %s, got: %s", tc.name, output)
			}
		})
	}
}

// TestRRuleStringAllDayConsistency 测试全天事件RRuleString的一致性
func TestRRuleStringAllDayConsistency(t *testing.T) {
	// 创建相同日期但不同时区的全天事件
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

	// 验证所有输出都相同（因为RRuleString不包含DTSTART，只包含RRULE部分）
	expected := "FREQ=WEEKLY;COUNT=3"
	for i, output := range outputs {
		if output != expected {
			t.Errorf("Option %d: expected %s, got %s", i+1, expected, output)
		}
	}
}

// TestRRuleStringAllDayWithUntilTimezone 测试不同时区的UNTIL参数处理
func TestRRuleStringAllDayWithUntilTimezone(t *testing.T) {
	// 创建相同日期但不同时区的全天事件，带UNTIL参数
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

	// 验证UNTIL参数使用DATE格式
	for i, output := range outputs {
		if !strings.Contains(output, "UNTIL=20230103") {
			t.Errorf("Option %d: expected UNTIL=20230103 in output, got: %s", i+1, output)
		}
		// 验证UNTIL不包含时间部分（T后面跟数字）
		if strings.Contains(output, "UNTIL=20230103T") {
			t.Errorf("Option %d: UNTIL should use DATE format (no time part) for all-day events, got: %s", i+1, output)
		}
	}
}
