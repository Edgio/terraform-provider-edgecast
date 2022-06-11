resource "edgecast_waf_custom_rule_set" "custom_rule_1" {
  account_number = "0001"
  name        = "Custom Rule Set 1"

  directive {
    sec_rule {
      name = "Sec Rule 1"
      action {
        id              = 66000000
        msg             = "Invalid user agent."
        transformations = ["NONE"]
      }
      operator {
        is_negated = false
        type       = "CONTAINS"
        value      = "bot"
      }
      variable {
        is_count = false
        match {
          is_negated = false
          is_regex   = false
          value      = "User-Agent"
        }
        match {
          is_negated = false
          is_regex   = false
          value      = "User-Agent"
        }
        type = "REQUEST_HEADERS"
      }
      variable {
        is_count = false
        match {
          is_negated = false
          is_regex   = false
          value      = "User-Agent"
        }
        type = "REQUEST_URI"
      }
      chained_rule {
        action {
          transformations = ["NONE"]
        }
        operator {
          is_negated = false
          type       = "CONTAINS"
          value      = "bot"
        }
        variable {
          is_count = false
          match {
            is_negated = false
            is_regex   = false
            value      = "User-Agent"
          }
          match {
            is_negated = false
            is_regex   = false
            value      = "User-Agent"
          }
          type = "REQUEST_HEADERS"
        }
      }
    }
  }
}
