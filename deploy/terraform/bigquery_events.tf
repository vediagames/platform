resource "google_bigquery_dataset" "events" {
  dataset_id  = "events"
  description = "Dataset for platform analytics and events."
  location    = var.region
  project     = var.project_id

  labels = {
    env = var.project_id
  }
}

resource "google_bigquery_table" "plays" {
  dataset_id = google_bigquery_dataset.events.dataset_id
  table_id   = "plays"
  project    = var.project_id

  time_partitioning {
    type = "DAY"
  }

  labels = {
    env = var.project_id
  }

  schema = file(var.bigquery_events_plays_schema_path)
}

resource "google_bigquery_table" "likes" {
  dataset_id = google_bigquery_dataset.events.dataset_id
  table_id   = "likes"
  project    = var.project_id

  time_partitioning {
    type = "DAY"
  }

  labels = {
    env = var.project_id
  }

  schema = file(var.bigquery_events_likes_schema_path)
}

resource "google_bigquery_table" "dislikes" {
  dataset_id = google_bigquery_dataset.events.dataset_id
  table_id   = "dislikes"
  project    = var.project_id

  time_partitioning {
    type = "DAY"
  }

  labels = {
    env = var.project_id
  }

  schema = file(var.bigquery_events_dislikes_schema_path)
}
