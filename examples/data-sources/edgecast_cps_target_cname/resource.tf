data "edgecast_cps_target_cname" "cname" {
  certificate_id       = "12380"
  wait_until_available = true
}