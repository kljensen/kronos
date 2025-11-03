package kronos

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddDuration(t *testing.T) {
	tests := []struct {
		name     string
		ref      time.Time
		duration Duration
		expected time.Time
	}{
		{
			name: "add years",
			ref:  time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
			duration: Duration{
				TimeunitYear: 2,
			},
			expected: time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			name: "add fractional years",
			ref:  time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
			duration: Duration{
				TimeunitYear: 1.5,
			},
			expected: time.Date(2021, 7, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			name: "add months",
			ref:  time.Date(2020, 1, 15, 12, 0, 0, 0, time.UTC),
			duration: Duration{
				TimeunitMonth: 3,
			},
			expected: time.Date(2020, 4, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name: "add months with overflow",
			ref:  time.Date(2020, 11, 15, 12, 0, 0, 0, time.UTC),
			duration: Duration{
				TimeunitMonth: 3,
			},
			expected: time.Date(2021, 2, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name: "add fractional months",
			ref:  time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
			duration: Duration{
				TimeunitMonth: 1.5,
			},
			// 1.5 months = 1 month + 0.5*4 weeks = 1 month + 2 weeks
			expected: time.Date(2020, 2, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name: "add quarters",
			ref:  time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
			duration: Duration{
				TimeunitQuarter: 2,
			},
			expected: time.Date(2020, 7, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			name: "add weeks",
			ref:  time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
			duration: Duration{
				TimeunitWeek: 2,
			},
			expected: time.Date(2020, 1, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name: "add days",
			ref:  time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
			duration: Duration{
				TimeunitDay: 10,
			},
			expected: time.Date(2020, 1, 11, 12, 0, 0, 0, time.UTC),
		},
		{
			name: "add hours",
			ref:  time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
			duration: Duration{
				TimeunitHour: 5,
			},
			expected: time.Date(2020, 1, 1, 17, 0, 0, 0, time.UTC),
		},
		{
			name: "add minutes",
			ref:  time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
			duration: Duration{
				TimeunitMinute: 30,
			},
			expected: time.Date(2020, 1, 1, 12, 30, 0, 0, time.UTC),
		},
		{
			name: "add seconds",
			ref:  time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
			duration: Duration{
				TimeunitSecond: 45,
			},
			expected: time.Date(2020, 1, 1, 12, 0, 45, 0, time.UTC),
		},
		{
			name: "add milliseconds",
			ref:  time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
			duration: Duration{
				TimeunitMillisecond: 500,
			},
			expected: time.Date(2020, 1, 1, 12, 0, 0, 500000000, time.UTC),
		},
		{
			name: "add mixed duration",
			ref:  time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
			duration: Duration{
				TimeunitYear:   1,
				TimeunitMonth:  2,
				TimeunitDay:    3,
				TimeunitHour:   4,
				TimeunitMinute: 5,
				TimeunitSecond: 6,
			},
			expected: time.Date(2021, 3, 4, 16, 5, 6, 0, time.UTC),
		},
		{
			name: "add negative duration",
			ref:  time.Date(2020, 6, 15, 12, 0, 0, 0, time.UTC),
			duration: Duration{
				TimeunitMonth: -3,
				TimeunitDay:   -5,
			},
			expected: time.Date(2020, 3, 10, 12, 0, 0, 0, time.UTC),
		},
		{
			name: "empty duration",
			ref:  time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
			duration: Duration{
				TimeunitDay:         0,
				TimeunitSecond:      0,
				TimeunitMillisecond: 0,
			},
			expected: time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			name: "leap year handling",
			ref:  time.Date(2020, 1, 31, 12, 0, 0, 0, time.UTC),
			duration: Duration{
				TimeunitMonth: 1,
			},
			// Jan 31 + 1 month = Mar 2 (Go's AddDate normalizes overflow)
			expected: time.Date(2020, 3, 2, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AddDuration(tt.ref, tt.duration)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestReverseDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration Duration
		expected Duration
	}{
		{
			name: "reverse positive duration",
			duration: Duration{
				TimeunitYear:  1,
				TimeunitMonth: 2,
				TimeunitDay:   3,
			},
			expected: Duration{
				TimeunitYear:  -1,
				TimeunitMonth: -2,
				TimeunitDay:   -3,
			},
		},
		{
			name: "reverse negative duration",
			duration: Duration{
				TimeunitHour:   -5,
				TimeunitMinute: -30,
			},
			expected: Duration{
				TimeunitHour:   5,
				TimeunitMinute: 30,
			},
		},
		{
			name: "reverse mixed duration",
			duration: Duration{
				TimeunitYear: 1,
				TimeunitDay:  -5,
			},
			expected: Duration{
				TimeunitYear: -1,
				TimeunitDay:  5,
			},
		},
		{
			name:     "reverse empty duration",
			duration: Duration{},
			expected: Duration{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ReverseDuration(tt.duration)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEmptyDuration(t *testing.T) {
	assert.NotNil(t, EmptyDuration)
	assert.Equal(t, float64(0), EmptyDuration[TimeunitDay])
	assert.Equal(t, float64(0), EmptyDuration[TimeunitSecond])
	assert.Equal(t, float64(0), EmptyDuration[TimeunitMillisecond])
}

func TestAddDurationDoesNotMutateInput(t *testing.T) {
	ref := time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
	duration := Duration{
		TimeunitDay: 5,
	}

	// Add duration
	AddDuration(ref, duration)

	// Verify original duration is unchanged
	assert.Equal(t, float64(5), duration[TimeunitDay])
	assert.Equal(t, 1, len(duration))

	// Verify original ref is unchanged
	assert.Equal(t, time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC), ref)
}
