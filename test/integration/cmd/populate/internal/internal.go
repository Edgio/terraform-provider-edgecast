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

func Check[T any](data T, err error) T {
	CheckError(err)
	return data
}

func Pointer[T any](t T) *T {
	return &t
}
