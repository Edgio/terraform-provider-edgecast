// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package helper

// Filter returns a new slice containing all items in the slice that satisfy the predicate f.
func Filter(vs []interface{}, f func(interface{}) bool) []interface{} {
	vsf := make([]interface{}, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

// IntArrayDiff finds the set difference of set1 relative to set2
func IntArrayDiff(set1 []int, set2 []int) []int {
	diff := []int{}

	for _, v1 := range set1 {
		foundInSet2 := false
		for _, v2 := range set2 {
			if v1 == v2 {
				foundInSet2 = true
				break
			}
		}

		if !foundInSet2 {
			diff = append(diff, v1)
		}
	}

	return diff
}
