resource "vmp_waf_custom_rule" "custom_rule_1" {
  account_number = "ca5db"
  name           = "Custom Rule 1"
  
  directive {
      sec_rule {
          name = "Sec Rule 1"
          action {
              id = 66000000 
              msg = "Invalid user agent."
              t = ["NONE"]
          }
          operator {
              is_negated = false
              type = "CONTAINS"
              value = "bot"
          }
          variable {
              is_count = false
              match {
                is_negated = false
                is_regex = false
                value = "User-Agent"
              } 
              match {
                  is_negated = false
                  is_regex = false
                  value = "customer-Agent"
              }
              type = "REQUEST_HEADERS"
          }
          variable {
              is_count = false
              match {
                is_negated = false
                is_regex = false
                value = "User-Agent"
              } 
              type = "RESPONSE_HEADERS"
          }
          chained_rule {
                action {
                id = 66000001 
                msg = "Invalid user agent - chained."
                t = ["NONE"]
            }
            operator {
                is_negated = false
                type = "CONTAINS"
                value = "bot"
            }
            variable {
                is_count = false
                match {
                    is_negated = false
                    is_regex = false
                    value = "User-Agent"
                } 
                match {
                    is_negated = false
                    is_regex = false
                    value = "customer-Agent"
                }
                type = "REQUEST_HEADERS"
            }
          }
        }
    }
}
