output "waf_bot_rule_set_id" {
  description = "waf_bot_rule_set_id"
  value       = edgecast_waf_bot_rule_set.bot_rule_set_1.*.id
}
