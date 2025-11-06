package common

import (
	"github.com/kljensen/kronos/internal/parsing"
)

// Deprecated: AbstractParserWithWordBoundary has been moved to internal/parsing.
// Use the builder functions in the parent package to create parsers instead.
type AbstractParserWithWordBoundary = parsing.AbstractParserWithWordBoundary

// Deprecated: NewAbstractParserWithWordBoundary has been moved to internal/parsing.
// This is provided for backward compatibility but should not be used directly.
var NewAbstractParserWithWordBoundary = parsing.NewAbstractParserWithWordBoundary
