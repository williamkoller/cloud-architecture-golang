terraform {
  required_version = ">= 1.5.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0"
    }
  }
}

provider "aws" {
  region  = var.aws_region
  profile = var.aws_profile
}

module "apigw" {
  source               = "./modules/apigw"
  env                  = var.env
  lambda_invoke_arn    = module.lambda.lambda_invoke_arn
  lambda_function_name = module.lambda.lambda_function_name
  depends_on = [module.lambda]
}

module "cloudwatch" {
  source               = "./modules/cloudwatch"
  env                  = var.env
  lambda_function_name = module.lambda.lambda_function_name
  sns_topic_arn        = module.sns.topic_arn
  depends_on           = [module.lambda, module.sns]
}

module "ecr" {
  source = "./modules/ecr"
  env    = var.env
}

module "iam" {
  source = "./modules/iam"
  env    = var.env
}

module "lambda" {
  source    = "./modules/lambda"
  env       = var.env
  image_tag = var.image_tag

  lambda_execution_role_arn  = module.iam.lambda_execution_role_arn
  ecr_repository_url         = module.ecr.repository_url
  iam_policy_attachment_name = module.iam.lambda_policy_attachment_name
  provisioned_concurrency    = var.provisioned_concurrency

  depends_on = [module.iam, module.ecr]
}

module "route53" {
  source = "./modules/route53"
  env    = var.env

  api_endpoint = module.apigw.api_endpoint
  depends_on   = [module.apigw]
  domain = var.custom_domain_name
}

module "sns" {
  source      = "./modules/sns"
  env         = var.env
  alert_email = var.alert_email
}

output "api_gateway_url" {
  value       = module.apigw.api_endpoint
  description = "URL da API Gateway"
}

output "api_gateway_id" {
  value       = module.apigw.api_id
  description = "ID da API Gateway"
}

output "lambda_function_name" {
  value       = module.lambda.lambda_function_name
  description = "Nome da função Lambda"
}

output "lambda_function_arn" {
  value       = module.lambda.lambda_function_arn
  description = "ARN da função Lambda"
}

output "ecr_repository_url" {
  value       = module.ecr.repository_url
  description = "URL do repositório ECR"
}

output "custom_domain_url" {
  value       = module.apigw.custom_domain_url
  description = "URL do domínio customizado (se configurado)"
}

output "health_check_id" {
  value       = module.route53.health_check_id
  description = "ID do health check Route53"
}

output "sns_topic_arn" {
  value       = module.sns.topic_arn
  description = "ARN do tópico SNS para notificações"
}

# output "prometheus_workspace_endpoint" {
#   value       = module.prometheus.workspace_prometheus_endpoint
#   description = "Endpoint do workspace Prometheus"
# }

# output "prometheus_workspace_id" {
#   value       = module.prometheus.workspace_id
#   description = "ID do workspace Prometheus"
# }

# output "grafana_endpoint" {
#   value       = module.grafana.grafana_endpoint
#   description = "Endpoint do workspace Grafana"
# }

# output "grafana_workspace_id" {
#   value       = module.grafana.grafana_id
#   description = "ID do workspace Grafana"
# }