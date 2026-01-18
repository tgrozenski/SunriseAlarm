resource "aws_dynamodb_table" "user_configs" {
  name         = "UserConfigs"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "deviceId"

  attribute {
    name = "deviceId"
    type = "S"
  }

  attribute {
    name = "alarmDateBucket"
    type = "S"
  }

  attribute {
    name = "nextAlarmTime"
    type = "S"
  }

  global_secondary_index {
    name            = "AlarmTimeIndex"
    hash_key        = "alarmDateBucket"
    range_key       = "nextAlarmTime"
    projection_type = "ALL"
  }
}
