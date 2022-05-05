output "customer_id" {
  description = "customerId"
  value       = edgecast_customer.test_customer.*.id
}