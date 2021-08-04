resource "ec_dns_group" "failover1" {
  account_number = "DE0B"
    group_type="cname"
    group_product_type="failover"
    name="failover-101"
    a {
      weight=100
      record {
        name="hot1"
        ttl=300
        rdata="10.10.1.3"
      }
    }
    a {
      weight=0
      record {
        name="cold1"
        ttl=300
        rdata="10.10.1.4"
      }
    }
}

resource "ec_dns_group" "loadbalancing1" {
  account_number = "DE0B"
  group_type="cname"
  group_product_type="loadbalancing"
  name="loadbal-200"
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
      name="lbg1"
      ttl=300
      rdata="10.10.3.5"
    }
  }
  a {
    weight=33
    record {
      name="lbg2"
      ttl=300
      rdata="10.10.3.6"
    }
  }
    a {
    weight=33
    record {
      name="lbg3"
      ttl=300
      rdata="10.10.3.7"
    }
  }
}

resource "ec_dns_group" "loadbalancing2" {
  account_number = "DE0B"
  group_type="cname"
  group_product_type="loadbalancing"
  name="loadbal-300"
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
      name="lbg1"
      ttl=300
      rdata="10.10.3.5"
    }
  }
  a {
    weight=33
    record {
      name="lbg2"
      ttl=300
      rdata="10.10.3.6"
    }
  }
    a {
    weight=33
    record {
      name="lbg3"
      ttl=300
      rdata="10.10.3.7"
    }
  }
}