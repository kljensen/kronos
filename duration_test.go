package kronos

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddDuration(t *testing.T) {
	tests := []struct {
		name        string
		ref         time.Time
		duration    Duration
		expected    time.Time
		expectError bool
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
			result, err := AddDuration(tt.ref, tt.duration)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
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
	_, err := AddDuration(ref, duration)
	assert.NoError(t, err)

	// Verify original duration is unchanged
	assert.Equal(t, float64(5), duration[TimeunitDay])
	assert.Equal(t, 1, len(duration))

	// Verify original ref is unchanged
	assert.Equal(t, time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC), ref)
}

func TestAddDuration_BoundsChecking(t *testing.T) {
	tests := []struct {
		name        string
		ref         time.Time
		duration    Duration
		expectError bool
		errorMsg    string
	}{
		{
			name: "extreme positive years",
			ref:  time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
			duration: Duration{
				TimeunitYear: 15000,
			},
			expectError: true,
			errorMsg:    "exceeds maximum",
		},
		{
			name: "extreme negative years",
			ref:  time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
			duration: Duration{
				TimeunitYear: -15000,
			},
			expectError: true,
			errorMsg:    "exceeds maximum",
		},
		{
			name: "result year too high",
			ref:  time.Date(5000, 1, 1, 12, 0, 0, 0, time.UTC),
			duration: Duration{
				TimeunitYear: 5000,
			},
			expectError: true,
			errorMsg:    "outside valid range",
		},
		{
			name: "result year too low",
			ref:  time.Date(500, 1, 1, 12, 0, 0, 0, time.UTC),
			duration: Duration{
				TimeunitYear: -500,
			},
			expectError: true,
			errorMsg:    "outside valid range",
		},
		{
			name: "input year too high",
			ref:  time.Date(10000, 1, 1, 12, 0, 0, 0, time.UTC),
			duration: Duration{
				TimeunitDay: 1,
			},
			expectError: true,
			errorMsg:    "outside valid range",
		},
		{
			name: "input year zero",
			ref:  time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC),
			duration: Duration{
				TimeunitDay: 1,
			},
			expectError: true,
			errorMsg:    "outside valid range",
		},
		{
			name: "extreme months",
			ref:  time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
			duration: Duration{
				TimeunitMonth: 150000,
			},
			expectError: true,
			errorMsg:    "exceeds maximum",
		},
		{
			name: "extreme days",
			ref:  time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
			duration: Duration{
				TimeunitDay: 5000000,
			},
			expectError: true,
			errorMsg:    "exceeds maximum",
		},
		{
			name: "extreme hours",
			ref:  time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
			duration: Duration{
				TimeunitHour: 100000000,
			},
			expectError: true,
			errorMsg:    "exceeds maximum",
		},
		{
			name: "NaN value",
			ref:  time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
			duration: Duration{
				TimeunitYear: math.NaN(),
			},
			expectError: true,
			errorMsg:    "invalid value",
		},
		{
			name: "positive infinity",
			ref:  time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
			duration: Duration{
				TimeunitYear: math.Inf(1),
			},
			expectError: true,
			errorMsg:    "invalid value",
		},
		{
			name: "negative infinity",
			ref:  time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
			duration: Duration{
				TimeunitYear: math.Inf(-1),
			},
			expectError: true,
			errorMsg:    "invalid value",
		},
		{
			name: "cascading overflow - decades to years",
			ref:  time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
			duration: Duration{
				TimeunitDecade: 1500,
			},
			expectError: true,
			errorMsg:    "cascading",
		},
		{
			name: "cascading overflow - years to months",
			ref:  time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
			duration: Duration{
				TimeunitYear:  5000,
				TimeunitMonth: 50000,
			},
			expectError: true,
			errorMsg:    "outside valid range",
		},
		{
			name: "float to int overflow",
			ref:  time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
			duration: Duration{
				TimeunitYear: math.MaxInt32 + 1.0,
			},
			expectError: true,
			errorMsg:    "exceeds maximum",
		},
		{
			name: "boundary - near max year",
			ref:  time.Date(9998, 1, 1, 12, 0, 0, 0, time.UTC),
			duration: Duration{
				TimeunitYear: 1,
			},
			expectError: false,
		},
		{
			name: "boundary - at max year",
			ref:  time.Date(9999, 1, 1, 12, 0, 0, 0, time.UTC),
			duration: Duration{
				TimeunitDay: 1,
			},
			expectError: false,
		},
		{
			name: "boundary - near min year",
			ref:  time.Date(2, 1, 1, 12, 0, 0, 0, time.UTC),
			duration: Duration{
				TimeunitYear: -1,
			},
			expectError: false,
		},
		{
			name: "boundary - at min year",
			ref:  time.Date(1, 1, 1, 12, 0, 0, 0, time.UTC),
			duration: Duration{
				TimeunitDay: 1,
			},
			expectError: false,
		},
		{
			name: "cumulative overflow from multiple units",
			ref:  time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
			duration: Duration{
				TimeunitYear:    5000,
				TimeunitMonth:   60000,
				TimeunitQuarter: 10000,
			},
			expectError: true,
			errorMsg:    "cascading",
		},
		{
			name: "large but valid duration",
			ref:  time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
			duration: Duration{
				TimeunitYear:  100,
				TimeunitMonth: 6,
				TimeunitDay:   15,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := AddDuration(tt.ref, tt.duration)
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateDate(t *testing.T) {
	tests := []struct {
		name        string
		date        time.Time
		expectError bool
	}{
		{
			name:        "valid date",
			date:        time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
			expectError: false,
		},
		{
			name:        "min year",
			date:        time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
			expectError: false,
		},
		{
			name:        "max year",
			date:        time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
			expectError: false,
		},
		{
			name:        "year too low",
			date:        time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC),
			expectError: true,
		},
		{
			name:        "year too high",
			date:        time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDate(tt.date)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateDuration(t *testing.T) {
	tests := []struct {
		name        string
		duration    Duration
		expectError bool
	}{
		{
			name: "valid duration",
			duration: Duration{
				TimeunitYear:  10,
				TimeunitMonth: 5,
				TimeunitDay:   3,
			},
			expectError: false,
		},
		{
			name: "NaN value",
			duration: Duration{
				TimeunitYear: math.NaN(),
			},
			expectError: true,
		},
		{
			name: "positive infinity",
			duration: Duration{
				TimeunitMonth: math.Inf(1),
			},
			expectError: true,
		},
		{
			name: "negative infinity",
			duration: Duration{
				TimeunitDay: math.Inf(-1),
			},
			expectError: true,
		},
		{
			name: "year exceeds maximum",
			duration: Duration{
				TimeunitYear: MaxYearsDuration + 1,
			},
			expectError: true,
		},
		{
			name: "month exceeds maximum",
			duration: Duration{
				TimeunitMonth: MaxMonthsDuration + 1,
			},
			expectError: true,
		},
		{
			name: "day exceeds maximum",
			duration: Duration{
				TimeunitDay: MaxDaysDuration + 1,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDuration(tt.duration)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCheckCascadingOverflow(t *testing.T) {
	tests := []struct {
		name        string
		duration    Duration
		expectError bool
	}{
		{
			name: "valid cascading",
			duration: Duration{
				TimeunitYear:  10,
				TimeunitMonth: 5,
			},
			expectError: false,
		},
		{
			name: "decade to year overflow",
			duration: Duration{
				TimeunitDecade: 1500,
			},
			expectError: true,
		},
		{
			name: "year to month overflow",
			duration: Duration{
				TimeunitYear:  5000,
				TimeunitMonth: 70000,
			},
			expectError: true,
		},
		{
			name: "quarter to month overflow",
			duration: Duration{
				TimeunitQuarter: 50000,
			},
			expectError: true,
		},
		{
			name: "week to day overflow",
			duration: Duration{
				TimeunitWeek: 600000,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkCascadingOverflow(tt.duration)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
