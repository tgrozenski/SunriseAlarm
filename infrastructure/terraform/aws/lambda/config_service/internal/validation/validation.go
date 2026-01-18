package validation

import (
	"errors"
	"myproject/internal/models"
	"regexp"
	"strings"
	"time"
)

var uuidRegex = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

func ValidateUserConfig(config *models.UserConfig) error {
	if config.DeviceID == "" {
		return errors.New("deviceId is required")
	}
	if len(config.DeviceID) > 255 {
		return errors.New("deviceId is too long")
	}
	if !uuidRegex.MatchString(config.DeviceID) {
		return errors.New("deviceId must be a valid UUID v4")
	}

	if config.Lat < -90 || config.Lat > 90 {
		return errors.New("lat must be between -90 and 90")
	}
	if config.Long < -180 || config.Long > 180 {
		return errors.New("long must be between -180 and 180")
	}

	if config.FCMToken == "" {
		return errors.New("fcmToken is required")
	}
	if strings.ContainsAny(config.FCMToken, " \t\n\r") {
		return errors.New("fcmToken cannot contain whitespace")
	}
	if len(config.FCMToken) > 1024 {
		return errors.New("fcmToken is too long")
	}

	if len(config.DayPreferences) != 7 {
		return errors.New("day_preferences must have exactly 7 values")
	}

	if config.TimeZone == "" {
		return errors.New("time_zone is required")
	}
	if _, err := time.LoadLocation(config.TimeZone); err != nil {
		return errors.New("time_zone must be a valid IANA time zone")
	}

	if config.Offset < -30 || config.Offset > 30 {
		return errors.New("offset must be between -30 and 30")
	}

	// Validate NextAlarmTime format if provided
	if config.NextAlarmTime != "" {
		if _, err := time.Parse(time.RFC3339, config.NextAlarmTime); err != nil {
			return errors.New("nextAlarmTime must be in RFC3339 format (e.g., 2026-01-15T14:32:00Z)")
		}
	}

	// Validate AlarmDateBucket format if provided
	if config.AlarmDateBucket != "" && config.AlarmDateBucket != "DISABLED" {
		if _, err := time.Parse("2006-01-02", config.AlarmDateBucket); err != nil {
			return errors.New("alarmDateBucket must be in YYYY-MM-DD format or 'DISABLED'")
		}
	}

	return nil
}
