package helper

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
					d.SetId(id)
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
			if res := read(ctx, d, m); res != nil {
				for _, e := range res {
					if e.Severity == diag.Error {
						return nil, fmt.Errorf("%s\n%s", e.Summary, e.Detail)
					}
				}
			}
			d.SetId(id)
			return []*schema.ResourceData{d}, nil
		},
	}
}
