output "waf_scopes_id" {
  description = "scopes_id"
  value       = ec_waf_scopes.scopes1.*.id
}