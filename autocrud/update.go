// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package autocrud

import (
	"context"
)

func Update(ctx context.Context, clientGetter KubernetesClientGetter, kind, apiVersion string, model any) error {
	return serverSideApply(ctx, clientGetter, apiVersion, kind, model)
}
