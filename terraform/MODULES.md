# 📦 Documentação dos Módulos Terraform

Este documento detalha cada módulo Terraform usado no projeto, suas responsabilidades, inputs, outputs e dependências.

## 🗂️ Estrutura dos Módulos

### 1. ECR Module (`modules/ecr/`)

**Responsabilidade**: Gerencia o repositório de containers no Amazon ECR.

**Recursos criados**:

- `aws_ecr_repository.lambda_repo` - Repositório para imagens Docker

**Inputs**:

```hcl
variable "env" {
  type        = string
  description = "Nome do ambiente (staging/production)"
}
```

**Outputs**:

```hcl
output "repository_url" {
  value = aws_ecr_repository.lambda_repo.repository_url
}

output "repository_name" {
  value = aws_ecr_repository.lambda_repo.name
}

output "repository_arn" {
  value = aws_ecr_repository.lambda_repo.arn
}
```

---

### 2. IAM Module (`modules/iam/`)

**Responsabilidade**: Gerencia roles e políticas IAM para execução da Lambda.

**Recursos criados**:

- `aws_iam_role.lambda_exec` - Role de execução da Lambda
- `aws_iam_role_policy_attachment.lambda_basic_logs` - Anexa política de logs

**Inputs**:

```hcl
variable "env" {
  type        = string
  description = "Nome do ambiente"
}
```

**Outputs**:

```hcl
output "lambda_execution_role_arn" {
  value = aws_iam_role.lambda_exec.arn
}

output "lambda_policy_attachment_name" {
  value = aws_iam_role_policy_attachment.lambda_basic_logs.id
}
```

---

### 3. Lambda Module (`modules/lambda/`)

**Responsabilidade**: Gerencia a função Lambda, alias e concorrência provisionada.

**Recursos criados**:

- `aws_lambda_function.golang_lambda` - Função Lambda principal
- `aws_lambda_alias.staging` - Alias para staging
- `aws_lambda_provisioned_concurrency_config.pc` - Concorrência provisionada (opcional)

**Inputs**:

```hcl
variable "env" {
  type        = string
  description = "Nome do ambiente"
}

variable "image_tag" {
  type        = string
  default     = "latest"
  description = "Tag da imagem Docker"
}

variable "lambda_execution_role_arn" {
  type        = string
  description = "ARN da role de execução do Lambda"
}

variable "ecr_repository_url" {
  type        = string
  description = "URL do repositório ECR"
}

variable "provisioned_concurrency" {
  type        = number
  description = "Número de execuções concorrentes provisionadas"
  default     = 0
}
```

**Outputs**:

```hcl
output "lambda_function_arn" {
  value = aws_lambda_function.golang_lambda.arn
}

output "lambda_function_name" {
  value = aws_lambda_function.golang_lambda.function_name
}

output "lambda_invoke_arn" {
  value = aws_lambda_function.golang_lambda.invoke_arn
}

output "lambda_alias_arn" {
  value = aws_lambda_alias.staging.arn
}
```

---

### 4. API Gateway Module (`modules/apigw/`)

**Responsabilidade**: Expõe a Lambda através de uma HTTP API.

**Recursos criados**:

- `aws_apigatewayv2_api.http_api` - HTTP API
- `aws_apigatewayv2_integration.lambda_integration` - Integração com Lambda
- `aws_apigatewayv2_route.health_route` - Rota /health
- `aws_apigatewayv2_route.any_root` - Rota ANY /
- `aws_apigatewayv2_route.any_proxy` - Rota ANY /{proxy+}
- `aws_apigatewayv2_stage.default_stage` - Stage padrão
- `aws_lambda_permission.api_invoke` - Permissão para API invocar Lambda

**Inputs**:

```hcl
variable "env" {
  type        = string
  description = "Nome do ambiente"
}

variable "lambda_invoke_arn" {
  type        = string
  description = "ARN de invocação da Lambda"
}

variable "lambda_function_name" {
  type        = string
  description = "Nome da função Lambda"
}
```

