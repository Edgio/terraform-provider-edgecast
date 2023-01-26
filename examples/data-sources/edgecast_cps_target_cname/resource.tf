data "edgecast_cps_target_cname" "mycert_cname" {
  certificate_id       = "12380"
  wait_until_available = true
}

output "target_cname" {
  value = data.edgecast_cps_target_cname.mycert_cname.value
}