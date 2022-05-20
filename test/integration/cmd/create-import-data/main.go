package main

import (
	"github.com/joho/godotenv"
	"terraform-provider-edgecast/test/integration/cmd/create-import-data/data"
	"terraform-provider-edgecast/test/integration/cmd/create-import-data/internal"
)

func main() {
	internal.CheckError(godotenv.Load())
	data.Create(
		createConfig(),
	)
}
