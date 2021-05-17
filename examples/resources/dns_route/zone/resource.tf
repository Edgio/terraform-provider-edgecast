resource "vmp_dns_zone" "anyk" {
  account_number = "DE0B"
	domain_name = "anyk.com."
  status = 1
	zone_type = 1
	is_customer_owned = true
	comment = "test chang"
	record_a {
    name="mail"
    ttl=3600
    rdata="10.10.10.45"
  }
  record_a {
			name="www"
      ttl=3600
      rdata="10.10.10.114"
  }
  record_a {
			name="news"
      ttl=3600
      rdata="10.10.10.200"
  }
  record_aaaa {
			name="www"
      ttl="3600"
      rdata="10:0:1::0:3"
  }
  record_cname {
			name="www"
      ttl=3600
      rdata="www.cooler.com"
  }
	record_mx {
    name="@"
    ttl=3600
    rdata="10 mail.cooler.com"
  }
  dnsroute_group {
    group_type="zone"
    group_product_type="failover"
    name="fo1"
    a {
      weight=100
      record {
        name="hot1"
        ttl=300
        rdata="10.10.1.1"
      }
    }
    a {
      weight=0
      record {
        name="cold1"
        ttl=300
        rdata="10.10.1.2"
      }
    }
  }
  dnsroute_group {
    group_type="zone"
    group_product_type="failover"
    name="fo2"
    a {
      weight=100
      record {
        name="hot3"
        ttl=300
        rdata="10.10.1.3"
      }
    }
    a {
      weight=0
      record {
        name="cold4"
        ttl=300
        rdata="10.10.1.4"
      }
    }
  }

}