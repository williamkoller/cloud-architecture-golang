resource "aws_ecr_repository" "lambda_repo" {
  name                 = "${var.env}-lambda-go"
  image_tag_mutability = "MUTABLE"
  image_scanning_configuration { scan_on_push = true }
  force_delete = true
}