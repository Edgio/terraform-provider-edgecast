# Verizon Media Terraform Provider
A Terraform provider for the Verizon Media Platform.

## Table of Contents

- [Background](#background)
- [Install](#install)
- [Development](#development)
- [Usage](#usage)
- [Security](#security)
- [Contribute](#contribute)
- [License](#license)

## Background

Terraform is a tool for developing, changing and versioning infrastructure safely and efficiently. One important reason people consider Terraform is to manage their infrastructure as code. With Terraform, you can store and version your configuration in GitHub (or your source code control system of choice). Once you learn Terraform's configuration syntax, you don't need to bother learning how to use providers' UIs or APIs—you just tell Terraform what you want and it figures out the rest.

## Install
This provider is automatically installed when you run `terraform init` on a Terraform configuration that contains a reference to the Verizon Media provider.

## Development
### Requirements
-    [Terraform](https://www.terraform.io/downloads.html) 0.13.x
-    [Terraform] on mac, move it to /usr/local/bin in order
-    [Go](https://golang.org/) 1.15 (also set up a GOPATH, as well as add $GOPATH/bin to your $PATH)

### Building The Provider
Follow the instructions to [install it as a plugin.](https://`www.terraform.io/docs/plugins/basics.html#installing-a-plugin). After placing it into your plugins directory, run `terraform init` to initialize it.

## Logging
- export TF_LOG=TRACE
- export TF_LOG_PATH=/somewhere/on/your/hard_drive/convenient/terraform.log

## Usage
Running `terraform init` will automatically download and install the plugin for your use as long as your terraform configuration references the provider like so:

```terraform
# You must provide either an API Token or IDS Credentials, but not both!
provider "vmp" {
	# your API token provided by Verizon Media
	api_token = "xxx"
	
	# OR your IDS credentials provided by Verizon Media
	ids_client_secret = "xxx"
	ids_client_id = "xxx"
	ids_scope = "scope1 scope2"
}
```
**Note for Verizon Media internal users:** you must also specify `partner_id` and `partner_user_id` in addition to your credentials.

There are two ways to set your data. One is to set your configuration data directly in the main.tf file.
For local testing, this works fine. **However, if you have multi-environments, we would recommend to use terraform.tfvars file.
This file contains all variables that you need in order to run the terraform. Please refer to the [Using Variable Files](#Using%20Variable%20Files) section.** 

### Provision MCC Accounts
Please use the Verizon Media API for retrieving specific IDs available for Services, Access Modules, and Delivery Regions.
A future version of this provider may provide Terraform data sources for these.

```terraform
resource "vmp_customer" "my_customer" {
  company_name = "My Company Name"
  service_level_code = "STND"
  services = [1,2,3]
  delivery_region =  1
  access_modules = [1,2,3]
}
```

### Provision MCC Users
```terraform
resource "vmp_customer_user" "test_customer_admin" {
  account_number = vmp_customer.my_customer.id
  first_name = "Admin"
  last_name = "User"
  email = "admin@mycompany.com"
  is_admin = true
}
```

### Configure Origins
```terraform
resource "vmp_origin" "images" {
    directory_name = "images3"
    host_header = "dev-images.mysite.com"
    http {
        load_balancing = "RR"
        hostnames = ["http://dev-images-customer-origin.mysite.com", "http://dev-images-customer-origin2.mysite.com"]
    }
}
```
### Configure CNAMEs
```terraform
resource "vmp_cname" "images" {
    name ="dev-images-customer-origin.mysite.com"
    type = 3
    origin_id = "${vmp_origin.images.id}"
    origin_type = 80
}
```

### Variable File Usage
```terraform
partner_info = {
    #for pointing to staging environment, leave null to default to production
    api_address = null
    # You must provide either an API Token or IDS credentials, but not both
    api_token = null
    ids_client_secret = null
    ids_client_id = null
    ids_scope = null
}

new_customer_info = {
    company_name = "Terraform test customer demo15" #Customer Name
    service_level_code = "STND" #Service Type
    #all available services=> 1:HTTP Large Object,2:HTTPS Large Object,3:HTTP Small Object,4:HTTPS Small Object,6:Windows,7:Advanced Reports,8:Real-Time Stats,9:Token Auth,10:Edge Performance Analytics,15:Origin Storage,16:RSYNC,19:ADN,20:Download Manager,21:ADNS,22:Dedicated Hosting,23:Edge Optimizer,25:DNS Route,26:DNS Zones,29:DNS Health Checks,31:Bandwidth By Report Code,32:DNS-Standard,33:DNS-Adaptive,34:DNS-APR,38:WAF,39:Analysis Engine,40:HTTP Rate Limiting,41:Basic Rules v4.0,42:Advanced Rules v4.0,43:Mobile Device Detection Rules v4.0,44:Rules Engine v4.0,47:Translate,48:Dynamic Cloud Packaging,49:Encrypted HLS,50:Origin Shield,51:Reports and Logs,52:Log Delivery,54:SSA,56:Encrypted Key Rotation,57:Real-Time Log Delivery,58:Report Builder,59:Dynamic Imaging,60:China Delivery,61:WAF Essential,62:Report Builder Users,63:Report Builder Rows,64:Report Builder Reports,65:Edge Functions,66:Certificate Provisioning,67:Edge-Insights,68:Edge Image Optimizer,69:Url Redirects,70:Azure Cloud Storage
    services = [1,9,15,19]
    #available access modules => 1:Storage, 4:Analytics, 5:Admin, 7:Customer Origin, 8:Purge/Load, 21:Users, 22:Company, 25:Country Filtering, 26:Token Auth, 27:Dashboard, 29:HTTP Large, 30:Edge CNAMEs, 32:Core Reports, 40:Token Auth, 46:Token Auth, 53:Cache Settings, 56:HTTP Large Object, 71:HTTP Streaming, 72:ADN, 73:Customer Origin, 74:Purge/Load, 75:Token Auth, 76:Country Filtering, 77:Edge CNAMEs, 78:Cache Settings, 79:Application Delivery Network, 81:Tools, 138:Query-String Caching, 139:Query-String Logging, 140:Compression, 144:Query-String Caching, 145:Query-String Logging, 146:Compression, 149:Smooth Streaming Player, 153:JW Player, 157:Raw Log Settings, 159:Traffic Summary, 160:Bandwidth, 161:Data Transferred, 162:Hits, 163:Cache Statuses, 164:Cache Hit Ratio, 166:CDN Storage, 168:Notes, 169:HTTP Large, 170:HTTPS Large, 171:HTTP Small, 172:HTTPS Small, 174:Flash, 175:ADN, 176:ADN SSL, 177:HTTP Large, 178:HTTPS Large, 179:HTTP Small, 180:HTTPS Small, 182:Flash, 183:ADN, 184:ADN SSL, 185:All Platforms, 186:HTTP Large, 187:HTTP Small, 189:Flash, 190:ADN, 191:All Platforms, 192:HTTP Large, 193:HTTP Small, 194:ADN, 195:All Platforms, 196:HTTP Large, 197:HTTP Small, 198:ADN, 204:Usage, 386:IPv4/IPv6, 387:Data Transferred, 409:Custom Reports, 410:Edge CNAMEs, 411:Notes, 412:All Platforms, 413:HTTP Large, 414:HTTP Small, 415:Flash, 416:ADN, 479:Token Generator, 501:Add Users, 502:Edit Users
    access_modules = [1,4,5,7,8,21,22,25,26,27,29,30,32,40,46,53,56,71,72,73,74,75,76,77,78,79,81,138,139,140,144,145,146,149,153,157,159,160,161,162,163,164,166,168,169,170,171,172,174,175,176,177,178,179,180,182,183,184,185,186,187,189,190,191,192,193,194,195,196,197,198,204,386,387,409,410,411,412,413,414,415,416,479,501,502]
    delivery_region = 1 # Delivery Region
}

new_admin_user = {
    first_name = "admin"
    last_name = "user"
    email = "adminuser15@test.com" #Please use correct email. If a temp email is provided, you would get errors.
    is_admin = true
}

origin_info = {
    origins = ["http://dev-images-customer-origin.mysite.com","http://dev-images-customer-origin2.mysite.com"] # one or more origins
    load_balancing = "RR" # load balance type
    host_header = "dev-images.mysite.com" 
    directory_name = "image14"
}

cname_info =  {
    cname = "dev-images-customer-origin15.edgec4tz.com"
    type = 3 #cname type
    origin_type = 80 #origin type
}
```


## Security

For those users who have been granted specific permission(s) by an account administrator to use the Verizon Media Platform, each usage requires the inclusion of a user specific token. Tokens can be created or revoked by the user via the Portal to ensure token security.

## Contribute

Please refer to [the contributing.md file](Contributing.md) for information about how to get involved. We welcome issues, questions, and pull requests.

## Maintainers
- Changgyu Oh: changgyu.oh@verizonmedia.com
- Steven Paz: steven.paz@verizonmedia.com

## License
This project is licensed under the terms of the [Apache 2.0](LICENSE) open source license.