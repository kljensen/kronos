package refiners

import "github.com/kljensen/kronos"

// Filter is a special type of Refiner that filters results based on validity.
type Filter interface {
	kronos.Refiner
	IsValid(context *kronos.ParsingContext, result *kronos.ParsingResult) bool
}

// BaseFilter provides a default Refine implementation for filters.
type BaseFilter struct{}

func (f *BaseFilter) Refine(context *kronos.ParsingContext, results []*kronos.ParsingResult) []*kronos.ParsingResult {
	filtered := make([]*kronos.ParsingResult, 0, len(results))
	for _, result := range results {
		// Each concrete filter must implement IsValid
		if filter, ok := interface{}(f).(Filter); ok {
			if filter.IsValid(context, result) {
				filtered = append(filtered, result)
			}
		}
	}
	return filtered
}

// MergingRefiner is a special type of Refiner that merges consecutive results.
type MergingRefiner interface {
	kronos.Refiner
	ShouldMergeResults(textBetween string, current, next *kronos.ParsingResult, context *kronos.ParsingContext) bool
	MergeResults(textBetween string, current, next *kronos.ParsingResult, context *kronos.ParsingContext) *kronos.ParsingResult
}

// BaseMergingRefiner provides a default Refine implementation for merging refiners.
type BaseMergingRefiner struct{}

func (m *BaseMergingRefiner) Refine(context *kronos.ParsingContext, results []*kronos.ParsingResult) []*kronos.ParsingResult {
	if len(results) < 2 {
		return results
	}

	merged := make([]*kronos.ParsingResult, 0, len(results))
	current := results[0]

	for i := 1; i < len(results); i++ {
		next := results[i]
		start := current.Index() + len(current.Text())
		end := next.Index()
		textBetween, okRange := kronos.SafeSlice(context.Text(), start, end)
		if !okRange {
			merged = append(merged, current)
			current = next
			continue
		}

		// Get the concrete implementation
		merger, ok := interface{}(m).(MergingRefiner)
		if !ok || !merger.ShouldMergeResults(textBetween, current, next, context) {
			merged = append(merged, current)
			current = next
		} else {
			current = merger.MergeResults(textBetween, current, next, context)
		}
	}

	if current != nil {
		merged = append(merged, current)
	}

	return merged
}
