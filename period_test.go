package kronos

import (
	"testing"
	"time"
)

// TestPeriodString tests the String() method of Period
func TestPeriodString(t *testing.T) {
	tests := []struct {
		period   Period
		expected string
	}{
		{PeriodYear, "year"},
		{PeriodMonth, "month"},
		{PeriodWeek, "week"},
		{PeriodDay, "day"},
		{PeriodTime, "time"},
		{PeriodUnknown, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.period.String(); got != tt.expected {
				t.Errorf("Period.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestDeterminePeriodFromDuration tests the DeterminePeriodFromDuration function
func TestDeterminePeriodFromDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration Duration
		expected Period
	}{
		{
			name:     "nil duration",
			duration: nil,
			expected: PeriodDay,
		},
		{
			name:     "seconds",
			duration: Duration{TimeunitSecond: 30},
			expected: PeriodTime,
		},
		{
			name:     "minutes",
			duration: Duration{TimeunitMinute: 15},
			expected: PeriodTime,
		},
		{
			name:     "hours",
			duration: Duration{TimeunitHour: 2},
			expected: PeriodTime,
		},
		{
			name:     "days",
			duration: Duration{TimeunitDay: 3},
			expected: PeriodDay,
		},
		{
			name:     "weeks",
			duration: Duration{TimeunitWeek: 1},
			expected: PeriodWeek,
		},
		{
			name:     "months",
			duration: Duration{TimeunitMonth: 2},
			expected: PeriodMonth,
		},
		{
			name:     "years",
			duration: Duration{TimeunitYear: 1},
			expected: PeriodYear,
		},
		{
			name:     "decades",
			duration: Duration{TimeunitDecade: 1},
			expected: PeriodYear,
		},
		{
			name:     "quarters",
			duration: Duration{TimeunitQuarter: 2},
			expected: PeriodYear,
		},
		{
			name:     "mixed year and month",
			duration: Duration{TimeunitYear: 1, TimeunitMonth: 2},
			expected: PeriodMonth, // Finest granularity
		},
		{
			name:     "mixed month and day",
			duration: Duration{TimeunitMonth: 1, TimeunitDay: 3},
			expected: PeriodDay, // Finest granularity
		},
		{
			name:     "mixed day and hours",
			duration: Duration{TimeunitDay: 1, TimeunitHour: 2},
			expected: PeriodTime, // Finest granularity
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := determinePeriodFromDuration(tt.duration); got != tt.expected {
				t.Errorf("determinePeriodFromDuration() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestDeterminePeriodFromComponents tests the DeterminePeriodFromComponents function
func TestDeterminePeriodFromComponents(t *testing.T) {
	refTime := time.Date(2020, 3, 15, 14, 30, 0, 0, time.UTC)
	ref := newReferenceWithTimezone(refTime, nil)

	tests := []struct {
		name     string
		setup    func() *ParsingComponents
		expected Period
	}{
		{
			name: "nil components",
			setup: func() *ParsingComponents {
				return nil
			},
			expected: PeriodUnknown,
		},
		{
			name: "time components certain",
			setup: func() *ParsingComponents {
				pc := newParsingComponents(ref, nil)
				pc.Assign(ComponentHour, 10)
				pc.Assign(ComponentMinute, 30)
				return pc
			},
			expected: PeriodTime,
		},
		{
			name: "day certain",
			setup: func() *ParsingComponents {
				pc := newParsingComponents(ref, nil)
				pc.Assign(ComponentYear, 2020)
				pc.Assign(ComponentMonth, 3)
				pc.Assign(ComponentDay, 15)
				return pc
			},
			expected: PeriodDay,
		},
		{
			name: "weekday certain without day",
			setup: func() *ParsingComponents {
				pc := newParsingComponents(ref, nil)
				pc.Assign(ComponentWeekday, int(WeekdayMonday))
				return pc
			},
			expected: PeriodWeek,
		},
		{
			name: "month certain without day",
			setup: func() *ParsingComponents {
				pc := newParsingComponents(ref, nil)
				pc.Assign(ComponentYear, 2020)
				pc.Assign(ComponentMonth, 3)
				return pc
			},
			expected: PeriodMonth,
		},
		{
			name: "only year certain",
			setup: func() *ParsingComponents {
				pc := newParsingComponents(ref, nil)
				pc.Assign(ComponentYear, 2020)
				return pc
			},
			expected: PeriodYear,
		},
		{
			name: "nothing certain",
			setup: func() *ParsingComponents {
				pc := newParsingComponents(ref, nil)
				return pc
			},
			expected: PeriodUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pc := tt.setup()
			if got := determinePeriodFromComponents(pc); got != tt.expected {
				t.Errorf("determinePeriodFromComponents() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestCasualReferencesPeriod tests that casual reference functions set appropriate periods
func TestCasualReferencesPeriod(t *testing.T) {
	refTime := time.Date(2020, 3, 15, 14, 30, 0, 0, time.UTC)
	ref := newReferenceWithTimezone(refTime, nil)

	tests := []struct {
		name     string
		fn       func() *ParsingComponents
		expected Period
	}{
		{"Now", func() *ParsingComponents { return now(ref) }, PeriodTime},
		{"Today", func() *ParsingComponents { return today(ref) }, PeriodDay},
		{"Yesterday", func() *ParsingComponents { return yesterday(ref) }, PeriodDay},
		{"Tomorrow", func() *ParsingComponents { return tomorrow(ref) }, PeriodDay},
		{"TheDayAfter", func() *ParsingComponents { return theDayAfter(ref, 2) }, PeriodDay},
		{"TheDayBefore", func() *ParsingComponents { return theDayBefore(ref, 1) }, PeriodDay},
		{"Tonight", func() *ParsingComponents { return tonight(ref) }, PeriodDay},
		{"LastNight", func() *ParsingComponents { return lastNight(ref) }, PeriodDay},
		{"Morning", func() *ParsingComponents { return morning(ref) }, PeriodTime},
		{"Afternoon", func() *ParsingComponents { return afternoon(ref) }, PeriodTime},
		{"Evening", func() *ParsingComponents { return evening(ref) }, PeriodTime},
		{"Midnight", func() *ParsingComponents { return midnight(ref) }, PeriodTime},
		{"Noon", func() *ParsingComponents { return noon(ref) }, PeriodTime},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pc := tt.fn()
			if got := pc.Period(); got != tt.expected {
				t.Errorf("%s() period = %v, want %v", tt.name, got, tt.expected)
			}
		})
	}
}

// TestRelativeDatePeriod tests that CreateRelativeFromReference sets appropriate periods
func TestRelativeDatePeriod(t *testing.T) {
	refTime := time.Date(2020, 3, 15, 14, 30, 0, 0, time.UTC)
	ref := newReferenceWithTimezone(refTime, nil)

	tests := []struct {
		name     string
		duration Duration
		expected Period
	}{
		{
			name:     "2 hours ago",
			duration: Duration{TimeunitHour: -2},
			expected: PeriodTime,
		},
		{
			name:     "3 days ago",
			duration: Duration{TimeunitDay: -3},
			expected: PeriodDay,
		},
		{
			name:     "1 week ago",
			duration: Duration{TimeunitWeek: -1},
			expected: PeriodWeek,
		},
		{
			name:     "2 months ago",
			duration: Duration{TimeunitMonth: -2},
			expected: PeriodMonth,
		},
		{
			name:     "1 year ago",
			duration: Duration{TimeunitYear: -1},
			expected: PeriodYear,
		},
		{
			name:     "1 year 2 months ago",
			duration: Duration{TimeunitYear: -1, TimeunitMonth: -2},
			expected: PeriodMonth, // Finest granularity
		},
		{
			name:     "1 week 3 days ago",
			duration: Duration{TimeunitWeek: -1, TimeunitDay: -3},
			expected: PeriodDay, // Finest granularity
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pc := createRelativeFromReference(ref, tt.duration)
			if got := pc.Period(); got != tt.expected {
				t.Errorf("createRelativeFromReference(%v) period = %v, want %v",
					tt.duration, got, tt.expected)
			}
		})
	}
}

// TestPeriodClone tests that Period is properly copied when cloning ParsingComponents
func TestPeriodClone(t *testing.T) {
	refTime := time.Date(2020, 3, 15, 14, 30, 0, 0, time.UTC)
	ref := newReferenceWithTimezone(refTime, nil)

	original := newParsingComponents(ref, nil)
	original.SetPeriod(PeriodMonth)

	clone := original.Clone()

	if clone.Period() != PeriodMonth {
		t.Errorf("Clone() period = %v, want %v", clone.Period(), PeriodMonth)
	}

	// Verify that modifying the clone doesn't affect the original
	clone.SetPeriod(PeriodDay)

	if original.Period() != PeriodMonth {
		t.Errorf("After modifying clone, original period = %v, want %v",
			original.Period(), PeriodMonth)
	}
}

// TestPeriodGetterSetter tests the Period getter and setter methods
func TestPeriodGetterSetter(t *testing.T) {
	refTime := time.Date(2020, 3, 15, 14, 30, 0, 0, time.UTC)
	ref := newReferenceWithTimezone(refTime, nil)

	pc := newParsingComponents(ref, nil)

	// Default should be PeriodUnknown
	if pc.Period() != PeriodUnknown {
		t.Errorf("Default period = %v, want %v", pc.Period(), PeriodUnknown)
	}

	// Test setting various periods
	periods := []Period{
		PeriodYear,
		PeriodMonth,
		PeriodWeek,
		PeriodDay,
		PeriodTime,
	}

	for _, period := range periods {
		t.Run(period.String(), func(t *testing.T) {
			pc.SetPeriod(period)
			if got := pc.Period(); got != period {
				t.Errorf("After SetPeriod(%v), Period() = %v, want %v",
					period, got, period)
			}
		})
	}
}
