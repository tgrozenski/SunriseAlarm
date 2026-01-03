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

- A validate config funciton, responsible for:

  - deviceID is a UUID string
  - Lattitude and longitude are floats within possible ranges
  - FCM token is valid format (don't need to dry run it, this will add too much latency), i.e. no whitespace, not absurdly long, and not empty
  - Day preferences is a list of exactly 7 booleans
  - Time zone is a valid IANA time zone
  - Enabled is a boolean

- Validation errors return 400 with body: {"error": "description of what failed"}

- The handler will need to listen and serve on all these endpoints and set up anything else necessary

- One test file that will hit these endpoints with test configs
  - We will test:
    - All endpoints while with mock dynamoDB (UNIT)
    - All helper funcitons (UNIT)
