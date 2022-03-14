package helper

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func parseAccountId(s string) (account, id string) {
	if strings.Contains(s, ":") == false {
		return s, ""
	}
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

func AccountIDImporter(read schema.ReadContextFunc) *schema.ResourceImporter {
	return Importer(read, "account_number", "id")
}

func AccountKeyImporter(key string, read schema.ReadContextFunc) *schema.ResourceImporter {
	return &schema.ResourceImporter{
		StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
			account, id := parseAccountId(d.Id())
			if account == "" || id == "" {
				return nil, fmt.Errorf("invalid id: %s", d.Id())
			}
			d.Set(key, account)
			d.SetId(id)

			_ = read(ctx, d, m)

			return []*schema.ResourceData{d}, nil
		},
	}
}

func parseKeys(s string) []string {
	if s == "" {
		return nil
	}
	if strings.Contains(s, ":") == false {
		return []string{s}
	}
	return strings.Split(s, ":")
}

func Importer(read schema.ReadContextFunc, keys ...string) *schema.ResourceImporter {
	return &schema.ResourceImporter{
		StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
			vals := parseKeys(d.Id())

			if len(keys) == 0 || len(vals) == 0 {
				return []*schema.ResourceData{d}, nil
			}

			for i, key := range keys {
				if strings.EqualFold(key, "id") {
					d.SetId(keys[i])
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
