# Sunrise Alarm Config Service API

## Overview
The Sunrise Alarm Config Service is a Go web server deployed as an AWS Lambda function that manages user configuration for sunrise alarms. It provides endpoints for storing/retrieving user configurations and triggering sunrise alarm notifications via Firebase Cloud Messaging (FCM).

## Base URL
The service is accessible at the Lambda Function URL (deployed via AWS Lambda). The base URL will be provided after deployment.

## Authentication
No authentication is required for the public endpoints. The service uses device ID (UUID) to identify users.

## Endpoints

### 1. GET `/config/{deviceId}`
Retrieves a user's configuration by device ID.

#### Path Parameters
| Parameter | Type   | Required | Description                     |
|-----------|--------|----------|---------------------------------|
| deviceId  | string | Yes      | UUID v4 format (e.g., `550e8400-e29b-41d4-a716-446655440000`) |

#### Responses
**200 OK** - Configuration found
```json
{
  "deviceId": "550e8400-e29b-41d4-a716-446655440000",
  "lat": 34.0522,
  "long": -118.2437,
  "fcmToken": "dKj8X9nR3qP5vL2m:APA91b...",
  "day_preferences": [false, false, false, true, true, true, true],
  "time_zone": "America/Los_Angeles",
  "offset": 30,
  "nextAlarmTime": "2026-01-15T14:32:00Z",
  "alarmDateBucket": "2026-01-15",
  "enabled": true
}
```

**400 Bad Request** - Invalid device ID format
```json
{
  "error": "deviceId is required"
}
```

**404 Not Found** - Configuration not found
```json
{
  "error": "config not found"
}
```

**500 Internal Server Error**
```json
{
  "error": "failed to retrieve config"
}
```

---

### 2. POST `/config`
Creates or updates a user configuration (upsert). The service automatically calculates the next alarm time based on latitude, longitude, timezone, and day preferences.

#### Request Body
```json
{
  "deviceId": "550e8400-e29b-41d4-a716-446655440000",
  "lat": 34.0522,
  "long": -118.2437,
  "fcmToken": "dKj8X9nR3qP5vL2m:APA91b...",
  "day_preferences": [false, false, false, true, true, true, true],
  "time_zone": "America/Los_Angeles",
  "offset": 30,
  "enabled": true,
  "nextAlarmTime": "2026-01-15T14:32:00Z",
  "alarmDateBucket": "2026-01-15"
}
```

**Note:** `nextAlarmTime` and `alarmDateBucket` are optional fields used for testing. When provided with `enabled=true`, they allow manual override of calculated alarm times.

#### Field Validation
| Field | Type | Required | Validation Rules |
|-------|------|----------|------------------|
| deviceId | string | Yes | UUID v4 format, max 255 chars |
| lat | float64 | Yes | Between -90 and 90 inclusive |
| long | float64 | Yes | Between -180 and 180 inclusive |
| fcmToken | string | Yes | No whitespace, max 1024 chars, not empty |
| day_preferences | boolean[] | Yes | Exactly 7 values (Sunday to Saturday) |
| time_zone | string | Yes | Valid IANA time zone (e.g., "America/Los_Angeles") |
| offset | integer | Yes | Between -30 and 30 inclusive (minutes offset from sunrise) |
| enabled | boolean | Yes | true/false |
| nextAlarmTime | string | No | RFC3339 format (e.g., "2026-01-15T14:32:00Z"), required if alarmDateBucket provided |
| alarmDateBucket | string | No | YYYY-MM-DD format or "DISABLED", required if nextAlarmTime provided |

**Note:** When `enabled` is `false`:
- `nextAlarmTime` is set to empty string
- `alarmDateBucket` is set to "DISABLED"
- Alarm will not trigger
- Any provided `nextAlarmTime` or `alarmDateBucket` values are ignored

**Manual Override for Testing:** When `enabled` is `true` and `nextAlarmTime` is provided:
- Uses provided `nextAlarmTime` (RFC3339 format) instead of calculating from sunrise
- Uses provided `alarmDateBucket` (YYYY-MM-DD format)
- Both fields must be provided together or neither
- Allows integration tests to set specific alarm times for predictable testing

#### Responses
**200 OK** - Configuration saved successfully
- No response body

**400 Bad Request** - Validation error
```json
{
  "error": "deviceId must be a valid UUID v4"
}
```

**400 Bad Request** - Invalid JSON
```json
{
  "error": "invalid JSON"
}
```

**500 Internal Server Error**
```json
{
  "error": "failed to save config"
}
```

---

### 3. GET `/check_alarm`
**Internal Endpoint** - Called by EventBridge cron every minute to check for alarms that need to fire.

Queries the DynamoDB GSI `AlarmTimeIndex` for alarms where:
- `alarmDateBucket` = today's date (UTC format "2006-01-02")
- `nextAlarmTime` is between (now - 1 minute) and (now + 1 minute)

For each matching alarm:
1. Sends an FCM data notification to the user's device
2. Recalculates the next alarm time and updates the configuration

