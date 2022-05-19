output "waf_bot_rule_id" {
  description = "bot_rule_id"
  value       = ec_waf_bot_rule.bot_rule_1.*.id
}
