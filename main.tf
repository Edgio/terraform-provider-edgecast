# Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.
terraform {
  required_providers {
    ec = {
      version = "0.4.2"
      source  = "github.com/terraform-providers/ec"
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
    api_address_legacy = string
    ids_address = string
  })
}

##########################################
# Providers
##########################################
provider "ec" {
    api_token = var.credentials.api_token
    ids_client_secret = var.credentials.ids_client_secret
    ids_client_id = var.credentials.ids_client_id
    ids_scope = var.credentials.ids_scope
    api_address = var.credentials.api_address
    api_address_legacy = var.credentials.api_address_legacy
    ids_address = var.credentials.ids_address
    partner_id = 2970
}

##########################################
# Resources
##########################################
resource "ec_customer" "test_customer_02" {
  company_name       = "Test Company 02"
  service_level_code = "STND"
  services           = [1,25,42,38] #all available services=> 1:HTTP Large Object,2:HTTPS Large Object,3:HTTP Small Object,4:HTTPS Small Object,6:Windows,7:Advanced Reports,8:Real-Time Stats,9:Token Auth,10:Edge Performance Analytics,15:Origin Storage,16:RSYNC,19:ADN,20:Download Manager,21:ADNS,22:Dedicated Hosting,23:Edge Optimizer,25:DNS Route,26:DNS Zones,29:DNS Health Checks,31:Bandwidth By Report Code,32:DNS-Standard,33:DNS-Adaptive,34:DNS-APR,38:WAF,39:Analysis Engine,40:HTTP Rate Limiting,41:Basic Rules v4.0,42:Advanced Rules v4.0,43:Mobile Device Detection Rules v4.0,44:Rules Engine v4.0,47:Translate,48:Dynamic Cloud Packaging,49:Encrypted HLS,50:Origin Shield,51:Reports and Logs,52:Log Delivery,54:SSA,56:Encrypted Key Rotation,57:Real-Time Log Delivery,58:Report Builder,59:Dynamic Imaging,60:China Delivery,61:WAF Essential,62:Report Builder Users,63:Report Builder Rows,64:Report Builder Reports,65:Edge Functions,66:Certificate Provisioning,67:Edge-Insights,68:Edge Image Optimizer,69:Url Redirects,70:Azure Cloud Storage
  delivery_region    = 1              # 1:Global + Premium Asia,2:North America and Europe,3:Global Standard,5:Global + Premium Asia + China,6:Global + Premium Asia + India,7:Global + Premium Asia + China + India,8:Global + Premium Asia + LATAM,9:Global + Premium Asia + Premium China + LATAM
  #available access modules => 1:Storage, 4:Analytics, 5:Admin, 7:Customer Origin, 8:Purge/Load, 21:Users, 22:Company, 25:Country Filtering, 26:Token Auth, 27:Dashboard, 29:HTTP Large, 30:Edge CNAMEs, 32:Core Reports, 40:Token Auth, 46:Token Auth, 53:Cache Settings, 56:HTTP Large Object, 71:HTTP Streaming, 72:ADN, 73:Customer Origin, 74:Purge/Load, 75:Token Auth, 76:Country Filtering, 77:Edge CNAMEs, 78:Cache Settings, 79:Application Delivery Network, 81:Tools, 138:Query-String Caching, 139:Query-String Logging, 140:Compression, 144:Query-String Caching, 145:Query-String Logging, 146:Compression, 149:Smooth Streaming Player, 153:JW Player, 157:Raw Log Settings, 159:Traffic Summary, 160:Bandwidth, 161:Data Transferred, 162:Hits, 163:Cache Statuses, 164:Cache Hit Ratio, 166:CDN Storage, 168:Notes, 169:HTTP Large, 170:HTTPS Large, 171:HTTP Small, 172:HTTPS Small, 174:Flash, 175:ADN, 176:ADN SSL, 177:HTTP Large, 178:HTTPS Large, 179:HTTP Small, 180:HTTPS Small, 182:Flash, 183:ADN, 184:ADN SSL, 185:All Platforms, 186:HTTP Large, 187:HTTP Small, 189:Flash, 190:ADN, 191:All Platforms, 192:HTTP Large, 193:HTTP Small, 194:ADN, 195:All Platforms, 196:HTTP Large, 197:HTTP Small, 198:ADN, 204:Usage, 386:IPv4/IPv6, 387:Data Transferred, 409:Custom Reports, 410:Edge CNAMEs, 411:Notes, 412:All Platforms, 413:HTTP Large, 414:HTTP Small, 415:Flash, 416:ADN, 479:Token Generator, 501:Add Users, 502:Edit Users
  access_modules = [1, 4, 5, 7, 8, 21, 22, 25, 26, 27, 29, 30, 32, 40, 46, 53, 56, 71, 72, 73, 74, 75, 76, 77, 78, 79, 81, 138, 139, 140, 144, 145, 146, 149, 153, 157, 159, 160, 161, 162, 163, 164, 166, 168, 169, 170, 171, 172, 174, 175, 176, 177, 178, 179, 180, 182, 183, 184, 185, 186, 187, 189, 190, 191, 192, 193, 194, 195, 196, 197, 198, 204, 386, 387, 409, 410, 411, 412, 413, 414, 415, 416, 479, 501, 502]

  #optional parameters
  bandwidth_usage_limit        = 0
  data_transferred_usage_limit = 0
  address1                     = "123 Main Street"
  address2                     = ""
  billing_account_tag          = ""
  billing_address1             = ""
  billing_address2             = ""
  billing_city                 = ""
  billing_contact_email        = ""
  billing_contact_fax          = ""
  billing_contact_first_name   = ""
  billing_contact_last_name    = ""
  billing_contact_mobile       = ""
  billing_contact_phone        = ""
  billing_contact_title        = ""
  billing_country              = ""
  billing_rate_info            = ""
  billing_state                = ""
  billing_zip                  = ""
  city                         = "Los Angeles"
  contact_email                = "terraformtestuser002@sharedectest.com"
  contact_fax                  = ""
  contact_first_name           = ""
  contact_last_name            = ""
  contact_mobile               = ""
  contact_phone                = ""
  contact_title                = ""
  country                      = ""
  notes                        = ""
  state                        = "CA"
  website                      = ""
  zip                          = "90210"
}

resource "ec_customer_user" "my_customer_user" {
  account_number = ec_customer.test_customer_02.id
  first_name     = "AdminFN"
  last_name      = "AdminLN"
  email          = "terraformadmin0002@sharedectest.com"
  is_admin       = true # cannot be modified after user creation

  #optional
  address1 = ""
  address2 = ""
  city = ""
  country = ""
  fax = ""
  mobile = ""
  phone = ""
  time_zone_id = "0"
  title = ""
  zip = ""
}

resource "ec_origin" "origin_images" {
    account_number = ec_customer.test_customer_02.id
    directory_name = "images"
    media_type = "httplarge"
    host_header = "images.mysite.com"
    http {
        load_balancing = "RR"
        hostnames = ["images-origin-1.mysite.com","images-origin-2.mysite.com"]
    }
}

