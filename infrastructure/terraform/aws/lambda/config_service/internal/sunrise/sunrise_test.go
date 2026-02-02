package sunrise

import (
	"fmt"
	"testing"
	"time"

	"github.com/nathan-osman/go-sunrise"
	"myproject/internal/models"
)

func absDuration(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}

func TestCalculateNextSunrise(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name          string
		config        models.UserConfig
		wantError     bool
		wantFuture    bool
		wantDayOfWeek int // 0=Sunday
		checkOffset   bool
	}{
		{
			name: "LA weekend enabled",
			config: models.UserConfig{
				Lat:            34.0522,
				Long:           -118.2437,
				TimeZone:       "America/Los_Angeles",
				DayPreferences: []bool{true, false, false, false, false, false, true}, // Sun & Sat
				Offset:         0,
			},
			wantFuture: true,
			wantError:  false,
		},
		{
			name: "NYC weekday with offset +30",
			config: models.UserConfig{
				Lat:            40.7128,
				Long:           -74.0060,
				TimeZone:       "America/New_York",
				DayPreferences: []bool{false, true, true, true, true, true, false}, // Mon-Fri
				Offset:         30,
			},
			wantFuture:  true,
			wantError:   false,
			checkOffset: true,
		},
		{
			name: "Chicago with negative offset",
			config: models.UserConfig{
				Lat:            41.8781,
				Long:           -87.6298,
				TimeZone:       "America/Chicago",
				DayPreferences: []bool{true, true, true, true, true, true, true}, // every day
				Offset:         -15,
			},
			wantFuture:  true,
			wantError:   false,
			checkOffset: true,
		},
		{
			name: "no enabled days",
			config: models.UserConfig{
				Lat:            34.0522,
				Long:           -118.2437,
				TimeZone:       "America/Los_Angeles",
				DayPreferences: []bool{false, false, false, false, false, false, false},
				Offset:         0,
			},
			wantError:  true, // should return error
			wantFuture: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CalculateNextSunrise(&tt.config)

			if tt.wantError {
				if err == nil {
					t.Errorf("CalculateNextSunrise() expected error but got none")
				}
				return // Skip further checks
			}
			if err != nil {
				t.Errorf("CalculateNextSunrise() unexpected error: %v", err)
				return
			}

			if tt.wantFuture && result.Before(now) {
				t.Errorf("CalculateNextSunrise() returned past time: %v", result)
			}

			// Verify day of week matches preferences
			if len(tt.config.DayPreferences) == 7 {
				weekday := int(result.Weekday())
				if !tt.config.DayPreferences[weekday] {
					t.Errorf("CalculateNextSunrise() returned day %v which is not enabled", result.Weekday())
				}
			}

			// Verify offset is applied
			if tt.checkOffset && tt.config.Offset != 0 {
				// Create a copy of config with offset 0
				zeroOffsetConfig := tt.config
				zeroOffsetConfig.Offset = 0
				zeroOffsetResult, err := CalculateNextSunrise(&zeroOffsetConfig)
				if err != nil {
					t.Errorf("CalculateNextSunrise() with offset 0 failed: %v", err)
					return
				}
				// If both results selected the same day (within 24 hours), verify offset difference
				if absDuration(result.Sub(zeroOffsetResult)) < 24*time.Hour {
					expectedDiff := time.Duration(tt.config.Offset) * time.Minute
					actualDiff := result.Sub(zeroOffsetResult)
					tolerance := 1 * time.Second
					if absDuration(actualDiff-expectedDiff) > tolerance {
						t.Errorf("CalculateNextSunrise() offset not applied correctly: expected difference %v, got %v", expectedDiff, actualDiff)
					}
				}
				// If different days selected (edge case), skip offset verification
			}
		})
	}
}

