// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.
package originv3

import (
	"testing"

	"github.com/EdgeCast/ec-sdk-go/edgecast/originv3"
	"github.com/go-test/deep"
)

func TestFlattenOriginShieldEdgeNodes(t *testing.T) {
	t.Parallel()

	regionID1 := int32(1)
	bypassCode1 := "code 1"
	bypassName1 := "name 1"
	region1 := "region 1"

	regionID2 := int32(2)
	bypassCode2 := "code 2"
	bypassName2 := "name 2"
	region2 := "region 2"

	pop1ID := int32(1)
	pop1Code := "pop code 1"
	pop1City := "pop city 1"
	pop2ID := int32(2)
	pop2Code := "pop code 2"
	pop2City := "pop city 2"
	pop3ID := int32(3)
	pop3Code := "pop code 3"
	pop3City := "pop city 3"
	isPCICertified := false

	tests := []struct {
		name string
		args []originv3.OriginShieldEdgeNode
		want []map[string]interface{}
	}{
		{
			name: "Happy Path",
			args: []originv3.OriginShieldEdgeNode{
				{
					BypassCode: &bypassCode1,
					BypassName: &bypassName1,
					RegionId:   &regionID1,
					RegionName: &region1,
					Pops: []originv3.OriginShieldPop{
						{
							Id:             &pop1ID,
							Code:           &pop1Code,
							City:           &pop1City,
							IsPciCertified: &isPCICertified,
						},
						{
							Id:             &pop2ID,
							Code:           &pop2Code,
							City:           &pop2City,
							IsPciCertified: &isPCICertified,
						},
					},
				},
				{
					BypassCode: &bypassCode2,
					BypassName: &bypassName2,
					RegionId:   &regionID2,
					RegionName: &region2,
					Pops: []originv3.OriginShieldPop{
						{
							Id:             &pop3ID,
							Code:           &pop3Code,
							City:           &pop3City,
							IsPciCertified: &isPCICertified,
						},
					},
				},
			},
			want: []map[string]interface{}{
				{
					"bypass_code": &bypassCode1,
					"bypass_name": &bypassName1,
					"region_id":   &regionID1,
					"region_name": &region1,
					"pops": []map[string]interface{}{
						{
							"id":               &pop1ID,
							"code":             &pop1Code,
							"city":             &pop1City,
							"is_pci_certified": &isPCICertified,
						},
						{
							"id":               &pop2ID,
							"code":             &pop2Code,
							"city":             &pop2City,
							"is_pci_certified": &isPCICertified,
						},
					},
				},
				{
					"bypass_code": &bypassCode2,
					"bypass_name": &bypassName2,
					"region_id":   &regionID2,
					"region_name": &region2,
					"pops": []map[string]interface{}{
						{
							"id":               &pop3ID,
							"code":             &pop3Code,
							"city":             &pop3City,
							"is_pci_certified": &isPCICertified,
						},
					},
				},
			},
		},
		{
			name: "No Pops",
			args: []originv3.OriginShieldEdgeNode{
				{
					BypassCode: &bypassCode1,
					BypassName: &bypassName1,
					RegionId:   &regionID1,
					RegionName: &region1,
					Pops:       nil,
				},
				{
					BypassCode: &bypassCode2,
					BypassName: &bypassName2,
					RegionId:   &regionID2,
					RegionName: &region2,
					Pops:       nil,
				},
			},
			want: []map[string]interface{}{
				{
					"bypass_code": &bypassCode1,
					"bypass_name": &bypassName1,
					"region_id":   &regionID1,
					"region_name": &region1,
					"pops":        []map[string]interface{}{},
				},
				{
					"bypass_code": &bypassCode2,
					"bypass_name": &bypassName2,
					"region_id":   &regionID2,
					"region_name": &region2,
					"pops":        []map[string]interface{}{},
				},
			},
		},
		{
			name: "No results",
			args: nil,
			want: make([]map[string]interface{}, 0),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := flattenOriginShieldEdgeNodes(tt.args)

			diffs := deep.Equal(got, tt.want)
			if len(diffs) > 0 {
				t.Logf("got %v, want %v", got, tt.want)
				t.Errorf("Differences: %v", diffs)
			}
		})
	}
}

func TestFlattenProtocolTypes(t *testing.T) {
	t.Parallel()

	id1 := int32(1)
	id2 := int32(2)
	name1 := "protocol name 1"
	name2 := "protocol name 2"
	tests := []struct {
		name string
		args []originv3.ProtocolType
		want []map[string]interface{}
	}{
		{
			name: "Happy Path",
			args: []originv3.ProtocolType{
				{
					Id:   &id1,
					Name: &name1,
				},
				{
					Id:   &id2,
					Name: &name2,
				},
			},
			want: []map[string]interface{}{
				{
					"id":   &id1,
					"name": &name1,
				},
				{
					"id":   &id2,
					"name": &name2,
				},
			},
		},
		{
			name: "No results",
			args: nil,
			want: make([]map[string]interface{}, 0),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := flattenProtocolTypes(tt.args)

			diffs := deep.Equal(got, tt.want)
			if len(diffs) > 0 {
				t.Logf("got %v, want %v", got, tt.want)
				t.Errorf("Differences: %v", diffs)
			}
		})
	}
}

func TestFlattenHostnameResolutionMethods(t *testing.T) {
	t.Parallel()

	id1 := int32(1)
	id2 := int32(2)
	name1 := "protocol name 1"
	name2 := "protocol name 2"
	value1 := "value 1"
	value2 := "value 2"
	tests := []struct {
		name string
		args []originv3.NetworkType
		want []map[string]interface{}
	}{
		{
			name: "Happy Path",
			args: []originv3.NetworkType{
				{
					Id:    &id1,
					Name:  &name1,
					Value: &value1,
				},
				{
					Id:    &id2,
					Name:  &name2,
					Value: &value2,
				},
			},
			want: []map[string]interface{}{
				{
					"id":    &id1,
					"name":  &name1,
					"value": &value1,
				},
				{
					"id":    &id2,
					"name":  &name2,
					"value": &value2,
				},
			},
		},
		{
			name: "No results",
			args: nil,
			want: make([]map[string]interface{}, 0),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := flattenHostnameResolutionMethods(tt.args)

			diffs := deep.Equal(got, tt.want)
			if len(diffs) > 0 {
				t.Logf("got %v, want %v", got, tt.want)
				t.Errorf("Differences: %v", diffs)
			}
		})
	}
}
