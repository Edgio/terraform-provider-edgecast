data "edgecast_cps_cert_request_cancel_actions" "cancel_request_actions" {
}

output "fetched_info_cert_request_cancel_actions_all" {
  value = data.edgecast_cps_cert_request_cancel_actions.cancel_request_actions.items
}