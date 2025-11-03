// Package refiners contains all refiner implementations for the Kronos library.
//
// Refiner implementations will be moved here from the common/refiners package.
// This includes refiners for:
//   - Merging adjacent results (date + time, date ranges)
//   - Date preference handling (past/future disambiguation)
//   - Timezone handling
//   - Forward date filtering
//   - Overlap removal
//   - Unlikely format filtering
//
// These implementations are hidden from users, who configure refinement
// behavior through the builder API.
package refiners
