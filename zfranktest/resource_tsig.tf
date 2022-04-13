resource "ec_dns_tsig" "tsig1" {
  account_number = "4FDBB"
  alias = "Test terraform key 1"
  key_name = "key1"
  key_value = "HFNASHDJJKQWHKJ1234"
  algorithm_name = "HMAC-SHA512"
}
