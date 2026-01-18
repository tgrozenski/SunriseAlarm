package models

type UserConfig struct {
	DeviceID        string  `json:"deviceId" dynamodbav:"deviceId"`
	Lat             float64 `json:"lat" dynamodbav:"lat"`
	Long            float64 `json:"long" dynamodbav:"long"`
	FCMToken        string  `json:"fcmToken" dynamodbav:"fcmToken"`
	DayPreferences  []bool  `json:"day_preferences" dynamodbav:"day_preferences"`
	TimeZone        string  `json:"time_zone" dynamodbav:"time_zone"`
	Offset          int     `json:"offset" dynamodbav:"offset"`
	NextAlarmTime   string  `json:"nextAlarmTime" dynamodbav:"nextAlarmTime,omitempty"`
	AlarmDateBucket string  `json:"alarmDateBucket" dynamodbav:"alarmDateBucket,omitempty"`
	Enabled         bool    `json:"enabled" dynamodbav:"enabled"`
}
