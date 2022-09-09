// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package main

import (
	"flag"
	"terraform-provider-edgecast/edgecast"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

// Generate docs for website
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	var debugMode bool

	flag.BoolVar(
		&debugMode,
		"debug",
		false,
		"set to true to run the provider with support for debuggers",
	)
	flag.Parse()

	opts := &plugin.ServeOpts{ProviderFunc: edgecast.Provider}

	if debugMode {
		opts.Debug = true
		opts.ProviderAddr = "github.com/terraform-providers/edgecast"
	}

	plugin.Serve(opts)
}
