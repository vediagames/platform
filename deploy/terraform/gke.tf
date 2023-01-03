resource "google_container_cluster" "primary" {
  name     = "primary-gke"
  location = var.zone
  project  = var.project_id

  remove_default_node_pool = true
  initial_node_count       = 1

  network    = google_compute_network.primary.name
  subnetwork = google_compute_subnetwork.primary.name
}

resource "google_container_node_pool" "primary_nodes" {
  name       = google_container_cluster.primary.name
  location   = var.zone
  cluster    = google_container_cluster.primary.name
  node_count = var.gke_num_nodes
  project    = var.project_id

  node_config {
    disk_size_gb = 25
    oauth_scopes = [
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]

    labels = {
      env  = var.project_id
      type = "primary"
    }

    machine_type = "e2-small"
    tags         = ["gke-node", "${google_container_cluster.primary.name}-node"]
    metadata = {
      disable-legacy-endpoints = "true"
    }
  }
}
