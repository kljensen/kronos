package common

import (
	"testing"
	"time"

	kronos "github.com/kljensen/kronos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestISOFormatParser(t *testing.T) {
	parser := NewISOFormatParser()
	refDate := time.Date(2012, 8, 8, 0, 0, 0, 0, time.Local)

	t.Run("basic ISO date YYYY-MM-DD", func(t *testing.T) {
		text := "Let's finish this before this 2013-2-7."
		context := kronos.XNewParsingContext(text, refDate, nil)
		pattern := parser.Pattern(context)
		match := pattern.FindStringSubmatch(text)

		require.NotNil(t, match)
		result := parser.Extract(context, match)
		require.NotNil(t, result)

		// Extract may return ParsingComponents or ParsingResultWithBoundary
		var components *kronos.ParsingComponents
		switch v := result.(type) {
		case *kronos.ParsingComponents:
			components = v
		case *kronos.ParsingResultWithBoundary:
			components = v.Components
		default:
			t.Fatalf("Unexpected result type: %T", result)
		}

		assert.Equal(t, 2013, *components.Get(kronos.ComponentYear))
		assert.Equal(t, 2, *components.Get(kronos.ComponentMonth))
		assert.Equal(t, 7, *components.Get(kronos.ComponentDay))

		// Time should not be set
		assert.False(t, components.IsCertain(kronos.ComponentHour))

		// Check tag
		tags := components.Tags()
		assert.True(t, tags["parser/ISOFormatParser"])
	})

	t.Run("ISO datetime with timezone offset", func(t *testing.T) {
		text := "1994-11-05T08:15:30-05:30"
		context := kronos.XNewParsingContext(text, refDate, nil)
		pattern := parser.Pattern(context)
		match := pattern.FindStringSubmatch(text)

		require.NotNil(t, match)
		result := parser.Extract(context, match)
		require.NotNil(t, result)

		// Extract may return ParsingComponents or ParsingResultWithBoundary
		var components *kronos.ParsingComponents
		switch v := result.(type) {
		case *kronos.ParsingComponents:
			components = v
		case *kronos.ParsingResultWithBoundary:
			components = v.Components
		default:
			t.Fatalf("Unexpected result type: %T", result)
		}
		assert.Equal(t, 1994, *components.Get(kronos.ComponentYear))
		assert.Equal(t, 11, *components.Get(kronos.ComponentMonth))
		assert.Equal(t, 5, *components.Get(kronos.ComponentDay))
		assert.Equal(t, 8, *components.Get(kronos.ComponentHour))
		assert.Equal(t, 15, *components.Get(kronos.ComponentMinute))
		assert.Equal(t, 30, *components.Get(kronos.ComponentSecond))
		assert.Equal(t, -330, *components.Get(kronos.ComponentTimezoneOffset))
	})

	t.Run("ISO datetime with Z timezone", func(t *testing.T) {
		text := "1994-11-05T13:15:30Z"
		context := kronos.XNewParsingContext(text, refDate, nil)
		pattern := parser.Pattern(context)
		match := pattern.FindStringSubmatch(text)

		require.NotNil(t, match)
		result := parser.Extract(context, match)
		require.NotNil(t, result)

		// Extract may return ParsingComponents or ParsingResultWithBoundary
		var components *kronos.ParsingComponents
		switch v := result.(type) {
		case *kronos.ParsingComponents:
			components = v
		case *kronos.ParsingResultWithBoundary:
			components = v.Components
		default:
			t.Fatalf("Unexpected result type: %T", result)
		}
		assert.Equal(t, 1994, *components.Get(kronos.ComponentYear))
		assert.Equal(t, 11, *components.Get(kronos.ComponentMonth))
		assert.Equal(t, 5, *components.Get(kronos.ComponentDay))
		assert.Equal(t, 13, *components.Get(kronos.ComponentHour))
		assert.Equal(t, 15, *components.Get(kronos.ComponentMinute))
		assert.Equal(t, 30, *components.Get(kronos.ComponentSecond))
		assert.Equal(t, 0, *components.Get(kronos.ComponentTimezoneOffset))
	})

	t.Run("ISO datetime with milliseconds", func(t *testing.T) {
		text := "2016-05-07T23:45:00.487+01:00"
		context := kronos.XNewParsingContext(text, refDate, nil)
		pattern := parser.Pattern(context)
		match := pattern.FindStringSubmatch(text)

		require.NotNil(t, match)
		result := parser.Extract(context, match)
		require.NotNil(t, result)

		// Extract may return ParsingComponents or ParsingResultWithBoundary
		var components *kronos.ParsingComponents
		switch v := result.(type) {
		case *kronos.ParsingComponents:
			components = v
		case *kronos.ParsingResultWithBoundary:
			components = v.Components
		default:
			t.Fatalf("Unexpected result type: %T", result)
		}
		assert.Equal(t, 2016, *components.Get(kronos.ComponentYear))
		assert.Equal(t, 5, *components.Get(kronos.ComponentMonth))
		assert.Equal(t, 7, *components.Get(kronos.ComponentDay))
		assert.Equal(t, 23, *components.Get(kronos.ComponentHour))
		assert.Equal(t, 45, *components.Get(kronos.ComponentMinute))
		assert.Equal(t, 0, *components.Get(kronos.ComponentSecond))
		assert.Equal(t, 487, *components.Get(kronos.ComponentMillisecond))
		assert.Equal(t, 60, *components.Get(kronos.ComponentTimezoneOffset))
	})

	t.Run("ISO datetime with microseconds (4 digits)", func(t *testing.T) {
		text := "2016-05-07T23:45:00.4876+01:00"
		context := kronos.XNewParsingContext(text, refDate, nil)
		pattern := parser.Pattern(context)
		match := pattern.FindStringSubmatch(text)

		require.NotNil(t, match)
		result := parser.Extract(context, match)
		require.NotNil(t, result)

		// Extract may return ParsingComponents or ParsingResultWithBoundary
		var components *kronos.ParsingComponents
		switch v := result.(type) {
		case *kronos.ParsingComponents:
			components = v
		case *kronos.ParsingResultWithBoundary:
			components = v.Components
		default:
			t.Fatalf("Unexpected result type: %T", result)
		}
		// Should truncate to 3 digits (milliseconds)
		assert.Equal(t, 487, *components.Get(kronos.ComponentMillisecond))
	})

	t.Run("ISO datetime without seconds", func(t *testing.T) {
		text := "2016-05-07T23:45+01:00"
		context := kronos.XNewParsingContext(text, refDate, nil)
		pattern := parser.Pattern(context)
		match := pattern.FindStringSubmatch(text)

		require.NotNil(t, match)
		result := parser.Extract(context, match)
		require.NotNil(t, result)

		// Extract may return ParsingComponents or ParsingResultWithBoundary
		var components *kronos.ParsingComponents
		switch v := result.(type) {
		case *kronos.ParsingComponents:
			components = v
		case *kronos.ParsingResultWithBoundary:
			components = v.Components
		default:
			t.Fatalf("Unexpected result type: %T", result)
		}
		assert.Equal(t, 2016, *components.Get(kronos.ComponentYear))
		assert.Equal(t, 5, *components.Get(kronos.ComponentMonth))
		assert.Equal(t, 7, *components.Get(kronos.ComponentDay))
		assert.Equal(t, 23, *components.Get(kronos.ComponentHour))
		assert.Equal(t, 45, *components.Get(kronos.ComponentMinute))
		assert.False(t, components.IsCertain(kronos.ComponentSecond))
		assert.Equal(t, 60, *components.Get(kronos.ComponentTimezoneOffset))
	})

	t.Run("ISO datetime without timezone", func(t *testing.T) {
		text := "1994-11-05T13:15:30"
		context := kronos.XNewParsingContext(text, refDate, nil)
		pattern := parser.Pattern(context)
		match := pattern.FindStringSubmatch(text)

		require.NotNil(t, match)
		result := parser.Extract(context, match)
		require.NotNil(t, result)

		// Extract may return ParsingComponents or ParsingResultWithBoundary
		var components *kronos.ParsingComponents
		switch v := result.(type) {
		case *kronos.ParsingComponents:
			components = v
		case *kronos.ParsingResultWithBoundary:
			components = v.Components
		default:
			t.Fatalf("Unexpected result type: %T", result)
		}
		assert.Equal(t, 1994, *components.Get(kronos.ComponentYear))
		assert.Equal(t, 11, *components.Get(kronos.ComponentMonth))
		assert.Equal(t, 5, *components.Get(kronos.ComponentDay))
		assert.Equal(t, 13, *components.Get(kronos.ComponentHour))
		assert.Equal(t, 15, *components.Get(kronos.ComponentMinute))
		assert.Equal(t, 30, *components.Get(kronos.ComponentSecond))
		assert.False(t, components.IsCertain(kronos.ComponentTimezoneOffset))
	})

	t.Run("with leading dash", func(t *testing.T) {
		text := "- 1994-11-05T13:15:30Z"
		context := kronos.XNewParsingContext(text, refDate, nil)
		pattern := parser.Pattern(context)
		match := pattern.FindStringSubmatch(text)

		require.NotNil(t, match)
		result := parser.Extract(context, match)
		require.NotNil(t, result)

		// Extract may return ParsingComponents or ParsingResultWithBoundary
		var components *kronos.ParsingComponents
		switch v := result.(type) {
		case *kronos.ParsingComponents:
			components = v
		case *kronos.ParsingResultWithBoundary:
			components = v.Components
		default:
			t.Fatalf("Unexpected result type: %T", result)
		}
		assert.Equal(t, 1994, *components.Get(kronos.ComponentYear))
		assert.Equal(t, 11, *components.Get(kronos.ComponentMonth))
		assert.Equal(t, 5, *components.Get(kronos.ComponentDay))
	})

	t.Run("word boundary prevents false match", func(t *testing.T) {
		text := "Version1994-11-05T13:15:30Z"
		context := kronos.XNewParsingContext(text, refDate, nil)
		pattern := parser.Pattern(context)
		match := pattern.FindStringSubmatch(text)

		// Should not match because there's no word boundary before the date
		assert.Nil(t, match)
	})

	t.Run("single digit month and day", func(t *testing.T) {
		text := "2023-1-5"
		context := kronos.XNewParsingContext(text, refDate, nil)
		pattern := parser.Pattern(context)
		match := pattern.FindStringSubmatch(text)

		require.NotNil(t, match)
		result := parser.Extract(context, match)
		require.NotNil(t, result)

		// Extract may return ParsingComponents or ParsingResultWithBoundary
		var components *kronos.ParsingComponents
		switch v := result.(type) {
		case *kronos.ParsingComponents:
			components = v
		case *kronos.ParsingResultWithBoundary:
			components = v.Components
		default:
			t.Fatalf("Unexpected result type: %T", result)
		}
		assert.Equal(t, 2023, *components.Get(kronos.ComponentYear))
		assert.Equal(t, 1, *components.Get(kronos.ComponentMonth))
		assert.Equal(t, 5, *components.Get(kronos.ComponentDay))
	})
}
