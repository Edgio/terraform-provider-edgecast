// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.

package helper

import "strconv"

const (
	int64Base    int = 10
	int64BitSize int = 64
)

func ParseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, int64Base, int64BitSize)
}
