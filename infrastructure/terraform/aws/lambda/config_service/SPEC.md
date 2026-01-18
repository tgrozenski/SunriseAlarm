# Go Config Service

This is the spec for the our Go config service.

### Background:

We have an mobile app that sets sunrise alarms. One of the challenges with an app like this is that we need a way to store user preferences. We're going to build out a simple go web server that will run on AWS lambda. This 'config service' will handle the simple responsablilty of handing off user configs. We can implement this very straightforwardly, these configs will not be secret at all since they will not have any identifying user information. For simplicity, since these are just configs, we'll just use a deviceID pattern to identify users. This ensures that user's will distinct and the client application will be able to specify the user they need. So primarily this is a CRUD service, but all it really needs to do is Create and Update. This will be written as a standard HTTP server in go, we'll use the Lambda Web Adapter for production.

Our DynamoDB table store Map type objects that look like this:

```json
Item:
{
    "deviceId": "550e8400-e29b-41d4-a716-446655440000",
    "lat": 34.0522,
    "long": -118.2437,
    "fcmToken": "dKj8X9nR3qP5vL2m:APA91b...",
    "day_preferences": [false, false, false, true, true, true, true],
    "time_zone": "America/Los_Angeles",
    "offset" : 30,
    "nextAlarmTime": "2026-01-15T14:32:00Z",
    "alarmDateBucket": "2026-01-15",
    "enabled": true
}
```


The table will be named UserConfigs and the Partition Key will be deviceId. We will choose the containerized web server pattern to keep things simple and portable. Since we will test with mocking the DynamoDB client needs to injected with dependency injection rather than hardcoded.

Ensure the following for the config service:

### What we'll need
- A get config endpoint, receives an HTTP GET request with a deviceID as a URL parameter
  - If a config doesn't exist return a 404 not found
  - 200 with the JSON config if found

- An Upsert config endpoint, receives an HTTP POST request with a JSON body to update a current config or create a new one
  - Returns 200 with no response to indicate success (the frontend should already have the user's choice)
  - Calculates nextAlarmTime and alarmDateBucket based on lat/long/timezone/offset/day_preferences (use CalculateSunrise func)
  - Stores complete config including computed fields

  - Upsert Special Case:
    - When enabled = false:
      - Set alarmDateBucket to "DISABLED"
      - Set nextAlarmTime to empty string
      - This ensures disabled alarms are excluded from GSI queries

  - Manual Override for Testing:
    - When enabled = true AND nextAlarmTime is provided in request:
      - Use provided nextAlarmTime (must be RFC3339 format: "2026-01-15T14:32:00Z")
      - Use provided alarmDateBucket (must be YYYY-MM-DD format or "DISABLED")
      - Both fields must be provided together or neither
      - This allows integration tests to set specific alarm times for predictable testing
    - When enabled = true AND nextAlarmTime is NOT provided:
      - Calculate nextAlarmTime and alarmDateBucket using sunrise.ComputeAlarmFields()

- A validate config funciton, responsible for:

  - deviceID is a UUID string
  - Lattitude and longitude are floats within possible ranges
  - FCM token is valid format (don't need to dry run it, this will add too much latency), i.e. no whitespace, not absurdly long, and not empty
  - Day preferences is a list of exactly 7 booleans
  - Time zone is a valid IANA time zone
  - Enabled is a boolean
  - Offset is an integer between -30 and 30

- Validation errors return 400 with body: {"error": "description of what failed"}

- The handler will need to listen and serve on all these endpoints and set up anything else necessary

- One test file that will hit these endpoints with test configs
  - We will test:
    - All endpoints while with mock dynamoDB (UNIT)
    - All helper funcitons (UNIT)

### Notification Service

Additionally, we'll need a notification service. In order to keep overhead low we'll include this as part of the same web server to reduce code duplication.

This service is responsible triggering alarms set off by eventbridge crons, it will then be responsible for checking the config to determine the next time to reschedule itself. To activate alarms we'll use polling. An eventbridge cron will poll a check alarm endpoint to check if any alarms need to go off that minute. This endpoint handles sending the data notification which triggers the alarm activity in the mobile client app.

In order to efficiently poll we won't poll at all hours of the day, but rather between an interval of all possible US sunrise times. We have:

```terraform
schedule_expression = "cron(* 8-16 * * ? *)"
```

> KISS - Micheal Scott

We will primarily need this new endpoint:

GET check_alarm -> response 200 {"fired": n}
  - Queries AlarmTimeIndex GSI where:
    - alarmDateBucket = today's date (UTC, format "2006-01-02")
    - nextAlarmTime BETWEEN (now - 1 minute) AND (now + 1 minute)
  - If > 0 alarms are found:
    - Reach into secret manager to get a service key named `LAMBDA_SERVICE_KEY`
    ```terraform
    resource "aws_secretsmanager_secret" "lambda_service_key" {
      name = "LAMBDA_SERVICE_KEY"
    }
      ```
    - For each alarm:
      - Use that service key to send a DATA notification through FCM to the mobile client app
        - Try with exponential backoff
        - FCM payload is a JSON object
        ```json
        {
          "message": {
            "token": "user's_fcm_token",
            "data": {
              "type": "SUNRISE_ALARM",
              "deviceId": "abc-123"
            },
            "android": {
              "priority": "high"
            }
          }
        }
        ```
      - Recalculate the the next alarm time and store it in the config

### Additional Functions

- We will need a calculate sunrise funtion that takes the lat/long and time zone and returns the next sunrise time.

  ```go
  func CalculateSunrise(config UserConfig) time.Time {
    // TODO: implement
  }
```
```

- I want to see unit tests for this function for various sunrises in the continental US. You can leave the actual sunrise times with a placeholder, I'll source the actual data myself.

### Testing

For integration testing, the service supports manual override of `nextAlarmTime` and `alarmDateBucket` fields:

- **Manual Override**: When `enabled=true`, the POST `/config` endpoint accepts optional `nextAlarmTime` (RFC3339) and `alarmDateBucket` (YYYY-MM-DD or "DISABLED") fields
- **Test Idempotence**: Tests must clean up after themselves using AWS CLI since no delete endpoint exists
- **FCM Testing**: FCM notifications will fail in tests (no valid service key) → tests should expect `fired: 0` in `/check_alarm` responses
- **Test Framework**: Integration tests should be written as a proper test framework with:
  - Assertion functions (status codes, JSON field validation)
  - Pass/fail tracking and reporting
  - Test data cleanup using `aws dynamodb delete-item`

Example test scenario using manual overrides:
1. Create config with specific `nextAlarmTime` set to current time ±1 minute
2. Call `/check_alarm` endpoint
3. Verify alarm is detected (`fired: 1` expected, but FCM fails so `fired: 0` in practice)
4. Clean up test data with AWS CLI
