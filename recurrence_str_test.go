package rrule

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func parseRecurrence(input string) (*Recurrence, error) {
	return ParseRRuleString(input)
}

func rruleFromOption(t *testing.T, option ROption) string {
	r, err := New(option)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return strings.TrimPrefix(r.RRuleString(), "RRULE:")
}

func TestCompatibility(t *testing.T) {
	str := "FREQ=WEEKLY;DTSTART=20120201T093000Z;INTERVAL=5;WKST=TU;COUNT=2;UNTIL=20130130T230000Z;BYSETPOS=2;BYMONTH=3;BYYEARDAY=95;BYWEEKNO=1;BYDAY=MO,+2FR;BYHOUR=9;BYMINUTE=30;BYSECOND=0;BYEASTER=-1"
	r, _ := parseRecurrence(str)
	want := "DTSTART:20120201T093000Z\nRRULE:FREQ=WEEKLY;INTERVAL=5;WKST=TU;COUNT=2;UNTIL=20130130T230000Z;BYSETPOS=2;BYMONTH=3;BYYEARDAY=95;BYWEEKNO=1;BYDAY=MO,FR;BYHOUR=9;BYMINUTE=30;BYSECOND=0;BYEASTER=-1"
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
		if _, err := Parse(ss...); err == nil {
			t.Error("Expected parse error for rules: ", ss)
		}
	}
}

func TestStrSetEmptySliceParse(t *testing.T) {
	s, err := Parse()
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
		s, err := Parse(input...)
		if err != nil {
			t.Error(err)
		}
		d := s.GetRDate()[0]
		if !d.Equal(time.Date(2018, 02, 23, 0, 0, 0, 0, time.UTC)) {
			t.Error("Bad time parsed: ", d)
		}
	})

	t.Run("IgnoreTimezoneForValueDate", func(t *testing.T) {
		input := []string{
			"RDATE;VALUE=DATE;TZID=America/Denver:20180223",
		}
		s, err := Parse(input...)
		if err != nil {
			t.Error(err)
		}
		d := s.GetRDate()[0]
		if !d.Equal(time.Date(2018, 02, 23, 0, 0, 0, 0, time.UTC)) {
			t.Error("Bad time parsed: ", d)
		}
	})
}

func TestSetStrCompatibility(t *testing.T) {
	badInputStrs := []string{
		"",
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
	if rruleStr != "DTSTART;TZID=America/New_York:20180101T090000\nRRULE:FREQ=DAILY;UNTIL=20180517T235959Z" {
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
		s, err := Parse(input...)
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
		s, err := Parse(input...)
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
		s, err := Parse(input...)
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
		s, err := Parse(input...)
		if err != nil {
			t.Error(err)
		}

		sRRule := s.String()

		if sRRule != expected {
			t.Errorf("DTSTART output not valid. Expected: \n%s \n Got: \n%s", expected, sRRule)
		}
	})

	t.Run("DefaultZoneIsUTC", func(t *testing.T) {
		input := []string{
			"RDATE;VALUE=DATE-TIME:20180223T100000",
		}
		s, err := Parse(input...)
		if err != nil {
			t.Error(err)
		}
		d := s.GetRDate()[0]
		if !d.Equal(time.Date(2018, 02, 23, 10, 0, 0, 0, time.UTC)) {
			t.Error("Bad time parsed: ", d)
		}
	})
}

func TestRFCSetToString(t *testing.T) {
	nyLoc, _ := time.LoadLocation("America/New_York")
	dtStart := time.Date(2018, 1, 1, 9, 0, 0, 0, nyLoc)

	r, err := New(ROption{Freq: MONTHLY, Dtstart: dtStart})
	if err != nil {
		t.Fatal(err)
	}
	want := "DTSTART;TZID=America/New_York:20180101T090000\nRRULE:FREQ=MONTHLY"
	if r.String() != want {
		t.Errorf("Expected RFC string %s, got %v", want, r.String())
	}

	expectedSetStr := "DTSTART;TZID=America/New_York:20180101T090000\nRRULE:FREQ=MONTHLY"

	set, err := New(ROption{Freq: MONTHLY, Dtstart: dtStart})
	if err != nil {
		t.Fatal(err)
	}
	if set.String() != expectedSetStr {
		t.Errorf("Expected RFC Set string %s, got %s", expectedSetStr, set.String())
	}
}

