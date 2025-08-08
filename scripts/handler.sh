#!/bin/bash

echo "🔨 Compilando função Go..."
GOOS=linux GOARCH=amd64 go build -o bootstrap app/main.go
zip bootstrap.zip bootstrap
echo "✅ Pronto: bootstrap.zip gerado"
