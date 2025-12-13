# Role for Lambda
resource "aws_iam_role" "lambda_assume_role" {
  name = "firebase_caller"

  assume_role_policy = jsonencode({
    Version     = "2012-10-17"
    Statement   = [{
      Effect    = "Allow"
      Principal = { Service = "lambda.amazonaws.com" }
      Action    = "sts:AssumeRole"
    }]
  })
}

# Lambda general policy, allows creating logs, reading secrets, and publishing to SNS
resource "aws_iam_role_policy" "lambda_general" {
  role         = aws_iam_role.lambda_assume_role.arn
  policy       = jsonencode({
    Version    = "2012-10-17"
    Statement  = [{
      Effect   = "Allow"
      Action   = ["secretsmanager:GetSecretValue"]
      Resource = "*"
    },
    {
      "Effect" : "Allow"
      "Action" : [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource" : "*"
    },
    {
      "Sid":"AllowPublishToMyTopic",
      "Effect":"Allow",
      "Action":"sns:Publish",
      "Resource":aws_sns_topic.observability.arn
    }]
  })
}

# Package the Lambda function code
data "archive_file" "lambda_src" {
  type        = "zip"
  source_dir  = "src"
  output_path = "lambda/function.zip"
}

# Lambda function
resource "aws_lambda_function" "firebase_caller" {
  filename         = data.archive_file.lambda_src.output_path
  function_name    = "firebase_caller"
  role             = aws_iam_role.lambda_assume_role.arn
  handler          = "main.handler"
  source_code_hash = data.archive_file.lambda_src.output_base64sha256

  runtime = "python3.14"
}

output "lambda_arn" {
  description = "The Lambda's ARN"
  value = aws_lambda_function.firebase_caller.arn
}