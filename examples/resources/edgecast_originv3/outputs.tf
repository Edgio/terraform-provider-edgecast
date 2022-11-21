output "origin_group_id" {
  description = "origin_group_id"
  value       = edgecast_originv3_group.group_httplarge_1.*.id
}