// Copyright 2022 Edgecast Inc., Licensed under the terms of the Apache 2.0
// license. See LICENSE file in project root for terms.

package cps

import (
	"context"

	"github.com/EdgeCast/ec-sdk-go/edgecast/cps"
	"github.com/EdgeCast/ec-sdk-go/edgecast/cps/appendix"
	"github.com/EdgeCast/ec-sdk-go/edgecast/cps/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceCertValidationLevels() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceCertValidationLevelsRead,
		Schema:      namedEntitySchema("Certificate Validation Level"),
	}
}

func DataSourceCertValidationLevelsRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	return DataSourceNamedEntityRead(ctx, d, m, callGetCertValidationLevels)
}

func callGetCertValidationLevels(
	svc *cps.CpsService,
	d *schema.ResourceData,
) (*models.HyperionCollectionNamedEntity, error) {
	params := appendix.NewAppendixGetValidationTypesParams()

	resp, err := svc.Appendix.AppendixGetValidationTypes(params)
	if err != nil {
		return nil, err
	}

	return &resp.HyperionCollectionNamedEntity, nil
}
