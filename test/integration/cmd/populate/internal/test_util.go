// Copyright 2023 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.
package internal

import (
	"fmt"
	"time"
)

// PrependTimeStamp adds a timestamp to the beginning of a string.
func Unique(s string) string {
	n := fmt.Sprintf("%d", time.Now().Unix())
	n = n[len(n)-4:]
	return fmt.Sprintf("%s%s", n, s)
}
