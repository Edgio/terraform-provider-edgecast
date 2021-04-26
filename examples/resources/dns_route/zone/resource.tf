resource "vmp_dns_zone" "glory_com" {
  account_number = ""
	domain_name = "glory.com."
  status = 1
	zone_type = 1
	is_customer_owned = 1
	comment = "sset trial"
	record_a {
    name="mail"
    ttl="3600"
    rdata="10.10.10.11"
    verify_id=1
  }
  record_a {
			name="www"
      ttl="3600"
      rdata="10.10.10.1"
      verify_id=2
  }
  record_aaaa {
			name="www"
      ttl="3600"
      rdata="10:0:1::0:2"
      verify_id=2
  }
  record_cname {
			name="www"
      ttl="3600"
      rdata="www.twice.com"
      verify_id=4
  }
	record_mx {
    name="@"
    ttl="3600"
    rdata="10 mail.twice.com"
    verify_id=1
  }
}