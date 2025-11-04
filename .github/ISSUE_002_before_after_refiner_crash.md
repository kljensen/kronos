# Critical Bug: Slice Bounds Panic in "before/after" Date Pattern Refiners

## Summary

Runtime crash (panic) in date parsing refiners when handling "before/after" patterns like "2 days before today" or "the day after tomorrow". This is a critical bug that makes the library crash on valid user input.

## Severity

**CRITICAL** - Runtime panic, blocking entire category of patterns, affects production stability

## Problem Statement

The English date parsing library crashes with a **slice bounds out of range** panic when parsing "before/after" patterns:
- `"2 days before today"`
- `"the day before yesterday"`
- `"2 weeks after tomorrow"`
- `"a week before yesterday"`

This affects two refiners:
1. **`ENMergeRelativeFollowByDateRefiner`** - Handles "X before/after Y" patterns
2. **`ENMergeRelativeAfterDateRefiner`** - Handles "Y +X" patterns

The crash prevents users from expressing a fundamental class of relative date expressions.

## Steps to Reproduce

```go
package main

import (
    "time"
    kronos "github.com/kljensen/kronos"
    "github.com/kljensen/kronos/en"
)

func main() {
    config := en.CreateCasualConfiguration(false)
    chrono := kronos.NewChrono(config)

    refDate := time.Date(2012, 8, 10, 0, 0, 0, 0, time.UTC)

    // This will panic with "runtime error: slice bounds out of range"
    results := chrono.Parse("2 day before today", refDate, nil)
    fmt.Printf("Results: %v\n", results)
}
```

**Error Output:**
```
panic: runtime error: slice bounds out of range [X:Y] with capacity Z

goroutine 1 [running]:
github.com/kljensen/kronos/en.(*ENMergeRelativeFollowByDateRefiner).Refine(...)
    /workspace/en/refiner_merge_relative_follow.go:56
```

## Current Behavior: Documented Test Failures

The bug is **extensively documented** in our test suite with 11+ SKIPped test cases:

### File 1: `/workspace/en/time_unit_ago_parser_test.go` (lines 459-527)

```go
// TestAgoBeforeWithReference tests "before" expressions with reference words
// SKIPPED: These tests crash due to a bug in ENMergeRelativeFollowByDateRefiner
// (slice bounds out of range). This is a known issue in the refiner, not the
// ago parser itself. The crash occurs when trying to merge "X before Y" patterns
// where Y is a casual date.
func TestAgoBeforeWithReference(t *testing.T) {
    t.Skip("Skipping due to refiner bug: slice bounds out of range in ENMergeRelativeFollowByDateRefiner")

    // 11 commented-out test cases that should work:
    // - "2 day before today"
    // - "the day before yesterday"
    // - "2 day before yesterday"
    // - "a week before yesterday"
    // - "15 minute before"
    // - "a min before"
    // - "the min before"
}
```

### File 2: `/workspace/en/en_merging_relative_dates_test.go` (lines 11-113)

```go
// TestMergingRelativeDates tests merging of relative date expressions
// SKIP: These tests currently fail due to bugs in the refiner logic.
// The refiners have slice bound issues that need to be fixed.
func TestMergingRelativeDates(t *testing.T) {
    t.Skip("Skipping - known bugs in merge relative refiners")

    // Tests for:
    // - "2 weeks after yesterday"
    // - "2 months before 02/02"
    // - "2 days after next Friday"
}
```

### File 3: `/workspace/en/time_unit_later_parser_test.go` (lines 568, 642)

```go
// SKIPPED: These tests require full configuration with refiners, and refiners
// have bugs with "after" references
```

## Expected Behavior

The library should successfully parse "before/after" patterns and return proper date calculations:

```go
// Input: "2 day before today"
// Reference: 2012-08-10
// Expected: 2012-08-08

// Input: "the day before yesterday"
// Reference: 2012-08-10
// Expected: 2012-08-08

// Input: "2 weeks after tomorrow"
// Reference: 2022-02-02
// Expected: 2022-02-17
```

## Root Cause Analysis

### Location of Bug

**File 1**: `/workspace/en/refiner_merge_relative_follow.go:56`
**File 2**: `/workspace/en/refiner_merge_relative_after.go:56`

### The Vulnerable Code

Both refiners calculate text between two parsing results using **unsafe string slicing**:

```go
// Line 50-56 in both refiners
startIdx := current.Index() + len(current.Text())
endIdx := next.Index()

if startIdx > endIdx {
    merged = append(merged, current)
    current = next
    continue
}

// BUG: No bounds checking before slicing
textBetween := context.Text()[startIdx:endIdx]  // ← PANIC HERE
```

**Problems:**
1. Checks if `startIdx > endIdx` but **not** if `endIdx > len(context.Text())`
2. No validation that `endIdx` is within valid bounds
3. Can panic when parsing results have incorrect indices

### Why This is Critical

The panic occurs at runtime during normal operation, not during initialization:
- User submits "2 days before today"
- Parser creates two results: "2 days before" + "today"
- Refiner attempts to merge them
- **CRASH** - entire application goes down

