# Edgecast Terraform Provider
A Terraform provider for the Edgecast Platform.

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
This provider is automatically installed when you run `terraform init` on a Terraform configuration that contains a reference to the Edgecast provider.

## Development
### Requirements
-    [Terraform](https://www.terraform.io/downloads.html) 0.13.x
-    [Terraform] on mac, move it to /usr/local/bin in order
-    [Go](https://golang.org/) 1.15 (also set up a GOPATH, as well as add $GOPATH/bin to your $PATH)

### Building The Provider
Follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin). After placing it into your plugins directory, run `terraform init` to initialize it.

## Logging
You can set the `TF_LOG` and `TF_LOG_PATH` environment variables to enable logging for Terraform. See the [official documentation](https://www.terraform.io/docs/internals/debugging.html) for details.

For example, on MAC OS, running the following two commands will enable logging **for your current terminal session**:
```
export TF_LOG=TRACE
export TF_LOG_PATH=/somewhere/on/your/hard_drive/convenient/terraform.log
```

## Usage
The detailed documentation for the provider and specific resources can be found on the [Terraform provider registry](https://registry.terraform.io/providers/EdgeCast/ec/latest/docs).

## Security

For those users who have been granted specific permission(s) by an account administrator to use the Edgecast Platform, each usage requires the inclusion of a user specific token. Tokens can be created or revoked by the user via the Portal to ensure token security.

## Contribute

Please refer to [the contributing.md file](Contributing.md) for information about how to get involved. We welcome issues, questions, and pull requests.

## Maintainers
- Changgyu Oh: changgyu.oh@edgecast.com
- Steven Paz: steven.paz@edgecast.com

## License
This project is licensed under the terms of the [Apache 2.0](LICENSE) open source license.