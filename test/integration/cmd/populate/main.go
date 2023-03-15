// Copyright 2023 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.
package main

import (
	"terraform-provider-edgecast/test/integration/cmd/populate/config"
	"terraform-provider-edgecast/test/integration/cmd/populate/data"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	dp := data.DataPopulator{
		Config: config.NewConfig(),
	}

	dp.Populate()
}
