output "waf_custom_rule_id" {
  description = "custom_rule_id"
  value       = ec_waf_custom_rule.custom_rule_1.*.id
}