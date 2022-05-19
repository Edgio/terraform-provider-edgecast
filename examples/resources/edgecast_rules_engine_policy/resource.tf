resource "edgecast_rules_engine_policy" "my_policy" {
  policy    = file("policy.json")
  deploy_to = "staging" # Valid values are "production" and "staging"

  # Below are optional arguments when using PCC credentials
  # account_number = "A12345"
  # customeruserid = 1
  # portaltypeid   = 1
}
