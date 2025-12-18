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
  role         = aws_iam_role.lambda_assume_role.name
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
    },
    {
    Effect   = "Allow"
    Action   = [
      "ecr:GetDownloadUrlForLayer",
      "ecr:BatchGetImage"
    ]
      Resource = aws_ecr_repository.firebase_caller.arn
    },
    {
      Effect   = "Allow"
      Action   = "ecr:GetAuthorizationToken"
      Resource = "*"
    }]
  })
}

# Ensure Lambda is allowed to pull from ECR
resource "aws_iam_role_policy_attachment" "lambda_ecr" {
  role       = aws_iam_role.lambda_assume_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# ECR Repository to store lambda container image
resource "aws_ecr_repository" "firebase_caller" {
  name         = "firebase-caller"
  force_delete = true
}

locals {
  image_tag = md5(join("", [
    filemd5("${path.module}/lambda/requirements.txt"),
    filemd5("${path.module}/lambda/Dockerfile"),
    filemd5("${path.module}/lambda/src/main.py")
  ]))
}

resource "null_resource" "docker_build_push" {
  triggers = {
    image_tag = local.image_tag
  }

  provisioner "local-exec" {
    command = <<EOF
      aws ecr get-login-password --region us-west-1 | docker login --username AWS --password-stdin ${aws_ecr_repository.firebase_caller.repository_url}
      docker build -t ${aws_ecr_repository.firebase_caller.repository_url}:${local.image_tag} ${path.module}/lambda
      docker push ${aws_ecr_repository.firebase_caller.repository_url}:${local.image_tag}
    EOF
  }

  depends_on = [aws_ecr_repository.firebase_caller]
}

resource "aws_lambda_function" "firebase_caller" {
  function_name = "firebase_caller"
  role          = aws_iam_role.lambda_assume_role.arn
  package_type  = "Image"
  image_uri     = "${aws_ecr_repository.firebase_caller.repository_url}:${local.image_tag}"
  timeout       = 30
  memory_size   = 256

  depends_on = [null_resource.docker_build_push]
}

output "lambda_arn" {
  description = "The Lambda's ARN"
  value = aws_lambda_function.firebase_caller.arn
}
