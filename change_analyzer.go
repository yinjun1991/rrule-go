package rrule

import (
	"fmt"
	"slices"
	"sort"
	"time"
)

// RecurrenceChangeType defines recurrence change types.
type RecurrenceChangeType int

const (
	// NoChange indicates no changes.
	NoChange RecurrenceChangeType = iota
	// FullRebuild indicates a full rebuild of all instances is required.
	FullRebuild
	// PartialUpdate indicates a partial update is sufficient.
	PartialUpdate
)

// RecurrenceChangeAnalysis holds recurrence change analysis results.
type RecurrenceChangeAnalysis struct {
	ChangeType RecurrenceChangeType
	// For PartialUpdate, the cutoff time for deleting instances.
	DeleteAfter *time.Time
	// For an extended UNTIL, the start time for generating new instances.
	GenerateFrom *time.Time
	// For an extended UNTIL, the end time for generating new instances.
	GenerateUntil *time.Time
	// New EXDATE entries.
	NewExDates []time.Time
	// Removed EXDATE entries.
	RemovedExDates []time.Time
	// Description of the changes.
	Description string
}

// RecurrenceDiffer provides recurrence change analysis.
type RecurrenceDiffer struct{}

// NewRecurrenceDiffer creates a new change analyzer.
func NewRecurrenceDiffer() *RecurrenceDiffer {
	return &RecurrenceDiffer{}
}

// AnalyzeChanges analyzes changes between two rule sets.
func (a *RecurrenceDiffer) AnalyzeChanges(oldRuleSet, newRuleSet []string) (*RecurrenceChangeAnalysis, error) {
	// If rules are identical, there are no changes.
	if !HasRRuleChanges(oldRuleSet, newRuleSet) {
		return &RecurrenceChangeAnalysis{
			ChangeType:  NoChange,
			Description: "No changes detected",
		}, nil
	}

	// Parse old and new rules; use allDay=false by default for analysis.
	// The actual allDay state should be provided by the caller per Event/Task.
	oldHelper, err := a.parseRulesForAnalysis(oldRuleSet)
	if err != nil {
		return nil, err
	}

	newHelper, err := a.parseRulesForAnalysis(newRuleSet)
	if err != nil {
		return nil, err
	}

	// If either ruleset is empty, a full rebuild is required.
	if oldHelper == nil || newHelper == nil {
		return &RecurrenceChangeAnalysis{
			ChangeType:  FullRebuild,
			Description: "Rule added or removed",
		}, nil
	}

	// Check whether a full rebuild is required.
	if a.requiresFullRebuild(oldHelper, newHelper) {
		return &RecurrenceChangeAnalysis{
			ChangeType:  FullRebuild,
			Description: "Core recurrence pattern changed, full rebuild required",
		}, nil
	}

	// Analyze partial updates.
	return a.analyzePartialUpdate(oldHelper, newHelper)
}

// parseRulesForAnalysis parses rules for analysis using a default dtstart.
func (a *RecurrenceDiffer) parseRulesForAnalysis(ruleset []string) (*Recurrence, error) {
	if len(ruleset) == 0 {
		return nil, nil
	}

	// Parse with the default dtstart.
	defaultDTStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	normalized, err := normalizeRecurrenceStrings(ruleset)
	if err != nil {
		return nil, err
	}

	set, err := StrSliceToRRuleSet(normalized)
	if err == nil {
		if set.GetDTStart().IsZero() {
			set.DTStart(defaultDTStart)
		}
		return set, nil
	}

	derivedDTStart, deriveErr := deriveMatchingDTStart(ruleset, defaultDTStart, false)
	if deriveErr != nil {
		return nil, err
	}

	set, err = StrSliceToRRuleSet(normalized)
	if err != nil {
		return nil, err
	}
	set.DTStart(derivedDTStart)
	return set, nil
}

