//nolint:staticcheck // SA1019: Must use deprecated types during transition
package parsers

import (
	"regexp"
	"strconv"
	"strings"

	kronos "github.com/kljensen/kronos"
	"github.com/kljensen/kronos/internal/helpers"
)

// SlashDateFormatParser parses date formats with slash, dot, or dash separators.
// For examples:
// - 7/10
// - 7/12/2020
// - 7.12.2020
// - 7-12-2020
type SlashDateFormatParser struct {
	littleEndian     bool
	groupNumberMonth int
	groupNumberDay   int
}

// Pattern for slash/dot/dash date formats
var slashPattern = regexp.MustCompile(
	`([^\d]|^)` +
		`([0-3]{0,1}[0-9]{1})[\/\.\-]([0-3]{0,1}[0-9]{1})` +
		`(?:[\/\.\-]([0-9]{4}|[0-9]{2}))?` +
		`(\W|$)`,
)

const (
	slashOpeningGroup       = 1
	slashEndingGroup        = 5
	slashFirstNumbersGroup  = 2
	slashSecondNumbersGroup = 3
	slashYearGroup          = 4
)

// NewSlashDateFormatParser creates a new slash date format parser.
// littleEndian determines whether to interpret dates as DD/MM (true) or MM/DD (false).
func NewSlashDateFormatParser(littleEndian bool) *SlashDateFormatParser {
	groupMonth := slashFirstNumbersGroup
	groupDay := slashSecondNumbersGroup

	if littleEndian {
		groupMonth = slashSecondNumbersGroup
		groupDay = slashFirstNumbersGroup
	}

	return &SlashDateFormatParser{
		littleEndian:     littleEndian,
		groupNumberMonth: groupMonth,
		groupNumberDay:   groupDay,
	}
}

// Pattern implements Parser.Pattern.
func (p *SlashDateFormatParser) Pattern(context *kronos.InternalParsingContext) *regexp.Regexp {
	return slashPattern
}