**Outputs**:

```hcl
output "api_endpoint" {
  value = aws_apigatewayv2_api.http_api.api_endpoint
}

output "api_id" {
  value = aws_apigatewayv2_api.http_api.id
}

output "api_execution_arn" {
  value = aws_apigatewayv2_api.http_api.execution_arn
}
```

---

### 5. CloudWatch Module (`modules/cloudwatch/`)

**Responsabilidade**: Gerencia logs e alarmes de monitoramento.

**Recursos criados**:

- `aws_cloudwatch_log_group.lambda_logs` - Grupo de logs da Lambda
- `aws_cloudwatch_metric_alarm.lambda_errors` - Alarme para erros

**Inputs**:

```hcl
variable "env" {
  type        = string
  description = "Nome do ambiente"
}

variable "lambda_function_name" {
  type        = string
  description = "Nome da função Lambda"
}

variable "sns_topic_arn" {
  type        = string
  description = "ARN do tópico SNS para alertas"
}
```

---

### 6. SNS Module (`modules/sns/`)

**Responsabilidade**: Gerencia notificações e alertas.

**Recursos criados**:

- `aws_sns_topic.alerts` - Tópico de alertas (condicional)
- `aws_sns_topic_subscription.email` - Inscrição por email (condicional)

**Inputs**:

```hcl
variable "env" {
  type        = string
  description = "Nome do ambiente"
}

variable "alert_email" {
  type        = string
  description = "Email para receber alertas"
}
```

**Outputs**:

```hcl
output "topic_arn" {
  value = length(aws_sns_topic.alerts) > 0 ? aws_sns_topic.alerts[0].arn : ""
}
```

---

### 7. Route53 Module (`modules/route53/`)

**Responsabilidade**: Gerencia health checks da aplicação.

**Recursos criados**:

- `aws_route53_health_check.lambda_health` - Health check HTTP

**Inputs**:

```hcl
variable "env" {
  type        = string
  description = "Nome do ambiente"
}

variable "api_endpoint" {
  type        = string
  description = "URL da API Gateway"
}
```

**Outputs**:

```hcl
output "health_check_id" {
  value = aws_route53_health_check.lambda_health.id
}
```

---

## 🔗 Dependências entre Módulos

```
main.tf orquestra as seguintes dependências:

ECR Module (independente)
  ↓
IAM Module (independente)
  ↓
Lambda Module (depende: ECR + IAM)
  ↓
API Gateway Module (depende: Lambda)
  ↓
CloudWatch Module (depende: Lambda + SNS)
Route53 Module (depende: API Gateway)

SNS Module (independente, usado pelo CloudWatch)
```

## 🚀 Como Usar os Módulos

### Exemplo de uso individual:

```hcl
module "lambda" {
  source = "./modules/lambda"

  env                        = "staging"
  image_tag                  = "latest"
  lambda_execution_role_arn  = "arn:aws:iam::123456789012:role/lambda-role"
  ecr_repository_url         = "123456789012.dkr.ecr.us-east-1.amazonaws.com/repo"
  provisioned_concurrency    = 0
}
```

### Exemplo de orquestração completa:

Veja o arquivo `main.tf` para um exemplo completo de como os módulos são orquestrados juntos.

---

## 🛠️ Manutenção

Para adicionar novos módulos:

1. Crie o diretório `modules/nome-do-modulo/`
2. Adicione os arquivos: `main.tf`, `variables.tf`, `outputs.tf`
3. Defina os recursos no `main.tf`
4. Exponha as saídas necessárias no `outputs.tf`
5. Documente as variáveis no `variables.tf`
6. Adicione o módulo ao `main.tf` principal
7. Atualize esta documentação

---

## 📚 Referências

- [Terraform Module Documentation](https://www.terraform.io/docs/language/modules/index.html)
- [AWS Provider Documentation](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)
- [Terraform Best Practices](https://www.terraform.io/docs/cloud/guides/recommended-practices/index.html)
