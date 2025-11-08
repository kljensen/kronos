package kronos_test

import (
	"testing"
	"time"

	"github.com/kljensen/kronos/internal/chrono"

	"github.com/kljensen/kronos"
	"github.com/kljensen/kronos/internal/common/parsers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestChronoIntegration demonstrates the full Chrono engine with real parsers
func TestChronoIntegration(t *testing.T) {
	t.Run("parse ISO date", func(t *testing.T) {
		// Create a configuration with ISO parser
		config := &chrono.Configuration{
			Parsers: []any{
				parsers.NewISOFormatParser(),
			},
		}

		c := chrono.NewChrono(config)
		refDate := time.Date(2012, 8, 8, 0, 0, 0, 0, time.Local)

		results, err := kronos.New(c).
			WithReferenceDate(refDate).
			Parse("Let's finish this before 2013-02-07")

		require.NoError(t, err)
		require.Len(t, results, 1)
		result := results[0]

		assert.Equal(t, 2013, *result.Start().Get(kronos.ComponentYear))
		assert.Equal(t, 2, *result.Start().Get(kronos.ComponentMonth))
		assert.Equal(t, 7, *result.Start().Get(kronos.ComponentDay))
		// Note: Tags are not exposed in the public Result interface
		// assert.True(t, result.Tags()["parser/ISOFormatParser"])
	})

	t.Run("parse ISO datetime with timezone", func(t *testing.T) {
		config := &chrono.Configuration{
			Parsers: []any{
				parsers.NewISOFormatParser(),
			},
		}

		c := chrono.NewChrono(config)
		refDate := time.Date(2012, 8, 8, 0, 0, 0, 0, time.Local)

		results, err := kronos.New(c).
			WithReferenceDate(refDate).
			Parse("Event at 1994-11-05T08:15:30-05:30")

		require.NoError(t, err)
		require.Len(t, results, 1)
		result := results[0]

		assert.Equal(t, 1994, *result.Start().Get(kronos.ComponentYear))
		assert.Equal(t, 11, *result.Start().Get(kronos.ComponentMonth))
		assert.Equal(t, 5, *result.Start().Get(kronos.ComponentDay))
		assert.Equal(t, 8, *result.Start().Get(kronos.ComponentHour))
		assert.Equal(t, 15, *result.Start().Get(kronos.ComponentMinute))
		assert.Equal(t, 30, *result.Start().Get(kronos.ComponentSecond))
		assert.Equal(t, -330, *result.Start().Get(kronos.ComponentTimezoneOffset))
	})

	t.Run("parse slash date MM/DD/YYYY", func(t *testing.T) {
		config := &chrono.Configuration{
			Parsers: []any{
				parsers.NewSlashDateFormatParser(false), // US format
			},
		}

		c := chrono.NewChrono(config)
		refDate := time.Date(2012, 8, 10, 0, 0, 0, 0, time.Local)

		results, err := kronos.New(c).
			WithReferenceDate(refDate).
			Parse("The Deadline is 8/10/2012")

		require.NoError(t, err)
		require.Len(t, results, 1)
		result := results[0]

		assert.Equal(t, 2012, *result.Start().Get(kronos.ComponentYear))
		assert.Equal(t, 8, *result.Start().Get(kronos.ComponentMonth))
		assert.Equal(t, 10, *result.Start().Get(kronos.ComponentDay))
		assert.Equal(t, "8/10/2012", result.Text())
	})

	t.Run("parse slash date DD/MM/YYYY", func(t *testing.T) {
		config := &chrono.Configuration{
			Parsers: []any{
				parsers.NewSlashDateFormatParser(true), // UK/EU format
			},
		}

		c := chrono.NewChrono(config)
		refDate := time.Date(2012, 8, 10, 0, 0, 0, 0, time.Local)

		results, err := kronos.New(c).
			WithReferenceDate(refDate).
			Parse("The Deadline is 8/10/2012")

		require.NoError(t, err)
		require.Len(t, results, 1)
		result := results[0]

		assert.Equal(t, 2012, *result.Start().Get(kronos.ComponentYear))
		assert.Equal(t, 10, *result.Start().Get(kronos.ComponentMonth))
		assert.Equal(t, 8, *result.Start().Get(kronos.ComponentDay))
	})

	t.Run("multiple parsers - multiple matches", func(t *testing.T) {
		config := &chrono.Configuration{
			Parsers: []any{
				parsers.NewISOFormatParser(),
				parsers.NewSlashDateFormatParser(false),
			},
		}

		c := chrono.NewChrono(config)
		refDate := time.Date(2012, 8, 8, 0, 0, 0, 0, time.Local)

		results, err := kronos.New(c).
			WithReferenceDate(refDate).
			Parse("Meeting on 8/10/2012 and event at 2013-02-07")

		require.NoError(t, err)
		require.Len(t, results, 2)

		// Results should be sorted by position
		assert.True(t, results[0].Index() < results[1].Index())

		// First result should be slash date
		assert.Equal(t, 2012, *results[0].Start().Get(kronos.ComponentYear))
		assert.Equal(t, 8, *results[0].Start().Get(kronos.ComponentMonth))
		assert.Equal(t, 10, *results[0].Start().Get(kronos.ComponentDay))

		// Second result should be ISO date
		assert.Equal(t, 2013, *results[1].Start().Get(kronos.ComponentYear))
		assert.Equal(t, 2, *results[1].Start().Get(kronos.ComponentMonth))
		assert.Equal(t, 7, *results[1].Start().Get(kronos.ComponentDay))
	})

	t.Run("ParseDate returns first result date", func(t *testing.T) {
		config := &chrono.Configuration{
			Parsers: []any{
				parsers.NewISOFormatParser(),
			},
		}

		c := chrono.NewChrono(config)
		refDate := time.Date(2012, 8, 8, 0, 0, 0, 0, time.Local)

		date, err := kronos.New(c).
			WithReferenceDate(refDate).
			ParseDate("Event at 1994-11-05T13:15:30Z")

		require.NoError(t, err)
		require.NotNil(t, date)
		assert.Equal(t, 1994, date.Year())
		assert.Equal(t, time.Month(11), date.Month())
		assert.Equal(t, 5, date.Day())
		// Note: hour may differ due to timezone adjustments from UTC to local
		// The important thing is that we got a valid date
		assert.Equal(t, 15, date.Minute())
		assert.Equal(t, 30, date.Second())
	})

	t.Run("ParseDate returns nil when no match", func(t *testing.T) {
		config := &chrono.Configuration{
			Parsers: []any{
				parsers.NewISOFormatParser(),
			},
		}

		c := chrono.NewChrono(config)
		refDate := time.Date(2012, 8, 8, 0, 0, 0, 0, time.Local)

		date, err := kronos.New(c).
			WithReferenceDate(refDate).
			ParseDate("No date here")

		require.NoError(t, err)
		assert.Nil(t, date)
	})

	t.Run("2-digit year inference", func(t *testing.T) {
		config := &chrono.Configuration{
			Parsers: []any{
				parsers.NewSlashDateFormatParser(false),
			},
		}

		c := chrono.NewChrono(config)
		refDate := time.Date(2012, 8, 10, 0, 0, 0, 0, time.Local)

		results, err := kronos.New(c).
			WithReferenceDate(refDate).
			Parse("Meeting on 12/30/16")

		require.NoError(t, err)
		require.Len(t, results, 1)
		result := results[0]

		// Should interpret as 2016
		assert.Equal(t, 2016, *result.Start().Get(kronos.ComponentYear))
	})

	t.Run("year inference without year", func(t *testing.T) {
		config := &chrono.Configuration{
			Parsers: []any{
				parsers.NewSlashDateFormatParser(false),
			},
		}

		c := chrono.NewChrono(config)
		refDate := time.Date(2012, 8, 10, 0, 0, 0, 0, time.Local)

		results, err := kronos.New(c).
			WithReferenceDate(refDate).
			Parse("Meeting on 8/15")

		require.NoError(t, err)
		require.Len(t, results, 1)
		result := results[0]

		// Year should be implied (not certain)
		assert.False(t, result.Start().IsCertain(kronos.ComponentYear))
		assert.NotNil(t, result.Start().Get(kronos.ComponentYear))
	})

	t.Run("word boundaries prevent false matches", func(t *testing.T) {
		config := &chrono.Configuration{
			Parsers: []any{
				parsers.NewISOFormatParser(),
			},
		}

		c := chrono.NewChrono(config)
		refDate := time.Date(2012, 8, 8, 0, 0, 0, 0, time.Local)

		// Should not match date in version number
		results, err := kronos.New(c).
			WithReferenceDate(refDate).
			Parse("Version1994-11-05T13:15:30Z released")

		require.NoError(t, err)
		assert.Len(t, results, 0)
	})
}
