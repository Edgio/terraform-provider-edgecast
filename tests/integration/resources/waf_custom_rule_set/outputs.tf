output "waf_custom_rule_set_id" {
  description = "custom_rule_id"
  value       = ec_waf_custom_rule_set.custom_rule_1.*.id
}
