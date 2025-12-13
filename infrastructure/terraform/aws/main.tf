terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
    }
  }
  # cloud { # Uncomment after dev
  #   organization = "example-org-5e1658" 

  #   workspaces { 
  #     name = "sunrise-app" 
  #   } 
  # }
}

# Configure the AWS Provider
provider "aws" {
  region     = "us-west-1"
}