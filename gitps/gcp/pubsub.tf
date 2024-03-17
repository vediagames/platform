resource "google_pubsub_schema" "example" {
  name       = "example"
  type       = "PROTOCOL_BUFFER"
  definition = file(var.pubsub_example_schema_path)
}

resource "google_pubsub_topic" "example" {
  name = "example-topic"

  depends_on = [google_pubsub_schema.example]
  schema_settings {
    schema   = "projects/${var.project_id}/schemas/${google_pubsub_schema.example.name}"
    encoding = "JSON"
  }
}
#
#resource "google_pubsub_subscription" "example" {
#  name  = "example-subscription"
#  topic = google_pubsub_topic.example.name
#
#  ack_deadline_seconds = 20
#
#  labels = {
#    foo = "bar"
#  }
#
#  push_config {
#    push_endpoint = "https://example.com/push"
#
#    attributes = {
#      x-goog-version = "v1"
#    }
#  }
#}
