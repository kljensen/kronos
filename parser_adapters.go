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

// Ensure adapters implement parser interfaces
var _ parser.Context = (*contextAdapter)(nil)
var _ parser.Reference = (*referenceWithTimezone)(nil)
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
	// Internal method returns *parsingComponents, which is what we want
	return ca.ctx.CreateParsingComponents(components)
}

func (ca *contextAdapter) CreateParsingResult(index int, textOrEndIndex any, args ...any) parser.Result {
	result := ca.ctx.CreateParsingResult(index, textOrEndIndex, args...)
	return &parserResultAdapter{result: result}
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

// parserResultAdapter wraps *parsingResult to implement parser.Result
// This is used when public parsers/refiners need to work with results.
type parserResultAdapter struct {
	result *parsingResult
}

// Ensure parserResultAdapter implements parser.Result
var _ parser.Result = (*parserResultAdapter)(nil)

func (ra *parserResultAdapter) Index() int                { return ra.result.Index() }
func (ra *parserResultAdapter) Text() string              { return ra.result.Text() }
func (ra *parserResultAdapter) Start() any                { return ra.result.Start() }
func (ra *parserResultAdapter) End() any                  { return ra.result.End() }
func (ra *parserResultAdapter) Date() time.Time           { return ra.result.Date() }
func (ra *parserResultAdapter) Clone() parser.Result      { return &parserResultAdapter{result: ra.result.Clone()} }
func (ra *parserResultAdapter) SetIndex(index int)        { ra.result.SetIndex(index) }
func (ra *parserResultAdapter) SetStart(start any)        { ra.result.SetStart(start) }
func (ra *parserResultAdapter) AddTag(tag string) parser.Result { ra.result.AddTag(tag); return ra }
func (ra *parserResultAdapter) Tags() map[string]bool     { return ra.result.Tags() }
func (ra *parserResultAdapter) RefDate() time.Time        { return ra.result.RefDate() }
