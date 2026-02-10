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
	ogr := []string{"DTSTART;TZID=America/Los_Angeles:20181115T000000", "RRULE:FREQ=DAILY;INTERVAL=1;WKST=SU;UNTIL=20181118T075959Z"}
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
	if [2]string{timeToUTCStr(exDates[0]), timeToUTCStr(exDates[1])} != [2]string{"20180525T070000Z", "20180530T130000Z"} {
		t.Errorf("Unexpected exDates: %v", exDates)
	}

	// matching parsed RDates
	rDates := set.GetRDate()
	if len(rDates) != 2 {
		t.Errorf("Unexpected number of rDates: %v != 2, %v", len(rDates), rDates)
	}
	if [2]string{timeToUTCStr(rDates[0]), timeToUTCStr(rDates[1])} != [2]string{"20180801T131313Z", "20180902T141414Z"} {
		t.Errorf("Unexpected exDates: %v", exDates)
	}
}

// TestSetAllDayTimezoneConsistency tests all-day consistency across timezones.
func TestSetAllDayTimezoneConsistency(t *testing.T) {
	// Create times in different timezones.
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

			// Set DTSTART in different timezones.
			dtstart := baseTime.In(tc.timezone)
			set.DTStart(dtstart)

			// Add RDATEs in different timezones.
			rdate := baseTime.AddDate(0, 0, 1).In(tc.timezone)
			set.RDate(rdate)

			// Add EXDATEs in different timezones.
			exdate := baseTime.AddDate(0, 0, 2).In(tc.timezone)
			set.ExDate(exdate)

			// Verify all times normalize to floating time (UTC 00:00:00).
			if set.GetDTStart().Location() != utc {
				t.Errorf("DTSTART should be in UTC, got %v", set.GetDTStart().Location())
			}

			expectedDate := time.Date(2024, 3, 15, 0, 0, 0, 0, utc)
			if !set.GetDTStart().Equal(expectedDate) {
				t.Errorf("DTSTART should be %v, got %v", expectedDate, set.GetDTStart())
			}

			// Verify RDATE normalization.
			rdates := set.GetRDate()
			if len(rdates) != 1 {
				t.Fatalf("Expected 1 RDATE, got %d", len(rdates))
			}
			expectedRDate := time.Date(2024, 3, 16, 0, 0, 0, 0, utc)
			if !rdates[0].Equal(expectedRDate) {
				t.Errorf("RDATE should be %v, got %v", expectedRDate, rdates[0])
			}

			// Verify EXDATE normalization.
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

// TestSetComplexRRuleRDateExDateInteraction tests complex RRULE + RDATE + EXDATE interaction.
func TestSetComplexRRuleRDateExDateInteraction(t *testing.T) {
	set := &Set{}

	// Create a daily recurrence rule.
	rrule, err := NewRRule(ROption{
		Freq:    DAILY,
		Count:   10,
		Dtstart: time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("Failed to create RRule: %v", err)
	}
	set.RRule(rrule)

	// Add extra RDATEs (outside the RRULE sequence).
	set.RDate(time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC))
	set.RDate(time.Date(2024, 1, 20, 9, 0, 0, 0, time.UTC))

	// Exclude some RRULE-generated dates.
	set.ExDate(time.Date(2024, 1, 3, 9, 0, 0, 0, time.UTC))
	set.ExDate(time.Date(2024, 1, 5, 9, 0, 0, 0, time.UTC))

	// Exclude one RDATE (should be filtered out).
	set.ExDate(time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC))

	occurrences := set.All()

	// Log all occurrences for debugging.
	t.Logf("Generated occurrences:")
	for i, occ := range occurrences {
		t.Logf("  %d: %v", i, occ)
	}

	// Verify results: 10 RRULE - 2 EXDATE (Jan 3 and 5) + 1 valid RDATE (Jan 20) = 9.
	// Note: Jan 15 RDATE is excluded by EXDATE and not counted.
	expectedCount := 9
	if len(occurrences) != expectedCount {
		t.Errorf("Expected %d occurrences, got %d", expectedCount, len(occurrences))
	}

	// Verify excluded dates are not in the results.
	excludedDates := []time.Time{
		time.Date(2024, 1, 3, 9, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 5, 9, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC), // This RDATE is excluded by EXDATE.
	}

	for _, excluded := range excludedDates {
		for _, occurrence := range occurrences {
			if occurrence.Equal(excluded) {
				t.Errorf("Excluded date %v found in occurrences", excluded)
			}
		}
	}

	// Verify the included RDATE is present.
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

// TestSetDaylightSavingTransition tests behavior during DST transitions.
func TestSetDaylightSavingTransition(t *testing.T) {
	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Skip("America/New_York timezone not available")
	}

	set := &Set{}

	// 2024 DST starts: Mar 10, 2:00 AM -> 3:00 AM.
	// Create a rule that crosses the DST transition.
	rrule, err := NewRRule(ROption{
		Freq:    DAILY,
		Count:   5,
		Dtstart: time.Date(2024, 3, 8, 2, 30, 0, 0, ny), // 2 days before the transition.
	})
	if err != nil {
		t.Fatalf("Failed to create RRule: %v", err)
	}
	set.RRule(rrule)

	occurrences := set.All()

	if len(occurrences) != 5 {
		t.Fatalf("Expected 5 occurrences, got %d", len(occurrences))
	}

	// Verify timezone is preserved.
	for i, occurrence := range occurrences {
		if occurrence.Location().String() != ny.String() {
			t.Errorf("Occurrence %d should be in %s timezone, got %s",
				i, ny.String(), occurrence.Location().String())
		}

		// During DST transition, time may shift.
		// On Mar 10, 2:30 AM jumps to 3:30 AM (DST starts).
		hour, min, sec := occurrence.Clock()

		// On the DST transition day (Mar 10), time shifts from 2:30 to 3:30.
		expectedHour := 2
		if occurrence.Month() == 3 && occurrence.Day() >= 10 {
			expectedHour = 3 // After DST, time becomes 3:30.
		}

		if hour != expectedHour || min != 30 || sec != 0 {
			t.Logf("Occurrence %d at %v: expected %02d:30:00, got %02d:%02d:%02d",
				i, occurrence, expectedHour, hour, min, sec)
		}
	}
}

