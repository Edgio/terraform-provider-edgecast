package helper

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func readSchema(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}
func createResourceData(t *testing.T, keys ...string) (*schema.ResourceImporter, *schema.ResourceData) {
	i := Import(readSchema, keys...)
	s := make(map[string]*schema.Schema)

	for _, key := range keys {
		if strings.EqualFold(key, "id") {
			continue
		}
		s[key] = &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		}
	}
	rd := schema.TestResourceDataRaw(t, s, map[string]interface{}{})
	return i, rd
}

func TestImporter_account_number_and_id(t *testing.T) {
	expect := assert.New(t)
	i, rd := createResourceData(t, "account_number", "id")
	expect.NotNil(i)
	expect.NotNil(rd)
	rd.SetId("123:456")
	rds, err := i.StateContext(context.Background(), rd, nil)
	expect.NoError(err)
	expect.NotNil(rds)
	expect.Equal("456", rd.Id())
	expect.Equal("123", rd.Get("account_number"))
}

func TestImporter_account_number_and_group_product_and_id(t *testing.T) {
	expect := assert.New(t)
	i, rd := createResourceData(t, "account_number", "group_product_type", "id")
	expect.NotNil(i)
	expect.NotNil(rd)
	rd.SetId("123:456:789")
	rds, err := i.StateContext(context.Background(), rd, nil)
	expect.NoError(err)
	expect.NotNil(rds)
	expect.Equal("789", rd.Id())
	expect.Equal("123", rd.Get("account_number"))
	expect.Equal("456", rd.Get("group_product_type"))
}

func TestImporter_account_number_and_media_type_and_id(t *testing.T) {
	expect := assert.New(t)
	i, rd := createResourceData(t, "account_number", "media_type_id", "id")
	expect.NotNil(i)
	expect.NotNil(rd)
	rd.SetId("123:456:789")
	rds, err := i.StateContext(context.Background(), rd, nil)
	expect.NoError(err)
	expect.NotNil(rds)
	expect.Equal("789", rd.Id())
	expect.Equal("123", rd.Get("account_number"))
	expect.Equal("456", rd.Get("media_type_id"))
}

func TestImporter_account_number(t *testing.T) {
	expect := assert.New(t)
	i, rd := createResourceData(t, "account_number")
	expect.NotNil(i)
	expect.NotNil(rd)
	rd.SetId("123")
	rds, err := i.StateContext(context.Background(), rd, nil)
	expect.NoError(err)
	expect.NotNil(rds)

	expect.Equal("123", rd.Get("account_number"))
}

func TestImporter_id(t *testing.T) {
	expect := assert.New(t)
	i, rd := createResourceData(t, "id")
	expect.NotNil(i)
	expect.NotNil(rd)
	rd.SetId("123")
	rds, err := i.StateContext(context.Background(), rd, nil)
	expect.NoError(err)
	expect.NotNil(rds)
	expect.Equal("123", rd.Id())
}

func TestImporter_NoKeys(t *testing.T) {
	expect := assert.New(t)
	i, rd := createResourceData(t)
	expect.NotNil(i)
	expect.NotNil(rd)
	rds, err := i.StateContext(context.Background(), rd, nil)
	expect.NoError(err)
	expect.NotNil(rds)
	expect.Empty(rd.Id())
}
