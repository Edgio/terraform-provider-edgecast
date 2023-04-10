resource "edgecast_cps_certificate" "certificate_2" {

	certificate_label = "sinclair cert update issue 2"
	description = "sinclair cert update issue 2"
	auto_renew = true
	certificate_authority = "DigiCert"
	validation_type = "DV"
	dcv_method = "DnsTxtToken"
	domain {
		is_common_name = true
		name =  "stevenpaz.com"
	}
	
        notification_setting {
                notification_type = "CertificateExpiring"
                enabled = true
                emails = ["steven.paz@edgecast.com"]
        } 
}
