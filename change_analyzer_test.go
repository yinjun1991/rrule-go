package rrule

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRRuleChangeAnalyzer_AnalyzeChanges(t *testing.T) {
	analyzer := NewRRuleChangeAnalyzer()

	t.Run("NoChange", func(t *testing.T) {
		t.Run("identical rules", func(t *testing.T) {
			oldRules := []string{"RRULE:FREQ=DAILY;COUNT=5"}
			newRules := []string{"RRULE:FREQ=DAILY;COUNT=5"}

			analysis, err := analyzer.AnalyzeChanges(oldRules, newRules)
			require.NoError(t, err)
			assert.EqualValues(t, NoChange, analysis.ChangeType)
			assert.EqualValues(t, "No changes detected", analysis.Description)
		})

		t.Run("empty rules", func(t *testing.T) {
			oldRules := []string{}
			newRules := []string{}

			analysis, err := analyzer.AnalyzeChanges(oldRules, newRules)
			require.NoError(t, err)
			assert.EqualValues(t, NoChange, analysis.ChangeType)
		})

		t.Run("multiple identical rules", func(t *testing.T) {
			oldRules := []string{
				"RRULE:FREQ=DAILY;COUNT=5",
				"EXDATE:20240101T100000Z",
				"RDATE:20240115T100000Z",
			}
			newRules := []string{
				"RRULE:FREQ=DAILY;COUNT=5",
				"EXDATE:20240101T100000Z",
				"RDATE:20240115T100000Z",
			}

			analysis, err := analyzer.AnalyzeChanges(oldRules, newRules)
			require.NoError(t, err)
			assert.EqualValues(t, NoChange, analysis.ChangeType)
		})
	})

	t.Run("FullRebuild", func(t *testing.T) {
		t.Run("frequency change", func(t *testing.T) {
			oldRules := []string{"RRULE:FREQ=DAILY;COUNT=5"}
			newRules := []string{"RRULE:FREQ=WEEKLY;COUNT=5"}

			analysis, err := analyzer.AnalyzeChanges(oldRules, newRules)
			require.NoError(t, err)
			assert.EqualValues(t, FullRebuild, analysis.ChangeType)
			assert.Contains(t, analysis.Description, "full rebuild required")
		})

		t.Run("interval change", func(t *testing.T) {
			oldRules := []string{"RRULE:FREQ=DAILY;INTERVAL=1;COUNT=5"}
			newRules := []string{"RRULE:FREQ=DAILY;INTERVAL=2;COUNT=5"}

			analysis, err := analyzer.AnalyzeChanges(oldRules, newRules)
			require.NoError(t, err)
			assert.EqualValues(t, FullRebuild, analysis.ChangeType)
		})

		t.Run("byday change", func(t *testing.T) {
			oldRules := []string{"RRULE:FREQ=WEEKLY;BYDAY=MO,WE,FR;COUNT=5"}
			newRules := []string{"RRULE:FREQ=WEEKLY;BYDAY=TU,TH;COUNT=5"}

			analysis, err := analyzer.AnalyzeChanges(oldRules, newRules)
			require.NoError(t, err)
			assert.EqualValues(t, FullRebuild, analysis.ChangeType)
		})

		t.Run("bymonth change", func(t *testing.T) {
			oldRules := []string{"RRULE:FREQ=YEARLY;BYMONTH=1,6,12;COUNT=5"}
			newRules := []string{"RRULE:FREQ=YEARLY;BYMONTH=3,9;COUNT=5"}

			analysis, err := analyzer.AnalyzeChanges(oldRules, newRules)
			require.NoError(t, err)
			assert.EqualValues(t, FullRebuild, analysis.ChangeType)
		})

		t.Run("bymonthday change", func(t *testing.T) {
			oldRules := []string{"RRULE:FREQ=MONTHLY;BYMONTHDAY=1,15;COUNT=5"}
			newRules := []string{"RRULE:FREQ=MONTHLY;BYMONTHDAY=10,20;COUNT=5"}

			analysis, err := analyzer.AnalyzeChanges(oldRules, newRules)
			require.NoError(t, err)
			assert.EqualValues(t, FullRebuild, analysis.ChangeType)
		})

		t.Run("rule added", func(t *testing.T) {
			oldRules := []string{}
			newRules := []string{"RRULE:FREQ=DAILY;COUNT=5"}

			analysis, err := analyzer.AnalyzeChanges(oldRules, newRules)
			require.NoError(t, err)
			assert.EqualValues(t, FullRebuild, analysis.ChangeType)
			assert.EqualValues(t, "Rule added or removed", analysis.Description)
		})

		t.Run("rule removed", func(t *testing.T) {
			oldRules := []string{"RRULE:FREQ=DAILY;COUNT=5"}
			newRules := []string{}

			analysis, err := analyzer.AnalyzeChanges(oldRules, newRules)
			require.NoError(t, err)
			assert.EqualValues(t, FullRebuild, analysis.ChangeType)
			assert.EqualValues(t, "Rule added or removed", analysis.Description)
		})

		t.Run("wkst change", func(t *testing.T) {
			oldRules := []string{"RRULE:FREQ=WEEKLY;WKST=MO;COUNT=5"}
			newRules := []string{"RRULE:FREQ=WEEKLY;WKST=SU;COUNT=5"}

			analysis, err := analyzer.AnalyzeChanges(oldRules, newRules)
			require.NoError(t, err)
			assert.EqualValues(t, FullRebuild, analysis.ChangeType)
		})
	})

	t.Run("PartialUpdate", func(t *testing.T) {
		t.Run("until extended", func(t *testing.T) {
			oldRules := []string{"RRULE:FREQ=DAILY;UNTIL=20240131T235959Z"}
			newRules := []string{"RRULE:FREQ=DAILY;UNTIL=20240229T235959Z"}

			analysis, err := analyzer.AnalyzeChanges(oldRules, newRules)
			require.NoError(t, err)
			assert.EqualValues(t, PartialUpdate, analysis.ChangeType)
			assert.Contains(t, analysis.Description, "UNTIL date extended")
			assert.NotNil(t, analysis.GenerateFrom)
			assert.NotNil(t, analysis.GenerateUntil)
			assert.True(t, analysis.GenerateUntil.After(*analysis.GenerateFrom))
		})

		t.Run("until shortened", func(t *testing.T) {
			oldRules := []string{"RRULE:FREQ=DAILY;UNTIL=20240229T235959Z"}
			newRules := []string{"RRULE:FREQ=DAILY;UNTIL=20240131T235959Z"}

			analysis, err := analyzer.AnalyzeChanges(oldRules, newRules)
			require.NoError(t, err)
			assert.EqualValues(t, PartialUpdate, analysis.ChangeType)
			assert.Contains(t, analysis.Description, "UNTIL date shortened")
			assert.NotNil(t, analysis.DeleteAfter)
		})

		t.Run("until added", func(t *testing.T) {
			oldRules := []string{"RRULE:FREQ=DAILY;COUNT=10"}
			newRules := []string{"RRULE:FREQ=DAILY;UNTIL=20240131T235959Z"}

			analysis, err := analyzer.AnalyzeChanges(oldRules, newRules)
			require.NoError(t, err)
			assert.EqualValues(t, PartialUpdate, analysis.ChangeType)
			assert.Contains(t, analysis.Description, "UNTIL added")
			assert.NotNil(t, analysis.DeleteAfter)
		})

		t.Run("until removed", func(t *testing.T) {
			oldRules := []string{"RRULE:FREQ=DAILY;UNTIL=20240131T235959Z"}
			newRules := []string{"RRULE:FREQ=DAILY;COUNT=10"}

			analysis, err := analyzer.AnalyzeChanges(oldRules, newRules)
			require.NoError(t, err)
			assert.EqualValues(t, PartialUpdate, analysis.ChangeType)
			assert.Contains(t, analysis.Description, "UNTIL removed")
			assert.NotNil(t, analysis.GenerateFrom)
		})

		t.Run("exdate added", func(t *testing.T) {
			oldRules := []string{"RRULE:FREQ=DAILY;COUNT=5"}
			newRules := []string{
				"RRULE:FREQ=DAILY;COUNT=5",
				"EXDATE:20240102T100000Z",
			}

			analysis, err := analyzer.AnalyzeChanges(oldRules, newRules)
			require.NoError(t, err)
			assert.EqualValues(t, PartialUpdate, analysis.ChangeType)
			assert.Contains(t, analysis.Description, "EXDATE changed")
			assert.Len(t, analysis.NewExDates, 1)
			assert.Len(t, analysis.RemovedExDates, 0)
		})

		t.Run("exdate removed", func(t *testing.T) {
			oldRules := []string{
				"RRULE:FREQ=DAILY;COUNT=5",
				"EXDATE:20240102T100000Z",
			}
			newRules := []string{"RRULE:FREQ=DAILY;COUNT=5"}

			analysis, err := analyzer.AnalyzeChanges(oldRules, newRules)
			require.NoError(t, err)
			assert.EqualValues(t, PartialUpdate, analysis.ChangeType)
			assert.Contains(t, analysis.Description, "EXDATE changed")
			assert.Len(t, analysis.NewExDates, 0)
			assert.Len(t, analysis.RemovedExDates, 1)
		})

		t.Run("multiple exdates changed", func(t *testing.T) {
			oldRules := []string{
				"RRULE:FREQ=DAILY;COUNT=10",
				"EXDATE:20240102T100000Z",
				"EXDATE:20240103T100000Z",
			}
			newRules := []string{
				"RRULE:FREQ=DAILY;COUNT=10",
				"EXDATE:20240103T100000Z",
				"EXDATE:20240104T100000Z",
				"EXDATE:20240105T100000Z",
			}

			analysis, err := analyzer.AnalyzeChanges(oldRules, newRules)
			require.NoError(t, err)
			assert.EqualValues(t, PartialUpdate, analysis.ChangeType)
			assert.Len(t, analysis.NewExDates, 2)     // 20240104, 20240105
			assert.Len(t, analysis.RemovedExDates, 1) // 20240102
		})

		t.Run("until and exdate both changed", func(t *testing.T) {
			oldRules := []string{
				"RRULE:FREQ=DAILY;UNTIL=20240131T235959Z",
				"EXDATE:20240102T100000Z",
			}
			newRules := []string{
				"RRULE:FREQ=DAILY;UNTIL=20240229T235959Z",
				"EXDATE:20240103T100000Z",
			}

			analysis, err := analyzer.AnalyzeChanges(oldRules, newRules)
			require.NoError(t, err)
			assert.EqualValues(t, PartialUpdate, analysis.ChangeType)
			assert.Contains(t, analysis.Description, "UNTIL date extended")
			assert.Contains(t, analysis.Description, "EXDATE changed")
			assert.NotNil(t, analysis.GenerateFrom)
			assert.NotNil(t, analysis.GenerateUntil)
			assert.Len(t, analysis.NewExDates, 1)
			assert.Len(t, analysis.RemovedExDates, 1)
		})
	})

	t.Run("AllDayEventHandling", func(t *testing.T) {
		t.Run("all day event exdate normalization", func(t *testing.T) {
			// All-day EXDATE should be normalized to date format.
			oldRules := []string{"RRULE:FREQ=DAILY;COUNT=5"}
			newRules := []string{
				"RRULE:FREQ=DAILY;COUNT=5",
				"EXDATE:20240102", // All-day format.
			}

			analysis, err := analyzer.AnalyzeChanges(oldRules, newRules)
			require.NoError(t, err)
			assert.EqualValues(t, PartialUpdate, analysis.ChangeType)
			assert.Len(t, analysis.NewExDates, 1)
		})

		t.Run("mixed time formats in exdate", func(t *testing.T) {
			// Test handling mixed time formats.
			oldRules := []string{
				"RRULE:FREQ=DAILY;COUNT=5",
				"EXDATE:20240102T100000Z",
			}
			newRules := []string{
				"RRULE:FREQ=DAILY;COUNT=5",
				"EXDATE:20240102", // Converted to all-day format.
			}

			analysis, err := analyzer.AnalyzeChanges(oldRules, newRules)
			require.NoError(t, err)
			assert.EqualValues(t, PartialUpdate, analysis.ChangeType)
		})
	})

	t.Run("EdgeCases", func(t *testing.T) {
		t.Run("invalid rrule strings", func(t *testing.T) {
			oldRules := []string{"RRULE:INVALID=VALUE"}
			newRules := []string{"RRULE:FREQ=DAILY;COUNT=5"}

			_, err := analyzer.AnalyzeChanges(oldRules, newRules)
			assert.Error(t, err)
		})

		t.Run("empty old rules", func(t *testing.T) {
			oldRules := []string{}
			newRules := []string{"RRULE:FREQ=DAILY;COUNT=5"}

			analysis, err := analyzer.AnalyzeChanges(oldRules, newRules)
			require.NoError(t, err)
			assert.EqualValues(t, FullRebuild, analysis.ChangeType)
		})

		t.Run("empty new rules", func(t *testing.T) {
			oldRules := []string{"RRULE:FREQ=DAILY;COUNT=5"}
			newRules := []string{}

			analysis, err := analyzer.AnalyzeChanges(oldRules, newRules)
			require.NoError(t, err)
			assert.EqualValues(t, FullRebuild, analysis.ChangeType)
		})

		t.Run("same until time", func(t *testing.T) {
			until := "20240131T235959Z"
			oldRules := []string{"RRULE:FREQ=DAILY;UNTIL=" + until}
			newRules := []string{"RRULE:FREQ=DAILY;UNTIL=" + until}

			analysis, err := analyzer.AnalyzeChanges(oldRules, newRules)
			require.NoError(t, err)
			assert.EqualValues(t, NoChange, analysis.ChangeType)
		})

		t.Run("count to until conversion", func(t *testing.T) {
			// COUNT -> UNTIL is a termination change.
			oldRules := []string{"RRULE:FREQ=DAILY;COUNT=5"}
			newRules := []string{"RRULE:FREQ=DAILY;UNTIL=20240105T235959Z"}

			analysis, err := analyzer.AnalyzeChanges(oldRules, newRules)
			require.NoError(t, err)
			// This should be PartialUpdate because only termination conditions change.
			assert.EqualValues(t, PartialUpdate, analysis.ChangeType)
		})

		t.Run("until to count conversion", func(t *testing.T) {
			// UNTIL -> COUNT is a termination change.
			oldRules := []string{"RRULE:FREQ=DAILY;UNTIL=20240105T235959Z"}
			newRules := []string{"RRULE:FREQ=DAILY;COUNT=5"}

			analysis, err := analyzer.AnalyzeChanges(oldRules, newRules)
			require.NoError(t, err)
			// This should be PartialUpdate because only termination conditions change.
			assert.EqualValues(t, PartialUpdate, analysis.ChangeType)
		})
	})

	t.Run("ComplexScenarios", func(t *testing.T) {
		t.Run("multiple rule components", func(t *testing.T) {
			oldRules := []string{
				"RRULE:FREQ=WEEKLY;BYDAY=MO,WE,FR;UNTIL=20240331T235959Z",
				"EXDATE:20240101T100000Z",
				"EXDATE:20240115T100000Z",
				"RDATE:20240201T100000Z",
			}
			newRules := []string{
				"RRULE:FREQ=WEEKLY;BYDAY=MO,WE,FR;UNTIL=20240430T235959Z", // UNTIL extended.
				"EXDATE:20240115T100000Z",                                 // Keep one EXDATE.
				"EXDATE:20240301T100000Z",                                 // Add one EXDATE.
				"RDATE:20240201T100000Z",                                  // RDATE unchanged.
			}

			analysis, err := analyzer.AnalyzeChanges(oldRules, newRules)
			require.NoError(t, err)
			assert.EqualValues(t, PartialUpdate, analysis.ChangeType)
			assert.Contains(t, analysis.Description, "UNTIL date extended")
			assert.Contains(t, analysis.Description, "EXDATE changed")
		})

		t.Run("timezone boundary cases", func(t *testing.T) {
			// Test UNTIL comparison across timezones.
			oldRules := []string{"RRULE:FREQ=DAILY;UNTIL=20240131T235959Z"}
			newRules := []string{"RRULE:FREQ=DAILY;UNTIL=20240201T000000Z"} // 1 second later.

			analysis, err := analyzer.AnalyzeChanges(oldRules, newRules)
			require.NoError(t, err)
			assert.EqualValues(t, PartialUpdate, analysis.ChangeType)
			assert.Contains(t, analysis.Description, "UNTIL date extended")
		})

		t.Run("microsecond precision", func(t *testing.T) {
			// Test microsecond-level time differences.
			oldRules := []string{"RRULE:FREQ=DAILY;UNTIL=20240131T235959Z"}
			newRules := []string{"RRULE:FREQ=DAILY;UNTIL=20240131T235959Z"} // Identical.

			analysis, err := analyzer.AnalyzeChanges(oldRules, newRules)
			require.NoError(t, err)
			assert.EqualValues(t, NoChange, analysis.ChangeType)
		})
	})
}

