package sunrise

import (
	"time"

	"github.com/nathan-osman/go-sunrise"
	"myproject/internal/models"
)

func ComputeAlarmFields(config *models.UserConfig) (nextAlarmTime, alarmDateBucket string) {
	if !config.Enabled {
		return "", "DISABLED"
	}

	alarmTime := CalculateNextSunrise(config)
	nextAlarmTime = alarmTime.UTC().Format(time.RFC3339)
	alarmDateBucket = alarmTime.UTC().Format("2006-01-02")
	return
}

func CalculateNextSunrise(config *models.UserConfig) time.Time {
	loc, err := time.LoadLocation(config.TimeZone)
	if err != nil {
		loc = time.UTC
	}

	now := time.Now().In(loc)

	// Start searching from today
	for dayOffset := 0; dayOffset < 365; dayOffset++ {
		date := now.AddDate(0, 0, dayOffset)
		weekday := int(date.Weekday())

		// Sunday = 0, Monday = 1, ..., Saturday = 6
		// day_preferences[0] is Sunday, [1] Monday, etc.
		if dayOffset < len(config.DayPreferences) && !config.DayPreferences[weekday] {
			continue
		}

		// Calculate sunrise for this date at the given latitude/longitude
		sunriseUTC, _ := sunrise.SunriseSunset(
			config.Lat,
			config.Long,
			date.Year(),
			date.Month(),
			date.Day(),
		)

		// Convert sunriseUTC (which is actually in local apparent time?)
		// The library returns sunrise time in UTC for the given date at the location.
		// Apply offset (minutes) to the sunrise time
		alarmTime := sunriseUTC.Add(time.Duration(config.Offset) * time.Minute)

		// If the alarm time is in the past relative to now (UTC), continue to next day
		if alarmTime.Before(time.Now().UTC()) {
			continue
		}

		return alarmTime
	}

	// Fallback: return a far future date (should not happen with day preferences)
	return time.Now().UTC().AddDate(1, 0, 0)
}
