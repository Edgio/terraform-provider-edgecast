// Copyright 2021 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.
package data

import (
	"github.com/EdgeCast/ec-sdk-go/edgecast"
	"github.com/EdgeCast/ec-sdk-go/edgecast/origin"
	"github.com/EdgeCast/ec-sdk-go/edgecast/shared/enums"
	"terraform-provider-edgecast/test/integration/cmd/populate/internal"
)

func createOriginData(cfg edgecast.SDKConfig) (id int) {
	svc := internal.Check(origin.New(cfg))
	id = createOrigin(svc)
	return
}

func createOrigin(svc *origin.OriginService) int {
	params := origin.AddOriginParams{
		AccountNumber: account(),
		MediaTypeID:   enums.HttpLarge,
		Origin: origin.Origin{
			DirectoryName:   unique("www"),
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
