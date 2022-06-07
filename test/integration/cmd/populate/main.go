// Copyright 2021 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.
package main

import (
	"github.com/joho/godotenv"
	"terraform-provider-edgecast/test/integration/cmd/populate/data"
)

func main() {
	_ = godotenv.Load()
	data.Create(
		createConfig(),
	)
}
