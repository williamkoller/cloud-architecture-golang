#!/bin/bash

set -e

ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text --region us-east-1)
REPO_NAME="staging-lambda-go"
IMAGE_TAG="latest"
IMAGE_URI="$ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/$REPO_NAME:$IMAGE_TAG"

echo "👉 Limpando binário antigo..."
rm -f bootstrap

echo "🚀 Recompilando Go com CGO desabilitado..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bootstrap app/main.go

echo "🧹 Limpando cache do Docker..."
docker builder prune --all --force

echo "🐳 Build da imagem..."
docker build --no-cache -t $REPO_NAME .

echo "🏷️ Tag para o ECR..."
docker tag $REPO_NAME:latest $IMAGE_URI

echo "🔐 Login no ECR..."
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin $ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com

echo "📤 Push para o ECR..."
docker push $IMAGE_URI

echo "✅ Imagem enviada com sucesso para o ECR: $IMAGE_URI"
