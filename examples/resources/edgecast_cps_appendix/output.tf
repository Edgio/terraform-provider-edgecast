
output "fetched_info_countrycode_Bermuda" {
  
  value = {
    for output in data.edgecast_cps_countrycodes.countrycodes.items : output.country =>
    output.two_letter_code if contains(["Bermuda"], output.country)
  }
}

output "fetched_info_countrycode_all" {
  value = data.edgecast_cps_countrycodes.countrycodes.items
}

output "fetched_info_countrycode_none" {
  value = {
    for output in data.edgecast_cps_countrycodes.countrycodes.items : output.country =>
    output.two_letter_code if contains(["random"], output.country)
  }
}

output "fetched_info_cert_order_status_all" {
  value = data.edgecast_cps_cert_order_statuses.order_statuses.items
}

output "fetched_info_cert_request_cancel_actions_all" {
  value = data.edgecast_cps_cert_request_cancel_actions.cancel_request_actions.items
}

output "fetched_info_request_status_all" {
  value = data.edgecast_cps_cert_request_statuses.cert_statuses.items
}

output "fetched_info_cert_validation_levels_all" {
  value = data.edgecast_cps_cert_validation_levels.validation_levels.items
}

output "fetched_info_dcv_types_all" {
  value = data.edgecast_cps_dcv_types.dcv_types.items
}

output "fetched_info_domain_status_all" {
  value = data.edgecast_cps_domain_statuses.domain_statuses.items
}

output "fetched_info_validation_status_all" {
  value = data.edgecast_cps_validation_statuses.validation_statuses.items
}



