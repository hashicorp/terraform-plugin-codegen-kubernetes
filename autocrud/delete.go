package autocrud

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sschema "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
)

const waitForDeletionSleepTime = 1 * time.Second

func Delete(ctx context.Context, clientGetter KubernetesClientGetter, kind, apiVersion string, req resource.DeleteRequest, waitForDeletion bool) error {
	client, err := clientGetter.DynamicClient()
	if err != nil {
		return err
	}
	discoveryClient, err := clientGetter.DiscoveryClient()
	if err != nil {
		return err
	}
	agr, err := restmapper.GetAPIGroupResources(discoveryClient)
	if err != nil {
		return err
	}
	gvk := k8sschema.FromAPIVersionAndKind(apiVersion, kind)
	restMapper := restmapper.NewDiscoveryRESTMapper(agr)
	mapping, err := restMapper.RESTMapping(gvk.GroupKind(), apiVersion)
	if err != nil {
		return err
	}

	var id string
	req.State.GetAttribute(ctx, path.Root("id"), &id)
	name, namespace := parseID(id)

	var resourceInterface dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		if namespace == "" {
			namespace = "default"
		}
		resourceInterface = client.Resource(mapping.Resource).Namespace(namespace)
	} else {
		resourceInterface = client.Resource(mapping.Resource)
	}

	err = resourceInterface.Delete(ctx, name, v1.DeleteOptions{})
	if err != nil {
		return err
	}

	// TODO move this into its own function called waitForDeletion
	if waitForDeletion {
		// TODO could look at usiung resourceInterface.Watch here instead of polling
		tflog.Debug(ctx, "Waiting for resource to be deleted")

		for {
			_, err := resourceInterface.Get(ctx, name, v1.GetOptions{})
			if err != nil {
				if errors.IsNotFound(err) {
					return nil
				} else {
					return err
				}
			}

			if err := ctx.Err(); err != nil {
				return err
			}

			if deadline, ok := ctx.Deadline(); ok {
				if time.Now().After(deadline) {
					return fmt.Errorf("timed out waiting for deletion")
				}
			}

			time.Sleep(waitForDeletionSleepTime)
		}
	}

	return nil
}
