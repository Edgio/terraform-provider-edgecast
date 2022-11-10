data "edgecast_cps_dcv_types" "dcv_types" {
}

output "fetched_info_dcv_types_all" {
  value = data.edgecast_cps_dcv_types.dcv_types.items
}