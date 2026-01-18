# Sunrise Alarm IaC Implementation Plan

## Overview

We will build a cloud native application. That is, most of the computing and logic will reside in the cloud. The android client app will be more of a simple controller and will recieve events (alarm ringing) from the cloud. Users will have a set and forget experience. We'll use cross-cloud integration to allow us to manage firebase through IaC without being forced to use GCP for the rest of the infrastructure.

---


TODO: 

- Write the spec for the notification service

## What We'll Need

## IaC

- Terraform with AWS and GCP providers

- Terraform cloud for remote state

- Terraform CLI and wrapper script to run infrastructure

### In the cloud

#### AWS

- Secrets Manager stores the service key needed by lambda to access GCP

- Dynamo DB for key value storage

- Lambda containerized config service (use lambda web adapter), stores the user's config (via the client application) in an Dynmamo DB table

- EventBridge schedules the lambda function to run at a specific time

- SNS topic as a central point for failures, all lambda instances have a policy which allows SNS publish

- Lambda containerized notification service, is triggered by an eventbride cron job
  - This triggers the alarm in the android app and recalculates next alarm and reschecules next alarm

According to timeanddate.com

summer-solstice-2026 -> June 20 or 21
easternmost point in continental US -> West Quoddy Head, Maine, USA
according to [timeanddate.com](https://www.timeanddate.com/sun/@4982705?month=6&year=2026) we have that sunrise will be 4:41 local time

Westernmost point in the continental US is -> Cape Alava Washington, 8:06 on the 21 and 22
[timeanddate.com](https://www.timeanddate.com/sun/@11822487?month=12&year=2026)

Thus UTC time is 8:41 UTC, latest is 16:06 UTC

```terraform
resource "aws_cloudwatch_event_rule" "sunrise_poll" {
  name                = "sunrise-poll"
  schedule_expression = "cron(0/1 8-16 * * ? *)"  # 08:00 - 16:00 UTC
}
```

---

#### GCP

- Manages our Firebase config

- Need the capability to send a Data message

- A service account for lambda to use

### New Services to build
- 

- Lambda config service, stores the user's config (via the client application) in an dynamoDB table (go)

- Lambda containerized notification service, triggered by an eventbridge cron job (go)

- Use docker compose to test these services

- Rehaul the android app, will be very simple UI wise probably just a few kotlin files

- Onboarding flow:

> First app launch
    ↓
Generate UUID (or hash device identifiers)
    ↓
Store locally (SharedPreferences / DataStore)
    ↓
Send to backend as device_id
    ↓
All future requests include this device_id

> Need to handle tokens rotating with firebase

> IMPORTANT: We use a FCM data message to wake the app, then locally create a notification with setFullScreenIntent() to trigger the alarm activity. This should make the alarm trigger more reliably accross different OEMs and in different power modes.
