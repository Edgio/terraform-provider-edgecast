# Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

terraform {
  required_providers {
    vmp = {
      version = "0.0.8"
      source = "VerizonDigital/vmp"
    }
  }
}

##########################################
# Variables
##########################################

variable "provider_config" {
  type = object ({
    api_address = string
    api_token = string
    ids_client_secret = string
    ids_client_id = string
    ids_scope = string
  })
}
variable "httplarge_policy" {
  type = string
  default = ""
}

variable "customer_info" {
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

##########################################
# Providers
##########################################

provider "vmp" {
    api_address = var.provider_config.api_address
    api_token = var.provider_config.api_token
    ids_client_secret = var.provider_config.ids_client_secret
    ids_client_id = var.provider_config.ids_client_id
    ids_scope = var.provider_config.ids_scope
}


resource "vmp_rules_engine_policy" "httplarge_policy"{
  account_number = var.customer_info.account_number
  customeruserid = var.customer_info.customeruserid
  portaltypeid = var.customer_info.portaltypeid
  policy = var.httplarge_policy
}