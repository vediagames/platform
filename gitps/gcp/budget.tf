resource "google_billing_budget" "budget_40e" {
  billing_account = var.billing_account
  display_name    = "Billing budget for 40 EUR"
  budget_filter {
    projects = ["projects/148457546981", "projects/392196479804", "projects/827850532057", "projects/925609030723"]
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
}
