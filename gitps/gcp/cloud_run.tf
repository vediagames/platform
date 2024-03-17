resource "google_cloud_run_v2_service" "imagor" {
  name     = "imagor"
  location = "us-central1"
  ingress  = "INGRESS_TRAFFIC_ALL"

  template {
    scaling {
      max_instance_count = 3
    }

    containers {
      image = "shumc/imagor:1.4.2"
      env {
        name  = "IMAGOR_SECRET"
        value = "UNp%ETVpympqP9AyTXvR%CRTZah^qExUJA*Fx5%*WtQN6ExLMX^"
      }
      env {
        name  = "IMAGOR_SIGNER_TYPE"
        value = "sha256"
      }
      env {
        name  = "IMAGOR_SIGNER_TRUNCATE"
        value = "40"
      }
    }
  }

  traffic {
    type    = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
    percent = 100
  }
}

resource "google_cloud_run_v2_service_iam_binding" "imagor_binding" {
  project  = google_cloud_run_v2_service.imagor.project
  location = google_cloud_run_v2_service.imagor.location
  name     = google_cloud_run_v2_service.imagor.name
  role     = "roles/run.invoker"
  members = [
    "allUsers",
  ]
}
