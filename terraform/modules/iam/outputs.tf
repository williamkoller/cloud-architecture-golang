output "lambda_execution_role_arn" {
  value = aws_iam_role.lambda_exec.arn
}

output "lambda_policy_attachment_name" {
  value = aws_iam_role_policy_attachment.lambda_basic_logs.id
}

output "api_gateway_logs_role_arn" {
  value       = aws_iam_role.api_gateway_logs.arn
  description = "ARN da role para logs da API Gateway"
}
