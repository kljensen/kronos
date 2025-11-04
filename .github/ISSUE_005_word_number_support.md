# Add Word Number Support for Time Quantifiers

## Summary

Support written numbers (one, two, three, etc.) in time expressions: "three days ago", "two weeks from now", "fifteen minutes ago".

## Problem

**Current**: Only "a/an" supported  
**Expected**: Support common word numbers (minimally 1-20, ideally with tens)

```go
❌ "three days ago"      → NOT PARSED
❌ "two weeks from now"  → NOT PARSED
❌ "fifteen minutes ago" → NOT PARSED
```

## Research

**All major libraries support this**:
- Chrono (JS): "three seconds ago" 
- Chronic (Ruby): "thirty-three days from now"
- dateparser (Python): "nine hours ago", "five years ago"
- Natty (Java): "four weeks ago", "twenty-five years"

## Current State

Infrastructure EXISTS but incomplete:
- `IntegerWordDictionary` has 1-12 (only used for ordinals)
- `NumberWordDictionary` missing cardinal numbers
- `ParseNumberPattern` ready to use word numbers

**File**: `/workspace/en/constants.go:119-145`

## Solution

Add to `NumberWordDictionary`:

```go
// ADD 13-19
"thirteen": 13.0, "fourteen": 14.0, "fifteen": 15.0,
"sixteen": 16.0, "seventeen": 17.0, "eighteen": 18.0, "nineteen": 19.0,

// ADD tens
"twenty": 20.0, "thirty": 30.0, "forty": 40.0,
"fifty": 50.0, "sixty": 60.0, "seventy": 70.0,
"eighty": 80.0, "ninety": 90.0,
```

## Test Cases

```go
"one day ago", "two weeks from now", "three minutes ago"
"five hours later", "twelve months ago"
"thirteen days ago", "fifteen minutes ago"
"twenty hours ago", "thirty days from now"
```

## Scope

**Minimum**: 1-20 + tens (20, 30, 40, 50, 60, 70, 80, 90)  
**Rationale**: Covers 90%+ usage. Users write "23 days" not "twenty-three days"

**Out of scope**: Compound numbers (twenty-three), hundreds, millions

## Files to Modify

1. `/workspace/en/constants.go:119-145` - Add dictionary entries
2. `/workspace/en/time_unit_ago_parser_test.go:188-201` - Uncomment SKIPped tests

## Effort

**Priority**: MEDIUM  
**Complexity**: MEDIUM  
**Effort**: 2-4 hours (mostly testing)

**Breakdown**:
- Dictionary entries: 30 minutes
- Enable tests: 30 minutes  
- Add comprehensive tests: 2-3 hours

## Acceptance Criteria

- [ ] Parse "one" through "twelve"
- [ ] Parse "thirteen" through "nineteen"
- [ ] Parse tens (twenty, thirty, etc.)
- [ ] Works across all time unit parsers (ago, later, within)
- [ ] Uncomment SKIPped test at line 188

---

**Labels**: `enhancement`, `parser`, `word-numbers`
**Milestone**: v1.next
