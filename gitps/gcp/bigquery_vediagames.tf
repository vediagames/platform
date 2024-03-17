resource "google_bigquery_dataset" "vediagames" {
  dataset_id  = "vediagames"
  description = "Dataset for vediagames.com platform analytics and events."
  location    = var.region
  project     = var.project_id

  labels = {
    env = var.project_id
  }
}

resource "google_bigquery_table" "sessions" {
  dataset_id = google_bigquery_dataset.vediagames.dataset_id
  table_id   = "sessions"
  project    = var.project_id

  time_partitioning {
    type = "DAY"
  }

  labels = {
    env = var.project_id
  }

  schema = file(var.bigquery_sessions_schema_path)
}
