// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.
package test

import (
	"net/url"
	"terraform-provider-vmp/vmp/api"
	"testing"
	//tfPlan "github.com/hashicorp/terraform/plans/planfile"
)

func TestClientConfig(t *testing.T) {
	apiToken := "<apiToken>"
	accountNumber := "C1B6"
	idsClientID := "<idsClientID>"
	idsClientSecret := "<idsClientSecret>"
	idsScope := "<idsScope>"
	apiURL := "http://dev-api.edgecast.com"
	idsURL := "id-dev.vdms.io"

	myConfig, err := api.NewClientConfig(apiToken, accountNumber, idsClientID, idsClientSecret, idsScope, apiURL, idsURL)

	if err != nil {
		t.Errorf("Can't use ClientConfig constructor. Error details:%s", err)
	}
	if myConfig.APIToken != apiToken {
		t.Errorf("config has different apiToken value from input: input value:%s config value: %s", apiToken, myConfig.APIURL)
	}
	if myConfig.AccountNumber != accountNumber {
		t.Errorf("config has different apiToken value from input: input value:%s config value: %s", accountNumber, myConfig.AccountNumber)
	}
	if myConfig.IdsClientID != idsClientID {
		t.Errorf("config has different apiToken value from input: input value:%s config value: %s", idsClientID, myConfig.IdsClientID)
	}
	if myConfig.IdsClientSecret != idsClientSecret {
		t.Errorf("config has different apiToken value from input: input value:%s config value: %s", idsClientSecret, myConfig.IdsClientSecret)
	}
	if myConfig.IdsScope != idsScope {
		t.Errorf("config has different apiToken value from input: input value:%s config value: %s", idsScope, myConfig.IdsScope)
	}
	api, _ := url.Parse(apiURL)
	if myConfig.APIURL.Host != api.Host {
		t.Errorf("config has different apiToken value from input: input value:%s config value: %s", apiURL, myConfig.APIURL)
	}
	ids, _ := url.Parse((idsURL))
	if myConfig.IdsURL.Host != ids.Host {
		t.Errorf("config has different apiToken value from input: input value:%s config value: %s", idsURL, myConfig.IdsURL)
	}
}
