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

data "edgecast_cps_target_cname" "mycert_cname" {
  certificate_id       = edgecast_cps_certificate.my_cert.id
  wait_until_available = var.wait_for_ec_cps_data_sources
  wait_timeout = "20m"
}

output "target_cname" {
  value = data.edgecast_cps_target_cname.mycert_cname.value
}