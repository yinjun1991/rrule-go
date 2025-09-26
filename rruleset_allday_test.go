package rrule

import (
	"strings"
	"testing"
	"time"
)

func TestSetAllDay_BasicFunctionality(t *testing.T) {
	set := &Set{}

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

func TestSetAllDay_DTStartNormalization(t *testing.T) {
	set := &Set{}
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

func TestSetAllDay_RDateNormalization(t *testing.T) {
	set := &Set{}
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

func TestSetAllDay_SetRDatesNormalization(t *testing.T) {
	set := &Set{}
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

func TestSetAllDay_ExDateNormalization(t *testing.T) {
	set := &Set{}
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

func TestSetAllDay_SetExDatesNormalization(t *testing.T) {
	set := &Set{}
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

func TestSetAllDay_ExistingTimesNormalization(t *testing.T) {
	set := &Set{}

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

func TestSetAllDay_RecurrenceSerialization(t *testing.T) {
	set := &Set{}
	set.SetAllDay(true)

	// Set up all-day event data
	dtstart := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)
	rdate := time.Date(2024, 7, 5, 0, 0, 0, 0, time.UTC)
	exdate := time.Date(2024, 7, 10, 0, 0, 0, 0, time.UTC)

	set.DTStart(dtstart)
	set.RDate(rdate)
	set.ExDate(exdate)

	// Test Recurrence() serialization
	recurrence := set.Recurrence(true)

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

func TestSetAllDay_StringSerialization(t *testing.T) {
	set := &Set{}
	set.SetAllDay(true)

	// Set up all-day event data
	dtstart := time.Date(2024, 8, 15, 0, 0, 0, 0, time.UTC)
	set.DTStart(dtstart)

	// Test String() serialization
	str := set.String(true)

	// Should contain VALUE=DATE format as per RFC 5545
	if !strings.Contains(str, "DTSTART;VALUE=DATE:20240815") {
		t.Errorf("String() should contain VALUE=DATE format, got: %s", str)
	}

	// Should not contain time part for all-day events
	if strings.Contains(str, "20240815T000000") {
		t.Errorf("String() should not contain time part for all-day events, got: %s", str)
	}
}

func TestSetAllDay_NonAllDayPreservesTimezone(t *testing.T) {
	set := &Set{}
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

func TestSetAllDay_EdgeCases(t *testing.T) {
	set := &Set{}
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

// TestAllDaySetStringWithRDate 测试全天事件 Set 使用 RDATE 的字符串序列化
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
			set := &Set{}
			set.SetAllDay(true)
			set.DTStart(tc.dtstart)

			for _, rdate := range tc.rdates {
				set.RDate(rdate)
			}

			output := set.String(true)
			t.Logf("RDATE test %s output: %s", tc.name, output)

			// 验证所有期望的字符串都存在
			for _, expected := range tc.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected %s in output, got: %s", expected, output)
				}
			}

			// 验证 RDATE 使用 VALUE=DATE 格式且不包含时间部分
			lines := strings.Split(output, "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "RDATE") {
					if !strings.Contains(line, "VALUE=DATE") {
						t.Errorf("RDATE should use VALUE=DATE format for all-day events, got: %s", line)
					}
					// 检查是否包含时间部分（T 后面跟数字）
					if strings.Contains(line, "T") && !strings.Contains(line, "VALUE=DATE") {
						t.Errorf("RDATE should not contain time part for all-day events, got: %s", line)
					}
				}
			}
		})
	}
}

// TestAllDaySetStringWithExDate 测试全天事件 Set 使用 EXDATE 的字符串序列化
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
			set := &Set{}
			set.SetAllDay(true)
			set.DTStart(tc.dtstart)

			for _, exdate := range tc.exdates {
				set.ExDate(exdate)
			}

			output := set.String(true)
			t.Logf("EXDATE test %s output: %s", tc.name, output)

			// 验证所有期望的字符串都存在
			for _, expected := range tc.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected %s in output, got: %s", expected, output)
				}
			}

			// 验证 EXDATE 使用 VALUE=DATE 格式且不包含时间部分
			lines := strings.Split(output, "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "EXDATE") {
					if !strings.Contains(line, "VALUE=DATE") {
						t.Errorf("EXDATE should use VALUE=DATE format for all-day events, got: %s", line)
					}
					// 检查是否包含时间部分（T 后面跟数字）
					if strings.Contains(line, "T") && !strings.Contains(line, "VALUE=DATE") {
						t.Errorf("EXDATE should not contain time part for all-day events, got: %s", line)
					}
				}
			}
		})
	}
}

// TestAllDaySetStringComplex 测试全天事件 Set 的复合场景（RRULE + RDATE + EXDATE + UNTIL）
func TestAllDaySetStringComplex(t *testing.T) {
	set := &Set{}
	set.SetAllDay(true)

	// 设置 DTSTART
	dtstart := time.Date(2023, 9, 1, 14, 30, 0, 0, time.FixedZone("EST", -5*3600))
	set.DTStart(dtstart)

	// 添加 RRULE with UNTIL
	until := time.Date(2023, 9, 30, 23, 59, 59, 0, time.UTC)
	rrule, err := NewRRule(ROption{
		Freq:    WEEKLY,
		AllDay:  true,
		Dtstart: dtstart,
		Until:   until,
	})
	if err != nil {
		t.Fatal(err)
	}
	set.RRule(rrule)

	// 添加 RDATE
	set.RDate(time.Date(2023, 9, 15, 16, 0, 0, 0, time.FixedZone("JST", 9*3600)))
	set.RDate(time.Date(2023, 9, 25, 8, 30, 0, 0, time.UTC))

	// 添加 EXDATE
	set.ExDate(time.Date(2023, 9, 8, 12, 0, 0, 0, time.UTC))
	set.ExDate(time.Date(2023, 9, 22, 20, 15, 0, 0, time.FixedZone("CET", 1*3600)))

	output := set.String(true)
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

	// 验证所有日期相关字段都使用 DATE 格式，不包含时间部分
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if (strings.HasPrefix(line, "DTSTART") ||
			strings.HasPrefix(line, "RDATE") ||
			strings.HasPrefix(line, "EXDATE")) &&
			strings.Contains(line, "T") &&
			!strings.Contains(line, "VALUE=DATE") {
			t.Errorf("Date field should use VALUE=DATE format for all-day events, got: %s", line)
		}
		// UNTIL 在 RRULE 中应该使用 DATE 格式（不包含 T）
		// 检查 UNTIL 是否包含时间部分（T 后面跟数字）
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
}
