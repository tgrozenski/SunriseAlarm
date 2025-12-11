module "eventbridge" {
  source = "terraform-aws-modules/eventbridge/aws"

  create_bus = false

  rules = {
    crons = {
      description         = "Trigger for a Lambda"
      schedule_expression = "cron(0 12 ? * SAT *)"
    }
  }

  targets = {
    crons = [
      {
        name  = "lambda-notification-sender"
        arn   = aws_lambda_function.firebase_caller.arn
        input = jsonencode({"job": "cron-by-rate"})
      }
    ]
  }
}

# Ensure eventbridge has permission to invoke lambda
resource "aws_lambda_permission" "allow_eventbridge" {
  statement_id  = "AllowEventBridge"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.firebase_caller.function_name
  principal     = "events.amazonaws.com"
}