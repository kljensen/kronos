//nolint:staticcheck // SA1019: Must use deprecated types during transition
package parsers

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/internal/en/data"
	"github.com/kljensen/kronos/internal/helpers"
	"github.com/kljensen/kronos/internal/parsing"
)

// ENWeekdayParser parses standalone weekdays with optional modifiers:
// "Monday", "on Friday", "this Tuesday", "next Wednesday", "last Thursday"
// Also handles special cases: "weekend", "weekday"
type ENWeekdayParser struct {
	*parsing.AbstractParserWithWordBoundary
}

// NewENWeekdayParser creates a new parser for weekday expressions
func NewENWeekdayParser() *ENWeekdayParser {
	parser := &ENWeekdayParser{}

	parser.AbstractParserWithWordBoundary = parsing.NewAbstractParserWithWordBoundary(
		func(context *kronos.InternalParsingContext) *regexp.Regexp {
			// Unicode parentheses (（ and ）) are literal characters, not escaped
			pattern := `(?:(?:,|\(|（)\s*)?` +
				`(?:on\s*?)?` +
				`(?:([0-9]+)\s+)?` + // Capture count for "N weekends ago"
				`(?:(this|last|past|next)\s*)?` +
				`(` + data.WeekdayPattern + `|weekend|weekday)s?` + // Allow optional plural "weekends"
				`(?:\s+(ago|from\s+now))?` + // Direction for counting
				`(?:\s*(?:,|\)|）))?` +
				`(?:\s*(this|last|past|next)\s*week)?` +
				`(?:\s|$|\b)`

			return regexp.MustCompile("(?i)" + pattern)
		},
		func(context *kronos.InternalParsingContext, match []string) any {
			if len(match) < 4 {
				return nil
			}

			// Capture groups:
			// [1] = count (e.g., "2" in "2 weekends ago")
			// [2] = modifier (this/last/past/next)
			// [3] = weekday/weekend/weekday word
			// [4] = direction (ago/from now)
			// [5] = postfix week modifier

			countStr := ""
			prefix := ""
			postfix := ""
			direction := ""

			if len(match) > 1 && match[1] != "" {
				countStr = match[1]
			}
			if len(match) > 2 && match[2] != "" {
				prefix = match[2]
			}
			if len(match) > 4 && match[4] != "" {
				direction = strings.ToLower(match[4])
			}
			if len(match) > 5 && match[5] != "" {
				postfix = match[5]
			}

			modifierWord := prefix
			if modifierWord == "" {
				modifierWord = postfix
			}
			modifierWord = strings.ToLower(modifierWord)

			modifier := ""
			switch modifierWord {
			case "last", "past":
				modifier = "last"
			case "next":
				modifier = "next"
			case "this":
				modifier = "this"
			}

			weekdayWord := strings.ToLower(strings.TrimSuffix(match[3], "s"))

			// Handle special cases that may return early
			if weekdayWord == "weekend" && (countStr != "" || direction != "") {
				refDate := context.Reference().GetDateWithAdjustedTimezone()
				return handleWeekendCounting(context, refDate, countStr, direction)
			}

			weekday := resolveWeekday(weekdayWord, modifier, context)
			if weekday < 0 {
				return nil
			}

			// Create components with weekday
			components := context.CreateParsingComponents(nil)
			refDate := context.Reference().GetDateWithAdjustedTimezone()

			// Calculate days offset to target weekday
			var modPtr *string
			if modifier != "" {
				modPtr = &modifier
			} else if context.Option().ForwardDate() {
				// When ForwardDate is enabled and no modifier is given,
				// treat it as "this" (forward)
				thisModifier := "this"
				modPtr = &thisModifier
			}
			daysOffset := helpers.GetDaysToWeekday(refDate, weekday, modPtr)
			targetDate := refDate.AddDate(0, 0, daysOffset)

			// Day/month/year are implied (uncertain) - they're calculated from the weekday
			components.ImplySimilarDate(targetDate)
			// Only weekday is certain
			components.Assign(kronos.ComponentWeekday, int(weekday))
			// Set period to week-level since we're parsing weekday
			components.SetPeriod(kronos.PeriodWeek)

			return components
		},
		nil,
	)

	return parser
}

// resolveWeekday determines the target weekday based on the word and modifier
// Returns -1 if the word is not recognized
func resolveWeekday(weekdayWord, modifier string, context *kronos.InternalParsingContext) time.Weekday {
	// Regular weekday from dictionary
	if wd, ok := data.WeekdayDictionary[weekdayWord]; ok || weekdayWord == "sunday" {
		return wd
	}

	// Weekend: "this/next weekend" means Saturday, "last weekend" means Sunday
	if weekdayWord == "weekend" {
		if modifier == "last" {
			return time.Sunday
		}
		return time.Saturday
	}

	// Weekday: any day except weekend
	if weekdayWord == "weekday" {
		refDate := context.Reference().GetDateWithAdjustedTimezone()
		refWeekday := refDate.Weekday()

		if refWeekday == time.Sunday || refWeekday == time.Saturday {
			if modifier == "last" {
				return time.Friday
			}
			return time.Monday
		}

		// On a weekday, find the next/last weekday
		wd := int(refWeekday) - 1
		if modifier == "last" {
			wd--
		} else {
			wd++
		}
		return time.Weekday((wd % 5) + 1)
	}

	return -1
}

// handleWeekendCounting calculates the date for "N weekends ago" or "N weekends from now"
func handleWeekendCounting(context *kronos.InternalParsingContext, refDate time.Time, countStr, direction string) *kronos.InternalParsingComponents {
	count := 1
	if countStr != "" {
		var err error
		count, err = strconv.Atoi(countStr)
		if err != nil || count < 1 {
			return nil
		}
	}

	var targetWeekday time.Weekday
	var daysOffset int
	refWeekday := refDate.Weekday()

	if direction == "ago" {
		targetWeekday = time.Sunday
		daysToLastSunday := int(refWeekday)
		if daysToLastSunday == 0 {
			daysToLastSunday = 7
		}
		daysOffset = -(daysToLastSunday + 7*(count-1))
	} else if direction == "from now" {
		targetWeekday = time.Saturday
		daysToNextSaturday := (int(time.Saturday) - int(refWeekday) + 7) % 7
		if daysToNextSaturday == 0 {
			daysToNextSaturday = 7
		}
		daysOffset = daysToNextSaturday + 7*(count-1)
	} else {
		return nil
	}

	targetDate := refDate.AddDate(0, 0, daysOffset)
	components := context.CreateParsingComponents(nil)
	components.ImplySimilarDate(targetDate)
	components.Assign(kronos.ComponentWeekday, int(targetWeekday))
	components.SetPeriod(kronos.PeriodWeek)

	return components
}
