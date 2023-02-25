// Copyright 2023 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package waf_bot_manager

import (
	"context"
	"errors"
	"fmt"
	"log"

	"terraform-provider-edgecast/edgecast/helper"
	"terraform-provider-edgecast/edgecast/internal"

	sdkbotmanager "github.com/EdgeCast/ec-sdk-go/edgecast/waf_bot_manager"
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

	// Call API
	cparams := sdkbotmanager.NewCreateBotManagerParams()
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
	config, ok := m.(internal.ProviderConfig)
	if !ok {
		return diag.Errorf("failed to load configuration")
	}

	svc, err := buildBotManagerService(config)
	if err != nil {
		return diag.FromErr(err)
	}

	customerID := d.Get("customer_id").(string)
	botManagerID := d.Id()

	log.Printf(
		"[INFO] Retrieving Bot Manager ID %s for Account >> %s",
		botManagerID,
		customerID)

	params := sdkbotmanager.NewGetBotManagerParams()
	params.BotManagerId = botManagerID
	params.CustId = customerID

	resp, err := svc.BotManagers.GetBotManager(params)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Retrieved Bot Manager: %# v\n", pretty.Formatter(resp))

	err = FlattenBotManager(d, resp)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
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

	customerID := d.Get("customer_id").(string)
	botManagerID := d.Id()

	log.Printf("[INFO] Deleting WAF Bot Manager ID %s for Account >> %s",
		botManagerID,
		customerID,
	)

	config := m.(internal.ProviderConfig)

	svc, err := buildBotManagerService(config)
	if err != nil {
		return diag.FromErr(err)
	}

	params := sdkbotmanager.NewDeleteBotManagerParams()
	params.CustId = customerID
	params.BotManagerId = botManagerID

	err = svc.BotManagers.DeleteBotManager(params)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Successfully deleted WAF Bot Manager: %s", botManagerID)
	d.SetId("")

	return diags
}

