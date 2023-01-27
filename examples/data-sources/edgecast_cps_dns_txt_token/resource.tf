resource "edgecast_cps_certificate" "my_cert" {
  certificate_label     = "retry demo cert 15"
  description           = "DV certificate for stevenpaz.com"
  auto_renew            = true
  certificate_authority = "DigiCert"
  validation_type       = "DV"
  dcv_method            = "DnsTxtToken"
  domain {
    is_common_name = true
    name           = "stevenpaz.com"
  }
}

data "edgecast_cps_dns_txt_token" "token" {
  certificate_id       = edgecast_cps_certificate.my_cert.id
  wait_until_available = true

   timeouts {
    read = "4h"
  }
}

output "dns_txt_token" {
  value = data.edgecast_cps_dns_txt_token.token.value
}