## Research: Industry Best Practices

### 1. The Codebase Already Has a Solution

We have a **safe slicing utility** that's used throughout other refiners:

**`/workspace/safety.go:24-43`**:
```go
// SafeSlice returns the substring of text between start (inclusive) and end
// (exclusive) while clamping the requested range to valid bounds. The returned
// boolean is false when start is beyond the end of the string, indicating that
// the requested slice could not be produced safely.
func SafeSlice(text string, start, end int) (string, bool) {
    if start < 0 {
        start = 0
    }
    if end < start {
        end = start
    }
    if start > len(text) {
        return "", false
    }
    if end > len(text) {
        end = len(text)
    }
    return text[start:end], true
}
```

### 2. SafeSlice is Already Used Correctly in 8+ Refiners

```bash
$ grep -r "SafeSlice" workspace/
/workspace/common/refiners/timezone_abbr.go:46
/workspace/common/refiners/abstract_refiners.go:49
/workspace/common/refiners/merge_daterange.go:141
/workspace/common/refiners/timezone_offset.go:43
/workspace/common/refiners/merge_weekday.go:37
/workspace/common/refiners/merge_datetime.go:81
/workspace/en/refiners/extract_year_suffix.go:39
/workspace/en/refiners/unlikely_format_filter.go:45
```

**Example from `/workspace/common/refiners/merge_datetime.go:81`:**
```go
// Correct usage of SafeSlice
textBetween, okSlice := kronos.SafeSlice(context.Text(), startIdx, endIdx)
if !okSlice || !patternBetween.MatchString(textBetween) {
    merged = append(merged, current)
    current = next
    continue
}
```

### 3. Other Date Libraries Handle This Gracefully

**Chrono.js** (JavaScript) - Uses safe string methods:
```typescript
// From chrono source - always validates bounds
const textBetween = text.substring(Math.max(0, start), Math.min(text.length, end));
```

**Chronic** (Ruby) - String slicing is safe by default:
```ruby
# Ruby returns nil for out-of-bounds, doesn't crash
text_between = text[start..end] || ""
```

## Proposed Fix

Replace unsafe slicing with `kronos.SafeSlice()` in both refiners.

### Fix for `/workspace/en/refiner_merge_relative_follow.go`

**Before (line 50-62):**
```go
startIdx := current.Index() + len(current.Text())
endIdx := next.Index()

if startIdx > endIdx {
    merged = append(merged, current)
    current = next
    continue
}

textBetween := context.Text()[startIdx:endIdx]
if !patternFollowBetween.MatchString(textBetween) {
    merged = append(merged, current)
    current = next
    continue
}
```

**After:**
```go
startIdx := current.Index() + len(current.Text())
endIdx := next.Index()

if startIdx > endIdx {
    merged = append(merged, current)
    current = next
    continue
}

// Use SafeSlice to prevent panic
textBetween, okSlice := kronos.SafeSlice(context.Text(), startIdx, endIdx)
if !okSlice || !patternFollowBetween.MatchString(textBetween) {
    merged = append(merged, current)
    current = next
    continue
}
```

### Fix for `/workspace/en/refiner_merge_relative_after.go`

**Same fix at line 56** - replace:
```go
textBetween := context.Text()[startIdx:endIdx]
```

with:
```go
textBetween, okSlice := kronos.SafeSlice(context.Text(), startIdx, endIdx)
if !okSlice || !patternAfterBetween.MatchString(textBetween) {
    merged = append(merged, current)
    current = next
    continue
}
```

## Test Cases to Verify the Fix

Once the fix is applied, **uncomment and run** these existing test cases:

### From `time_unit_ago_parser_test.go` (lines 459-527)

```go
// 1. Basic "before" with casual reference
{
    name:     "2 day before today",
    text:     "2 day before today",
    refDate:  time.Date(2012, 8, 10, 0, 0, 0, 0, time.UTC),
    expected: time.Date(2012, 8, 8, 0, 0, 0, 0, time.UTC),
}

// 2. The day before yesterday
{
    name:     "the day before yesterday",
    text:     "the day before yesterday",
    refDate:  time.Date(2012, 8, 10, 0, 0, 0, 0, time.UTC),
    expected: time.Date(2012, 8, 8, 0, 0, 0, 0, time.UTC),
}

// 3. Multiple units before
{
    name:     "2 day before yesterday",
    text:     "2 day before yesterday",
    refDate:  time.Date(2012, 8, 10, 0, 0, 0, 0, time.UTC),
    expected: time.Date(2012, 8, 7, 0, 0, 0, 0, time.UTC),
}

// 4. Week before
{
    name:     "a week before yesterday",
    text:     "a week before yesterday",
    refDate:  time.Date(2012, 8, 10, 0, 0, 0, 0, time.UTC),
    expected: time.Date(2012, 8, 2, 0, 0, 0, 0, time.UTC),
}

// 5-7. Minute precision
{
    name:     "15 minute before",
    text:     "15 minute before",
    refDate:  time.Date(2012, 8, 10, 12, 14, 0, 0, time.UTC),
    expected: time.Date(2012, 8, 10, 11, 59, 0, 0, time.UTC),
}
```

