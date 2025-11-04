# Business Quarter Support (Q1-Q4 Notation)

## Summary

Add support for business quarter references: "Q1 2024", "Q3", "1st quarter 2023", "this quarter".

## Problem

**Current**: Only relative quarters work ("next quarter")  
**Missing**: Absolute quarter notation (Q1-Q4) and "this quarter"

```go
✅ Parse("next quarter")  // Works - adds 3 months
❌ Parse("Q1 2024")       // NO MATCH
❌ Parse("Q3")            // NO MATCH  
❌ Parse("this quarter")  // Parser missing switch case
```

## Research

### Chrono (JS) - `tmp/test-cases/chrono_en_relative.test.ts:139-176`

```javascript
// Lines 139-163
"next quarter" (Jan 22) → Apr 22  // +3 months
"next qtr" (Oct 22) → Jan 22      // +3 months, wraps year
"next two quarter" (Jan 22) → Jul 22  // +6 months
```

**Behavior**: Adds 3 months per quarter, keeps same day-of-month

## Current Infrastructure

**Already exists**:
- ✅ `TimeunitQuarter` constant (`/workspace/types.go:38`)
- ✅ Quarter-to-months conversion (`/workspace/duration.go:46-50`)
- ✅ "quarter" in dictionaries (`/workspace/en/constants.go:234-236`)
- ✅ "quarter" in relative parser (`/workspace/en/relative_date_parser.go:24`)

**Missing**:
- ❌ Q1/Q2/Q3/Q4 notation parser
- ❌ "this quarter" switch case

## Design Decisions

### Calendar Quarters (Standard)
- **Q1**: January-March → resolves to Jan 1
- **Q2**: April-June → resolves to Apr 1  
- **Q3**: July-September → resolves to Jul 1
- **Q4**: October-December → resolves to Oct 1

**Rationale**: Standard business convention. Fiscal quarters vary by organization.

### Resolution
- "Q1" without year → use current year (with preference settings)
- "Q1 2024" → explicit year
- Returns start date (first day of quarter)

## Test Cases

```go
Parse("Q1", 2024-06-15) → 2024-01-01
Parse("Q2", 2024-06-15) → 2024-04-01
Parse("Q1 2024")        → 2024-01-01
Parse("2024 Q3")        → 2024-07-01
Parse("1st quarter 2024") → 2024-01-01
Parse("fourth quarter") → 2024-10-01 (current year)

// "this quarter"
Parse("this quarter", 2024-02-15) → 2024-01-01 (in Q1)
Parse("this quarter", 2024-08-20) → 2024-07-01 (in Q3)
```

## Implementation

### 1. Fix "this quarter" (Quick Win)

**File**: `/workspace/en/relative_date_parser.go:81`

Add to switch statement:

```go
case kronos.TimeunitQuarter:
    // Calculate current quarter (1-4)
    currentQuarter := (int(refDate.Month()) - 1) / 3 + 1
    startMonth := (currentQuarter - 1) * 3 + 1
    
    date := time.Date(refDate.Year(), time.Month(startMonth), 1, 
                     0, 0, 0, 0, refDate.Location())
    components.Assign(kronos.ComponentYear, date.Year())
    components.Assign(kronos.ComponentMonth, int(date.Month()))
    components.Imply(kronos.ComponentDay, 1)
    components.SetPeriod(kronos.PeriodMonth)
```

### 2. Create Quarter Parser (Full Feature)

**New file**: `/workspace/en/quarter_parser.go`

```go
// Pattern: Q[1-4] or ordinal quarter with optional year
pattern := `(?i)(Q[1-4]|(?:first|second|third|fourth|1st|2nd|3rd|4th)\s+quarter)(?:\s+(\d{4}))?`

// QuarterDictionary
var QuarterToMonth = map[int]int{
    1: 1,   // Q1 → January
    2: 4,   // Q2 → April
    3: 7,   // Q3 → July
    4: 10,  // Q4 → October
}
```

## Effort

**Priority**: MEDIUM (business use case)  
**Complexity**: MEDIUM  
**Effort**: 2-4 hours

**Phased approach**:
1. Fix "this quarter" (15 min) - Quick win
2. Add Q1-Q4 notation (1-2 hours)
3. Add written quarters (30 min)

## Acceptance Criteria

- [ ] Parse Q1, Q2, Q3, Q4
- [ ] Parse "Q1 2024", "2024 Q1"
- [ ] Parse "1st quarter 2024", "fourth quarter"
- [ ] Fix "this quarter" (add switch case)
- [ ] 10 test cases pass

---

**Labels**: `enhancement`, `quarter`, `business`
**Milestone**: v1.next
