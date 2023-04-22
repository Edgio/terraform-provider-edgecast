// Copyright 2021 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package customer

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"terraform-provider-edgecast/edgecast/internal"

	"github.com/EdgeCast/ec-sdk-go/edgecast/customer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceCustomer() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceCustomerCreate,
		ReadContext:   ResourceCustomerRead,
		UpdateContext: ResourceCustomerUpdate,
		DeleteContext: ResourceCustomerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"company_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"service_level_code": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"bandwidth_usage_limit": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0",
			},
			"data_transferred_usage_limit": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0",
			},
			"account_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"address1": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"address2": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"billing_account_tag": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"billing_address1": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"billing_address2": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"billing_city": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"billing_contact_email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"billing_contact_fax": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"billing_contact_first_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"billing_contact_last_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"billing_contact_mobile": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"billing_contact_phone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"billing_contact_title": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"billing_country": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"billing_rate_info": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"billing_state": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"billing_zip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"city": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"contact_email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"contact_fax": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"contact_first_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"contact_last_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"contact_mobile": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"contact_phone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"contact_title": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"country": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"notes": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"website": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"usage_limit_update_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"partner_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"partner_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"wholesale_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"wholesale_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"delivery_region": {
				Type:     schema.TypeInt,
				Optional: true,
			},
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
			"partner_user_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func ResourceCustomerCreate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {

	// Create Customer object
	newCustomer, err := getCustomerCreateUpdate(d)

	if err != nil {
		return diag.FromErr(err)
	}

	// Initialize Customer Service
	config := m.(internal.ProviderConfig)
	customerService, err := buildCustomerService(config)
	if err != nil {
		d.SetId("") // Terraform requires an empty ID for failed creation
		return diag.FromErr(err)
	}

	// Call Add Customer API
	addCustomerParams := customer.NewAddCustomerParams()
	addCustomerParams.Customer = *newCustomer
	accountNumber, err := customerService.AddCustomer(addCustomerParams)

	if err != nil {
		d.SetId("") // Terraform requires an empty ID for failed creation
		return diag.FromErr(err)
	}

	d.SetId(accountNumber)

	// Retrieve newly created Customer object
	getCustomerParams := customer.NewGetCustomerParams()
	getCustomerParams.AccountNumber = accountNumber
	customerObj, err := customerService.GetCustomer(*getCustomerParams)
	if err != nil {
		return diag.FromErr(err)
	}

	// Enable Services defined for newly created Customer
	if attr, ok := d.GetOk("services"); ok {
		attrList := attr.([]interface{})
		enableFlag := 1

		providedServiceIDs := make([]int, len(attrList))

		for i := range attrList {
			providedServiceIDs[i] = attrList[i].(int)
		}

		params := customer.NewUpdateCustomerServicesParams()
		params.Customer = *customerObj
		params.ServiceIDs = providedServiceIDs
		params.Status = enableFlag

		err = customerService.UpdateCustomerServices(*params)

		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Enable Delivery Regions defined for newly created Customer
	if attr, ok := d.GetOk("delivery_region"); ok {
		params := customer.NewUpdateCustomerDeliveryRegionParams()
		params.Customer = *customerObj
		params.DeliveryRegionID = attr.(int)
		err = customerService.UpdateCustomerDeliveryRegion(*params)

		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Enable Access Modules defined for newly created Customer
	if attr, ok := d.GetOk("access_modules"); ok {
		attrList := attr.([]interface{})
		enableFlag := 1

		var accessModuleIDs []int
		for _, v := range attrList {
			accessModuleIDs = append(accessModuleIDs, v.(int))
		}

		params := customer.NewUpdateCustomerAccessModuleParams()
		params.Customer = *customerObj
		params.AccessModuleIDs = accessModuleIDs
		params.Status = enableFlag
		err = customerService.UpdateCustomerAccessModule(*params)

		if err != nil {
			return diag.FromErr(err)
		}
	}

	return ResourceCustomerRead(ctx, d, m)
}

func ResourceCustomerRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	var diags diag.Diagnostics

	accountNumber := d.Id()
	fmt.Printf("GetCustomer>>[CustomerID]:%s", accountNumber)

	// Initialize Customer Service
	config := m.(internal.ProviderConfig)
	customerService, err := buildCustomerService(config)
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

	// Update Terraform state with retrieved Customer data
	d.SetId(customerObj.HexID)
	d.Set("status", customerObj.Status)
	d.Set("address1", customerObj.Address1)
	d.Set("address2", customerObj.Address2)
	d.Set("bandwidth_usage_limit", strconv.FormatInt(
		customerObj.BandwidthUsageLimit, 10,
	))
	d.Set("billing_account_tag", customerObj.BillingAccountTag)
	d.Set("billing_address1", customerObj.BillingAddress1)
	d.Set("billing_address2", customerObj.BillingAddress2)
	d.Set("billing_city", customerObj.BillingCity)
	d.Set("billing_contact_email", customerObj.BillingContactEmail)
	d.Set("billing_contact_fax", customerObj.BillingContactFax)
	d.Set("billing_contact_first_name", customerObj.BillingContactFirstName)
	d.Set("billing_contact_last_name", customerObj.BillingContactLastName)
	d.Set("billing_contact_mobile", customerObj.BillingContactMobile)
	d.Set("billing_contact_phone", customerObj.BillingContactPhone)
	d.Set("billing_contact_title", customerObj.BillingContactTitle)
	d.Set("billing_country", customerObj.BillingCountry)
	d.Set("billing_rate_info", customerObj.BillingRateInfo)
	d.Set("billing_state", customerObj.BillingState)
	d.Set("billing_zip", customerObj.BillingZIP)
	d.Set("city", customerObj.City)
	d.Set("company_name", customerObj.CompanyName)
	d.Set("contact_email", customerObj.ContactEmail)
	d.Set("contact_fax", customerObj.ContactFax)
	d.Set("contact_first_name", customerObj.ContactFirstName)
	d.Set("contact_last_name", customerObj.ContactLastName)
	d.Set("contact_mobile", customerObj.ContactMobile)
	d.Set("contact_phone", customerObj.ContactPhone)
	d.Set("contact_title", customerObj.ContactTitle)
	d.Set("country", customerObj.Country)
	d.Set("data_transferred_usage_limit", strconv.FormatInt(
		customerObj.DataTransferredUsageLimit, 10,
	))
	d.Set("notes", customerObj.Notes)
	d.Set("service_level_code", customerObj.ServiceLevelCode)
	d.Set("state", customerObj.State)
	d.Set("website", customerObj.Website)
	d.Set("zip", customerObj.ZIP)
	d.Set("usage_limit_update_date", customerObj.UsageLimitUpdateDate)
	d.Set("partner_id", customerObj.PartnerID)
	d.Set("partner_name", customerObj.PartnerName)
	d.Set("wholesale_id", customerObj.WholesaleID)
	d.Set("wholesale_name", customerObj.WholesaleName)
	d.Set("partner_user_id", customerObj.PartnerUserID)

	// Update Terraform state with retrieved Customer Services data for Customer
	getServicesParams := customer.NewGetCustomerServicesParams()
	getServicesParams.Customer = *customerObj
	if services, err := customerService.GetCustomerServices(
		*getServicesParams,
	); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "error retrieving customer services",
			Detail:   err.Error(),
		})
	} else {
		serviceIds := []int{}

		for _, s := range *services {
			// return only those services with Status = 1
			if s.Status == 1 {
				serviceIds = append(serviceIds, s.ID)
			}
		}

		// order matters for terraform state, so we'll sort
		sort.Ints(serviceIds)
		d.Set("services", serviceIds)
	}

	// Update Terraform state with retrieved Delivery Regions data for Customer
	getRegionParams := customer.NewGetCustomerDeliveryRegionParams()
	getRegionParams.Customer = *customerObj
	if deliveryRegion, err := customerService.GetCustomerDeliveryRegion(
		*getRegionParams,
	); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "error retrieving customer delivery region",
			Detail:   err.Error(),
		})
	} else {
		d.Set("delivery_region", deliveryRegion)
	}

	// Update Terraform state with retrieved Access Modules data for Customer
	getAccessModuleParams := customer.NewGetCustomerAccessModulesParams()
	getAccessModuleParams.Customer = *customerObj
	if accessModules, err := customerService.GetCustomerAccessModules(
		*getAccessModuleParams,
	); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "error retrieving customer access modules",
			Detail:   err.Error(),
		})
	} else {
		accessModuleIds := []int{}

		for _, a := range *accessModules {
			accessModuleIds = append(accessModuleIds, a.ID)
		}

		// order matters for terraform state, so we'll sort
		sort.Ints(accessModuleIds)
		d.Set("access_modules", accessModuleIds)
	}

	return diags
}

