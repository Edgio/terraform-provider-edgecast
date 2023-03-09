// Copyright 2023 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.
package data

import (
	"terraform-provider-edgecast/test/integration/cmd/populate/config"
	"terraform-provider-edgecast/test/integration/cmd/populate/internal"

	"github.com/EdgeCast/ec-sdk-go/edgecast/edgecname"
)

func createEdgeCnameData(cfg config.Config) CNAMEResult {
	svc := internal.Check(edgecname.New(cfg.SDKConfig))
	edgeCnameID := createEdgeCname(svc, cfg.AccountNumber)

	return CNAMEResult{edgeCnameID}
}

func createEdgeCname(
	svc *edgecname.EdgeCnameService,
	accountNumber string,
) int {
	params := edgecname.AddEdgeCnameParams{
		AccountNumber: accountNumber,
		EdgeCname: edgecname.EdgeCname{
			Name:        internal.Unique("abc.asd"),
			OriginID:    -1,
			MediaTypeID: 3,
		},
	}

	return *internal.Check(svc.AddEdgeCname(params))
}
