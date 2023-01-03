resource "google_artifact_registry_repository" "primary" {
  location      = var.region
  repository_id = "primary"
  description   = "Docker container registry"
  format        = "DOCKER"
}
