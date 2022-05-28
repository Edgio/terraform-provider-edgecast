package dnsroute

import (
	"reflect"
	"sort"
	"testing"

	"github.com/EdgeCast/ec-sdk-go/edgecast/routedns"
)

func TestExpandMasterServers(t *testing.T) {

	cases := []struct {
		name          string
		input         []interface{}
		expectedPtr   []routedns.MasterServer
		expectSuccess bool
	}{
		{
			name: "Happy path",
			input: []interface{}{
				map[string]interface{}{
					"id":        1,
					"name":      "masterServer01",
					"ipaddress": "10.11.12.13",
				},
				map[string]interface{}{
					"id":        2,
					"name":      "masterServer02",
					"ipaddress": "10.11.12.14",
				},
			},
			expectedPtr: []routedns.MasterServer{
				{
					ID:        1,
					Name:      "masterServer01",
					IPAddress: "10.11.12.13",
				},
				{
					ID:        2,
					Name:      "masterServer02",
					IPAddress: "10.11.12.14",
				},
			},
			expectSuccess: true,
		},
		{
			name:          "Happy path - None Defined",
			input:         []interface{}{},
			expectedPtr:   []routedns.MasterServer{},
			expectSuccess: true,
		},
		{
			name:          "Edge case - nil input",
			input:         nil,
			expectedPtr:   []routedns.MasterServer{},
			expectSuccess: true,
		},
	}

	for _, v := range cases {
		actualPtr := expandMasterServers(v.input)

		actual := actualPtr
		expected := v.expectedPtr

		// array equality depends on order, sort by Name
		sort.Slice(actual, func(i, j int) bool {
			return actual[i].Name < actual[j].Name
		})

		if !reflect.DeepEqual(expected, actual) {
			t.Fatalf(
				"%s: Expected %+v but got %+v",
				v.name, expected, actual,
			)
		}

	}
}

func TestFlattenMasterServers(t *testing.T) {
	master01 := routedns.MasterServer{
		ID:        1,
		Name:      "master01",
		IPAddress: "10.11.12.13",
	}
	master02 := routedns.MasterServer{
		ID:        2,
		Name:      "master02",
		IPAddress: "10.11.12.14",
	}
	masters := make([]routedns.MasterServer, 0)
	masters = append(masters, master01, master02)

	cases := []struct {
		name          string
		input         routedns.MasterServerGroupAddGetOK
		expected      []map[string]interface{}
		expectSuccess bool
	}{
		{

			name: "Happy path",
			input: routedns.MasterServerGroupAddGetOK{
				MasterServerGroup: routedns.MasterServerGroup{
					Masters: masters,
				},
			},
			expected: []map[string]interface{}{
				{
					"id":        1,
					"name":      "master01",
					"ipaddress": "10.11.12.13",
				},
				{
					"id":        2,
					"name":      "master02",
					"ipaddress": "10.11.12.14",
				},
			},
			expectSuccess: true,
		},
		{
			name: "Nil path",
			input: routedns.MasterServerGroupAddGetOK{
				MasterServerGroup: routedns.MasterServerGroup{
					Masters: nil,
				},
			},
			expected:      nil,
			expectSuccess: true,
		},
	}

	for _, v := range cases {

		actualGroups := (v.input)

		for i, actual := range actualGroups.MasterServerGroup.Masters {

			expected := v.expected[i]

			if actual.ID != expected["id"] {
				t.Fatalf(
					"%s: Expected %+v but got %+v",
					v.name, expected["id"], actual.ID,
				)
			}

			if actual.Name != expected["name"] {
				t.Fatalf(
					"%s: Expected %+v but got %+v",
					v.name, expected["name"], actual.Name,
				)
			}

			if actual.IPAddress != expected["ipaddress"] {
				t.Fatalf(
					"%s: Expected %+v but got %+v",
					v.name, expected["ipaddress"], actual.IPAddress,
				)
			}
		}
	}
}
