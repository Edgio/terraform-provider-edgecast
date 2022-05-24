package main

import (
	"github.com/joho/godotenv"
	"terraform-provider-edgecast/test/integration/cmd/populate/data"
)

func main() {
	_ = godotenv.Load()
	data.Fix(
		createConfig(),
	)
}
