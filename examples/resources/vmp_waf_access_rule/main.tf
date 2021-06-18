terraform {
  required_providers {
    vmp = {
      version = "0.2.0"
      source = "VerizonDigital/vmp"
    }
  }
}

##########################################
# Providers
##########################################
provider "vmp" {
    api_token = ""
    ids_client_secret = "var.credentials.ids_client_secret"
    ids_client_id = "var.credentials.ids_client_id"
    ids_scope = "var.credentials.ids_scope"
}
##########################################
# Resources
##########################################
resource "vmp_waf_access_rule" "access_rule_1" {
  account_number = "ABCDF"
  name = "Access Rule #1"
  response_header_name = "header-name"
  allowed_http_methods = ["GET", "POST"]
  allowed_request_content_types = ["application/json"]
  disallowed_extensions = ["ext1","ext2"]
  disallowed_headers = ["header1","header2"]

  asn {
      accesslist = [12, 200, 465]
      blacklist = [13, 201, 466]
      whitelist = [14, 202, 467]
  }

  cookie {
      accesslist = ["al1", "al2"]
      blacklist = ["bl1", "bl2"]
      whitelist = ["wl1", "wl2"]
  }

  country {
      accesslist = ["al1", "al2"]
      blacklist = ["val1", "val2"]
      whitelist = ["val3", "val4"]
  }

  ip {
      accesslist = ["al1", "al2"]
      blacklist = ["val1", "val2"]
      whitelist = ["val3", "val4"]
  }

  referer {
      accesslist = ["al1", "al2"]
      blacklist = ["val1", "val2"]
      whitelist = ["val3", "val4"]
  }

  url {
      accesslist = ["al1", "al2"]
      blacklist = ["val1", "val2"]
      whitelist = ["val3", "val4"]
  }

  user_agent {
      accesslist = ["al1", "al2"]
      blacklist = ["val1", "val2"]
      whitelist = ["val3", "val4"]
  }
}
