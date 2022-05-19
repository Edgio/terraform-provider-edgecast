
resource "ec_waf_bot_rule" "bot_rule_1" {
  customer_id = "<customer id>"
  name        = "Bot Rule Terraform Test"

  directive {
    include = "r3010_ec_bot_challenge_reputation.conf.json"
  }

  directive {
    sec_rule {
      name = "Sec Rule 2"
      action {
        id              = 77000000
        msg             = "Invalid user agent"
        transformations = ["NONE"]
      }
      operator {
        is_negated = false
        type       = "CONTAINS"
        value      = "myheadervalue"
      }
      variable {
        is_count = false
        type     = "REQUEST_HEADERS"
        match {
          is_negated = false
          is_regex   = false
          value      = "myheader"
        }
      }

      chained_rule {
        action {
          transformations = ["LOWERCASE"]
        }
        operator {
          is_negated = false
          type       = "BEGINSWITH"
          value      = "bot"
        }
        variable {
          is_count = false
          type     = "REQUEST_HEADERS"
          match {
            is_negated = false
            is_regex   = false
            value      = "User-Agent"
          }
        }
      }
    }
  }
}
