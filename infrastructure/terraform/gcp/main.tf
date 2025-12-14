terraform {
  required_providers {
    google = {
        source  = "hashicorp/google"
        version = "7.13.0"
    }
    google-beta = {
      source  = "hashicorp/google-beta"
      version = "7.13.0"
    }
  }
#   cloud { # uncomment this after dev
#     organization = "example-org-5e1658" 

#     workspaces { 
#       name = "sunrise-app" 
#     } 
#   }

}
# Configure the GCP Providers
provider "google-beta" {
    project = var.project_id
    region  = "us-west2"
    zone    = "us-west2-a"
}

provider "google" {
    project = var.project_id
    region  = "us-west2"
    zone    = "us-west2-a"
}

# Enable required APIs
locals {
  apis = [
    "iam.googleapis.com",
    "firebase.googleapis.com",
    "cloudresourcemanager.googleapis.com",
  ]
}

resource "google_project_service" "enabled" {
  for_each = toset(local.apis)

  project            = var.project_id
  service            = each.value
  disable_on_destroy = false
}
