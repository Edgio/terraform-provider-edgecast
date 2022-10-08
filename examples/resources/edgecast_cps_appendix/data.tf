data "edgecast_cps_countrycodes" "bermuda" {
    name = "Bermuda"
}

data "edgecast_cps_countrycodes" "all" {}

data "edgecast_cps_cert_validation_levels" "validation_levels" {}

data "edgecast_cps_domain_statuses" "domain_statuses" {}

data "edgecast_cps_dcv_types" "dcv_types" {}

data "edgecast_cps_validation_statuses" "validation_statuses" {}

data "edgecast_cps_cert_request_cancel_actions" "cancel_request_actions" {}

data "edgecast_cps_cert_order_statuses" "order_statuses" {}

data "edgecast_cps_cert_request_statuses" "cert_statuses" {}