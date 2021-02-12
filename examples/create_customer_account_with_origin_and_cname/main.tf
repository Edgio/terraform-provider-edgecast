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

variable "credentials" {
  type = object ({
    api_token = string
    ids_client_secret = string
    ids_client_id = string
    ids_scope = string
    api_address = string
    ids_address = string
  })
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
    media_type = string
  })
  default = {
    origins = []
    load_balancing = "RR"
    host_header = ""
    directory_name = ""
    media_type = "httplarge"
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
    api_token = var.credentials.api_token
    ids_client_secret = var.credentials.ids_client_secret
    ids_client_id = var.credentials.ids_client_id
    ids_scope = var.credentials.ids_scope
    api_address = var.credentials.api_address
    ids_address = var.credentials.ids_address
}

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
}

resource "vmp_customer_user" "test_customer_admin" {
  account_number = vmp_customer.test_customer.id
  first_name = var.new_admin_user.first_name
  last_name = var.new_admin_user.last_name
  email = var.new_admin_user.email != "" ? var.new_admin_user.email : "admin@${vmp_customer.test_customer.id}.com"
  is_admin = var.new_admin_user.is_admin
}

resource "vmp_origin" "origin_images" {
    account_number = vmp_customer.test_customer.id
    directory_name = var.origin_info.directory_name
    media_type = var.origin_info.media_type
    host_header = var.origin_info.host_header
    http {
        load_balancing = var.origin_info.load_balancing
        hostnames = var.origin_info.origins
    }
}

resource "vmp_cname" "cname_images" {
    account_number = vmp_customer.test_customer.id
    name = var.cname_info.cname
    type = var.cname_info.type
    origin_id = vmp_origin.origin_images.id
    origin_type = var.cname_info.origin_type
}
