

resource "edgecast_cps_certificate" "my_cert" {
	certificate_label = "test for retry 4"
	description = "DV certificate for cdn.example.com"
	auto_renew = true
	certificate_authority = "DigiCert"
	validation_type = "DV"
	dcv_method = "DnsTxtToken"
	domain {
		is_common_name = true
		name =  "testdomain2.frankintosh.us"
	}
}


data "edgecast_cps_dns_txt_token" "token" {
  certificate_id       = edgecast_cps_certificate.my_cert.id
  wait_until_available = true
  
  # timeouts {
  #   read = "10m"
  # }
}

# output "token" {
#   value = data.edgecast_cps_dns_txt_token.token
# }
