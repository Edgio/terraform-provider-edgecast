output "customer_id" {
  description = "customerId"
  value       = edgecast_customer.test_customer_01.*.id
}