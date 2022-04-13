resource "ec_dns_masterservergroup" "master_server_group"{
  account_number = "4FDBB"
  master_server_group_name = "msg99"
  masters {
        name="ns1.test.com"
        ipaddress="10.10.10.1"
        }
  masters {
        name="ns2.test.com"
        ipaddress="10.10.10.2"
        }
  masters {
        name="ns3.test.com"
        ipaddress="10.10.10.3"
        }
  masters {
        name="ns4.test.com"
        ipaddress="10.10.10.4"
        } 
}
