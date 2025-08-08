output "topic_arn" {
  value = length(aws_sns_topic.alerts) > 0 ? aws_sns_topic.alerts[0].arn : ""
}