func TestCalculateNextSunrise_RespectsDayPreferences(t *testing.T) {
	config := models.UserConfig{
		Lat:            34.0522,
		Long:           -118.2437,
		TimeZone:       "America/Los_Angeles",
		DayPreferences: []bool{false, true, false, false, false, false, false}, // Monday only
		Offset:         0,
	}

	result, err := CalculateNextSunrise(&config)
	if err != nil {
		t.Errorf("CalculateNextSunrise() unexpected error: %v", err)
		return
	}
	weekday := int(result.Weekday())
	if weekday != 1 {
		t.Errorf("CalculateNextSunrise() returned day %v, expected Monday (1)", weekday)
	}
}

func TestCalculateNextSunrise_OffsetApplied(t *testing.T) {
	baseConfig := models.UserConfig{
		Lat:            34.0522,
		Long:           -118.2437,
		TimeZone:       "America/Los_Angeles",
		DayPreferences: []bool{true, true, true, true, true, true, true},
		Offset:         0,
	}

	// Test various offsets, including edge cases
	offsets := []int{30, -15, 0, 30, -30}
	for _, offset := range offsets {
		t.Run(fmt.Sprintf("offset_%d", offset), func(t *testing.T) {
			config := baseConfig
			config.Offset = offset
			result, err := CalculateNextSunrise(&config)
			if err != nil {
				t.Errorf("CalculateNextSunrise() unexpected error: %v", err)
				return
			}
			if result.IsZero() {
				t.Error("CalculateNextSunrise() returned zero time")
				return
			}

			// Compare with offset 0 to verify offset application
			zeroResult, err := CalculateNextSunrise(&baseConfig)
			if err != nil {
				t.Errorf("CalculateNextSunrise() with offset 0 failed: %v", err)
				return
			}
			// If same day selected (within 24 hours), verify offset difference
			if absDuration(result.Sub(zeroResult)) < 24*time.Hour {
				expectedDiff := time.Duration(offset) * time.Minute
				actualDiff := result.Sub(zeroResult)
				tolerance := 1 * time.Second
				if absDuration(actualDiff-expectedDiff) > tolerance {
					t.Errorf("CalculateNextSunrise() offset not applied correctly: expected difference %v, got %v", expectedDiff, actualDiff)
				}
			}
			// If different days selected (edge case), skip offset verification
			// This can happen if offset pushes alarm across day boundary relative to now
		})
	}
}

func TestSunriseCalculationAgainstKnownValues(t *testing.T) {
	tests := []struct {
		name           string
		city           string
		state          string
		latitude       float64
		longitude      float64
		timeZone       string
		year           int
		month          time.Month
		day            int
		expectedHour   int
		expectedMinute int
		toleranceSec   int
	}{
		{
			name:           "Los Angeles CA January 1",
			city:           "Los Angeles",
			state:          "CA",
			latitude:       34.0522,
			longitude:      -118.2437,
			timeZone:       "America/Los_Angeles",
			year:           2026,
			month:          time.January,
			day:            1,
			expectedHour:   6,
			expectedMinute: 58,
			toleranceSec:   120, // 2 minutes tolerance
		},
		{
			name:           "Portland OR December 31",
			city:           "Portland",
			state:          "OR",
			latitude:       45.5152,
			longitude:      -122.6784,
			timeZone:       "America/Los_Angeles", // Portland is in Pacific Time
			year:           2025,
			month:          time.December,
			day:            31,
			expectedHour:   7,
			expectedMinute: 50,
			toleranceSec:   120,
		},
		{
			name:           "Augusta ME January 10",
			city:           "Augusta",
			state:          "ME",
			latitude:       44.3106,
			longitude:      -69.7795,
			timeZone:       "America/New_York", // Augusta is in Eastern Time
			year:           2026,
			month:          time.January,
			day:            10,
			expectedHour:   7,
			expectedMinute: 13,
			toleranceSec:   120,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Calculate sunrise using the library
			sunriseUTC, _ := sunrise.SunriseSunset(
				tt.latitude,
				tt.longitude,
				tt.year,
				tt.month,
				tt.day,
			)

			// Convert to local timezone
			loc, err := time.LoadLocation(tt.timeZone)
			if err != nil {
				t.Errorf("failed to load location %s: %v", tt.timeZone, err)
				return
			}
			sunriseLocal := sunriseUTC.In(loc)

			// Create expected time
			expectedTime := time.Date(
				tt.year,
				tt.month,
				tt.day,
				tt.expectedHour,
				tt.expectedMinute,
				0, 0, loc,
			)

			// Calculate difference
			diff := sunriseLocal.Sub(expectedTime)
			if diff < 0 {
				diff = -diff
			}

			if diff > time.Duration(tt.toleranceSec)*time.Second {
				t.Errorf(
					"sunrise calculation mismatch for %s: expected %s, got %s (difference %v, tolerance %vs)",
					tt.name,
					expectedTime.Format("15:04"),
					sunriseLocal.Format("15:04"),
					diff,
					tt.toleranceSec,
				)
			}
		})
	}
}

