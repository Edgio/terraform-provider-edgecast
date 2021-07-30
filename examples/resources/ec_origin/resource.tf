resource "ec_origin" "origin_images" {
    account_number = "A1234"
    directory_name = "images"
    media_type = "httplarge"
    host_header = "images.mysite.com"
    http {
        load_balancing = "RR"
        hostnames = ["images-origin-1.mysite.com","images-origin-2.mysite.com"]
    }
}