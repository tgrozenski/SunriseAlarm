resource "aws_sns_topic" "observability" {
  name = "observability-topic"
}

# Subscribing myself to get updates
resource "aws_sns_topic_subscription" "get_observability_email" {
  topic_arn = aws_sns_topic.observability.arn
  protocol  = "email"
  endpoint  = "tyler.grozenski@gmail.com"
}

output "sns_topic_arn" {
  description = "The SNS topic's arn"
  value = aws_sns_topic.observability.arn
}