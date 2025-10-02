package rrule

import (
    "testing"
    "time"
)

func TestStrToDtStart_InvalidTZID(t *testing.T) {
    // Missing zone name
    if _, err := StrToDtStart("TZID=:", time.UTC); err == nil {
        t.Error("expected error for bad TZID parameter format")
    }
    // Unknown zone
    if _, err := StrToDtStart("TZID=Bad/Zone:20250101T010203", time.UTC); err == nil {
        t.Error("expected error for unknown TZID zone")
    }
}

func TestStrToDates_UnsupportedValueParam(t *testing.T) {
    // VALUE=PERIOD is not supported by StrToDatesInLoc
    if _, err := StrToDates("VALUE=PERIOD:20180223T100000Z/20180223T120000Z"); err == nil {
        t.Error("expected error for unsupported VALUE=PERIOD")
    }
}
