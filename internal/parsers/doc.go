// Package parsers contains all parser implementations for the Kronos library.
//
// Parser implementations will be moved here from the root and en packages.
// This includes parsers for:
//   - ISO 8601 formats
//   - Slash date formats (MM/DD/YYYY, DD/MM/YYYY)
//   - Month name formats (Jan 1, January 1st, etc.)
//   - Relative dates (yesterday, tomorrow, last week)
//   - Time expressions (3pm, 15:30, noon)
//   - Casual references (now, today)
//
// These implementations are hidden from users, who only interact with
// the builder-based Parser API.
package parsers
