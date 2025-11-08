// Package parsers provides shared utilities and parsers for date/time parsing.
//
//nolint:staticcheck // SA1019: Must use deprecated types during transition
package parsers

import (
	"strings"

	"github.com/kljensen/kronos/internal/helpers"
)

// parseMeridiem parses AM/PM indicator and adjusts hour accordingly.
// Returns the adjusted hour, meridiem value, and success status.
// For AM: hour 12 becomes 0 (midnight). For PM: hours 1-11 add 12.
func parseMeridiem(ampmStr string, hour int) (int, *helpers.Meridiem, bool) {
	if hour > maxHour12Format {
		return 0, nil, false
	}

	ampm := strings.ToLower(string(ampmStr[0]))
	if ampm == "a" {
		m := helpers.MeridiemAM
		if hour == maxHour12Format {
			hour = 0
		}
		return hour, &m, true
	}

	if ampm == "p" {
		m := helpers.MeridiemPM
		if hour != maxHour12Format {
			hour += maxHour12Format
		}
		return hour, &m, true
	}

	return hour, nil, true
}
