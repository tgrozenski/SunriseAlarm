#!/usr/bin/env python3
"""
Integration tests for the Sunrise Alarm Config Service.

This test suite validates the HTTP API endpoints for the Go Config Service
deployed as an AWS Lambda function. It tests CRUD operations, validation,
manual override features, and the alarm checking endpoint.

Environment Variables:
    BASE_URL: The base URL of the deployed service (default: test Lambda URL)
    DYNAMODB_TABLE: DynamoDB table name (default: UserConfigs)
    AWS_REGION: AWS region (default: us-west-1)

Dependencies:
    pip install pytest requests boto3
"""

import os
import json
import uuid
import subprocess
from datetime import datetime, timedelta
from typing import List, Dict, Any, Optional, Tuple, Generator
import sys
import time

# Third-party imports
import pytest
import requests

# Try to import boto3 for cleanup, but make it optional
try:
    import boto3
    from botocore.exceptions import ClientError

    BOTO3_AVAILABLE = True
except ImportError:
    BOTO3_AVAILABLE = False

# ==========================================
# Configuration
# ==========================================

BASE_URL = os.environ.get(
    "BASE_URL", "https://hf3g6mqofv6mzoeqwhwbupduhu0phvge.lambda-url.us-west-1.on.aws"
)
DYNAMODB_TABLE = os.environ.get("DYNAMODB_TABLE", "UserConfigs")
AWS_REGION = os.environ.get("AWS_REGION", "us-west-1")

# Timeout for HTTP requests (seconds)
REQUEST_TIMEOUT = 30

# ==========================================
# Test Helper Classes
# ==========================================


class TestContext:
    """Context for tracking test data and cleanup."""

    def __init__(self):
        self.devices_to_cleanup: List[str] = []
        self.last_request: Optional[str] = None
        self.last_response: Optional[str] = None

    def register_cleanup(self, device_id: str) -> None:
        """Register a device ID for cleanup after tests."""
        self.devices_to_cleanup.append(device_id)

    def cleanup_all(self) -> None:
        """Clean up all registered test data."""
        if not self.devices_to_cleanup:
            return

        print(f"Cleaning up {len(self.devices_to_cleanup)} test devices...")

        # Try boto3 first
        if BOTO3_AVAILABLE:
            self._cleanup_with_boto3()
        else:
            # Fall back to AWS CLI
            self._cleanup_with_cli()

        self.devices_to_cleanup.clear()
        print("Cleanup complete")

    def _cleanup_with_boto3(self) -> None:
        """Clean up using boto3."""
        if not BOTO3_AVAILABLE:
            self._cleanup_with_cli()
            return

        try:
            dynamodb = boto3.client("dynamodb", region_name=AWS_REGION)
            for device_id in self.devices_to_cleanup:
                try:
                    dynamodb.delete_item(
                        TableName=DYNAMODB_TABLE, Key={"deviceId": {"S": device_id}}
                    )
                except ClientError as e:
                    print(f"Warning: Failed to delete device {device_id}: {e}")
        except Exception as e:
            print(f"Warning: boto3 cleanup failed: {e}")
            # Fall back to CLI
            self._cleanup_with_cli()

    def _cleanup_with_cli(self) -> None:
        """Clean up using AWS CLI."""
        for device_id in self.devices_to_cleanup:
            try:
                subprocess.run(
                    [
                        "aws",
                        "dynamodb",
                        "delete-item",
                        "--table-name",
                        DYNAMODB_TABLE,
                        "--region",
                        AWS_REGION,
                        "--key",
                        json.dumps({"deviceId": {"S": device_id}}),
                        "--output",
                        "json",
                    ],
                    capture_output=True,
                    timeout=10,
                    check=False,
                )
            except (subprocess.SubprocessError, FileNotFoundError):
                pass  # AWS CLI not available or command failed


# ==========================================
# Fixtures
# ==========================================


@pytest.fixture(scope="session")
def test_context() -> TestContext:
    """Session-level fixture for test context."""
    return TestContext()


@pytest.fixture(scope="function", autouse=True)
def auto_cleanup(test_context: TestContext) -> Generator[None, None, None]:
    """Auto-cleanup after each test."""
    yield
    test_context.cleanup_all()


# ==========================================
# Assertion Helpers
# ==========================================


