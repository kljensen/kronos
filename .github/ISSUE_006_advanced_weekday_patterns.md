# Advanced Weekday Patterns: "after next", "before last", "coming"

## Summary

Support advanced weekday modifiers: "the Monday after next", "Tuesday before last", "this coming Monday", "this upcoming Friday".

## Problem

**Current**: Basic next/last/this works  
**Missing**: Double-hop modifiers and coming/upcoming synonyms

```go
❌ "the Monday after next" → 2 weeks forward
❌ "Tuesday before last"   → 2 weeks backward
❌ "this coming Monday"    → synonym for "next Monday"
```

## Research

### Natty (Java) - `tmp/test-cases/natty_DateTest.java`

```java
// Lines 101-102, 127-129, 133
"the saturday after next"  // March 19, 2011
"the monday after next"    // March 14, 2011
"tuesday before last"      // February 15, 2011
"this coming monday"       // March 7, 2011
```

Reference: Monday, February 28, 2011

## Test Cases

| Input | Reference | Expected |
|-------|-----------|----------|
| "the saturday after next" | Mon Feb 28 | Sat Mar 19 (+19 days) |
| "the monday after next" | Mon Feb 28 | Mon Mar 14 (+14 days) |
| "tuesday before last" | Mon Feb 28 | Tue Feb 15 (-13 days) |
| "this coming monday" | Mon Feb 28 | Mon Mar 7 (next Mon) |
| "this upcoming friday" | Mon Feb 28 | Fri Mar 4 (next Fri) |

## Implementation

### Regex Update (`/workspace/en/weekday_parser.go:25-31`)

```go
// Add:
`(?:the\s+)?` +                                    // optional "the"
`(?:(this|last|past|next|coming|upcoming)\s*)?` + // add coming/upcoming
`(Monday|Tuesday|...|weekend|weekday)` +
`(?:\s+(after\s+next|before\s+last))?` +          // NEW extended modifier
```

### Logic

```go
// "coming/upcoming" → treat as "next"
if modifierWord == "coming" || modifierWord == "upcoming" {
    modifier = "next"
}

// Calculate base offset
daysOffset := kronos.GetDaysToWeekday(refDate, weekday, modPtr)

// Apply extended modifiers
switch extendedModifier {
case "after next":  daysOffset += 7  // One more week
case "before last": daysOffset -= 7  // One week earlier
}
```

## Effort

**Priority**: LOW-MEDIUM (nice to have)  
**Complexity**: MEDIUM  
**Effort**: 2-3 hours

## Acceptance Criteria

- [ ] Parse "after next" patterns (2 weeks forward)
- [ ] Parse "before last" patterns (2 weeks back)
- [ ] Parse "coming/upcoming" as "next" synonyms
- [ ] 10 test cases pass

---

**Labels**: `enhancement`, `weekday`, `low-priority`
**Milestone**: Future