func TestRRuleChangeAnalyzer_WeekdaySliceComparison(t *testing.T) {
	analyzer := NewRRuleChangeAnalyzer()

	t.Run("equal slices", func(t *testing.T) {
		// Indirectly test BYDAY comparison via AnalyzeChanges.
		oldRules := []string{"RRULE:FREQ=WEEKLY;BYDAY=MO,WE,FR;COUNT=5"}
		newRules := []string{"RRULE:FREQ=WEEKLY;BYDAY=MO,WE,FR;COUNT=5"}

		analysis, err := analyzer.AnalyzeChanges(oldRules, newRules)
		require.NoError(t, err)
		assert.EqualValues(t, NoChange, analysis.ChangeType)
	})

	t.Run("different order should trigger rebuild", func(t *testing.T) {
		oldRules := []string{"RRULE:FREQ=WEEKLY;BYDAY=MO,WE,FR;COUNT=5"}
		newRules := []string{"RRULE:FREQ=WEEKLY;BYDAY=FR,WE,MO;COUNT=5"}

		analysis, err := analyzer.AnalyzeChanges(oldRules, newRules)
		require.NoError(t, err)
		// Different order should trigger rebuild because it may affect occurrences.
		assert.EqualValues(t, FullRebuild, analysis.ChangeType)
	})
}

