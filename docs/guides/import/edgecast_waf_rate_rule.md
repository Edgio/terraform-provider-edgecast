# Import
## edgecast_waf_rate_rule

To import a resource, first write a resource block for it in your configuration, establishing the name by which it will be known to Terraform:

```terraform
resource "edgecast_waf_rate_rule" "example" {
  
}
```

## Usage
Now terraform import can be run to attach an existing instance to this resource configuration:


```shell
terraform import edgecast_waf_rate_rule.example ACCOUNT_NUMBER:ID   
```
|                 |                                                                   |
|:----------------|-------------------------------------------------------------------|
| ACCOUNT_NUMBER  | The EdgeCast account number the customer user is associated with. |
| ID | The id of the waf rate rule to import.                           | 

As a result of the above command, the resource is recorded in the state file. You can now run terraform plan to see how the configuration compares to the imported resource, and make any adjustments to the configuration to align with the current (or desired) state of the imported object.