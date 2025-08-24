variable "alias" {
  description = "Alias para o workspace Prometheus"
  type        = string
}

variable "environment" {
  description = "Ambiente (staging, production, etc)"
  type        = string
  default     = "staging"
}