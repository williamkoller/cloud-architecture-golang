output "repository_url" {
  value = aws_ecr_repository.lambda_repo.repository_url
}

output "repository_name" {
  value = aws_ecr_repository.lambda_repo.name
}

output "repository_arn" {
  value = aws_ecr_repository.lambda_repo.arn
} 