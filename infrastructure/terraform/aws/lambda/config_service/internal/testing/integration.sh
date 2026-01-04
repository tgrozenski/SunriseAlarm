#!/bin/bash

BASE_URL="${BASE_URL:-https://hf3g6mqofv6mzoeqwhwbupduhu0phvge.lambda-url.us-west-1.on.aws}"

echo "=========================================="
echo "Config Service Integration Tests"
echo "Base URL: $BASE_URL"
echo "=========================================="

# ==========================================
# Happy Path - Create and Retrieve
# ==========================================

echo ""
echo ">>> TEST: Create a valid config"
echo ">>> EXPECTED: 200 OK (config created successfully)"
curl -s -w "\nHTTP Status: %{http_code}\n" -X POST "$BASE_URL/config" \
  -H "Content-Type: application/json" \
  -d '{
    "deviceId": "550e8400-e29b-41d4-a716-446655440000",
    "lat": 34.0522,
    "long": -118.2437,
    "fcmToken": "dKj8X9nR3qP5vL2mAPA91btest1234567890",
    "day_preferences": [false, false, false, true, true, true, true],
    "time_zone": "America/Los_Angeles",
    "enabled": true
  }'

echo ""
echo ">>> TEST: Retrieve the config we just created"
echo ">>> EXPECTED: 200 OK with config JSON body"
curl -s -w "\nHTTP Status: %{http_code}\n" "$BASE_URL/config/550e8400-e29b-41d4-a716-446655440000"

# ==========================================
# Not Found
# ==========================================

echo ""
echo ">>> TEST: Retrieve a config that doesn't exist"
echo ">>> EXPECTED: 404 Not Found"
curl -s -w "\nHTTP Status: %{http_code}\n" "$BASE_URL/config/00000000-0000-0000-0000-000000000000"

# ==========================================
# Validation: Invalid UUID
# ==========================================

echo ""
echo ">>> TEST: Invalid deviceId (not a UUID)"
echo ">>> EXPECTED: 400 Bad Request - deviceId must be a valid UUID"
curl -s -w "\nHTTP Status: %{http_code}\n" -X POST "$BASE_URL/config" \
  -H "Content-Type: application/json" \
  -d '{
    "deviceId": "not-a-valid-uuid",
    "lat": 34.0522,
    "long": -118.2437,
    "fcmToken": "validtoken123",
    "day_preferences": [false, false, false, true, true, true, true],
    "time_zone": "America/Los_Angeles",
    "enabled": true
  }'

# ==========================================
# Validation: Latitude out of range
# ==========================================

echo ""
echo ">>> TEST: Latitude too high (> 90)"
echo ">>> EXPECTED: 400 Bad Request - latitude must be between -90 and 90"
curl -s -w "\nHTTP Status: %{http_code}\n" -X POST "$BASE_URL/config" \
  -H "Content-Type: application/json" \
  -d '{
    "deviceId": "550e8400-e29b-41d4-a716-446655440001",
    "lat": 91.0,
    "long": -118.2437,
    "fcmToken": "validtoken123",
    "day_preferences": [false, false, false, true, true, true, true],
    "time_zone": "America/Los_Angeles",
    "enabled": true
  }'

echo ""
echo ">>> TEST: Latitude too low (< -90)"
echo ">>> EXPECTED: 400 Bad Request - latitude must be between -90 and 90"
curl -s -w "\nHTTP Status: %{http_code}\n" -X POST "$BASE_URL/config" \
  -H "Content-Type: application/json" \
  -d '{
    "deviceId": "550e8400-e29b-41d4-a716-446655440001",
    "lat": -91.0,
    "long": -118.2437,
    "fcmToken": "validtoken123",
    "day_preferences": [false, false, false, true, true, true, true],
    "time_zone": "America/Los_Angeles",
    "enabled": true
  }'

# ==========================================
# Validation: Longitude out of range
# ==========================================

