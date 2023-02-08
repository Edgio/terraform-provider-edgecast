// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.
package internal

import (
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// ValidateDuration validates both numeric string duration values for terraform
// properties. e.g. '6m', '6h'
func ValidateDuration(val any, p cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	switch raw := val.(type) {
	case time.Duration, int64, float64:
		return diags
	case string:
		_, err := time.ParseDuration(raw)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Summary:  "Invalid duration value",
				Severity: diag.Error,
				Detail:   err.Error(),
			})
		}
	default:
		diags = append(diags, diag.Diagnostic{
			Summary:  "Invalid duration value",
			Severity: diag.Error,
		})
	}

	return diags
}
