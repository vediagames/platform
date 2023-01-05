variable "project_id" {
  description = "GCP project ID"
}

variable "region" {
  description = "GCP region"
}

variable "zone" {
  description = "GCP region zone"
}

variable "credentials_path" {
  description = "Path to GCP credentials (JSON)"
}

variable "project" {
  description = "Project name (not necessarily related to GCP)"
}

variable "billing_account" {
  description = "GCP billing Account"
}

variable "pagerduty_gcp_billing_service_key" {
  description = "PagerDuty Events API v2 key for GCP Notification Channel"
}

variable "gke_num_nodes" {
  default     = 2
  description = "Amount of worker nodes"
}

variable "bigquery_events_plays_schema_path" {
  description = "File path for BigQuery schema for events.plays (JSON)"
}

variable "bigquery_events_likes_schema_path" {
  description = "File path for BigQuery schema for events.likes (JSON)"
}

variable "bigquery_events_dislikes_schema_path" {
  description = "File path for BigQuery schema for events.dislikes (JSON)"
}

variable "pubsub_example_schema_path" {
  description = "File path for Pub/Sub message schema for 'example' topic (JSON)"
}

variable "authorized_source_ranges" {
  type = list(string)
  description = "A list of CIDR addresses that are authorized to connect to GKE"
}
