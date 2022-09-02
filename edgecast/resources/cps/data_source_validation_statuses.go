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

func DataSourceValidationStatuses() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceValidationStatusesRead,
		Schema:      namedEntitySchema("Validation Status"),
	}
}

func DataSourceValidationStatusesRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	return DataSourceNamedEntityRead(ctx, d, m, callGetValidationStatuses)
}

func callGetValidationStatuses(
	svc *cps.CpsService,
	d *schema.ResourceData,
) (*models.HyperionCollectionNamedEntity, error) {
	params := appendix.NewAppendixGetValidationStatusesParams()

	resp, err := svc.Appendix.AppendixGetValidationStatuses(params)
	if err != nil {
		return nil, err
	}

	return &resp.HyperionCollectionNamedEntity, nil
}
