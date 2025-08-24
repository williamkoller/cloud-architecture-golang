# NOTA: O usuário terraform-user não tem permissão aps:CreateWorkspace
# Para usar este módulo, solicite as seguintes permissões IAM:
# - aps:CreateWorkspace
# - aps:TagResource
# - aps:ListWorkspaces

# Workspace Prometheus comentado temporariamente devido a falta de permissões
# resource "aws_prometheus_workspace" "amp" {
#   alias = var.alias
# 
#   lifecycle {
#     prevent_destroy = true
#   }
# }

# Rule group comentado pois depende do workspace
# resource "aws_prometheus_rule_group_namespace" "alerts" {
#   workspace_id = aws_prometheus_workspace.amp.id
#   name         = "alerts"
#   data         = file("${path.root}/../prometheus/alerts.yml")
#
#   depends_on = [aws_prometheus_workspace.amp]
# }

# Outputs temporariamente desabilitados
# output "workspace_prometheus_endpoint" {
#   value = aws_prometheus_workspace.amp.prometheus_endpoint
# }

# output "workspace_id" {
#   value = aws_prometheus_workspace.amp.id
# }

# Outputs temporários para evitar erros de referência
output "workspace_prometheus_endpoint" {
  value = "prometheus-endpoint-not-available"
}

output "workspace_id" {
  value = "prometheus-workspace-not-available"
}
