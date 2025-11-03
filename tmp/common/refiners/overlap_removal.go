package refiners

import (
	kronos "github.com/kljensen/kronos"
)

// OverlapRemovalRefiner removes overlapping parse results.
// When two results overlap, it keeps the longer/more specific one.
type OverlapRemovalRefiner struct{}

// NewOverlapRemovalRefiner creates a new OverlapRemovalRefiner
func NewOverlapRemovalRefiner() *OverlapRemovalRefiner {
	return &OverlapRemovalRefiner{}
}

// Refine removes overlapping results
func (r *OverlapRemovalRefiner) Refine(context *kronos.ParsingContext, results []*kronos.ParsingResult) []*kronos.ParsingResult {
	if len(results) < 2 {
		return results
	}

	filteredResults := make([]*kronos.ParsingResult, 0, len(results))
	prevResult := results[0]

	for i := 1; i < len(results); i++ {
		result := results[i]

		// Check if results overlap
		if result.Index >= prevResult.Index+len(prevResult.Text) {
			// No overlap, keep previous result and move to current
			filteredResults = append(filteredResults, prevResult)
			prevResult = result
			continue
		}

		// Results overlap - keep the longer one
		var kept, removed *kronos.ParsingResult
		if len(result.Text) > len(prevResult.Text) {
			kept = result
			removed = prevResult
		} else {
			kept = prevResult
			removed = result
		}

		if context.Option.Debug {
			context.DebugLog("OverlapRemovalRefiner removing %s by %s", removed, kept)
		}

		prevResult = kept
	}

	// Add the last result
	if prevResult != nil {
		filteredResults = append(filteredResults, prevResult)
	}

	return filteredResults
}
