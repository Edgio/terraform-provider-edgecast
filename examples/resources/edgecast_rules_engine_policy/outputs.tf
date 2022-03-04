output "policy_id" {
  description = "policyid"
  value       = edgecast_rules_engine_policy.httplarge_policy.*.id
}