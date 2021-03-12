resource "vmp_customer_user" "my_customer_user" {
  account_number = "A1234"
  first_name     = "Admin_First_Name"
  last_name      = "Admin_Last_Name"
  email          = "Admin@mysite.com"
  is_admin       = true # cannot be modified after user creation

  #optional
  address1 = ""
  address2 = ""
  city = ""
  country = ""
  fax = ""
  mobile = ""
  phone = ""
  time_zone_id = ""
  title = ""
  zip = ""
}
