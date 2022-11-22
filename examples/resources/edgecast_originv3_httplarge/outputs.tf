output "origin_group_id" {
  description = "origin_group_id"
  value       = edgecast_originv3_httplarge.group_1.*.id
}