// Copyright 2021 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.

package customer

import (
	"context"
	"log"
	"strconv"
	"terraform-provider-ec/ec/api"

	"github.com/EdgeCast/ec-sdk-go/edgecast/customer"
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

	accountNumber := d.Get("account_number").(string)
	log.Printf("[INFO] Creating customer user for Account>> [AccountNumber]: %s", accountNumber)

	config := m.(**api.ClientConfig)
	(*config).AccountNumber = accountNumber

	customerService, err := buildCustomerService(**config)
	if err != nil {
		return diag.FromErr(err)
	}

	customer, err := customerService.GetCustomer(accountNumber)
	if err != nil {
		return diag.FromErr(err)
	}

	customerUser := getCustomerUserFromData(d)
	customerUserID, err := customerService.AddCustomerUser(customer, customerUser)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Successfully created user, ID=%d", customerUserID)

	d.SetId(strconv.Itoa(customerUserID))

	return ResourceCustomerUserRead(ctx, d, m)
}

func ResourceCustomerUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	accountNumber := d.Get("account_number").(string)
	customerUser := getCustomerUserFromData(d)

	customerUserID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Updating customer user %d for account %s", customerUserID, accountNumber)

	config := m.(**api.ClientConfig)
	(*config).AccountNumber = accountNumber

	customerService, err := buildCustomerService(**config)
	if err != nil {
		return diag.FromErr(err)
	}

	customer, err := customerService.GetCustomer(accountNumber)
	if err != nil {
		return diag.FromErr(err)
	}

	err = customerService.UpdateCustomerUser(customer, customerUser)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	return ResourceCustomerUserRead(ctx, d, m)
}

func ResourceCustomerUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	accountNumber := d.Get("account_number").(string)
	customerUserID, err := strconv.Atoi(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Retreiving Customer User %d for Account %s", customerUserID, accountNumber)

	config := m.(**api.ClientConfig)
	(*config).AccountNumber = accountNumber
	customerService, err := buildCustomerService(**config)
	if err != nil {
		return diag.FromErr(err)
	}

	customer, err := customerService.GetCustomer(accountNumber)
	if err != nil {
		return diag.FromErr(err)
	}

	customerUser, err := customerService.GetCustomerUser(customer, customerUserID)

	if err != nil {
		return diag.FromErr(err)
	}

	isAdminOld, isAdminNew := d.GetChange("is_admin")

	if isAdminOld.(bool) != isAdminNew.(bool) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "note: is_admin is a write-only property; modifications will be ignored",
		})
	}

	d.Set("address1", customerUser.Address1)
	d.Set("address2", customerUser.Address2)
	d.Set("city", customerUser.City)
	d.Set("country", customerUser.Country)
	d.Set("custom_id", customerUser.CustomID)
	d.Set("email", customerUser.Email)
	d.Set("fax", customerUser.Fax)
	d.Set("first_name", customerUser.FirstName)
	d.Set("is_admin", customerUser.IsAdmin)
	d.Set("last_name", customerUser.LastName)
	d.Set("mobile", customerUser.Mobile)
	d.Set("phone", customerUser.Phone)
	d.Set("state", customerUser.State)
	d.Set("title", customerUser.Title)
	d.Set("zip", customerUser.Zip)
	d.Set("last_login_date", customerUser.LastLoginDate)

	return diags
}

func ResourceCustomerUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	accountNumber := d.Get("account_number").(string)
	customerUserID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleting Customer User %d for Account %s", customerUserID, accountNumber)

	config := m.(**api.ClientConfig)
	(*config).AccountNumber = accountNumber
	customerService, err := buildCustomerService(**config)
	if err != nil {
		return diag.FromErr(err)
	}

	customer, err := customerService.GetCustomer(accountNumber)
	if err != nil {
		return diag.FromErr(err)
	}

	customerUser, err := customerService.GetCustomerUser(customer, customerUserID)
	if err != nil {
		return diag.FromErr(err)
	}

	err = customerService.DeleteCustomerUser(customer, *customerUser)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func getCustomerUserFromData(d *schema.ResourceData) *customer.CustomerUser {
	var isAdmin int8 = 0

	if attr, ok := d.GetOk("is_admin"); ok {
		if attr.(bool) {
			isAdmin = 1
		}
	}

	return &customer.CustomerUser{
		Address1:  d.Get("address1").(string),
		Address2:  d.Get("address2").(string),
		City:      d.Get("city").(string),
		Country:   d.Get("country").(string),
		Email:     d.Get("email").(string),
		Fax:       d.Get("fax").(string),
		FirstName: d.Get("first_name").(string),
		IsAdmin:   isAdmin,
		LastName:  d.Get("last_name").(string),
		Mobile:    d.Get("mobile").(string),
		Phone:     d.Get("phone").(string),
		State:     d.Get("state").(string),
		Title:     d.Get("title").(string),
		Zip:       d.Get("zip").(string),
	}
}