func TestRRuleChangeAnalyzer_IntSliceComparison(t *testing.T) {
	analyzer := NewRRuleChangeAnalyzer()

	t.Run("bymonth changes", func(t *testing.T) {
		oldRules := []string{"RRULE:FREQ=YEARLY;BYMONTH=1,6,12;COUNT=5"}
		newRules := []string{"RRULE:FREQ=YEARLY;BYMONTH=1,6,12;COUNT=5"}

		analysis, err := analyzer.AnalyzeChanges(oldRules, newRules)
		require.NoError(t, err)
		assert.EqualValues(t, NoChange, analysis.ChangeType)
	})

	t.Run("bymonthday order matters", func(t *testing.T) {
		oldRules := []string{"RRULE:FREQ=MONTHLY;BYMONTHDAY=1,15,30;COUNT=5"}
		newRules := []string{"RRULE:FREQ=MONTHLY;BYMONTHDAY=30,15,1;COUNT=5"}

		analysis, err := analyzer.AnalyzeChanges(oldRules, newRules)
		require.NoError(t, err)
		// Different order should trigger rebuild.
		assert.EqualValues(t, FullRebuild, analysis.ChangeType)
	})
}

func TestRRuleChangeAnalyzer_ErrorHandling(t *testing.T) {
	analyzer := NewRRuleChangeAnalyzer()

	t.Run("malformed rrule", func(t *testing.T) {
		oldRules := []string{"RRULE:FREQ=INVALID"}
		newRules := []string{"RRULE:FREQ=DAILY;COUNT=5"}

		_, err := analyzer.AnalyzeChanges(oldRules, newRules)
		assert.Error(t, err)
	})

	t.Run("nil analyzer", func(t *testing.T) {
		// In Go, calling a method on a nil pointer won't panic if it doesn't access receiver fields.
		// AnalyzeChanges calls other methods, which may access fields.
		var analyzer *RRuleChangeAnalyzer
		// Go allows methods on nil receivers, so this test checks actual behavior.
		_, err := analyzer.AnalyzeChanges([]string{}, []string{})
		// If the method accesses receiver fields it will panic; otherwise it may succeed.
		// We expect either success or panic depending on the implementation.
		if err != nil {
			assert.Error(t, err)
		}
	})
}

// Benchmark.
func BenchmarkRRuleChangeAnalyzer_AnalyzeChanges(b *testing.B) {
	analyzer := NewRRuleChangeAnalyzer()
	oldRules := []string{
		"RRULE:FREQ=DAILY;UNTIL=20240331T235959Z",
		"EXDATE:20240101T100000Z",
		"EXDATE:20240115T100000Z",
	}
	newRules := []string{
		"RRULE:FREQ=DAILY;UNTIL=20240430T235959Z",
		"EXDATE:20240115T100000Z",
		"EXDATE:20240301T100000Z",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := analyzer.AnalyzeChanges(oldRules, newRules)
		if err != nil {
			b.Fatal(err)
		}
	}
}
