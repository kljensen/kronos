//nolint:staticcheck // SA1019: Must use deprecated types during transition
package parsers

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/internal/parsing"
	"github.com/kljensen/kronos/internal/en/data"
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
		func(context *kronos.ParsingContext) *regexp.Regexp {
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
		func(context *kronos.ParsingContext, match []string) interface{} {
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

			weekdayWord := strings.ToLower(match[3])
			// Remove plural 's' if present
			weekdayWord = strings.TrimSuffix(weekdayWord, "s")

			var weekday time.Weekday

			switch {
			case data.WeekdayDictionary[weekdayWord] != 0 || weekdayWord == "sunday":
				// Regular weekday from dictionary
				weekday = data.WeekdayDictionary[weekdayWord]
			case weekdayWord == "weekend":
				// Handle weekend counting (e.g., "2 weekends ago", "1 weekend from now")
				if countStr != "" || direction != "" {
					refDate := context.Reference().GetDateWithAdjustedTimezone()
					components := handleWeekendCounting(context, refDate, countStr, direction)
					if components != nil {
						return components
					}
					return nil
				}

				// "This/next weekend" means the coming Saturday,
				// "last weekend" means last Sunday
				if modifier == "last" {
					weekday = time.Sunday
				} else {
					weekday = time.Saturday
				}
			case weekdayWord == "weekday":
				// Weekday means any day of the week except weekend
				refDate := context.Reference().GetDateWithAdjustedTimezone()
				refWeekday := time.Weekday(refDate.Weekday())

				if refWeekday == time.Sunday || refWeekday == time.Saturday {
					if modifier == "last" {
						weekday = time.Friday
					} else {
						weekday = time.Monday
					}
				} else {
					// On a weekday, find the next/last weekday
					wd := int(refWeekday) - 1
					if modifier == "last" {
						wd--
					} else {
						wd++
					}
					wd = (wd % 5) + 1
					weekday = time.Weekday(wd)
				}
			default:
				return nil
			}

			// Create components with weekday
			components := context.CreateParsingComponents(nil)
			refDate := context.Reference().GetDateWithAdjustedTimezone()

			// Calculate days offset to target weekday
			var modPtr *string
			if modifier != "" {
				modPtr = &modifier
			} else if context.Option().ForwardDate {
				// When ForwardDate is enabled and no modifier is given,
				// treat it as "this" (forward)
				thisModifier := "this"
				modPtr = &thisModifier
			}
			daysOffset := kronos.XGetDaysToWeekday(refDate, weekday, modPtr)
			targetDate := refDate.AddDate(0, 0, daysOffset)

			// Day/month/year are implied (uncertain) - they're calculated from the weekday
			kronos.XImplySimilarDate(components, targetDate)
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

// handleWeekendCounting calculates the date for "N weekends ago" or "N weekends from now"
func handleWeekendCounting(context *kronos.ParsingContext, refDate time.Time, countStr, direction string) *kronos.ParsingComponents {
	// Parse count, default to 1 if not specified
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

	switch direction {
	case "ago":
		// Looking backward: resolve to Sunday
		targetWeekday = time.Sunday

		// Calculate days back to the Nth previous Sunday
		// First, find days to last Sunday
		daysToLastSunday := int(refWeekday)
		if daysToLastSunday == 0 {
			// Already on Sunday, go back 7 days to previous Sunday
			daysToLastSunday = 7
		}

		// Then go back (count-1) more weeks
		daysOffset = -(daysToLastSunday + 7*(count-1))
	case "from now":
		// Looking forward: resolve to Saturday
		targetWeekday = time.Saturday

		// Calculate days forward to the Nth next Saturday
		daysToNextSaturday := (int(time.Saturday) - int(refWeekday) + 7) % 7
		if daysToNextSaturday == 0 {
			// Already on Saturday, go forward 7 days to next Saturday
			daysToNextSaturday = 7
		}

		// Then go forward (count-1) more weeks
		daysOffset = daysToNextSaturday + 7*(count-1)
	default:
		return nil
	}

	targetDate := refDate.AddDate(0, 0, daysOffset)

	// Create components
	components := context.CreateParsingComponents(nil)
	kronos.XImplySimilarDate(components, targetDate)
	components.Assign(kronos.ComponentWeekday, int(targetWeekday))
	components.SetPeriod(kronos.PeriodWeek)

	return components
}
