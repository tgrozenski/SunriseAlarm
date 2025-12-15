output "secret_manager_arn" {
  description = "secret manager's arn"
  value = aws_secretsmanager_secret.firebase_secrets.arn
}

# Define the secret
resource "aws_secretsmanager_secret" "lambda_service_key" {
  name = "LAMBDA_SERVICE_KEY"
}

# Create primary version of the secret
resource "aws_secretsmanager_secret_version" "v1" {
  secret_id     = aws_secretsmanager_secret.lambda_service_key.id
  secret_string = base64decode(data.terraform_remote_state.gcp.outputs.service_account_key)
}
