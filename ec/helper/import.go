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
	if strings.Contains(s, ":") == false {
		return []string{s}
	}
	return strings.Split(s, ":")
}

func Import(read schema.ReadContextFunc, keys ...string) *schema.ResourceImporter {
	return &schema.ResourceImporter{
		StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
			vals := parseKeys(d.Id())

			if len(keys) == 0 || len(vals) == 0 {
				return []*schema.ResourceData{d}, nil
			}

			for i, key := range keys {
				if strings.EqualFold(key, "id") {
					d.SetId(vals[i])
					continue
				}
				if i < len(vals) {
					_ = d.Set(key, vals[i])
				}
			}

			_ = read(ctx, d, m)

			return []*schema.ResourceData{d}, nil
		},
	}
}
