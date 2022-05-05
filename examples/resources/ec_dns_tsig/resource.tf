resource "edgecast_dns_tsig" "tsig1" {
  account_number = "A1234"
  alias = "Test terraform keys"
  key_name = "key1"
  key_value = "HFNASHDJJKQWHKJ1234"
  algorithm_name = "HMAC-SHA512"
}
