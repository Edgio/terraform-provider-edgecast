// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.

package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	timeStampNumberBase = 10
)

// CreationErrorf is a helper function for errors encountered while
// creating a new resource.
func CreationErrorf(
	d *schema.ResourceData,
	format string,
	a ...interface{},
) diag.Diagnostics {
	d.SetId("")

	return diag.Errorf(format, a...)
}

// CreationError is a helper function for errors encountered while
// creating a new resource.
func CreationError(d *schema.ResourceData, err error) diag.Diagnostics {
	d.SetId("")

	return diag.FromErr(err)
}

// CreationErrors is a helper function for errors encountered while
// creating a new resource.
func CreationErrors(
	d *schema.ResourceData,
	msg string,
	errs []error,
) diag.Diagnostics {
	d.SetId("")

	return DiagsFromErrors(msg, errs)
}

// DiagFromError wraps an error in a diag.Diagnostic with an additional message.
func DiagFromError(msg string, err error) diag.Diagnostics {
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Error,
			Summary:  msg,
			Detail:   err.Error(),
		},
	}
}

// DiagsFromErrors creates a diag.Diagnostics instance from multiple errors and
// a base message.
func DiagsFromErrors(msg string, errs []error) diag.Diagnostics {
	// start with a base message
	var diags diag.Diagnostics
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  msg,
	})

	for _, err := range errs {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  err.Error(),
		})
	}

	return diags
}

// GetUnixTimeStamp gets the current unix time and returns it as a string.
func GetUnixTimeStamp() string {
	return strconv.FormatInt(time.Now().Unix(), timeStampNumberBase)
}

// ConvertTFCollectionToSlice converts Terraform's
// TypeList and TypeSet collections into a []interface{}.
func ConvertTFCollectionToSlice(v interface{}) ([]interface{}, error) {
	// TF TypeList results in a interface{} whose underlying type is
	// []interface{}, so we just cast and return it
	if listItems, ok := v.([]interface{}); ok {
		return listItems, nil
	}

	// TF TypeSet results in a *schema.Set, which must be converted to a slice
	// using List()
	if set, ok := v.(*schema.Set); ok {
		return set.List(), nil
	}

	return nil, errors.New("input was not a TF collection")
}

// ConvertTFCollectionToStrings converts Terraform's
// TypeList and TypeSet collections into a []string.
func ConvertTFCollectionToStrings(v interface{}) ([]string, error) {
	// If the underlying type is already []string, just return it.
	if strings, ok := v.([]string); ok {
		return strings, nil
	}

	listItems, err := ConvertTFCollectionToSlice(v)

	if err == nil {
		return ConvertSliceToStrings(listItems)
	}

	return nil, err
}

// ConvertSingletonSetToMap converts a interface{} that is actually a Terraform
// schema.TypeSet into a Map. This is useful when using the schema.TypeSet
// with MaxItems=1 workaround.
func ConvertSingletonSetToMap(attr interface{}) (map[string]interface{}, error) {
	if attr == nil {
		return nil, fmt.Errorf("set was nil")
	}

	// Must convert the set to a slice/array to work with the values
	if set, ok := attr.(*schema.Set); ok {
		arr := set.List()

		// Terraform will normally not allow more than one item because we
		// specify MaxItems=1, but we'll check again here
		if len(arr) > 1 {
			return nil, fmt.Errorf(
				"set must contain exactly one item, actual #: %d", len(arr))
		}

		if len(arr) == 0 {
			return make(map[string]interface{}), nil
		}

		// The single item in the list
		// is expected to be a map[string][]interface{}
		if entryMap, ok := arr[0].(map[string]interface{}); ok {
			return entryMap, nil
		} else {
			return nil, fmt.Errorf(
				"%v must be of type map[string]interface{}, actual: %T",
				arr[0],
				arr[0])
		}
	}

	return nil, fmt.Errorf("attr was not a *schema.Set")
}

// NewTerraformSet is useful for unit testing when mocking Terraform.
func NewTerraformSet(items []interface{}) *schema.Set {
	return schema.NewSet(dummySetFunc, items)
}

// ConvertToInt converts an interface{} to int
// If v is nil or not an int, 0 is returned.
func ConvertToInt(v interface{}) int {
	i, _ := v.(int)

	return i
}

// ConvertToIntPointer converts an interface{} to *int
// If v is nil or not an integer, a nil pointer is returned.
func ConvertToIntPointer(v interface{}, assertNaturalNumber bool) *int {
	if v == nil {
		return nil
	}

	if i, ok := v.(int); ok {
		if assertNaturalNumber && i <= 0 {
			return nil
		}

		return &i
	}

	return nil
}

// ConvertToBoolPointer converts an interface{} to *bool
// If v is nil or not a bool, a nil pointer is returned.
func ConvertToBoolPointer(v interface{}) *bool {
	if v == nil {
		return nil
	}

	if b, ok := v.(bool); ok {
		return &b
	}

	return nil
}

// ConvertToString converts an interface{} to string
// If v is nil or not a string, an empty string is returned.
func ConvertToString(v interface{}) string {
	s, _ := v.(string)

	return s
}

// ConvertToStringPointer converts an interface{} to *string
// If v is nil or not a string, a nil pointer is returned.
func ConvertToStringPointer(v interface{}, excludeWhiteSpace bool) *string {
	if s, ok := v.(string); ok {
		if excludeWhiteSpace && len(strings.TrimSpace(s)) == 0 {
			return nil
		}

		return &s
	}

	return nil
}

// ConvertToStringSlicePointer converts an interface{} to *[]string
// If v is nil or not a string, a nil pointer is returned. If a value within
// v is not a string, it will be returned as an empty string.
func ConvertToStringsPointer(v interface{}, excludeEmpty bool) *[]string {
	if v == nil {
		return nil
	}

	if l, ok := v.([]interface{}); ok {
		if excludeEmpty && len(l) == 0 {
			return nil
		}

		s := make([]string, len(l))

		for i, v := range l {
			vs, ok := v.(string)
			if !ok {
				vs = ""
			}

			s[i] = vs
		}

		return &s
	}

	return nil
}

// ConvertToMapPointer converts an interface{} to *map[string]string{}
// If v is nil or not a string, a nil pointer is returned.
func ConvertToStringMapPointer(v interface{}, excludeEmpty bool) *map[string]string {
	if v == nil {
		return nil
	}

	if m, ok := v.(map[string]interface{}); ok {
		if excludeEmpty && len(m) == 0 {
			return nil
		}

		result := make(map[string]string)

		for key, value := range m {
			stringValue, _ := value.(string)
			result[key] = stringValue
		}

		return &result
	}

	return nil
}

// StringIsNotEmptyJSON is a SchemaValidateFunc which tests to make sure the
// supplied string is not an empty JSON object i.e. "{}".
func StringIsNotEmptyJSON(
	i interface{},
	k string,
) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(
			errors,
			fmt.Errorf("expected type of %s to be string", k))

		return warnings, errors
	}

	var data map[string]interface{}

	if err := json.Unmarshal([]byte(v), &data); err != nil {
		errors = append(
			errors,
			fmt.Errorf("%q contains invalid JSON: %s", k, err))

		return warnings, errors
	}

	if len(data) == 0 {
		errors = append(
			errors,
			fmt.Errorf("%q contains empty JSON: '%s'", k, v))
	}

	return warnings, errors
}

// dummySetFunc is to be used when imitating Terraform
// in unit tests by using schema.NewSet.
func dummySetFunc(i interface{}) int {
	return random.Random(math.MinInt32, math.MaxInt32)
}
