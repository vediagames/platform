resource "google_compute_network" "primary" {
  name    = "primary"
  project = var.project_id
}

resource "google_compute_subnetwork" "gke_primary" {
  name                     = "gke-primary"
  region                   = var.region
  network                  = google_compute_network.primary.name
  ip_cidr_range            = "10.10.0.0/24"
  project                  = var.project_id
  private_ip_google_access = true
}

resource "google_compute_network" "production" {
  name                    = "production"
  project                 = var.project_id
  auto_create_subnetworks = true
}
