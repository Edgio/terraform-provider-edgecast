data "edgecast_cps_cert_validation_levels" "validation_levels" {
}

output "fetched_info_cert_validation_levels_all" {
  value = data.edgecast_cps_cert_validation_levels.validation_levels.items
}