func ExpandBotManager(
	d *schema.ResourceData,
) (*sdkbotmanager.BotManager, []error) {
	if d == nil {
		return nil, []error{errors.New("no data to read")}
	}

	errs := make([]error, 0)
	botManagerState := &sdkbotmanager.BotManager{}

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
			errs = append(
				errs,
				errors.New("spoof_bot_action_type not a string"))
		}
	}

	if v, ok := d.GetOk("actions"); ok {
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
// actions into the ActionObj API Model.
func ExpandActions(
	attr interface{},
) (*sdkbotmanager.ActionObj, error) {
	raw, ok := attr.([]any)
	if !ok {
		return nil, errors.New("attr was not a TypeList")
	}

	if len(raw) == 0 {
		return nil, nil
	}

	if len(raw) > 1 {
		return nil, errors.New("only one action is allowed")
	}

	curr := raw[0].(map[string]any)

	// Empty map.
	if len(curr) == 0 {
		return nil, nil
	}

	action := sdkbotmanager.ActionObj{}

	if curr["alert"] != nil {
		alert, err := ExpandAlert(curr["alert"])
		if err != nil {
			return nil, err
		}
		action.ALERT = alert
	}
	if curr["custom_response"] != nil {
		customResponse, err := ExpandCustomResponse(curr["custom_response"])
		if err != nil {
			return nil, err
		}
		action.CUSTOM_RESPONSE = customResponse
	}
	if curr["block_request"] != nil {
		blockRequest, err := ExpandBlockRequest(curr["block_request"])
		if err != nil {
			return nil, err
		}
		action.BLOCK_REQUEST = blockRequest
	}
	if curr["redirect_302"] != nil {
		redirect302, err := ExpandRedirect302(curr["redirect_302"])
		if err != nil {
			return nil, err
		}
		action.REDIRECT302 = redirect302
	}
	if curr["browser_challenge"] != nil {
		bc, err := ExpandBrowserChallenge(curr["browser_challenge"])
		if err != nil {
			return nil, err
		}
		action.BROWSER_CHALLENGE = bc
	}

	return &action, nil
}

func ExpandAlert(attr interface{}) (*sdkbotmanager.AlertAction, error) {
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

	alert := sdkbotmanager.AlertAction{}

	id := curr["id"].(string)
	alert.Id = &id

	name := curr["name"].(string)
	alert.Name = &name

	return &alert, nil
}

func ExpandCustomResponse(
	attr interface{},
) (*sdkbotmanager.CustomResponseAction, error) {
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

	customResponse := sdkbotmanager.CustomResponseAction{}

	id := curr["id"].(string)
	customResponse.Id = &id

	name := curr["name"].(string)
	customResponse.Name = &name

	responseBodyBase64 := curr["response_body_base64"].(string)
	customResponse.ResponseBodyBase64 = &responseBodyBase64

	status := curr["status"].(int)
	statusint32 := int32(status)
	customResponse.Status = &statusint32

	if curr["response_headers"] != nil {
		responseHeaders := helper.ConvertToStringMapPointer(
			curr["response_headers"], true,
		)
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

func ExpandBlockRequest(
	attr interface{},
) (*sdkbotmanager.BlockRequestAction, error) {
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

	blockRequest := sdkbotmanager.BlockRequestAction{}

	id := curr["id"].(string)
	blockRequest.Id = &id

	name := curr["name"].(string)
	blockRequest.Name = &name

	return &blockRequest, nil
}

func ExpandRedirect302(
	attr interface{},
) (*sdkbotmanager.RedirectAction, error) {
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

	redirect := sdkbotmanager.RedirectAction{}

	id := curr["id"].(string)
	redirect.Id = &id

	name := curr["name"].(string)
	redirect.Name = &name

	url := curr["url"].(string)
	redirect.Url = &url

	return &redirect, nil
}

func ExpandBrowserChallenge(
	attr interface{},
) (*sdkbotmanager.BrowserChallengeAction, error) {
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

	browserChallenge := sdkbotmanager.BrowserChallengeAction{}

	id := curr["id"].(string)
	browserChallenge.Id = &id

	name := curr["name"].(string)
	browserChallenge.Name = &name

	isCustomChallenge := curr["is_custom_challenge"].(bool)
	browserChallenge.IsCustomChallenge = &isCustomChallenge

	responseBodyBase64 := curr["response_body_base64"].(string)
	browserChallenge.ResponseBodyBase64 = &responseBodyBase64

	validfor := curr["valid_for_sec"].(int)
	validforInt32 := int32(validfor)
	browserChallenge.ValidForSec = &validforInt32

	status := curr["status"].(int)
	statusInt32 := int32(status)
	browserChallenge.Status = &statusInt32

	return &browserChallenge, nil
}

func ExpandKnownBots(
	attr interface{},
) ([]sdkbotmanager.KnownBotObj, error) {
	if items, ok := attr.([]interface{}); ok {
		knownBots := make([]sdkbotmanager.KnownBotObj, 0)

		for _, item := range items {
			curr := item.(map[string]interface{})

			knownBot := sdkbotmanager.KnownBotObj{
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

func FlattenBotManager(
	d *schema.ResourceData,
	bm *sdkbotmanager.BotManager,
) error {
	if bm == nil {
		return fmt.Errorf("bot manager is nil")
	}

	if name, ok := bm.GetNameOk(); ok {
		d.Set("name", *name)
	}

	if botsProdID, ok := bm.GetBotsProdIdOk(); ok {
		d.Set("bots_prod_id", *botsProdID)
	}

	if actions, ok := bm.GetActionsOk(); ok {
		d.Set("actions", FlattenActions(*actions))
	}

	if exceptionCookie, ok := bm.GetExceptionCookieOk(); ok {
		d.Set("exception_cookie", exceptionCookie)
	}

	if exceptionJa3, ok := bm.GetExceptionJa3Ok(); ok {
		d.Set("exception_ja3", exceptionJa3)
	}

	if exceptionURL, ok := bm.GetExceptionUrlOk(); ok {
		d.Set("exception_url", exceptionURL)
	}

	if exceptionUserAgent, ok := bm.GetExceptionUserAgentOk(); ok {
		d.Set("exception_user_agent", exceptionUserAgent)
	}

	if inspectKnownBots, ok := bm.GetInspectKnownBotsOk(); ok {
		d.Set("inspect_known_bots", *inspectKnownBots)
	}

	if knownBots, ok := bm.GetKnownBotsOk(); ok {
		d.Set("known_bots", FlattenKnownBots(knownBots))
	}

	if lastModifiedDate, ok := bm.GetLastModifiedDateOk(); ok {
		d.Set("last_modified_date", *lastModifiedDate)
	}

	if lastModifiedBy, ok := bm.GetLastModifiedByOk(); ok {
		d.Set("last_modified_by", *lastModifiedBy)
	}

	if spoofBotActionType, ok := bm.GetSpoofBotActionTypeOk(); ok {
		d.Set("spoof_bot_action_type", *spoofBotActionType)
	}

	return nil
}

func FlattenKnownBots(kb []sdkbotmanager.KnownBotObj) interface{} {
	flattened := make([]map[string]interface{}, len(kb), len(kb))

	for i, v := range kb {
		fb := make(map[string]interface{})

		if at, ok := v.GetActionTypeOk(); ok {
			fb["action_type"] = *at
		}

		if bt, ok := v.GetBotTokenOk(); ok {
			fb["bot_token"] = *bt
		}

		flattened[i] = fb
	}

	return flattened
}

func FlattenActions(action sdkbotmanager.ActionObj) interface{} {
	m := make(map[string]interface{})

	if alert, ok := action.GetALERTOk(); ok {
		m["alert"] = FlattenAlert(*alert)
	}

	if cr, ok := action.GetCUSTOM_RESPONSEOk(); ok {
		m["custom_response"] = FlattenCustomResponse(*cr)
	}

	if br, ok := action.GetBLOCK_REQUESTOk(); ok {
		m["block_request"] = FlattenBlockRequest(*br)
	}

	if r, ok := action.GetREDIRECT302Ok(); ok {
		m["redirect_302"] = FlattenRedirect(*r)
	}

	if bc, ok := action.GetBROWSER_CHALLENGEOk(); ok {
		m["browser_challenge"] = FlattenBrowserChallenge(*bc)
	}

	return []map[string]interface{}{m}
}

func FlattenAlert(alert sdkbotmanager.AlertAction) interface{} {
	flattened := make(map[string]interface{})

	if id, ok := alert.GetIdOk(); ok {
		flattened["id"] = *id
	}

	if name, ok := alert.GetNameOk(); ok {
		flattened["name"] = *name
	}

	return flattened
}

func FlattenCustomResponse(cr sdkbotmanager.CustomResponseAction) interface{} {
	flattened := make(map[string]interface{})

	if id, ok := cr.GetIdOk(); ok {
		flattened["id"] = *id
	}

	if name, ok := cr.GetNameOk(); ok {
		flattened["name"] = *name
	}

	if responseBodyBase64, ok := cr.GetResponseBodyBase64Ok(); ok {
		flattened["response_body_base64"] = *responseBodyBase64
	}

	if status, ok := cr.GetStatusOk(); ok {
		flattened["status"] = *status
	}

	if responseHeaders, ok := cr.GetResponseHeadersOk(); ok {
		flattened["response_headers"] = *responseHeaders
	}

	return flattened
}

func FlattenBlockRequest(br sdkbotmanager.BlockRequestAction) interface{} {
	flattened := make(map[string]interface{})

	if id, ok := br.GetIdOk(); ok {
		flattened["id"] = *id
	}

	if name, ok := br.GetNameOk(); ok {
		flattened["name"] = *name
	}

	return flattened
}

func FlattenRedirect(r sdkbotmanager.RedirectAction) interface{} {
	flattened := make(map[string]interface{})

	if id, ok := r.GetIdOk(); ok {
		flattened["id"] = *id
	}

	if name, ok := r.GetNameOk(); ok {
		flattened["name"] = *name
	}

	if url, ok := r.GetUrlOk(); ok {
		flattened["url"] = *url
	}

	return flattened
}

func FlattenBrowserChallenge(
	bc sdkbotmanager.BrowserChallengeAction,
) interface{} {
	flattened := make(map[string]interface{})

	if id, ok := bc.GetIdOk(); ok {
		flattened["id"] = *id
	}

	if name, ok := bc.GetNameOk(); ok {
		flattened["name"] = *name
	}

	if isCustomChallenge, ok := bc.GetIsCustomChallengeOk(); ok {
		flattened["is_custom_challenge"] = *isCustomChallenge
	}

	if responseBodyBase64, ok := bc.GetResponseBodyBase64Ok(); ok {
		flattened["response_body_base64"] = *responseBodyBase64
	}

	if validForSec, ok := bc.GetValidForSecOk(); ok {
		flattened["valid_for_sec"] = *validForSec
	}

	if status, ok := bc.GetStatusOk(); ok {
		flattened["status"] = *status
	}

	return flattened
}
