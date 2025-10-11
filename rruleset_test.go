// 2017-2022, Teambition. All rights reserved.

package rrule

import (
	"strings"
	"testing"
	"time"
)

func TestSet(t *testing.T) {
	set := Set{}
	r, _ := NewRRule(ROption{Freq: YEARLY, Count: 2, Byweekday: []Weekday{TU},
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	set.RRule(r)
	value := set.All()
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 9, 9, 0, 0, 0, time.UTC)}
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSetOverlapping(t *testing.T) {
	set := Set{}
	r, _ := NewRRule(ROption{Freq: YEARLY,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	set.RRule(r)
	v1 := set.All()
	if len(v1) > 300 || len(v1) < 200 {
		t.Errorf("No default Util time")
	}
}

func TestSetString(t *testing.T) {
	moscow, _ := time.LoadLocation("Europe/Moscow")
	newYork, _ := time.LoadLocation("America/New_York")
	tehran, _ := time.LoadLocation("Asia/Tehran")

	set := Set{}
	r, _ := NewRRule(ROption{Freq: YEARLY, Count: 1, Byweekday: []Weekday{TU},
		Dtstart: time.Date(1997, 9, 2, 8, 0, 0, 0, time.UTC)})
	set.RRule(r)
	set.ExDate(time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC))
	set.ExDate(time.Date(1997, 9, 11, 9, 0, 0, 0, time.UTC).In(moscow))
	set.ExDate(time.Date(1997, 9, 18, 9, 0, 0, 0, time.UTC).In(newYork))
	set.RDate(time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC).In(tehran))
	set.RDate(time.Date(1997, 9, 9, 9, 0, 0, 0, time.UTC))

	want := `DTSTART:19970902T080000Z
RRULE:FREQ=YEARLY;COUNT=1;BYDAY=TU
RDATE;TZID=Asia/Tehran:19970904T133000
RDATE:19970909T090000Z
EXDATE:19970904T090000Z
EXDATE;TZID=Europe/Moscow:19970911T130000
EXDATE;TZID=America/New_York:19970918T050000`
	value := set.String(true)
	if want != value {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSetDTStart(t *testing.T) {
	set := Set{}
	r, _ := NewRRule(ROption{Freq: YEARLY, Count: 1, Byweekday: []Weekday{TU},
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	set.RRule(r)
	set.ExDate(time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC))
	set.ExDate(time.Date(1997, 9, 11, 9, 0, 0, 0, time.UTC))
	set.ExDate(time.Date(1997, 9, 18, 9, 0, 0, 0, time.UTC))
	set.RDate(time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC))
	set.RDate(time.Date(1997, 9, 9, 9, 0, 0, 0, time.UTC))

	nyLoc, _ := time.LoadLocation("America/New_York")
	set.DTStart(time.Date(1997, 9, 3, 9, 0, 0, 0, nyLoc))

	want := `DTSTART;TZID=America/New_York:19970903T090000
RRULE:FREQ=YEARLY;COUNT=1;BYDAY=TU
RDATE:19970904T090000Z
RDATE:19970909T090000Z
EXDATE:19970904T090000Z
EXDATE:19970911T090000Z
EXDATE:19970918T090000Z`
	value := set.String(true)
	if want != value {
		t.Errorf("get \n%v\n want \n%v\n", value, want)
	}

	sset, err := StrToRRuleSet(set.String(true))
	if err != nil {
		t.Errorf("Could not create RSET from set output")
	}
	if sset.String(true) != set.String(true) {
		t.Errorf("RSET created from set output different than original set, %s", sset.String(true))
	}
}

func TestSetRecurrence(t *testing.T) {
	set := Set{}
	r, _ := NewRRule(ROption{Freq: YEARLY, Count: 1, Byweekday: []Weekday{TU},
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	set.RRule(r)
	value := set.Recurrence(true)
	if len(value) != 2 {
		t.Errorf("Wrong length for recurrence got=%v want=%v", len(value), 2)
	}
	want := "DTSTART:19970902T090000Z\nRRULE:FREQ=YEARLY;COUNT=1;BYDAY=TU"
	if set.String(true) != want {
		t.Errorf("get %s, want %v", set.String(true), want)
	}
}

func TestSetDate(t *testing.T) {
	set := Set{}
	r, _ := NewRRule(ROption{Freq: YEARLY, Count: 1, Byweekday: []Weekday{TU},
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	set.RRule(r)
	set.RDate(time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC))
	set.RDate(time.Date(1997, 9, 9, 9, 0, 0, 0, time.UTC))
	value := set.All()
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 9, 9, 0, 0, 0, time.UTC)}
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSetRDates(t *testing.T) {
	set := Set{}
	r, _ := NewRRule(ROption{Freq: YEARLY, Count: 1, Byweekday: []Weekday{TU},
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	set.RRule(r)
	set.SetRDates([]time.Time{
		time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 9, 9, 0, 0, 0, time.UTC),
	})
	value := set.All()
	want := []time.Time{
		time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 9, 9, 0, 0, 0, time.UTC),
	}
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSetExDate(t *testing.T) {
	set := Set{}
	r, _ := NewRRule(ROption{Freq: YEARLY, Count: 6, Byweekday: []Weekday{TU, TH},
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	set.RRule(r)
	set.ExDate(time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC))
	set.ExDate(time.Date(1997, 9, 11, 9, 0, 0, 0, time.UTC))
	set.ExDate(time.Date(1997, 9, 18, 9, 0, 0, 0, time.UTC))
	value := set.All()
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 9, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 16, 9, 0, 0, 0, time.UTC)}
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSetExDates(t *testing.T) {
	set := Set{}
	r, _ := NewRRule(ROption{Freq: YEARLY, Count: 6, Byweekday: []Weekday{TU, TH},
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	set.RRule(r)
	set.SetExDates([]time.Time{
		time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 11, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 18, 9, 0, 0, 0, time.UTC),
	})
	value := set.All()
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 9, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 16, 9, 0, 0, 0, time.UTC)}
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSetExDateRevOrder(t *testing.T) {
	set := Set{}
	r, _ := NewRRule(ROption{Freq: MONTHLY, Count: 5, Bymonthday: []int{10},
		Dtstart: time.Date(2004, 1, 1, 9, 0, 0, 0, time.UTC)})
	set.RRule(r)
	set.ExDate(time.Date(2004, 4, 10, 9, 0, 0, 0, time.UTC))
	set.ExDate(time.Date(2004, 2, 10, 9, 0, 0, 0, time.UTC))
	value := set.All()
	want := []time.Time{time.Date(2004, 1, 10, 9, 0, 0, 0, time.UTC),
		time.Date(2004, 3, 10, 9, 0, 0, 0, time.UTC),
		time.Date(2004, 5, 10, 9, 0, 0, 0, time.UTC)}
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSetDateAndExDate(t *testing.T) {
	set := Set{}
	set.RDate(time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC))
	set.RDate(time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC))
	set.RDate(time.Date(1997, 9, 9, 9, 0, 0, 0, time.UTC))
	set.RDate(time.Date(1997, 9, 11, 9, 0, 0, 0, time.UTC))
	set.RDate(time.Date(1997, 9, 16, 9, 0, 0, 0, time.UTC))
	set.RDate(time.Date(1997, 9, 18, 9, 0, 0, 0, time.UTC))
	set.ExDate(time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC))
	set.ExDate(time.Date(1997, 9, 11, 9, 0, 0, 0, time.UTC))
	set.ExDate(time.Date(1997, 9, 18, 9, 0, 0, 0, time.UTC))
	value := set.All()
	want := []time.Time{time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 9, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 16, 9, 0, 0, 0, time.UTC)}
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSetBefore(t *testing.T) {
	set := Set{}
	r, _ := NewRRule(ROption{Freq: DAILY, Count: 7,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	set.RRule(r)
	want := time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC)
	value := set.Before(time.Date(1997, 9, 5, 9, 0, 0, 0, time.UTC), false)
	if value != want {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSetBeforeInc(t *testing.T) {
	set := Set{}
	r, _ := NewRRule(ROption{Freq: DAILY, Count: 7,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	set.RRule(r)
	want := time.Date(1997, 9, 5, 9, 0, 0, 0, time.UTC)
	value := set.Before(time.Date(1997, 9, 5, 9, 0, 0, 0, time.UTC), true)
	if value != want {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSetAfter(t *testing.T) {
	set := Set{}
	r, _ := NewRRule(ROption{Freq: DAILY, Count: 7,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	set.RRule(r)
	want := time.Date(1997, 9, 5, 9, 0, 0, 0, time.UTC)
	value := set.After(time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC), false)
	if value != want {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSetAfterInc(t *testing.T) {
	set := Set{}
	r, _ := NewRRule(ROption{Freq: DAILY, Count: 7,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	set.RRule(r)
	want := time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC)
	value := set.After(time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC), true)
	if value != want {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSetBetween(t *testing.T) {
	set := Set{}
	r, _ := NewRRule(ROption{Freq: DAILY, Count: 7,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	set.RRule(r)
	value := set.Between(time.Date(1997, 9, 3, 9, 0, 0, 0, time.UTC), time.Date(1997, 9, 6, 9, 0, 0, 0, time.UTC), false)
	want := []time.Time{time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 5, 9, 0, 0, 0, time.UTC)}
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSetBetweenInc(t *testing.T) {
	set := Set{}
	r, _ := NewRRule(ROption{Freq: DAILY, Count: 7,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC)})
	set.RRule(r)
	value := set.Between(time.Date(1997, 9, 3, 9, 0, 0, 0, time.UTC), time.Date(1997, 9, 6, 9, 0, 0, 0, time.UTC), true)
	want := []time.Time{time.Date(1997, 9, 3, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 4, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 5, 9, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 6, 9, 0, 0, 0, time.UTC)}
	if !timesEqual(value, want) {
		t.Errorf("get %v, want %v", value, want)
	}
}

func TestSetTrickyTimeZones(t *testing.T) {
	set := Set{}

	moscow, _ := time.LoadLocation("Europe/Moscow")
	newYork, _ := time.LoadLocation("America/New_York")
	tehran, _ := time.LoadLocation("Asia/Tehran")

	r, _ := NewRRule(ROption{
		Freq:    DAILY,
		Count:   4,
		Dtstart: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).In(moscow),
	})
	set.RRule(r)

	set.ExDate(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).In(newYork))
	set.ExDate(time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC).In(tehran))
	set.ExDate(time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC).In(moscow))
	set.ExDate(time.Date(2000, 1, 4, 0, 0, 0, 0, time.UTC))

	occurrences := set.All()

	if len(occurrences) > 0 {
		t.Errorf("No all occurrences excluded by ExDate: [%+v]", occurrences)
	}
}

func TestSetDtStart(t *testing.T) {
	ogr := []string{"DTSTART;TZID=America/Los_Angeles:20181115T000000", "RRULE:FREQ=DAILY;INTERVAL=1;WKST=SU;UNTIL=20181117T235959"}
	set, _ := StrSliceToRRuleSet(ogr)

	ogoc := set.All()
	set.DTStart(set.GetDTStart().AddDate(0, 0, 1))

	noc := set.All()
	if len(noc) != len(ogoc)-1 {
		t.Fatalf("As per the new DTStart the new occurences should exactly be one less that the original, new :%d original: %d", len(noc), len(ogoc))
	}

	for i := range noc {
		if noc[i] != ogoc[i+1] {
			t.Errorf("New occurences should just offset by one, mismatch at %d, expected: %+v, actual: %+v", i, ogoc[i+1], noc[i])
		}
	}
}

func TestRuleSetChangeDTStartTimezoneRespected(t *testing.T) {
	/*
		https://golang.org/pkg/time/#LoadLocation

		"The time zone database needed by LoadLocation may not be present on all systems, especially non-Unix systems.
		LoadLocation looks in the directory or uncompressed zip file named by the ZONEINFO environment variable,
		if any, then looks in known installation locations on Unix systems, and finally looks in
		$GOROOT/lib/time/zoneinfo.zip."
	*/
	loc, err := time.LoadLocation("CET")
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	ruleSet := &Set{}
	rule, err := NewRRule(
		ROption{
			Freq:     DAILY,
			Count:    10,
			Wkst:     MO,
			Byhour:   []int{10},
			Byminute: []int{0},
			Bysecond: []int{0},
			Dtstart:  time.Date(2019, 3, 6, 0, 0, 0, 0, loc),
		},
	)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
	ruleSet.RRule(rule)
	ruleSet.DTStart(time.Date(2019, 3, 6, 0, 0, 0, 0, time.UTC))

	events := ruleSet.All()
	if len(events) != 10 {
		t.Fatal("expected", 10, "got", len(events))
	}

	for _, e := range events {
		if e.Location().String() != "UTC" {
			t.Fatal("expected", "UTC", "got", e.Location().String())
		}
	}
}

func TestSetStr(t *testing.T) {
	setStr := "RRULE:FREQ=DAILY;UNTIL=20180517T235959Z\n" +
		"EXDATE;VALUE=DATE-TIME:20180525T070000Z,20180530T130000Z\n" +
		"RDATE;VALUE=DATE-TIME:20180801T131313Z,20180902T141414Z\n"

	set, err := StrToRRuleSet(setStr)
	if err != nil {
		t.Fatalf("StrToRRuleSet(%s) returned error: %v", setStr, err)
	}

	rule := set.GetRRule()
	if rule == nil {
		t.Errorf("Unexpected rrule parsed")
	}
	if rule.String() != "FREQ=DAILY;UNTIL=20180517T235959Z" {
		t.Errorf("Unexpected rrule: %s", rule.String())
	}

	// matching parsed EXDates
	exDates := set.GetExDate()
	if len(exDates) != 2 {
		t.Errorf("Unexpected number of exDates: %v != 2, %v", len(exDates), exDates)
	}
	if [2]string{timeToStr(exDates[0]), timeToStr(exDates[1])} != [2]string{"20180525T070000Z", "20180530T130000Z"} {
		t.Errorf("Unexpected exDates: %v", exDates)
	}

	// matching parsed RDates
	rDates := set.GetRDate()
	if len(rDates) != 2 {
		t.Errorf("Unexpected number of rDates: %v != 2, %v", len(rDates), rDates)
	}
	if [2]string{timeToStr(rDates[0]), timeToStr(rDates[1])} != [2]string{"20180801T131313Z", "20180902T141414Z"} {
		t.Errorf("Unexpected exDates: %v", exDates)
	}
}

// TestSetAllDayTimezoneConsistency 测试全天事件在不同时区下的一致性
func TestSetAllDayTimezoneConsistency(t *testing.T) {
	// 创建不同时区的时间
	utc := time.UTC
	ny, _ := time.LoadLocation("America/New_York")
	tokyo, _ := time.LoadLocation("Asia/Tokyo")

	baseTime := time.Date(2024, 3, 15, 14, 30, 45, 0, utc)

	testCases := []struct {
		name     string
		timezone *time.Location
	}{
		{"UTC", utc},
		{"New_York", ny},
		{"Tokyo", tokyo},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			set := &Set{}
			set.SetAllDay(true)

			// 设置不同时区的 DTSTART
			dtstart := baseTime.In(tc.timezone)
			set.DTStart(dtstart)

			// 添加不同时区的 RDATE
			rdate := baseTime.AddDate(0, 0, 1).In(tc.timezone)
			set.RDate(rdate)

			// 添加不同时区的 EXDATE
			exdate := baseTime.AddDate(0, 0, 2).In(tc.timezone)
			set.ExDate(exdate)

			// 验证所有时间都被标准化为浮动时间（UTC 00:00:00）
			if set.GetDTStart().Location() != utc {
				t.Errorf("DTSTART should be in UTC, got %v", set.GetDTStart().Location())
			}

			expectedDate := time.Date(2024, 3, 15, 0, 0, 0, 0, utc)
			if !set.GetDTStart().Equal(expectedDate) {
				t.Errorf("DTSTART should be %v, got %v", expectedDate, set.GetDTStart())
			}

			// 验证 RDATE 标准化
			rdates := set.GetRDate()
			if len(rdates) != 1 {
				t.Fatalf("Expected 1 RDATE, got %d", len(rdates))
			}
			expectedRDate := time.Date(2024, 3, 16, 0, 0, 0, 0, utc)
			if !rdates[0].Equal(expectedRDate) {
				t.Errorf("RDATE should be %v, got %v", expectedRDate, rdates[0])
			}

			// 验证 EXDATE 标准化
			exdates := set.GetExDate()
			if len(exdates) != 1 {
				t.Fatalf("Expected 1 EXDATE, got %d", len(exdates))
			}
			expectedExDate := time.Date(2024, 3, 17, 0, 0, 0, 0, utc)
			if !exdates[0].Equal(expectedExDate) {
				t.Errorf("EXDATE should be %v, got %v", expectedExDate, exdates[0])
			}
		})
	}
}

// TestSetComplexRRuleRDateExDateInteraction 测试复杂的 RRULE + RDATE + EXDATE 交互
func TestSetComplexRRuleRDateExDateInteraction(t *testing.T) {
	set := &Set{}

	// 创建每日循环规则
	rrule, err := NewRRule(ROption{
		Freq:    DAILY,
		Count:   10,
		Dtstart: time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("Failed to create RRule: %v", err)
	}
	set.RRule(rrule)

	// 添加额外的 RDATE（不在 RRULE 生成的序列中）
	set.RDate(time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC))
	set.RDate(time.Date(2024, 1, 20, 9, 0, 0, 0, time.UTC))

	// 排除一些 RRULE 生成的日期
	set.ExDate(time.Date(2024, 1, 3, 9, 0, 0, 0, time.UTC))
	set.ExDate(time.Date(2024, 1, 5, 9, 0, 0, 0, time.UTC))

	// 排除一个 RDATE（应该被过滤掉）
	set.ExDate(time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC))

	occurrences := set.All()

	// 让我们先打印出所有事件来调试
	t.Logf("Generated occurrences:")
	for i, occ := range occurrences {
		t.Logf("  %d: %v", i, occ)
	}

	// 验证结果：10个RRULE - 2个EXDATE（1月3日和5日） + 1个有效RDATE（1月20日） = 9个
	// 注意：1月15日的RDATE被EXDATE排除了，所以不计入
	expectedCount := 9
	if len(occurrences) != expectedCount {
		t.Errorf("Expected %d occurrences, got %d", expectedCount, len(occurrences))
	}

	// 验证排除的日期不在结果中
	excludedDates := []time.Time{
		time.Date(2024, 1, 3, 9, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 5, 9, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC), // 这个RDATE被EXDATE排除
	}

	for _, excluded := range excludedDates {
		for _, occurrence := range occurrences {
			if occurrence.Equal(excluded) {
				t.Errorf("Excluded date %v found in occurrences", excluded)
			}
		}
	}

	// 验证包含的 RDATE 在结果中
	expectedRDate := time.Date(2024, 1, 20, 9, 0, 0, 0, time.UTC)
	found := false
	for _, occurrence := range occurrences {
		if occurrence.Equal(expectedRDate) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected RDATE %v not found in occurrences", expectedRDate)
	}
}

// TestSetDaylightSavingTransition 测试夏令时转换期间的行为
func TestSetDaylightSavingTransition(t *testing.T) {
	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Skip("America/New_York timezone not available")
	}

	set := &Set{}

	// 2024年夏令时开始：3月10日 2:00 AM -> 3:00 AM
	// 创建跨越夏令时转换的循环规则
	rrule, err := NewRRule(ROption{
		Freq:    DAILY,
		Count:   5,
		Dtstart: time.Date(2024, 3, 8, 2, 30, 0, 0, ny), // 夏令时转换前2天
	})
	if err != nil {
		t.Fatalf("Failed to create RRule: %v", err)
	}
	set.RRule(rrule)

	occurrences := set.All()

	if len(occurrences) != 5 {
		t.Fatalf("Expected 5 occurrences, got %d", len(occurrences))
	}

	// 验证时区信息保持一致
	for i, occurrence := range occurrences {
		if occurrence.Location().String() != ny.String() {
			t.Errorf("Occurrence %d should be in %s timezone, got %s",
				i, ny.String(), occurrence.Location().String())
		}

		// 在夏令时转换期间，时间可能会发生变化
		// 3月10日 2:30 AM 会跳到 3:30 AM（夏令时开始）
		hour, min, sec := occurrence.Clock()

		// 对于夏令时转换日（3月10日），时间会从2:30变为3:30
		expectedHour := 2
		if occurrence.Month() == 3 && occurrence.Day() >= 10 {
			expectedHour = 3 // 夏令时后时间变为3:30
		}

		if hour != expectedHour || min != 30 || sec != 0 {
			t.Logf("Occurrence %d at %v: expected %02d:30:00, got %02d:%02d:%02d",
				i, occurrence, expectedHour, hour, min, sec)
		}
	}
}

// TestSetAllDayDynamicToggle 测试动态切换全天/非全天状态
func TestSetAllDayDynamicToggle(t *testing.T) {
	set := &Set{}

	// 初始设置为非全天事件
	dtstart := time.Date(2024, 6, 15, 14, 30, 45, 123456789, time.UTC)
	set.DTStart(dtstart)

	rdate := time.Date(2024, 6, 16, 10, 15, 30, 987654321, time.UTC)
	set.RDate(rdate)

	exdate := time.Date(2024, 6, 17, 16, 45, 20, 555666777, time.UTC)
	set.ExDate(exdate)

	// 验证初始状态（非全天）
	if set.IsAllDay() {
		t.Error("Set should not be all-day initially")
	}

	// 验证时间精度保持到秒
	if set.GetDTStart().Nanosecond() != 0 {
		t.Error("Non-all-day DTSTART should be truncated to seconds")
	}

	// 切换到全天事件
	set.SetAllDay(true)

	// 验证状态切换
	if !set.IsAllDay() {
		t.Error("Set should be all-day after SetAllDay(true)")
	}

	// 验证所有时间被标准化为浮动时间
	expectedDTStart := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	if !set.GetDTStart().Equal(expectedDTStart) {
		t.Errorf("All-day DTSTART should be %v, got %v", expectedDTStart, set.GetDTStart())
	}

	rdates := set.GetRDate()
	expectedRDate := time.Date(2024, 6, 16, 0, 0, 0, 0, time.UTC)
	if len(rdates) != 1 || !rdates[0].Equal(expectedRDate) {
		t.Errorf("All-day RDATE should be %v, got %v", expectedRDate, rdates)
	}

	exdates := set.GetExDate()
	expectedExDate := time.Date(2024, 6, 17, 0, 0, 0, 0, time.UTC)
	if len(exdates) != 1 || !exdates[0].Equal(expectedExDate) {
		t.Errorf("All-day EXDATE should be %v, got %v", expectedExDate, exdates)
	}

	// 切换回非全天事件
	set.SetAllDay(false)

	// 验证状态切换
	if set.IsAllDay() {
		t.Error("Set should not be all-day after SetAllDay(false)")
	}

	// 注意：切换回非全天后，时间仍然是标准化的（00:00:00 UTC）
	// 这是预期行为，因为原始时区信息已丢失
}

// TestSetIteratorConsistency 测试迭代器与批量方法的一致性
func TestSetIteratorConsistency(t *testing.T) {
	testCases := []struct {
		name   string
		setup  func() *Set
		allDay bool
	}{
		{
			name: "NonAllDay_WithRRule",
			setup: func() *Set {
				set := &Set{}
				rrule, _ := NewRRule(ROption{
					Freq:    WEEKLY,
					Count:   5,
					Dtstart: time.Date(2024, 2, 1, 10, 0, 0, 0, time.UTC),
				})
				set.RRule(rrule)
				return set
			},
			allDay: false,
		},
		{
			name: "AllDay_WithRDateExDate",
			setup: func() *Set {
				set := &Set{}
				set.SetAllDay(true)
				set.DTStart(time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC))
				set.RDate(time.Date(2024, 2, 5, 0, 0, 0, 0, time.UTC))
				set.RDate(time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC))
				set.ExDate(time.Date(2024, 2, 5, 0, 0, 0, 0, time.UTC))
				return set
			},
			allDay: true,
		},
		{
			name: "Complex_Mixed",
			setup: func() *Set {
				set := &Set{}
				rrule, _ := NewRRule(ROption{
					Freq:    DAILY,
					Count:   7,
					Dtstart: time.Date(2024, 3, 1, 15, 30, 0, 0, time.UTC),
				})
				set.RRule(rrule)
				set.RDate(time.Date(2024, 3, 10, 15, 30, 0, 0, time.UTC))
				set.ExDate(time.Date(2024, 3, 3, 15, 30, 0, 0, time.UTC))
				return set
			},
			allDay: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			set := tc.setup()

			// 使用 All() 方法获取所有事件
			allOccurrences := set.All()

			// 使用迭代器手动收集所有事件
			var iteratorOccurrences []time.Time
			iterator := set.Iterator()
			for {
				dt, ok := iterator()
				if !ok {
					break
				}
				iteratorOccurrences = append(iteratorOccurrences, dt)
			}

			// 验证数量一致
			if len(allOccurrences) != len(iteratorOccurrences) {
				t.Errorf("All() returned %d occurrences, iterator returned %d",
					len(allOccurrences), len(iteratorOccurrences))
			}

			// 验证内容一致
			for i, expected := range allOccurrences {
				if i >= len(iteratorOccurrences) {
					t.Errorf("Iterator missing occurrence at index %d: %v", i, expected)
					continue
				}

				actual := iteratorOccurrences[i]
				if !expected.Equal(actual) {
					t.Errorf("Occurrence %d mismatch: All()=%v, Iterator()=%v",
						i, expected, actual)
				}
			}
		})
	}
}

// TestSetStringRoundTrip 测试字符串序列化和反序列化的往返一致性
func TestSetStringRoundTrip(t *testing.T) {
	testCases := []struct {
		name   string
		setup  func() *Set
		allDay bool
	}{
		{
			name: "AllDay_Complete",
			setup: func() *Set {
				set := &Set{}
				set.SetAllDay(true)
				set.DTStart(time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC))

				rrule, _ := NewRRule(ROption{
					Freq:    WEEKLY,
					Count:   4,
					Dtstart: time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC),
				})
				set.RRule(rrule)

				set.RDate(time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC))
				set.ExDate(time.Date(2024, 7, 11, 0, 0, 0, 0, time.UTC))

				return set
			},
			allDay: true,
		},
		{
			name: "NonAllDay_WithTimezone",
			setup: func() *Set {
				ny, _ := time.LoadLocation("America/New_York")
				set := &Set{}
				set.DTStart(time.Date(2024, 7, 4, 14, 30, 0, 0, ny))

				rrule, _ := NewRRule(ROption{
					Freq:    DAILY,
					Count:   3,
					Dtstart: time.Date(2024, 7, 4, 14, 30, 0, 0, ny),
				})
				set.RRule(rrule)

				set.RDate(time.Date(2024, 7, 10, 14, 30, 0, 0, ny))
				set.ExDate(time.Date(2024, 7, 5, 14, 30, 0, 0, ny))

				return set
			},
			allDay: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			originalSet := tc.setup()

			// 序列化为字符串
			setString := originalSet.String(true)

			// 反序列化回 Set
			parsedSet, err := StrToRRuleSet(setString)
			if err != nil {
				t.Fatalf("Failed to parse set string: %v", err)
			}

			// 设置 AllDay 状态（解析器可能不会自动检测）
			if tc.allDay {
				parsedSet.SetAllDay(true)
			}

			// 比较原始和解析后的 Set
			originalOccurrences := originalSet.All()
			parsedOccurrences := parsedSet.All()

			if len(originalOccurrences) != len(parsedOccurrences) {
				t.Errorf("Occurrence count mismatch: original=%d, parsed=%d",
					len(originalOccurrences), len(parsedOccurrences))
			}

			// 验证每个事件
			for i, original := range originalOccurrences {
				if i >= len(parsedOccurrences) {
					t.Errorf("Missing occurrence at index %d: %v", i, original)
					continue
				}

				parsed := parsedOccurrences[i]

				// 对于全天事件，只比较日期部分
				if tc.allDay {
					if original.Year() != parsed.Year() ||
						original.Month() != parsed.Month() ||
						original.Day() != parsed.Day() {
						t.Errorf("All-day occurrence %d date mismatch: original=%v, parsed=%v",
							i, original, parsed)
					}
				} else {
					// 对于非全天事件，比较完整时间（考虑时区转换）
					if !original.Equal(parsed) {
						t.Errorf("Non-all-day occurrence %d mismatch: original=%v, parsed=%v",
							i, original, parsed)
					}
				}
			}
		})
	}
}

// TestSetEdgeCases 测试边界情况和错误处理
func TestSetEdgeCases(t *testing.T) {
	t.Run("EmptySet", func(t *testing.T) {
		set := &Set{}

		occurrences := set.All()
		if len(occurrences) != 0 {
			t.Errorf("Empty set should return no occurrences, got %d", len(occurrences))
		}

		// 测试迭代器
		iterator := set.Iterator()
		if dt, ok := iterator(); ok {
			t.Errorf("Empty set iterator should return false, got %v", dt)
		}
	})

	t.Run("OnlyRDates", func(t *testing.T) {
		set := &Set{}

		dates := []time.Time{
			time.Date(2024, 5, 1, 10, 0, 0, 0, time.UTC),
			time.Date(2024, 5, 15, 10, 0, 0, 0, time.UTC),
			time.Date(2024, 5, 30, 10, 0, 0, 0, time.UTC),
		}

		for _, date := range dates {
			set.RDate(date)
		}

		occurrences := set.All()
		if len(occurrences) != 3 {
			t.Errorf("Expected 3 occurrences from RDates, got %d", len(occurrences))
		}

		// 验证排序
		for i := 1; i < len(occurrences); i++ {
			if occurrences[i].Before(occurrences[i-1]) {
				t.Errorf("Occurrences should be sorted, but %v is before %v",
					occurrences[i], occurrences[i-1])
			}
		}
	})

	t.Run("AllExcluded", func(t *testing.T) {
		set := &Set{}

		// 创建生成3个事件的规则
		rrule, err := NewRRule(ROption{
			Freq:    DAILY,
			Count:   3,
			Dtstart: time.Date(2024, 8, 1, 12, 0, 0, 0, time.UTC),
		})
		if err != nil {
			t.Fatalf("Failed to create RRule: %v", err)
		}
		set.RRule(rrule)

		// 排除所有生成的事件
		set.ExDate(time.Date(2024, 8, 1, 12, 0, 0, 0, time.UTC))
		set.ExDate(time.Date(2024, 8, 2, 12, 0, 0, 0, time.UTC))
		set.ExDate(time.Date(2024, 8, 3, 12, 0, 0, 0, time.UTC))

		occurrences := set.All()
		if len(occurrences) != 0 {
			t.Errorf("All events should be excluded, got %d occurrences", len(occurrences))
		}
	})

	t.Run("DuplicateRDates", func(t *testing.T) {
		set := &Set{}

		duplicateDate := time.Date(2024, 9, 15, 16, 30, 0, 0, time.UTC)

		// 添加重复的 RDATE
		set.RDate(duplicateDate)
		set.RDate(duplicateDate)
		set.RDate(duplicateDate)

		occurrences := set.All()
		if len(occurrences) != 1 {
			t.Errorf("Duplicate RDates should be deduplicated, expected 1 occurrence, got %d",
				len(occurrences))
		}

		if !occurrences[0].Equal(duplicateDate) {
			t.Errorf("Expected occurrence %v, got %v", duplicateDate, occurrences[0])
		}
	})
}

// TestSetPerformance 测试性能基准
func TestSetPerformance(t *testing.T) {
	t.Run("LargeRDateSet", func(t *testing.T) {
		set := &Set{}

		// 添加大量 RDATE
		baseDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		for i := 0; i < 1000; i++ {
			set.RDate(baseDate.AddDate(0, 0, i))
		}

		start := time.Now()
		occurrences := set.All()
		duration := time.Since(start)

		if len(occurrences) != 1000 {
			t.Errorf("Expected 1000 occurrences, got %d", len(occurrences))
		}

		// 性能检查：应该在合理时间内完成（这里设置为100ms）
		if duration > 100*time.Millisecond {
			t.Logf("Performance warning: Large RDate set took %v", duration)
		}
	})

	t.Run("ComplexSetWithManyExclusions", func(t *testing.T) {
		set := &Set{}

		// 创建生成大量事件的规则
		rrule, err := NewRRule(ROption{
			Freq:    DAILY,
			Count:   500,
			Dtstart: time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
		})
		if err != nil {
			t.Fatalf("Failed to create RRule: %v", err)
		}
		set.RRule(rrule)

		// 排除一半的事件
		baseDate := time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC)
		for i := 0; i < 250; i += 2 {
			set.ExDate(baseDate.AddDate(0, 0, i))
		}

		start := time.Now()
		occurrences := set.All()
		duration := time.Since(start)

		expectedCount := 500 - 125 // 500个事件 - 125个排除的事件
		if len(occurrences) != expectedCount {
			t.Errorf("Expected %d occurrences, got %d", expectedCount, len(occurrences))
		}

		t.Logf("Complex set with exclusions took %v", duration)
	})
}

func TestSetRRulePreservesTimezoneForTimedEvents(t *testing.T) {
	newYork, _ := time.LoadLocation("America/New_York")
	losAngeles, _ := time.LoadLocation("America/Los_Angeles")

	rruleStart := time.Date(2024, 11, 1, 9, 30, 0, 0, newYork)
	rrule, err := NewRRule(ROption{
		Freq:    DAILY,
		Count:   2,
		Dtstart: rruleStart,
	})
	if err != nil {
		t.Fatalf("Failed to create RRule: %v", err)
	}

	set := &Set{}
	set.RRule(rrule)

	set.RDate(time.Date(2024, 11, 3, 14, 45, 0, 0, losAngeles))
	set.ExDate(time.Date(2024, 11, 2, 9, 30, 0, 0, newYork))

	if set.IsAllDay() {
		t.Fatal("Timed events should not flip the set into all-day mode")
	}
	if set.GetDTStart().Location() != newYork {
		t.Fatalf("DTSTART should retain original timezone, got %v", set.GetDTStart().Location())
	}

	want := []time.Time{
		time.Date(2024, 11, 1, 9, 30, 0, 0, newYork),
		time.Date(2024, 11, 3, 14, 45, 0, 0, losAngeles),
	}
	got := set.All()
	if !timesEqual(got, want) {
		t.Fatalf("Unexpected iterator results, want %v got %v", want, got)
	}

	recurrence := set.Recurrence(true)
	ensureContains := func(substr string) {
		for _, line := range recurrence {
			if strings.Contains(line, substr) {
				return
			}
		}
		t.Fatalf("Recurrence output missing %q in %v", substr, recurrence)
	}

	ensureContains("DTSTART;TZID=America/New_York:20241101T093000")
	ensureContains("RDATE;TZID=America/Los_Angeles:20241103T144500")
	ensureContains("EXDATE;TZID=America/New_York:20241102T093000")
}
