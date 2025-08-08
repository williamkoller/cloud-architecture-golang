variable "env" {
  type = string
}

variable "image_tag" {
  type    = string
  default = "latest"
}

variable "lambda_execution_role_arn" {
  type        = string
  description = "ARN da role de execução do Lambda"
}

variable "ecr_repository_url" {
  type        = string
  description = "URL do repositório ECR"
}

variable "iam_policy_attachment_name" {
  type        = string
  description = "Nome do attachment da política IAM para dependência"
}

variable "provisioned_concurrency" {
  type        = number
  description = "Número de execuções concorrentes provisionadas"
  default     = 0
}