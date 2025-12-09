// Copyright IBM Corp. 2024
// SPDX-License-Identifier: MPL-2.0

package autocrud

import (
	"context"
)

func Create(ctx context.Context, clientGetter KubernetesClientGetter, apiVersion, kind string, model any) error {
	return serverSideApply(ctx, clientGetter, apiVersion, kind, model)
}
