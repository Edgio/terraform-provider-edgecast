// Copyright Edgecast, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package customer

import (
	"context"
	"log"
	"strconv"
	"terraform-provider-ec/ec/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceCustomerUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceCustomerUserCreate,
		ReadContext:   ResourceCustomerUserRead,
		UpdateContext: ResourceCustomerUserUpdate,
		DeleteContext: ResourceCustomerUserDelete,

		Schema: map[string]*schema.Schema{
			"account_number": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"address1":   {Type: schema.TypeString, Optional: true},
			"address2":   {Type: schema.TypeString, Optional: true},
			"city":       {Type: schema.TypeString, Optional: true},
			"country":    {Type: schema.TypeString, Optional: true},
			"custom_id":  {Type: schema.TypeString, Computed: true},
			"email":      {Type: schema.TypeString, Optional: true},
			"fax":        {Type: schema.TypeString, Optional: true},
			"first_name": {Type: schema.TypeString, Optional: true},
			"is_admin": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// IsAdmin is a write-only field, so suppress the diff if being changed
					// if the resource is new, do not suppress
					return len(old) > 0
				},
			},
			"last_name":       {Type: schema.TypeString, Optional: true},
			"mobile":          {Type: schema.TypeString, Optional: true},
			"phone":           {Type: schema.TypeString, Optional: true},
			"state":           {Type: schema.TypeString, Optional: true},
			"time_zone_id":    {Type: schema.TypeInt, Optional: true},
			"title":           {Type: schema.TypeString, Optional: true},
			"zip":             {Type: schema.TypeString, Optional: true},
			"last_login_date": {Type: schema.TypeString, Computed: true},
		},
	}
}

func ResourceCustomerUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(**api.ClientConfig)
	accountNumber := d.Get("account_number").(string)
	(*config).AccountNumber = accountNumber
	log.Printf("[INFO] Creating user for Account>> [AccountNumber]: %s", accountNumber)

	customerUser := getCustomerUserFromData(d)

	apiClient := api.NewUserAPIClient(*config)

	customerUserID, err := apiClient.AddCustomerUser(accountNumber, customerUser)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Successfully created user, ID=%d", customerUserID)

	d.SetId(strconv.Itoa(customerUserID))

	return ResourceCustomerUserRead(ctx, d, m)
}

func ResourceCustomerUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(**api.ClientConfig)
	accountNumber := d.Get("account_number").(string)
	(*config).AccountNumber = accountNumber
	apiClient := api.NewUserAPIClient(*config)

	customerUser := getCustomerUserFromData(d)
	customerUserID, err := strconv.Atoi(d.Id())

	log.Printf("[INFO] Updating Customer User %d for Account %s", customerUserID, accountNumber)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	err = apiClient.UpdateCustomerUser(accountNumber, customerUserID, customerUser)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	return ResourceCustomerUserRead(ctx, d, m)
}

func ResourceCustomerUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(**api.ClientConfig)
	accountNumber := d.Get("account_number").(string)
	(*config).AccountNumber = accountNumber
	apiClient := api.NewUserAPIClient(*config)

	customerUserID, err := strconv.Atoi(d.Id())

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Retreiving Customer User %d for Account %s", customerUserID, accountNumber)
	resp, err := apiClient.GetCustomerUser(accountNumber, customerUserID)

	isAdminOld, isAdminNew := d.GetChange("is_admin")

	if isAdminOld.(bool) != isAdminNew.(bool) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "note: is_admin is a write-only property; modifications will be ignored",
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

func ResourceCustomerUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(**api.ClientConfig)
	accountNumber := d.Get("account_number").(string)
	(*config).AccountNumber = accountNumber
	apiClient := api.NewUserAPIClient(*config)

	customerUserID, err := strconv.Atoi(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleting Customer User %d for Account %s", customerUserID, accountNumber)
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

	var timeZoneID *int

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
		Phone:      d.Get("phone").(string),
		State:      d.Get("state").(string),
		TimeZoneID: timeZoneID,
		Title:      d.Get("title").(string),
		ZIP:        d.Get("zip").(string),
	}
}
