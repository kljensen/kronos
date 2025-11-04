# Support Compound Relative Expressions with Commas and "and" Connectors

## Summary

Add support for natural language compound time expressions using commas and "and" connectors: "1 year, 2 months ago", "1 month and 5 days ago", "1 year, 1 month, 1 week, 1 day, 1 hour and 1 minute ago".

## Problem Statement

Users naturally write compound relative dates with commas and "and" connectors, but our parser only fully supports space-separated formats.

### What Works ✓
- ✅ Space-separated: `"1 year 2 months ago"`
- ✅ Simple "and": `"1 year and 2 months ago"`

### What Doesn't Work ✗
- ❌ Comma-separated: `"1 year, 2 months ago"` (only parses "1 year ago")
- ❌ Mixed: `"1 year, 1 month, 1 week, 1 day, 1 hour and 1 minute ago"`

## Research: Other Libraries

### Python's dateparser
```python
"1 year, 09 months, 01 weeks"  # Comma-separated
"1 year, 1 month, 1 week, 1 day, 1 hour and 1 minute ago"  # Full compound
```

### Go's go-dateparser (`tmp/test-cases/go-dateparser_parser-relative_test.go`)
```go
// Lines 71-74, 90-92, 103
{"1 year, 09 months,01 weeks"}
{"1 decade and 11 months"}
{"1 year, 1 month, 1 week, 1 day, 1 hour and 1 minute ago"}
```

### JavaScript's chrono (`tmp/test-cases/chrono_en_time_units_ago.test.ts`)
```javascript
// Lines 237-279
"15 hours 29 min ago"  // Space-separated (works)
"1d 21 h 25m ago"      // Abbreviated space-separated
```

## Test Cases

```go
// Basic comma
"1 year, 2 months ago" → 1 year + 2 months in past
"2 weeks, 3 days ago"  → 2 weeks + 3 days in past

// Mixed comma and "and"
"1 year, 1 month and 1 week ago"
"1 year, 1 month, 1 week, 1 day, 1 hour and 1 minute ago"
"2 years, 3 months and 5 days ago"

// Edge cases
"1 year,2 months ago"  // No space after comma
"1 year , 2 months ago" // Extra spaces
"1y, 2m, 3d ago"       // Abbreviated with commas

// Future tense
"in 1 year, 2 months"
"in 1 year and 2 months"
```

## Root Cause

File: `/workspace/en/time_unit_ago_parser.go:36`

**Current pattern requires spaces** (`\s+`) between units:
```go
`(?:\s+(?:(?:[0-9]+(?:[.,][0-9]+)?|half|a|an|the|few|couple|several)\s*(?:an?\s+)?)?` +
timeUnitPattern + `)*)\s{0,5}(?:ago|before|earlier)(?:\s|$|\b)`
```

## Proposed Fix

### Option A: Update Regex Pattern

```go
// New separator pattern: optional comma, optional "and", with spaces
separatorPattern := `(?:\s*,\s*|\s+)(?:and\s+)?`

pattern := approximationPattern +
    `((?:(?:[0-9]+(?:[.,][0-9]+)?|half|a|an|the|few|couple|several)\s*(?:an?\s+)?)?` +
    timeUnitPattern +
    `(?:` + separatorPattern + `(?:(?:[0-9]+(?:[.,][0-9]+)?|half|a|an|the|few|couple|several)\s*(?:an?\s+)?)?` +
    timeUnitPattern + `)*)\s{0,5}(?:ago|before|earlier)(?:\s|$|\b)`
```

### Option B: Normalize Before Parsing

```go
func ParseDuration(text string) kronos.Duration {
    // Strip commas and "and" before parsing
    normalized := regexp.MustCompile(`\s*,\s*|\s+and\s+`).ReplaceAllString(text, " ")
    // Continue with existing pattern...
}
```

**Recommendation**: Option B (simpler, requires fewer parser changes)

## Files to Modify

1. `/workspace/en/time_unit_ago_parser.go` - Update pattern or add normalization
2. `/workspace/en/time_unit_later_parser.go` - Same for "in X, Y"
3. `/workspace/en/time_unit_within_parser.go` - Same for "within X, Y"
4. `/workspace/en/constants.go` - Update `ParseDuration` if using normalization

## Priority & Complexity

**Priority**: MEDIUM  
**Justification**: Quality-of-life improvement, workaround exists (use spaces)

**Complexity**: LOW-MEDIUM  
**Estimated Effort**: 2-4 hours
- Regex update or normalization: 1-2 hours
- Testing: 1-2 hours

## Acceptance Criteria

- [ ] All 12 test cases pass
- [ ] Works with strict and non-strict modes
- [ ] Works with abbreviated units (1y, 2m)
- [ ] Works with fractional units
- [ ] No regressions in existing tests

---

**Labels**: `enhancement`, `parser`, `natural-language`
**Milestone**: v1.next
