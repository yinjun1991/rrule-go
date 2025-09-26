package rrule

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// TestROptionFrequencies 测试所有频率类型的基本功能
func TestROptionFrequencies(t *testing.T) {
	testCases := []struct {
		name string
		freq Frequency
		expected string
	}{
		{"YEARLY", YEARLY, "FREQ=YEARLY"},
		{"MONTHLY", MONTHLY, "FREQ=MONTHLY"},
		{"WEEKLY", WEEKLY, "FREQ=WEEKLY"},
		{"DAILY", DAILY, "FREQ=DAILY"},
		{"HOURLY", HOURLY, "FREQ=HOURLY"},
		{"MINUTELY", MINUTELY, "FREQ=MINUTELY"},
		{"SECONDLY", SECONDLY, "FREQ=SECONDLY"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			option := ROption{
				Freq: tc.freq,
			}
			output := option.RRuleString()
			if output != tc.expected {
				t.Errorf("Expected %s, got: %s", tc.expected, output)
			}
		})
	}
}

// TestROptionBasicParameters 测试基本参数
func TestROptionBasicParameters(t *testing.T) {
	testCases := []struct {
		name string
		option ROption
		expected string
	}{
		{
			"With Interval",
			ROption{Freq: DAILY, Interval: 2},
			"FREQ=DAILY;INTERVAL=2",
		},
		{
			"With Count",
			ROption{Freq: WEEKLY, Count: 5},
			"FREQ=WEEKLY;COUNT=5",
		},
		{
			"With WKST",
			ROption{Freq: WEEKLY, Wkst: SU},
			"FREQ=WEEKLY;WKST=SU",
		},
		{
			"Complex combination",
			ROption{Freq: MONTHLY, Interval: 3, Count: 10, Wkst: TU},
			"FREQ=MONTHLY;INTERVAL=3;WKST=TU;COUNT=10",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output := tc.option.RRuleString()
			if output != tc.expected {
				t.Errorf("Expected %s, got: %s", tc.expected, output)
			}
		})
	}
}

// TestROptionByRules 测试 BY* 规则
func TestROptionByRules(t *testing.T) {
	testCases := []struct {
		name string
		option ROption
		expected string
	}{
		{
			"BYMONTH",
			ROption{Freq: YEARLY, Bymonth: []int{1, 6, 12}},
			"FREQ=YEARLY;BYMONTH=1,6,12",
		},
		{
			"BYMONTHDAY",
			ROption{Freq: MONTHLY, Bymonthday: []int{1, 15, -1}},
			"FREQ=MONTHLY;BYMONTHDAY=1,15,-1",
		},
		{
			"BYDAY",
			ROption{Freq: WEEKLY, Byweekday: []Weekday{MO, WE, FR}},
			"FREQ=WEEKLY;BYDAY=MO,WE,FR",
		},
		{
			"BYHOUR",
			ROption{Freq: DAILY, Byhour: []int{9, 12, 18}},
			"FREQ=DAILY;BYHOUR=9,12,18",
		},
		{
			"BYMINUTE",
			ROption{Freq: HOURLY, Byminute: []int{0, 30}},
			"FREQ=HOURLY;BYMINUTE=0,30",
		},
		{
			"BYSECOND",
			ROption{Freq: MINUTELY, Bysecond: []int{0, 15, 30, 45}},
			"FREQ=MINUTELY;BYSECOND=0,15,30,45",
		},
		{
			"Multiple BY rules",
			ROption{
				Freq: MONTHLY,
				Bymonthday: []int{1, 15},
				Byhour: []int{9, 17},
				Byminute: []int{0},
			},
			"FREQ=MONTHLY;BYMONTHDAY=1,15;BYHOUR=9,17;BYMINUTE=0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output := tc.option.RRuleString()
			if output != tc.expected {
				t.Errorf("Expected %s, got: %s", tc.expected, output)
			}
		})
	}
}