class AssertionError(Exception):
    """Custom exception for assertion failures with detailed context."""

    def __init__(
        self,
        message: str,
        expected: Any = None,
        actual: Any = None,
        response: Any = None,
        request: Optional[str] = None,
    ):
        self.message = message
        self.expected = expected
        self.actual = actual
        self.response = response
        self.request = request

        # Build detailed error message
        details = [f"❌ {message}"]
        if expected is not None:
            details.append(f"Expected: {expected}")
        if actual is not None:
            details.append(f"Actual: {actual}")
        if response is not None:
            details.append(f"Response: {response}")
        if request is not None:
            details.append(f"Request: {request}")

        super().__init__("\n".join(details))


def assert_status(
    response: requests.Response, expected_status: int, context: TestContext
) -> None:
    """Assert that response has expected status code."""
    if response.status_code != expected_status:
        raise AssertionError(
            f"Status code mismatch",
            expected=f"Status {expected_status}",
            actual=f"Status {response.status_code}",
            response=response.text,
            request=context.last_request,
        )


def assert_json_field(
    response: requests.Response,
    field_path: str,
    expected_value: Any,
    context: TestContext,
) -> None:
    """Assert that JSON response contains field with expected value."""
    try:
        data = response.json()
    except json.JSONDecodeError:
        raise AssertionError(
            f"Response is not valid JSON",
            expected=f"Field {field_path} = {expected_value}",
            actual="Invalid JSON",
            response=response.text,
            request=context.last_request,
        )

    # Navigate nested fields using dot notation
    value = data
    for part in field_path.strip(".").split("."):
        if part:
            if isinstance(value, dict) and part in value:
                value = value[part]
            else:
                raise AssertionError(
                    f"Field not found in response",
                    expected=f"Field {field_path} = {expected_value}",
                    actual=f"Field {part} not found",
                    response=data,
                    request=context.last_request,
                )

    if value != expected_value:
        raise AssertionError(
            f"Field value mismatch",
            expected=f"Field {field_path} = {expected_value}",
            actual=f"Field {field_path} = {value}",
            response=data,
            request=context.last_request,
        )


def assert_json_contains(
    response: requests.Response, field_path: str, context: TestContext
) -> None:
    """Assert that JSON response contains specified field."""
    try:
        data = response.json()
    except json.JSONDecodeError:
        raise AssertionError(
            f"Response is not valid JSON",
            expected=f"JSON contains field {field_path}",
            actual="Invalid JSON",
            response=response.text,
            request=context.last_request,
        )

    # Check if field exists
    value = data
    for part in field_path.strip(".").split("."):
        if part:
            if isinstance(value, dict) and part in value:
                value = value[part]
            else:
                raise AssertionError(
                    f"Field not found in response",
                    expected=f"JSON contains field {field_path}",
                    actual=f"Field {part} not found",
                    response=data,
                    request=context.last_request,
                )


# ==========================================
# HTTP Request Helper
# ==========================================


def make_request(
    method: str,
    endpoint: str,
    data: Optional[Dict] = None,
    context: Optional[TestContext] = None,
) -> requests.Response:
    """
    Make HTTP request to the service.

    Args:
        method: HTTP method (GET, POST, etc.)
        endpoint: API endpoint path
        data: Request body data (will be JSON serialized)
        context: TestContext for tracking

    Returns:
        requests.Response object
    """
    url = f"{BASE_URL.rstrip('/')}/{endpoint.lstrip('/')}"

    # Store request info for debugging
    request_info = f"{method} {url}"
    if data:
        request_info += f"\nData: {json.dumps(data, indent=2)}"

    if context:
        context.last_request = request_info

    try:
        headers = {"Content-Type": "application/json"}

        if method.upper() == "GET":
            response = requests.get(url, headers=headers, timeout=REQUEST_TIMEOUT)
        elif method.upper() == "POST":
            response = requests.post(
                url,
                headers=headers,
                data=json.dumps(data) if data else None,
                timeout=REQUEST_TIMEOUT,
            )
        else:
            raise ValueError(f"Unsupported HTTP method: {method}")

        if context:
            context.last_response = response.text

        return response

    except requests.exceptions.RequestException as e:
        error_msg = f"HTTP request failed: {e}"
        if context:
            raise AssertionError(error_msg, request=request_info)
        else:
            raise RuntimeError(error_msg)


# ==========================================
# Test Data Generation
# ==========================================


def generate_uuid() -> str:
    """Generate a random UUID v4."""
    return str(uuid.uuid4())


