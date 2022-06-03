// Copyright 2021 Edgecast Inc., Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.
package data

import (
	"github.com/EdgeCast/ec-sdk-go/edgecast"
	"github.com/EdgeCast/ec-sdk-go/edgecast/edgecname"
	"terraform-provider-edgecast/test/integration/cmd/populate/internal"
)

func createEdgeCnameData(cfg edgecast.SDKConfig) (edgeCnameID int) {
	svc := internal.Check(edgecname.New(cfg))
	edgeCnameID = createEdgeCname(svc)
	return
}

func createEdgeCname(svc *edgecname.EdgeCnameService) int {
	params := edgecname.AddEdgeCnameParams{
		AccountNumber: account(),
		EdgeCname: edgecname.EdgeCname{
			Name:        unique("abc.asd"),
			OriginID:    -1,
			MediaTypeID: 3,
		},
	}

	return *internal.Check(svc.AddEdgeCname(params))
}
