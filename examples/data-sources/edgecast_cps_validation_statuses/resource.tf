data "edgecast_cps_validation_statuses" "validation_statuses" {
}

output "fetched_info_validation_status_all" {
  value = data.edgecast_cps_validation_statuses.validation_statuses.items
}