def get_current_time_utc() -> str:
    """Get current UTC time in RFC3339 format."""
    return datetime.utcnow().strftime("%Y-%m-%dT%H:%M:%SZ")


def get_future_time_utc(minutes: int = 5) -> str:
    """Get future UTC time in RFC3339 format."""
    future = datetime.utcnow() + timedelta(minutes=minutes)
    return future.strftime("%Y-%m-%dT%H:%M:%SZ")


def get_past_time_utc(minutes: int = 5) -> str:
    """Get past UTC time in RFC3339 format."""
    past = datetime.utcnow() - timedelta(minutes=minutes)
    return past.strftime("%Y-%m-%dT%H:%M:%SZ")


def get_time_within_window(seconds_offset: int = -30) -> str:
    """Get UTC time within the check_alarm query window (now ± 1 minute)."""
    target = datetime.utcnow() + timedelta(seconds=seconds_offset)
    return target.strftime("%Y-%m-%dT%H:%M:%SZ")


def get_today_date() -> str:
    """Get today's date in YYYY-MM-DD format."""
    return datetime.utcnow().strftime("%Y-%m-%d")


# ==========================================
# Test Cases
# ==========================================


class TestHappyPath:
    """Tests for successful create and retrieve operations."""

    def test_create_and_retrieve_config(self, test_context: TestContext) -> None:
        """Create a valid config and retrieve it."""
        device_id = generate_uuid()
        test_context.register_cleanup(device_id)

        # Create config
        response = make_request(
            "POST",
            "config",
            {
                "deviceId": device_id,
                "lat": 34.0522,
                "long": -118.2437,
                "fcmToken": "dKj8X9nR3qP5vL2mAPA91btest1234567890",
                "day_preferences": [False, False, False, True, True, True, True],
                "time_zone": "America/Los_Angeles",
                "enabled": True,
            },
            test_context,
        )

        assert_status(response, 200, test_context)

        # Retrieve config
        response = make_request("GET", f"config/{device_id}", context=test_context)

        assert_status(response, 200, test_context)
        assert_json_field(response, "deviceId", device_id, test_context)
        assert_json_field(response, "lat", 34.0522, test_context)
        assert_json_field(response, "time_zone", "America/Los_Angeles", test_context)
        assert_json_field(response, "enabled", True, test_context)

    def test_not_found(self, test_context: TestContext) -> None:
        """Attempt to retrieve a config that doesn't exist."""
        device_id = "00000000-0000-0000-0000-000000000000"

        response = make_request("GET", f"config/{device_id}", context=test_context)

        assert_status(response, 404, test_context)


class TestValidation:
    """Tests for input validation."""

    @pytest.fixture
    def base_config(self) -> Dict:
        """Base configuration for validation tests."""
        return {
            "deviceId": generate_uuid(),
            "lat": 34.0522,
            "long": -118.2437,
            "fcmToken": "validtoken123",
            "day_preferences": [False, False, False, True, True, True, True],
            "time_zone": "America/Los_Angeles",
            "enabled": True,
        }

    def test_latitude_too_high(
        self, test_context: TestContext, base_config: Dict
    ) -> None:
        """Test latitude validation (> 90)."""
        base_config["lat"] = 91.0
        test_context.register_cleanup(base_config["deviceId"])

        response = make_request("POST", "config", base_config, test_context)
        assert_status(response, 400, test_context)

    def test_latitude_too_low(
        self, test_context: TestContext, base_config: Dict
    ) -> None:
        """Test latitude validation (< -90)."""
        base_config["lat"] = -91.0
        test_context.register_cleanup(base_config["deviceId"])

        response = make_request("POST", "config", base_config, test_context)
        assert_status(response, 400, test_context)

    def test_longitude_too_high(
        self, test_context: TestContext, base_config: Dict
    ) -> None:
        """Test longitude validation (> 180)."""
        base_config["long"] = 181.0
        test_context.register_cleanup(base_config["deviceId"])

        response = make_request("POST", "config", base_config, test_context)
        assert_status(response, 400, test_context)

    def test_longitude_too_low(
        self, test_context: TestContext, base_config: Dict
    ) -> None:
        """Test longitude validation (< -180)."""
        base_config["long"] = -181.0
        test_context.register_cleanup(base_config["deviceId"])

        response = make_request("POST", "config", base_config, test_context)
        assert_status(response, 400, test_context)

    def test_empty_fcm_token(
        self, test_context: TestContext, base_config: Dict
    ) -> None:
        """Test empty FCM token validation."""
        base_config["fcmToken"] = ""
        test_context.register_cleanup(base_config["deviceId"])

        response = make_request("POST", "config", base_config, test_context)
        assert_status(response, 400, test_context)

    def test_day_preferences_too_few(
        self, test_context: TestContext, base_config: Dict
    ) -> None:
        """Test day_preferences validation (too few elements)."""
        base_config["day_preferences"] = [False, False, False]
        test_context.register_cleanup(base_config["deviceId"])

        response = make_request("POST", "config", base_config, test_context)
        assert_status(response, 400, test_context)

    def test_day_preferences_too_many(
        self, test_context: TestContext, base_config: Dict
    ) -> None:
        """Test day_preferences validation (too many elements)."""
        base_config["day_preferences"] = [False] * 10
        test_context.register_cleanup(base_config["deviceId"])

        response = make_request("POST", "config", base_config, test_context)
        assert_status(response, 400, test_context)

    def test_invalid_time_zone(
        self, test_context: TestContext, base_config: Dict
    ) -> None:
        """Test invalid time zone validation."""
        base_config["time_zone"] = "Not/A_Real_Timezone"
        test_context.register_cleanup(base_config["deviceId"])

        response = make_request("POST", "config", base_config, test_context)
        assert_status(response, 400, test_context)


