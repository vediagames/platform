variable "project_id" {
  type = string
  description = "GCP project ID"
}

variable "region" {
  type = string
  description = "GCP region"
}

variable "zone" {
  type = string
  description = "GCP region zone"
}

variable "credentials_path" {
  type = string
  description = "Path to GCP credentials (JSON)"
}

variable "project" {
  type = string
  description = "Project name (not necessarily related to GCP)"
}

variable "billing_account" {
  type = string
  description = "GCP billing Account"
}

variable "pagerduty_gcp_billing_service_key" {
  type = string
  description = "PagerDuty Events API v2 key for GCP Notification Channel"
}

variable "gke_num_nodes" {
  default     = 2
  description = "Amount of worker nodes"
}

variable "bigquery_sessions_schema_path" {
  type = string
  description = "File path for BigQuery schema for sessions (JSON)"
}

variable "bigquery_games_schema_path" {
  type = string
  description = "Path to the schema file for the games table."
}

variable "bigquery_game_texts_schema_path" {
  type = string
  description = "Path to the schema file for the game_texts table."
}

variable "bigquery_game_tags_schema_path" {
  type = string
  description = "Path to the schema file for the game_tags table."
}

variable "bigquery_available_languages_schema_path" {
  type = string
  description = "Path to the schema file for the available_languages table."
}

variable "bigquery_categories_schema_path" {
  type = string
  description = "Path to the schema file for the categories table."
}

variable "bigquery_category_texts_schema_path" {
  type = string
  description = "Path to the schema file for the category_texts table."
}

variable "bigquery_game_categories_schema_path" {
  type = string
  description = "Path to the schema file for the game_categories table."
}

variable "bigquery_tag_texts_schema_path" {
  type = string
  description = "Path to the schema file for the tag_texts table."
}

variable "bigquery_tags_schema_path" {
  type = string
  description = "Path to the schema file for the tags table."
}

variable "pubsub_example_schema_path" {
  type = string
  description = "File path for Pub/Sub message schema for 'example' topic (JSON)"
}

variable "authorized_source_ranges" {
  type = list(string)
  description = "A list of CIDR addresses that are authorized to connect to GKE"
}
