// Copyright Edgecast, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package main

import (
	"terraform-provider-ec/ec"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

// Generate docs for website
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return ec.Provider()
		},
	})
}
