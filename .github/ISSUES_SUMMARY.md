# Date Parser Enhancement Issues - Summary

This document summarizes 7 well-researched GitHub issues for improving the English date parsing library based on comprehensive analysis of 9 major date parsing libraries.

## 📊 Overview

| Priority | Issues | Total Effort | Impact |
|----------|--------|--------------|--------|
| **HIGH** (Critical) | 2 | 3.5-4.5 hrs | Fixes + Core features |
| **MEDIUM-HIGH** | 2 | 6-12 hrs | High-value additions |
| **MEDIUM** | 2 | 4-8 hrs | Quality improvements |
| **LOW-MEDIUM** | 1 | 2-3 hrs | Nice polish |
| **TOTAL** | **7** | **15.5-27.5 hrs** | **74% → 85%+ coverage** |

---

## 🔥 Priority 1: HIGH - Critical Fixes

### [ISSUE #001](ISSUE_001_this_next_last_week_month_year.md): Support "this/next/last week/month/year" Patterns
**Status**: Feature 80% implemented, needs completion
**Effort**: 2.5 hours
**Impact**: Most requested natural language pattern

**What's Missing**:
- "next/last day" and "next/last hour" not supported
- Month boundary bug (preserves day instead of resetting to 1st)

**Files**: `en/relative_date_parser.go`

**Quick Win**: Add 5 dictionary entries + 30 lines of switch cases

---

### [ISSUE #002](ISSUE_002_before_after_refiner_crash.md): 🐛 Slice Bounds Panic in "before/after" Refiners
**Status**: Runtime crash - CRITICAL
**Effort**: 1-2 hours
**Impact**: 11+ test cases SKIPped, crashes on user input

**Root Cause**: Unsafe string slicing without bounds checking

**Fix**: Replace `context.Text()[startIdx:endIdx]` with `kronos.SafeSlice()` (2 lines changed)

**Files**: `en/refiner_merge_relative_follow.go:56`, `en/refiner_merge_relative_after.go:56`

**Severity**: Production crash bug - fix immediately

---

## 📈 Priority 2: MEDIUM-HIGH - High Value Additions

### [ISSUE #003](ISSUE_003_weekend_support.md): Add "weekend" Support
**Effort**: 4-8 hours
**Impact**: Extremely common casual reference

**What's Needed**:
- "this weekend", "next weekend", "last weekend" ✅ (basic support exists)
- "N weekends ago/from now" ❌ (NEW)
- Span support (Saturday-Monday range) ❌ (future enhancement)

**Implementation**: Extend existing weekday parser with counting logic

---

### [ISSUE #004](ISSUE_004_compound_relative_expressions.md): Compound Relative Expressions with Commas
**Effort**: 2-4 hours
**Impact**: Natural language quality improvement

**What's Missing**:
- "1 year, 2 months ago" (comma-separated)
- "1 year, 1 month, 1 week, 1 day, 1 hour and 1 minute ago"

**Current**: Space-separated works, commas don't

**Fix Options**:
- A: Update regex to accept commas/and
- B: Normalize (strip commas) before parsing ⭐ Simpler

---

## 📝 Priority 3: MEDIUM - Nice to Have

### [ISSUE #005](ISSUE_005_word_number_support.md): Add Word Number Support
**Effort**: 2-4 hours
**Impact**: Better NLP, text parsing

**What's Needed**: "three days ago", "two weeks from now", "fifteen minutes"

**Current State**: Infrastructure exists, dictionary entries missing

**Solution**: Add 20 lines to `NumberWordDictionary` in `en/constants.go`

**Scope**: 1-20 + tens (twenty, thirty, ..., ninety) covers 90%+ usage

---

### [ISSUE #006](ISSUE_006_advanced_weekday_patterns.md): Advanced Weekday Patterns
**Effort**: 2-3 hours
**Impact**: Natural language variations

**What's Needed**:
- "the Monday after next" (2 weeks forward)
- "Tuesday before last" (2 weeks back)
- "coming Monday", "this upcoming Friday" (synonyms)

**Implementation**: Extend weekday parser regex, add offset logic

---

## 🔮 Priority 4: LOW-MEDIUM - Polish

### [ISSUE #007](ISSUE_007_quarter_support.md): Business Quarter Support (Q1-Q4)
**Effort**: 2-4 hours
**Impact**: Business/financial contexts

**What's Needed**:
- "Q1 2024", "Q2", "1st quarter 2023"
- Fix "this quarter" (missing switch case - 15 min quick win!)

**Infrastructure**: Quarter support already exists for "next quarter"

**Implementation**: New quarter parser + fix one switch case

---

