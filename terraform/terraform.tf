terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0"
    }
  }

  required_version = ">= 1.5.0"

  # in case you want to create a state in s3
  #backend "s3" {
    
  #}
}

provider "aws" {
  region  = var.aws_region
  profile = var.aws_profile # opcional, defina em variables.tf/ tfvars
}