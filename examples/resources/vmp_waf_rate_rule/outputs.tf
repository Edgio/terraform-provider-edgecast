output "waf_rate_rule_id" {
  description = "rate_rule_id"
  value       = vmp_waf_rate_rule.rate_rule_1.*.id
}