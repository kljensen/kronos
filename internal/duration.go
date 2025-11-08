package internal

// This file provides wrappers for duration functions to avoid import cycles.
// The actual implementations are in the kronos package (duration.go).
// These wrappers allow internal packages to use duration functions without
// importing the main kronos package (which would create an import cycle).

import (
	"time"

	"github.com/kljensen/kronos"
)

// AddDuration wraps kronos.InternalAddDuration for use by internal packages.
func AddDuration(ref time.Time, duration kronos.Duration) (time.Time, error) {
	return kronos.InternalAddDuration(ref, duration)
}

// ReverseDuration wraps kronos.InternalReverseDuration for use by internal packages.
func ReverseDuration(duration kronos.Duration) kronos.Duration {
	return kronos.InternalReverseDuration(duration)
}
