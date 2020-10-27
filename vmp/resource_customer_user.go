// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package vmp

import (
	"context"
	"fmt"
	"strconv"
	"terraform-provider-vmp/vmp/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceCustomerUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomerUserCreate,
		ReadContext:   resourceCustomerUserRead,
		UpdateContext: resourceCustomerUserUpdate,
		DeleteContext: resourceCustomerUserDelete,

		Schema: map[string]*schema.Schema{
			"account_number": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"address1":        &schema.Schema{Type: schema.TypeString, Optional: true},
			"address2":        &schema.Schema{Type: schema.TypeString, Optional: true},
			"city":            &schema.Schema{Type: schema.TypeString, Optional: true},
			"country":         &schema.Schema{Type: schema.TypeString, Optional: true},
			"custom_id":       &schema.Schema{Type: schema.TypeString, Computed: true},
			"email":           &schema.Schema{Type: schema.TypeString, Optional: true},
			"fax":             &schema.Schema{Type: schema.TypeString, Optional: true},
			"first_name":      &schema.Schema{Type: schema.TypeString, Optional: true},
			"is_admin":        &schema.Schema{Type: schema.TypeBool, Optional: true, Default: false},
			"last_name":       &schema.Schema{Type: schema.TypeString, Optional: true},
			"mobile":          &schema.Schema{Type: schema.TypeString, Optional: true},
			"password":        &schema.Schema{Type: schema.TypeString, Optional: true, Sensitive: true},
			"phone":           &schema.Schema{Type: schema.TypeString, Optional: true},
			"state":           &schema.Schema{Type: schema.TypeString, Optional: true},
			"time_zone_id":    &schema.Schema{Type: schema.TypeInt, Optional: true},
			"title":           &schema.Schema{Type: schema.TypeString, Optional: true},
			"zip":             &schema.Schema{Type: schema.TypeString, Optional: true},
			"last_login_date": &schema.Schema{Type: schema.TypeString, Computed: true},
		},
	}
}

func resourceCustomerUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*ProviderConfiguration)

	payload := getCustomerUserFromData(d)

	accountNumber := d.Get("account_number").(string)

	apiClient := api.NewUserAPIClient(config.APIClient, config.PartnerID)

	customerUserID, err := apiClient.AddCustomerUser(accountNumber, payload)

	if err != nil {
		// Terraform requires an empty ID for failed creation
		d.SetId("")
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(customerUserID))

	return resourceCustomerUserRead(ctx, d, m)
}

func resourceCustomerUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*ProviderConfiguration)

	apiClient := api.NewUserAPIClient(config.APIClient, config.PartnerID)

	payload := getCustomerUserFromData(d)
	accountNumber := d.Get("account_number").(string)
	customerUserID, err := strconv.Atoi(d.Id())

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	err = apiClient.UpdateCustomerUser(accountNumber, customerUserID, payload)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	return resourceCustomerUserRead(ctx, d, m)
}

func resourceCustomerUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(*ProviderConfiguration)

	apiClient := api.NewUserAPIClient(config.APIClient, config.PartnerID)

	customerUserID, err := strconv.Atoi(d.Id())

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	accountNumber := d.Get("account_number").(string)

	resp, err := apiClient.GetCustomerUser(accountNumber, customerUserID)

	isAdminOld, isAdminNew := d.GetChange("is_admin")

	if isAdminOld.(bool) != isAdminNew.(bool) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "is_admin is a write-only property",
			Detail:   fmt.Sprintf("please set is_admin back to %t in your configuration for vmp_customer_user id=%s", isAdminOld, d.Id()),
		})
	}

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.Set("address1", resp.Address1)
	d.Set("address2", resp.Address2)
	d.Set("city", resp.City)
	d.Set("country", resp.Country)
	d.Set("custom_id", resp.CustomID)
	d.Set("email", resp.Email)
	d.Set("fax", resp.Fax)
	d.Set("first_name", resp.FirstName)
	d.Set("is_admin", resp.IsAdmin)
	d.Set("last_name", resp.LastName)
	d.Set("mobile", resp.Mobile)
	d.Set("password", resp.Password)
	d.Set("phone", resp.Phone)
	d.Set("state", resp.State)
	d.Set("time_zone_id", resp.TimeZoneID)
	d.Set("title", resp.Title)
	d.Set("zip", resp.ZIP)
	d.Set("last_login_date", resp.LastLoginDate)

	return diags
}

func resourceCustomerUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(*ProviderConfiguration)

	apiClient := api.NewUserAPIClient(config.APIClient, config.PartnerID)

	accountNumber := d.Get("account_number").(string)
	customerUserID, err := strconv.Atoi(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	err = apiClient.DeleteCustomerUser(accountNumber, customerUserID)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func getCustomerUserFromData(d *schema.ResourceData) *api.CustomerUser {
	var isAdmin int8 = 0

	if attr, ok := d.GetOk("is_admin"); ok {
		if attr.(bool) {
			isAdmin = 1
		}
	}

	// Go does not support nullable primitives, need to assign these as pointers
	var password *string
	var timeZoneID *int

	if attr, ok := d.GetOk("password"); ok && len(attr.(string)) > 0 {
		p := attr.(string)
		password = &p
	}

	if attr, ok := d.GetOk("time_zone_id"); ok {
		t := attr.(int)
		timeZoneID = &t
	}

	return &api.CustomerUser{
		Address1:   d.Get("address1").(string),
		Address2:   d.Get("address2").(string),
		City:       d.Get("city").(string),
		Country:    d.Get("country").(string),
		Email:      d.Get("email").(string),
		Fax:        d.Get("fax").(string),
		FirstName:  d.Get("first_name").(string),
		IsAdmin:    isAdmin,
		LastName:   d.Get("last_name").(string),
		Mobile:     d.Get("mobile").(string),
		Password:   password,
		Phone:      d.Get("phone").(string),
		State:      d.Get("state").(string),
		TimeZoneID: timeZoneID,
		Title:      d.Get("title").(string),
		ZIP:        d.Get("zip").(string),
	}
}
