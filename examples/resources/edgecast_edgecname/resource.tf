resource "edgecast_edgecname" "edgecname_cdn_origin" {
    account_number = "0001"
    name = "cdn.example.com"
    dir_path = "/marketing"
    enable_custom_reports = 1 # True = 1, False = 0
    media_type_id = 3 # HTTPLarge = 3, HTTP Small = 8, ADN = 14
    origin_id = -1
}

resource "edgecast_edgecname" "edgecname_customer_origin" {
    account_number = "0001"
    name = "sales.example.com"
    dir_path = "/sales"
    enable_custom_reports = 1 # True = 1, False = 0
    media_type_id = 3 # HTTPLarge = 3, HTTP Small = 8, ADN = 14
    origin_id = 100
}