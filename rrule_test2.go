package rrule

import (
	"testing"
	"time"
)

func TestCompatibility(t *testing.T) {
	str := "FREQ=WEEKLY;DTSTART=20120201T093000Z;INTERVAL=5;WKST=TU;COUNT=2;UNTIL=20130130T230000Z;BYSETPOS=2;BYMONTH=3;BYYEARDAY=95;BYWEEKNO=1;BYDAY=MO,+2FR;BYHOUR=9;BYMINUTE=30;BYSECOND=0;BYEASTER=-1"
	r, _ := StrToRRule(str)
	want := "DTSTART:20120201T093000Z\nRRULE:FREQ=WEEKLY;INTERVAL=5;WKST=TU;COUNT=2;UNTIL=20130130T230000Z;BYSETPOS=2;BYMONTH=3;BYYEARDAY=95;BYWEEKNO=1;BYDAY=MO,+2FR;BYHOUR=9;BYMINUTE=30;BYSECOND=0;BYEASTER=-1"
	if s := r.String(); s != want {
		t.Errorf("StrToRRule(%q).String() = %q, want %q", str, s, want)
	}
	r, _ = StrToRRule(want)
	if s := r.String(); s != want {
		t.Errorf("StrToRRule(%q).String() = %q, want %q", want, want, want)
	}
}

func TestInvalidString(t *testing.T) {
	cases := []string{
		"",
		"    ",
		"FREQ",
		"FREQ=HELLO",
		"BYMONTH=",
		"FREQ=WEEKLY;HELLO=WORLD",
		"FREQ=WEEKLY;BYMONTHDAY=I",
		"FREQ=WEEKLY;BYDAY=M",
		"FREQ=WEEKLY;BYDAY=MQ",
		"FREQ=WEEKLY;BYDAY=+MO",
		"BYDAY=MO",
	}
	for _, item := range cases {
		if _, e := StrToRRule(item); e == nil {
			t.Errorf("StrToRRule(%q) = nil, want error", item)
		}
	}
}

func TestStrSetParseErrors(t *testing.T) {
	inputs := [][]string{
		{"RRULE:XXX"},
		{"RDATE;TZD=X:1"},
	}

	for _, ss := range inputs {
		if _, err := StrSliceToRRuleSet(ss); err == nil {
			t.Error("Expected parse error for rules: ", ss)
		}
	}
}

func TestStrSetEmptySliceParse(t *testing.T) {
	s, err := StrSliceToRRuleSet([]string{})
	if err != nil {
		t.Error(err)
	}
	if s == nil {
		t.Error("Empty set should not be nil")
	}
}

func TestRDateValueDateStr(t *testing.T) {
	t.Run("DefaultToUTC", func(t *testing.T) {
		input := []string{
			"RDATE;VALUE=DATE:20180223",
		}
		s, err := StrSliceToRRuleSet(input)
		if err != nil {
			t.Error(err)
		}
		d := s.GetRDate()[0]
		if !d.Equal(time.Date(2018, 02, 23, 0, 0, 0, 0, time.UTC)) {
			t.Error("Bad time parsed: ", d)
		}
	})

	t.Run("PreserveExplicitTimezone", func(t *testing.T) {
		denver, _ := time.LoadLocation("America/Denver")
		input := []string{
			"RDATE;VALUE=DATE;TZID=America/Denver:20180223",
		}
		s, err := StrSliceToRRuleSet(input)
		if err != nil {
			t.Error(err)
		}
		d := s.GetRDate()[0]
		if !d.Equal(time.Date(2018, 02, 23, 0, 0, 0, 0, denver)) {
			t.Error("Bad time parsed: ", d)
		}
	})
}

