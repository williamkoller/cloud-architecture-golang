resource "aws_grafana_workspace" "grafana" {
  name                      = var.name
  account_access_type       = "CURRENT_ACCOUNT"
  authentication_providers = ["SAML"]
  role_arn                  = var.iam_role_arn
  data_sources              = ["CLOUDWATCH", "PROMETHEUS"]
  permission_type           = "SERVICE_MANAGED"

  tags = {
    Name        = var.name
    Environment = var.environment
    Project     = "cloud-architecture-golang"
  }

  lifecycle {
    prevent_destroy = true
    ignore_changes = [
      name
    ]
  }
}

output "grafana_endpoint" {
  value = aws_grafana_workspace.grafana.endpoint
}

output "grafana_id" {
  value = aws_grafana_workspace.grafana.id
}
