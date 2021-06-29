output "waf_access_rule_id" {
  description = "access_rule_id"
  value       = vmp_waf_access_rule.access_rule_1.*.id
}