echo ""
echo ">>> TEST: Longitude too high (> 180)"
echo ">>> EXPECTED: 400 Bad Request - longitude must be between -180 and 180"
curl -s -w "\nHTTP Status: %{http_code}\n" -X POST "$BASE_URL/config" \
  -H "Content-Type: application/json" \
  -d '{
    "deviceId": "550e8400-e29b-41d4-a716-446655440001",
    "lat": 34.0522,
    "long": 181.0,
    "fcmToken": "validtoken123",
    "day_preferences": [false, false, false, true, true, true, true],
    "time_zone": "America/Los_Angeles",
    "enabled": true
  }'

echo ""
echo ">>> TEST: Longitude too low (< -180)"
echo ">>> EXPECTED: 400 Bad Request - longitude must be between -180 and 180"
curl -s -w "\nHTTP Status: %{http_code}\n" -X POST "$BASE_URL/config" \
  -H "Content-Type: application/json" \
  -d '{
    "deviceId": "550e8400-e29b-41d4-a716-446655440001",
    "lat": 34.0522,
    "long": -181.0,
    "fcmToken": "validtoken123",
    "day_preferences": [false, false, false, true, true, true, true],
    "time_zone": "America/Los_Angeles",
    "enabled": true
  }'

# ==========================================
# Validation: FCM Token
# ==========================================

echo ""
echo ">>> TEST: Empty FCM token"
echo ">>> EXPECTED: 400 Bad Request - fcmToken cannot be empty"
curl -s -w "\nHTTP Status: %{http_code}\n" -X POST "$BASE_URL/config" \
  -H "Content-Type: application/json" \
  -d '{
    "deviceId": "550e8400-e29b-41d4-a716-446655440001",
    "lat": 34.0522,
    "long": -118.2437,
    "fcmToken": "",
    "day_preferences": [false, false, false, true, true, true, true],
    "time_zone": "America/Los_Angeles",
    "enabled": true
  }'

echo ""
echo ">>> TEST: FCM token with whitespace"
echo ">>> EXPECTED: 400 Bad Request - fcmToken cannot contain whitespace"
curl -s -w "\nHTTP Status: %{http_code}\n" -X POST "$BASE_URL/config" \
  -H "Content-Type: application/json" \
  -d '{
    "deviceId": "550e8400-e29b-41d4-a716-446655440001",
    "lat": 34.0522,
    "long": -118.2437,
    "fcmToken": "token with spaces",
    "day_preferences": [false, false, false, true, true, true, true],
    "time_zone": "America/Los_Angeles",
    "enabled": true
  }'

echo ""
echo ">>> TEST: FCM token with tabs/newlines"
echo ">>> EXPECTED: 400 Bad Request - fcmToken cannot contain whitespace"
curl -s -w "\nHTTP Status: %{http_code}\n" -X POST "$BASE_URL/config" \
  -H "Content-Type: application/json" \
  -d '{
    "deviceId": "550e8400-e29b-41d4-a716-446655440001",
    "lat": 34.0522,
    "long": -118.2437,
    "fcmToken": "token\twith\ttabs",
    "day_preferences": [false, false, false, true, true, true, true],
    "time_zone": "America/Los_Angeles",
    "enabled": true
  }'

# ==========================================
# Validation: Day Preferences
# ==========================================

echo ""
echo ">>> TEST: day_preferences with too few elements (3 instead of 7)"
echo ">>> EXPECTED: 400 Bad Request - day_preferences must have exactly 7 elements"
curl -s -w "\nHTTP Status: %{http_code}\n" -X POST "$BASE_URL/config" \
  -H "Content-Type: application/json" \
  -d '{
    "deviceId": "550e8400-e29b-41d4-a716-446655440001",
    "lat": 34.0522,
    "long": -118.2437,
    "fcmToken": "validtoken123",
    "day_preferences": [false, false, false],
    "time_zone": "America/Los_Angeles",
    "enabled": true
  }'