func TestCalculateNextSunrise_SkipsPastSunrise(t *testing.T) {
	// Test that when sunrise has already passed today, it skips to next enabled day
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		t.Fatalf("failed to load location: %v", err)
	}

	// Use January 1, 2026 in Los Angeles (sunrise ~06:58)
	testDate := time.Date(2026, time.January, 1, 0, 0, 0, 0, loc)

	// Calculate sunrise for January 1
	sunriseUTC, _ := sunrise.SunriseSunset(
		34.0522,   // LA latitude
		-118.2437, // LA longitude
		testDate.Year(),
		testDate.Month(),
		testDate.Day(),
	)

	// Set "now" to be 1 hour after sunrise (07:58 local, which is 15:58 UTC)
	// January 1, 2026 07:58 AM PST = January 1, 2026 15:58 UTC
	nowUTC := sunriseUTC.Add(1 * time.Hour)

	// Create config with all days enabled
	config := models.UserConfig{
		Lat:            34.0522,
		Long:           -118.2437,
		TimeZone:       "America/Los_Angeles",
		DayPreferences: []bool{true, true, true, true, true, true, true},
		Offset:         0,
	}

	// Calculate next sunrise given that "now" is after today's sunrise
	result, err := calculateNextSunriseAtTime(&config, nowUTC)
	if err != nil {
		t.Errorf("calculateNextSunriseAtTime() unexpected error: %v", err)
		return
	}

	if result.IsZero() {
		t.Error("calculateNextSunriseAtTime() returned zero time")
		return
	}

	// Result should be for January 2, 2026 (next day)
	expectedDate := time.Date(2026, time.January, 2, 0, 0, 0, 0, loc)
	// Calculate sunrise for January 2
	expectedSunriseUTC, _ := sunrise.SunriseSunset(
		34.0522,
		-118.2437,
		expectedDate.Year(),
		expectedDate.Month(),
		expectedDate.Day(),
	)

	// Allow small tolerance for floating point calculations
	tolerance := 1 * time.Second
	if absDuration(result.Sub(expectedSunriseUTC)) > tolerance {
		t.Errorf("calculateNextSunriseAtTime() did not skip to next day: expected %v (Jan 2 sunrise), got %v",
			expectedSunriseUTC.Format("2006-01-02 15:04:05 MST"),
			result.Format("2006-01-02 15:04:05 MST"))
	}
}

