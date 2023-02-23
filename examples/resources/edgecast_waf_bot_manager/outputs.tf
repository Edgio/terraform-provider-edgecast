output "botmanager_id" {
  description = "botmanager_id"
  value       = edgecast_waf_botmanager.botmanager_1.*.id
}
