package dnsroute

import (
	"reflect"
	"sort"
	"testing"

	"github.com/EdgeCast/ec-sdk-go/edgecast/routedns"
)

func TestExpandZoneCompositionCreate(t *testing.T) {

	cases := []struct {
		name          string
		input         interface{}
		expectedPtr   routedns.ZoneComposition
		expectSuccess bool
	}{
		{
			name: "Happy path",
			input: map[string]interface{}{
				"id":        1,
				"name":      "masterServer01",
				"ipaddress": "10.11.12.13",
			},
			expectedPtr:   createZoneComposition(),
			expectSuccess: true,
		},
		{
			name:          "Happy path - None Defined",
			input:         []interface{}{},
			expectedPtr:   routedns.ZoneComposition{},
			expectSuccess: true,
		},
		{
			name:          "Edge case - nil input",
			input:         nil,
			expectedPtr:   routedns.ZoneComposition{},
			expectSuccess: true,
		},
	}

	for _, v := range cases {
		actualPtr, error := expandZoneCompositionCreate(v.input)

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

func TestFlattenZones(t *testing.T) {
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
	masters := make([]routedns.MasterServer, 0, 2)
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

		expected := v.input
		actual := flattenMasterServers(v.input)

		for i, master := range expected.MasterServerGroup.Masters {

			actualMap := actual[i].(map[string]interface{})

			if master.ID != actualMap["id"] {
				t.Fatalf(
					"%s: Expected %+v but got %+v",
					v.name, master.ID, actualMap["id"],
				)
			}

			if master.Name != actualMap["name"] {
				t.Fatalf(
					"%s: Expected %+v but got %+v",
					v.name, master.Name, actualMap["name"],
				)
			}

			if master.IPAddress != actualMap["ipaddress"] {
				t.Fatalf(
					"%s: Expected %+v but got %+v",
					v.name, master.IPAddress, actualMap["ipaddress"],
				)
			}
		}
	}
}

func createZoneComposition() routedns.ZoneComposition {
	tsigID1 := routedns.MasterServerTSIGIDs{
		MasterServer: routedns.MasterServerID{ID: 1},
		TSIG:         routedns.TSIGID{ID: 1},
	}
	tsigID2 := routedns.MasterServerTSIGIDs{
		MasterServer: routedns.MasterServerID{ID: 2},
		TSIG:         routedns.TSIGID{ID: 2},
	}
	tsigIDs := make([]routedns.MasterServerTSIGIDs, 0, 2)
	tsigIDs = append(tsigIDs, tsigID1, tsigID2)

	sz := routedns.SecondaryZone{
		DomainName: "testdomain.com",
		Status:     1,
	}

	szl := make([]routedns.SecondaryZone, 0, 1)
	szl = append(szl, sz)

	return routedns.ZoneComposition{
		MasterGroupID:     1,
		MasterServerTSIGs: tsigIDs,
		Zones:             szl,
	}
}
