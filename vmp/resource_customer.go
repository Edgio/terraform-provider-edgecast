// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package vmp

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"terraform-provider-vmp/vmp/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCustomer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomerCreate,
		ReadContext:   resourceCustomerRead,
		UpdateContext: resourceCustomerUpdate,
		DeleteContext: resourceCustomerDelete,

		Schema: map[string]*schema.Schema{
			"company_name":                 {Type: schema.TypeString, Required: true},
			"status":                       {Type: schema.TypeInt, Computed: true},
			"service_level_code":           {Type: schema.TypeString, Required: true},
			"bandwidth_usage_limit":        {Type: schema.TypeString, Optional: true, Default: "0"},
			"data_transferred_usage_limit": {Type: schema.TypeString, Optional: true, Default: "0"},
			"account_id":                   {Type: schema.TypeString, Optional: true},
			"address1":                     {Type: schema.TypeString, Optional: true},
			"address2":                     {Type: schema.TypeString, Optional: true},
			"billing_account_tag":          {Type: schema.TypeString, Optional: true},
			"billing_address1":             {Type: schema.TypeString, Optional: true},
			"billing_address2":             {Type: schema.TypeString, Optional: true},
			"billing_city":                 {Type: schema.TypeString, Optional: true},
			"billing_contact_email":        {Type: schema.TypeString, Optional: true},
			"billing_contact_fax":          {Type: schema.TypeString, Optional: true},
			"billing_contact_first_name":   {Type: schema.TypeString, Optional: true},
			"billing_contact_last_name":    {Type: schema.TypeString, Optional: true},
			"billing_contact_mobile":       {Type: schema.TypeString, Optional: true},
			"billing_contact_phone":        {Type: schema.TypeString, Optional: true},
			"billing_contact_title":        {Type: schema.TypeString, Optional: true},
			"billing_country":              {Type: schema.TypeString, Optional: true},
			"billing_rate_info":            {Type: schema.TypeString, Optional: true},
			"billing_state":                {Type: schema.TypeString, Optional: true},
			"billing_zip":                  {Type: schema.TypeString, Optional: true},
			"city":                         {Type: schema.TypeString, Optional: true},
			"contact_email":                {Type: schema.TypeString, Optional: true},
			"contact_fax":                  {Type: schema.TypeString, Optional: true},
			"contact_first_name":           {Type: schema.TypeString, Optional: true},
			"contact_last_name":            {Type: schema.TypeString, Optional: true},
			"contact_mobile":               {Type: schema.TypeString, Optional: true},
			"contact_phone":                {Type: schema.TypeString, Optional: true},
			"contact_title":                {Type: schema.TypeString, Optional: true},
			"country":                      {Type: schema.TypeString, Optional: true},
			"notes":                        {Type: schema.TypeString, Optional: true},
			"state":                        {Type: schema.TypeString, Optional: true},
			"website":                      {Type: schema.TypeString, Optional: true},
			"zip":                          {Type: schema.TypeString, Optional: true},
			"usage_limit_update_date":      {Type: schema.TypeString, Computed: true},
			"partner_id":                   {Type: schema.TypeInt, Computed: true},
			"partner_name":                 {Type: schema.TypeString, Computed: true},
			"wholesale_id":                 {Type: schema.TypeInt, Computed: true},
			"wholesale_name":               {Type: schema.TypeString, Computed: true},
			"delivery_region":              {Type: schema.TypeInt, Optional: true},
			"services": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"access_modules": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func getCustomerCreateUpdate(d *schema.ResourceData) (*api.CustomerCreateUpdate, error) {
	var bandwidthUsageLimit int64
	var dataTransferredUsageLimit int64

	// Terraform schemas only support Int32, so gotta do a little more heavy lifting for int64 and int8
	if attr, ok := d.GetOk("bandwidth_usage_limit"); ok && len(attr.(string)) > 0 {
		parsed, err := strconv.ParseInt(attr.(string), 10, 64)

		if err != nil {
			return nil, fmt.Errorf("bandwidth_usage_limit should be a 64-bit integer")
		}

		bandwidthUsageLimit = parsed
	}

	if attr, ok := d.GetOk("data_transferred_usage_limit"); ok && len(attr.(string)) > 0 {
		parsed, err := strconv.ParseInt(attr.(string), 10, 64)

		if err != nil {
			return nil, fmt.Errorf("data_transferred_usage_limit should be a 64-bit integer")
		}

		dataTransferredUsageLimit = parsed
	}

	return &api.CustomerCreateUpdate{
		CompanyName:               d.Get("company_name").(string),
		Status:                    1, // not user configurable
		AccountID:                 d.Get("account_id").(string),
		Address1:                  d.Get("address1").(string),
		Address2:                  d.Get("address2").(string),
		BandwidthUsageLimit:       bandwidthUsageLimit,
		BillingAccountTag:         d.Get("billing_account_tag").(string),
		BillingAddress1:           d.Get("billing_address1").(string),
		BillingAddress2:           d.Get("billing_address2").(string),
		BillingCity:               d.Get("billing_city").(string),
		BillingContactEmail:       d.Get("billing_contact_email").(string),
		BillingContactFax:         d.Get("billing_contact_fax").(string),
		BillingContactFirstName:   d.Get("billing_contact_first_name").(string),
		BillingContactLastName:    d.Get("billing_contact_last_name").(string),
		BillingContactMobile:      d.Get("billing_contact_mobile").(string),
		BillingContactPhone:       d.Get("billing_contact_phone").(string),
		BillingContactTitle:       d.Get("billing_contact_title").(string),
		BillingCountry:            d.Get("billing_country").(string),
		BillingRateInfo:           d.Get("billing_rate_info").(string),
		BillingState:              d.Get("billing_state").(string),
		BillingZIP:                d.Get("billing_zip").(string),
		City:                      d.Get("city").(string),
		ContactEmail:              d.Get("contact_email").(string),
		ContactFax:                d.Get("contact_fax").(string),
		ContactFirstName:          d.Get("contact_first_name").(string),
		ContactLastName:           d.Get("contact_last_name").(string),
		ContactMobile:             d.Get("contact_mobile").(string),
		ContactPhone:              d.Get("contact_phone").(string),
		ContactTitle:              d.Get("contact_title").(string),
		Country:                   d.Get("country").(string),
		DataTransferredUsageLimit: dataTransferredUsageLimit,
		Notes:                     d.Get("notes").(string),
		ServiceLevelCode:          d.Get("service_level_code").(string),
		State:                     d.Get("state").(string),
		Website:                   d.Get("website").(string),
		ZIP:                       d.Get("zip").(string),
	}, nil
}

func resourceCustomerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(**api.ClientConfig)

	payload, err := getCustomerCreateUpdate(d)

	if err != nil {
		return diag.FromErr(err)
	}

	customerAPIClient := api.NewCustomerAPIClient(*config)

	accountNumber, err := customerAPIClient.AddCustomer(payload)

	if err != nil {
		// Terraform requires an empty ID for failed creation
		d.SetId("")
		return diag.FromErr(err)
	}

	d.SetId(accountNumber)

	if attr, ok := d.GetOk("services"); ok {
		attrList := attr.([]interface{})

		providedServiceIDs := make([]int, len(attrList))

		for i := range attrList {
			providedServiceIDs[i] = attrList[i].(int)
		}

		err = customerAPIClient.UpdateCustomerServices(accountNumber, providedServiceIDs, 1)
	}

	if attr, ok := d.GetOk("delivery_region"); ok {
		deliveryRegion := attr.(int)
		err = customerAPIClient.UpdateCustomerDeliveryRegion(accountNumber, deliveryRegion)

		if err != nil {
			return diag.FromErr(err)
		}
	}

	if attr, ok := d.GetOk("access_modules"); ok {
		attrList := attr.([]interface{})

		for _, v := range attrList {
			accessModuleID := v.(int)

			err = customerAPIClient.UpdateCustomerAccessModule(accountNumber, accessModuleID)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceCustomerRead(ctx, d, m)
}

func resourceCustomerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(**api.ClientConfig)
	accountNumber := d.Id()
	fmt.Printf("GetCustomer>>[CustomerID]:%s", accountNumber)
	(*config).AccountNumber = accountNumber

	customerAPIClient := api.NewCustomerAPIClient(*config)

	resp, err := customerAPIClient.GetCustomer(accountNumber)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.SetId(resp.HexID)
	d.Set("status", resp.Status)
	d.Set("address1", resp.Address1)
	d.Set("address2", resp.Address2)
	d.Set("bandwidth_usage_limit", strconv.FormatInt(resp.BandwidthUsageLimit, 10))
	d.Set("billing_account_tag", resp.BillingAccountTag)
	d.Set("billing_address1", resp.BillingAddress1)
	d.Set("billing_address2", resp.BillingAddress2)
	d.Set("billing_city", resp.BillingCity)
	d.Set("billing_contact_email", resp.BillingContactEmail)
	d.Set("billing_contact_fax", resp.BillingContactFax)
	d.Set("billing_contact_first_name", resp.BillingContactFirstName)
	d.Set("billing_contact_last_name", resp.BillingContactLastName)
	d.Set("billing_contact_mobile", resp.BillingContactMobile)
	d.Set("billing_contact_phone", resp.BillingContactPhone)
	d.Set("billing_contact_title", resp.BillingContactTitle)
	d.Set("billing_country", resp.BillingCountry)
	d.Set("billing_rate_info", resp.BillingRateInfo)
	d.Set("billing_state", resp.BillingState)
	d.Set("billing_zip", resp.BillingZIP)
	d.Set("city", resp.City)
	d.Set("company_name", resp.CompanyName)
	d.Set("contact_email", resp.ContactEmail)
	d.Set("contact_fax", resp.ContactFax)
	d.Set("contact_first_name", resp.ContactFirstName)
	d.Set("contact_last_name", resp.ContactLastName)
	d.Set("contact_mobile", resp.ContactMobile)
	d.Set("contact_phone", resp.ContactPhone)
	d.Set("contact_title", resp.ContactTitle)
	d.Set("country", resp.Country)
	d.Set("data_transferred_usage_limit", strconv.FormatInt(resp.DataTransferredUsageLimit, 10))
	d.Set("notes", resp.Notes)
	d.Set("service_level_code", resp.ServiceLevelCode)
	d.Set("state", resp.State)
	d.Set("website", resp.Website)
	d.Set("zip", resp.ZIP)
	d.Set("usage_limit_update_date", resp.UsageLimitUpdateDate)
	d.Set("partner_id", resp.PartnerID)
	d.Set("partner_name", resp.PartnerName)
	d.Set("wholesale_id", resp.WholesaleID)
	d.Set("wholesale_name", resp.WholesaleName)

	if services, err := customerAPIClient.GetCustomerServices(accountNumber); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "error retrieving customer services",
			Detail:   err.Error(),
		})
	} else {
		serviceIds := []int{}

		for _, s := range services {
			// return only those services with Status = 1
			if s.Status == 1 {
				serviceIds = append(serviceIds, s.ID)
			}
		}

		// order matters for terraform state, so we'll sort
		sort.Ints(serviceIds)
		d.Set("services", serviceIds)
	}

	if deliveryRegion, err := customerAPIClient.GetCustomerDeliveryRegion(accountNumber); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "error retrieving customer delivery region",
			Detail:   err.Error(),
		})
	} else {
		d.Set("delivery_region", deliveryRegion)
	}

	// Uncomment below when new API endpoint is up on production
	// if accessModules, err := customerAPIClient.GetCustomerAccessModules(accountNumber); err != nil {
	// 	diags = append(diags, diag.Diagnostic{
	// 		Severity: diag.Error,
	// 		Summary:  "error retrieving customer access modules",
	// 		Detail:   err.Error(),
	// 	})
	// } else {
	// 	accessModuleIds := []int{}

	// 	for _, a := range *accessModules {
	// 		accessModuleIds = append(accessModuleIds, a.ID)
	// 	}

	// 	// order matters for terraform state, so we'll sort
	// 	sort.Ints(accessModuleIds)
	// 	d.Set("access_modules", accessModuleIds)
	// }

	return diags
}

func resourceCustomerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Not Yet Implemented
	return resourceCustomerRead(ctx, d, m)
}

func resourceCustomerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(**api.ClientConfig)
	accountNumber := d.Id()
	(*config).AccountNumber = accountNumber

	customerAPIClient := api.NewCustomerAPIClient(*config)

	err := customerAPIClient.DeleteCustomer(accountNumber)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func selectServiceIDs(arr []api.Service) []int {
	serviceIDs := make([]int, len(arr))

	for i := range arr {
		serviceIDs[i] = arr[i].ID
	}

	return serviceIDs
}
