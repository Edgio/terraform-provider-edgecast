// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.

package helper

import (
	"fmt"
	"strconv"
)

const (
	int64Base    int = 10
	int64BitSize int = 64
)

func ParseInt64(s string) (int64, error) {
	i, err := strconv.ParseInt(s, int64Base, int64BitSize)
	if err != nil {
		return i, fmt.Errorf("error parsing int64: %w", err)
	}

	return i, nil
}