#### FCM Notification Payload
```json
{
  "message": {
    "token": "user_fcm_token",
    "data": {
      "type": "SUNRISE_ALARM",
      "deviceId": "550e8400-e29b-41d4-a716-446655440000"
    },
    "android": {
      "priority": "high"
    }
  }
}
```

#### Responses
**200 OK** - Success (with or without alarms)
```json
{
  "fired": 2
}
```

**500 Internal Server Error**
```json
{
  "error": "failed to query alarms"
}
```

---

### 4. GET `/health`
Health check endpoint.

#### Responses
**200 OK**
```
OK
```

## Data Types

### UserConfig Object
| Field | Type | Description |
|-------|------|-------------|
| deviceId | string | Unique device identifier (UUID v4) |
| lat | float64 | Latitude for sunrise calculation |
| long | float64 | Longitude for sunrise calculation |
| fcmToken | string | Firebase Cloud Messaging token for notifications |
| day_preferences | boolean[] | Array of 7 booleans for days of the week (Sunday to Saturday) |
| time_zone | string | IANA time zone for alarm time calculation |
| offset | integer | Minutes offset from sunrise (-30 to +30) |
| nextAlarmTime | string | RFC3339 timestamp of next alarm (calculated or manually overridden for testing) |
| alarmDateBucket | string | Date bucket for GSI querying (YYYY-MM-DD or "DISABLED") |
| enabled | boolean | Whether alarm is active |

### Day Preferences Format
The `day_preferences` array represents days of the week:
- Index 0: Sunday
- Index 1: Monday
- Index 2: Tuesday
- Index 3: Wednesday
- Index 4: Thursday
- Index 5: Friday
- Index 6: Saturday

Example: `[false, true, true, true, true, true, false]` means alarm is enabled Monday-Friday only.

## Error Handling
All error responses follow the same format:
```json
{
  "error": "human-readable error message"
}
```

Common HTTP status codes:
- **400** - Bad Request (validation errors, invalid JSON)
- **404** - Not Found (config not found)
- **500** - Internal Server Error (database errors, FCM failures)

## Examples

### 1. Create/Update Configuration
```bash
curl -X POST https://[lambda-url]/config \
  -H "Content-Type: application/json" \
  -d '{
    "deviceId": "550e8400-e29b-41d4-a716-446655440000",
    "lat": 34.0522,
    "long": -118.2437,
    "fcmToken": "dKj8X9nR3qP5vL2m:APA91b...",
    "day_preferences": [false, false, false, true, true, true, true],
    "time_zone": "America/Los_Angeles",
    "offset": 30,
    "enabled": true
  }'
```

### 2. Retrieve Configuration
```bash
curl https://[lambda-url]/config/550e8400-e29b-41d4-a716-446655440000
```

### 3. Disable Alarm
```bash
curl -X POST https://[lambda-url]/config \
  -H "Content-Type: application/json" \
  -d '{
    "deviceId": "550e8400-e29b-41d4-a716-446655440000",
    "lat": 34.0522,
    "long": -118.2437,
    "fcmToken": "dKj8X9nR3qP5vL2m:APA91b...",
    "day_preferences": [false, false, false, true, true, true, true],
    "time_zone": "America/Los_Angeles",
    "offset": 30,
    "enabled": false
  }'
```

### 4. Health Check
```bash
curl https://[lambda-url]/health
```

## Implementation Notes

### Sunrise Calculation
The service uses the `go-sunrise` library to calculate sunrise times based on:
- Latitude/longitude
- Current date and time zone
- Day preferences (skips disabled days)
- Offset (minutes before/after sunrise)

The next alarm time is calculated by finding the next enabled day with a sunrise time that hasn't passed.

### Database Schema
- **Table**: `UserConfigs`
- **Partition Key**: `deviceId` (string)
- **GSI**: `AlarmTimeIndex`
  - Hash Key: `alarmDateBucket` (string)
  - Range Key: `nextAlarmTime` (string)
- **Special Value**: When `enabled=false`, `alarmDateBucket` is set to "DISABLED" to exclude from GSI queries

### EventBridge Integration
- **Schedule**: Runs every minute (`rate(1 minute)`)
- **Target**: Lambda function `/check_alarm` endpoint
- **Purpose**: Polls for alarms that need to fire in the current minute window

### FCM Integration
- **Secret**: `LAMBDA_SERVICE_KEY` stored in AWS Secrets Manager
- **Retry**: Exponential backoff with 3 retries
- **Priority**: High priority for Android notifications
- **Payload Type**: Data notification (`SUNRISE_ALARM`)

### Testing with Manual Overrides
For integration testing, the service supports manual override of alarm times:
- **Optional Fields**: `nextAlarmTime` (RFC3339) and `alarmDateBucket` (YYYY-MM-DD) can be provided when `enabled=true`
- **Use Case**: Allows tests to set specific alarm times for predictable testing of the `/check_alarm` endpoint
- **Validation**: Both fields must be provided together or neither; validated for correct format
- **Disabled State**: When `enabled=false`, any provided override values are ignored and fields are cleared
- **Test Framework**: See `integration.sh` for a complete test framework with assertions and cleanup