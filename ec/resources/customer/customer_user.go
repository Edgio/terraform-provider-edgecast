// Copyright 2021 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.

package customer

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"terraform-provider-edgecast/ec/api"

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
				DiffSuppressFunc: func(
					k,
					old,
					new string,
					d *schema.ResourceData,
				) bool {
					// IsAdmin is a write-only field, so suppress the diff if
					// being changed. If the resource is new, do not suppress
					return len(old) > 0
				},
			},
			"last_name":       {Type: schema.TypeString, Optional: true},
			"mobile":          {Type: schema.TypeString, Optional: true},
			"phone":           {Type: schema.TypeString, Optional: true},
			"state":           {Type: schema.TypeString, Optional: true},
			"title":           {Type: schema.TypeString, Optional: true},
			"zip":             {Type: schema.TypeString, Optional: true},
			"last_login_date": {Type: schema.TypeString, Computed: true},
		},
	}
}

func ResourceCustomerUserCreate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {

	accountNumber := d.Get("account_number").(string)
	log.Printf(
		"[INFO] Creating customer user for [Account Number]: %s",
		accountNumber,
	)

	// Initialize Customer Service
	config := m.(**api.ClientConfig)
	customerService, err := buildCustomerService(**config)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	// Retrieve Customer object from API
	getCustomerParams := customer.NewGetCustomerParams()
	getCustomerParams.AccountNumber = accountNumber
	customerObj, err := customerService.GetCustomer(*getCustomerParams)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	// Call Add Customer User API
	addCustUserParams := customer.NewAddCustomerUserParams()
	addCustUserParams.Customer = *customerObj
	addCustUserParams.CustomerUser = *getCustomerUserFromData(d)
	customerUserID, err := customerService.AddCustomerUser(*addCustUserParams)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	log.Printf(
		"[INFO] Created [Customer User]: %d for [Account Number]: %s",
		customerUserID,
		accountNumber,
	)

	d.SetId(strconv.Itoa(customerUserID))

	return ResourceCustomerUserRead(ctx, d, m)
}

func ResourceCustomerUserUpdate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	accountNumber := d.Get("account_number").(string)
	customerUser := getCustomerUserFromData(d)

	customerUserID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf(
		"[INFO] Updating [Customer User]: %d for [Account Number]: %s",
		customerUserID,
		accountNumber,
	)

	// Initialize Customer Service
	config := m.(**api.ClientConfig)
	customerService, err := buildCustomerService(**config)
	if err != nil {
		return diag.FromErr(err)
	}

	// Retrieve Customer object from API
	getCustomerParams := customer.NewGetCustomerParams()
	getCustomerParams.AccountNumber = accountNumber
	customerObj, err := customerService.GetCustomer(*getCustomerParams)
	if err != nil {
		return diag.FromErr(err)
	}

	// Retrieve Customer User object from API
	getCustomerUserParams := customer.NewGetCustomerUserParams()
	getCustomerUserParams.Customer = *customerObj
	getCustomerUserParams.CustomerUserID = customerUserID
	customerUserObj, err := customerService.GetCustomerUser(
		*getCustomerUserParams,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	// Update Customer User object with changed data
	customerUserObj.Address1 = customerUser.Address1
	customerUserObj.Address2 = customerUser.Address2
	customerUserObj.City = customerUser.City
	customerUserObj.Country = customerUser.Country
	customerUserObj.Email = customerUser.Email
	customerUserObj.Fax = customerUser.Fax
	customerUserObj.FirstName = customerUser.FirstName
	customerUserObj.LastName = customerUser.LastName
	customerUserObj.Mobile = customerUser.Mobile
	customerUserObj.Phone = customerUser.Phone
	customerUserObj.State = customerUser.State
	customerUserObj.Title = customerUser.Title
	customerUserObj.ZIP = customerUser.ZIP

	// Call Update Customer User API
	updateCustUserParams := customer.NewUpdateCustomerUserParams()
	updateCustUserParams.Customer = *customerObj
	updateCustUserParams.CustomerUser = *customerUserObj
	err = customerService.UpdateCustomerUser(*updateCustUserParams)

	if err != nil {
		return diag.FromErr(err)
	}

	return ResourceCustomerUserRead(ctx, d, m)
}

func ResourceCustomerUserRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	var diags diag.Diagnostics

	accountNumber := d.Get("account_number").(string)
	customerUserID, err := strconv.Atoi(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Retreiving [Customer User]: %d for [Account Number]: %s",
		customerUserID,
		accountNumber,
	)

	// Initialize Customer Service
	config := m.(**api.ClientConfig)
	customerService, err := buildCustomerService(**config)
	if err != nil {
		return diag.FromErr(err)
	}

	// Retrieve Customer object from API
	getCustomerParams := customer.NewGetCustomerParams()
	getCustomerParams.AccountNumber = accountNumber
	customerObj, err := customerService.GetCustomer(*getCustomerParams)
	if err != nil {
		return diag.FromErr(err)
	}

	// Retrieve Customer User object from API
	getCustUserParams := customer.NewGetCustomerUserParams()
	getCustUserParams.Customer = *customerObj
	getCustUserParams.CustomerUserID = customerUserID
	customerUser, err := customerService.GetCustomerUser(*getCustUserParams)

	if err != nil {
		return diag.FromErr(err)
	}

	// Process special is_admin field
	isAdminOld, isAdminNew := d.GetChange("is_admin")

	if isAdminOld.(bool) != isAdminNew.(bool) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary: "note: is_admin is a write-only property; modifications " +
				"will be ignored",
		})
	}

	// Update Terraform state with retrieved Customer User data
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
	d.Set("zip", customerUser.ZIP)
	d.Set("last_login_date", customerUser.LastLoginDate)

	return diags
}

func ResourceCustomerUserDelete(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	var diags diag.Diagnostics

	accountNumber := d.Get("account_number").(string)
	customerUserID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.Get("is_admin").(bool) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary: fmt.Sprintf(
				"[Customer User]: %d for [Account Number]: %s is an admin "+
					"user and cannot be deleted. This user will only be "+
					"deleted if the parent Customer is deleted.",
				customerUserID,
				accountNumber,
			),
		})
		return diags
	}

	log.Printf(
		"[INFO] Deleting [Customer User]: %d for [Account Number]: %s",
		customerUserID,
		accountNumber,
	)

	// Initialize Customer Service
	config := m.(**api.ClientConfig)
	customerService, err := buildCustomerService(**config)
	if err != nil {
		return diag.FromErr(err)
	}

	// Retrieve Customer object from API
	getCustomerParams := customer.NewGetCustomerParams()
	getCustomerParams.AccountNumber = accountNumber
	customerObj, err := customerService.GetCustomer(*getCustomerParams)
	if err != nil {
		return diag.FromErr(err)
	}

	// Retrieve Customer User object from API
	getCustUserParams := customer.NewGetCustomerUserParams()
	getCustUserParams.Customer = *customerObj
	getCustUserParams.CustomerUserID = customerUserID
	customerUser, err := customerService.GetCustomerUser(*getCustUserParams)
	if err != nil {
		return diag.FromErr(err)
	}

	// Call Delete Customer User API
	deleteCustUserParams := customer.NewDeleteCustomerUserParams()
	deleteCustUserParams.Customer = *customerObj
	deleteCustUserParams.CustomerUser = *customerUser
	err = customerService.DeleteCustomerUser(*deleteCustUserParams)
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
		ZIP:       d.Get("zip").(string),
	}
}
