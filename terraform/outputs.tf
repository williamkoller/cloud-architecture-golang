#The name of the Lambda function.
output "lambda_function_name" {
  value       = aws_lambda_function.golang_lambda.function_name
}
output "api_gateway_endpoint" {
  description = "The base URL for the API Gateway."
  value       = aws_apigatewayv2_api.http_api.api_endpoint
}

output "ecr_repository_url" {
  description = "The URL of the ECR repository."
  value       = aws_ecr_repository.lambda_repo.repository_url
}

output "lambda_iam_role_arn" {
  description = "The ARN of the IAM role used by the Lambda function."
  value       = aws_iam_role.lambda_exec.arn
}

output "sns_topic_arn" {
  description = "The ARN of the SNS topic for alerts. This will be null if no alert_email is provided."
  value       = one(aws_sns_topic.alerts[*.arn)
}