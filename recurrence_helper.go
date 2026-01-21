package rrule

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

// RecurrenceSetHelper provides a wrapper around Set with full CRUD operations
// for RDATE and EXDATE, addressing the limitation that rrule-go only provides add methods
type RecurrenceSetHelper struct {
	set     *Set
	allDay  bool
	rdates  []time.Time // maintain RDATE list for delete operations
	exdates []time.Time // maintain EXDATE list for delete operations
}

// DTStartAlignmentError captures detailed context when DTSTART does not align with the recurrence rule.
type DTStartAlignmentError struct {
	DTStart              time.Time
	FirstOccurrence      time.Time
	LastOccurrenceBefore time.Time
	AllDay               bool
	reason               dtstartAlignmentReason
}

type dtstartAlignmentReason int

const (
	dtstartAlignmentReasonFirstMismatch dtstartAlignmentReason = iota
	dtstartAlignmentReasonNoOccurrenceOnOrAfter
	dtstartAlignmentReasonNoOccurrences
)

// Error implements the error interface with user-friendly messages while preserving RFC5545 terminology.
func (e *DTStartAlignmentError) Error() string {
	if e == nil {
		return "dtstart alignment error"
	}

	switch e.reason {
	case dtstartAlignmentReasonFirstMismatch:
		return fmt.Sprintf(
			"dtstart %s does not align with recurrence set: first occurrence on or after dtstart is %s",
			formatTimeForError(e.DTStart, e.AllDay),
			formatTimeForError(e.FirstOccurrence, e.AllDay),
		)
	case dtstartAlignmentReasonNoOccurrenceOnOrAfter:
		return fmt.Sprintf(
			"dtstart %s does not align with recurrence set: last occurrence before dtstart is %s and no occurrence exists on or after dtstart",
			formatTimeForError(e.DTStart, e.AllDay),
			formatTimeForError(e.LastOccurrenceBefore, e.AllDay),
		)
	default:
		return fmt.Sprintf(
			"dtstart %s does not align with recurrence set: recurrence set produces no occurrences on or after dtstart",
			formatTimeForError(e.DTStart, e.AllDay),
		)
	}
}

// AsDTStartAlignmentError extracts DTStartAlignmentError from an error chain.
func AsDTStartAlignmentError(err error) (*DTStartAlignmentError, bool) {
	if err == nil {
		return nil, false
	}
	var alignmentErr *DTStartAlignmentError
	if errors.As(err, &alignmentErr) {
		return alignmentErr, true
	}
	return nil, false
}

// ParseRecurrenceSet creates a RecurrenceSetHelper from RRule string array
// This is the main entry point for converting Event/Task.RRule to a manageable object
// dtstart is required for proper recurrence expansion
//
// Handles the following formats:
// - RRULE strings with or without "RRULE:" prefix (only first RRULE is used)
// - RDATE strings with "RDATE:" prefix
// - EXDATE strings with "EXDATE:" prefix
func ParseRecurrenceSet(
	ruleset []string,
	dtstart time.Time,
	allDay bool,
) (*RecurrenceSetHelper, error) {
	if len(ruleset) == 0 {
		return nil, fmt.Errorf("empty rrule strings")
	}

	// Process and normalize the input strings
	normalizedStrings, err := normalizeRecurrenceStrings(ruleset)
	if err != nil {
		return nil, fmt.Errorf("failed to normalize recurrence strings: %w", err)
	}

	// Use rrule-go's built-in parser with normalized strings
	set, err := StrSliceToRRuleSet(normalizedStrings)
	if err != nil {
		return nil, fmt.Errorf("failed to parse rrule strings: %w", err)
	}

	// Set DTStart - this is critical for proper recurrence expansion
	// Without DTStart, the set will use current time which may result in 0 occurrences
	if allDay {
		// For all-day events, normalize to date only (00:00:00 UTC)
		year, month, day := dtstart.Date()
		normalizedDTStart := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		set.DTStart(normalizedDTStart)
	} else {
		// For non-all-day events, use the provided dtstart
		set.DTStart(dtstart)
	}

	// Extract existing rdates and exdates for our internal tracking
	rdates := set.GetRDate()
	exdates := set.GetExDate()

	helper := &RecurrenceSetHelper{
		set:     set,
		allDay:  allDay,
		rdates:  make([]time.Time, len(rdates)),
		exdates: make([]time.Time, len(exdates)),
	}

	// Copy the slices to avoid reference issues
	copy(helper.rdates, rdates)
	copy(helper.exdates, exdates)

	return helper, nil
}

