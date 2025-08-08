resource "aws_apigatewayv2_api" "http_api" {
  name          = "${var.env}-cloud-architecture-golang-api"
  protocol_type = "HTTP"
  description   = "API Gateway para Lambda Go - ${var.env}"
  
  cors_configuration {
    allow_credentials = false
    allow_headers     = ["content-type", "x-amz-date", "authorization", "x-api-key", "x-amz-security-token", "x-amz-user-agent"]
    allow_methods     = ["*"]
    allow_origins     = ["*"]
    expose_headers    = ["date", "keep-alive"]
    max_age          = 86400
  }
  
  tags = {
    Name        = "${var.env}-cloud-architecture-golang-api"
    Environment = var.env
    Project     = "cloud-architecture-golang"
  }
}

resource "aws_apigatewayv2_integration" "lambda_integration" {
  api_id                  = aws_apigatewayv2_api.http_api.id
  integration_type        = "AWS_PROXY"
  integration_uri         = var.lambda_invoke_arn
  payload_format_version  = "2.0"
  integration_method      = "POST"
}

resource "aws_apigatewayv2_route" "health_route" {
  api_id    = aws_apigatewayv2_api.http_api.id
  route_key = "GET /health"
  target    = "integrations/${aws_apigatewayv2_integration.lambda_integration.id}"
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

resource "aws_apigatewayv2_stage" "default_stage" {
  api_id      = aws_apigatewayv2_api.http_api.id
  name        = "$default"
  auto_deploy = true
  
  tags = {
    Name        = "${var.env}-default-stage"
    Environment = var.env
  }
}

# Domínio customizado (opcional)
resource "aws_apigatewayv2_domain_name" "custom_domain" {
  count       = var.custom_domain_name != "" ? 1 : 0
  domain_name = var.custom_domain_name

  domain_name_configuration {
    certificate_arn = var.certificate_arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }

  tags = {
    Name        = "${var.env}-custom-domain"
    Environment = var.env
  }
}

# Mapeamento do domínio customizado (opcional)
resource "aws_apigatewayv2_api_mapping" "custom_domain_mapping" {
  count       = var.custom_domain_name != "" ? 1 : 0
  api_id      = aws_apigatewayv2_api.http_api.id
  domain_name = aws_apigatewayv2_domain_name.custom_domain[0].id
  stage       = aws_apigatewayv2_stage.default_stage.id
}

resource "aws_lambda_permission" "api_invoke" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.http_api.execution_arn}/*/*"
}