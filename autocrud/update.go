package autocrud

import (
	"context"
)

func Update(ctx context.Context, clientGetter KubernetesClientGetter, kind, apiVersion string, model any) error {
	return serverSideApply(ctx, clientGetter, apiVersion, kind, model)
}
