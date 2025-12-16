terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
    }
  }
  cloud {

    organization = "example-org-5e1658"

    workspaces {
      name = "sunrise-app-aws"
    }
  }
}

# Remote State, necessary for placing GCP access key in secret manager
data "terraform_remote_state" "gcp" {
  backend = "remote"

  config = {
    organization = "example-org-5e1658"
    workspaces = {
      name = "sunrise-app-gcp"
    }
  }
}

# Configure the AWS Provider
provider "aws" {
  region     = "us-west-1"
}
