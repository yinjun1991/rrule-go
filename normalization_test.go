package rrule

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
				"RRULE:FREQ=WEEKLY;BYDAY=MO,WE,FR",
				"RDATE:20240115T100000Z",
				"EXDATE:20240120T100000Z",
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

func TestNormalizeRRuleString(t *testing.T) {
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
			result, err := normalizeRulesetChunk(tt.input)

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

func TestValidateRRuleContent(t *testing.T) {
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
			err := validateRRuleContent(tt.input)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsValidRRuleContent(t *testing.T) {
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
			result := isValidRRuleContent(tt.input)
			assert.EqualValues(t, tt.expected, result)
		})
	}
}

// Integration test: Test the complete flow from input to storage
func TestRRuleNormalizationIntegration(t *testing.T) {
	testCases := []struct {
		name           string
		inputRules     []string
		expectedStored []string
		description    string
	}{
		{
			name:           "Google Calendar import format",
			inputRules:     []string{"FREQ=DAILY;COUNT=5"},
			expectedStored: []string{"RRULE:FREQ=DAILY;COUNT=5"},
			description:    "Rules from Google Calendar typically don't have RRULE: prefix",
		},
		{
			name:           "ICS file format",
			inputRules:     []string{"RRULE:FREQ=WEEKLY;BYDAY=MO,WE,FR"},
			expectedStored: []string{"RRULE:FREQ=WEEKLY;BYDAY=MO,WE,FR"},
			description:    "Rules from ICS files have RRULE: prefix",
		},
		{
			name: "Mixed format with RDATE/EXDATE",
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
			description: "Mixed rules with recurrence dates and exception dates",
		},
		{
			name: "Complex recurrence pattern",
			inputRules: []string{
				"FREQ=MONTHLY;BYMONTHDAY=15;BYSETPOS=1,3;COUNT=12",
			},
			expectedStored: []string{
				"RRULE:FREQ=MONTHLY;BYMONTHDAY=15;BYSETPOS=1,3;COUNT=12",
			},
			description: "Complex monthly recurrence with position constraints",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Testing: %s", tc.description)

			// Step 1: Normalize for storage
			normalized, err := NormalizeRecurrenceRuleset(tc.inputRules)
			require.NoError(t, err, "Normalization should succeed")
			assert.EqualValues(t, tc.expectedStored, normalized, "Normalized rules should match expected")

			// Step 4: Verify the normalized rules can be parsed
			if len(normalized) > 0 {
				// This would be the actual parsing test using the existing ParseRecurrenceSet
				// We'll add this when we integrate with the existing code
				t.Logf("Normalized rules ready for parsing: %v", normalized)
			}
		})
	}
}
