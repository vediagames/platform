resource "google_compute_network" "primary" {
  name    = "primary-network"
  project = var.project_id
}

resource "google_compute_subnetwork" "primary" {
  name          = "primary-subnet"
  region        = var.region
  network       = google_compute_network.primary.name
  ip_cidr_range = "10.10.0.0/24"
  project       = var.project_id
}