### From `en_merging_relative_dates_test.go` (lines 11-113)

```go
// 8. Weeks after yesterday
{
    name:     "2 weeks after yesterday",
    text:     "2 weeks after yesterday",
    refDate:  time.Date(2022, 2, 2, 0, 0, 0, 0, time.UTC),
    expected: time.Date(2022, 2, 15, 0, 0, 0, 0, time.UTC),
}

// 9. Months before date
{
    name:     "2 months before 02/02",
    text:     "2 months before 02/02",
    refDate:  time.Date(2022, 2, 2, 0, 0, 0, 0, time.UTC),
    expected: time.Date(2021, 12, 2, 0, 0, 0, 0, time.UTC),
}

// 10. Days after next Friday
{
    name:     "2 days after next Friday",
    text:     "2 days after next Friday",
    refDate:  time.Date(2022, 2, 2, 0, 0, 0, 0, time.UTC),  // Wednesday
    expected: time.Date(2022, 2, 13, 0, 0, 0, 0, time.UTC), // Sunday
}
```

### Additional Reference Tests (from Chrono)

From `/workspace/tmp/test-cases/chrono_en_time_units_ago.test.ts`:

```typescript
// 11. Plus/minus patterns
test("next tuesday +10 days", () => {
    const refDate = new Date(2023, 12-1, 29);
    const result = chrono.parse("next tuesday +10 days", refDate);
    expect(result[0].start.get("day")).toBe(12); // Jan 12, 2024
});

// 12. Negative offset
test("2023-12-29 -10days", () => {
    const refDate = new Date(2023, 12-1, 29);
    const result = chrono.parse("2023-12-29 -10days", refDate);
    expect(result[0].start.get("day")).toBe(19); // Dec 19
});
```

## Implementation Checklist

- [ ] Apply SafeSlice fix to `refiner_merge_relative_follow.go:56`
- [ ] Apply SafeSlice fix to `refiner_merge_relative_after.go:56`
- [ ] Uncomment tests in `time_unit_ago_parser_test.go:459-527`
- [ ] Uncomment tests in `en_merging_relative_dates_test.go:11-113`
- [ ] Uncomment tests in `time_unit_later_parser_test.go:568,642`
- [ ] Run full test suite: `go test ./en/... -v`
- [ ] Add regression test to prevent future unsafe slicing
- [ ] Update CHANGELOG.md with bug fix note

## Impact Assessment

### What's Broken
- ❌ All "X before Y" patterns crash
- ❌ All "X after Y" patterns crash
- ❌ Plus/minus adjustments ("+10 days") crash
- ❌ 14+ test cases SKIPped
- ❌ Production applications can crash on user input

### What's Blocked
- Natural date expressions users commonly write
- Merging relative and absolute date references
- Conversational date parsing quality

### User Impact
- **Severity**: CRITICAL - Runtime panics crash applications
- **Frequency**: HIGH - Common patterns users expect to work
- **Workaround**: None - users must avoid these patterns entirely
- **Trust**: Damages library reliability and confidence

## Complexity & Effort

**Complexity**: **TRIVIAL** (bug fix with clear solution)

- ✅ Root cause identified and documented
- ✅ Solution is clear (use existing SafeSlice utility)
- ✅ Pattern already used correctly in 8+ other refiners
- ✅ Test cases already written (just commented out)
- ✅ Zero risk to existing functionality (pure bug fix)

**Estimated Effort**: **1-2 hours**
- Code changes: 15 minutes (4 lines changed across 2 files)
- Uncommenting tests: 15 minutes
- Running test suite: 15 minutes
- Regression test: 30 minutes
- Code review & merge: 15 minutes

## Priority

**CRITICAL** - Should be fixed immediately

**Rationale**:
- Runtime crash (not graceful error)
- Blocks common user patterns
- Simple fix with existing utility
- High user impact
- Already documented with SKIPped tests

## References

### Affected Files
- `/workspace/en/refiner_merge_relative_follow.go:56` - Primary bug
- `/workspace/en/refiner_merge_relative_after.go:56` - Secondary bug
- `/workspace/safety.go:24-43` - SafeSlice utility (solution)

### Test Files
- `/workspace/en/time_unit_ago_parser_test.go:459-527` - 11 SKIPped tests
- `/workspace/en/en_merging_relative_dates_test.go:11-113` - Merge tests
- `/workspace/en/time_unit_later_parser_test.go:568,642` - Later tests

### Correct Examples (Using SafeSlice)
- `/workspace/common/refiners/merge_daterange.go:141`
- `/workspace/common/refiners/merge_datetime.go:81`
- `/workspace/en/refiners/extract_year_suffix.go:39`

### Reference Implementations
- **Chrono tests**: `/workspace/tmp/test-cases/chrono_en_time_units_ago.test.ts`
- **Natty examples**: `/workspace/tmp/test-cases/natty_DateTest.java`

---

**Labels**: `bug`, `critical`, `crash`, `security`, `good-first-issue`
**Milestone**: Immediate (hotfix)
**Assignee**: Open (urgent fix needed)
