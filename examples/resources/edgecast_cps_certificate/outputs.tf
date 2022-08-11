output "certificate_id" {
  description = "certificate_id"
  value       = edgecast_cps_certificate.certificate_1.*.id
}
