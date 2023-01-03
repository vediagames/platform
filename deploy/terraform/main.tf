terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "4.47.0"
    }
  }
}

provider "google" {
  credentials = file(var.credentials_path)

  project = var.project_id
  region  = var.region
  zone    = var.zone
}
