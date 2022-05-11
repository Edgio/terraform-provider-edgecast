# Import
## edgecast_customer_user

To import a resource, first write a resource block for it in your configuration, establishing the name by which it will be known to Terraform:

```terraform
resource "edgecast_customer_user" "example" {
  
}
```

## Usage
Now terraform import can be run to attach an existing instance to this resource configuration:


```shell
terraform import edgecast_customer_user.example ACCOUNT_NUMBER:ID   
```
|                 |                                                                   |
|:----------------|-------------------------------------------------------------------|
| ACCOUNT_NUMBER  | The EdgeCast account number the customer user is associated with. |
| ID | The customer user id to import.                                   | 

As a result of the above command, the resource is recorded in the state file. You can now run terraform plan to see how the configuration compares to the imported resource, and make any adjustments to the configuration to align with the current (or desired) state of the imported object.
