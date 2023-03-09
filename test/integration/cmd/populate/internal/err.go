// Copyright 2023 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.
package internal

import (
	"fmt"
	"os"
)

// CheckError will inspect an error and, if it exists, prints its details and
// immediately exists the application.
func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Check will inspect data and the error created at the time of the data's
// retrieval. If the error exists, it prints the details and immediately exists
// the application.
func Check[T any](data T, err error) T {
	CheckError(err)
	return data
}

// UnChecked returns data and discards the error created at the time of the
// data's retrieval.
func UnChecked[T any](data T, _ error) T {
	return data
}
