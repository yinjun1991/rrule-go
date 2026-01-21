# RFC 5545 VEVENT and VTODO Reference

## Overview

This document is based on RFC 5545. It explains the rules, examples, and differences between VEVENT (events) and VTODO (tasks) as they relate to recurrence rules (RRULE).

## 1. Baseline comparison

### 1.1 Time property comparison

| Property | VEVENT | VTODO | Notes |
|------|--------|-------|------|
| DTSTART | Required (unless METHOD is specified) | Optional | Start time |
| DTEND | Optional | Not supported | Event end time (exclusive) |
| DUE | Not supported | Optional | Task due time (inclusive) |
| DURATION | Optional | Optional | Duration |

### 1.2 Time property constraints

**VEVENT constraints:**
- DTEND and DURATION cannot coexist
- If there is no DTEND and no DURATION, the event ends at DTSTART
- All-day events: if DTSTART is DATE, DTEND must also be DATE

**VTODO constraints:**
- DUE and DURATION cannot coexist
- If DURATION exists, DTSTART is required
- Allowed combinations: only DUE, only DTSTART, DTSTART + DUE, DTSTART + DURATION
- Time properties may be omitted entirely (ongoing tasks)

## 2. All-day event/task time semantics

### 2.1 All-day event (VEVENT)

For an all-day event on 2025-09-10:

```ics
BEGIN:VEVENT
DTSTART;VALUE=DATE:20250910
DTEND;VALUE=DATE:20250911
SUMMARY:All-day event example
END:VEVENT
```

**Key points:**
- DTSTART: 2025-09-10 (inclusive)
- DTEND: 2025-09-11 (exclusive, the event ends on 2025-09-10)
- DTEND is exclusive, so add 1 day

### 2.2 All-day task (VTODO)

For an all-day task on 2025-09-10:

```ics
BEGIN:VTODO
DTSTART;VALUE=DATE:20250910
DUE;VALUE=DATE:20250910
SUMMARY:All-day task example
END:VTODO
```

**Key points:**
- DTSTART: 2025-09-10 (inclusive)
- DUE: 2025-09-10 (inclusive, due on 2025-09-10)
- DUE is inclusive, so no +1 day

### 2.3 Summary

| Type | All day on 2025-09-10 | End representation | Notes |
|------|-----------------|--------------|------|
| VEVENT | DTEND=20250911 | Exclusive | Add 1 day to represent end |
| VTODO | DUE=20250910 | Inclusive | Use the same day as the due date |

## 3. RRULE recurrence rules

### 3.1 Basic rules

RRULE can appear in VEVENT, VTODO, and VJOURNAL components:

```
RRULE:FREQ=DAILY;COUNT=10
RRULE:FREQ=WEEKLY;UNTIL=20251224T000000Z
RRULE:FREQ=MONTHLY;INTERVAL=2;BYDAY=1MO
```

### 3.2 Importance of DTSTART

**Core principles:**
- RRULE calculations are based on DTSTART
- DTSTART defines the first instance in the recurrence set
- Missing RRULE fields are derived from DTSTART

**Example:**
```ics
DTSTART:20250101T090000
RRULE:FREQ=WEEKLY
# Result: every Monday at 9:00 AM (weekday and time derived from DTSTART)
```

### 3.3 VTODO without DTSTART

**Spec note:**
Per RFC 5545, VTODO may omit DTSTART and DUE. In that case:

1. **RRULE is not supported**: without DTSTART there is no base time
2. **Special behavior**: VTODO is associated with each consecutive calendar date until completion
3. **Implementation differences**: calendar apps may handle this differently

**Example:**
```ics
BEGIN:VTODO
UID:todo-without-time@example.com
SUMMARY:Ongoing task
DESCRIPTION:Task with no time constraints
STATUS:NEEDS-ACTION
END:VTODO
```

## 4. UNTIL behavior

### 4.1 General rules for UNTIL

UNTIL behaves the same in VEVENT and VTODO:

```
RRULE:FREQ=DAILY;UNTIL=20251224T000000Z
```

**Key points:**
- UNTIL specifies the inclusive end of the recurrence
- It must match the DTSTART value type (DATE or DATE-TIME)
- If DTSTART is local time, UNTIL should be in UTC

### 4.2 UNTIL for all-day events/tasks

**All-day event:**
```ics
DTSTART;VALUE=DATE:20250101
RRULE:FREQ=DAILY;UNTIL=20251231
```

**All-day task:**
```ics
DTSTART;VALUE=DATE:20250101
RRULE:FREQ=DAILY;UNTIL=20251231
```

**Note:**
- All-day UNTIL uses DATE format
- DTEND/DUE differences do not affect UNTIL; it only controls recurrence generation

## 5. Complete examples

### 5.1 Daily meeting (VEVENT)

```ics
BEGIN:VEVENT
UID:daily-meeting@example.com
DTSTART;TZID=Asia/Shanghai:20250101T090000
DTEND;TZID=Asia/Shanghai:20250101T100000
RRULE:FREQ=DAILY;COUNT=30
SUMMARY:Daily standup
END:VEVENT
```

