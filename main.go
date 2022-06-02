// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package main

import (
	"context"
	"flag"
	"log"
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
		err := plugin.Debug(
			context.Background(),
			"github.com/terraform-providers/edgecast",
			opts,
		)
		if err != nil {
			log.Fatal(err.Error())
		}
		return
	}

	plugin.Serve(opts)
}
