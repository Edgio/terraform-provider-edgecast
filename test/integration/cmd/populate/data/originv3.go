// Copyright 2023 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.
package data

import (
	"terraform-provider-edgecast/test/integration/cmd/populate/config"
	"terraform-provider-edgecast/test/integration/cmd/populate/internal"

	"github.com/EdgeCast/ec-sdk-go/edgecast/originv3"
	"github.com/EdgeCast/ec-sdk-go/edgecast/shared/enums"
)

func createOriginV3Data(cfg config.Config) OriginV3Result {
	svc := internal.Check(originv3.New(cfg.SDKConfig))
	groupID := createOriginV3Group(svc)
	originID := createOriginV3(svc, groupID)

	return OriginV3Result{
		GroupIdV3:  groupID,
		OriginIdV3: originID,
	}
}

func createOriginV3Group(svc *originv3.Service) int32 {
	tlsSettings := originv3.TlsSettings{
		PublicKeysToVerify: []string{
			"ff8b4a82b08ea5f7be124e6b4363c00d7462655f",
			"c571398b01fce46a8a177abdd6174dfee6137358",
		},
	}

	tlsSettings.SetAllowSelfSigned(false)
	tlsSettings.SetSniHostname("origin.example.com")
	grp := originv3.CustomerOriginGroupHTTPRequest{
		Name:        internal.Unique("IntegrationTestGroup"),
		TlsSettings: &tlsSettings,
	}

	grp.SetHostHeader("override-hostheader.example.com")
	grp.SetNetworkTypeId(2)          // Prefer IPv6 over IPv4
	grp.SetStrictPciCertified(false) // Allow non-PCI regions

	params := originv3.AddHttpLargeGroupParams{
		CustomerOriginGroupHTTPRequest: grp,
	}

	resp := internal.Check(svc.HttpLargeOnly.AddHttpLargeGroup(params))
	return *resp.Id
}

func createOriginV3(svc *originv3.Service, groupID int32) int32 {
	addOriginParams := originv3.NewAddOriginParams()
	addOriginParams.MediaType = enums.HttpLarge.String()
	originRequest := originv3.NewCustomerOriginRequest(
		"cdn-origin-example.com",
		false,
		groupID,
	)
	addOriginParams.CustomerOriginRequest = *originRequest

	resp := internal.Check(svc.Common.AddOrigin(addOriginParams))
	return *resp.Id
}