// Extract implements Parser.Extract.
func (p *SlashDateFormatParser) Extract(context *kronos.InternalParsingContext, match []string) interface{} {
	// Get the match boundaries
	fullMatch := match[0]
	opening := match[slashOpeningGroup]
	ending := match[slashEndingGroup]

	// The actual date text is the full match without the opening/ending boundaries
	matchText := fullMatch[len(opening) : len(fullMatch)-len(ending)]

	// Determine month and day group numbers based on DateOrder from context.
	// Start with the construction-time default (based on littleEndian parameter).
	groupMonth := p.groupNumberMonth
	groupDay := p.groupNumberDay

	// Allow DateOrder from context option to override the construction-time default.
	// This enables the new builder API (DateOrder setting) to work correctly.
	dateOrder := context.Option().DateOrder

	// DateOrderDMY means day first (little-endian, DD/MM)
	// DateOrderMDY means month first (middle-endian, MM/DD)
	// DateOrderYMD falls back to construction-time default (not applicable for slash dates)
	if dateOrder == kronos.DateOrderDMY {
		groupMonth = slashSecondNumbersGroup
		groupDay = slashFirstNumbersGroup
	} else if dateOrder == kronos.DateOrderMDY && !p.littleEndian {
		// Only override to MDY if the construction-time default was also MDY.
		// This maintains backward compatibility with the old API.
		groupMonth = slashFirstNumbersGroup
		groupDay = slashSecondNumbersGroup
	}
	// Otherwise (DateOrderYMD or littleEndian with MDY default), use construction-time default values

	// Check if this match is part of an invalid 3-component date like "4/13/1"
	// This can happen in two ways:
	// 1. Opening is empty and preceded by "digit/"
	// 2. Opening is "/" and this is the second part of an invalid date
	matchStartInText := strings.Index(context.Text(), fullMatch)

	if len(opening) == 0 {
		// Case 1: Opening boundary is empty (matched ^), check if there's a digit before
		if matchStartInText > 0 {
			prevChar := context.Text()[matchStartInText-1]
			if prevChar >= '0' && prevChar <= '9' {
				return nil
			}
			// Check for pattern like "/13/1" where this is part of an invalid date
			if prevChar == '/' || prevChar == '.' || prevChar == '-' {
				// Look further back for digits
				if matchStartInText > 1 {
					beforeSep := context.Text()[matchStartInText-2]
					if beforeSep >= '0' && beforeSep <= '9' {
						// This looks like part of a 3-component date
						// Check if the current match has no year, which would make it invalid
						if match[slashYearGroup] == "" {
							return nil
						}
					}
				}
			}
		}
	} else if opening == "/" || opening == "." || opening == "-" {
		// Case 2: Opening is a separator, check if this is part of an invalid 3-component date
		// like "/13/1" in "4/13/1"
		if matchStartInText > 0 {
			beforeSep := context.Text()[matchStartInText-1]
			if beforeSep >= '0' && beforeSep <= '9' {
				// This is part of a date like "4/13/1"
				// Check if the current match has no year, which would make it invalid
				if match[slashYearGroup] == "" {
					return nil
				}
			}
		}
	}

	// Check if there's an invalid year pattern after the date (e.g., "4/13/1")
	// If the year wasn't matched but the ending is a slash, it means there's a third
	// component that didn't match the year pattern (invalid format like single digit)
	if match[slashYearGroup] == "" && (ending == "/" || ending == "." || ending == "-") {
		return nil
	}

	// Skip version numbers like "1.12" or "1.12.12"
	if regexp.MustCompile(`^\d\.\d$`).MatchString(matchText) ||
		regexp.MustCompile(`^\d\.\d{1,2}\.\d{1,2}\s*$`).MatchString(matchText) {
		return nil
	}

	// Dates without year must use slash (not dot or dash)
	if match[slashYearGroup] == "" && !strings.Contains(matchText, "/") {
		return nil
	}

	// Parse month and day using the determined group numbers
	month, err := strconv.Atoi(match[groupMonth])
	if err != nil {
		return nil
	}
	day, err := strconv.Atoi(match[groupDay])
	if err != nil {
		return nil
	}

	// Validate and swap if needed
	if month < 1 || month > 12 {
		if month > 12 {
			// Try swapping
			if day >= 1 && day <= 12 && month <= 31 {
				month, day = day, month
			} else {
				return nil
			}
		} else {
			return nil
		}
	}

	if day < 1 || day > 31 {
		return nil
	}

	// Create the parsing components
	components := helpers.NewParsingComponents(context.Reference(), map[kronos.Component]int{
		kronos.ComponentDay:   day,
		kronos.ComponentMonth: month,
	})
	components.AddTag("parser/SlashDateFormatParser")
	// Set period to day-level since this parser extracts a specific date
	components.SetPeriod(kronos.PeriodDay)

	// Handle year
	if match[slashYearGroup] != "" {
		rawYear, err := strconv.Atoi(match[slashYearGroup])
		if err != nil {
			return nil
		}
		year := helpers.FindMostLikelyADYear(rawYear)
		components.Assign(kronos.ComponentYear, year)
	} else {
		// Use preference-aware year selection, which handles Feb 29 leap years
		year := helpers.FindYearClosestToRefWithPreference(context.RefDate(), day, month, context.Option().PreferDatesFrom)
		components.Imply(kronos.ComponentYear, year)
	}

	// If there's a boundary, return ParsingResultWithBoundary
	// Otherwise, return ParsingComponents for chrono.go to handle
	if len(opening) > 0 {
		return &kronos.InternalParsingResultWithBoundary{
			Components:         components,
			AdjustedText:       matchText,
			BoundaryLen:        len(opening),
			IncludeBoundaryIdx: true, // Index should point past the boundary
		}
	}

	// No boundary, return components directly
	return components
}
