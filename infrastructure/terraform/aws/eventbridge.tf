# EventBridge rule to poll check_alarm endpoint every minute
resource "aws_cloudwatch_event_rule" "check_alarm_cron" {
  name                = "check-alarm-cron"
  description         = "Triggers every minute to poll the check_alarm endpoint"
  schedule_expression = "cron(* 8-16 * * ? *)"
  state               = "DISABLED"  # Use "ENABLED" for prod
  # schedule_expression = "rate(1 minute)"
  # state = "ENABLED"
}

resource "aws_cloudwatch_event_target" "check_alarm_target" {
  rule = aws_cloudwatch_event_rule.check_alarm_cron.name
  arn  = aws_lambda_function.config_service.arn

  input_transformer {
    input_template = jsonencode({
      httpMethod = "GET"
      path       = "/check_alarm"
      headers    = {}
      body       = ""
    })
  }
}

# Allow EventBride to invoke lambda
resource "aws_lambda_permission" "allow_eventbridge" {
  statement_id  = "AllowExecutionFromEventBridge"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.config_service.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.check_alarm_cron.arn
}

