# Feature Request: Weekend Reference Support

## Summary

Add support for parsing "weekend" references (this weekend, next weekend, last weekend, N weekends ago/from now) - one of the most common casual date expressions in everyday conversation.

## Problem Statement

Users frequently reference weekends when scheduling and discussing plans. Our library currently lacks support for these natural patterns despite having a partially implemented foundation in the weekday parser.

### What Users Say
- "Let's meet this weekend"
- "I was out of town last weekend"
- "We're planning a trip next weekend"
- "2 weekends ago we went hiking"
- "Can we reschedule to 1 weekend from now?"

### Current State
The library does **not** fully recognize "weekend" as a date reference:

```go
Parse("this weekend", refDate)       // → Partial support (returns Saturday)
Parse("last weekend", refDate)       // → Partial support (returns Sunday)
Parse("2 weekends ago", refDate)     // → NO MATCH
Parse("1 weekend from now", refDate) // → NO MATCH
```

## Current Implementation Analysis

### What Exists Now

Looking at `/workspace/en/weekday_parser.go:109-118`:

```go
} else if weekdayWord == "weekend" {
    // "This/next weekend" means the coming Saturday,
    // "last weekend" means last Sunday
    if modifier == "last" {
        weekday = kronos.WeekdaySunday
    } else {
        weekday = kronos.WeekdaySaturday
    }
}
```

**Current capabilities:**
- ✅ Basic "this/next/last weekend" works
- ✅ "last weekend" → Sunday (sensible default)
- ✅ "this/next weekend" → Saturday (forward-looking)
- ❌ No "N weekends ago/from now"
- ❌ No span/range support
- ❌ Limited context-dependent behavior

## Research: How Other Libraries Handle Weekends

### 1. Chrono (JavaScript/TypeScript) - 10K+ GitHub stars

**Test File**: `/workspace/tmp/test-cases/chrono_en_weekday.test.ts`

```typescript
// Lines 222-238: Weekend support
test("Test - Weekday expressions with weekend", () => {
    // Reference: Friday, October 18, 2024 at 12:00
    testSingleCase(chrono, "last weekend", new Date(2024, 10-1, 18, 12), (result) => {
        expect(result.start.get("year")).toBe(2024);
        expect(result.start.get("month")).toBe(10);
        expect(result.start.get("day")).toBe(13);  // Sunday, Oct 13
        expect(result.start.get("weekday")).toBe(0); // Sunday
    });

    testSingleCase(chrono, "this weekend", new Date(2024, 10-1, 18, 12), (result) => {
        expect(result.start.get("year")).toBe(2024);
        expect(result.start.get("month")).toBe(10);
        expect(result.start.get("day")).toBe(19);  // Saturday, Oct 19
        expect(result.start.get("weekday")).toBe(6); // Saturday
    });

    testSingleCase(chrono, "next weekend", new Date(2024, 10-1, 18, 12), (result) => {
        expect(result.start.get("year")).toBe(2024);
        expect(result.start.get("month")).toBe(10);
        expect(result.start.get("day")).toBe(26);  // Saturday, Oct 26
        expect(result.start.get("weekday")).toBe(6); // Saturday
    });
});
```

**Key behaviors:**
- "last weekend" → Previous Sunday
- "this weekend" → Upcoming Saturday (same week if before weekend, next week if currently weekend)
- "next weekend" → Saturday of next week

### 2. Chronic (Ruby) - Natural language date parser

**Test File**: `/workspace/tmp/test-cases/chronic_test_parsing.rb`

```ruby
# Lines 706-715: Weekend with context
def test_handle_r_weekend
  time = parse_now("this weekend", :context => :future)
  assert_equal Time.local(2006, 8, 20), time  # Sunday

  time = parse_now("this weekend", :context => :past)
  assert_equal Time.local(2006, 8, 13), time  # Previous Sunday

  time = parse_now("last weekend")
  assert_equal Time.local(2006, 8, 13), time  # Sunday
end

# Lines 945-946: Counting weekends
time = parse_now("2 weekends ago")
assert_equal Time.local(2006, 8, 5), time  # Saturday 2 weeks back

# Lines 980-984: Future weekends
time = parse_now("1 weekend from now")
assert_equal Time.local(2006, 8, 19), time  # Saturday

time = parse_now("2 weekends from now")
assert_equal Time.local(2006, 8, 26), time  # Saturday
```

**Key behaviors:**
- Context-aware: "this weekend" can be past or future based on context
- Counting support: "N weekends ago/from now"
- Span mode: Can return weekend as a range (Saturday 19:00 - Monday 00:00)

### 3. Comparison Summary

