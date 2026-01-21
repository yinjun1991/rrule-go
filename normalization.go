package rrule

import (
	"fmt"
	"strings"
)

// NormalizeRecurrenceRuleset normalizes recurrence ruleset strings for storage.
// Ensures all RRULE strings have the proper "RRULE:" prefix for consistency.
// This is critical for:
// 1. Database data integrity
// 2. Google Calendar integration
// 3. ICS generation
// 4. Domain event parsing
func NormalizeRecurrenceRuleset(ruleset []string) ([]string, error) {
	if len(ruleset) == 0 {
		return nil, nil
	}

	normalized := make([]string, 0, len(ruleset))

	for _, rule := range ruleset {
		rule = strings.TrimSpace(rule)
		if rule == "" {
			continue // Skip empty strings
		}

		// Validate and normalize the rule
		normalizedRule, err := normalizeRulesetChunk(rule)
		if err != nil {
			return nil, fmt.Errorf("invalid rrule string '%s': %w", rule, err)
		}

		normalized = append(normalized, normalizedRule)
	}

	return normalized, nil
}

// normalizeRulesetChunk normalizes a single RRule string
func normalizeRulesetChunk(rule string) (string, error) {
	rule = strings.TrimSpace(rule)

	// Convert to uppercase for consistency
	upperRule := strings.ToUpper(rule)

	// Handle different rule types
	if strings.HasPrefix(upperRule, "RRULE:") {
		// Already has RRULE prefix, validate the content
		content := rule[6:] // Remove "RRULE:" prefix
		if err := validateRRuleContent(content); err != nil {
			return "", err
		}
		return rule, nil
	} else if strings.HasPrefix(upperRule, "RDATE:") || strings.HasPrefix(upperRule, "RDATE;") ||
		strings.HasPrefix(upperRule, "EXDATE:") || strings.HasPrefix(upperRule, "EXDATE;") {
		// RDATE and EXDATE are valid as-is (including those with TZID parameters)
		return rule, nil
	} else if isValidRRuleContent(rule) {
		// Missing RRULE prefix, add it
		if err := validateRRuleContent(rule); err != nil {
			return "", err
		}
		return "RRULE:" + rule, nil
	} else {
		return "", fmt.Errorf("unrecognized rule format")
	}
}

// validateRRuleContent validates the content part of an RRULE
func validateRRuleContent(content string) error {
	content = strings.TrimSpace(content)
	if content == "" {
		return fmt.Errorf("empty rrule content")
	}

	// Basic validation: must contain FREQ parameter
	upperContent := strings.ToUpper(content)
	if !strings.Contains(upperContent, "FREQ=") {
		return fmt.Errorf("rrule must contain FREQ parameter")
	}

	// Additional validation can be added here
	// For now, we rely on the rrule-go library for detailed validation
	return nil
}

// isValidRRuleContent checks if a string looks like valid RRULE content
func isValidRRuleContent(content string) bool {
	content = strings.TrimSpace(content)
	upperContent := strings.ToUpper(content)

	// Must contain FREQ and not start with known prefixes
	return strings.Contains(upperContent, "FREQ=") &&
		!strings.HasPrefix(upperContent, "RRULE:") &&
		!strings.HasPrefix(upperContent, "RDATE:") &&
		!strings.HasPrefix(upperContent, "EXDATE:")
}
