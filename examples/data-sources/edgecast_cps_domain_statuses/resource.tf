data "edgecast_cps_domain_statuses" "domain_statuses" {
}

output "fetched_info_domain_status_all" {
  value = data.edgecast_cps_domain_statuses.domain_statuses.items
}