| Feature | Kronos (Current) | Chrono | Chronic |
|---------|------------------|--------|---------|
| "this weekend" | ✅ Saturday | ✅ Saturday | ✅ Context-aware |
| "next weekend" | ✅ Saturday | ✅ Saturday | ✅ Saturday |
| "last weekend" | ✅ Sunday | ✅ Sunday | ✅ Sunday |
| "N weekends ago" | ❌ | ❌ | ✅ |
| "N weekends from now" | ❌ | ❌ | ✅ |
| Span/range mode | ❌ | ❌ | ✅ Sat-Mon |
| Context switching | ❌ | ⚠️ Limited | ✅ Full |

## Design Decisions

### 1. What day should "weekend" resolve to?

**Decision: Context-dependent (matches observed behavior)**

- **"this/next weekend"** → Saturday (forward-looking planning)
- **"last weekend"** → Sunday (looking back at weekend that passed)
- **"N weekends from now"** → Saturday (planning future)
- **"N weekends ago"** → Saturday (counting back)

**Rationale**:
- Asymmetric but intuitive
- Matches both Chrono and Chronic behavior
- Saturday = start of weekend (planning perspective)
- Sunday = end of weekend (retrospective perspective)

### 2. What is the "weekend" span?

**Decision: Weekend represents Saturday-Sunday, resolve to single day by default**

When parsed as a point in time:
- Returns Saturday or Sunday depending on context (see above)

When parsed as a span (future feature):
- Could return range: Saturday 00:00 - Monday 00:00
- Set `Period = PeriodWeekend` (new period type)

**Rationale**:
- Consistent with how we handle "March" (returns March 1, not a range)
- Span support can be added later without breaking changes

### 3. How should "this weekend" behave based on current day?

**Decision: Simple forward-looking logic**

| Current Day | "this weekend" | "next weekend" | "last weekend" |
|-------------|----------------|----------------|----------------|
| Monday-Friday | Upcoming Saturday | Saturday after next | Previous Sunday |
| Saturday | Today (Saturday) | Saturday next week | Previous Sunday |
| Sunday | Today (Sunday) | Saturday next week | Previous Sunday |

**Rationale**: Keeps implementation simple, matches user intuition

## Proposed Test Cases

Based on reference implementations, here are **15 test cases**:

### Basic Modifiers

```go
// Reference: Friday, Oct 18, 2024 12:00
{
    name:     "last weekend from Friday",
    text:     "last weekend",
    refDate:  time.Date(2024, 10, 18, 12, 0, 0, 0, time.UTC),
    expected: time.Date(2024, 10, 13, 12, 0, 0, 0, time.UTC), // Sunday
},
{
    name:     "this weekend from Friday",
    text:     "this weekend",
    refDate:  time.Date(2024, 10, 18, 12, 0, 0, 0, time.UTC),
    expected: time.Date(2024, 10, 19, 12, 0, 0, 0, time.UTC), // Saturday
},
{
    name:     "next weekend from Friday",
    text:     "next weekend",
    refDate:  time.Date(2024, 10, 18, 12, 0, 0, 0, time.UTC),
    expected: time.Date(2024, 10, 26, 12, 0, 0, 0, time.UTC), // Saturday
},

// Reference: Wednesday, Aug 16, 2006 14:00
{
    name:     "this weekend from Wednesday",
    text:     "this weekend",
    refDate:  time.Date(2006, 8, 16, 14, 0, 0, 0, time.UTC),
    expected: time.Date(2006, 8, 19, 14, 0, 0, 0, time.UTC), // Saturday
},
{
    name:     "last weekend from Wednesday",
    text:     "last weekend",
    refDate:  time.Date(2006, 8, 16, 14, 0, 0, 0, time.UTC),
    expected: time.Date(2006, 8, 13, 14, 0, 0, 0, time.UTC), // Sunday
},
```

### During Weekend

```go
// Reference: Saturday, Aug 19, 2006 14:00
{
    name:     "this weekend from Saturday",
    text:     "this weekend",
    refDate:  time.Date(2006, 8, 19, 14, 0, 0, 0, time.UTC),
    expected: time.Date(2006, 8, 19, 14, 0, 0, 0, time.UTC), // Today (Saturday)
},

// Reference: Sunday, Aug 20, 2006 14:00
{
    name:     "this weekend from Sunday",
    text:     "this weekend",
    refDate:  time.Date(2006, 8, 20, 14, 0, 0, 0, time.UTC),
    expected: time.Date(2006, 8, 20, 14, 0, 0, 0, time.UTC), // Today (Sunday)
},
{
    name:     "next weekend from Sunday",
    text:     "next weekend",
    refDate:  time.Date(2006, 8, 20, 14, 0, 0, 0, time.UTC),
    expected: time.Date(2006, 8, 26, 14, 0, 0, 0, time.UTC), // Saturday
},
```

### Relative Counting (NEW FEATURE)

