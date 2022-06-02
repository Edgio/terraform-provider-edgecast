# Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
# See LICENSE file in project root for terms.

terraform {
  required_providers {
    edgecast = {
      version = "0.5.1"
      source  = "github.com/terraform-providers/edgecast"
    }
  }
}

##########################################
# Variables
##########################################
variable "credentials" {
  type = object({
    api_token          = string
    ids_client_secret  = string
    ids_client_id      = string
    ids_scope          = string
    api_address        = string
    api_address_legacy = string
    ids_address        = string
  })
}

variable "account_number" {
  type = string
}

variable "policy_contents" {
  type = string
  default = "{\r\n    \"name\": \"test policy 2022-02-08.3\u00DF\",\r\n    \"description\": \"This is a test policy!\",\r\n    \"platform\": \"http_large\",\r\n    \"rules\": [\r\n        {\r\n            \"name\": \"rule1\",\r\n            \"description\": \"This is a test rule.\",\r\n            \"matches\": [\r\n                {\r\n                    \"type\": \"match.always\",\r\n                    \"features\": [\r\n                        {\r\n                            \"type\": \"feature.comment\",\r\n                            \"value\": \"Update this comment!\"\r\n                        }\r\n                    ]\r\n                }\r\n            ]\r\n        }\r\n    ]\r\n}"
}

data "local_file" "policy" {
  filename = "policy.json"
}

##########################################
# Providers
##########################################
provider "edgecast" {
  api_token          = var.credentials.api_token
  ids_client_secret  = var.credentials.ids_client_secret
  ids_client_id      = var.credentials.ids_client_id
  ids_scope          = var.credentials.ids_scope
  api_address        = var.credentials.api_address
  api_address_legacy = var.credentials.api_address_legacy
  ids_address        = var.credentials.ids_address
}
