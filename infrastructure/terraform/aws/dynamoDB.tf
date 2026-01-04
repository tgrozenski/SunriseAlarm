resource "aws_dynamodb_table" "user_configs" {
  name           = "UserConfigs"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "deviceId"

  attribute {
    name = "deviceId"
    type = "S"
  }
}
