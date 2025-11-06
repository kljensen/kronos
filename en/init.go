package en

import (
	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/common"
	"github.com/kljensen/kronos/internal/en/parsers"
)

// init registers all English parsers with the global registry.
// This allows parsers to be discovered and used by the settings system.
func init() {
	registerCommonParsers()
	registerEnglishParsers()
}

// registerCommonParsers registers common parsers used across languages.
func registerCommonParsers() {
	kronos.Register("iso8601", kronos.ParserInfo{
		Description: "ISO 8601 date/time formats (e.g., 2020-03-15, 2020-03-15T14:30:00Z)",
		Priority:    100, // High priority - most specific
		Tags:        []string{"iso", "formal"},
	}, func() kronos.Parser {
		return common.NewISOFormatParser()
	})

	kronos.Register("slash_date", kronos.ParserInfo{
		Description: "Slash-separated dates (e.g., 3/15/2020, 15/3/2020)",
		Priority:    80,
		Tags:        []string{"formal", "numeric"},
	}, func() kronos.Parser {
		return common.NewSlashDateFormatParser(false) // US format by default
	})

	kronos.Register("slash_date_little_endian", kronos.ParserInfo{
		Description: "Little-endian slash dates (e.g., 15/3/2020)",
		Priority:    80,
		Tags:        []string{"formal", "numeric", "uk"},
	}, func() kronos.Parser {
		return common.NewSlashDateFormatParser(true)
	})
}

// registerEnglishParsers registers all English-specific parsers.
func registerEnglishParsers() {
	kronos.Register("en_year_month_day", kronos.ParserInfo{
		Description: "Year-month-day patterns (e.g., 2020-03-15)",
		Priority:    90,
		Tags:        []string{"formal", "english"},
	}, func() kronos.Parser {
		return parsers.NewENYearMonthDayParser(false)
	})

	kronos.Register("en_time_unit_within", kronos.ParserInfo{
		Description: "Time units within a period (e.g., '5 days within this month')",
		Priority:    70,
		Tags:        []string{"relative", "english"},
	}, func() kronos.Parser {
		return parsers.NewENTimeUnitWithinFormatParser(false)
	})

	kronos.Register("en_month_name_little_endian", kronos.ParserInfo{
		Description: "Little-endian month names (e.g., '15 March 2020')",
		Priority:    75,
		Tags:        []string{"formal", "english", "uk"},
	}, func() kronos.Parser {
		return parsers.NewENMonthNameLittleEndianParser()
	})

	kronos.Register("en_month_name_middle_endian", kronos.ParserInfo{
		Description: "Middle-endian month names (e.g., 'March 15, 2020')",
		Priority:    75,
		Tags:        []string{"formal", "english", "us"},
	}, func() kronos.Parser {
		return parsers.NewENMonthNameMiddleEndianParser(false)
	})

	kronos.Register("en_weekday", kronos.ParserInfo{
		Description: "Weekday names (e.g., 'Monday', 'next Friday')",
		Priority:    65,
		Tags:        []string{"casual", "english"},
	}, func() kronos.Parser {
		return parsers.NewENWeekdayParser()
	})

	kronos.Register("en_slash_month", kronos.ParserInfo{
		Description: "Slash month formats (e.g., '3/15')",
		Priority:    60,
		Tags:        []string{"informal", "english"},
	}, func() kronos.Parser {
		return parsers.NewENSlashMonthFormatParser()
	})

	kronos.Register("en_time_expression", kronos.ParserInfo{
		Description: "Time expressions (e.g., '3:30 PM', '14:30', 'noon')",
		Priority:    70,
		Tags:        []string{"time", "english"},
	}, func() kronos.Parser {
		return parsers.NewENTimeExpressionParser(false)
	})

	kronos.Register("en_time_unit_ago", kronos.ParserInfo{
		Description: "Time units ago (e.g., '2 days ago', '3 hours ago')",
		Priority:    65,
		Tags:        []string{"relative", "english"},
	}, func() kronos.Parser {
		return parsers.NewENTimeUnitAgoFormatParser(false)
	})

	kronos.Register("en_time_unit_later", kronos.ParserInfo{
		Description: "Time units later (e.g., 'in 2 days', '3 hours from now')",
		Priority:    65,
		Tags:        []string{"relative", "english"},
	}, func() kronos.Parser {
		return parsers.NewENTimeUnitLaterFormatParser(false)
	})

	kronos.Register("en_casual_date", kronos.ParserInfo{
		Description: "Casual date expressions (e.g., 'today', 'tomorrow', 'yesterday')",
		Priority:    60,
		Tags:        []string{"casual", "english"},
	}, func() kronos.Parser {
		return parsers.NewENCasualDateParser()
	})

	kronos.Register("en_casual_time", kronos.ParserInfo{
		Description: "Casual time expressions (e.g., 'now', 'tonight', 'this morning')",
		Priority:    60,
		Tags:        []string{"casual", "english", "time"},
	}, func() kronos.Parser {
		return parsers.NewENCasualTimeParser()
	})

	kronos.Register("en_month_name", kronos.ParserInfo{
		Description: "Month names alone (e.g., 'March', 'December')",
		Priority:    50,
		Tags:        []string{"casual", "english"},
	}, func() kronos.Parser {
		return parsers.NewENMonthNameParser()
	})

	kronos.Register("en_relative_date", kronos.ParserInfo{
		Description: "Relative date expressions (e.g., 'last week', 'next month')",
		Priority:    55,
		Tags:        []string{"relative", "english"},
	}, func() kronos.Parser {
		return parsers.NewENRelativeDateFormatParser()
	})

	kronos.Register("en_time_unit_casual_relative", kronos.ParserInfo{
		Description: "Casual relative time units (e.g., 'this week', 'last year')",
		Priority:    55,
		Tags:        []string{"relative", "casual", "english"},
	}, func() kronos.Parser {
		return parsers.NewENTimeUnitCasualRelativeFormatParser(true)
	})

	kronos.Register("en_compact", kronos.ParserInfo{
		Description: "Compact numeric formats (e.g., '20200315', '031520')",
		Priority:    30, // Low priority - catch-all
		Tags:        []string{"numeric", "informal", "english"},
	}, func() kronos.Parser {
		return parsers.NewENCompactFormatParser()
	})
}

// SetDefaultParserOrder sets the default parser order in the global registry.
// This order is used when no custom order is specified in settings.
func SetDefaultParserOrder() {
	kronos.GlobalRegistry.SetDefaultOrder([]string{
		// Most specific first
		"iso8601",
		"en_year_month_day",
		"slash_date",
		"slash_date_little_endian",
		"en_month_name_little_endian",
		"en_month_name_middle_endian",
		"en_time_expression",
		"en_time_unit_within",
		"en_time_unit_ago",
		"en_time_unit_later",
		"en_weekday",
		"en_casual_date",
		"en_casual_time",
		"en_slash_month",
		"en_relative_date",
		"en_time_unit_casual_relative",
		"en_month_name",
		// Least specific last
		"en_compact",
	})
}

func init() {
	SetDefaultParserOrder()
}
