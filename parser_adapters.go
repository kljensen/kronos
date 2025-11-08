package kronos

import (
	"time"

	"github.com/kljensen/kronos/parser"
)

// ============================================================================
// Adapters to implement parser package interfaces
// ============================================================================

// optionAdapter wraps parsingOption to implement parser.Option
type optionAdapter struct {
	opt parsingOption
}

func (o *optionAdapter) ForwardDate() bool       { return o.opt.ForwardDate }
func (o *optionAdapter) Preference() int         { return int(o.opt.Preference) }
func (o *optionAdapter) DateOrder() int          { return int(o.opt.DateOrder) }
func (o *optionAdapter) Timezones() map[string]any { return o.opt.Timezones }
func (o *optionAdapter) Debug() func(string) {
	return o.opt.Debug
}

// contextAdapter wraps parsingContext to implement parser.Context
type contextAdapter struct {
	ctx *parsingContext
}

// Ensure contextAdapter implements parser.Context
var _ parser.Context = (*contextAdapter)(nil)
var _ parser.Reference = (*referenceWithTimezone)(nil)
var _ parser.Result = (*parsingResult)(nil)
var _ parser.Option = (*optionAdapter)(nil)

// Implement parser.Context interface
func (ca *contextAdapter) Text() string {
	return ca.ctx.Text()
}

func (ca *contextAdapter) Reference() parser.Reference {
	return ca.ctx.reference
}

func (ca *contextAdapter) RefDate() time.Time {
	return ca.ctx.RefDate()
}

func (ca *contextAdapter) Option() parser.Option {
	return &optionAdapter{opt: ca.ctx.option}
}

func (ca *contextAdapter) Settings() parser.Settings {
	if ca.ctx.settings == nil {
		return nil
	}
	return &settingsAdapter{settings: ca.ctx.settings}
}

func (ca *contextAdapter) CreateParsingComponents(components any) any {
	return ca.ctx.CreateParsingComponents(components)
}

func (ca *contextAdapter) CreateParsingResult(index int, textOrEndIndex any, args ...any) parser.Result {
	return ca.ctx.CreateParsingResult(index, textOrEndIndex, args...)
}

// settingsAdapter wraps Settings to implement parser.Settings
type settingsAdapter struct {
	settings *Settings
}

func (s *settingsAdapter) ForwardDate() bool         { opt := s.settings.toparsingOption(nil); return opt.ForwardDate }
func (s *settingsAdapter) Preference() int           { opt := s.settings.toparsingOption(nil); return int(opt.Preference) }
func (s *settingsAdapter) DateOrder() int            { opt := s.settings.toparsingOption(nil); return int(opt.DateOrder) }
func (s *settingsAdapter) Timezones() map[string]any { opt := s.settings.toparsingOption(nil); return opt.Timezones }
func (s *settingsAdapter) Debug() func(string) {
	opt := s.settings.toparsingOption(nil)
	return opt.Debug
}

// CreateParsingComponents is already defined in parsingContext

// CreateParsingResult is already defined in parsingContext

// parsingResult.SetStart signature has been updated to accept any,
// so it now satisfies the parser.Result interface directly.
