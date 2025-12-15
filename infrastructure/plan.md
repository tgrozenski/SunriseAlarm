# Sunrise Alarm IaC Implementation Plan

## Overview

Multi-cloud infrastructure using Terraform to provision a weekly notification system for the SunriseAlarm Android app. AWS handles scheduling and serverless logic, GCP handles Firebase and push notifications.

---

## Phase 1: AWS Infrastructure

### 1.1 Project Setup
- [ ✓ ] Initialize Terraform project structure
- [ ✓ ] Configure AWS provider
- [ ✓ ] Set up remote state backend (S3 + DynamoDB for locking, or Terraform Cloud)
- [ ✓ ] Define input variables (region, environment, naming conventions)

### 1.2 Lambda Function
- [ ✓ ] Create IAM role for Lambda execution
- [ ✓ ] Attach policies for CloudWatch Logs, Secrets Manager read access, SNS publish
- [ ✓ ] Create Lambda function resource (placeholder code initially)
- [ ] Configure environment variables (GCP project ID, FCM topic name)
- [ ] Set appropriate timeout and memory allocation

### 1.3 EventBridge Scheduler
- [ ✓ ] Create EventBridge rule with cron expression for Saturday (e.g., )
- [ ✓ ] Create EventBridge target pointing to Lambda
- [ ✓ ] Add permissions for EventBridge to invoke Lambda

### 1.4 SNS Notifications
- [ ✓ ] Create SNS topic for job status notifications
- [ ✓ ] Create SNS subscription (email or SMS to your address)
- [ ✓ ] Update Lambda IAM role to allow SNS publish
Note: Lambda code will publish success/failure messages to this topic

### 1.5 Secrets Manager
- [ ✓ ] Create Secrets Manager secret resource (empty placeholder)
- [ ] This will be populated by GCP service account key in Phase 3
- [ ] Configure secret rotation policy (optional, for later)

### 1.6 Outputs
- [ ✓ ] Output Lambda ARN
- [ ✓ ] Output SNS topic ARN
- [ ✓ ] Output Secrets Manager secret ARN (needed for Phase 3)

---

## Phase 2: GCP/Firebase Infrastructure

### 2.1 GCP Setup
- [ ✓ ] Configure Google provider and Google Beta provider
- [ ✓ ] Enable required APIs (Firebase, Cloud Resource Manager, IAM)
- [ ✓ ] Create or reference existing GCP project

### 2.2 Firebase Project
- [ ✓ ] Enable Firebase on GCP project ( resource)
- [ ✓ ] Create Firebase Android app ( resource)
- [ ✓ ] Use your app's package name: 
- [ ] Extract for your Android project (Coming back to this)

### 2.3 Service Account for FCM Access
- [ ✓ ] Create dedicated service account for Lambda to use
- [ ✓ ] Assign  or minimal FCM permissions
- [ ✓ ] Generate service account key ( resource)

### 2.4 Outputs
- [ ✓ ] Output Firebase project ID
- [ ✓ ] Output service account email
- [ ✓ ] Output service account key (sensitive, for Phase 3)
- [ ✓ ] Automate copying of google-services.json into application directory

---

## Phase 3: Cross-Cloud Integration

### 3.1 Automatic Key Injection
- [ ✓ ] Use secrets manager to store GCP service account key
- [ ✓ ] Reference the key output from GCP module as the secret value (This will be in the lambda code)
- [ ] Ensure dependency ordering (GCP resources created before AWS secret version)

### 3.2 Verification
- [ ] Terraform plan shows correct cross-cloud dependencies
- [ ] Secret contains valid GCP service account JSON
- [ ] Lambda can retrieve secret at runtime

---

Program flow: EventBridge -> Lambda -> (retrieve secret manager) -> SNS + FCM -> Android app

## Phase 4: Lambda Function Code

> Not Terraform, but needed to complete the system

### 4.1 Implementation
- [ ] Retrieve GCP service account key from Secrets Manager
- [ ] Exchange service account key for OAuth2 access token
- [ ] Send FCM message to topic (e.g., )
- [ ] Publish result to SNS topic (success or failure with details)

### 4.2 FCM Message Structure
- [ ] Define notification payload (title, body, click action)
- [ ] Include intent data for your app's alarm reset flow
- [ ] Test with FCM HTTP v1 API

---

## Phase 5: Android App Integration

> Not Terraform, but needed to complete the system

### 5.1 Firebase Setup
- [ ] Add  to app (from Terraform output)
- [ ] Add Firebase Messaging dependency
- [ ] Subscribe to FCM topic on app launch ()

### 5.2 Notification Handling
- [ ] Implement 
- [ ] Handle notification tap → open alarm reset flow
- [ ] Read stored preferences from local storage
- [ ] Implement one-click recalculation UI

---

## Notes

- The clever part: GCP service account key flows automatically into AWS Secrets Manager through Terraform's dependency graph. One  provisions both clouds and wires them together.
- FCM topic messaging means no database needed—users self-subscribe from the app.
- SNS gives you visibility into whether the job runs; you can extend this later to include Firebase delivery receipts if needed.
- Consider adding a  vs  workspace or variable for testing.
