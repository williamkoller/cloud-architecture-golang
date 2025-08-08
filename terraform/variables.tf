variable "env" {
  type = string
}

variable "aws_region" {
  type    = string
  default = "us-east-1"
}

variable "aws_profile" {
  type    = string
  default = null
}

variable "account_id" {
  type = string
}

variable "alert_email" {
  type    = string
  default = "041codekoller@gmail.com"
}

variable "alert_topic_arn" {
  type    = string
  default = ""
}

variable "provisioned_concurrency" {
  type    = number
  default = 0
}

variable "image_tag" {
  type        = string
  description = "Tag da imagem no ECR"
  default     = "latest"

  validation {
    condition     = length(var.image_tag) > 0
    error_message = "image_tag não pode ser vazio."
  }
}

# Variáveis para domínio customizado (opcional)
variable "custom_domain_name" {
  type        = string
  description = "Nome do domínio customizado para a API (ex: api.seudominio.com)"
  default     = "cloud-architecture-golang.com"
}

variable "certificate_arn" {
  type        = string
  description = "ARN do certificado SSL/TLS para o domínio customizado"
  default     = ""
}