// TestROptionNonAllDayWithTimezone 测试非全天事件的时区处理
func TestROptionNonAllDayWithTimezone(t *testing.T) {
	testCases := []struct {
		name string
		tz *time.Location
		dtstart time.Time
		until time.Time
	}{
		{
			"UTC",
			time.UTC,
			time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
			time.Date(2023, 1, 31, 10, 0, 0, 0, time.UTC),
		},
		{
			"EST",
			time.FixedZone("EST", -5*3600),
			time.Date(2023, 6, 1, 14, 30, 0, 0, time.FixedZone("EST", -5*3600)),
			time.Date(2023, 6, 30, 14, 30, 0, 0, time.FixedZone("EST", -5*3600)),
		},
		{
			"JST",
			time.FixedZone("JST", 9*3600),
			time.Date(2023, 12, 1, 9, 0, 0, 0, time.FixedZone("JST", 9*3600)),
			time.Date(2023, 12, 25, 9, 0, 0, 0, time.FixedZone("JST", 9*3600)),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			option := ROption{
				Freq: DAILY,
				Dtstart: tc.dtstart,
				Until: tc.until,
				AllDay: false,
			}

			output := option.String()
			t.Logf("Non-AllDay %s output: %s", tc.name, output)

			// 验证 DTSTART 包含时区信息
			if !strings.Contains(output, "DTSTART") {
				t.Errorf("Expected DTSTART in output for %s", tc.name)
			}

			// 验证 UNTIL 包含时间部分
			if !strings.Contains(output, "UNTIL=") {
				t.Errorf("Expected UNTIL in output for %s", tc.name)
			}

			// 非全天事件的 UNTIL 应该包含时间部分
			lines := strings.Split(output, "\n")
			for _, line := range lines {
				if strings.Contains(line, "UNTIL=") {
					parts := strings.Split(line, "UNTIL=")
					if len(parts) > 1 {
						untilValue := strings.Split(parts[1], ";")[0]
						if !strings.Contains(untilValue, "T") {
							t.Errorf("Non-AllDay UNTIL should contain time part, got: %s", untilValue)
						}
					}
				}
			}
		})
	}
}

// TestROptionStringParsing 测试字符串解析功能
func TestROptionStringParsing(t *testing.T) {
	testCases := []struct {
		name string
		input string
		expectedError bool
		validateFunc func(*testing.T, *ROption)
	}{
		{
			"Simple RRULE",
			"RRULE:FREQ=DAILY;COUNT=5",
			false,
			func(t *testing.T, opt *ROption) {
				if opt.Freq != DAILY {
					t.Errorf("Expected DAILY, got %v", opt.Freq)
				}
				if opt.Count != 5 {
					t.Errorf("Expected Count=5, got %d", opt.Count)
				}
			},
		},
		{
			"RRULE with DTSTART",
			"DTSTART:20230101T100000Z\nRRULE:FREQ=WEEKLY;INTERVAL=2",
			false,
			func(t *testing.T, opt *ROption) {
				if opt.Freq != WEEKLY {
					t.Errorf("Expected WEEKLY, got %v", opt.Freq)
				}
				if opt.Interval != 2 {
					t.Errorf("Expected Interval=2, got %d", opt.Interval)
				}
				expected := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
				if !opt.Dtstart.Equal(expected) {
					t.Errorf("Expected Dtstart=%v, got %v", expected, opt.Dtstart)
				}
			},
		},
		{
			"Complex RRULE",
			"RRULE:FREQ=MONTHLY;BYMONTHDAY=1,15;BYHOUR=9,17;COUNT=10",
			false,
			func(t *testing.T, opt *ROption) {
				if opt.Freq != MONTHLY {
					t.Errorf("Expected MONTHLY, got %v", opt.Freq)
				}
				if len(opt.Bymonthday) != 2 || opt.Bymonthday[0] != 1 || opt.Bymonthday[1] != 15 {
					t.Errorf("Expected Bymonthday=[1,15], got %v", opt.Bymonthday)
				}
				if len(opt.Byhour) != 2 || opt.Byhour[0] != 9 || opt.Byhour[1] != 17 {
					t.Errorf("Expected Byhour=[9,17], got %v", opt.Byhour)
				}
			},
		},
		{
			"Missing FREQ - should error",
			"RRULE:COUNT=5;INTERVAL=2",
			true,
			nil,
		},
		{
			"Invalid format - should error",
			"INVALID_RRULE_FORMAT",
			true,
			nil,
		},
		{
			"Empty value - should error",
			"RRULE:FREQ=;COUNT=5",
			true,
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opt, err := StrToROption(tc.input)

			if tc.expectedError {
				if err == nil {
					t.Errorf("Expected error for input: %s", tc.input)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for input %s: %v", tc.input, err)
				return
			}

			if tc.validateFunc != nil {
				tc.validateFunc(t, opt)
			}
		})
	}
}

