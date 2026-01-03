package models

type UserConfig struct {
	DeviceID       string  `json:"deviceId" dynamodbav:"deviceId"`
	Lat            float64 `json:"lat" dynamodbav:"lat"`
	Long           float64 `json:"long" dynamodbav:"long"`
	FCMToken       string  `json:"fcmToken" dynamodbav:"fcmToken"`
	DayPreferences []bool  `json:"day_preferences" dynamodbav:"day_preferences"`
	TimeZone       string  `json:"time_zone" dynamodbav:"time_zone"`
	Enabled        bool    `json:"enabled" dynamodbav:"enabled"`
}
