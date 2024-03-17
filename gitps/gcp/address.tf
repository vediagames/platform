resource "google_compute_address" "primary_lb_ip" {
  name = "primary-lb-ipv4"
}

resource "google_compute_address" "gke_production_lb_ip" {
  name   = "gke-production-lb-ip"
  region = "us-central1"
}

