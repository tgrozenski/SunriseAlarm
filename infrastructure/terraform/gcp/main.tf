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
    project = "sunrise-app-1374823832"
    region  = "us-west2"
    zone    = "us-west2-a"
}
provider "google" {
    project = "sunrise-app-1374823832"
    region  = "us-west2"
    zone    = "us-west2-a"
}

resource "google_compute_instance" "vm_instance" {
  name         = "terraform-instance"
  machine_type = "e2-micro"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }

  network_interface {
    # A default network is created for all GCP projects
    network = "default"
    access_config {
    }
  }
}