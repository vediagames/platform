resource "circleci_context" "github_creds" {
  name = "github_creds"
}

resource "circleci_context_environment_variable" "github_creds" {
  for_each = {
    GITHUB_TOKEN    = var.github_creds_token
    GITHUB_USERNAME = var.github_creds_username
  }

  variable   = each.key
  value      = each.value
  context_id = circleci_context.github_creds.id
}
