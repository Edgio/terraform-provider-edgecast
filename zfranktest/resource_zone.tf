resource "ec_dns_zone" "anyl" {
  account_number = "4FDBB"
	domain_name = "anyl2.com."
  status = 1
	zone_type = 1
	is_customer_owned = true
	comment = "test chang"
	record_a {
    name="mail"
    ttl=3600
    rdata="10.10.10.44"
  }
  record_a {
			name="www"
      ttl=3600
      rdata="10.10.10.114"
  }
  record_a {
			name="news2"
      ttl=3600
      rdata="10.10.10.100"
  }
  record_a {
			name="news3"
      ttl=3600
      rdata="10.10.10.200"
  }
  record_aaaa {
			name="www"
      ttl="3600"
      rdata="10:0:1::0:3"
  }
  record_cname {
			name="www2"
      ttl=3600
      rdata="www.cooler.com"
  }
  record_cname {
			name="www3"
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
        ttl=300
        rdata="10.10.1.11"
      }
    }
    a {
      weight=0
      record {
        ttl=300
        rdata="10.10.1.13"
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
        ttl=300
        rdata="10.10.2.21"
      }
    }
    a {
      weight=0
      record {
        ttl=300
        rdata="10.10.2.22"
      }
    }
  }
  

  dnsroute_group {
    group_type="zone"
    group_product_type="loadbalancing"
    name="lbg"
    a {
      weight=33
      health_check {
        check_interval=150
        check_type_id=1
        content_verification="10"
        email_notification_address="notice@glory1.com"
        failed_check_threshold=10
        http_method_id=1
        # ip_address="" # IP address only required when check_type_id 3,4
        ip_version=1
        # port_number=80 # Port only required when check_type_id 3,4
        reintegration_method_id=1
        uri="www.yahoo.com"
        timeout=100
      }
      record {
        ttl=300
        rdata="10.10.3.1"
      }
    }
    a {
      weight=33
      record {
        ttl=300
        rdata="10.10.3.2"
      }
    }
  }
}