func ResourceCustomerUpdate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	// Not Yet Implemented
	return ResourceCustomerRead(ctx, d, m)
}

func ResourceCustomerDelete(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	var diags diag.Diagnostics

	accountNumber := d.Id()
	fmt.Printf("DeleteCustomer>>[CustomerID]:%s", accountNumber)

	// Initialize Customer Service
	config := m.(internal.ProviderConfig)
	customerService, err := buildCustomerService(config)
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

	// Call Delete Customer API
	deleteCustomerParams := customer.NewDeleteCustomerParams()
	deleteCustomerParams.Customer = *customerObj
	err = customerService.DeleteCustomer(*deleteCustomerParams)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func getCustomerCreateUpdate(
	d *schema.ResourceData,
) (*customer.Customer, error) {
	var bandwidthUsageLimit int64
	var dataTransferredUsageLimit int64

	// Terraform schemas only support Int32, so gotta do a little more heavy
	// lifting for int64 and int8
	if attr, ok := d.GetOk(
		"bandwidth_usage_limit",
	); ok && len(attr.(string)) > 0 {
		parsed, err := strconv.ParseInt(attr.(string), 10, 64)

		if err != nil {
			return nil, fmt.Errorf(
				"bandwidth_usage_limit should be a 64-bit integer",
			)
		}

		bandwidthUsageLimit = parsed
	}

	if attr, ok := d.GetOk(
		"data_transferred_usage_limit",
	); ok && len(attr.(string)) > 0 {
		parsed, err := strconv.ParseInt(attr.(string), 10, 64)

		if err != nil {
			return nil, fmt.Errorf(
				"data_transferred_usage_limit should be a 64-bit integer",
			)
		}

		dataTransferredUsageLimit = parsed
	}

	return &customer.Customer{
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
		PartnerUserID:             d.Get("partner_user_id").(int),
	}, nil
}
