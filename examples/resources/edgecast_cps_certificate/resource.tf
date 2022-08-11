resource "edgecast_cps_certificate" "certificate_1" {

	certificate_label = "cdn example tf 4"
	description = "cdn example 2"
	auto_renew = true
	certificate_authority = "DigiCert"
	validation_type = "OV"
	organization {
		city =               "L.A."
		company_address =    "111 fantastic way"
		company_name =    "Test Co."
		contact_email =       "user3@test.com"
		contact_first_name =   "test3"
		contact_last_name =  "user"
		contact_phone =   "111-111-1111"
		contact_title =      "N/A"
		country =      "US"
		organizational_unit = "Dept1"
		state =             "CA"
		zip_code =            "90001"
	}
	dcv_method= "Email"
	domains {
				is_common_name = true
				name =  "testssdomain.com"
			}
}