### 5.2 Daily task (VTODO)

```ics
BEGIN:VTODO
UID:daily-task@example.com
DTSTART;TZID=Asia/Shanghai:20250101T090000
DUE;TZID=Asia/Shanghai:20250101T180000
RRULE:FREQ=DAILY;COUNT=30
SUMMARY:Daily work task
END:VTODO
```

### 5.3 Complex recurrence example

**Task on the last weekday of each month:**
```ics
BEGIN:VTODO
UID:monthly-report@example.com
DTSTART;VALUE=DATE:20250131
DUE;VALUE=DATE:20250131
RRULE:FREQ=MONTHLY;BYDAY=-1MO,-1TU,-1WE,-1TH,-1FR;BYSETPOS=-1
SUMMARY:Monthly report
END:VTODO
```

## 6. Implementation notes

### 6.1 Timezone handling

1. **Local time with timezone**: recommended for recurring events/tasks
2. **UTC time**: for cross-timezone scenarios
3. **Floating time**: for all-day events/tasks

### 6.2 All-day handling

1. **VEVENT**: DTEND = DTSTART + 1 day
2. **VTODO**: DUE = due date (no +1 day)
3. **Time format**: use VALUE=DATE

### 6.3 Recurrence calculation

1. **Base time**: always DTSTART
2. **Time derivation**: missing time fields are derived from DTSTART
3. **Exceptions**: use EXDATE to exclude specific instances

## 7. Summary of differences

| Aspect | VEVENT | VTODO |
|------|--------|-------|
| End time property | DTEND (exclusive) | DUE (inclusive) |
| All-day end time | Add 1 day | Same day |
| Time properties required | DTSTART required | All time properties optional |
| Support without time properties | Not supported | Supported (ongoing task) |
| RRULE support | Requires DTSTART | Requires DTSTART |
| Recurrence base | DTSTART | DTSTART (if present) |

## 8. Best practices

1. **Time properties**: prefer DTEND/DUE over DURATION
2. **Timezone consistency**: keep DTSTART, DTEND/DUE, and UNTIL in the same timezone
3. **All-day handling**: account for DTEND vs DUE semantics
4. **Recurrence validation**: keep DTSTART aligned with RRULE
5. **Exceptions**: use EXDATE and RDATE appropriately

## 9. Recurrence storage guidance

### 9.1 Common vendor approach

Based on research into Google Calendar, Apple Calendar, and Microsoft Outlook, the following structure is recommended:

#### Google Calendar API style (recommended)
```go
type Event struct {
    // Time fields stored separately.
    DTStart   time.Time `json:"dtstart"`
    DTEnd     time.Time `json:"dtend,omitempty"`
    DUE       time.Time `json:"due,omitempty"`

    // Recurrence rules as a []string array.
    Recurrence []string `json:"recurrence,omitempty"`
    // Example:
    // ["RRULE:FREQ=WEEKLY;BYDAY=MO,WE,FR", "EXDATE:20250115T100000Z", "RDATE:20250120T100000Z"]
}
```

### 9.2 Recurrence array composition

**Correct grouping:**
1. `RRULE:...` - recurrence rule (including COUNT/UNTIL)
2. `RDATE:...` - additional dates
3. `EXDATE:...` - excluded dates
4. **Store time fields separately** (not in the recurrence array)

### 9.3 Design principles

1. **Follow RFC 5545**: time properties and recurrence rules are separate in iCalendar
2. **Align with major vendors**: Google/Apple/Microsoft separate time fields from recurrence rules
3. **Ease of processing**:
   - time fields define the base event
   - recurrence rules generate repeated instances
   - exceptions remain explicit

### 9.4 Implementation example

```go
// Recurrence rule parsing
func ParseRecurrence(recurrence []string) (*RecurrenceRule, error) {
    var rrule *RRule
    var rdates []time.Time
    var exdates []time.Time

    for _, line := range recurrence {
        switch {
        case strings.HasPrefix(line, "RRULE:"):
            rrule = parseRRule(line[6:])
        case strings.HasPrefix(line, "RDATE:"):
            rdates = append(rdates, parseRDate(line[6:]))
        case strings.HasPrefix(line, "EXDATE:"):
            exdates = append(exdates, parseExDate(line[7:]))
        }
    }

    return &RecurrenceRule{
        RRule:   rrule,
        RDates:  rdates,
        ExDates: exdates,
    }, nil
}
```

### 9.5 Key takeaways

1. **Separate time fields**: keep `DTSTART`/`DTEND`/`DUE` out of the `recurrence` array
2. **Follow RFC 5545**: each string is a full iCalendar property line
3. **Extensible**: additional recurrence properties can be added later
4. **Compatibility**: matches Google Calendar API formatting

---

*This document is based on RFC 5545 and is intended as a technical reference for calendar application development.*
