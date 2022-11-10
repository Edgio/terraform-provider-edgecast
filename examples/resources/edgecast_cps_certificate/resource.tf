resource "edgecast_cps_certificate" "certificate_1" {

	# Note: certificate_label must be unique, including deleted certificates.
	certificate_label = "cdn example com EV certificate"
	description = "EV certificate for cdn.example.com"
	auto_renew = true
	certificate_authority = "DigiCert"
	validation_type = "EV"
	organization {
		city =               "L.A."
		company_address =    "111 fantastic way"
		company_name =    "Example Co."
		contact_email =       "user3@example.com"
		contact_first_name =   "test3"
		contact_last_name =  "user"
		contact_phone =   "111-111-1111"
		contact_title =      "N/A"
		country =      "US"
		organizational_unit = "Dept1"
		state =             "CA"
		zip_code =            "90001"
		additional_contact{
			first_name	= "contact1"
			last_name	= "contactlastname1"
			email	= "first.lastname@example.com"
			phone	= "111-111-1111"				
			title	= "contactManager"
			contact_type	= "EvApprover"
		}
		additional_contact{
			first_name	= "contact2"
			last_name	= "contactlastname2"
			email	= "first.lastname@example.com"
			phone	= "111-111-2222"				
			title	= "contactAccount"
			contact_type	= "TechnicalContact"
		}
	}
	dcv_method = "Email"
	domain {
		is_common_name = true
		name =  "cdn.example.com"
	}
	domain {
		is_common_name = false
		name =  "cdn2.example.com"
	}
	notification_setting {
		notification_type = "CertificateRenewal"
		enabled = true
		emails = ["first.lastname@example.com"]
	}
	notification_setting {
		notification_type = "CertificateExpiring"
		enabled = true
		emails = ["first.lastname@example.com"]
	}
	notification_setting {
		notification_type = "PendingValidations"
		enabled = true
		emails = ["first.lastname@example.com"]
	}
}

resource "edgecast_cps_certificate" "certificate_2" {

	certificate_label = "cdn example com DV certificate"
	description = "DV certificate for cdn.example.com"
	auto_renew = true
	certificate_authority = "DigiCert"
	validation_type = "DV"
	dcv_method = "DnsTxtToken"
	domain {
		is_common_name = true
		name =  "cdn3.example.com"
	}
	domain {
		is_common_name = false
		name =  "cdn4.example.com"
	}
        notification_setting {
                notification_type = "CertificateExpiring"
                enabled = true
                emails = ["joe@example.com"]
        } 
}
