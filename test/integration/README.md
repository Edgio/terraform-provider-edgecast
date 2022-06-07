# Integration Tests

## Configuration
You can specify the environment you wish to run the tests against by specifying configuration values in the `.env` in the root of this project. 


```bash
  ACCOUNT_NUMBER      = ""
  
  API_ADDRESS         = ""
  API_ADDRESS_LEGACY  = ""
  API_TOKEN           = ""
  
  IDS_ADDRESS         = ""
  IDS_CLIENT_ID       = ""
  IDS_CLIENT_SECRET   = ""
  IDS_SCOPE           = ""
  
  # SUPPLYING THESE VALUES ENABLE IMPORT TESTING ON THE RESOURCE
  CUSTOMER_USER_IMPORT_ID             = ""
  DNS_MASTER_ZONE_GROUP_IMPORT_ID     = ""
  DNS_SECONDARY_ZONE_GROUP_IMPORT_ID  = ""
  DNS_GROUP_IMPORT_ID                 = ""
  DNS_TSIG_IMPORT_ID                  = "" 
  DNS_ZONE_IMPORT_ID                  = ""
  ORIGIN_IMPORT_ID                    = ""
  RULES_ENGINE_IMPORT_ID              = ""
  WAF_ACCESS_RULE_IMPORT_ID           = "" 
  WAF_BOT_RULE_SET_ID                 = ""
  WAF_CUSTOM_RULE_IMPORT_ID           = ""
  WAF_MANAGED_RULE_IMPORT_ID          = ""
  WAF_RATE_RULE_IMPORT_ID             = ""
  WAF_SCOPES_IMPORT_ID                = ""
```

## Usage
Integration tests are run with the following command in the root directory of the project:

```
task integration_test:run
```

This will run all the tests within the integration test suite.

## Creating Tests
These are the steps for adding a new test to the integration test suite.


#### Creating the test directory

Create a directory for the resource you wish to test within the `/test/integration/resources` directory. 

```shell
$ ls -l resources
total 0
drwxr-xr-x   6 user  staff  192 Apr  1 00:00 customer
drwxr-xr-x   8 user  staff  256 Apr  1 00:00 customer_user
drwxr-xr-x   8 user  staff  256 Apr  1 00:00 dns_group
drwxr-xr-x   9 user  staff  288 Apr  1 00:00 dns_masterservergroup
drwxr-xr-x   9 user  staff  288 Apr  1 00:00 dns_secondaryzonegroup
drwxr-xr-x   8 user  staff  256 Apr  1 00:00 dns_tsig
drwxr-xr-x   9 user  staff  288 Apr  1 00:00 dns_zone
drwxr-xr-x   8 user  staff  256 Apr  1 00:00 edgecname
drwxr-xr-x   8 user  staff  256 Apr  1 00:00 origin
drwxr-xr-x  10 user  staff  320 Apr  1 00:00 rules_engine_policy
drwxr-xr-x   9 user  staff  288 Apr  1 00:00 waf_access_rule
drwxr-xr-x   9 user  staff  288 Apr  1 00:00 waf_custom_rule_set
drwxr-xr-x   9 user  staff  288 Apr  1 00:00 waf_managed_rule
drwxr-xr-x   9 user  staff  288 Apr  1 00:00 waf_rate_rule
drwxr-xr-x   9 user  staff  288 Apr  1 00:00 waf_scopes

```
_a listing of integration tests at the time this documentation was created._


### Creating the tests
Within the newly-created directory, perform the following steps:

   1. **Core Files**: Create main file that will be used for all steps within the suite. This will specify the exact provider and version you wish to run tests against, for example:
   ```terraform
    terraform {
      required_providers {
        edgecast = {
          version = "0.5.1"
          source  = "EdgeCast/edgecast"
        }
      }
    }
    
    ##########################################
    # Variables
    ##########################################
    variable "credentials" {
      type = object({
        api_token         = string
        ids_client_secret = string
        ids_client_id     = string
        ids_scope         = string
      })
    }
    
    ##########################################
    # Providers
    ##########################################
    provider "edgecast" {
      api_token         = var.credentials.api_token
      ids_client_secret = var.credentials.ids_client_secret
      ids_client_id     = var.credentials.ids_client_id
      ids_scope         = var.credentials.ids_scope
    }

   ```
_An example of a typical `main.tf`, you can find more examples within the `/examples` directory._ 

If applicable include an `outputs.tf` file declaring any expected output values, for example:
```terraform
    output "waf_scopes_id" {
      description = "scopes_id"
      value       = edgecast_waf_scopes.scopes1.*.id
    }
```

   2. **Include Create Step**: create a `create.tf.step` file. This should contain the initial state of the resource, for example:
```terraform
    resource "edgecast_dns_tsig" "tsig1" {
      account_number = "A1234"
      alias = "Test terraform keys"
      key_name = "key1"
      key_value = "HFNASHDJJKQWHKJ1234"
      algorithm_name = "HMAC-SHA512"
    }
```
_A `create.tf.step` file for `dns_tsig`._


3. **Include Update Step**: create an `update.tf.step` file. This should contain an updated state of the resource, for example:
```terraform
    resource "edgecast_dns_tsig" "tsig1" {
      account_number = "A1234"
      alias = "updated"
      key_name = "key1"
      key_value = "updated"
      algorithm_name = "HMAC-SHA512"
    }
```
_An `updated.tf.step` file for `dns_tsig`._

4. **Include test into the integration test suite**: Modify the `Taskfile.yaml` to include the newly created test directory, for example:
```yaml
includes:
  dns_tsig:
    taskfile: Integration.Taskfile.yaml
    dir: ./resources/dns_tsig
```
Then include the test in the suite by modifying the `run` task, for example:
```yaml
tasks:
  run:
    desc: "run will execute the entire integration test suite"
    cmds:
      - task: dns_tsig:default
```
And now, the new test will be included in the run whenever the' run' task is executed.


### Populating Test Data 

Integration test data can be generated by running the following command.

```shell
task integration_test:populate
```

This will create test data for all resources within terraform-provider-edgecast using the `ACCOUNT_NUMBER` specified in your `.env` file.