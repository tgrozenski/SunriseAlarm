terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
    }
  }
  cloud { 
    organization = "example-org-5e1658" 

    workspaces { 
      name = "sunrise-app" 
    } 
  }
}

variable "aws_access_key" {}
variable "aws_secret_key" {}

# Configure the AWS Provider
provider "aws" {
  region     = "us-west-1"
  access_key = var.aws_access_key
  secret_key = var.aws_secret_key
}

# TODO: write all lambda permissions here, attach once
data "aws_iam_policy_document" "assume_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "firebase_caller" {
  name               = "lambda_execution_role"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

# Attach policy to create CloudWatch logs
resource "aws_iam_role_policy_attachment" "lambda_logs" {
  role       = aws_iam_role.firebase_caller.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# Secrets Manager Read Only
resource "aws_iam_role_policy" "secrets_read" {
  role = aws_iam_role.firebase_caller.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect   = "Allow"
      Action   = ["secretsmanager:GetSecretValue"]
      Resource = "*"
    }]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_read_secrets" {
  role       = aws_iam_role.firebase_caller.name
  policy_arn = aws_iam_role_policy.secrets_read
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
  role             = aws_iam_role.firebase_caller.arn
  handler          = "main.handler"
  source_code_hash = data.archive_file.lambda_src.output_base64sha256

  runtime = "python3.14"
}