# ===== CloudWatch =====
resource "aws_cloudwatch_log_group" "lambda_logs" {
  name              = "/aws/lambda/${aws_lambda_function.golang_lambda.function_name}"
  retention_in_days = 14
}

resource "aws_cloudwatch_metric_alarm" "lambda_errors" {
  alarm_name          = "${var.env}-lambda-errors"
  metric_name         = "Errors"
  namespace           = "AWS/Lambda"
  statistic           = "Sum"
  period              = 60
  evaluation_periods  = 1
  threshold           = 1
  comparison_operator = "GreaterThanThreshold"
  alarm_description   = "Erros na Lambda"
  dimensions = { FunctionName = aws_lambda_function.golang_lambda.function_name }
  alarm_actions = var.alert_topic_arn == "" ? [] : [var.alert_topic_arn]
}
