package rulesengine

import (
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"testing"
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
