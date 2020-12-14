# Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

terraform {
  required_providers {
    vmp = {
      version = "0.1"
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
    partner_user_id = number
    partner_id = number
  })
  default = {
    api_address = ""
    api_token = ""
    ids_client_secret = ""
    ids_client_id = ""
    ids_scope = ""
    partner_user_id = 0
    partner_id = 0
  }
}

variable "new_customer_info" {
  type = object({
    company_name = string
    service_level_code = string
    services = list(number)
    access_modules = list(number)
    delivery_region = number
  })
  default = {
    company_name = "new customer1"
    service_level_code = "STND"
    services = []
    access_modules = []
    delivery_region = 1
  }
}

variable "new_admin_user" {
  type = object({
    first_name = string
    last_name = string
    email = string
    is_admin = bool
  })
  default = {
    first_name = "admin"
    last_name = "user"
    email = ""
    is_admin = true
  }
}

variable "origin_info" {
  type = object({
    origins = list(string)
    load_balancing = string
    host_header = string
    directory_name = string
  })
  default = {
    origins = []
    load_balancing = "RR"
    host_header = ""
    directory_name = ""
  }
}

variable "cname_info" {
  type = object({
    cname = string
    type = number
    origin_type = number
  })
  default = {
    cname = ""
    type = 3
    origin_type = 80
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
	  partner_user_id = var.partner_info.partner_user_id
    partner_id = var.partner_info.partner_id
}
##########################################
# Data Sources - Read-only data from VM APIs
##########################################

# data "vmp_customer_services" "rules_engine" {
# }

# output "rules_engine_services" {
#   value = data.vmp_customer_services.rules_engine.services
# }

##########################################
# Resources
##########################################

resource "vmp_customer" "test_customer" {
  company_name = var.new_customer_info.company_name
  service_level_code = var.new_customer_info.service_level_code
  services = var.new_customer_info.services
  delivery_region =  var.new_customer_info.delivery_region

  #################################################################################
  # Optional Params
  #################################################################################

  # TODO use data source to get ID (need a new endpoint for this)
  access_modules = var.new_customer_info.access_modules

  # TODO: may need to change these into their own resource type... but need new endpoints
  # HTTP Large
  domain {
      type = 1 # TODO use data source to get ID
      url = "http://wpc.5F534.omegacdn.net"
    }
  
  # ADN
  domain {
      type = 30 # TODO use data source to get ID
      url = "http://adn.5F534.omegacdn.net"
    }
}

resource "vmp_customer_user" "test_customer_admin" {
  account_number = vmp_customer.test_customer.id
  first_name = var.new_admin_user.first_name
  last_name = var.new_admin_user.last_name
  email = var.new_admin_user.email != "" ? var.new_admin_user.email : "admin@${vmp_customer.test_customer.id}.com"
  is_admin = var.new_admin_user.is_admin
}

resource "vmp_origin" "images" {
    account_number = vmp_customer.test_customer.id
    directory_name = var.origin_info.directory_name
    host_header = var.origin_info.host_header
    http {
        load_balancing = var.origin_info.load_balancing
        hostnames = var.origin_info.origins
    }
}

resource "vmp_cname" "images" {
    account_number = vmp_customer.test_customer.id
    name = var.cname_info.cname
    type = var.cname_info.type
    origin_id = vmp_origin.images.id
    origin_type = var.cname_info.origin_type
}