# ===== ECR =====
resource "aws_ecr_repository" "lambda_repo" {
  name                 = "${var.env}-lambda-go"
  image_tag_mutability = "MUTABLE"
  image_scanning_configuration { scan_on_push = true }
}

# Dica: a imagem deve existir no ECR antes de criar a Lambda (docker build/tag/push)