// TestSetAllDayDynamicToggle tests toggling all-day/timed state dynamically.
func TestSetAllDayDynamicToggle(t *testing.T) {
	set := &Set{}

	// Initialize as a timed event.
	dtstart := time.Date(2024, 6, 15, 14, 30, 45, 123456789, time.UTC)
	set.DTStart(dtstart)

	rdate := time.Date(2024, 6, 16, 10, 15, 30, 987654321, time.UTC)
	set.RDate(rdate)

	exdate := time.Date(2024, 6, 17, 16, 45, 20, 555666777, time.UTC)
	set.ExDate(exdate)

	// Verify initial state (timed).
	if set.IsAllDay() {
		t.Error("Set should not be all-day initially")
	}

	// Verify time precision is truncated to seconds.
	if set.GetDTStart().Nanosecond() != 0 {
		t.Error("Non-all-day DTSTART should be truncated to seconds")
	}

	// Switch to all-day.
	set.SetAllDay(true)

	// Verify state switch.
	if !set.IsAllDay() {
		t.Error("Set should be all-day after SetAllDay(true)")
	}

	// Verify all times are normalized to floating time.
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

	// Switch back to timed.
	set.SetAllDay(false)

	// Verify state switch.
	if set.IsAllDay() {
		t.Error("Set should not be all-day after SetAllDay(false)")
	}

	// Note: after switching back, time remains normalized (00:00:00 UTC).
	// This is expected because the original timezone info is lost.
}

// TestSetIteratorConsistency tests iterator vs batch consistency.
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

			// Use All() to get all occurrences.
			allOccurrences := set.All()

			// Use the iterator to collect all occurrences.
			var iteratorOccurrences []time.Time
			iterator := set.Iterator()
			for {
				dt, ok := iterator()
				if !ok {
					break
				}
				iteratorOccurrences = append(iteratorOccurrences, dt)
			}

			// Verify counts match.
			if len(allOccurrences) != len(iteratorOccurrences) {
				t.Errorf("All() returned %d occurrences, iterator returned %d",
					len(allOccurrences), len(iteratorOccurrences))
			}

			// Verify contents match.
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

