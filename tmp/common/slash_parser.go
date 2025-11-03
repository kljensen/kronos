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
	text := context.Text()
	fullMatch := match[0]
	opening := match[slashOpeningGroup]
	ending := match[slashEndingGroup]

	// Calculate actual index by finding where the match starts in the text
	index := strings.Index(text, fullMatch)
	if index < 0 {
		return nil
	}
	index += len(opening)
	indexEnd := index + len(fullMatch) - len(opening) - len(ending)

	// Check for digits before/after to avoid matching parts of larger numbers
	if index > 0 {
		textBefore := text[:index]
		if regexp.MustCompile(`\d/?$`).MatchString(textBefore) {
			return nil
		}
	}

	if indexEnd < len(text) {
		textAfter := text[indexEnd:]
		if regexp.MustCompile(`^/?\d`).MatchString(textAfter) {
			return nil
		}
	}

	matchText := text[index:indexEnd]

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
	components := map[kronos.Component]int{
		kronos.ComponentDay:   day,
		kronos.ComponentMonth: month,
	}

	// Handle year
	var result *kronos.ParsingResult
	if match[slashYearGroup] != "" {
		rawYear, _ := strconv.Atoi(match[slashYearGroup])
		year := kronos.FindMostLikelyADYear(rawYear)
		components[kronos.ComponentYear] = year
		result = context.CreateParsingResult(index, matchText, components)
	} else {
		result = context.CreateParsingResult(index, matchText, components)
		year := kronos.FindYearClosestToRef(context.RefDate(), day, month)
		// Get the start components as concrete type
		startComponents := result.Start().(*kronos.ParsingComponents)
		startComponents.Imply(kronos.ComponentYear, year)
	}

	return result.AddTag("parser/SlashDateFormatParser")
}
