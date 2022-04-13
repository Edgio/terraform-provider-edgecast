resource "ec_customer_user" "my_customer_user" {
  account_number = "4FDBB"
  first_name     = "Terraform"
  last_name      = "Test1"
  email          = "terraformuser1@sharedectest.com"
  is_admin       = false # cannot be modified after user creation

  #optional
  address1 = ""
  address2 = ""
  city = ""
  state = ""
  country = ""
  fax = ""
  mobile = ""
  phone = ""
  title = ""
  zip = ""
}
