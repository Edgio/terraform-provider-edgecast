// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package waf_bot_manager

import (
	"context"
	"errors"
	"fmt"
	"log"
	"terraform-provider-edgecast/edgecast/helper"
	"terraform-provider-edgecast/edgecast/internal"

	botmanager "github.com/EdgeCast/ec-sdk-go/edgecast/waf_bot_manager"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kr/pretty"
)

func ResourceBotManagerCreate(ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	config, ok := m.(internal.ProviderConfig)
	if !ok {
		return helper.CreationErrorf(d, "failed to load configuration")
	}

	svc, err := buildBotManagerService(config)
	if err != nil {
		return diag.FromErr(err)
	}

	customerID := d.Get("customer_id").(string)

	log.Printf("[INFO] Creating Bot Manager for Account >> %s", customerID)

	// Read from TF state.
	botManagerState, errs := ExpandBotManager(d)
	if len(errs) > 0 {
		return helper.DiagsFromErrors("error parsing bot manager", errs)
	}

	//Call API
	cparams := botmanager.NewCreateBotManagerParams()
	cparams.CustId = customerID
	cparams.BotManagerInfo = *botManagerState

	cresp, err := svc.BotManagers.CreateBotManager(cparams)
	if err != nil {
		return helper.CreationError(d, err)
	}

	botManagerID := cresp.Id
	log.Printf("[INFO] bot manager created: %# v\n", pretty.Formatter(cresp))
	log.Printf("[INFO] bot manager id: %d\n", botManagerID)

	d.SetId(*botManagerID)

	return ResourceBotManagerRead(ctx, d, m)
}

