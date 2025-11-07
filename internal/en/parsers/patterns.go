package parsers

// Common regex patterns used across multiple parsers
const (
	// approximationPattern matches optional approximation words at the beginning
	// Includes: ~, about, around, roughly, approximately, approx, circa
	approximationPattern = `(?:~\s*|(?:about|around|roughly|approximately|approx|circa)\s+)?`

	// wordNumbers matches common number words
	wordNumbers = `half|dozen|several|couple|few|ninety|eighty|seventy|sixty|fifty|forty|thirty|twenty|nineteen|eighteen|seventeen|sixteen|fifteen|fourteen|thirteen|twelve|eleven|ten|nine|eight|seven|six|five|four|three|two|one|a|an|the`

	// unitSeparator matches separators between time units (space, comma, "and")
	unitSeparator = `(?:\s*,\s*|\s+and\s+|\s+)`
)
