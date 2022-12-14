
data "edgecast_originv3_httplarge_origin_shield_pops" "all" {
}

data "edgecast_originv3_httplarge_origin_shield_pops" "pop" {
    code = "BLR"
}

data "edgecast_originv3_httplarge_origin_shield_pops" "region" {
    code = "BYAP"
}