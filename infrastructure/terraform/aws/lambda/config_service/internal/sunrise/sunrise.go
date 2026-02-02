package sunrise

import (
	"errors"
	"time"

	"github.com/nathan-osman/go-sunrise"
	"myproject/internal/models"
)

func ComputeAlarmFields(config *models.UserConfig) (nextAlarmTime, alarmDateBucket string, err error) {
	if !config.Enabled {
		return "", "DISABLED", nil
	}

	alarmTime, err := CalculateNextSunrise(config)
	if err != nil {
		return "", "", err
	}
	nextAlarmTime = alarmTime.UTC().Format(time.RFC3339)
	alarmDateBucket = alarmTime.UTC().Format("2006-01-02")
	return nextAlarmTime, alarmDateBucket, nil
}

func CalculateNextSunrise(config *models.UserConfig) (time.Time, error) {
	return calculateNextSunriseAtTime(config, time.Now().UTC())
}

func calculateNextSunriseAtTime(config *models.UserConfig, nowUTC time.Time) (time.Time, error) {
	loc, err := time.LoadLocation(config.TimeZone)
	if err != nil {
		loc = time.UTC
	}

	now := nowUTC.In(loc)

	// Start searching from today up to next week (8 days total)
	for dayOffset := 0; dayOffset <= 7; dayOffset++ {
		date := now.AddDate(0, 0, dayOffset)
		weekday := int(date.Weekday())

		// Sunday = 0, Monday = 1, ..., Saturday = 6
		// day_preferences[0] is Sunday, [1] Monday, etc.
		if !config.DayPreferences[weekday] {
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
		if alarmTime.Before(nowUTC) {
			continue
		}

		return alarmTime, nil
	}

	// Should not reach here with proper validation (at least one day enabled)
	return time.Time{}, errors.New("no enabled day found within next week")
}
