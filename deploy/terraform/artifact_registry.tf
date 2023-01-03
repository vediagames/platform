resource "google_artifact_registry_repository" "primary" {
  location      = var.zone
  repository_id = "primary"
  description   = "Docker container registry"
  format        = "DOCKER"
}
