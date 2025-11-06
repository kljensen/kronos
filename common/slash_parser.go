package common

import (
	"github.com/kljensen/kronos/internal/common/parsers"
)

// Deprecated: SlashDateFormatParser has been moved to internal/common/parsers.
// Use NewEnglishConfiguration() or NewConfiguration() from the parent package instead.
type SlashDateFormatParser = parsers.SlashDateFormatParser

// Deprecated: NewSlashDateFormatParser has been moved to internal/common/parsers.
// This is provided for backward compatibility but should not be used directly.
var NewSlashDateFormatParser = parsers.NewSlashDateFormatParser
