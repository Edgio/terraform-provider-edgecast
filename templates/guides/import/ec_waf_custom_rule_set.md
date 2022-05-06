# Import
## edgecast_waf_custom_rule_set

To import a resource, first write a resource block for it in your configuration, establishing the name by which it will be known to Terraform:

```terraform
resource "edgecast_waf_custom_rule_set" "example" {
  
}
```

## Usage
Now terraform import can be run to attach an existing instance to this resource configuration:


```shell
terraform import edgecast_waf_custom_rule_set.example CUSTOMER_ID:ID   
```
|             |                                                                                                                  |
|:------------|------------------------------------------------------------------------------------------------------------------|
| CUSTOMER_ID | The customer id of the waf custom ruleset to import. |
| ID          | The waf custom ruleset id to import.                                                                             | 

As a result of the above command, the resource is recorded in the state file. You can now run terraform plan to see how the configuration compares to the imported resource, and make any adjustments to the configuration to align with the current (or desired) state of the imported object.
