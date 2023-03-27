terraform {
  required_providers {
    edgecast = {
      version = "1.2.1"
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
    api_address_legacy = string
    ids_address        = string
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
  api_address        = var.credentials.api_address
  api_address_legacy = var.credentials.api_address_legacy
  ids_address        = var.credentials.ids_address
  partner_id         = 3249
  partner_user_id    = 7613
}
