resource "edgecast_dns_tsig" "tsig1" {
      account_number = "DE0B"
	alias = "tsig42"
	key_name ="kn1"
	key_value="0c94515c15e5095b8a87a50ba0df3bf38ed05fe6"
	algorithm_name = "HMAC-MD5"
}

resource "edgecast_dns_masterservergroup" "master_server_group"{
  account_number = "DE0B"
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
      account_number = "DE0B"
      name="second48"
      zone_composition {
            master_group_id = edgecast_dns_masterservergroup.master_server_group.id
            zones {
                  domain_name = "second48.com"
                  status=1
                  zone_type= 2
                  comment="comment2"
            }
            master_server_tsigs{
                  master_server {
                        master_server_id = edgecast_dns_masterservergroup.master_server_group.masters[0].id
                  }
                  tsig {
                        tsig_id = edgecast_dns_tsig.tsig1.id
                  }
            }
            master_server_tsigs{
                  master_server {
                        master_server_id = edgecast_dns_masterservergroup.master_server_group.masters[1].id
                  }
                  tsig {
                        tsig_id = edgecast_dns_tsig.tsig1.id
                  }
            }

      }
}