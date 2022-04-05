output "policy_id" {
  description = "policyid"
  value       = ec_rules_engine_policy.my_policy.*.id
}