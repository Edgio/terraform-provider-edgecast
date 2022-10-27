# Edgecast Terraform Provider
A Terraform provider for the Edgecast Platform.

## Table of Contents

- [Background](#background)
- [Install](#install)
- [Development](#development)
- [Usage](#usage)
- [Security](#security)
- [Structure](#structure)
- [Contribute](#contribute)
- [License](#license)

## Background

Terraform is a tool for developing, changing and versioning infrastructure 
safely and efficiently. It allows the management of infrastructure as code. 
With Terraform, you can store and version your configuration in GitHub (or your 
source code control system of choice). Thanks to Terraform's configuration 
syntax, there is no need to write custom code to use APIs. Simply describe your
infrastructure in a file and Terraform will figure out the rest.

## Install
This provider is automatically installed when you run `terraform init` in a 
directory that contains a Terraform configuration file that references the 
Edgecast provider.

## Using the Provider
Reference this provider in a Terraform Configuration file (e.g. `main.tf`):

```
terraform {
  required_providers {
    edgecast = {
       version = "0.6.0"
      source  = "Edgio/edgecast"
    }
  }
}
```

Then, use it in a provider block, passing in any credentials provided to you:
```
provider "edgecast" {
  api_token          = "YOUR_API_TOKEN"
  ids_client_secret  = "IDS_SECRET"
  ids_client_id      = "IDS_CLIENT_ID"
  ids_scope          = "IDS_SCOPE"
}
```

Below this, you can start defining resources. For example:
```
resource "edgecast_origin" "origin_images" {
    account_number = "A1234"
    directory_name = "images"
    media_type = "httplarge"
    host_header = "images.mysite.com"
    http {
        load_balancing = "RR"
        hostnames = ["images-origin-1.mysite.com","images-origin-2.mysite.com"]
    }
}
```
Then follow the usual flow for Terraform:
1. Run `terraform init` in a command line or terminal window.
2. Run `terraform plan out="tf.plan"` and inspect the detected changes.
3. Run `terraform apply tf.plan`

## Development
### Requirements
-    [Terraform](https://www.terraform.io/downloads.html) 0.13.x
-    [Terraform] on Mac, move it to /usr/local/bin; Windows add it to your path
-    [Go](https://golang.org/) 1.15
It is strongly recommended to create a GOPATH environment variable that points
to the directory containing your Go installation, then add $GOPATH/bin to your 
$PATH.

### Testing Your Local Code
If you wish to modify the source code and test it with a configuration file,
there are special steps involved.

To simply build the code on any machine, run `go build`.

Actually using your local version provider with Terraform is more complicated.
You must:

1. Install the provider to your machine's Terraform plugin directory.

On a Mac, open a terminal window and change directory to the root folder of
the source code. Run `make install`. The binary will be built and moved to your
local Terraform plugin directory (`~/.terraform.d/plugins`).

On a Windows machine, do the same from a command line or powershell window. 
However, instead of the makefile, run the `install_win.bat` script. The provider
exe will be moved to `%APPDATA%\terraform.d\plugins`.

2. Reference the correct version of the provider in your `.tf` test file.
If you do not do this, Terraform will download the provider from the remote
Terraform Registry instead of using your locally installed provider. 

Make note of the version that is used from within the the `Makefile` or 
`install_win.bat`. You must use the same version within your Terraform 
configuration file. Also, use `"github.com/terraform-providers/edgecast"` as the
source.

Example:
```
terraform {
  required_providers {
    edgecast = {
       version = "0.6.0"
      source  = "Edgio/edgecast"
    }
  }
}
```

## Logging
You can set the `TF_LOG` and `TF_LOG_PATH` environment variables to enable 
logging for Terraform. 
See the [official documentation](https://www.terraform.io/internals/debugging) 
for details.

For example, on MAC OS, running the following two commands will enable logging 
**for your current terminal session**:
```
export TF_LOG=TRACE
export TF_LOG_PATH=/somewhere/on/your/hard_drive/convenient/terraform.log
```

## Usage
The detailed documentation for the provider and specific resources can be found 
on the [Terraform provider registry](https://registry.terraform.io/providers/Edgio/edgecast/latest/docs).

## Security

For those users who have been granted specific permission(s) by an account 
administrator to use the Edgecast Platform, each usage requires the inclusion of 
a user specific token. Tokens can be created or revoked by the user via the 
Portal to ensure token security.
See [Authentication and Authorization](https://developer.edgecast.com/cdn/api/index.html#Introduction/Authentication.htm) 
for more details on the type of token required, and how to acquire a token.


## Structure

```
.
├── edgecast
    package containing edgecast terraform provider resources and functionality 
    to manage and provision edgecast configurations in terraform
│   ├── api
        base client and service specific client files needed for the provider 
        to interact with EdgeCast APIs
        please add new service specific client files here
│   ├── helper
        package containing helper methods
        please add new helper methods here
│   ├── resources
        resource files for individual services
        please add new service specific resource files here
│   └── provider
        edgecast terraform provider
├── examples
    example files to get started managing and provisioning edgecast 
    configurations in terraform
├── docs
    holds detailed documentation for resources and steps 
    to manage and provision edgecast configurations in terraform
├── templates
    holds terraform templates
│   ├── data-sources
│   └── resources
├── tools
    package containing build, terraform-docs and other plug in tools
└── unit-tests
    package containing unit test files
```

## Contribute

Please refer to [the contributing.md file](Contributing.md) for information 
about how to get involved. We welcome issues, questions, and pull requests.

## Maintainers
- Steven Paz: steven.paz@edgecast.com
- Shikha Saluja: shikha.saluja@edgecast.com
- Frank Contreras: frank.contreras@edgecast.com

## License
This project is licensed under the terms of the [Apache 2.0](LICENSE) open 
source license.

## Resources
[CDN Reference Documentation](https://docs.edgecast.com/cdn/index.html) - This 
is a useful resource for learning about the EdgeCast CDN. It is a good starting 
point before using this provider.

[API Documentation](https://docs.edgecast.com/cdn/index.html#REST-API.htm%3FTocPath%3D_____8) - For developers that want to interact directly with the EdgeCast CDN API, refer 
to this documentation. It contains all of the available operations as well as 
their inputs and outputs.

[Examples](https://github.com/EdgeCast/terraform-provider-edgecast/tree/Master/examples) - Examples to get started can be found here.

[Submit an Issue](https://github.com/EdgeCast/terraform-provider-edgecast/issues) - Found a bug? Want to request a feature? Please do so here.