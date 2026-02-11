package rrule

import (
	"strings"
	"testing"
	"time"
)

func parseRecurrence(input string) (*Recurrence, error) {
	option, err := StrToROption(input)
	if err != nil {
		return nil, err
	}
	return newRecurrence(*option)
}

func TestCompatibility(t *testing.T) {
	str := "FREQ=WEEKLY;DTSTART=20120201T093000Z;INTERVAL=5;WKST=TU;COUNT=2;UNTIL=20130130T230000Z;BYSETPOS=2;BYMONTH=3;BYYEARDAY=95;BYWEEKNO=1;BYDAY=MO,+2FR;BYHOUR=9;BYMINUTE=30;BYSECOND=0;BYEASTER=-1"
	r, _ := parseRecurrence(str)
	want := "DTSTART:20120201T093000Z\nRRULE:FREQ=WEEKLY;INTERVAL=5;WKST=TU;COUNT=2;UNTIL=20130130T230000Z;BYSETPOS=2;BYMONTH=3;BYYEARDAY=95;BYWEEKNO=1;BYDAY=MO,+2FR;BYHOUR=9;BYMINUTE=30;BYSECOND=0;BYEASTER=-1"
	if s := r.String(); s != want {
		t.Errorf("parseRecurrence(%q).String() = %q, want %q", str, s, want)
	}
	r, _ = parseRecurrence(want)
	if s := r.String(); s != want {
		t.Errorf("parseRecurrence(%q).String() = %q, want %q", want, want, want)
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
		if _, e := parseRecurrence(item); e == nil {
			t.Errorf("parseRecurrence(%q) = nil, want error", item)
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

	var rruleLines []string
	for _, line := range set.Strings() {
		if strings.HasPrefix(line, "DTSTART") || strings.HasPrefix(line, "RRULE") {
			rruleLines = append(rruleLines, line)
		}
	}
	rruleStr := strings.Join(rruleLines, "\n")
	if rruleStr != "DTSTART;TZID=America/New_York:20180101T090000\nRRULE:FREQ=WEEKLY;INTERVAL=2;BYDAY=MO,TU" {
		t.Errorf("Unexpected rrule: %s", rruleStr)
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
	setStr := set.String()
	setFromSetStr, _ := StrToRRuleSet(setStr)

	if setStr != setFromSetStr.String() {
		t.Errorf("Expected string output\n %s \nbut got\n %s\n", setStr, setFromSetStr.String())
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

		sRRule := s.String()

		if sRRule != expected {
			t.Errorf("DTSTART output not valid. Expected: \n%s \n Got: \n%s", expected, sRRule)
		}
	})
	t.Run("DtstartTimeZoneValidOutputForExDate", func(t *testing.T) {
		input := []string{
			"DTSTART;TZID=Europe/Moscow:20180220T090000",
			"EXDATE;VALUE=DATE-TIME:20180223T100000",
		}
		expected := "DTSTART;TZID=Europe/Moscow:20180220T090000\nEXDATE;TZID=Europe/Moscow:20180223T100000"
		s, err := StrSliceToRRuleSet(input)
		if err != nil {
			t.Error(err)
		}
		sRRule := s.String()
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

		sRRule := s.String()

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

	r := New(ROption{Freq: MONTHLY, Dtstart: dtStart})
	if r == nil {
		t.Fatal("failed to create recurrence")
	}
	want := "DTSTART;TZID=America/New_York:20180101T090000\nRRULE:FREQ=MONTHLY"
	if r.String() != want {
		t.Errorf("Expected RFC string %s, got %v", want, r.String())
	}

	expectedSetStr := "DTSTART;TZID=America/New_York:20180101T090000\nRRULE:FREQ=MONTHLY"

	set := New(ROption{Freq: MONTHLY, Dtstart: dtStart})
	if set.String() != expectedSetStr {
		t.Errorf("Expected RFC Set string %s, got %s", expectedSetStr, set.String())
	}
}

func TestRFCRuleToStr(t *testing.T) {
	nyLoc, _ := time.LoadLocation("America/New_York")
	dtStart := time.Date(2018, 1, 1, 9, 0, 0, 0, nyLoc)

	r := New(ROption{Freq: MONTHLY, Dtstart: dtStart})
	if r == nil {
		t.Fatal("failed to create recurrence")
	}
	want := "DTSTART;TZID=America/New_York:20180101T090000\nRRULE:FREQ=MONTHLY"
	if r.String() != want {
		t.Errorf("Expected RFC string %s, got %v", want, r.String())
	}
}

// TestRRuleStringAllDayUntil tests RRuleString() handling of UNTIL for all-day events.
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
			// All-day event with an UNTIL parameter.
			option := ROption{
				Freq:    DAILY,
				AllDay:  true,
				Dtstart: time.Date(2023, 1, 1, 14, 30, 0, 0, tc.tz),
				Until:   time.Date(2023, 1, 5, 16, 45, 30, 0, tc.tz),
			}

			output := rruleStringFromOption(&option)
			t.Logf("Timezone %s RRuleString output: %s", tc.name, output)

			// Verify the UNTIL parameter is handled correctly.
			if !strings.Contains(output, "UNTIL=") {
				t.Errorf("Expected UNTIL parameter in output for timezone %s, got: %s", tc.name, output)
			}

			// Verify UNTIL uses DATE format (no time part).
			if !strings.Contains(output, "UNTIL=20230105") {
				t.Errorf("Expected UNTIL=20230105 in output for timezone %s, got: %s", tc.name, output)
			}

			// Verify UNTIL has no time part (no "T" time).
			if strings.Contains(output, "UNTIL=20230105T") {
				t.Errorf("UNTIL should use DATE format (no time part) for all-day events in timezone %s, got: %s", tc.name, output)
			}
		})
	}
}

