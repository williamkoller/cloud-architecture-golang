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

variable "image_tag" {
  type    = string
  default = "latest"
}

variable "alert_email" {
  type    = string
  default = ""
}

variable "alert_topic_arn" {
  type    = string
  default = ""
}
