package internal

import (
	"fmt"
	"os"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Check(data any, err error) any {
	CheckError(err)
	return data
}
