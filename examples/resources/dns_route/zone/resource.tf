resource "vmp_dns_zone" "middlem" {
  account_number = "DE0B"
	domain_name = "middlem.com."
  status = 1
	zone_type = 1
	is_customer_owned = true
	comment = "test chang"
	record_a {
    name="mail"
    ttl=3600
    rdata="10.10.10.13"
  }
  record_a {
			name="www"
      ttl=3600
      rdata="10.10.10.24"
  }
  record_aaaa {
			name="www"
      ttl="3600"
      rdata="10:0:1::0:2"
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
        name="fo1"
        ttl=300
        rdata="10.10.10.2"
      }
    }
    a {
      weight=0
      record {
        name="fo1"
        ttl=300
        rdata="10.10.10.123"
      }
    }
  }
  dnsroute_group {
    group_type="zone"
    group_product_type="failover"
    name="fo2"
    a {
      weight=100
      health_check {
        check_interval=300
        check_type_id=1
        content_verification="10"
        email_notification_address="notice@glory1.com"
        failed_check_threshold=10
        http_method_id=1
        ip_address=""
        ip_version=1
        port_number="80"
        reintegration_method_id=1
        status= 4
        status_name="Unknown"
        uri="www.yahoo.com"
      }
      record {
        name="fo2"
        ttl=300
        rdata="10.10.10.6"
      }
    }
    a {
      weight=0
      record {
        fixed_group_id=-1
        group_id=-1
        fixed_record_id=-1
        name="fo2"
        ttl=300
        rdata="10.10.10.9"
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
        check_interval=300
        check_type_id=1
        content_verification="10"
        email_notification_address="notice@glory1.com"
        failed_check_threshold=10
        http_method_id=1
        ip_address=""
        ip_version=1
        port_number="80"
        reintegration_method_id=1
        status= 4
        status_name="Unknown"
        uri="www.yahoo.com"
        timeout=100
      }
      record {
        fixed_group_id=-1
        group_id=-1
        fixed_record_id=-1
        name="lbg1"
        ttl=300
        rdata="10.10.10.12"
      }
    }
    a {
      weight=33
      record {
        name="lbg2"
        ttl=300
        rdata="10.10.10.13"
      }
    }
  }
}