// requiresFullRebuild checks whether a full rebuild is required.
func (a *RecurrenceDiffer) requiresFullRebuild(oldSet, newSet *Recurrence) bool {
	if oldSet == nil || newSet == nil || !oldSet.hasRule || !newSet.hasRule {
		return true
	}

	// Check DTSTART changes; this affects all occurrences.
	oldDTStart := oldSet.GetDTStart()
	newDTStart := newSet.GetDTStart()
	if !oldDTStart.Equal(newDTStart) {
		return true
	}

	// Check all-day vs timed transitions; time semantics differ, so rebuild.
	if oldSet.IsAllDay() != newSet.IsAllDay() {
		return true
	}

	// Check core property changes.
	if oldSet.freq != newSet.freq {
		return true
	}

	if oldSet.interval != newSet.interval {
		return true
	}

	// NOTE: COUNT and UNTIL changes can be handled with partial updates.
	// They only affect termination, not the recurrence pattern itself.

	// Check BYDAY changes.
	if !slices.Equal(oldSet.byweekday, newSet.byweekday) {
		return true
	}

	if !equalWeekdaySlices(oldSet.bynweekday, newSet.bynweekday) {
		return true
	}

	// Check other BY* changes.
	if !slices.Equal(oldSet.bymonth, newSet.bymonth) {
		return true
	}

	if !slices.Equal(oldSet.bymonthday, newSet.bymonthday) {
		return true
	}

	if !slices.Equal(oldSet.bynmonthday, newSet.bynmonthday) {
		return true
	}

	if !slices.Equal(oldSet.byyearday, newSet.byyearday) {
		return true
	}

	if !slices.Equal(oldSet.byweekno, newSet.byweekno) {
		return true
	}

	if !slices.Equal(oldSet.byhour, newSet.byhour) {
		return true
	}

	if !slices.Equal(oldSet.byminute, newSet.byminute) {
		return true
	}

	if !slices.Equal(oldSet.bysecond, newSet.bysecond) {
		return true
	}

	if !slices.Equal(oldSet.bysetpos, newSet.bysetpos) {
		return true
	}

	if oldSet.wkst != newSet.wkst {
		return true
	}

	if !slices.Equal(oldSet.byeaster, newSet.byeaster) {
		return true
	}

	// Check RDATE changes; they require a full rebuild.
	if a.hasRDateChange(oldSet, newSet) {
		return true
	}

	return false
}

// analyzePartialUpdate analyzes partial updates.
func (a *RecurrenceDiffer) analyzePartialUpdate(oldSet, newSet *Recurrence) (*RecurrenceChangeAnalysis, error) {
	analysis := &RecurrenceChangeAnalysis{
		ChangeType: PartialUpdate,
	}

	// Check UNTIL changes.
	if a.hasUntilChange(oldSet, newSet) {
		oldUntil := ruleUntilValue(oldSet)
		newUntil := ruleUntilValue(newSet)

		if oldUntil != nil && newUntil != nil {
			if newUntil.After(*oldUntil) {
				// UNTIL extended.
				analysis.GenerateFrom = oldUntil
				analysis.GenerateUntil = newUntil
				analysis.Description = "UNTIL date extended, generating new occurrences"
			} else {
				// UNTIL shortened.
				analysis.DeleteAfter = newUntil
				analysis.Description = "UNTIL date shortened, removing occurrences after new end date"
			}
		} else if oldUntil != nil && newUntil == nil {
			// UNTIL removed.
			analysis.GenerateFrom = oldUntil
			analysis.Description = "UNTIL removed"
		} else if oldUntil == nil && newUntil != nil {
			// UNTIL added.
			analysis.DeleteAfter = newUntil
			analysis.Description = "UNTIL added"
		}
	}

	// Check EXDATE changes.
	if a.hasExDateChange(oldSet, newSet) {
		addedExDates, removedExDates := a.compareExDates(oldSet.GetExDate(), newSet.GetExDate())
		analysis.NewExDates = addedExDates
		analysis.RemovedExDates = removedExDates
		if analysis.Description == "" {
			analysis.Description = "EXDATE changed"
		} else {
			analysis.Description += " and EXDATE changed"
		}
	}

	return analysis, nil
}

