output "api_endpoint" {
  value       = aws_apigatewayv2_api.http_api.api_endpoint
  description = "URL padrão da API Gateway"
}

output "api_id" {
  value       = aws_apigatewayv2_api.http_api.id
  description = "ID da API Gateway"
}

output "api_execution_arn" {
  value       = aws_apigatewayv2_api.http_api.execution_arn
  description = "ARN de execução da API Gateway"
}

output "custom_domain_url" {
  value       = var.custom_domain_name != "" ? "https://${var.custom_domain_name}" : ""
  description = "URL do domínio customizado (se configurado)"
}

output "custom_domain_target" {
  value       = var.custom_domain_name != "" ? aws_apigatewayv2_domain_name.custom_domain[0].domain_name_configuration[0].target_domain_name : ""
  description = "Target domain para configuração do DNS (se domínio customizado configurado)"
}