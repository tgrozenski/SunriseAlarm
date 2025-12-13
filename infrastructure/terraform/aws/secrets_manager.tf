resource "aws_secretsmanager_secret" "firebase_secrets" {
  name = "firebase_secrets"
}

output "secret_manager_arn" {
  description = "secret manager's arn"
  value = aws_secretsmanager_secret.firebase_secrets.arn
}