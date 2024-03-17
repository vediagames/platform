terraform {
  required_providers {
    circleci = {
      source  = "mrolla/circleci"
      version = "0.6.1"
    }
  }
}

provider "circleci" {
  api_token    = var.circleci_token
  vcs_type     = "github"
  organization = "vediagames"
}
