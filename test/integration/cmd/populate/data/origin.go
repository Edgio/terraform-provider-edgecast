// Copyright 2023 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.
package data

import (
	"terraform-provider-edgecast/test/integration/cmd/populate/config"
	"terraform-provider-edgecast/test/integration/cmd/populate/internal"

	"github.com/EdgeCast/ec-sdk-go/edgecast/origin"
	"github.com/EdgeCast/ec-sdk-go/edgecast/shared/enums"
)

func createOriginData(cfg config.Config) OriginResult {
	svc := internal.Check(origin.New(cfg.SDKConfig))
	id := createOrigin(svc, cfg.AccountNumber)
	return OriginResult{id}
}

func createOrigin(svc *origin.OriginService, accountNumber string) int {
	params := origin.AddOriginParams{
		AccountNumber: accountNumber,
		MediaTypeID:   enums.HttpLarge,
		Origin: origin.Origin{
			DirectoryName:   internal.Unique("www"),
			FollowRedirects: false,
			HostHeader:      "home.example.com:80",
			HTTPHostnames: []origin.Hostname{
				{
					Name: "http://app.example.com:80",
				},
			},
			HTTPLoadBalancing: "RR",
			HTTPSHostnames: []origin.Hostname{
				{
					Name: "https://app.example.com:443",
				},
			},
			HTTPSLoadBalancing: "PF",
			ValidationURL:      "http://home.example.com:80/images/test.gif",
		},
	}

	return *internal.Check(svc.AddOrigin(params))
}
