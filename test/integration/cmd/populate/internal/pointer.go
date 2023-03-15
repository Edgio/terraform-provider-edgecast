// Copyright 2023 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.
package internal

// Pointer wraps any instance in a pointer.
func Pointer[T any](t T) *T {
	return &t
}