func TestSetStrCompatibility(t *testing.T) {
	badInputStrs := []string{
		"",
		"FREQ=DAILY;UNTIL=20180517T235959Z",
		"DTSTART:;",
		"RRULE:;",
	}

	for _, badInputStr := range badInputStrs {
		_, err := StrToRRuleSet(badInputStr)
		if err == nil {
			t.Fatalf("StrToRRuleSet(%s) didn't return error", badInputStr)
		}
	}

	inputStr := "DTSTART;TZID=America/New_York:20180101T090000\n" +
		"RRULE:FREQ=DAILY;UNTIL=20180517T235959Z\n" +
		"RRULE:FREQ=WEEKLY;INTERVAL=2;BYDAY=MO,TU\n" +
		"EXRULE:FREQ=MONTHLY;UNTIL=20180520;BYMONTHDAY=1,2,3\n" +
		"EXDATE;VALUE=DATE-TIME:20180525T070000Z,20180530T130000Z\n" +
		"RDATE;VALUE=DATE-TIME:20180801T131313Z,20180902T141414Z\n"

	set, err := StrToRRuleSet(inputStr)
	if err != nil {
		t.Fatalf("StrToRRuleSet(%s) returned error: %v", inputStr, err)
	}

	nyLoc, _ := time.LoadLocation("America/New_York")
	dtWantTime := time.Date(2018, 1, 1, 9, 0, 0, 0, nyLoc)

	rrule := set.GetRRule()
	if rrule.String() != "DTSTART;TZID=America/New_York:20180101T090000\nRRULE:FREQ=WEEKLY;INTERVAL=2;BYDAY=MO,TU" {
		t.Errorf("Unexpected rrule: %s", rrule.String())
	}
	if !dtWantTime.Equal(rrule.dtstart) {
		t.Fatalf("Expected RRule dtstart to be %v got %v", dtWantTime, rrule.dtstart)
	}
	if !dtWantTime.Equal(set.GetDTStart()) {
		t.Fatalf("Expected Set dtstart to be %v got %v", dtWantTime, set.GetDTStart())
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

	dtWantAfter := time.Date(2018, 1, 2, 9, 0, 0, 0, nyLoc)
	dtAfter := set.After(dtWantTime, false)
	if !dtWantAfter.Equal(dtAfter) {
		t.Errorf("Next time wrong should be %s but is %s", dtWantAfter, dtAfter)
	}

	// String to set to string comparison
	setStr := set.String(true)
	setFromSetStr, _ := StrToRRuleSet(setStr)

	if setStr != setFromSetStr.String(true) {
		t.Errorf("Expected string output\n %s \nbut got\n %s\n", setStr, setFromSetStr.String(true))
	}
}

func TestSetParseLocalTimes(t *testing.T) {
	moscow, _ := time.LoadLocation("Europe/Moscow")

	t.Run("DtstartTimeZoneIsUsed", func(t *testing.T) {
		input := []string{
			"DTSTART;TZID=Europe/Moscow:20180220T090000",
			"RDATE;VALUE=DATE-TIME:20180223T100000",
		}
		s, err := StrSliceToRRuleSet(input)
		if err != nil {
			t.Error(err)
		}
		d := s.GetRDate()[0]
		if !d.Equal(time.Date(2018, 02, 23, 10, 0, 0, 0, moscow)) {
			t.Error("Bad time parsed: ", d)
		}
	})

	t.Run("DtstartTimeZoneValidOutput", func(t *testing.T) {
		input := []string{
			"DTSTART;TZID=Europe/Moscow:20180220T090000",
			"RDATE;VALUE=DATE-TIME:20180223T100000",
		}
		expected := "DTSTART;TZID=Europe/Moscow:20180220T090000\nRDATE;TZID=Europe/Moscow:20180223T100000"
		s, err := StrSliceToRRuleSet(input)
		if err != nil {
			t.Error(err)
		}

		sRRule := s.String(true)

		if sRRule != expected {
			t.Errorf("DTSTART output not valid. Expected: \n%s \n Got: \n%s", expected, sRRule)
		}
	})

	t.Run("DtstartUTCValidOutput", func(t *testing.T) {
		input := []string{
			"DTSTART:20180220T090000Z",
			"RDATE;VALUE=DATE-TIME:20180223T100000",
		}
		expected := "DTSTART:20180220T090000Z\nRDATE:20180223T100000Z"
		s, err := StrSliceToRRuleSet(input)
		if err != nil {
			t.Error(err)
		}

		sRRule := s.String(true)

		if sRRule != expected {
			t.Errorf("DTSTART output not valid. Expected: \n%s \n Got: \n%s", expected, sRRule)
		}
	})

	t.Run("SpecifiedDefaultZoneIsUsed", func(t *testing.T) {
		input := []string{
			"RDATE;VALUE=DATE-TIME:20180223T100000",
		}
		s, err := StrSliceToRRuleSetInLoc(input, moscow)
		if err != nil {
			t.Error(err)
		}
		d := s.GetRDate()[0]
		if !d.Equal(time.Date(2018, 02, 23, 10, 0, 0, 0, moscow)) {
			t.Error("Bad time parsed: ", d)
		}
	})
}

func TestRFCSetToString(t *testing.T) {
	nyLoc, _ := time.LoadLocation("America/New_York")
	dtStart := time.Date(2018, 1, 1, 9, 0, 0, 0, nyLoc)

	r, _ := NewRRule(ROption{Freq: MONTHLY, Dtstart: dtStart})
	want := "DTSTART;TZID=America/New_York:20180101T090000\nRRULE:FREQ=MONTHLY"
	if r.String() != want {
		t.Errorf("Expected RFC string %s, got %v", want, r.String())
	}

	expectedSetStr := "DTSTART;TZID=America/New_York:20180101T090000\nRRULE:FREQ=MONTHLY"

	set := Set{}
	set.RRule(r)
	set.DTStart(dtStart)
	if set.String(true) != expectedSetStr {
		t.Errorf("Expected RFC Set string %s, got %s", expectedSetStr, set.String(true))
	}
}

func TestRFCRuleToStr(t *testing.T) {
	nyLoc, _ := time.LoadLocation("America/New_York")
	dtStart := time.Date(2018, 1, 1, 9, 0, 0, 0, nyLoc)

	r, _ := NewRRule(ROption{Freq: MONTHLY, Dtstart: dtStart})
	want := "DTSTART;TZID=America/New_York:20180101T090000\nRRULE:FREQ=MONTHLY"
	if r.String() != want {
		t.Errorf("Expected RFC string %s, got %v", want, r.String())
	}
}
