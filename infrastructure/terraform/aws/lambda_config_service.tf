# Create execution role for lambda to run
resource "aws_iam_role" "lambda_assume_role" {
  name = "config_service"

  assume_role_policy = jsonencode({
  Version       = "2012-10-17"
    Statement   = [{
      Effect    = "Allow"
      Principal = { Service = "lambda.amazonaws.com" }
      Action    = "sts:AssumeRole"
    }]
  })
}

# Permissions for lambda to access dynamoDB, and create CloudWatch logs
resource "aws_iam_role_policy" "config_service" {
  role = aws_iam_role.lambda_assume_role.name
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:Query"
        ]
        Resource = [
          aws_dynamodb_table.user_configs.arn,
          "${aws_dynamodb_table.user_configs.arn}/index/AlarmTimeIndex"
        ]
      },
      {
        Effect   = "Allow"
        Action   = "secretsmanager:GetSecretValue"
        Resource = aws_secretsmanager_secret.lambda_service_key.arn
      }
    ]
  })
}

# Needed to push to ECR and put cloudwatch logs
resource "aws_iam_role_policy_attachment" "lambda_ecr" {
  role       = aws_iam_role.lambda_assume_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_lambda_function" "config_service" {
  function_name = "config_service"
  role          = aws_iam_role.lambda_assume_role.arn
  package_type  = "Image"
  image_uri     = "${aws_ecr_repository.sunrise_repo.repository_url}:${local.config_service_tag}"
  timeout       = 30
  memory_size   = 256

  environment {
    variables = {
      PORT = "8080"
      AWS_LWA_PASS_THROUGH_PATH = "/check_alarm" # Pass eventbridge traffic to this endpoint
      AWS_LWA_READINESS_CHECK_PATH = "/health"
    }
  }

  depends_on = [null_resource.docker_build_push_config_service]
}

# service URL
resource "aws_lambda_function_url" "config_service" {
  function_name      = aws_lambda_function.config_service.function_name
  authorization_type = "NONE"
}

# Allow public access to the function
resource "aws_lambda_permission" "function_url_public_access" {
  statement_id           = "FunctionURLAllowPublicAccess2"
  action                 = "lambda:InvokeFunctionUrl"
  function_name          = aws_lambda_function.config_service.function_name
  principal              = "*"
  function_url_auth_type = "NONE"
}

# Allow public invoking to the function
resource "aws_lambda_permission" "function_url_invoke_access" {
  statement_id  = "FunctionURLInvokeAllowPublicAccess"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.config_service.function_name
  principal     = "*"
}

output "config_service_url" {
  value = aws_lambda_function_url.config_service.function_url
}
