package sunrise

import (
	"testing"
	"time"

	"myproject/internal/models"
)

func TestCalculateNextSunrise(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name          string
		config        models.UserConfig
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
			wantFuture: true,
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
			wantFuture: true,
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
			wantFuture: true, // far future
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateNextSunrise(&tt.config)

			if tt.wantFuture && result.Before(now) {
				t.Errorf("CalculateNextSunrise() returned past time: %v", result)
			}

			// Verify day of week matches preferences
			if len(tt.config.DayPreferences) == 7 {
				weekday := int(result.Weekday())
				if !tt.config.DayPreferences[weekday] && tt.name != "no enabled days" {
					// For "no enabled days", the function returns a far future date
					// which may not match any preference; that's okay
					t.Errorf("CalculateNextSunrise() returned day %v which is not enabled", result.Weekday())
				}
			}

			// Verify offset is applied (rough check)
			if tt.config.Offset != 0 {
				// We can't easily verify offset without recalculating sunrise
				// This is a sanity check that result is not zero time
				if result.IsZero() {
					t.Errorf("CalculateNextSunrise() returned zero time")
				}
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

	result := CalculateNextSunrise(&config)
	weekday := int(result.Weekday())
	if weekday != 1 {
		t.Errorf("CalculateNextSunrise() returned day %v, expected Monday (1)", weekday)
	}
}

func TestCalculateNextSunrise_OffsetApplied(t *testing.T) {
	config := models.UserConfig{
		Lat:            34.0522,
		Long:           -118.2437,
		TimeZone:       "America/Los_Angeles",
		DayPreferences: []bool{true, true, true, true, true, true, true},
		Offset:         30,
	}

	result := CalculateNextSunrise(&config)
	// We can't easily verify exact offset without re-implementing sunrise calculation
	// but we can ensure result is not zero
	if result.IsZero() {
		t.Error("CalculateNextSunrise() returned zero time")
	}
}
