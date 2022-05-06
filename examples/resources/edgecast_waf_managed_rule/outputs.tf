output "waf_manged_rule_id" {
  description = "managed_rule_id"
  value       = edgecast_waf_managed_rule.managed_rule_1.*.id
}