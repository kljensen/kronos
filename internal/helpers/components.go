package helpers

import (
	"time"

	"github.com/kljensen/kronos"
)

// AssignSimilarDate assigns (force updates) the parsing components to the same day as the target.
// This sets year, month, and day as certain (known) values.
func AssignSimilarDate(components *kronos.InternalParsingComponents, date time.Time) {
	components.AssignSimilarDate(date)
}

// AssignSimilarTime assigns (force updates) the parsing components to the same time as the target.
// This sets hour, minute, second, millisecond, microsecond, nanosecond, and meridiem as certain (known) values.
func AssignSimilarTime(components *kronos.InternalParsingComponents, date time.Time) {
	components.AssignSimilarTime(date)
}

// ImplySimilarDate implies (weakly updates) the parsing components to the same day as the target.
// This sets year, month, and day as implied values (only if not already certain).
func ImplySimilarDate(components *kronos.InternalParsingComponents, date time.Time) {
	components.ImplySimilarDate(date)
}

// ImplySimilarTime implies (weakly updates) the parsing components to the same time as the target.
// This sets hour, minute, second, millisecond, microsecond, nanosecond, and meridiem as implied values (only if not already certain).
func ImplySimilarTime(components *kronos.InternalParsingComponents, date time.Time) {
	components.ImplySimilarTime(date)
}

// MergeDateTimeComponent merges date and time components.
func MergeDateTimeComponent(dateComp, timeComp *kronos.InternalParsingComponents) *kronos.InternalParsingComponents {
	return kronos.InternalMergeDateTimeComponent(dateComp, timeComp)
}
