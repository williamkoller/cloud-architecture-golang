variable "alert_topic_arn" {
  type    = string
  default = ""
}

variable "env" {
  type = string
}

variable "lambda_function_name" {
  type        = string
  description = "Lambda function name"
}

variable "sns_topic_arn" {
  type        = string
  description = "SNS topic ARN for alerts"
}