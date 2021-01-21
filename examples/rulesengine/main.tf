# Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

terraform {
  required_providers {
    vmp = {
      version = "0.0.1"
      source = "github.com/terraform-providers/vmp"
    }
  }
}

##########################################
# Variables
##########################################

variable "partner_info" {
  type = object ({
    api_address = string
    api_token = string
    ids_client_secret = string
    ids_client_id = string
    ids_scope = string
  })
}
variable "httplarge_policy_request" {
  type = object({
    policyid = number
    policy = string
    customerid = string
    customeruserid = number
    portaltypeid = number
  })
  default = {
    policyid = 0
    policy = ""
    customerid = ""
    customeruserid = 0
    portaltypeid = 1
  }
}
##########################################
# Providers
##########################################

provider "vmp" {
    api_address = var.partner_info.api_address
    api_token = var.partner_info.api_token
    ids_client_secret = var.partner_info.ids_client_secret
    ids_client_id = var.partner_info.ids_client_id
    ids_scope = var.partner_info.ids_scope
}


resource "vmp_httplarge_re_policy" "httplarge_policy"{
  policyid = var.httplarge_policy_request.policyid
  customerid = var.httplarge_policy_request.customerid
  customeruserid = var.httplarge_policy_request.customeruserid
  portaltypeid = var.httplarge_policy_request.portaltypeid
  policy = var.httplarge_policy_request.policy #
}
// resource "vmp_httplarge_re_policy" "httplarge_policy"{
//   policyid = var.httplarge_policy_request.policyid
//   customerid = var.httplarge_policy_request.customerid
//   customeruserid = var.httplarge_policy_request.customeruserid
//   portaltypeid = var.httplarge_policy_request.portaltypeid
//   policy = jsonencode(jsondecode(file("${path.module}/httplarge-policy.json")))
// }