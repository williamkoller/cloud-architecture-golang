output "lambda_function_name" {
  value = aws_lambda_function.golang_lambda.function_name
}

output "lambda_alias_arn" {
  value = aws_lambda_alias.staging.arn
}
