# Create execution role for lambda to run
resource "aws_iam_role" "lambda_assume_role" {
  name = "config_service"

  assume_role_policy = jsonencode({
    Version     = "2012-10-17"
    Statement   = [{
      Effect    = "Allow"
      Principal = { Service = "lambda.amazonaws.com" }
      Action    = "sts:AssumeRole"
    }]
  })
}

# Permissions for lambda to access dynamoDB
resource "aws_iam_role_policy" "config_service" {
  role         = aws_iam_role.lambda_assume_role.name
  policy       = jsonencode({
    Version    = "2012-10-17"
    Statement  = [
      {
        Effect = "Allow"
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem"
        ]
        Resource = aws_dynamodb_table.user_configs.arn
      }
    ]
  })
}

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

  layers = [
    "arn:aws:lambda:${var.aws_region}:753240598075:layer:LambdaAdapterLayerX86:20"
  ]

  environment {
    variables = {
      PORT    = "8080"
    }
  }

  depends_on = [null_resource.docker_build_push_config_service]
}
