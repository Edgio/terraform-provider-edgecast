# Onboarding: terraform-provider-edgecast

## Purpose
This document is a collection of all onboarding-related information, tips & 
tricks, etc. for first-time contributors.

## Learning Terraform as a User
Hashicorp provides extensive documentation and learning guides for learning how 
to use Terraform at [Hashicorp Learn](https://learn.hashicorp.com/terraform). 

Once you are familiar with how to use Terraform, you can learn how to develop a 
Terraform provider yourself.

## Learning Go
Before you can implement a Terraform Provider, you must learn Go. The internet 
contains a plethora of learning resources for the Go language. Below are sources 
that we find useful:
- [Go Get Started Guide](https://go.dev/learn/)
    - Also lists additional learning resources.
- [Pluralsight (License Required)](https://www.pluralsight.com/)
    - We recommend the [Getting Started With Go](https://app.pluralsight.com/library/courses/getting-started-with-go/) course for an excellent guide to setting up your 
    local development environment.
- [Go by Example](https://gobyexample.com/)
    - Ideal for those who learn best by doing.
- [The Go Playground](https://go.dev/play/)
    - Sandbox for playing with Go code.
- [Go Docs](https://golang.org/doc/)
    - Landing page for many Go resources e.g. learning, concepts, etc.
- [Effective Go](https://go.dev/doc/effective_go)
    - Outlines community-accepted code style and best practices.
- [Our Go Guide](https://github.com/EdgeCast/ec-sdk-go/blob/main/Go.md)
    - Best practices and code style our team adheres to.

## Developing for Terraform
Once you have familiarized yourself with Go, you can learn how to develop a
custom provider. You can find the Hashicorp Learn tutorial 
[here](https://learn.hashicorp.com/collections/terraform/providers).

## Local Environment Setup
If you have a Pluralsight license, we recommend watching the Getting Started 
With Go course listed above in the Learning Resources.

Below is a breakdown of the setup:
1. Install Go - [download](https://go.dev/dl/)
2. Install Git - [download](https://git-scm.com/downloads)
3. Install Visual Studio Code - [download](https://code.visualstudio.com/download)
4. Set up Git access for your machine
    - Use a Git client like [Sourcetree](https://www.sourcetreeapp.com/) or
    [GitHub Desktop](https://desktop.github.com/)
    - [Create an SSH Key](https://docs.github.com/en/authentication/connecting-to-github-with-ssh) 
    for your GitHub account and add it to your local machine.
5. Install Visual Studio Code Extension
    - In VS Code, open the Extensions menu and search for "Go". The developer 
    should be "Go Team at Google".
6. Install the HashiCorp Terraform Extension.
    - In VS Code, open the Extensions menu and search for "Go". The developer 
    should be "Hashicorp".
7. Install any suggested Go and Terraform tools as prompted by Visual Studio
8. Clone the [terraform-provider-edgecast](https://github.com/EdgeCast/terraform-provider-edgecast) 
repository.
8. Refer to the [README.md](README.md), [Contributing.md](Contributing.md), and 
[Architecture.md](Architecture.md) files.

## Development Workflow
Create a branch off of main and begin coding!

### Running Your Local Code
Refer to the "Testing Your Local Code" section of the [readme](README.md).

### Debugging
Debugging Terraform code is a bit more complicated than just clicking Run and 
Debug in VS Code. 
Please carefully read through 
[Hashicorp's article on debugging](https://www.terraform.io/plugin/debugging#visual-studio-code)
as these instructions may change.

### Default Environment Configuration
Please note that the provider is configured to point to the EdgeCast production 
environment. This is fine if you own a test account. Developers employed at 
EdgeCast may wish to point to a different environment. You can do so globally by 
modifying the URLs in the provider configuration.

For example:
```terraform
provider "edgecast" {
    api_token = "MY_API_TOKEN"
    ids_client_secret = "MY_IDS_SECRET"
    ids_client_id = "MY_IDS_ID"
    ids_scope = "MY_IDS_SCOPE"
    api_address = "https://api.vdms.io"
    api_address_legacy = "https://api.edgecast.com"
    ids_address = "https://id.vdms.io"
}
```

### Testing
#### Unit Testing
Ensure that all unit tests pass before submitting a pull request. 

1. Open a terminal and navigate to the `edgecast` directory.
2. Run `go test ./…`

Please create or modify unit tests when modifying or adding to any of the code 
in `edgecast`. 

Create functions for "flattening" and "expanding" data and write unit 
tests for them. Flattening refers to the conversion of data returned from the 
EdgeCast SDK into a form that Terraform can consume. Expanding is the opposite - 
it refers to the conversion of raw Terraform state into a form that the EdgeCast 
SDK can consume. 

Refer to `edgecast/resources/origin/origin_test.go` for a test example. There are 
helper functions in `edgecast/helper` that can assist in writing tests.

```go
// origin.go
func expandHostname(attr interface{}) (*[]origin.Hostname, error) {
    // ...
}

// origin_test.go
func TestExpandHostname(t *testing.T) {
    // arrange
    input := helper.NewTerraformSet(...)

    // act
    actual, err := expandHostname(input)

    // assert...
}

```

#### Regression Testing
Consider the scope of changes in your PR and whether it is necessary to run some 
or all of the example files located in the [examples](examples) folder as 
end-to-end tests to identify regressions.

### Release
Create a new release in GitHub with the appropriate vM.m.r semantic version e.g. 
v0.1.8 via the releases tab. Ensure to create a tag on the release screen using 
the same name. This process will be replaced with a GitHub action in the future.

In your release, be sure to include each section below. If there are no changes
for a section, omit it. Refer to existing 
[releases](https://github.com/EdgeCast/terraform-provider-edgecast/releases).
- Breaking Changes
    - Alert consumers of the provider of any changes that can break their code.
- New Features
- Bug Fixes and Enhancements
    - Enhancements can include performance improvements, code optimization, etc.

Once your release is complete, verify that the new version has been pushed to
the Terraform Registry for users to consume. You can do so by checking the 
version [here](https://registry.terraform.io/providers/EdgeCast/edgecast).