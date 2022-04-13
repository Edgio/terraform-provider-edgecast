terraform {
  required_providers {
    ec = {
      version = "0.5.0"
      source  = "EdgeCast/ec"
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
  })
}

variable "account_number" {
  type = string
}

##########################################
# Providers
##########################################
provider "ec" {
  api_token          = var.credentials.api_token
  ids_client_secret  = var.credentials.ids_client_secret
  ids_client_id      = var.credentials.ids_client_id
  ids_scope          = var.credentials.ids_scope
}
