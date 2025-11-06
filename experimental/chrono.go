//nolint:staticcheck // SA1019: Must use deprecated types during transition
package experimental

import (
	"github.com/kljensen/kronos"
)

// Chrono is the main parsing engine that coordinates multiple parsers and refiners.
// Re-exported from the main package for the experimental API.
type Chrono = kronos.Chrono

// NewChrono creates a new Chrono instance with the given configuration.
// Re-exported from the main package for the experimental API.
var NewChrono = kronos.NewChrono

// Configuration holds the parsers and refiners for Chrono.
// Re-exported from the main package for the experimental API.
type Configuration = kronos.Configuration

// Parser is an abstraction for Chrono parsers.
// Re-exported from the main package for the experimental API.
type Parser = kronos.Parser

// Refiner is an abstraction for Chrono refiners.
// Re-exported from the main package for the experimental API.
type Refiner = kronos.Refiner

// Pipeline manages the parsing process with configurable parsers and settings.
// Re-exported from the main package for the experimental API.
type Pipeline = kronos.Pipeline

// NewPipeline creates a new parsing pipeline with the given configuration and settings.
// Re-exported from the main package for the experimental API.
var NewPipeline = kronos.NewPipeline

// NewPipelineWithSettings creates a pipeline using settings to determine parsers.
// Re-exported from the main package for the experimental API.
var NewPipelineWithSettings = kronos.NewPipelineWithSettings

// ParseWithSettings is a convenience function that creates a pipeline and executes it.
// Re-exported from the main package for the experimental API.
var ParseWithSettings = kronos.ParseWithSettings
