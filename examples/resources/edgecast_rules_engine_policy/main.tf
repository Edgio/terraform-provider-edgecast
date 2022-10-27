# Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
# See LICENSE file in project root for terms.

terraform {
  required_providers {
    edgecast = {
       version = "0.6.0"
      source  = "Edgio/edgecast"
    }
  }
}

##########################################m
# Variables
##########################################

variable "credentials" {
  sensitive = true
  type = object({
    ids_client_secret = string
    ids_client_id = string
    ids_scope = string
  })
}

##########################################
# Providers
##########################################

provider "edgecast" {
    ids_client_secret = var.credentials.ids_client_secret
    ids_client_id = var.credentials.ids_client_id
    ids_scope = var.credentials.ids_scope
}
