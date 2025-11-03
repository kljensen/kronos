# Kronos Edge Case Behavior Documentation

This document describes how Kronos handles various edge cases and invalid inputs, documenting behavioral differences from strict parsing.

## Summary

Kronos is a **lenient parser** designed to extract dates from natural text. This means it prioritizes finding valid date information even in imperfect input, rather than strictly rejecting all malformed data.

## Key Behaviors

### ✅ What Kronos Validates (Strict Checking)

1. **Invalid Dates**: Properly rejects impossible dates
   - February 30, 2020 ❌
   - February 29, 2021 (non-leap year) ❌
   - April 31, 2020 ❌
   - June 31, 2020 ❌
   - November 31, 2020 ❌
   - September 31, 2020 ❌
   - Month 13 (2020-13-15) ❌
   - Day 32 (2020-01-32) ❌
   - Day 0 (2020-01-00) ❌
   - Month 0 (2020-00-15) ❌

2. **Leap Year Rules**: Correctly validates all leap year edge cases
   - February 29, 2000 ✅ (divisible by 400)
   - February 29, 1900 ❌ (divisible by 100 but not 400)
   - February 29, 2004 ✅ (divisible by 4)
   - February 29, 2024 ✅ (divisible by 4)
   - February 29, 2100 ❌ (divisible by 100 but not 400)

3. **Invalid Time Values**: Rejects impossible time components
   - Hour 25 (25:00) ❌
   - Hour 24 with minutes (24:01) ❌
   - Minute 60 (12:60) ❌
   - Second 60 (12:30:60) ❌
   - 14 PM ❌
   - 13 PM ❌
   - 24:00:00 ❌

4. **Boundary Times**: Correctly handles valid boundary conditions
   - 00:00:00 (midnight) ✅
   - 23:59:59 (last second of day) ✅
   - 12:00:00 (noon) ✅

5. **Malformed Input**: Rejects garbage input
   - Empty strings ❌
   - Whitespace only ❌
   - Random characters ❌
   - Special characters only ❌
   - Invalid separators (2020@01@15) ❌

6. **Robustness**: Never panics on any input
   - Binary data ✅ (returns nil)
   - Emoji ✅ (returns nil)
   - SQL injection attempts ✅ (returns nil)
   - XSS attempts ✅ (returns nil)
   - Extremely long strings ✅ (returns nil)

7. **Large Offsets**: Handles extreme relative dates without overflow
   - 1000 years ago ✅
   - 10000 days ago ✅
   - 100000 hours ago ✅

### 🟡 Lenient Behaviors (Not Errors)

These behaviors reflect Kronos' design as a **forgiving parser** that extracts what it can:

1. **Negative Signs in Times**: `-1:30` → parses as `1:30`
   - The minus sign is ignored, and the time is extracted
   - **Acceptable**: Helps with text like "arrived at -1:30 delay"

2. **Mixed Separators**: `2020-01/15` → parses successfully
   - Kronos can handle inconsistent date separators
   - **Acceptable**: Real-world text often has formatting inconsistencies

3. **Partial Dates**: Extracts valid portions from partial input
   - `March 45` → parses just `March` (ignores invalid day)
   - `15 of Marchtember` → parses `15 of March` (ignores trailing text)
   - **Acceptable**: Useful for extracting dates from imperfect text

4. **Year 0000**: `0000-01-01` → treated as year 2000
   - Two-digit year conversion applies
   - **Acceptable**: Follows year inference rules

### ❌ Known Limitations (Not Currently Supported)

1. **Zero Offsets**: Not recognized by the parser
   - `0 days ago` ❌
   - `0 hours ago` ❌
   - `0 minutes ago` ❌
   - **Reason**: Pattern matcher requires non-zero numbers

2. **Leap Seconds**: Not supported
   - `23:59:60` ❌
   - **Reason**: Go's time package doesn't support leap seconds

## Test Coverage

The edge case test suite includes:

- **12 invalid date tests**: Various impossible dates
- **8 leap year tests**: All leap year edge cases
- **6 invalid time tests**: Impossible time values
- **4 boundary time tests**: Midnight, noon, etc.
- **8 malformed input tests**: Empty, whitespace, garbage
- **3 extreme date tests**: Very old/future years
- **4 large offset tests**: Extreme relative dates
- **17 panic safety tests**: Binary, emoji, injection attempts

**Total**: 60+ edge case tests covering all categories from issue #88

## Comparison to Other Parsers

### Kronos vs Strict Parsers

**Kronos** (lenient approach):
- Extracts dates from natural text
- Handles inconsistent formatting
- Useful for real-world data extraction

**Strict Parsers** (validation approach):
- Reject any malformed input
- Require exact format matching
- Useful for data validation

### When to Use Kronos

✅ Extracting dates from emails, logs, documents
✅ Processing user-entered text
✅ Handling dates in various formats
✅ Real-world data with inconsistencies

❌ Strict data validation
❌ When you need to reject partial matches
❌ When format ambiguity is unacceptable

## Testing Methodology

All edge cases are tested with:
1. **Positive cases**: Verify valid input parses correctly
2. **Negative cases**: Verify invalid input is rejected
3. **Boundary cases**: Test limits of valid ranges
4. **Panic safety**: Ensure no panics on any input

Tests follow existing patterns and use the casual configuration to match real-world usage.
