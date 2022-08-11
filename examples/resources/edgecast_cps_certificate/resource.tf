resource "edgecast_cps_certificate" "certificate_1" {

	certificate_label = "cdn example tf ev"
	description = "cdn example"
	auto_renew = true
	certificate_authority = "DigiCert"
	validation_type = "EV"
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
		additional_contacts{
			first_name	= "contact1"
			last_name	= "contactlastname1"
			email	= "first.lastname@testuser.com"
			phone	= "111-111-1111"				
			title	= "contactManager"
			contact_type	= "EvApprover"
		}
		additional_contacts{
			first_name	= "contact2"
			last_name	= "contactlastname2"
			email	= "first.lastname@testuser.com"
			phone	= "111-111-2222"				
			title	= "contactAccount"
			contact_type	= "TechnicalContact"
		}
	}
	dcv_method = "Email"
	domains {
		is_common_name = true
		name =  "testdomain1.com"
	}
	domains {
		is_common_name = false
		name =  "testdomain2.com"
	}
}

resource "edgecast_cps_certificate" "certificate_2" {

	certificate_label = "cdn example tf dv"
	description = "cdn example dv"
	auto_renew = true
	certificate_authority = "DigiCert"
	validation_type = "DV"
	dcv_method = "DnsTxtToken"
	domains {
		is_common_name = true
		name =  "testdomain3.com"
	}
	domains {
		is_common_name = false
		name =  "testdomain4.com"
	}
}