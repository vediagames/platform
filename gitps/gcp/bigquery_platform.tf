resource "google_bigquery_dataset" "platform" {
  dataset_id  = "platform"
  description = "Dataset for we gaming platform-related content."
  location    = var.region
  project     = var.project_id

  labels = {
    env = var.project_id
  }
}

resource "google_bigquery_table" "games" {
  dataset_id = google_bigquery_dataset.platform.dataset_id
  table_id   = "games"
  project    = var.project_id

  time_partitioning {
    type = "DAY"
  }

  labels = {
    env = var.project_id
  }

  schema = file(var.bigquery_games_schema_path)
}

resource "google_bigquery_table" "game_texts" {
  dataset_id = google_bigquery_dataset.platform.dataset_id
  table_id   = "game_texts"
  project    = var.project_id

  time_partitioning {
    type = "DAY"
  }

  labels = {
    env = var.project_id
  }

  schema = file(var.bigquery_game_texts_schema_path)
}

resource "google_bigquery_table" "game_tags" {
  dataset_id = google_bigquery_dataset.platform.dataset_id
  table_id   = "game_tags"
  project    = var.project_id

  time_partitioning {
    type = "DAY"
  }

  labels = {
    env = var.project_id
  }

  schema = file(var.bigquery_game_tags_schema_path)
}

resource "google_bigquery_table" "available_languages" {
  dataset_id = google_bigquery_dataset.platform.dataset_id
  table_id   = "available_languages"
  project    = var.project_id

  time_partitioning {
    type = "DAY"
  }

  labels = {
    env = var.project_id
  }

  schema = file(var.bigquery_available_languages_schema_path)
}

resource "google_bigquery_table" "categories" {
  dataset_id = google_bigquery_dataset.platform.dataset_id
  table_id   = "categories"
  project    = var.project_id

  time_partitioning {
    type = "DAY"
  }

  labels = {
    env = var.project_id
  }

  schema = file(var.bigquery_categories_schema_path)
}

resource "google_bigquery_table" "category_texts" {
  dataset_id = google_bigquery_dataset.platform.dataset_id
  table_id   = "category_texts"
  project    = var.project_id

  time_partitioning {
    type = "DAY"
  }

  labels = {
    env = var.project_id
  }

  schema = file(var.bigquery_category_texts_schema_path)
}

resource "google_bigquery_table" "game_categories" {
  dataset_id = google_bigquery_dataset.platform.dataset_id
  table_id   = "game_categories"
  project    = var.project_id

  time_partitioning {
    type = "DAY"
  }

  labels = {
    env = var.project_id
  }

  schema = file(var.bigquery_game_categories_schema_path)
}

resource "google_bigquery_table" "tag_texts" {
  dataset_id = google_bigquery_dataset.platform.dataset_id
  table_id   = "tag_texts"
  project    = var.project_id

  time_partitioning {
    type = "DAY"
  }

  labels = {
    env = var.project_id
  }

  schema = file(var.bigquery_tag_texts_schema_path)
}

resource "google_bigquery_table" "tags" {
  dataset_id = google_bigquery_dataset.platform.dataset_id
  table_id   = "tags"
  project    = var.project_id

  time_partitioning {
    type = "DAY"
  }

  labels = {
    env = var.project_id
  }

  schema = file(var.bigquery_tags_schema_path)
}
