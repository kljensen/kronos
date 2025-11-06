package common

import (
	"github.com/kljensen/kronos/internal/common/parsers"
)

// Deprecated: ISOFormatParser has been moved to internal/common/parsers.
// Use NewEnglishConfiguration() or NewConfiguration() from the parent package instead.
type ISOFormatParser = parsers.ISOFormatParser

// Deprecated: NewISOFormatParser has been moved to internal/common/parsers.
// This is provided for backward compatibility but should not be used directly.
var NewISOFormatParser = parsers.NewISOFormatParser
