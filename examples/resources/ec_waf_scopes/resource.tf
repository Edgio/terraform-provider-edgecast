resource "ec_waf_scopes" "scopes1" {
  account_number = var.account_number

  scope {
    name           = "Scopes Test"
    
    host {
      is_case_insensitive = false
      type                = "EM"
      values              = ["site1.com/path1","site2.com"]
    }

    limit {
      id                   = "<Rate Rule ID #1>"
      duration_sec         = 60
      enf_type             = "CUSTOM_RESPONSE"
      name                 = "limit 1 custom response"
      status               = 404
      response_body_base64 = "SGVsbG8sIHdvcmxkIQo="
      response_headers = {
        "header1" = "val1"
        "header2" = "val1"
      }
    }

    limit {
      id           = "<Rate Rule ID #2>"
      duration_sec = 300
      enf_type     = "DROP_REQUEST"
      name         = "limit 2 drop request"
    }

    limit {
      id           = "<Rate Rule ID #3>"
      duration_sec = 10
      enf_type     = "REDIRECT_302"
      name         = "limit 3 redirect"
      url          = "https://mysite.com/redirected"
    }

    path {
      is_case_insensitive = false
      is_negated          = false
      type                = "GLOB"
      value               = "*"
    }

    acl_audit_action {
      name = "audit action"
      type = "ALERT"
    }

    acl_audit_id = "<Access Rule ID>"

    acl_prod_action {
      valid_for_sec = 60
      enf_type      = "ALERT"
    }

    acl_prod_id = "<Access Rule ID>"

    profile_audit_action {
      name = "profile action"
      type = "ALERT"
    }

    profile_audit_id = "<Managed Rule ID>"

    profile_prod_action {
      valid_for_sec        = 60
      enf_type             = "CUSTOM_RESPONSE"
      name                 = "profile action"
      response_body_base64 = "SGVsbG8sIHdvcmxkIQo="
      response_headers = {
        "header1" = "val1"
        "header2" = "val1"
      }
      status = 404
    }

    profile_prod_id = "<Managed Rule ID>"

    rule_audit_action {
      name = "custom rule action"
      type = "ALERT"
    }

    rule_audit_id = "<Custom Rule ID>"

    rule_prod_action {
      valid_for_sec        = 60
      enf_type             = "BLOCK_REQUEST"
      name                 = "profile action"
    }

    rule_prod_id = "<Custom Rule ID>"
  }

}
