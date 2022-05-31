package helper

import (
	"context"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func parseKeys(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return strings.Split(s, ":")
}

// Import provides a schema.ResourceImporter to parse keys using the ReadContextFunc
func Import(read schema.ReadContextFunc, keys ...string) *schema.ResourceImporter {
	return &schema.ResourceImporter{
		StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

			vals := parseKeys(d.Id())
			var id string
			if len(keys) == 0 || len(vals) == 0 {
				return []*schema.ResourceData{d}, nil
			}

			for i, key := range keys {
				if strings.EqualFold(key, "id") {
					id = vals[i]
					continue
				}
				if i < len(vals) {
					err := d.Set(key, vals[i])
					if err != nil {
						v, _ := strconv.Atoi(vals[i])
						_ = d.Set(key, v)
					}
				}
			}
			_ = read(ctx, d, m)
			d.SetId(id)
			return []*schema.ResourceData{d}, nil
		},
	}
}