// ValidateDTStartAlignment ensures the provided DTSTART belongs to the recurrence rule.
// RFC 5545 mandates that DTSTART must represent the first instance produced by the rule.
func (h *RecurrenceSetHelper) ValidateDTStartAlignment(ignoreExdate bool) error {
	if h == nil || h.set == nil {
		return fmt.Errorf("recurrence set not initialized")
	}

	dtstart := h.DTStart()
	if dtstart.IsZero() {
		return fmt.Errorf("recurrence rule requires a non-zero DTSTART")
	}

	normalizedDTStart := h.normalizeTime(dtstart)

	// RFC 5545 requires DTSTART to be an actual occurrence produced by the recurrence set.
	// We therefore look at the first generated occurrence on/after the supplied DTSTART:
	//   • If the iterator returns the same instant we are aligned.
	//   • If not, a hand-crafted RDATE might still supply the first instance, so we check
	//     the stored RDATE list before rejecting the input.
	// By relying on the set-level iterator (instead of the first RRULE only) we correctly
	// honour irregular sequences such as "first instance via RDATE, rest via RRULE".
	firstOnOrAfter := h.set.After(dtstart, true)
	if !firstOnOrAfter.IsZero() {
		normalizedFirst := h.normalizeTime(firstOnOrAfter)
		if normalizedFirst.Equal(normalizedDTStart) {
			return nil
		}
	}

	if h.containsRDate(normalizedDTStart) {
		return nil
	}

	if ignoreExdate && h.containsExDate(normalizedDTStart) {
		return nil
	}

	if firstOnOrAfter.IsZero() {
		if lastBefore := h.set.Before(dtstart, true); !lastBefore.IsZero() {
			return &DTStartAlignmentError{
				DTStart:              normalizedDTStart,
				LastOccurrenceBefore: h.normalizeTime(lastBefore),
				AllDay:               h.allDay,
				reason:               dtstartAlignmentReasonNoOccurrenceOnOrAfter,
			}
		}

		return &DTStartAlignmentError{
			DTStart: normalizedDTStart,
			AllDay:  h.allDay,
			reason:  dtstartAlignmentReasonNoOccurrences,
		}
	}

	return &DTStartAlignmentError{
		DTStart:         normalizedDTStart,
		FirstOccurrence: h.normalizeTime(firstOnOrAfter),
		AllDay:          h.allDay,
		reason:          dtstartAlignmentReasonFirstMismatch,
	}
}

func (h *RecurrenceSetHelper) containsRDate(target time.Time) bool {
	for _, rdate := range h.rdates {
		if h.normalizeTime(rdate).Equal(target) {
			return true
		}
	}
	return false
}

func (h *RecurrenceSetHelper) containsExDate(target time.Time) bool {
	for _, exdate := range h.exdates {
		if h.normalizeTime(exdate).Equal(target) {
			return true
		}
	}
	return false
}

func formatTimeForError(t time.Time, allDay bool) string {
	if t.IsZero() {
		return "0001-01-01"
	}
	if allDay {
		return t.Format("2006-01-02")
	}
	return t.Format(time.RFC3339)
}

// AddRDate adds a recurrence date to the set
func (h *RecurrenceSetHelper) AddRDate(rdate time.Time) error {
	// Normalize time based on allDay flag
	normalizedTime := h.normalizeTime(rdate)

	// Check if already exists
	for _, existing := range h.rdates {
		if existing.Equal(normalizedTime) {
			return nil // Already exists, no need to add
		}
	}

	// Add to our tracking list
	h.rdates = append(h.rdates, normalizedTime)

	// Rebuild the set
	return h.rebuildSet()
}

// RemoveRDate removes a recurrence date from the set
func (h *RecurrenceSetHelper) RemoveRDate(rdate time.Time) error {
	normalizedTime := h.normalizeTime(rdate)

	// Find and remove from our tracking list
	for i, existing := range h.rdates {
		if existing.Equal(normalizedTime) {
			// Remove from slice
			h.rdates = append(h.rdates[:i], h.rdates[i+1:]...)
			// Rebuild the set
			return h.rebuildSet()
		}
	}

	return nil // Not found, nothing to remove
}

// AddExDate adds an exclusion date to the set
func (h *RecurrenceSetHelper) AddExDate(exdate time.Time) error {
	normalizedTime := h.normalizeTime(exdate)

	// Check if already exists
	for _, existing := range h.exdates {
		if existing.Equal(normalizedTime) {
			return nil // Already exists, no need to add
		}
	}

	// Add to our tracking list
	h.exdates = append(h.exdates, normalizedTime)

	// Rebuild the set
	return h.rebuildSet()
}

// Next returns the next occurrence after the given time.
func (h *RecurrenceSetHelper) Next(after time.Time) (time.Time, error) {
	if h.set == nil {
		return time.Time{}, fmt.Errorf("recurrence set not initialized")
	}
	occ := h.set.After(after, false)
	if occ.IsZero() {
		return time.Time{}, fmt.Errorf("no further occurrences")
	}
	return h.normalizeTime(occ), nil
}