// TestROptionStringParsingInLocation 测试带时区的字符串解析
func TestROptionStringParsingInLocation(t *testing.T) {
	estLoc := time.FixedZone("EST", -5*3600)
	jstLoc := time.FixedZone("JST", 9*3600)

	testCases := []struct {
		name string
		input string
		location *time.Location
		validateFunc func(*testing.T, *ROption)
	}{
		{
			"Local time in EST",
			"DTSTART:20230601T140000\nRRULE:FREQ=DAILY;COUNT=3",
			estLoc,
			func(t *testing.T, opt *ROption) {
				expected := time.Date(2023, 6, 1, 14, 0, 0, 0, estLoc)
				if !opt.Dtstart.Equal(expected) {
					t.Errorf("Expected Dtstart in EST: %v, got %v", expected, opt.Dtstart)
				}
			},
		},
		{
			"UNTIL in JST",
			"RRULE:FREQ=WEEKLY;UNTIL=20231225T090000",
			jstLoc,
			func(t *testing.T, opt *ROption) {
				expected := time.Date(2023, 12, 25, 9, 0, 0, 0, jstLoc)
				if !opt.Until.Equal(expected) {
					t.Errorf("Expected Until in JST: %v, got %v", expected, opt.Until)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opt, err := StrToROptionInLocation(tc.input, tc.location)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if tc.validateFunc != nil {
				tc.validateFunc(t, opt)
			}
		})
	}
}

// TestROptionEdgeCases 测试边界情况和错误处理
func TestROptionEdgeCases(t *testing.T) {
	t.Run("Zero values", func(t *testing.T) {
		option := ROption{Freq: DAILY}
		output := option.RRuleString()
		expected := "FREQ=DAILY"
		if output != expected {
			t.Errorf("Expected %s, got: %s", expected, output)
		}
	})

	t.Run("Empty DTSTART", func(t *testing.T) {
		option := ROption{Freq: WEEKLY, Count: 3}
		output := option.String()
		// 当 DTSTART 为零值时，String() 应该只返回 RRuleString()
		expected := "FREQ=WEEKLY;COUNT=3"
		if output != expected {
			t.Errorf("Expected %s, got: %s", expected, output)
		}
	})

	t.Run("Negative values in BY rules", func(t *testing.T) {
		option := ROption{
			Freq: MONTHLY,
			Bymonthday: []int{-1, -7, 15}, // 负值表示从月末倒数
		}
		output := option.RRuleString()
		expected := "FREQ=MONTHLY;BYMONTHDAY=-1,-7,15"
		if output != expected {
			t.Errorf("Expected %s, got: %s", expected, output)
		}
	})

	t.Run("Large numbers", func(t *testing.T) {
		option := ROption{
			Freq: YEARLY,
			Interval: 999,
			Count: 9999,
			Byyearday: []int{1, 100, 365, -1},
		}
		output := option.RRuleString()
		if !strings.Contains(output, "INTERVAL=999") {
			t.Errorf("Expected INTERVAL=999 in output: %s", output)
		}
		if !strings.Contains(output, "COUNT=9999") {
			t.Errorf("Expected COUNT=9999 in output: %s", output)
		}
	})
}

// TestROptionRoundTrip 测试往返转换的一致性
func TestROptionRoundTrip(t *testing.T) {
	testCases := []struct {
		name string
		option ROption
	}{
		{
			"Simple daily",
			ROption{
				Freq: DAILY,
				Count: 5,
				Dtstart: time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
			},
		},
		{
			"Complex monthly",
			ROption{
				Freq: MONTHLY,
				Interval: 2,
				Bymonthday: []int{1, 15},
				Byhour: []int{9, 17},
				Wkst: SU,
				Dtstart: time.Date(2023, 6, 1, 9, 0, 0, 0, time.UTC),
				Until: time.Date(2024, 6, 1, 9, 0, 0, 0, time.UTC),
			},
		},
		{
			"All-day event",
			ROption{
				Freq: WEEKLY,
				Byweekday: []Weekday{MO, WE, FR},
				Count: 10,
				AllDay: true,
				Dtstart: time.Date(2023, 3, 1, 14, 30, 0, 0, time.FixedZone("EST", -5*3600)),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 第一步：ROption -> String
			originalStr := tc.option.String()
			t.Logf("Original string: %s", originalStr)

			// 第二步：String -> ROption
			parsedOption, err := StrToROption(originalStr)
			if err != nil {
				t.Errorf("Failed to parse string: %v", err)
				return
			}

			// 第三步：验证关键字段的一致性
			if parsedOption.Freq != tc.option.Freq {
				t.Errorf("Freq mismatch: expected %v, got %v", tc.option.Freq, parsedOption.Freq)
			}

			if parsedOption.Count != tc.option.Count {
				t.Errorf("Count mismatch: expected %d, got %d", tc.option.Count, parsedOption.Count)
			}

			if parsedOption.Interval != tc.option.Interval {
				t.Errorf("Interval mismatch: expected %d, got %d", tc.option.Interval, parsedOption.Interval)
			}

			// 对于全天事件，只比较日期部分
			if tc.option.AllDay {
				origY, origM, origD := tc.option.Dtstart.Date()
				parsedY, parsedM, parsedD := parsedOption.Dtstart.Date()
				if origY != parsedY || origM != parsedM || origD != parsedD {
					t.Errorf("AllDay Dtstart date mismatch: expected %v-%v-%v, got %v-%v-%v",
						origY, origM, origD, parsedY, parsedM, parsedD)
				}
			} else if !tc.option.Dtstart.IsZero() {
				// 对于非全天事件，比较完整的时间（允许一定的时区转换误差）
				if !tc.option.Dtstart.Equal(parsedOption.Dtstart) {
					t.Errorf("Dtstart mismatch: expected %v, got %v", tc.option.Dtstart, parsedOption.Dtstart)
				}
			}

			// 第四步：再次序列化，验证字符串的一致性（排除全天事件的时区差异）
			if !tc.option.AllDay {
				reparsedStr := parsedOption.String()
				t.Logf("Reparsed string: %s", reparsedStr)

				// 对于非全天事件，RRULE 部分应该完全一致
				origRRule := tc.option.RRuleString()
				parsedRRule := parsedOption.RRuleString()
				if origRRule != parsedRRule {
					t.Errorf("RRule mismatch:\nOriginal:  %s\nReparsed:  %s", origRRule, parsedRRule)
				}
			}
		})
	}
}

// TestAllDayStringOutput 测试全天事件的String()和RRuleString()输出
func TestAllDayStringOutput(t *testing.T) {
	// 创建一个全天事件
	option := ROption{
		Freq:    DAILY,
		Count:   3,
		AllDay:  true,
		Dtstart: time.Date(2023, 1, 1, 14, 30, 0, 0, time.FixedZone("EST", -5*3600)),
	}

	// 测试String()方法输出
	t.Run("String() output", func(t *testing.T) {
		output := option.String()
		t.Logf("String() output: %s", output)

		// 验证DTSTART是否使用DATE格式（符合RFC 5545规范）
		if !strings.Contains(output, "DTSTART;VALUE=DATE:20230101") {
			t.Errorf("Expected DTSTART;VALUE=DATE:20230101 in output, got: %s", output)
		}
	})

	// 测试RRuleString()方法输出
	t.Run("RRuleString() output", func(t *testing.T) {
		output := option.RRuleString()
		t.Logf("RRuleString() output: %s", output)

		// RRuleString不包含DTSTART，只包含RRULE部分
		expected := "FREQ=DAILY;COUNT=3"
		if output != expected {
			t.Errorf("Expected %s, got: %s", expected, output)
		}
	})
}

// TestAllDayStringWithTimezone 测试不同时区的全天事件序列化
func TestAllDayStringWithTimezone(t *testing.T) {
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
			option := ROption{
				Freq:    WEEKLY,
				Count:   2,
				AllDay:  true,
				Dtstart: time.Date(2023, 6, 15, 16, 45, 30, 0, tc.tz),
			}

			output := option.String()
			t.Logf("Timezone %s output: %s", tc.name, output)

			// 所有时区的全天事件都应该序列化为相同的DATE格式（符合RFC 5545规范）
			expectedDTSTART := "DTSTART;VALUE=DATE:20230615"
			if !strings.Contains(output, expectedDTSTART) {
				t.Errorf("Expected %s in output for timezone %s, got: %s",
					expectedDTSTART, tc.name, output)
			}
		})
	}
}

