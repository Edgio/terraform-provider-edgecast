resource "edgecast_rules_engine_policy" "my_policy" {
  policy    = file("policy.json")
  deploy_to = "staging" # Valid values are "production" and "staging"
}