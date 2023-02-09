resource "edgecast_cps_certificate" "my_cert" {
  certificate_label     = "my cert 1"
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
  wait_timeout = "20m"
}

output "dns_txt_token" {
  value = data.edgecast_cps_dns_txt_token.token.value
}