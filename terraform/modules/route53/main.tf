resource "aws_route53_health_check" "lambda_health" {
  fqdn              = split("/", replace(var.api_endpoint, "https://", ""))[0]
  port              = 443
  type              = "HTTPS"
  resource_path     = "/health"
  failure_threshold = 3
  request_interval  = 30
  regions           = ["us-east-1", "us-west-2", "eu-west-1"]
}