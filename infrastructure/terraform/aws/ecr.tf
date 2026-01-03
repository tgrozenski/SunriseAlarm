# ECR Repository to store lambda container image
resource "aws_ecr_repository" "sunrise_repo" {
  name         = "sunrise_repo"
  force_delete = true
}

# Recursively build a hash from all go src files
locals {
  config_service_files = fileset("${path.module}/lambda/config_service", "**/*")
  config_service_tag = md5(join("", [
    for f in local.config_service_files : filemd5("${path.module}/lambda/config_service/${f}")
  ]))
}

resource "null_resource" "docker_build_push_config_service" {
  triggers    = {
    image_tag = local.config_service_tag
  }

  provisioner "local-exec" {
    command = <<EOF
      aws ecr get-login-password --region ${var.aws_region} | docker login --username AWS --password-stdin ${aws_ecr_repository.sunrise_repo.repository_url}
      docker build -t ${aws_ecr_repository.sunrise_repo.repository_url}:${local.config_service_tag} ${path.module}/lambda/config_service
      docker push ${aws_ecr_repository.sunrise_repo.repository_url}:${local.config_service_tag}
    EOF
  }

  depends_on = [aws_ecr_repository.sunrise_repo]
}

