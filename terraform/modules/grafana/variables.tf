variable "name" {
  description = "Nome do workspace Grafana"
  type        = string
}

variable "iam_role_arn" {
  description = "ARN da role IAM para o Grafana"
  type        = string
}

variable "prometheus_endpoint" {
  description = "Endpoint do workspace Prometheus"
  type        = string
  default     = ""
}

variable "environment" {
  description = "Ambiente (staging, production, etc)"
  type        = string
  default     = "staging"
}