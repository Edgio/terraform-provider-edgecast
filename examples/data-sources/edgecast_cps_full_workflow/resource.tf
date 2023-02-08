resource "edgecast_cps_certificate" "my_cert" {
  certificate_label     = "my certificate 1"
  description           = "DV certificate for somedomain.com"
  auto_renew            = true
  certificate_authority = "DigiCert"
  validation_type       = "DV"
  dcv_method            = "DnsTxtToken"
  domain {
    is_common_name = true
    name           = "somedomain.com"
  }
}

data "edgecast_cps_dns_txt_token" "token" {
  certificate_id       = edgecast_cps_certificate.my_cert.id
  wait_until_available = true
}

#
# manual step: prove domain control using the DCV token above
#

data "edgecast_cps_target_cname" "target_cname" {
  certificate_id       = edgecast_cps_certificate.my_cert.id
  wait_until_available = true
  wait_timeout = "6h"

  depends_on = [
    data.edgecast_cps_dns_txt_token.token
  ]
}

output "dns_txt_token" {
  value = data.edgecast_cps_dns_txt_token.token.value
}

output "target_cname" {
  value = data.edgecast_cps_target_cname.target_cname.value
}