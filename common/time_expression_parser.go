package common

import (
	"github.com/kljensen/kronos/internal/common/parsers"
)

// Deprecated: AbstractTimeExpressionParser has been moved to internal/common/parsers.
// Use NewEnglishConfiguration() or NewConfiguration() from the parent package instead.
type AbstractTimeExpressionParser = parsers.AbstractTimeExpressionParser

// Deprecated: NewAbstractTimeExpressionParser has been moved to internal/common/parsers.
// This is provided for backward compatibility but should not be used directly.
var NewAbstractTimeExpressionParser = parsers.NewAbstractTimeExpressionParser
