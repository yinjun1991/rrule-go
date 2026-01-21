package rrule

import (
	"fmt"
	"slices"
	"time"
)

// RRuleChangeType defines RRule change types.
type RRuleChangeType int

const (
	// NoChange indicates no changes.
	NoChange RRuleChangeType = iota
	// FullRebuild indicates a full rebuild of all instances is required.
	FullRebuild
	// PartialUpdate indicates a partial update is sufficient.
	PartialUpdate
)

// RRuleChangeAnalysis holds RRule change analysis results.
type RRuleChangeAnalysis struct {
	ChangeType RRuleChangeType
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

// RRuleChangeAnalyzer provides RRule change analysis.
type RRuleChangeAnalyzer struct{}

// NewRRuleChangeAnalyzer creates a new change analyzer.
func NewRRuleChangeAnalyzer() *RRuleChangeAnalyzer {
	return &RRuleChangeAnalyzer{}
}

// AnalyzeChanges analyzes changes between two rule sets.
func (a *RRuleChangeAnalyzer) AnalyzeChanges(oldRuleSet, newRuleSet []string) (*RRuleChangeAnalysis, error) {
	// If rules are identical, there are no changes.
	if !HasRRuleChanges(oldRuleSet, newRuleSet) {
		return &RRuleChangeAnalysis{
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
		return &RRuleChangeAnalysis{
			ChangeType:  FullRebuild,
			Description: "Rule added or removed",
		}, nil
	}

	// Check whether a full rebuild is required.
	if a.requiresFullRebuild(oldHelper, newHelper) {
		return &RRuleChangeAnalysis{
			ChangeType:  FullRebuild,
			Description: "Core recurrence pattern changed, full rebuild required",
		}, nil
	}

	// Analyze partial updates.
	return a.analyzePartialUpdate(oldHelper, newHelper)
}

// parseRulesForAnalysis parses rules for analysis using a default dtstart.
func (a *RRuleChangeAnalyzer) parseRulesForAnalysis(ruleset []string) (*RecurrenceSetHelper, error) {
	if len(ruleset) == 0 {
		return nil, nil
	}

	// Parse with the default dtstart.
	defaultDTStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	helper, err := ParseRecurrenceSet(ruleset, defaultDTStart, false)
	if err == nil {
		return helper, nil
	}

	derivedDTStart, deriveErr := deriveMatchingDTStart(ruleset, defaultDTStart, false)
	if deriveErr != nil {
		return nil, err
	}

	return ParseRecurrenceSet(ruleset, derivedDTStart, false)
}

// requiresFullRebuild checks whether a full rebuild is required.
func (a *RRuleChangeAnalyzer) requiresFullRebuild(oldHelper, newHelper *RecurrenceSetHelper) bool {
	// Get RRule objects.
	oldRRule := oldHelper.set.GetRRule()
	newRRule := newHelper.set.GetRRule()

	if oldRRule == nil || newRRule == nil {
		return true
	}

	// Check DTSTART changes; this affects all occurrences.
	oldDTStart := oldHelper.set.GetDTStart()
	newDTStart := newHelper.set.GetDTStart()
	if !oldDTStart.Equal(newDTStart) {
		return true
	}

	// Check all-day vs timed transitions; time semantics differ, so rebuild.
	if oldHelper.allDay != newHelper.allDay {
		return true
	}

	// Check core property changes.
	if oldRRule.Options.Freq != newRRule.Options.Freq {
		return true
	}

	if oldRRule.Options.Interval != newRRule.Options.Interval {
		return true
	}

	// NOTE: COUNT and UNTIL changes can be handled with partial updates.
	// They only affect termination, not the recurrence pattern itself.

	// Check BYDAY changes.
	if !slices.Equal(oldRRule.Options.Byweekday, newRRule.Options.Byweekday) {
		return true
	}

	// Check other BY* changes.
	if !slices.Equal(oldRRule.Options.Bymonth, newRRule.Options.Bymonth) {
		return true
	}

	if !slices.Equal(oldRRule.Options.Bymonthday, newRRule.Options.Bymonthday) {
		return true
	}

	if !slices.Equal(oldRRule.Options.Byyearday, newRRule.Options.Byyearday) {
		return true
	}

	if !slices.Equal(oldRRule.Options.Byweekno, newRRule.Options.Byweekno) {
		return true
	}

	if !slices.Equal(oldRRule.Options.Byhour, newRRule.Options.Byhour) {
		return true
	}

	if !slices.Equal(oldRRule.Options.Byminute, newRRule.Options.Byminute) {
		return true
	}

	if !slices.Equal(oldRRule.Options.Bysecond, newRRule.Options.Bysecond) {
		return true
	}

	if !slices.Equal(oldRRule.Options.Bysetpos, newRRule.Options.Bysetpos) {
		return true
	}

	if oldRRule.Options.Wkst != newRRule.Options.Wkst {
		return true
	}

	// Check RDATE changes; they require a full rebuild.
	if a.hasRDateChange(oldHelper, newHelper) {
		return true
	}

	return false
}

// analyzePartialUpdate analyzes partial updates.
func (a *RRuleChangeAnalyzer) analyzePartialUpdate(oldHelper, newHelper *RecurrenceSetHelper) (*RRuleChangeAnalysis, error) {
	analysis := &RRuleChangeAnalysis{
		ChangeType: PartialUpdate,
	}

	// Check UNTIL changes.
	if a.hasUntilChange(oldHelper, newHelper) {
		oldUntil := oldHelper.GetUntil()
		newUntil := newHelper.GetUntil()

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
	if a.hasExDateChange(oldHelper, newHelper) {
		addedExDates, removedExDates := a.compareExDates(oldHelper.GetExDates(), newHelper.GetExDates())
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
func (a *RRuleChangeAnalyzer) hasUntilChange(oldHelper, newHelper *RecurrenceSetHelper) bool {
	oldUntil := oldHelper.GetUntil()
	newUntil := newHelper.GetUntil()

	if oldUntil == nil && newUntil == nil {
		return false
	}
	if oldUntil == nil || newUntil == nil {
		return true
	}
	return !oldUntil.Equal(*newUntil)
}

// hasExDateChange checks whether EXDATE changed.
func (a *RRuleChangeAnalyzer) hasExDateChange(oldHelper, newHelper *RecurrenceSetHelper) bool {
	oldExDates := oldHelper.GetExDates()
	newExDates := newHelper.GetExDates()

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
func (a *RRuleChangeAnalyzer) compareExDates(oldExDates, newExDates []time.Time) ([]time.Time, []time.Time) {
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

// hasRDateChange checks RDATE changes between RecurrenceSetHelper instances.
func (a *RRuleChangeAnalyzer) hasRDateChange(oldHelper, newHelper *RecurrenceSetHelper) bool {
	oldRDates := oldHelper.GetRDates()
	newRDates := newHelper.GetRDates()

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