func TestCalculateNextSunrise_SkipsPastSunriseWithOffset(t *testing.T) {
	// Test that offset doesn't affect day-skip logic incorrectly
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		t.Fatalf("failed to load location: %v", err)
	}

	// Use January 1, 2026
	testDate := time.Date(2026, time.January, 1, 0, 0, 0, 0, loc)

	// Calculate sunrise for January 1
	sunriseUTC, _ := sunrise.SunriseSunset(
		34.0522,
		-118.2437,
		testDate.Year(),
		testDate.Month(),
		testDate.Day(),
	)

	// Apply a negative offset (earlier alarm)
	offset := -30 // 30 minutes before sunrise
	alarmTimeWithOffset := sunriseUTC.Add(time.Duration(offset) * time.Minute)

	// Set "now" to be 1 hour after the OFFSET alarm time
	// i.e., if sunrise is 06:58, offset alarm is 06:28, now is 07:28
	nowUTC := alarmTimeWithOffset.Add(1 * time.Hour)

	// Create config with negative offset
	config := models.UserConfig{
		Lat:            34.0522,
		Long:           -118.2437,
		TimeZone:       "America/Los_Angeles",
		DayPreferences: []bool{true, true, true, true, true, true, true},
		Offset:         offset,
	}

	result, err := calculateNextSunriseAtTime(&config, nowUTC)
	if err != nil {
		t.Errorf("calculateNextSunriseAtTime() with offset unexpected error: %v", err)
		return
	}

	// Should still skip to January 2
	expectedDate := time.Date(2026, time.January, 2, 0, 0, 0, 0, loc)
	expectedSunriseUTC, _ := sunrise.SunriseSunset(
		34.0522,
		-118.2437,
		expectedDate.Year(),
		expectedDate.Month(),
		expectedDate.Day(),
	)
	expectedAlarmTime := expectedSunriseUTC.Add(time.Duration(offset) * time.Minute)

	tolerance := 1 * time.Second
	if absDuration(result.Sub(expectedAlarmTime)) > tolerance {
		t.Errorf("calculateNextSunriseAtTime() with offset did not skip correctly: expected %v, got %v",
			expectedAlarmTime.Format("2006-01-02 15:04:05 MST"),
			result.Format("2006-01-02 15:04:05 MST"))
	}
}

func TestCalculateNextSunrise_SkipsToNextEnabledDay(t *testing.T) {
	// Test skipping multiple days when only specific weekdays are enabled
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		t.Fatalf("failed to load location: %v", err)
	}

	// January 1, 2026 is a Thursday (weekday 4)
	thursdayDate := time.Date(2026, time.January, 1, 0, 0, 0, 0, loc)

	// Calculate sunrise for Thursday
	thursdaySunriseUTC, _ := sunrise.SunriseSunset(
		34.0522,
		-118.2437,
		thursdayDate.Year(),
		thursdayDate.Month(),
		thursdayDate.Day(),
	)

	// Set "now" to be 1 hour after Thursday's sunrise
	nowUTC := thursdaySunriseUTC.Add(1 * time.Hour)

	// Create config with only Sundays enabled (weekday 0)
	// Thursday is weekday 4, so should skip to Sunday (3 days later: Friday, Saturday, Sunday)
	sundayOnlyPrefs := []bool{true, false, false, false, false, false, false} // Sun only

	config := models.UserConfig{
		Lat:            34.0522,
		Long:           -118.2437,
		TimeZone:       "America/Los_Angeles",
		DayPreferences: sundayOnlyPrefs,
		Offset:         0,
	}

	result, err := calculateNextSunriseAtTime(&config, nowUTC)
	if err != nil {
		t.Errorf("calculateNextSunriseAtTime() unexpected error: %v", err)
		return
	}

	// Next Sunday after Jan 1, 2026 is Jan 4, 2026
	expectedDate := time.Date(2026, time.January, 4, 0, 0, 0, 0, loc)
	expectedSunriseUTC, _ := sunrise.SunriseSunset(
		34.0522,
		-118.2437,
		expectedDate.Year(),
		expectedDate.Month(),
		expectedDate.Day(),
	)

	tolerance := 1 * time.Second
	if absDuration(result.Sub(expectedSunriseUTC)) > tolerance {
		t.Errorf("calculateNextSunriseAtTime() did not skip to next enabled day: expected %v (Jan 4, Sunday sunrise), got %v",
			expectedSunriseUTC.Format("2006-01-02 15:04:05 MST"),
			result.Format("2006-01-02 15:04:05 MST"))
	}

	// Verify result is a Sunday (weekday 0)
	if result.Weekday() != time.Sunday {
		t.Errorf("calculateNextSunriseAtTime() returned wrong weekday: expected Sunday, got %v", result.Weekday())
	}
}
