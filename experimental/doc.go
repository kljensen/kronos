// Package experimental is DEPRECATED and will be removed in a future version.
//
// # Migration Guide
//
// The experimental package has been deprecated because it was over-engineered
// and exposed too many internal APIs. Most users never needed this package.
//
// ## Pre-configured Chronos (MOVED)
//
// The pre-configured Chrono instances have been moved to the en package:
//
//   - experimental.EnglishCasualChrono() → en.CasualChrono()
//   - experimental.EnglishStrictChrono() → en.StrictChrono()
//   - experimental.EnglishGBChrono() → en.GBChrono()
//
// However, most users should use the builder API instead:
//
//   - en.New() for casual English parsing
//   - en.NewStrict() for strict English parsing
//   - en.NewGB() for British English parsing
//
// ## Advanced Options (REMOVED)
//
// The experimental options (WithDayPreference, WithTimezoneAware, etc.) have been
// removed. These exposed internal settings that most users never needed. If you were
// using these options, you have two alternatives:
//
// 1. Use the main ParserBuilder API which covers most use cases:
//
//     parser := en.New().
//         WithReferenceDate(refDate).
//         PreferPast().
//         DateOrder(kronos.DateOrderDMY)
//
// 2. If you truly need low-level control, access the Settings directly via
//    ParserBuilder.Settings() and modify them before calling Parse().
//
// ## Low-level APIs (USE INTERNAL PACKAGES)
//
// The low-level Chrono/Configuration/Parser/Refiner types are still available
// in the main kronos package. If you need even lower-level access, you can
// import the internal packages directly (though this is not recommended and
// not covered by API stability guarantees).
//
// ## Type Re-exports (USE MAIN PACKAGE)
//
// All type re-exports (Component, Duration, Settings, etc.) are available in
// the main kronos package. Simply import "github.com/kljensen/kronos" instead
// of "github.com/kljensen/kronos/experimental".
//
// # Rationale
//
// The experimental package was originally created to provide an "escape hatch"
// for advanced users and to help with API migration. However, it caused more
// confusion than it solved:
//
//   - It wasn't clear what was "experimental" vs. stable
//   - It exposed too many internal implementation details
//   - It duplicated types from the main package
//   - Most users never needed any of its features
//
// By removing it, we simplify the API surface and make the library easier to
// understand and maintain.
//
// Deprecated: Use the en package for English parsing and the main kronos package
// for types and configuration.
package experimental
