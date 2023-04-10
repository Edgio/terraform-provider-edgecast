
output "certificate_id_dv" {
  description = "certificate_id_dv"
  value       = edgecast_cps_certificate.certificate_2.*.id
}