echo ""
echo ">>> TEST: day_preferences with too many elements (10 instead of 7)"
echo ">>> EXPECTED: 400 Bad Request - day_preferences must have exactly 7 elements"
curl -s -w "\nHTTP Status: %{http_code}\n" -X POST "$BASE_URL/config" \
  -H "Content-Type: application/json" \
  -d '{
    "deviceId": "550e8400-e29b-41d4-a716-446655440001",
    "lat": 34.0522,
    "long": -118.2437,
    "fcmToken": "validtoken123",
    "day_preferences": [false, false, false, true, true, true, true, false, false, false],
    "time_zone": "America/Los_Angeles",
    "enabled": true
  }'

echo ""
echo ">>> TEST: day_preferences empty array"
echo ">>> EXPECTED: 400 Bad Request - day_preferences must have exactly 7 elements"
curl -s -w "\nHTTP Status: %{http_code}\n" -X POST "$BASE_URL/config" \
  -H "Content-Type: application/json" \
  -d '{
    "deviceId": "550e8400-e29b-41d4-a716-446655440001",
    "lat": 34.0522,
    "long": -118.2437,
    "fcmToken": "validtoken123",
    "day_preferences": [],
    "time_zone": "America/Los_Angeles",
    "enabled": true
  }'

# ==========================================
# Validation: Time Zone
# ==========================================

echo ""
echo ">>> TEST: Invalid time zone"
echo ">>> EXPECTED: 400 Bad Request - time_zone must be a valid IANA timezone"
curl -s -w "\nHTTP Status: %{http_code}\n" -X POST "$BASE_URL/config" \
  -H "Content-Type: application/json" \
  -d '{
    "deviceId": "550e8400-e29b-41d4-a716-446655440001",
    "lat": 34.0522,
    "long": -118.2437,
    "fcmToken": "validtoken123",
    "day_preferences": [false, false, false, true, true, true, true],
    "time_zone": "Not/A_Real_Timezone",
    "enabled": true
  }'

echo ""
echo ">>> TEST: Empty time zone"
echo ">>> EXPECTED: 400 Bad Request - time_zone must be a valid IANA timezone"
curl -s -w "\nHTTP Status: %{http_code}\n" -X POST "$BASE_URL/config" \
  -H "Content-Type: application/json" \
  -d '{
    "deviceId": "550e8400-e29b-41d4-a716-446655440001",
    "lat": 34.0522,
    "long": -118.2437,
    "fcmToken": "validtoken123",
    "day_preferences": [false, false, false, true, true, true, true],
    "time_zone": "",
    "enabled": true
  }'

# ==========================================
# Validation: Malformed JSON
# ==========================================

echo ""
echo ">>> TEST: Malformed JSON body"
echo ">>> EXPECTED: 400 Bad Request - invalid JSON"
curl -s -w "\nHTTP Status: %{http_code}\n" -X POST "$BASE_URL/config" \
  -H "Content-Type: application/json" \
  -d 'this is not json at all'

echo ""
echo ">>> TEST: Missing required fields"
echo ">>> EXPECTED: 400 Bad Request - missing required fields"
curl -s -w "\nHTTP Status: %{http_code}\n" -X POST "$BASE_URL/config" \
  -H "Content-Type: application/json" \
  -d '{
    "deviceId": "550e8400-e29b-41d4-a716-446655440001"
  }'

# ==========================================
# Update existing config
# ==========================================

echo ""
echo ">>> TEST: Update existing config (upsert)"
echo ">>> EXPECTED: 200 OK"
curl -s -w "\nHTTP Status: %{http_code}\n" -X POST "$BASE_URL/config" \
  -H "Content-Type: application/json" \
  -d '{
    "deviceId": "550e8400-e29b-41d4-a716-446655440000",
    "lat": 40.7128,
    "long": -74.0060,
    "fcmToken": "updatedtoken9876543210",
    "day_preferences": [true, true, true, true, true, false, false],
    "time_zone": "America/New_York",
    "enabled": false
  }'

echo ""
echo ">>> TEST: Verify config was updated"
echo ">>> EXPECTED: 200 OK with updated values (New York timezone, enabled=false)"
curl -s -w "\nHTTP Status: %{http_code}\n" "$BASE_URL/config/550e8400-e29b-41d4-a716-446655440000"

echo ""
echo "=========================================="
echo "Tests complete!"
echo "=========================================="
