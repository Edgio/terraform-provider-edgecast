resource "ec_dns_tsig" "tsigszg1" {
      account_number = "4FDBB"
	alias = "tsig42"
	key_name ="kn1"
	key_value="0c94515c15e5095b8a87a50ba0df3bf38ed05fe6"
	algorithm_name = "HMAC-MD5"
}

resource "ec_dns_masterservergroup" "master_server_group_szg"{
  account_number = "4FDBB"
  master_server_group_name = "msg42"
  masters {
        name="ns1.test.com"
        ipaddress="10.10.10.1"
        }
  masters {
        name="ns2.test.com"
        ipaddress="10.10.10.2"
        }
}

resource "ec_dns_secondaryzonegroup" "backup" {
      account_number = "4FDBB"
      name="second4"
      zone_composition {
            master_group_id = ec_dns_masterservergroup.master_server_group_szg.id
            zones {
                  domain_name = "second4.com"
                  status=1
                  zone_type=2
                  comment="comment2"
            }
            master_server_tsigs{
                  master_server {
                        master_server_id = ec_dns_masterservergroup.master_server_group_szg.masters[0].id
                  }
                  tsig {
                        tsig_id = ec_dns_tsig.tsigszg1.id
                  }
            }
            master_server_tsigs{
                  master_server {
                        master_server_id = ec_dns_masterservergroup.master_server_group_szg.masters[1].id
                  }
                  tsig {
                        tsig_id = ec_dns_tsig.tsigszg1.id
                  }
            }

      }
}