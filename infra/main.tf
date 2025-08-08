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

# ===== IAM =====
resource "aws_iam_role" "lambda_exec" {
  name = "${var.env}-lambda-exec-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Effect = "Allow",
      Principal = { Service = "lambda.amazonaws.com" },
      Action    = "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_basic_logs" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"

  depends_on = [ aws_iam_role.lambda_exec ]
}

# ===== ECR =====
resource "aws_ecr_repository" "lambda_repo" {
  name                 = "${var.env}-lambda-go"
  image_tag_mutability = "MUTABLE"
  image_scanning_configuration { scan_on_push = true }
}

# Dica: a imagem deve existir no ECR antes de criar a Lambda (docker build/tag/push)

# ===== Lambda (container) =====
resource "aws_lambda_function" "golang_lambda" {
  function_name = "${var.env}-golang-api"
  role          = aws_iam_role.lambda_exec.arn

  package_type = "Image"
  image_uri    = "${aws_ecr_repository.lambda_repo.repository_url}:${var.image_tag}"

  timeout     = 10
  memory_size = 128
  publish     = true

  depends_on  = [aws_iam_role_policy_attachment.lambda_basic_logs]
}

# ===== API Gateway HTTP =====
resource "aws_apigatewayv2_api" "http_api" {
  name          = "${var.env}-api"
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_integration" "lambda_integration" {
  api_id                  = aws_apigatewayv2_api.http_api.id
  integration_type        = "AWS_PROXY"
  integration_uri         = aws_lambda_function.golang_lambda.invoke_arn
  payload_format_version  = "2.0"
  integration_method      = "POST"
}

resource "aws_apigatewayv2_route" "health_route" {
  api_id    = aws_apigatewayv2_api.http_api.id
  route_key = "GET /health"
  target    = "integrations/${aws_apigatewayv2_integration.lambda_integration.id}"
}

resource "aws_apigatewayv2_stage" "default_stage" {
  api_id      = aws_apigatewayv2_api.http_api.id
  name        = "$default"
  auto_deploy = true
}

resource "aws_lambda_permission" "api_invoke" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.golang_lambda.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.http_api.execution_arn}/*/*"
}

# ===== CloudWatch =====
resource "aws_cloudwatch_log_group" "lambda_logs" {
  name              = "/aws/lambda/${aws_lambda_function.golang_lambda.function_name}"
  retention_in_days = 14
}

resource "aws_cloudwatch_metric_alarm" "lambda_errors" {
  alarm_name          = "${var.env}-lambda-errors"
  metric_name         = "Errors"
  namespace           = "AWS/Lambda"
  statistic           = "Sum"
  period              = 60
  evaluation_periods  = 1
  threshold           = 1
  comparison_operator = "GreaterThanThreshold"
  alarm_description   = "Erros na Lambda"
  dimensions = { FunctionName = aws_lambda_function.golang_lambda.function_name }
  alarm_actions = var.alert_topic_arn == "" ? [] : [var.alert_topic_arn]
}

# ===== SNS (opcional) =====
resource "aws_sns_topic" "alerts" {
  count = var.alert_email == "" ? 0 : 1
  name  = "${var.env}-lambda-alerts"
}

resource "aws_sns_topic_subscription" "email" {
  count     = var.alert_email == "" ? 0 : 1
  topic_arn = aws_sns_topic.alerts[0].arn
  protocol  = "email"
  endpoint  = var.alert_email
}

# ===== Route53 Health Check =====
resource "aws_route53_health_check" "lambda_health" {
  fqdn              = split("/", replace(aws_apigatewayv2_api.http_api.api_endpoint, "https://", ""))[0]
  port              = 443
  type              = "HTTPS"
  resource_path     = "/health"
  failure_threshold = 3
  request_interval  = 30
  regions           = ["us-east-1", "us-west-2", "eu-west-1"]
}

output "api_gateway_url" { value = aws_apigatewayv2_api.http_api.api_endpoint }
output "lambda_image_uri" { value = aws_lambda_function.golang_lambda.image_uri }


resource "aws_lambda_alias" "staging" {
  name             = "staging"
  function_name    = aws_lambda_function.golang_lambda.function_name
  function_version = aws_lambda_function.golang_lambda.version
}

resource "aws_lambda_provisioned_concurrency_config" "pc" {
  count = var.provisioned_concurrency > 0 ? 1 : 0

  function_name                     = aws_lambda_function.golang_lambda.function_name
  qualifier                         = aws_lambda_alias.staging.name
  provisioned_concurrent_executions = var.provisioned_concurrency
}

resource "aws_apigatewayv2_route" "any_root" {
  api_id    = aws_apigatewayv2_api.http_api.id
  route_key = "ANY /"
  target    = "integrations/${aws_apigatewayv2_integration.lambda_integration.id}"
}

resource "aws_apigatewayv2_route" "any_proxy" {
  api_id    = aws_apigatewayv2_api.http_api.id
  route_key = "ANY /{proxy+}"
  target    = "integrations/${aws_apigatewayv2_integration.lambda_integration.id}"
}