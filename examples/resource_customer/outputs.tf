output "customer_id" {
  description = "customerId"
  value       = vmp_customer.test_customer.*.id
}