// TestRRuleStringAllDayConsistency tests RRuleString consistency for all-day events.
func TestRRuleStringAllDayConsistency(t *testing.T) {
	// Create all-day events on the same date across timezones.
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
		output := rruleStringFromOption(&option)
		outputs = append(outputs, output)
		t.Logf("Option %d RRuleString output: %s", i+1, output)
	}

	// Verify outputs match since RRuleString omits DTSTART and includes only RRULE.
	expected := "FREQ=WEEKLY;COUNT=3"
	for i, output := range outputs {
		if output != expected {
			t.Errorf("Option %d: expected %s, got %s", i+1, expected, output)
		}
	}
}

// TestRRuleStringAllDayWithUntilTimezone tests UNTIL handling across timezones for all-day events.
func TestRRuleStringAllDayWithUntilTimezone(t *testing.T) {
	// Create all-day events on the same date across timezones with UNTIL.
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
		output := rruleStringFromOption(&option)
		outputs = append(outputs, output)
		t.Logf("Option %d RRuleString output: %s", i+1, output)
	}

	// Verify UNTIL uses DATE format.
	for i, output := range outputs {
		if !strings.Contains(output, "UNTIL=20230103") {
			t.Errorf("Option %d: expected UNTIL=20230103 in output, got: %s", i+1, output)
		}
		// Verify UNTIL has no time part (no "T" time).
		if strings.Contains(output, "UNTIL=20230103T") {
			t.Errorf("Option %d: UNTIL should use DATE format (no time part) for all-day events, got: %s", i+1, output)
		}
	}
}

