resource "aws_lambda_function" "golang_lambda" {
  function_name = "${var.env}-golang-api"
  role          = var.lambda_execution_role_arn

  package_type = "Image"
  image_uri    = "${var.ecr_repository_url}:${var.image_tag}"

  timeout       = 10
  memory_size   = 256
  publish       = true
  architectures = ["x86_64"]

  tags = {
    Name        = "${var.env}-golang-api"
    Environment = var.env
    Project     = "cloud-architecture-golang"
  }

  lifecycle {
    # Evita conflitos quando a função já existe
    ignore_changes = [
      function_name,
      package_type
    ]
  }
}

# Alias para staging (opcional, pode ser removido se não necessário)
resource "aws_lambda_alias" "staging" {
  name             = "staging"
  description      = "Lambda alias for staging environment"
  function_name    = aws_lambda_function.golang_lambda.function_name
  function_version = aws_lambda_function.golang_lambda.version
}

# Concorrência provisionada (opcional)
resource "aws_lambda_provisioned_concurrency_config" "pc" {
  count = var.provisioned_concurrency > 0 ? 1 : 0

  function_name                     = aws_lambda_function.golang_lambda.function_name
  qualifier                         = aws_lambda_alias.staging.name
  provisioned_concurrent_executions = var.provisioned_concurrency
}
