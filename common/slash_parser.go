package common

import (
	"regexp"
	"strconv"
	"strings"

	kronos "github.com/kljensen/kronos"
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
func (p *SlashDateFormatParser) Pattern(context *kronos.ParsingContext) *regexp.Regexp {
	return slashPattern
}

// Extract implements Parser.Extract.
func (p *SlashDateFormatParser) Extract(context *kronos.ParsingContext, match []string) interface{} {
	// Get the match boundaries
	fullMatch := match[0]
	opening := match[slashOpeningGroup]
	ending := match[slashEndingGroup]

	// The actual date text is the full match without the opening/ending boundaries
	matchText := fullMatch[len(opening) : len(fullMatch)-len(ending)]

	// If opening boundary is empty (matched ^), check if there's actually a digit before
	// this match in the original text. This prevents matching "3/31/2018" in "13/31/2018".
	if len(opening) == 0 {
		matchStartInText := strings.Index(context.Text(), fullMatch)
		if matchStartInText > 0 {
			prevChar := context.Text()[matchStartInText-1]
			if prevChar >= '0' && prevChar <= '9' {
				return nil
			}
		}
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

	// Parse month and day
	month, _ := strconv.Atoi(match[p.groupNumberMonth])
	day, _ := strconv.Atoi(match[p.groupNumberDay])

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
	components := kronos.NewParsingComponents(context.Reference(), map[kronos.Component]int{
		kronos.ComponentDay:   day,
		kronos.ComponentMonth: month,
	})

	// Handle year
	if match[slashYearGroup] != "" {
		rawYear, _ := strconv.Atoi(match[slashYearGroup])
		year := kronos.FindMostLikelyADYear(rawYear)
		components.Assign(kronos.ComponentYear, year)
	} else {
		year := kronos.FindYearClosestToRef(context.RefDate(), day, month)
		components.Imply(kronos.ComponentYear, year)
	}

	// If there's a boundary, return ParsingResultWithBoundary
	// Otherwise, return ParsingComponents for chrono.go to handle
	if len(opening) > 0 {
		return &kronos.ParsingResultWithBoundary{
			Components:   components,
			AdjustedText: matchText,
			BoundaryLen:  len(opening),
		}
	}

	// No boundary, return components directly
	return components
}
