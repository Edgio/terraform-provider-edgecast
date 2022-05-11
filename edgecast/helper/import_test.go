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

func TestImporter_AccountNumberWithID(t *testing.T) {
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

func TestImporter_OptionalKeys(t *testing.T) {
	expect := assert.New(t)
	i, rd := createResourceData(t, "account_number", "id", "portal_type_id", "customer_user_id")
	expect.NotNil(i)
	expect.NotNil(rd)
	rd.SetId("123:456")
	rds, err := i.StateContext(context.Background(), rd, nil)
	expect.NoError(err)
	expect.NotNil(rds)
	expect.Equal("456", rd.Id())
	expect.Equal("123", rd.Get("account_number"))
	expect.Equal("", rd.Get("portal_type_id"))
	expect.Equal("", rd.Get("customer_user_id"))

}

func TestImporter_AccountNumberWithID_MediaType(t *testing.T) {
	expect := assert.New(t)
	i, rd := createResourceData(t, "account_number", "id", "media_type_id")
	expect.NotNil(i)
	expect.NotNil(rd)
	rd.SetId("123:456:789")
	rds, err := i.StateContext(context.Background(), rd, nil)
	expect.NoError(err)
	expect.NotNil(rds)
	expect.Equal("456", rd.Id())
	expect.Equal("123", rd.Get("account_number"))
	expect.Equal("789", rd.Get("media_type_id"))
}

func TestImporter_AccountNumber(t *testing.T) {
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

func TestImporter_ID(t *testing.T) {
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
