output "policy_id" {
  description = "policyid"
  value       = edgecast_rules_engine_policy.my_policy.*.id
}