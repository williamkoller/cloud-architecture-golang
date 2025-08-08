output "lambda_function_arn" {
  value = aws_lambda_function.golang_lambda.arn
}

output "lambda_function_name" {
  value = aws_lambda_function.golang_lambda.function_name
}

output "lambda_invoke_arn" {
  value = aws_lambda_function.golang_lambda.invoke_arn
}

output "lambda_alias_arn" {
  value = aws_lambda_alias.staging.arn
}