func TestRFCRuleToStr(t *testing.T) {
	nyLoc, _ := time.LoadLocation("America/New_York")
	dtStart := time.Date(2018, 1, 1, 9, 0, 0, 0, nyLoc)

	r, err := New(ROption{Freq: MONTHLY, Dtstart: dtStart})
	if err != nil {
		t.Fatal(err)
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

			output := rruleFromOption(t, option)
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
		output := rruleFromOption(t, option)
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
		output := rruleFromOption(t, option)
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
	got := rruleFromOption(t, opt)
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
	if !strings.Contains(got, "RRULE:FREQ=DAILY") {
		t.Fatalf("expected RRULE:FREQ=DAILY, got %q", got)
	}
	if strings.Contains(got, "COUNT=") {
		t.Fatalf("expected COUNT to be omitted, got %q", got)
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
	got := rruleFromOption(t, opt)
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
	set, err := Parse(lines...)
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
	set, err := Parse(lines...)
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
	set, err := Parse(lines...)
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
	set, err := Parse(lines...)
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
	set, err := Parse(lines...)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ignored := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	if set.GetDTStart().Equal(ignored) {
		t.Errorf("expected DTSTART to be ignored when not first, got %v", set.GetDTStart())
	}
}

func TestNewWithDTStart_IgnoresDTSTARTInLines(t *testing.T) {
	dtstart := time.Date(2024, 1, 2, 9, 30, 0, 0, time.UTC)
	lines := []string{
		"DTSTART:20240101T000000Z",
		"RRULE:FREQ=DAILY;COUNT=2",
	}
	set, err := NewWithDTStart(dtstart, false, lines...)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !set.GetDTStart().Equal(dtstart) {
		t.Fatalf("expected dtstart %v, got %v", dtstart, set.GetDTStart())
	}
	if got := getRRULELine(set.String()); got != "RRULE:FREQ=DAILY;COUNT=2" {
		t.Fatalf("unexpected RRULE line: %s", got)
	}
}

func TestNewWithDTStart_AllDayUsesProvidedDate(t *testing.T) {
	dtstart := time.Date(2024, 3, 5, 18, 45, 0, 0, time.FixedZone("JST", 9*3600))
	lines := []string{
		"DTSTART:20240101T000000Z",
		"RRULE:FREQ=DAILY;BYHOUR=9;BYMINUTE=30;BYSECOND=0",
	}
	set, err := NewWithDTStart(dtstart, true, lines...)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wantStart := "DTSTART;VALUE=DATE:20240305"
	if out := set.DTStartString(); out != wantStart {
		t.Fatalf("expected %q, got %q", wantStart, out)
	}
	if rrule := getRRULELine(set.String()); rrule != "RRULE:FREQ=DAILY" {
		t.Fatalf("expected RRULE without time parts, got %q", rrule)
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
	got := rruleFromOption(t, opt)
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
	if !strings.Contains(got, "RRULE:FREQ=DAILY") {
		t.Fatalf("expected RRULE:FREQ=DAILY, got %q", got)
	}
	if strings.Contains(got, "INTERVAL=") {
		t.Fatalf("expected INTERVAL to be omitted, got %q", got)
	}
}

// WKST=MO (default) should be omitted from RRULE output on round-trip.
func TestParseRecurrence_String_OmitsDefaultWkst(t *testing.T) {
	r, err := parseRecurrence("RRULE:FREQ=WEEKLY;WKST=MO")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := r.String()
	if !strings.Contains(got, "RRULE:FREQ=WEEKLY") {
		t.Fatalf("expected RRULE:FREQ=WEEKLY, got %q", got)
	}
	if strings.Contains(got, "WKST=") {
		t.Fatalf("expected WKST to be omitted, got %q", got)
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
	expected := "RRULE:FREQ=DAILY;UNTIL=" + until.Format(DateTimeFormat)
	if !strings.Contains(got, expected) {
		t.Fatalf("expected %q, got %q", expected, got)
	}
	if strings.Contains(got, "COUNT=") {
		t.Fatalf("expected COUNT to be omitted, got %q", got)
	}
}

func TestParseRecurrence_String_OmitsCountZero(t *testing.T) {
	r, err := parseRecurrence("RRULE:FREQ=DAILY;COUNT=0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	option := r.ruleOptionFromState()
	if got, want := rruleFromOption(t, option), "FREQ=DAILY"; got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestSet_String_OmitsCountZero(t *testing.T) {
	lines := []string{
		"DTSTART:20240101T000000Z",
		"RRULE:FREQ=DAILY;COUNT=0",
	}
	set, err := Parse(lines...)
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
	set, err := New(ROption{
		Freq:    DAILY,
		Dtstart: time.Date(2024, 2, 1, 9, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatal(err)
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
	set, err := Parse(lines...)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := set.String()
	if !strings.Contains(out, "RDATE:20240104T090000Z,20240105T090000Z") {
		t.Fatalf("expected comma-separated RDATE output, got %q", out)
	}
	if !strings.Contains(out, "EXDATE:20240102T090000Z,20240103T090000Z") {
		t.Fatalf("expected comma-separated EXDATE output, got %q", out)
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

func TestNormalizeRecurrenceRuleset(t *testing.T) {
	tests := []struct {
		name        string
		input       []string
		expected    []string
		expectError bool
	}{
		{
			name:        "empty input",
			input:       []string{},
			expected:    nil,
			expectError: false,
		},
		{
			name:        "nil input",
			input:       nil,
			expected:    nil,
			expectError: false,
		},
		{
			name:     "already normalized RRULE",
			input:    []string{"RRULE:FREQ=DAILY;COUNT=5"},
			expected: []string{"RRULE:FREQ=DAILY;COUNT=5"},
		},
		{
			name:     "dtstart passthrough",
			input:    []string{"DTSTART:20240101T090000Z"},
			expected: []string{"DTSTART:20240101T090000Z"},
		},
		{
			name:     "missing RRULE prefix",
			input:    []string{"FREQ=DAILY;COUNT=5"},
			expected: []string{"RRULE:FREQ=DAILY;COUNT=5"},
		},
		{
			name:     "mixed case with missing prefix",
			input:    []string{"freq=daily;count=5"},
			expected: []string{"RRULE:freq=daily;count=5"},
		},
		{
			name: "multiple rules with mixed formats",
			input: []string{
				"RRULE:FREQ=DAILY;COUNT=5",
				"FREQ=WEEKLY;BYDAY=MO,WE,FR",
				"RDATE:20240115T100000Z",
				"EXDATE:20240120T100000Z",
			},
			expected: []string{
				"RRULE:FREQ=DAILY;COUNT=5",
				"RDATE:20240115T100000Z",
				"EXDATE:20240120T100000Z",
			},
		},
		{
			name: "multiple RRULE lines keeps first",
			input: []string{
				"RRULE:FREQ=DAILY;COUNT=5",
				"RRULE:FREQ=WEEKLY;BYDAY=MO",
				"RDATE:20240115T100000Z",
			},
			expected: []string{
				"RRULE:FREQ=DAILY;COUNT=5",
				"RDATE:20240115T100000Z",
			},
		},
		{
			name:     "RDATE without modification",
			input:    []string{"RDATE:20240115T100000Z"},
			expected: []string{"RDATE:20240115T100000Z"},
		},
		{
			name:     "EXDATE without modification",
			input:    []string{"EXDATE:20240120T100000Z"},
			expected: []string{"EXDATE:20240120T100000Z"},
		},
		{
			name:        "empty string in array",
			input:       []string{"FREQ=DAILY;COUNT=5", "", "RDATE:20240115T100000Z"},
			expected:    []string{"RRULE:FREQ=DAILY;COUNT=5", "RDATE:20240115T100000Z"},
			expectError: false,
		},
		{
			name:        "whitespace only string",
			input:       []string{"FREQ=DAILY;COUNT=5", "   ", "RDATE:20240115T100000Z"},
			expected:    []string{"RRULE:FREQ=DAILY;COUNT=5", "RDATE:20240115T100000Z"},
			expectError: false,
		},
		{
			name:        "invalid rule without FREQ",
			input:       []string{"COUNT=5"},
			expected:    nil,
			expectError: true,
		},
		{
			name:        "completely invalid format",
			input:       []string{"invalid rule format"},
			expected:    nil,
			expectError: true,
		},
		{
			name:        "empty RRULE content",
			input:       []string{"RRULE:"},
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NormalizeRecurrenceRuleset(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.expected, result)
			}
		})
	}
}

func TestNormalizeRecurrenceLine(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{
			name:     "already normalized RRULE",
			input:    "RRULE:FREQ=DAILY;COUNT=5",
			expected: "RRULE:FREQ=DAILY;COUNT=5",
		},
		{
			name:     "missing RRULE prefix",
			input:    "FREQ=DAILY;COUNT=5",
			expected: "RRULE:FREQ=DAILY;COUNT=5",
		},
		{
			name:     "RDATE rule",
			input:    "RDATE:20240115T100000Z",
			expected: "RDATE:20240115T100000Z",
		},
		{
			name:     "EXDATE rule",
			input:    "EXDATE:20240120T100000Z",
			expected: "EXDATE:20240120T100000Z",
		},
		{
			name:     "DTSTART rule passthrough",
			input:    "DTSTART:20240101T090000Z",
			expected: "DTSTART:20240101T090000Z",
		},
		{
			name:     "complex RRULE",
			input:    "FREQ=WEEKLY;BYDAY=MO,WE,FR;UNTIL=20241231T235959Z",
			expected: "RRULE:FREQ=WEEKLY;BYDAY=MO,WE,FR;UNTIL=20241231T235959Z",
		},
		{
			name:        "invalid rule without FREQ",
			input:       "COUNT=5;INTERVAL=2",
			expected:    "",
			expectError: true,
		},
		{
			name:        "empty string",
			input:       "",
			expected:    "",
			expectError: true,
		},
		{
			name:        "whitespace only",
			input:       "   ",
			expected:    "",
			expectError: true,
		},
		{
			name:        "empty RRULE content",
			input:       "RRULE:",
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := normalizeRecurrenceLine(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.expected, result)
			}
		})
	}
}

func TestValidateRRuleProperties(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "valid daily rule",
			input:       "FREQ=DAILY;COUNT=5",
			expectError: false,
		},
		{
			name:        "valid weekly rule",
			input:       "FREQ=WEEKLY;BYDAY=MO,WE,FR",
			expectError: false,
		},
		{
			name:        "valid monthly rule",
			input:       "FREQ=MONTHLY;BYMONTHDAY=15",
			expectError: false,
		},
		{
			name:        "missing FREQ parameter",
			input:       "COUNT=5;INTERVAL=2",
			expectError: true,
		},
		{
			name:        "empty content",
			input:       "",
			expectError: true,
		},
		{
			name:        "whitespace only",
			input:       "   ",
			expectError: true,
		},
		{
			name:        "case insensitive FREQ",
			input:       "freq=daily;count=5",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRRuleProperties(tt.input)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsRRuleProperties(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid RRULE content",
			input:    "FREQ=DAILY;COUNT=5",
			expected: true,
		},
		{
			name:     "valid complex RRULE content",
			input:    "FREQ=WEEKLY;BYDAY=MO,WE,FR;UNTIL=20241231T235959Z",
			expected: true,
		},
		{
			name:     "already has RRULE prefix",
			input:    "RRULE:FREQ=DAILY;COUNT=5",
			expected: false,
		},
		{
			name:     "RDATE format",
			input:    "RDATE:20240115T100000Z",
			expected: false,
		},
		{
			name:     "EXDATE format",
			input:    "EXDATE:20240120T100000Z",
			expected: false,
		},
		{
			name:     "missing FREQ",
			input:    "COUNT=5;INTERVAL=2",
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "case insensitive",
			input:    "freq=daily;count=5",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isRRuleProperties(tt.input)
			assert.EqualValues(t, tt.expected, result)
		})
	}
}

func TestRecurrenceRulesetNormalizationIntegration(t *testing.T) {
	testCases := []struct {
		name           string
		inputRules     []string
		expectedStored []string
	}{
		{
			name:           "missing RRULE prefix",
			inputRules:     []string{"FREQ=DAILY;COUNT=5"},
			expectedStored: []string{"RRULE:FREQ=DAILY;COUNT=5"},
		},
		{
			name:           "already normalized RRULE",
			inputRules:     []string{"RRULE:FREQ=WEEKLY;BYDAY=MO,WE,FR"},
			expectedStored: []string{"RRULE:FREQ=WEEKLY;BYDAY=MO,WE,FR"},
		},
		{
			name: "mixed RRULE and date exceptions",
			inputRules: []string{
				"FREQ=DAILY;COUNT=10",
				"RDATE:20240115T100000Z",
				"EXDATE:20240120T100000Z",
			},
			expectedStored: []string{
				"RRULE:FREQ=DAILY;COUNT=10",
				"RDATE:20240115T100000Z",
				"EXDATE:20240120T100000Z",
			},
		},
		{
			name: "complex RRULE parameters",
			inputRules: []string{
				"FREQ=MONTHLY;BYMONTHDAY=15;BYSETPOS=1,3;COUNT=12",
			},
			expectedStored: []string{
				"RRULE:FREQ=MONTHLY;BYMONTHDAY=15;BYSETPOS=1,3;COUNT=12",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			normalized, err := NormalizeRecurrenceRuleset(tc.inputRules)
			require.NoError(t, err, "Normalization should succeed")
			assert.EqualValues(t, tc.expectedStored, normalized, "Normalized rules should match expected")

			if len(normalized) > 0 {
				t.Logf("Normalized rules ready for parsing: %v", normalized)
			}
		})
	}
}