// Negative COUNT should not appear in the generated RRULE string.
func TestROption_RRuleString_OmitsNegativeCount(t *testing.T) {
	opt := ROption{Freq: DAILY, Count: -5}
	got := rruleStringFromOption(&opt)
	want := "FREQ=DAILY"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

// When parsing an RRULE with COUNT=-1, String() should omit COUNT.
func TestParseRecurrence_String_OmitsNegativeCount_NoDTSTART(t *testing.T) {
	r, err := parseRecurrence("RRULE:FREQ=DAILY;COUNT=-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := r.String()
	want := "RRULE:FREQ=DAILY"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

// Same as above but with DTSTART present. COUNT should be omitted from output.
func TestParseRecurrence_String_OmitsNegativeCount_WithDTSTART(t *testing.T) {
	dt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	in := "DTSTART:" + dt.Format(DateTimeFormat) + "\nRRULE:FREQ=DAILY;COUNT=-1"

	r, err := parseRecurrence(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := r.String()
	want := "DTSTART:" + dt.Format(DateTimeFormat) + "\nRRULE:FREQ=DAILY"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

// All-day RRULE output should omit BYHOUR/BYMINUTE/BYSECOND, even if provided.
func TestROption_RRuleString_AllDay_OmitsTimeParts(t *testing.T) {
	opt := ROption{
		Freq:     DAILY,
		AllDay:   true,
		Dtstart:  time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
		Byhour:   []int{9, 15},
		Byminute: []int{30},
		Bysecond: []int{0},
	}
	got := rruleStringFromOption(&opt)
	if strings.Contains(got, "BYHOUR=") || strings.Contains(got, "BYMINUTE=") || strings.Contains(got, "BYSECOND=") {
		t.Fatalf("all-day RRULE must not include time parts, got: %q", got)
	}
}

// Set serialization should also omit time parts for all-day rules.
func TestSet_String_AllDay_OmitsTimePartsInRRULE(t *testing.T) {
	lines := []string{
		"DTSTART;VALUE=DATE:20230101",
		"RRULE:FREQ=DAILY;BYHOUR=9;BYMINUTE=30;BYSECOND=0",
	}
	set, err := StrSliceToRRuleSet(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := set.String()
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(line, "RRULE:") {
			if strings.Contains(line, "BYHOUR=") || strings.Contains(line, "BYMINUTE=") || strings.Contains(line, "BYSECOND=") {
				t.Fatalf("all-day Set RRULE must not include time parts, got: %q", line)
			}
		}
	}
}

func TestSet_String_OmitsIgnoredParams_AllDay(t *testing.T) {
	lines := []string{
		"DTSTART;VALUE=DATE:20240101",
		"RRULE:FREQ=DAILY;COUNT=-1;INTERVAL=0;WKST=MO;BYHOUR=9;BYMINUTE=30;BYSECOND=0",
	}
	set, err := StrSliceToRRuleSet(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := set.String()
	rline := getRRULELine(out)
	if rline == "" {
		t.Fatalf("RRULE line not found in: %q", out)
	}
	// All-day: drop COUNT<=0, INTERVAL=0, default WKST=MO, and time parts
	if want := "RRULE:FREQ=DAILY"; rline != want {
		t.Fatalf("expected %q, got %q", want, rline)
	}
}

func TestSet_String_OmitsIgnoredParams_Timed(t *testing.T) {
	lines := []string{
		"DTSTART:20240101T100000Z",
		"RRULE:FREQ=DAILY;COUNT=-1;INTERVAL=0;WKST=MO;BYHOUR=9;BYMINUTE=30;BYSECOND=0",
	}
	set, err := StrSliceToRRuleSet(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := set.String()
	rline := getRRULELine(out)
	if rline == "" {
		t.Fatalf("RRULE line not found in: %q", out)
	}
	// Timed: drop COUNT<=0, INTERVAL=0, default WKST=MO, keep time parts
	if want := "RRULE:FREQ=DAILY;BYHOUR=9;BYMINUTE=30;BYSECOND=0"; rline != want {
		t.Fatalf("expected %q, got %q", want, rline)
	}
}

func TestSet_String_OmitsNegativeCount_KeepsUntil(t *testing.T) {
	lines := []string{
		"DTSTART:20240101T000000Z",
		"RRULE:FREQ=DAILY;COUNT=-1;UNTIL=20240103T000000Z",
	}
	set, err := StrSliceToRRuleSet(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rline := getRRULELine(set.String())
	if want := "RRULE:FREQ=DAILY;UNTIL=20240103T000000Z"; rline != want {
		t.Fatalf("expected %q, got %q", want, rline)
	}
}

// Confirms current behavior: DTSTART not first is ignored (non-fatal) per current parser.
func TestSetParse_DTSTARTNotFirst_Ignored(t *testing.T) {
	lines := []string{
		"RRULE:FREQ=DAILY;COUNT=1",
		"DTSTART:20250101T000000Z",
	}
	set, err := StrSliceToRRuleSet(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !set.GetDTStart().IsZero() {
		t.Errorf("expected dtstart to remain zero when DTSTART is not first, got %v", set.GetDTStart())
	}
}

// Ensures DTSTART;VALUE=DATE auto-detects all-day in StrToRRuleSet and propagates to RRULE
func TestStrToRRuleSet_AllDayDetectionAndPropagation(t *testing.T) {
	setStr := strings.Join([]string{
		"DTSTART;VALUE=DATE:20230901",
		"RRULE:FREQ=DAILY;COUNT=3;UNTIL=20230903",
		"RDATE;VALUE=DATE:20230902",
		"EXDATE;VALUE=DATE:20230901",
	}, "\n")

	set, err := StrToRRuleSet(setStr)
	if err != nil {
		t.Fatalf("StrToRRuleSet failed: %v", err)
	}

	// All-day should be auto-detected
	if !set.IsAllDay() {
		t.Fatal("expected set to be all-day after parsing DTSTART;VALUE=DATE")
	}

	// RRULE serialization should use DATE UNTIL and DTSTART should use VALUE=DATE
	out := set.String()
	if !strings.Contains(out, "DTSTART;VALUE=DATE:20230901") {
		t.Errorf("expected DTSTART;VALUE=DATE in output, got: %s", out)
	}
	if !strings.Contains(out, "RRULE:FREQ=DAILY;COUNT=3;UNTIL=20230903") {
		t.Errorf("expected UNTIL as DATE in RRULE, got: %s", out)
	}
	if strings.Contains(out, "UNTIL=20230903T") {
		t.Errorf("UNTIL must not include time part for all-day, got: %s", out)
	}

	// Verify occurrences: DTSTART=2023-09-01, COUNT=3 -> 1st, 2nd, 3rd; EXDATE removes 1st; RDATE adds 2nd (duplicate)
	want := []time.Time{
		time.Date(2023, 9, 2, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 9, 3, 0, 0, 0, 0, time.UTC),
	}
	got := set.All()
	if !timesEqual(got, want) {
		t.Errorf("occurrences mismatch, got %v, want %v", got, want)
	}
}

// UNTIL zero-value should not appear in RRULE output.
func TestROption_RRuleString_OmitsZeroUntil(t *testing.T) {
	opt := ROption{Freq: DAILY, Until: time.Time{}}
	got := rruleStringFromOption(&opt)
	want := "FREQ=DAILY"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

// INTERVAL=0 should be treated as default and omitted from RRULE output on round-trip.
func TestParseRecurrence_String_OmitsIntervalZero(t *testing.T) {
	r, err := parseRecurrence("RRULE:FREQ=DAILY;INTERVAL=0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := r.String()
	want := "RRULE:FREQ=DAILY"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

// WKST=MO (default) should be omitted from RRULE output on round-trip.
func TestParseRecurrence_String_OmitsDefaultWkst(t *testing.T) {
	r, err := parseRecurrence("RRULE:FREQ=WEEKLY;WKST=MO")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := r.String()
	want := "RRULE:FREQ=WEEKLY"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

// When COUNT is invalid (negative) but UNTIL is valid, omit COUNT and keep UNTIL.
func TestParseRecurrence_String_OmitsNegativeCount_KeepsUntil(t *testing.T) {
	until := time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)
	in := "RRULE:FREQ=DAILY;COUNT=-1;UNTIL=" + until.Format(DateTimeFormat)
	r, err := parseRecurrence(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := r.String()
	want := "RRULE:FREQ=DAILY;UNTIL=" + until.Format(DateTimeFormat)
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestParseRecurrence_String_OmitsCountZero(t *testing.T) {
	r, err := parseRecurrence("RRULE:FREQ=DAILY;COUNT=0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	option := r.ruleOptionFromState()
	if got, want := rruleStringFromOption(&option), "FREQ=DAILY"; got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestSet_String_OmitsCountZero(t *testing.T) {
	lines := []string{
		"DTSTART:20240101T000000Z",
		"RRULE:FREQ=DAILY;COUNT=0",
	}
	set, err := StrSliceToRRuleSet(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := set.String()
	if out == "" {
		t.Fatal("unexpected empty set output")
	}
	if idx := indexOfRRULE(out); idx == -1 {
		t.Fatalf("RRULE line not found in: %q", out)
	} else {
		rline := rrLineAt(out, idx)
		if rline != "RRULE:FREQ=DAILY" {
			t.Fatalf("expected RRULE to omit COUNT=0, got: %q", rline)
		}
	}
}

func TestSet_String_ExcludeDTSTARTWhenDisabled(t *testing.T) {
	set := New(ROption{
		Freq:    DAILY,
		Dtstart: time.Date(2024, 2, 1, 9, 0, 0, 0, time.UTC),
	})
	if set == nil {
		t.Fatal("failed to create recurrence")
	}
	got := set.RRuleString()
	want := "RRULE:FREQ=DAILY"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestSetParse_CommaSeparatedDates_StringSplitsLines(t *testing.T) {
	lines := []string{
		"DTSTART:20240101T090000Z",
		"RRULE:FREQ=DAILY;COUNT=3",
		"RDATE:20240104T090000Z,20240105T090000Z",
		"EXDATE:20240102T090000Z,20240103T090000Z",
	}
	set, err := StrSliceToRRuleSet(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := set.String()
	rdateCount := 0
	exdateCount := 0
	for _, line := range splitLines(out) {
		if strings.HasPrefix(line, "RDATE") {
			rdateCount++
			if strings.Contains(line, ",") {
				t.Fatalf("unexpected comma-separated RDATE output: %q", line)
			}
		}
		if strings.HasPrefix(line, "EXDATE") {
			exdateCount++
			if strings.Contains(line, ",") {
				t.Fatalf("unexpected comma-separated EXDATE output: %q", line)
			}
		}
	}
	if rdateCount != 2 {
		t.Fatalf("expected 2 RDATE lines, got %d in %q", rdateCount, out)
	}
	if exdateCount != 2 {
		t.Fatalf("expected 2 EXDATE lines, got %d in %q", exdateCount, out)
	}
}

// small helpers for the test file
func getRRULELine(s string) string {
	for _, line := range splitLines(s) {
		if strings.HasPrefix(line, "RRULE:") {
			return line
		}
	}
	return ""
}

func indexOfRRULE(s string) int {
	i := 0
	for _, line := range splitLines(s) {
		if len(line) >= 6 && line[:6] == "RRULE:" {
			return i
		}
		i++
	}
	return -1
}

func rrLineAt(s string, idx int) string {
	i := 0
	for _, line := range splitLines(s) {
		if i == idx {
			return line
		}
		i++
	}
	return ""
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	lines = append(lines, s[start:])
	return lines
}