// hasUntilChange checks whether UNTIL changed.
func (a *RecurrenceDiffer) hasUntilChange(oldSet, newSet *Recurrence) bool {
	oldUntil := ruleUntilValue(oldSet)
	newUntil := ruleUntilValue(newSet)

	if oldUntil == nil && newUntil == nil {
		return false
	}
	if oldUntil == nil || newUntil == nil {
		return true
	}
	return !oldUntil.Equal(*newUntil)
}

// hasExDateChange checks whether EXDATE changed.
func (a *RecurrenceDiffer) hasExDateChange(oldSet, newSet *Recurrence) bool {
	oldExDates := oldSet.GetExDate()
	newExDates := newSet.GetExDate()

	if len(oldExDates) != len(newExDates) {
		return true
	}

	// Create a map for comparison.
	oldMap := make(map[int64]bool)
	for _, date := range oldExDates {
		oldMap[date.Unix()] = true
	}

	for _, date := range newExDates {
		if !oldMap[date.Unix()] {
			return true
		}
	}

	return false
}

// compareExDates compares two EXDATE lists.
func (a *RecurrenceDiffer) compareExDates(oldExDates, newExDates []time.Time) ([]time.Time, []time.Time) {
	// Create a map for comparison.
	oldMap := make(map[int64]time.Time)
	newMap := make(map[int64]time.Time)

	for _, date := range oldExDates {
		oldMap[date.Unix()] = date
	}

	for _, date := range newExDates {
		newMap[date.Unix()] = date
	}

	// Find added EXDATEs.
	var addedExDates []time.Time
	for timestamp, date := range newMap {
		if _, exists := oldMap[timestamp]; !exists {
			addedExDates = append(addedExDates, date)
		}
	}

	// Find removed EXDATEs.
	var removedExDates []time.Time
	for timestamp, date := range oldMap {
		if _, exists := newMap[timestamp]; !exists {
			removedExDates = append(removedExDates, date)
		}
	}

	return addedExDates, removedExDates
}

func deriveMatchingDTStart(ruleset []string, anchor time.Time, allDay bool) (time.Time, error) {
	normalized, err := normalizeRecurrenceStrings(ruleset)
	if err != nil {
		return time.Time{}, err
	}

	set, err := StrSliceToRRuleSet(normalized)
	if err != nil {
		return time.Time{}, err
	}

	set.SetAllDay(allDay)
	set.DTStart(anchor)

	first := set.After(anchor.Add(-time.Second), false)
	if first.IsZero() {
		return time.Time{}, fmt.Errorf("unable to derive DTSTART from recurrence rule")
	}

	return first, nil
}

func (a *RecurrenceDiffer) hasRDateChange(oldSet, newSet *Recurrence) bool {
	oldRDates := oldSet.GetRDate()
	newRDates := newSet.GetRDate()

	if len(oldRDates) != len(newRDates) {
		return true
	}

	// Use a map for fast comparison.
	oldMap := make(map[int64]bool)
	for _, date := range oldRDates {
		oldMap[date.Unix()] = true
	}

	for _, date := range newRDates {
		if !oldMap[date.Unix()] {
			return true
		}
	}

	return false
}

func ruleUntilValue(set *Recurrence) *time.Time {
	if set == nil || !set.hasRule {
		return nil
	}
	maxUntil := set.dtstart.Add(time.Duration(1<<63 - 1))
	if set.until.Equal(maxUntil) {
		return nil
	}
	value := set.until
	return &value
}

func equalWeekdaySlices(a, b []Weekday) bool {
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

func HasRRuleChanges(oldRules, newRules []string) bool {
	normalizedOld, errOld := NormalizeRecurrenceRuleset(oldRules)
	normalizedNew, errNew := NormalizeRecurrenceRuleset(newRules)

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
