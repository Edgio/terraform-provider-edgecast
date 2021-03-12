resource "vmp_cname" "cname_images" {
    account_number = "<account_number>"
    name = "<cname_url>"
    type = 3 # HTTPLarge = 3, HTTP Small = 8, ADN = 14
    origin_id = 100
    origin_type = 80 # Customer Origin = 80
}