class TestUpdateOperations:
    """Tests for update/upsert operations."""

    def test_update_existing_config(self, test_context: TestContext) -> None:
        """Update an existing config (upsert)."""
        device_id = generate_uuid()
        test_context.register_cleanup(device_id)

        # First create
        make_request(
            "POST",
            "config",
            {
                "deviceId": device_id,
                "lat": 34.0522,
                "long": -118.2437,
                "fcmToken": "dKj8X9nR3qP5vL2mAPA91btest1234567890",
                "day_preferences": [False, False, False, True, True, True, True],
                "time_zone": "America/Los_Angeles",
                "enabled": True,
            },
            test_context,
        )

        # Update with new values
        response = make_request(
            "POST",
            "config",
            {
                "deviceId": device_id,
                "lat": 40.7128,
                "long": -74.0060,
                "fcmToken": "updatedtoken9876543210",
                "day_preferences": [True, True, True, True, True, False, False],
                "time_zone": "America/New_York",
                "enabled": False,
            },
            test_context,
        )

        assert_status(response, 200, test_context)

        # Verify update
        response = make_request("GET", f"config/{device_id}", context=test_context)

        assert_status(response, 200, test_context)
        assert_json_field(response, "time_zone", "America/New_York", test_context)
        assert_json_field(response, "enabled", False, test_context)
        assert_json_field(response, "lat", 40.7128, test_context)


class TestManualOverride:
    """Tests for manual alarm time override feature."""

    def test_manual_override_enabled_with_time(self, test_context: TestContext) -> None:
        """Create config with manual nextAlarmTime override."""
        device_id = generate_uuid()
        test_context.register_cleanup(device_id)

        alarm_time = get_future_time_utc()
        alarm_date = get_today_date()

        response = make_request(
            "POST",
            "config",
            {
                "deviceId": device_id,
                "lat": 34.0522,
                "long": -118.2437,
                "fcmToken": "dKj8X9nR3qP5vL2mAPA91btest1234567890",
                "day_preferences": [False, False, False, True, True, True, True],
                "time_zone": "America/Los_Angeles",
                "enabled": True,
                "nextAlarmTime": alarm_time,
                "alarmDateBucket": alarm_date,
            },
            test_context,
        )

        assert_status(response, 200, test_context)

        # Verify override fields are stored
        response = make_request("GET", f"config/{device_id}", context=test_context)

        assert_status(response, 200, test_context)
        # Note: These assertions may fail due to DynamoDB zero-value handling
        # Leaving them as-is to match original test behavior
        assert_json_field(response, "nextAlarmTime", alarm_time, test_context)
        assert_json_field(response, "alarmDateBucket", alarm_date, test_context)

    def test_manual_override_disabled_clears_time(
        self, test_context: TestContext
    ) -> None:
        """Update enabled config with manual time to disabled."""
        device_id = generate_uuid()
        test_context.register_cleanup(device_id)

        # First create enabled with manual time
        make_request(
            "POST",
            "config",
            {
                "deviceId": device_id,
                "lat": 34.0522,
                "long": -118.2437,
                "fcmToken": "dKj8X9nR3qP5vL2mAPA91btest1234567890",
                "day_preferences": [False, False, False, True, True, True, True],
                "time_zone": "America/Los_Angeles",
                "enabled": True,
                "nextAlarmTime": get_future_time_utc(),
                "alarmDateBucket": get_today_date(),
            },
            test_context,
        )

        # Update to disabled
        response = make_request(
            "POST",
            "config",
            {
                "deviceId": device_id,
                "lat": 34.0522,
                "long": -118.2437,
                "fcmToken": "dKj8X9nR3qP5vL2mAPA91btest1234567890",
                "day_preferences": [False, False, False, True, True, True, True],
                "time_zone": "America/Los_Angeles",
                "enabled": False,
            },
            test_context,
        )

        assert_status(response, 200, test_context)

        # Verify disabled config clears alarm time
        response = make_request("GET", f"config/{device_id}", context=test_context)

        assert_status(response, 200, test_context)
        # Note: These assertions may fail due to DynamoDB zero-value handling
        # Leaving them as-is to match original test behavior
        assert_json_field(response, "alarmDateBucket", "DISABLED", test_context)
        assert_json_field(response, "nextAlarmTime", "", test_context)