func ResourceBotManagerRead(ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {

	var diags diag.Diagnostics
	return diags
}

func ResourceBotManagerUpdate(ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {

	return ResourceBotManagerRead(ctx, d, m)
}

func ResourceBotManagerDelete(ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags
}

func ExpandBotManager(
	d *schema.ResourceData,
) (*botmanager.BotManager, []error) {
	if d == nil {
		return nil, []error{errors.New("no data to read")}
	}

	errs := make([]error, 0)
	botManagerState := &botmanager.BotManager{}

	if v, ok := d.GetOk("customer_id"); ok {
		if customerID, ok := v.(string); ok {
			botManagerState.CustomerId = &customerID
		} else {
			errs = append(errs, errors.New("customer_id not a string"))
		}
	}

	if v, ok := d.GetOk("name"); ok {
		if name, ok := v.(string); ok {
			botManagerState.Name = &name
		} else {
			errs = append(errs, errors.New("name not a string"))
		}
	}

	if v, ok := d.GetOk("bots_prod_id"); ok {
		if botsProdID, ok := v.(string); ok {
			botManagerState.BotsProdId = &botsProdID
		} else {
			errs = append(errs, errors.New("bots_prod_id not a string"))
		}
	}

	if v, ok := d.GetOk("inspect_known_bots"); ok {
		if inspectKnownBots, ok := v.(bool); ok {
			botManagerState.InspectKnownBots = &inspectKnownBots
		} else {
			errs = append(errs, errors.New("inspect_known_bots not a bool"))
		}
	}

	if v, ok := d.GetOk("last_modified_date"); ok {
		if lastModifiedDt, ok := v.(string); ok {
			botManagerState.LastModifiedDate = &lastModifiedDt
		} else {
			errs = append(errs, errors.New("last_modified_date not a string"))
		}
	}

	if v, ok := d.GetOk("last_modified_by"); ok {
		if lastModifiedBy, ok := v.(string); ok {
			botManagerState.LastModifiedBy = &lastModifiedBy
		} else {
			errs = append(errs, errors.New("last_modified_by not a string"))
		}
	}

	if v, ok := d.GetOk("spoof_bot_action_type"); ok {
		if spoofBotActionType, ok := v.(string); ok {
			botManagerState.SpoofBotActionType = &spoofBotActionType
		} else {
			errs = append(errs, errors.New("spoof_bot_action_type not a string"))
		}
	}

	if v, ok := d.GetOk("action"); ok {
		if actions, err := ExpandActions(v); err == nil {
			botManagerState.Actions = actions
		} else {
			errs = append(errs, fmt.Errorf("error parsing actions: %w", err))
		}
	}

	if v, ok := d.GetOk("exception_cookie"); ok {
		exceptionCookies, err := helper.ConvertTFCollectionToStrings(v)
		if err == nil {
			botManagerState.ExceptionCookie = exceptionCookies
		} else {
			err := fmt.Errorf(
				"error parsing exception cookie: %w",
				err)
			errs = append(errs, err)
		}
	}

	if v, ok := d.GetOk("exception_ja3"); ok {
		exceptionja3, err := helper.ConvertTFCollectionToStrings(v)
		if err == nil {
			botManagerState.ExceptionJa3 = exceptionja3
		} else {
			err := fmt.Errorf(
				"error parsing exception ja3: %w",
				err)
			errs = append(errs, err)
		}
	}

	if v, ok := d.GetOk("exception_url"); ok {
		exceptionurl, err := helper.ConvertTFCollectionToStrings(v)
		if err == nil {
			botManagerState.ExceptionUrl = exceptionurl
		} else {
			err := fmt.Errorf(
				"error parsing exception url: %w",
				err)
			errs = append(errs, err)
		}
	}

	if v, ok := d.GetOk("exception_user_agent"); ok {
		exceptionUserAgent, err := helper.ConvertTFCollectionToStrings(v)
		if err == nil {
			botManagerState.ExceptionUserAgent = exceptionUserAgent
		} else {
			err := fmt.Errorf(
				"error parsing exception user agent: %w",
				err)
			errs = append(errs, err)
		}
	}

	if v, ok := d.GetOk("known_bot"); ok {
		if knownBots, err := ExpandKnownBots(v); err == nil {
			botManagerState.KnownBots = knownBots
		} else {
			errs = append(errs, fmt.Errorf("error parsing known bots: %w", err))
		}
	}

	return botManagerState, errs
}

// ExpandActions converts the Terraform representation of
// action into the ActionObj API Model.
func ExpandActions(
	attr interface{},
) (*botmanager.ActionObj, error) {
	raw, ok := attr.([]any)
	if !ok {
		return nil, errors.New("attr was not a TypeList")
	}

	if len(raw) == 0 {
		return nil, nil
	}

	if len(raw) > 1 {
		return nil, errors.New("only one alert is allowed")
	}

	curr := raw[0].(map[string]any)

	// Empty map.
	if len(curr) == 0 {
		return nil, nil
	}

	action := botmanager.ActionObj{}

	if curr["alert"] != nil {
		alert, err :=
			ExpandAlert(curr["alert"])
		if err != nil {
			return nil, err
		}
		action.ALERT = alert
	}
	if curr["custom_response"] != nil {
		customResponse, err :=
			ExpandCustomResponse(curr["custom_response"])
		if err != nil {
			return nil, err
		}
		action.CUSTOM_RESPONSE = customResponse
	}
	if curr["block_request"] != nil {
		blockRequest, err :=
			ExpandBlockRequest(curr["block_request"])
		if err != nil {
			return nil, err
		}
		action.BLOCK_REQUEST = blockRequest
	}
	if curr["redirect_302"] != nil {
		redirect302, err :=
			ExpandRedirect302(curr["redirect_302"])
		if err != nil {
			return nil, err
		}
		action.REDIRECT302 = redirect302
	}
	if curr["browser_challenge"] != nil {
		browserChallenge, err :=
			ExpandBrowserChallenge(curr["browser_challenge"])
		if err != nil {
			return nil, err
		}
		action.BROWSER_CHALLENGE = browserChallenge
	}

	return &action, nil
}

func ExpandAlert(attr interface{}) (*botmanager.AlertAction, error) {
	raw, ok := attr.([]any)
	if !ok {
		return nil, errors.New("attr was not a TypeList")
	}

	if len(raw) == 0 {
		return nil, nil
	}

	if len(raw) > 1 {
		return nil, errors.New("only one alert is allowed")
	}

	curr := raw[0].(map[string]any)

	// Empty map.
	if len(curr) == 0 {
		return nil, nil
	}

	alert := botmanager.AlertAction{}

	id := curr["id"].(string)
	alert.Id = &id

	name := curr["name"].(string)
	alert.Name = &name

	enfType := curr["enf_type"].(string)
	alert.EnfType = &enfType

	return &alert, nil
}

func ExpandCustomResponse(attr interface{}) (*botmanager.CustomResponseAction, error) {
	raw, ok := attr.([]any)
	if !ok {
		return nil, errors.New("attr was not a TypeList")
	}

	if len(raw) == 0 {
		return nil, nil
	}

	if len(raw) > 1 {
		return nil, errors.New("only one custom_response is allowed")
	}

	curr := raw[0].(map[string]any)

	// Empty map.
	if len(curr) == 0 {
		return nil, nil
	}

	customResponse := botmanager.CustomResponseAction{}

	id := curr["id"].(string)
	customResponse.Id = &id

	name := curr["name"].(string)
	customResponse.Name = &name

	enfType := curr["enf_type"].(string)
	customResponse.EnfType = &enfType

	responseBodyBase64 := curr["response_body_base64"].(string)
	customResponse.ResponseBodyBase64 = &responseBodyBase64

	if status, ok := curr["status"].(uint32); ok {
		statusint := int32(status)
		customResponse.Status = &statusint
	}

	if curr["response_headers"] != nil {
		responseHeaders, err :=
			ExpandResponseHeaders(curr["response_headers"])
		if err != nil {
			return nil, err
		}
		customResponse.ResponseHeaders = responseHeaders
	}

	return &customResponse, nil
}

func ExpandResponseHeaders(
	attr interface{},
) (*map[string]string, error) {
	if items, ok := attr.([]interface{}); ok {
		responseHeaders := make(map[string]string, 0)

		for _, item := range items {
			curr := item.(map[string]interface{})

			responseHeaders[curr["key"].(string)] = curr["value"].(string)
		}

		return &responseHeaders, nil
	} else {
		return nil, errors.New(
			"ExpandResponseHeaders: attr input was not a []interface{}")
	}
}

func ExpandBlockRequest(attr interface{}) (*botmanager.BlockRequestAction, error) {
	raw, ok := attr.([]any)
	if !ok {
		return nil, errors.New("attr was not a TypeList")
	}

	if len(raw) == 0 {
		return nil, nil
	}

	if len(raw) > 1 {
		return nil, errors.New("only one block request is allowed")
	}

	curr := raw[0].(map[string]any)

	// Empty map.
	if len(curr) == 0 {
		return nil, nil
	}

	blockRequest := botmanager.BlockRequestAction{}

	id := curr["id"].(string)
	blockRequest.Id = &id

	name := curr["name"].(string)
	blockRequest.Name = &name

	enfType := curr["enf_type"].(string)
	blockRequest.EnfType = &enfType

	return &blockRequest, nil
}

func ExpandRedirect302(attr interface{}) (*botmanager.RedirectAction, error) {
	raw, ok := attr.([]any)
	if !ok {
		return nil, errors.New("attr was not a TypeList")
	}

	if len(raw) == 0 {
		return nil, nil
	}

	if len(raw) > 1 {
		return nil, errors.New("only one redirect action is allowed")
	}

	curr := raw[0].(map[string]any)

	// Empty map.
	if len(curr) == 0 {
		return nil, nil
	}

	redirect := botmanager.RedirectAction{}

	id := curr["id"].(string)
	redirect.Id = &id

	name := curr["name"].(string)
	redirect.Name = &name

	enfType := curr["enf_type"].(string)
	redirect.EnfType = &enfType

	url := curr["url"].(string)
	redirect.Url = &url

	return &redirect, nil
}

func ExpandBrowserChallenge(attr interface{}) (*botmanager.BrowserChallengeAction, error) {
	raw, ok := attr.([]any)
	if !ok {
		return nil, errors.New("attr was not a TypeList")
	}

	if len(raw) == 0 {
		return nil, nil
	}

	if len(raw) > 1 {
		return nil, errors.New("only one browser challenge action is allowed")
	}

	curr := raw[0].(map[string]any)

	// Empty map.
	if len(curr) == 0 {
		return nil, nil
	}

	browserChallenge := botmanager.BrowserChallengeAction{}

	id := curr["id"].(string)
	browserChallenge.Id = &id

	name := curr["name"].(string)
	browserChallenge.Name = &name

	enfType := curr["enf_type"].(string)
	browserChallenge.EnfType = &enfType

	isCustomChallenge := curr["is_custom_challenge"].(bool)
	browserChallenge.IsCustomChallenge = &isCustomChallenge

	responseBodyBase64 := curr["response_body_base64"].(string)
	browserChallenge.ResponseBodyBase64 = &responseBodyBase64

	if validfor, ok := curr["valid_for_sec"].(uint32); ok {
		validforInt := int32(validfor)
		browserChallenge.ValidForSec = &validforInt
	}

	if status, ok := curr["status"].(uint32); ok {
		statusint := int32(status)
		browserChallenge.Status = &statusint
	}

	return &browserChallenge, nil
}

// ExpandAdditionalContacts converts the Terraform representation of
// organization contacts into the OrganizationContact API Model.
func ExpandKnownBots(
	attr interface{},
) ([]botmanager.KnownBotObj, error) {
	if items, ok := attr.([]interface{}); ok {
		knownBots := make([]botmanager.KnownBotObj, 0)

		for _, item := range items {
			curr := item.(map[string]interface{})

			knownBot := botmanager.KnownBotObj{
				ActionType: curr["action_type"].(string),
				BotToken:   curr["bot_token"].(string),
			}

			knownBots = append(knownBots, knownBot)
		}

		return knownBots, nil
	} else {
		return nil, errors.New(
			"ExpandKNownBots: attr input was not a []interface{}")
	}
}
