[![Deploy Lambda Go to AWS](https://github.com/williamkoller/cloud-architecture-golang/actions/workflows/deploy.yml/badge.svg)](https://github.com/williamkoller/cloud-architecture-golang/actions/workflows/deploy.yml)

# 🚀 Lambda Go com Terraform na AWS

Este projeto implementa uma função **AWS Lambda** escrita em Go, empacotada como imagem Docker e provisionada automaticamente com **Terraform**.

A ideia é simples: você escreve seu código em Go ➡ compila ➡ empacota no Docker ➡ envia para o **Amazon ECR** ➡ Terraform provisiona tudo na AWS.

---

## 📦 Ferramentas AWS Utilizadas

- **AWS Lambda** – Executa a função Go.
- **Amazon ECR** – Armazena a imagem Docker da função.
- **API Gateway** – Disponibiliza a função como uma API HTTP.
- **CloudWatch Logs** – Armazena os logs da função.
- **SNS (Simple Notification Service)** – Envia alertas configurados.
- **IAM (Identity and Access Management)** – Controla permissões.

---

## 🔄 Fluxo do Projeto

1. Escrevemos o código Go.
2. Compilamos o binário (`bootstrap`).
3. Construímos a imagem Docker usando `public.ecr.aws/lambda/provided:al2`.
4. Fazemos push para o **Amazon ECR**.
5. O Terraform cria:
   - Função Lambda usando a imagem no ECR.
   - API Gateway para expor a função.
   - Permissões no IAM.
   - Alertas no SNS (opcional).
6. Chamamos a API e vemos a mágica acontecer ✨.

---

## 📋 Pré-requisitos

- Go 1.21+
- Docker
- AWS CLI configurado (`aws configure`)
- Terraform 1.5+

---

## ▶️ Como Executar

```bash
# 1️⃣ Compilar o Go e empacotar
./deploy.sh

# 2️⃣ Provisionar a infraestrutura
terraform init
terraform apply -var="env=staging" -var="alert_email=seu@email.com" -auto-approve
```
