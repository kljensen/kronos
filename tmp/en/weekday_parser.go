package en

import (
	"regexp"
	"strings"

	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/common"
)

// ENWeekdayParser parses standalone weekdays with optional modifiers:
// "Monday", "on Friday", "this Tuesday", "next Wednesday", "last Thursday"
// Also handles special cases: "weekend", "weekday"
type ENWeekdayParser struct {
	*common.AbstractParserWithWordBoundary
}

// NewENWeekdayParser creates a new parser for weekday expressions
func NewENWeekdayParser() *ENWeekdayParser {
	parser := &ENWeekdayParser{}

	parser.AbstractParserWithWordBoundary = common.NewAbstractParserWithWordBoundary(
		func(context *kronos.ParsingContext) *regexp.Regexp {
			// Unicode parentheses (（ and ）) are literal characters, not escaped
			pattern := `(?:(?:,|\(|（)\s*)?` +
				`(?:on\s*?)?` +
				`(?:(this|last|past|next)\s*)?` +
				`(` + WeekdayPattern + `|weekend|weekday)` +
				`(?:\s*(?:,|\)|）))?` +
				`(?:\s*(this|last|past|next)\s*week)?` +
				`(?:\s|$|\b)`

			return regexp.MustCompile("(?i)" + pattern)
		},
		func(context *kronos.ParsingContext, match []string) interface{} {
			if len(match) < 3 {
				return nil
			}

			prefix := ""
			postfix := ""
			if len(match) > 1 && match[1] != "" {
				prefix = match[1]
			}
			if len(match) > 3 && match[3] != "" {
				postfix = match[3]
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

			weekdayWord := strings.ToLower(match[2])
			var weekday kronos.Weekday

			if wd, ok := WeekdayDictionary[weekdayWord]; ok {
				weekday = wd
			} else if weekdayWord == "weekend" {
				// "This/next weekend" means the coming Saturday,
				// "last weekend" means last Sunday
				if modifier == "last" {
					weekday = kronos.WeekdaySunday
				} else {
					weekday = kronos.WeekdaySaturday
				}
			} else if weekdayWord == "weekday" {
				// Weekday means any day of the week except weekend
				refDate := context.Reference().GetDateWithAdjustedTimezone()
				refWeekday := kronos.Weekday(refDate.Weekday())

				if refWeekday == kronos.WeekdaySunday || refWeekday == kronos.WeekdaySaturday {
					if modifier == "last" {
						weekday = kronos.WeekdayFriday
					} else {
						weekday = kronos.WeekdayMonday
					}
				} else {
					// On a weekday, find the next/last weekday
					wd := int(refWeekday) - 1
					if modifier == "last" {
						wd = wd - 1
					} else {
						wd = wd + 1
					}
					wd = (wd % 5) + 1
					weekday = kronos.Weekday(wd)
				}
			} else {
				return nil
			}

			// Create components with weekday
			components := context.CreateParsingComponents(nil)
			refDate := context.Reference().GetDateWithAdjustedTimezone()

			// Calculate days offset to target weekday
			var modPtr *string
			if modifier != "" {
				modPtr = &modifier
			}
			daysOffset := kronos.GetDaysToWeekday(refDate, weekday, modPtr)
			targetDate := refDate.AddDate(0, 0, daysOffset)

			components.Assign(kronos.ComponentDay, targetDate.Day())
			components.Assign(kronos.ComponentMonth, int(targetDate.Month()))
			components.Assign(kronos.ComponentYear, targetDate.Year())
			components.Assign(kronos.ComponentWeekday, int(weekday))

			return components
		},
		nil,
	)

	return parser
}
