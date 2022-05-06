# Import
## edgecast_dns_group

To import a resource, first write a resource block for it in your configuration, establishing the name by which it will be known to Terraform:

```terraform
resource "edgecast_dns_group" "example" {
  
}
```

## Usage
Now terraform import can be run to attach an existing instance to this resource configuration:


```shell
terraform import edgecast_dns_group.example ACCOUNT_NUMBER:GROUP_PRODUCT_TYPE:ID   
```
|                 |                                                                   |
|:----------------|-------------------------------------------------------------------|
| ACCOUNT_NUMBER  | The EdgeCast account number the customer user is associated with. |
| GROUP_PRODUCT_TYPE  | The grouop product type of the dns group to import.               |
| ID | The dns group id to import.                                       | 

As a result of the above command, the resource is recorded in the state file. You can now run terraform plan to see how the configuration compares to the imported resource, and make any adjustments to the configuration to align with the current (or desired) state of the imported object.