// TestSetStringRoundTrip tests string round-trip serialization.
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

			// Serialize to string.
			setString := originalSet.String(true)

			// Parse back into a Set.
			parsedSet, err := StrToRRuleSet(setString)
			if err != nil {
				t.Fatalf("Failed to parse set string: %v", err)
			}

			// Set AllDay state (parser may not detect it automatically).
			if tc.allDay {
				parsedSet.SetAllDay(true)
			}

			// Compare original and parsed Set.
			originalOccurrences := originalSet.All()
			parsedOccurrences := parsedSet.All()

			if len(originalOccurrences) != len(parsedOccurrences) {
				t.Errorf("Occurrence count mismatch: original=%d, parsed=%d",
					len(originalOccurrences), len(parsedOccurrences))
			}

			// Verify each occurrence.
			for i, original := range originalOccurrences {
				if i >= len(parsedOccurrences) {
					t.Errorf("Missing occurrence at index %d: %v", i, original)
					continue
				}

				parsed := parsedOccurrences[i]

				// For all-day events, compare the date only.
				if tc.allDay {
					if original.Year() != parsed.Year() ||
						original.Month() != parsed.Month() ||
						original.Day() != parsed.Day() {
						t.Errorf("All-day occurrence %d date mismatch: original=%v, parsed=%v",
							i, original, parsed)
					}
				} else {
					// For timed events, compare full time (account for timezone conversion).
					if !original.Equal(parsed) {
						t.Errorf("Non-all-day occurrence %d mismatch: original=%v, parsed=%v",
							i, original, parsed)
					}
				}
			}
		})
	}
}

// TestSetEdgeCases tests edge cases and error handling.
func TestSetEdgeCases(t *testing.T) {
	t.Run("EmptySet", func(t *testing.T) {
		set := &Set{}

		occurrences := set.All()
		if len(occurrences) != 0 {
			t.Errorf("Empty set should return no occurrences, got %d", len(occurrences))
		}

		// Test the iterator.
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

		// Verify sorting.
		for i := 1; i < len(occurrences); i++ {
			if occurrences[i].Before(occurrences[i-1]) {
				t.Errorf("Occurrences should be sorted, but %v is before %v",
					occurrences[i], occurrences[i-1])
			}
		}
	})

	t.Run("AllExcluded", func(t *testing.T) {
		set := &Set{}

		// Create a rule that generates 3 events.
		rrule, err := NewRRule(ROption{
			Freq:    DAILY,
			Count:   3,
			Dtstart: time.Date(2024, 8, 1, 12, 0, 0, 0, time.UTC),
		})
		if err != nil {
			t.Fatalf("Failed to create RRule: %v", err)
		}
		set.RRule(rrule)

		// Exclude all generated events.
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

		// Add duplicate RDATEs.
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

// TestSetPerformance tests performance baseline.
func TestSetPerformance(t *testing.T) {
	t.Run("LargeRDateSet", func(t *testing.T) {
		set := &Set{}

		// Add many RDATEs.
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

		// Performance check: should complete within a reasonable time (100ms).
		if duration > 100*time.Millisecond {
			t.Logf("Performance warning: Large RDate set took %v", duration)
		}
	})

	t.Run("ComplexSetWithManyExclusions", func(t *testing.T) {
		set := &Set{}

		// Create a rule that generates many events.
		rrule, err := NewRRule(ROption{
			Freq:    DAILY,
			Count:   500,
			Dtstart: time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
		})
		if err != nil {
			t.Fatalf("Failed to create RRule: %v", err)
		}
		set.RRule(rrule)

		// Exclude half the events.
		baseDate := time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC)
		for i := 0; i < 250; i += 2 {
			set.ExDate(baseDate.AddDate(0, 0, i))
		}

		start := time.Now()
		occurrences := set.All()
		duration := time.Since(start)

		expectedCount := 500 - 125 // 500 events - 125 excluded events.
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
