variable "env" {
  type        = string
  description = "Environment name"
}

variable "lambda_invoke_arn" {
  type        = string
  description = "Lambda function invoke ARN"
}

variable "lambda_function_name" {
  type        = string
  description = "Lambda function name"
}

variable "custom_domain_name" {
  type        = string
  description = "Nome do domínio customizado (opcional)"
  default     = ""
}

variable "certificate_arn" {
  type        = string
  description = "ARN do certificado SSL para o domínio customizado (opcional)"
  default     = ""
}