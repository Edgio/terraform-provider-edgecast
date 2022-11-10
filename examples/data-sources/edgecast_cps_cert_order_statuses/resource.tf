data "edgecast_cps_cert_order_statuses" "order_statuses" {
}

output "fetched_info_cert_order_status_all" {
  value = data.edgecast_cps_cert_order_statuses.order_statuses.items
}