class TestCheckAlarm:
    """Tests for the check_alarm endpoint."""

    def test_check_alarm_endpoint(self, test_context: TestContext) -> None:
        """Test basic check_alarm endpoint."""
        response = make_request("GET", "check_alarm", context=test_context)

        # Note: This may return 404 if endpoint not deployed
        # Original test expects 200, keeping that expectation
        assert_status(response, 200, test_context)
        assert_json_contains(response, "fired", test_context)

    def test_check_alarm_with_manual_override(self, test_context: TestContext) -> None:
        """Test check_alarm with manually overridden past alarm time."""
        device_id = generate_uuid()
        test_context.register_cleanup(device_id)

        # Create config with alarm time in the past (should have fired)
        window_time = get_time_within_window()
        today_date = get_today_date()

        response = make_request(
            "POST",
            "config",
            {
                "deviceId": device_id,
                "lat": 34.0522,
                "long": -118.2437,
                "fcmToken": "dKj8X9nR3qP5vL2mAPA91btest1234567890",
                "day_preferences": [False, False, False, True, True, True, True],
                "time_zone": "America/Los_Angeles",
                "enabled": True,
                "nextAlarmTime": window_time,
                "alarmDateBucket": today_date,
            },
            test_context,
        )

        assert_status(response, 200, test_context)
        time.sleep(2)  # Allow eventual consistency for GSI

        # Check alarm (FCM will fail in tests, so fired will be 0)
        response = make_request("GET", "check_alarm", context=test_context)

        # Note: This may return 404 if endpoint not deployed
        # Original test expects 200, keeping that expectation
        assert_status(response, 200, test_context)
        assert_json_contains(response, "fired", test_context)

        # Verify the alarm time was recalculated
        response = make_request("GET", f"config/{device_id}", context=test_context)

        assert_status(response, 200, test_context)

        # Next alarm time should not be the window time we set
        data = response.json()
        current_next_alarm = data.get("nextAlarmTime")

        if current_next_alarm == window_time:
            raise AssertionError(
                "Alarm time was not recalculated",
                expected=f"Alarm time to be recalculated (not equal to {window_time})",
                actual=f"Alarm time is still {current_next_alarm}",
                response=data,
                request=test_context.last_request,
            )


# ==========================================
# Main entry point for running tests directly
# ==========================================

if __name__ == "__main__":
    # Check prerequisites (imports will fail here if missing)
    if not BOTO3_AVAILABLE:
        print("WARNING: 'boto3' not available. AWS cleanup will use CLI or be skipped.")

    # Run pytest programmatically
    print(f"Running integration tests against: {BASE_URL}")
    print(f"DynamoDB Table: {DYNAMODB_TABLE}")
    print(f"AWS Region: {AWS_REGION}")
    print()

    # Run tests and exit with appropriate code
    exit_code = pytest.main(
        [
            __file__,
            "-v",  # Verbose output
            "--tb=short",  # Short traceback format
        ]
    )

    sys.exit(exit_code)
