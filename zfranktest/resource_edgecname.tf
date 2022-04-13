resource "ec_edgecname" "edgecname_edgecast_origin_cdn" {
    account_number = "4FDBB"
    name = "terrform.cdnorigin.com"
    dir_path = "/ec/origin/path2"
    enable_custom_reports = 1 # True = 1, False = 0
    media_type_id = 3 # HTTPLarge = 3, HTTP Small = 8, ADN = 14
    origin_id = -1
}

resource "ec_edgecname" "edgecname_customer_origin_cust" {
    account_number = "4FDBB"
    name = "terraform.customerorigin.com"
    dir_path = "/origin/path/to/content2"
    enable_custom_reports = 1 # True = 1, False = 0
    media_type_id = 3 # HTTPLarge = 3, HTTP Small = 8, ADN = 14
    origin_id = 3169027
}