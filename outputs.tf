output "customer_id" {
  description = "customerId"
  value       = ec_customer.test_customer_02.*.id
}