// TestAllDayStringWithUntil 测试全天事件使用 UNTIL 的字符串序列化
func TestAllDayStringWithUntil(t *testing.T) {
	testCases := []struct {
		name     string
		dtstart  time.Time
		until    time.Time
		expected string
	}{
		{
			name:     "Same date UNTIL",
			dtstart:  time.Date(2023, 3, 1, 10, 0, 0, 0, time.UTC),
			until:    time.Date(2023, 3, 5, 23, 59, 59, 0, time.UTC),
			expected: "UNTIL=20230305",
		},
		{
			name:     "Different timezone UNTIL",
			dtstart:  time.Date(2023, 4, 10, 14, 30, 0, 0, time.FixedZone("EST", -5*3600)),
			until:    time.Date(2023, 4, 20, 8, 15, 0, 0, time.FixedZone("JST", 9*3600)),
			expected: "UNTIL=20230420",
		},
		{
			name:     "Cross month UNTIL",
			dtstart:  time.Date(2023, 2, 28, 12, 0, 0, 0, time.UTC),
			until:    time.Date(2023, 3, 15, 18, 45, 0, 0, time.UTC),
			expected: "UNTIL=20230315",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			option := ROption{
				Freq:    DAILY,
				AllDay:  true,
				Dtstart: tc.dtstart,
				Until:   tc.until,
			}

			output := option.String()
			t.Logf("UNTIL test %s output: %s", tc.name, output)

			// 验证 DTSTART 使用 DATE 格式
			year, month, day := tc.dtstart.Date()
			expectedDTSTART := fmt.Sprintf("DTSTART;VALUE=DATE:%04d%02d%02d", year, int(month), day)
			if !strings.Contains(output, expectedDTSTART) {
				t.Errorf("Expected %s in output, got: %s", expectedDTSTART, output)
			}

			// 验证 UNTIL 使用 DATE 格式（符合 RFC 5545 规范）
			if !strings.Contains(output, tc.expected) {
				t.Errorf("Expected %s in output, got: %s", tc.expected, output)
			}

			// 验证 UNTIL 不包含时间部分（T 后面跟数字）
			lines := strings.Split(output, "\n")
			for _, line := range lines {
				if strings.Contains(line, "UNTIL=") {
					// 提取 UNTIL 值
					parts := strings.Split(line, "UNTIL=")
					if len(parts) > 1 {
						untilValue := strings.Split(parts[1], ";")[0]  // 处理可能的后续参数
						untilValue = strings.Split(untilValue, " ")[0] // 处理可能的空格
						// 检查是否包含时间部分（T 后面跟数字）
						if strings.Contains(untilValue, "T") {
							t.Errorf("UNTIL should use DATE format (no time part) for all-day events, got: %s", untilValue)
						}
					}
				}
			}
		})
	}
}
