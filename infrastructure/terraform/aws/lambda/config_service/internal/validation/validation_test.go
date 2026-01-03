package validation

import (
	"myproject/internal/models"
	"testing"
)

func TestValidateUserConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  models.UserConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: models.UserConfig{
				DeviceID:       "550e8400-e29b-41d4-a716-446655440000",
				Lat:            34.0522,
				Long:           -118.2437,
				FCMToken:       "dKj8X9nR3qP5vL2m:APA91b",
				DayPreferences: []bool{false, false, false, true, true, true, true},
				TimeZone:       "America/Los_Angeles",
				Enabled:        true,
			},
			wantErr: false,
		},
		{
			name: "missing deviceId",
			config: models.UserConfig{
				DeviceID:       "",
				Lat:            34.0522,
				Long:           -118.2437,
				FCMToken:       "token",
				DayPreferences: []bool{false, false, false, true, true, true, true},
				TimeZone:       "America/Los_Angeles",
				Enabled:        true,
			},
			wantErr: true,
		},
		{
			name: "deviceId too long",
			config: models.UserConfig{
				DeviceID:       string(make([]byte, 256)),
				Lat:            34.0522,
				Long:           -118.2437,
				FCMToken:       "token",
				DayPreferences: []bool{false, false, false, true, true, true, true},
				TimeZone:       "America/Los_Angeles",
				Enabled:        true,
			},
			wantErr: true,
		},
		{
			name: "lat out of range",
			config: models.UserConfig{
				DeviceID:       "id",
				Lat:            100.0,
				Long:           -118.2437,
				FCMToken:       "token",
				DayPreferences: []bool{false, false, false, true, true, true, true},
				TimeZone:       "America/Los_Angeles",
				Enabled:        true,
			},
			wantErr: true,
		},
		{
			name: "long out of range",
			config: models.UserConfig{
				DeviceID:       "id",
				Lat:            34.0522,
				Long:           -200.0,
				FCMToken:       "token",
				DayPreferences: []bool{false, false, false, true, true, true, true},
				TimeZone:       "America/Los_Angeles",
				Enabled:        true,
			},
			wantErr: true,
		},
		{
			name: "missing fcmToken",
			config: models.UserConfig{
				DeviceID:       "id",
				Lat:            34.0522,
				Long:           -118.2437,
				FCMToken:       "",
				DayPreferences: []bool{false, false, false, true, true, true, true},
				TimeZone:       "America/Los_Angeles",
				Enabled:        true,
			},
			wantErr: true,
		},
		{
			name: "fcmToken with whitespace",
			config: models.UserConfig{
				DeviceID:       "id",
				Lat:            34.0522,
				Long:           -118.2437,
				FCMToken:       "token with space",
				DayPreferences: []bool{false, false, false, true, true, true, true},
				TimeZone:       "America/Los_Angeles",
				Enabled:        true,
			},
			wantErr: true,
		},
		{
			name: "fcmToken too long",
			config: models.UserConfig{
				DeviceID:       "id",
				Lat:            34.0522,
				Long:           -118.2437,
				FCMToken:       string(make([]byte, 1025)),
				DayPreferences: []bool{false, false, false, true, true, true, true},
				TimeZone:       "America/Los_Angeles",
				Enabled:        true,
			},
			wantErr: true,
		},
		{
			name: "day preferences wrong length",
			config: models.UserConfig{
				DeviceID:       "id",
				Lat:            34.0522,
				Long:           -118.2437,
				FCMToken:       "token",
				DayPreferences: []bool{true, true, true},
				TimeZone:       "America/Los_Angeles",
				Enabled:        true,
			},
			wantErr: true,
		},
		{
			name: "invalid time zone",
			config: models.UserConfig{
				DeviceID:       "id",
				Lat:            34.0522,
				Long:           -118.2437,
				FCMToken:       "token",
				DayPreferences: []bool{false, false, false, true, true, true, true},
				TimeZone:       "Invalid/Timezone",
				Enabled:        true,
			},
			wantErr: true,
		},
		{
			name: "missing time zone",
			config: models.UserConfig{
				DeviceID:       "id",
				Lat:            34.0522,
				Long:           -118.2437,
				FCMToken:       "token",
				DayPreferences: []bool{false, false, false, true, true, true, true},
				TimeZone:       "",
				Enabled:        true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUserConfig(&tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUserConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
