variable "env" {
  type        = string
  description = "The environment name (e.g., dev, staging, prod)."
}

variable "aws_region" {
  type        = string
  default     = "us-east-1"
  description = "The AWS region to deploy resources in."
}

variable "aws_profile" {
  type        = string
  default     = null
  description = "The AWS profile to use for authentication."
}

variable "account_id" {
  type        = string
  description = "The AWS account ID where the resources will be deployed."
}

variable "image_tag" {
  type        = string
  default     = "latest"
  description = "The Docker image tag to be deployed."
  validation {
    condition     = length(var.image_tag) > 0
    error_message = "image_tag não pode ser vazio."
  }
}

variable "alert_email" {
  type        = string
  default     = ""
  description = "The email address for SNS topic subscription for alerts. If provided, a new SNS topic and subscription will be created."
}

variable "alert_topic_arn" {
  type        = string
  default     = ""
  description = "The ARN of an existing SNS topic to send alerts to. If provided, this will be used instead of creating a new one."
}