## 🎯 Recommended Implementation Order

### Sprint 1: Quick Wins (6 hours, high impact)
1. **ISSUE #002** (1-2 hrs) - Fix crash bug first! 🚨
2. **ISSUE #001** (2.5 hrs) - Complete partially implemented feature
3. **ISSUE #005** (2-4 hrs) - Infrastructure ready, just add dictionary entries

### Sprint 2: High-Value Features (8-12 hours)
4. **ISSUE #003** (4-8 hrs) - Weekend support (common request)
5. **ISSUE #004** (2-4 hrs) - Polish compound expressions

### Sprint 3: Nice to Have (4-7 hours)
6. **ISSUE #007** (2-4 hrs) - Business quarters
7. **ISSUE #006** (2-3 hrs) - Advanced weekday patterns

---

## 📚 Research Methodology

All issues are based on comprehensive analysis of test suites from:

1. **Chrono** (JavaScript) - 10K+ stars, battle-tested
2. **Chronic** (Ruby) - Natural language focus
3. **dateparser** (Python) - 98% test coverage, 7K+ stars
4. **go-dateparser** (Go) - Direct competitor
5. **Natty** (Java) - ANTLR-based parser
6. **parsedatetime** (Python) - Human-readable dates
7. **any-date-parser** (JavaScript) - Wide format support
8. **moment.js** (JavaScript) - Industry standard
9. **date-fns** (TypeScript) - Modern functional approach

**Test files analyzed**: 1.6 MB, ~31,000 lines of tests across 54 files

---

## 📋 Issue Quality Standards

Each issue includes:
- ✅ **Standalone, PM-readable format** - No deep technical knowledge required
- ✅ **Problem statement with examples** - Clear what's broken/missing
- ✅ **Current vs expected behavior** - Side-by-side comparison
- ✅ **10-15 specific test cases** - Concrete examples with expected outputs
- ✅ **Implementation suggestions** - File paths, regex patterns, logic
- ✅ **Code examples** - Before/after code snippets
- ✅ **References to other libraries** - How competitors solve it
- ✅ **Research citations** - Line numbers in test files
- ✅ **Priority and complexity estimates** - Effort breakdown
- ✅ **Acceptance criteria** - Clear done conditions

---

## 📈 Expected Outcome

**Before**: 74% pattern coverage (excellent for English-only parser)

**After implementing all 7 issues**: 85%+ coverage

### Coverage Improvements by Category

| Category | Before | After | Improvement |
|----------|--------|-------|-------------|
| Casual References | 58% | 85% | +27% |
| Relative Time | 61% | 90% | +29% |
| Weekday Patterns | 73% | 95% | +22% |
| Month/Year Patterns | 67% | 85% | +18% |
| Slash Formats | 100% | 100% | ✅ Perfect |
| Time Expressions | 100% | 100% | ✅ Perfect |

---

## 🔗 Quick Links

- [Issue #001 - this/next/last week/month/year](ISSUE_001_this_next_last_week_month_year.md) - HIGH ⚡
- [Issue #002 - before/after crash bug](ISSUE_002_before_after_refiner_crash.md) - CRITICAL 🚨
- [Issue #003 - weekend support](ISSUE_003_weekend_support.md) - MEDIUM-HIGH 📅
- [Issue #004 - compound relative expressions](ISSUE_004_compound_relative_expressions.md) - MEDIUM
- [Issue #005 - word number support](ISSUE_005_word_number_support.md) - MEDIUM
- [Issue #006 - advanced weekday patterns](ISSUE_006_advanced_weekday_patterns.md) - LOW-MEDIUM
- [Issue #007 - quarter support](ISSUE_007_quarter_support.md) - MEDIUM 💼

---

## 🎓 Key Insights from Research

1. **Most libraries prioritize natural language** over format coverage
2. **"this/next/last" modifiers are table stakes** - all 9 libraries support them
3. **Word numbers (one, two, three) appear in 6+ libraries** - expected feature
4. **Weekend references are surprisingly common** - 3 libraries have dedicated support
5. **Comma-separated compounds are standard** in Python/Go ecosystems
6. **Quarter support is business-critical** - appears in 2 libraries, requested feature
7. **Slice safety is critical** - all mature libraries use bounds checking

---

## 📞 Contact & Attribution

Issues created by: Autonomous analysis using pragmatic-go-dev and general-purpose agents

Analysis date: 2025-11-04

Test corpus: 9 libraries, 54 files, 1.6 MB, ~31,000 lines

---

**Note**: All issues are production-ready and can be copied directly to GitHub. Each is standalone and requires no additional context beyond what's included in the issue itself.
