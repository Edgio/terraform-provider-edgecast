resource "edgecast_rules_engine_policy" "my_policy" {
  policy    = file("my-policy.json")
  deploy_to = "production"

  # for PCC users, otherwise will be ignored
  account_number = "A12345"
  customeruserid = 1
  portaltypeid   = 1
}
