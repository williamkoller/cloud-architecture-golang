resource "aws_grafana_workspace" "grafana" {
  name                      = var.name
  account_access_type       = "CURRENT_ACCOUNT"
  authentication_providers = ["SAML"]
  role_arn                  = var.iam_role_arn
  data_sources              = ["CLOUDWATCH", "PROMETHEUS"]
  permission_type           = "SERVICE_MANAGED"
  
  # Configuração simplificada compatível com Grafana 10.4
  # A configuração de datasources será feita após a criação via API/UI
  # configuration = jsonencode({
  #   datasources = {
  #     datasources = [
  #       {
  #         name   = "CloudWatch"
  #         type   = "cloudwatch"
  #         access = "proxy"
  #       },
  #       {
  #         name   = "Prometheus"
  #         type   = "prometheus"
  #         access = "proxy"
  #         url    = var.prometheus_endpoint
  #       }
  #     ]
  #   }
  # })

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

# Output para acesso
output "grafana_endpoint" {
  value = aws_grafana_workspace.grafana.endpoint
}

output "grafana_id" {
  value = aws_grafana_workspace.grafana.id
}
