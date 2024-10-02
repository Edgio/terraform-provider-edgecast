package rulesengine

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type JSONMap map[string]any

func (j JSONMap) String() string {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func Test_policyDiffSuppress(t *testing.T) {
	source := JSONMap{
		"name": "a",
		"a":    "b",
		"c":    "d",
	}
	diffNameShouldMatch := JSONMap{
		"name": "b",
		"a":    "b",
		"c":    "d",
	}
	diffNameShouldFail := JSONMap{
		"name": "b",
		"a":    "x",
		"c":    "x",
	}
	sameNameShouldMatch := JSONMap{
		"name": "a",
		"a":    "b",
		"c":    "d",
	}
	sameNameShouldFail := JSONMap{
		"name": "a",
		"a":    "x",
		"c":    "x",
	}

	type args struct {
		k   string
		old string
		new string
		in3 *schema.ResourceData
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Different Name Should Match",
			args: args{
				k:   "",
				old: source.String(),
				new: diffNameShouldMatch.String(),
				in3: nil,
			},
			want: true,
		},
		{
			name: "Different Name Should Fail",
			args: args{
				k:   "",
				old: source.String(),
				new: diffNameShouldFail.String(),
				in3: nil,
			},
			want: false,
		},

		{
			name: "Same Name Should Match",
			args: args{
				k:   "",
				old: source.String(),
				new: sameNameShouldMatch.String(),
				in3: nil,
			},
			want: true,
		},
		{
			name: "Different Name Should Fail",
			args: args{
				k:   "",
				old: source.String(),
				new: sameNameShouldFail.String(),
				in3: nil,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := policyDiffSuppress(tt.args.k, tt.args.old, tt.args.new, tt.args.in3); got != tt.want {
				t.Errorf("policyDiffSuppress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cleanMatches(t *testing.T) {

	policy := map[string]any{
		"@type":       "Policy",
		"name":        "tf--staging-http_large-1727827247",
		"policy_type": "customer",
		"state":       "locked",
		"created_at":  "2024-10-02T00:00:47Z",
		"rules": []any{
			map[string]interface{}{
				"name":        "CustOriginMatch",
				"description": "Rule to match origin to test raw values",
				"ordinal":     1,
				"created_at":  "2024-10-02T00:00:47Z",
				"@id":         "/rules-engine/v1.1/policies/4409533/rules/33753446",
				"@type":       "Rule",
				"id":          "33753446",
				"updated_at":  "2024-10-02T00:00:47Z",
				"matches": []any{
					map[string]any{
						"type":      "match.origin.customer-origin.literal",
						"ordinal":   1,
						"value":     "/804FDBB/test/",
						"raw_value": "/804FDBB/test/",
						"features": []any{
							map[string]any{
								"type":            "feature.url.url-redirect",
								"ordinal":         1,
								"code":            "301",
								"destination":     "/804FDBB/test/testdest",
								"raw_destination": "/804FDBB/test/testdest",
								"source":          "/804FDBB/test/testsource",
								"raw_source":      "/804FDBB/test/testsource",
							},
						},
					},
				},
			},
		},
		"@id":         "/rules-engine/v1.1/policies/4409533",
		"id":          "4409533",
		"description": "This is a test policy of PolicyCreate.",
		"platform":    "http_large",
		"updated_at":  "2024-10-02T00:00:47Z",
		"history":     []any{},
	}

	cleanedPolicy := map[string]any{
		"name": "tf--staging-http_large-1727827247",
		"rules": []map[string]interface{}{
			{
				"name":        "CustOriginMatch",
				"description": "Rule to match origin to test raw values",
				"matches": []map[string]interface{}{
					{
						"type":  "match.origin.customer-origin.literal",
						"value": "/804FDBB/test/",
						"features": []map[string]interface{}{
							{
								"type":        "feature.url.url-redirect",
								"code":        "301",
								"destination": "/804FDBB/test/testdest",
								"source":      "/804FDBB/test/testsource",
							},
						},
					},
				},
			},
		},
		"description": "This is a test policy of PolicyCreate.",
		"platform":    "http_large",
	}

	tests := []struct {
		name           string
		originalPolicy map[string]any
		expectedPolicy map[string]any
	}{
		{
			name:           "Unneeded Values Removed",
			originalPolicy: policy,
			expectedPolicy: cleanedPolicy,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cleanPolicy(tt.originalPolicy)
			if err != nil {
				t.Errorf("cleanPolicy() function call errored unexpectedly %v", err)
			}

			if !reflect.DeepEqual(tt.originalPolicy, tt.expectedPolicy) {
				t.Errorf("cleanPolicy() = got %v, want %v", tt.originalPolicy, tt.expectedPolicy)
			}
		})
	}
}
