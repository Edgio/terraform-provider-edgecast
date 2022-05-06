# Import
## edgecast_rules_engine_policy

To import a resource, first write a resource block for it in your configuration, establishing the name by which it will be known to Terraform:

```terraform
resource "edgecast_rules_engine_policy" "example" {
  
}
```

## Usage
Now terraform import can be run to attach an existing instance to this resource configuration:


```shell
terraform import edgecast_rules_engine_policy.example ACCOUNT_NUMBER:PORTAL_TYPE_ID:CUSTOMER_USER_ID:ID   
```
|                  |                                                                   |
|:-----------------|-------------------------------------------------------------------|
| ACCOUNT_NUMBER   | The EdgeCast account number the customer user is associated with. |
| PORTAL_TYPE_ID   | The portal type id of the cname to import.                        | 
| CUSTOMER_USER_ID | The customer user id of the edgecname to import.                  | 
| ID               | The edge cname id to import.                                      | 

As a result of the above command, the resource is recorded in the state file. You can now run terraform plan to see how the configuration compares to the imported resource, and make any adjustments to the configuration to align with the current (or desired) state of the imported object.
