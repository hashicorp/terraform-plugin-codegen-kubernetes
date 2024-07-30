// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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

func Delete(ctx context.Context, clientGetter KubernetesClientGetter, kind, apiVersion string, req resource.DeleteRequest, wait bool) error {
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
	gv, err := k8sschema.ParseGroupVersion(apiVersion)
	if err != nil {
		return err
	}

	restMapper := restmapper.NewDiscoveryRESTMapper(agr)
	mapping, err := restMapper.RESTMapping(gv.WithKind(kind).GroupKind(), gv.Version)
	if err != nil {
		return err
	}

	var id string
	req.State.GetAttribute(ctx, path.Root("id"), &id)
	namespace, name := parseID(id)

	var resourceInterface dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		if namespace == "" {
			namespace = "default"
		}
		resourceInterface = client.Resource(mapping.Resource).Namespace(namespace)
	} else {
		resourceInterface = client.Resource(mapping.Resource)
	}

	tflog.Debug(ctx, "Executing delete operation", map[string]any{
		"namespace": namespace,
		"name":      name,
	})

	err = resourceInterface.Delete(ctx, name, v1.DeleteOptions{})
	if err != nil {
		return err
	}

	if wait {
		tflog.Debug(ctx, "Waiting for resource to be deleted", map[string]any{
			"namespace": namespace,
			"name":      name,
		})
		return waitForDeletion(ctx, resourceInterface, name)
	}
	return nil
}

func waitForDeletion(ctx context.Context, r dynamic.ResourceInterface, name string) error {
	// TODO could look at usiung resourceInterface.Watch here instead of polling

	for {
		_, err := r.Get(ctx, name, v1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				break
			}
			return err
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

	return nil
}
