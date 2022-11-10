data "edgecast_cps_cert_request_statuses" "cert_statuses" {
}

output "fetched_info_request_status_all" {
  value = data.edgecast_cps_cert_request_statuses.cert_statuses.items
}