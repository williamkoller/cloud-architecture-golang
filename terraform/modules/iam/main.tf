resource "aws_iam_role" "lambda_exec" {
  name = "${var.env}-lambda-exec-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Effect    = "Allow",
      Principal = { Service = "lambda.amazonaws.com" },
      Action    = "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_basic_logs" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"

  depends_on = [aws_iam_role.lambda_exec]
}

# Role para API Gateway logs
resource "aws_iam_role" "api_gateway_logs" {
  name = "${var.env}-api-gateway-logs-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "apigateway.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "api_gateway_logs" {
  role       = aws_iam_role.api_gateway_logs.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonAPIGatewayPushToCloudWatchLogs"
}


resource "aws_iam_role" "grafana_role" {
  name = "grafana-service-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Principal = {
        Service = "grafana.amazonaws.com"
      }
      Action = "sts:AssumeRole"
    }]
  })

  lifecycle {
    prevent_destroy = true
    ignore_changes  = [name]
  }
}

resource "aws_iam_role_policy_attachment" "grafana_policy_attach" {
  role       = aws_iam_role.grafana_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonGrafanaCloudWatchAccess"

  lifecycle {
    prevent_destroy = true
  }
}
