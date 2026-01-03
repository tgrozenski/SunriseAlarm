package validation

import (
	"errors"
	"myproject/internal/models"
	"strings"
	"time"
)

func ValidateUserConfig(config *models.UserConfig) error {
	if config.DeviceID == "" {
		return errors.New("deviceId is required")
	}
	if len(config.DeviceID) > 255 {
		return errors.New("deviceId is too long")
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

	return nil
}
