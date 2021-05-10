package helper

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceFiltersSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},

				"values": {
					Type:     schema.TypeList,
					Required: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}

func IsJsonArray(input interface{}) bool {
	switch input.(type) {
	case []interface{}:
		return true
	default:
		return false
	}
}

func IsJSONString(s string) bool {
	var js string
	return json.Unmarshal([]byte(s), &js) != nil

}
