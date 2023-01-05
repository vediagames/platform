resource "google_container_cluster" "primary-autopilot" {
  name     = "primary-autopilot"
  location = var.region
  project  = var.project_id

  network    = google_compute_network.primary.name
  subnetwork = google_compute_subnetwork.gke_primary.name

  enable_autopilot = true

  private_cluster_config {
    enable_private_endpoint = false
    enable_private_nodes    = true
  }

  master_authorized_networks_config {
    dynamic "cidr_blocks" {
      for_each = var.authorized_source_ranges
      content {
        cidr_block = cidr_blocks.value
      }
    }
  }

  maintenance_policy {
    recurring_window {
      start_time = "2023-01-05T00:23:35Z"
      recurrence = "FREQ=WEEKLY"
      end_time   = "2050-01-01T00:00:00Z"
    }
  }

  ip_allocation_policy {
    cluster_secondary_range_name  = "pods"
    services_secondary_range_name = "services"
  }

  release_channel {
    channel = "REGULAR"
  }
}
