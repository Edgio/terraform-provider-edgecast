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
* edgecast_dns_group
* edgecast_dns_masterservergroup
* edgecast_dns_secondaryzonegroup
* edgecast_dns_tsig
* edgecast_dns_zone
* edgecast_edgecname
* edgecast_origin
* edgecast_waf_access_rule
* edgecast_waf_custom_rule_set
* edgecast_waf_managed_rule
* edgecast_waf_rate_rule
* edgecast_waf_scopes


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
    
    provider "edgecast" {
        api_token = var.token
    }

After which, you should define the `token` variable within the **terraform.tfvars** and exclude it from source control by adding it to your **.gitignore** file. The following sample **terraform.tfvars** file sets the `token` variable to a REST API token:

    token = "12345467890abcdefghijklmnopqrst"

## REST API (OAuth 2.0) Client Credentials
A REST API (OAuth 2.0) client consists of a client ID, a secret key, and a scope that authorizes this API client to perform specific actions (e.g., create configurations). 

[Learn more.](https://developer.edgecast.com/cdn/api/index.html#Identity/REST-API-OAuth-Client-Management.htm)

### Resources
Use REST API (OAuth 2.0) client credentials when working with the following resources and data sources:
* edgecast_rules_engine_policy
* edgecast_cps_certificate
* edgecast_cps_cert_order_statuses
* edgecast_cps_cert_request_cancel_actions
* edgecast_cps_cert_request_statuses
* edgecast_cps_cert_validation_levels
* edgecast_cps_countrycodes
* edgecast_cps_dcv_types
* edgecast_cps_domain_statuses
* edgecast_cps_validation_statuses

### Scope
Verify that the desired credentials belong to a REST API client that has been granted the scope (e.g., `ec.rules` and `sec.cps.certificates`) required to manage the desired resource. 

[View root scopes.](https://developer.edgecast.com/cdn/api/index.html#Identity/REST-API-OAuth-Client-Management.htm)

!> Our Terraform implementation only supports a single scope per REST API client. If you plan on managing both `edgecast_rules_engine_policy` and `edgecast_cps_certificate` resources, then you will need to use a separate working directory for each type of resource.

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
    
    provider "edgecast" {
        ids_client_id = var.ids_client_id
        ids_client_secret = var.ids_client_secret
        ids_scope = var.ids_scope
    }

After which, you should define these variable within the **terraform.tfvars** and exclude it from source control by adding it to your **.gitignore** file. The following sample **terraform.tfvars** file sets your REST API client credentials:

    ids_client_id = "12345667-890a-bcde-fghi-jklmnopqrstu"
    ids_client_secret = "abcdefghijklmnopqrstuvwxyz123456"
    ids_scope = "ec.rules"
