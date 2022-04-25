resource "ec_waf_rate_rule" "rate_rule_1" {
  account_number = "<Account Number>"
  name           = "Rate Rule #1"
  duration_sec   = 1
  keys           = ["IP", "USER_AGENT"]
  disabled       = false
  num            = 10

  condition_group {
    name = "Group 1"

    condition {
      target {
        type = "REMOTE_ADDR"
      }

      op {
        is_case_insensitive = false
        is_negated          = false
        type                = "IPMATCH"                  # only use op.type=IPMATCH when target.type=REMOTE_ADDR
        values              = ["10.10.2.3", "10.10.2.4"] # values array property is required when op.type=EM or op.type=IPMATCH
      }
    }

    condition {
      target {
        type = "REQUEST_METHOD"
      }

      op {
        is_case_insensitive = true
        is_negated          = true
        type                = "EM"
        values              = ["GET", "POST"] # values array property is required when op.type=EM or op.type=IPMATCH
      }
    }

    condition {
      target {
        type = "FILE_EXT"
      }

      op {
        is_case_insensitive = true
        is_negated          = false
        type                = "RX"
        value               = "(.*?)\\.(jpg|gif|doc|pdf)$" # value property is required when op.type=RX
      }
    }
  }
  
  condition_group {
    name = "Group 2"

    condition {
      target {
        type  = "REQUEST_HEADERS"
        value = "User-Agent" # only valid when type=REQUEST_HEADERS. 
      }

      op {
        is_case_insensitive = true
        is_negated          = false
        type                = "EM"
        values              = ["Mozilla/5.0", "Chrome/91.0.4472.114"] # values array property is required when op.type=EM or op.type=IPMATCH
      }

    }
    condition {
      target {
        type = "REQUEST_URI"
      }

      op {
        is_case_insensitive = true
        is_negated          = false
        type                = "EM"
        values              = ["/marketing", "/800001/mycustomerorigin"] # values array property is required when op.type=EM or op.type=IPMATCH
      }

    }
  }

}