```go
// Reference: Wednesday, Aug 16, 2006 14:00
{
    name:     "2 weekends ago",
    text:     "2 weekends ago",
    refDate:  time.Date(2006, 8, 16, 14, 0, 0, 0, time.UTC),
    expected: time.Date(2006, 8, 5, 14, 0, 0, 0, time.UTC), // Saturday 2 weeks back
},
{
    name:     "1 weekend from now",
    text:     "1 weekend from now",
    refDate:  time.Date(2006, 8, 16, 14, 0, 0, 0, time.UTC),
    expected: time.Date(2006, 8, 19, 14, 0, 0, 0, time.UTC), // This Saturday
},
{
    name:     "2 weekends from now",
    text:     "2 weekends from now",
    refDate:  time.Date(2006, 8, 16, 14, 0, 0, 0, time.UTC),
    expected: time.Date(2006, 8, 26, 14, 0, 0, 0, time.UTC), // Saturday in 2 weeks
},
{
    name:     "3 weekends ago",
    text:     "3 weekends ago",
    refDate:  time.Date(2006, 8, 16, 14, 0, 0, 0, time.UTC),
    expected: time.Date(2006, 7, 29, 14, 0, 0, 0, time.UTC), // Saturday 3 weeks back
},
```

### Edge Cases

```go
// Just after weekend ends
{
    name:     "this weekend from Monday",
    text:     "this weekend",
    refDate:  time.Date(2006, 8, 21, 9, 0, 0, 0, time.UTC), // Monday
    expected: time.Date(2006, 8, 26, 9, 0, 0, 0, time.UTC), // Next Saturday
},
{
    name:     "last weekend from Monday",
    text:     "last weekend",
    refDate:  time.Date(2006, 8, 21, 9, 0, 0, 0, time.UTC),
    expected: time.Date(2006, 8, 20, 9, 0, 0, 0, time.UTC), // Yesterday (Sunday)
},

// Year boundary
{
    name:     "next weekend crosses year",
    text:     "next weekend",
    refDate:  time.Date(2023, 12, 29, 12, 0, 0, 0, time.UTC), // Friday
    expected: time.Date(2023, 12, 30, 12, 0, 0, 0, time.UTC), // Saturday (this week)
},
```

## Implementation Suggestions

### Recommended Approach: Extend ENWeekdayParser

The weekday parser already handles "weekend" partially. Extend it to support counting.

**File**: `/workspace/en/weekday_parser.go`

### 1. Update Regex Pattern (line 25-31)

**Current:**
```go
pattern := `(?:(?:,|\(|（)\s*)?` +
    `(?:on\s*?)?` +
    `(?:(this|last|past|next)\s*)?` +
    `(` + WeekdayPattern + `|weekend|weekday)` +
    `(?:\s*(?:,|\)|）))?` +
    `(?:\s*(this|last|past|next)\s*week)?` +
    `(?:\s|$|\b)`
```

**Proposed:**
```go
pattern := `(?:(?:,|\(|（)\s*)?` +
    `(?:on\s*?)?` +
    `(?:(\d+)\s+)?` +                           // NEW: capture count "2 weekends"
    `(?:(this|last|past|next)\s*)?` +
    `(` + WeekdayPattern + `|weekend|weekday)` +
    `(?:\s+(ago|from\s+now))?` +                // NEW: direction for counting
    `(?:\s*(?:,|\)|）))?` +
    `(?:\s*(this|last|past|next)\s*week)?` +
    `(?:\s|$|\b)`
```

**Capture groups:**
1. Count (optional): "2" in "2 weekends ago"
2. Modifier: "this"/"last"/"next"
3. Weekday/weekend: "weekend"
4. Direction (new): "ago" or "from now"
5. Postfix week modifier: "next week"

### 2. Resolution Logic (around line 109-118)

**Current logic:**
```go
} else if weekdayWord == "weekend" {
    if modifier == "last" {
        weekday = kronos.WeekdaySunday
    } else {
        weekday = kronos.WeekdaySaturday
    }
}
```