// NextN returns up to n occurrences after the given time.
func (h *RecurrenceSetHelper) NextN(after time.Time, n int) ([]time.Time, error) {
	if n <= 0 {
		return nil, nil
	}
	occ := h.set.After(after, false)
	if occ.IsZero() {
		return nil, nil
	}
	result := make([]time.Time, 0, n)
	current := occ
	for len(result) < n && !current.IsZero() {
		result = append(result, h.normalizeTime(current))
		current = h.set.After(current, false)
	}
	return result, nil
}

// NextNIncludingDTStart returns up to n occurrences after the given time, guaranteeing that
// the configured DTSTART is treated as the first occurrence when the caller asks for times
// prior to it. This guards against provider-specific RRULE filters that exclude DTSTART.
func (h *RecurrenceSetHelper) NextNIncludingDTStart(after time.Time, n int) ([]time.Time, error) {
	occurrences, err := h.NextN(after, n)
	if err != nil {
		return nil, err
	}

	dtstart := h.DTStart()
	if dtstart.IsZero() || !after.Before(dtstart) {
		return occurrences, nil
	}

	if len(occurrences) == 0 {
		return []time.Time{dtstart}, nil
	}

	if dtstart.Before(occurrences[0]) {
		occurrences[0] = dtstart
	}

	return occurrences, nil
}

// RemoveExDate removes an exclusion date from the set
func (h *RecurrenceSetHelper) RemoveExDate(exdate time.Time) error {
	normalizedTime := h.normalizeTime(exdate)

	// Find and remove from our tracking list
	for i, existing := range h.exdates {
		if existing.Equal(normalizedTime) {
			// Remove from slice
			h.exdates = append(h.exdates[:i], h.exdates[i+1:]...)
			// Rebuild the set
			return h.rebuildSet()
		}
	}

	return nil // Not found, nothing to remove
}

// UpdateUntil updates the UNTIL parameter of the RRULE and clears COUNT (they are mutually exclusive)
func (h *RecurrenceSetHelper) UpdateUntil(until *time.Time) error {
	rrule := h.set.GetRRule()
	if rrule == nil {
		return fmt.Errorf("no rrule found in set")
	}

	if until != nil {
		normalizedUntil := h.normalizeTime(*until)
		rrule.Options.Until = normalizedUntil
		// Clear COUNT as UNTIL and COUNT are mutually exclusive per RFC 5545
		rrule.Options.Count = 0
	} else {
		// Clear UNTIL by setting it to zero time
		rrule.Options.Until = time.Time{}
	}

	return nil
}

// UpdateCount updates the COUNT parameter of the RRULE and clears UNTIL (they are mutually exclusive)
func (h *RecurrenceSetHelper) UpdateCount(count *int) error {
	rrule := h.set.GetRRule()
	if rrule == nil {
		return fmt.Errorf("no rrule found in set")
	}

	if count != nil && *count > 0 {
		rrule.Options.Count = *count
		// Clear UNTIL as UNTIL and COUNT are mutually exclusive per RFC 5545
		rrule.Options.Until = time.Time{}
	} else {
		// Clear COUNT by setting it to 0
		rrule.Options.Count = 0
	}

	return nil
}

// ToRRuleStrings converts the helper back to []string format for database storage
func (h *RecurrenceSetHelper) ToRRuleStrings() []string {
	// Use rrule-go's built-in method to generate RFC 5545 compliant strings
	// includeDTSTART=false because DTSTART is handled separately in Event/Task
	return h.set.Recurrence(false)
}

// Between returns occurrences between the given time range
func (h *RecurrenceSetHelper) Between(after, before time.Time) []time.Time {
	return h.set.Between(after, before, true)
}

// GetCount returns the COUNT value from the RRULE options (0 if not set)
func (h *RecurrenceSetHelper) GetCount() int {
	if h == nil || h.set == nil {
		return 0
	}
	r := h.set.GetRRule()
	if r == nil {
		return 0
	}
	return r.Options.Count
}

// GetRDates returns a copy of the current RDATE list
func (h *RecurrenceSetHelper) GetRDates() []time.Time {
	result := make([]time.Time, len(h.rdates))
	copy(result, h.rdates)
	return result
}

// GetExDates returns a copy of the current EXDATE list
func (h *RecurrenceSetHelper) GetExDates() []time.Time {
	result := make([]time.Time, len(h.exdates))
	copy(result, h.exdates)
	return result
}

// GetUntil returns the UNTIL time from the RRULE, or nil if not set
func (h *RecurrenceSetHelper) GetUntil() *time.Time {
	rrule := h.set.GetRRule()
	if rrule == nil {
		return nil
	}

	if rrule.Options.Until.IsZero() {
		return nil
	}

	return &rrule.Options.Until
}

