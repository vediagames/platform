resource "google_billing_budget" "budget_40e" {
  billing_account = var.billing_account
  display_name    = "Billing budget for 40 EUR"

  budget_filter {
    projects = ["projects/${var.project_number}"]
  }

  amount {
    specified_amount {
      currency_code = "EUR"
      units         = "40"
    }
  }

  threshold_rules {
    threshold_percent = 0.75
  }
  threshold_rules {
    threshold_percent = 0.9
    spend_basis       = "FORECASTED_SPEND"
  }

  all_updates_rule {
    monitoring_notification_channels = [
      google_monitoring_notification_channel.pagerduty_gcp_billing.id,
    ]
    disable_default_iam_recipients = true
  }
}
