resource "ec_rules_engine_policy" "my_policy" {
  policy    = file("PathNormalization.json")
  deploy_to = "staging" # Valid values are "production" and "staging"
  account_number = "4FDBB"

  # Below are optional arguments when using PCC credentials
  # account_number = "A12345"
  # customeruserid = 1
  # portaltypeid   = 1
}
