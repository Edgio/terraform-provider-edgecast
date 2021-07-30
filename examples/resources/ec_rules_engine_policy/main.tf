# Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

terraform {
  required_providers {
    ec = {
      version = "0.4.2"
      source  = "github.com/terraform-providers/ec"
    }
  }
}

##########################################m
# Variables
##########################################

variable credentials {
  type = object ({
    api_token = string
    ids_client_secret = string
    ids_client_id = string
    ids_scope = string
    api_address = string
    api_address_legacy = string
    ids_address = string
  })
}

variable rulesEngineEnvironment {
  type = string
}

variable test_customer_info {
  type = object({
    account_number = string
    customeruserid = string
    portaltypeid = string
  })
  default = {
    account_number = ""
    customeruserid = ""
    portaltypeid = "1"
  }
}

variable policy {
  type = string
  default = null
}
##########################################
# Providers
##########################################

provider "ec" {
    api_token = var.credentials.api_token
    ids_client_secret = var.credentials.ids_client_secret
    ids_client_id = var.credentials.ids_client_id
    ids_scope = var.credentials.ids_scope
    ids_address = var.credentials.ids_address
    api_address = var.credentials.api_address
    api_address_legacy = var.credentials.api_address_legacy
}

##########################################
# Resources
##########################################
resource "ec_rules_engine_policy" "httplarge_policy"{
  policy = var.policy != null ? var.policy : file("httplarge-policy.json")
  deploy_to = var.rulesEngineEnvironment

  account_number = var.test_customer_info.account_number
  customeruserid = var.test_customer_info.customeruserid
  portaltypeid = var.test_customer_info.portaltypeid
}
