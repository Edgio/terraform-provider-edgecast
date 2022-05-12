package helper

import (
	"context"
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

			if len(keys) == 0 || len(vals) == 0 {
				return []*schema.ResourceData{d}, nil
			}

			for i, key := range keys {
				if i < len(vals) {
					if strings.EqualFold(key, "id") {
						d.SetId(vals[i])
						continue
					}
					_ = d.Set(key, vals[i])
				}
			}

			_ = read(ctx, d, m)

			return []*schema.ResourceData{d}, nil
		},
	}
}
