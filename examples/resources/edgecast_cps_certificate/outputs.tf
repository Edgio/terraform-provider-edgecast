output "certificate_id_ev" {
  description = "certificate_id_ev"
  value       = edgecast_cps_certificate.certificate_1.*.id
}

output "certificate_id_dv" {
  description = "certificate_id_dv"
  value       = edgecast_cps_certificate.certificate_2.*.id
}
