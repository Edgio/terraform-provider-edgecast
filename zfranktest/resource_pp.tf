#resource "ec_rules_engine_policy" "my_policy" {
#  policy    = file("my-policy.json")
#  deploy_to = "production"

#  # for PCC users, otherwise will be ignored
#  account_number = "A12345"
#  customeruserid = 1
#  portaltypeid   = 1
#}

resource "ec_rules_engine_policy" "PathNormalization" {
  policy    = file("${path.module}/PathNormalization.json")
  deploy_to = "production"
  account_number = "4FDBB"

  # for PCC users, otherwise will be ignored
  #account_number =  var.account_no
  #customeruserid = 1
  #portaltypeid   = 1

}