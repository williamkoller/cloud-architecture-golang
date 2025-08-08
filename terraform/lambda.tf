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

resource "aws_lambda_permission" "api_invoke" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.golang_lambda.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.http_api.execution_arn}/*/*"
}
