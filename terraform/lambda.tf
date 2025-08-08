# ===== Lambda (container) =====
resource "aws_lambda_function" "golang_lambda" {
  function_name = "${var.env}-golang-api"
  role          = aws_iam_role.lambda_exec.arn

  package_type = "Image"
  image_uri    = "${var.account_id}.dkr.ecr.${var.aws_region}.amazonaws.com/${aws_ecr_repository.lambda_repo.name}:${var.image_tag}"

  timeout     = 10
  memory_size = 128

  # Se usava kms_key_arn e não tem a chave, comente/remova:
  # kms_key_arn = aws_kms_key.lambda_env.arn
  depends_on = [aws_iam_role_policy_attachment.lambda_basic_logs]
}

resource "aws_lambda_permission" "api_invoke" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.golang_lambda.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.http_api.execution_arn}/*/*"
}