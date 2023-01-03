resource "google_monitoring_notification_channel" "pagerduty_gcp_billing" {
  display_name = "PagerDuty GCP Billing Service"
  type         = "pagerduty"

  labels = {
    service_key = var.pagerduty_gcp_billing_service_key
  }
}
