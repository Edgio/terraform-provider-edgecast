---
page_title: "Authentication"
---

# Authentication
Edgecast supports authentication through REST API tokens and REST API (OAuth 2.0) client credentials. The type of authentication required to build your infrastructure varies by resource. 

!> Do not define sensitive information (e.g., API credentials) within a variable block. One method for securely defining variables is to set variable values from within a **.tfvars** file (e.g., terraform.tfvars) that is excluded from source control. 

[Learn how to protect sensitive input variables.](https://learn.hashicorp.com/tutorials/terraform/sensitive-variables)

## REST API Token
A REST API token is a unique alphanumeric value that identifies the user account through which the requested task will be performed.

[Learn more.](https://developer.edgecast.com/cdn/api/index.html#Introduction/Authentication.htm#RESTAPIToken)

### Resources
Use a REST API token to provision the following types of resources:
* ec_dns_group
* ec_dns_masterservergroup
* ec_dns_secondaryzonegroup
* ec_dns_tsig
* ec_dns_zone
* ec_edgecname
* ec_origin
* ec_waf_access_rule
* ec_waf_custom_rule_set
* ec_waf_managed_rule
* ec_waf_rate_rule
* ec_waf_scopes

### Requirements
A REST API token must meet the following requirements:
* It must be associated with a user that has sufficient privileges to manage the desired resource. 
* It must be granted all HTTP methods.

[Learn how to set up a REST API token.](https://docs.edgecast.com/cdn/#Getting_to_Know_the_Media_Control_Center/Web_Services_REST_API_Token.htm)

### Authentication
You may define a REST API token as a variable within your **main.tf** file as shown in the following excerpt: 

    ...
    variable "token" {
        description = "API key" 
        type = string
        sensitive = true
    }
    
    provider "ec" {
        api_token = var.token
    }

After which, you should define the `token` variable within the **terraform.tfvars** and exclude it from source control by adding it to your **.gitignore** file. The following sample **terraform.tfvars** file sets the `token` variable to a REST API token:

    token = "12345467890abcdefghijklmnopqrst"

## REST API (OAuth 2.0) Client Credentials
A REST API (OAuth 2.0) client consists of a client ID, a secret key, and a scope that authorizes this API client to perform specific actions (e.g., create configurations). 

[Learn more.](https://developer.edgecast.com/cdn/api/index.html#Identity/REST-API-OAuth-Client-Management.htm)

### Resources
Use REST API (OAuth 2.0) client credentials to provision the following type of resource:
* ec_rules_engine_policy

### Scope
Verify that the desired credentials belong to a REST API client that has been granted the scope (e.g., ec.rules) required to manage the desired resource. 

[View root scopes.](https://developer.edgecast.com/cdn/api/index.html#Identity/REST-API-OAuth-Client-Management.htm)

### Authentication
You may define REST API client credentials as variables within your **main.tf** file as shown in the following excerpt: 

    ...
    variable "ids_client_id" {
        description = "REST API client ID" 
        type = string
        sensitive = true
    }
    
    variable "ids_client_secret" {
        description = "REST API client secret" 
        type = string
        sensitive = true
    }
    
    variable "ids_scope" {
        description = "Scope" 
        type = string
        sensitive = true
    }
    
    provider "ec" {
        ids_client_id = var.ids_client_id
        ids_client_secret = var.ids_client_secret
        ids_scope = var.ids_scope
    }

After which, you should define these variable within the **terraform.tfvars** and exclude it from source control by adding it to your **.gitignore** file. The following sample **terraform.tfvars** file sets your REST API client credentials:

    ids_client_id = "12345667-890a-bcde-fghi-jklmnopqrstu"
    ids_client_secret = "abcdefghijklmnopqrstuvwxyz123456"
    ids_scope = "ec.rules"