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
  wait_timeout = "20m"

  # There is a known issue where wait_until_available=true causes the provider
  # to wait for this data source on terraform plan.
  #
  # Use the var wait_for_ec_cps_data_sources (defined in main.tf) to force the
  # provider to not wait for data sources on terraform plan. Note that its
  # default value is false
  #
  # When you are ready to apply, set the variable:
  #
  # terraform apply -var 'wait_for_ec_cps_data_sources=true'
  wait_until_available = var.wait_for_ec_cps_data_sources
}

output "target_cname" {
  value = data.edgecast_cps_target_cname.mycert_cname.value
}