// GetRRule returns the underlying *RRule for accessing Options
// This allows callers to access all RRule fields like Freq, Interval, Byweekday, etc.
func (h *RecurrenceSetHelper) GetRRule() *RRule {
	if h == nil || h.set == nil {
		return nil
	}
	return h.set.GetRRule()
}

// HasRRuleChanges checks if two RRule string arrays are different
func HasRRuleChanges(oldRules, newRules []string) bool {
	normalizedOld, errOld := NormalizeRecurrenceRuleset(oldRules)
	normalizedNew, errNew := NormalizeRecurrenceRuleset(newRules)

	// If either side fails to normalize, fall back to strict raw comparison.
	if errOld != nil || errNew != nil {
		return !stringSlicesEqual(oldRules, newRules)
	}

	if len(normalizedOld) != len(normalizedNew) {
		return true
	}

	if len(normalizedOld) == 0 {
		return false
	}

	sortedOld := make([]string, len(normalizedOld))
	copy(sortedOld, normalizedOld)
	sortedNew := make([]string, len(normalizedNew))
	copy(sortedNew, normalizedNew)

	sort.Strings(sortedOld)
	sort.Strings(sortedNew)

	for i := range sortedOld {
		if sortedOld[i] != sortedNew[i] {
			return true
		}
	}

	return false
}

func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// normalizeTime handles time normalization based on allDay flag
// Following RFC 5545: all-day events use floating time (represented as UTC)
func (h *RecurrenceSetHelper) normalizeTime(t time.Time) time.Time {
	if h.allDay {
		// All-day events: convert to floating time (UTC) as per RFC 5545
		year, month, day := t.Date()
		return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	}
	// Non all-day events: truncate to second precision
	return t.Truncate(time.Second)
}

// rebuildSet reconstructs the rrule.Set with current rdates and exdates
// This is necessary because rrule-go doesn't provide delete methods
func (h *RecurrenceSetHelper) rebuildSet() error {
	// Get the current rrule
	currentRRule := h.set.GetRRule()
	dtstart := h.set.GetDTStart()

	// Create a new set with the same configuration
	newSet := &Set{}

	// Set DTSTART if it exists
	if !dtstart.IsZero() {
		newSet.DTStart(dtstart)
	}

	// Set RRULE if it exists
	if currentRRule != nil {
		newSet.RRule(currentRRule)
	}

	// Set RDATEs using the batch method for efficiency
	if len(h.rdates) > 0 {
		newSet.SetRDates(h.rdates)
	}

	// Set EXDATEs using the batch method for efficiency
	if len(h.exdates) > 0 {
		newSet.SetExDates(h.exdates)
	}

	// Replace the old set
	h.set = newSet
	return nil
}

// String returns a human-readable representation of the recurrence set
func (h *RecurrenceSetHelper) String() string {
	parts := h.ToRRuleStrings()
	return strings.Join(parts, "\n")
}

// DTStart returns the normalized DTSTART associated with the recurrence set.
func (h *RecurrenceSetHelper) DTStart() time.Time {
	if h == nil || h.set == nil {
		return time.Time{}
	}

	dtstart := h.set.GetDTStart()
	if dtstart.IsZero() {
		return time.Time{}
	}

	return h.normalizeTime(dtstart)
}

// normalizeRecurrenceStrings processes input strings to ensure proper format.
// Uses NormalizeRecurrenceRuleset but adds logic to keep only the first RRULE.
// Handles:
// - RRULE strings with or without "RRULE:" prefix (only first RRULE is used)
// - RDATE strings with "RDATE:" prefix
// - EXDATE strings with "EXDATE:" prefix
func normalizeRecurrenceStrings(ruleset []string) ([]string, error) {
	if len(ruleset) == 0 {
		return nil, fmt.Errorf("empty input strings")
	}

	// First normalize using the standard function
	normalized, err := NormalizeRecurrenceRuleset(ruleset)
	if err != nil {
		return nil, err
	}

	// Then apply the "first RRULE only" rule for recurrence set parsing
	var result []string
	var foundRRule bool

	for _, str := range normalized {
		str = strings.TrimSpace(str)
		if str == "" {
			continue
		}

		// Handle RRULE strings - only keep the first one
		if strings.HasPrefix(strings.ToUpper(str), "RRULE:") {
			if !foundRRule {
				result = append(result, str)
				foundRRule = true
			}
			// Skip additional RRULE strings (only use the first one)
			continue
		}

		// Keep all RDATE and EXDATE strings
		result = append(result, str)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no valid recurrence strings found")
	}

	return result, nil
}
