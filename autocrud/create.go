package autocrud

import (
	"context"
)

func Create(ctx context.Context, clientGetter KubernetesClientGetter, apiVersion, kind string, model any) error {
	return serverSideApply(ctx, clientGetter, apiVersion, kind, model)
}
