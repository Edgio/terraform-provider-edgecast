// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package waf

const (
	errorIntExpand        string = "%v was not an int, actual: %T"
	errorBoolExpand       string = "%v was not a bool, actual: %T"
	errorStringExpand     string = "%v was not a string, actual: %T"
	errorInterfacesExpand string = "%v was not a []interface{}, actual: %T"
)
