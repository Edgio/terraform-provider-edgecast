# Copyright 2021 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
# See LICENSE file in project root for terms.

terraform {
  required_providers {
    edgecast = {
      version = "1.3.4"
      source  = "github.com/terraform-providers/edgecast"
    }
  }
}

##########################################
# Variables
##########################################

variable "credentials" {
  sensitive = true
  type = object({
    api_token          = string
    ids_client_secret  = string
    ids_client_id      = string
    ids_scope          = string
    api_address        = string
    ids_address        = string
    api_address_legacy = string
  })
}

variable "account_number" {
  type = string
}

##########################################
# Providers
##########################################

provider "edgecast" {
  api_token          = var.credentials.api_token
  ids_client_secret  = var.credentials.ids_client_secret
  ids_client_id      = var.credentials.ids_client_id
  ids_scope          = var.credentials.ids_scope
  ids_address        = var.credentials.ids_address
  api_address        = var.credentials.api_address
  api_address_legacy = var.credentials.api_address_legacy
}