**Proposed enhancement:**
```go
} else if weekdayWord == "weekend" {
    // Parse count if present (e.g., "2 weekends ago")
    count := 1
    if match[1] != "" {
        countVal, err := strconv.Atoi(match[1])
        if err == nil && countVal > 0 {
            count = countVal
        }
    }

    // Parse direction (ago/from now)
    direction := match[4] // "ago" or "from now"
    hasDirection := (direction != "")

    // Determine target weekday
    var targetWeekday kronos.Weekday
    if modifier == "last" || direction == "ago" {
        targetWeekday = kronos.WeekdaySunday  // Looking back
    } else {
        targetWeekday = kronos.WeekdaySaturday  // Looking forward
    }

    // Calculate offset
    var daysOffset int
    if hasDirection {
        // Counting weekends: each weekend = 7 days
        if direction == "ago" {
            // Go back N weeks to Saturday
            daysOffset = -7 * count
            // Adjust to previous Saturday from reference
            daysToLastSat := int(refDate.Weekday()) + 1
            if refDate.Weekday() == time.Saturday {
                daysToLastSat = 0
            }
            daysOffset -= daysToLastSat
        } else { // "from now"
            // Go forward N weeks to Saturday
            daysToNextSat := (int(time.Saturday) - int(refDate.Weekday()) + 7) % 7
            if daysToNextSat == 0 && count > 1 {
                daysToNextSat = 7
            }
            daysOffset = daysToNextSat + 7*(count-1)
        }
    } else {
        // Standard "this/next/last weekend" logic
        modPtr := &modifier
        daysOffset = kronos.GetDaysToWeekday(refDate, targetWeekday, modPtr)
    }

    weekday = targetWeekday
    // Continue with existing offset calculation...
}
```

### 3. Alternative: Extract to Dedicated Function

For cleaner code, consider:

```go
func calculateWeekendOffset(refDate time.Time, modifier string, count int, direction string) (int, kronos.Weekday) {
    // Weekend calculation logic here
    // Returns: days offset, target weekday (Sat or Sun)
}
```

### 4. Tests Structure

**New test file**: `/workspace/en/weekend_parser_test.go` (or add to `weekday_parser_test.go`)

```go
func TestENWeekdayParser_Weekend(t *testing.T) {
    // Group 1: Basic this/next/last (ensure no regression)
    // Group 2: Weekend from different days of week
    // Group 3: N weekends ago
    // Group 4: N weekends from now
    // Group 5: Edge cases (year boundaries, during weekend)
}
```

## Complexity & Effort

**Complexity**: **SMALL-MEDIUM**

**Breakdown:**
- Regex pattern update: Simple (add 2 capture groups)
- Weekend resolution logic: Medium (date arithmetic)
- Count parsing: Simple (Atoi conversion)
- Testing: Medium (15 test cases, edge cases)

**Estimated Effort**: **4-6 hours**
- Pattern update: 30 minutes
- Logic implementation: 2-3 hours
- Test writing: 1-2 hours
- Edge case testing: 1 hour

**Risk Level**: LOW
- Changes are additive (extends existing parser)
- Basic "weekend" support continues working
- No breaking changes to API

## Priority

**MEDIUM-HIGH**

**Justification:**
- ✅ **User demand**: Extremely common casual expression
- ✅ **Competitive parity**: Both Chrono and Chronic support this
- ✅ **Quick win**: Foundation already exists, just needs extension
- ✅ **Low risk**: Additive feature, no breaking changes
- ⚠️ **Workaround exists**: Users can say "Saturday" or "Sunday" explicitly

**Recommended for**: Next minor release (v1.x)

## Acceptance Criteria

- [ ] Parse "this weekend", "next weekend", "last weekend" (verify no regression)
- [ ] Parse "N weekends ago" (N = 1-99)
- [ ] Parse "N weekends from now" (N = 1-99)
- [ ] Resolve to Saturday for future-looking references
- [ ] Resolve to Sunday for past-looking references
- [ ] Handle edge cases: during weekend, year boundaries
- [ ] All 15 test cases pass
- [ ] Backward compatible with existing weekday parsing
- [ ] Documentation updated with weekend examples
- [ ] Period tracking set correctly (`PeriodDay` or future `PeriodWeekend`)

## Future Enhancements (Out of Scope)

- **Span mode**: Return weekend as a date range (Saturday-Monday)
- **Context awareness**: "this weekend" behaves differently based on past/future context
- **Fiscal weekends**: Support for regions where weekend is Friday-Saturday
- **"long weekend"**: Detect holidays and extend weekend range

## References

### Code Files
- **Current implementation**: `/workspace/en/weekday_parser.go:109-118`
- **Current tests**: `/workspace/en/weekday_parser_test.go` (848 lines)
- **Weekday utilities**: `/workspace/weekdays.go` (`GetDaysToWeekday`)

### Research Files
- **Chrono tests**: `/workspace/tmp/test-cases/chrono_en_weekday.test.ts:222-238`
- **Chronic tests**: `/workspace/tmp/test-cases/chronic_test_parsing.rb:706-715,945-946,980-984`

### Related Types
- `kronos.WeekdaySaturday` (constant)
- `kronos.WeekdaySunday` (constant)
- `kronos.PeriodDay` (for tracking granularity)

---

**Labels**: `enhancement`, `parser`, `weekend`, `natural-language`, `medium-priority`
**Milestone**: v1.next
**Estimated Effort**: 4-6 hours
