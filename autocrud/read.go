// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package autocrud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"k8s.io/apimachinery/pkg/api/meta"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sschema "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
)

func Read(ctx context.Context, clientGetter KubernetesClientGetter, kind, apiVersion string, req resource.ReadRequest, model any) error {
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

	tflog.Debug(ctx, "Reading resource", map[string]any{
		"name":      name,
		"namespace": namespace,
	})
	res, err := resourceInterface.Get(ctx, name, v1.GetOptions{})
	if err != nil {
		return err
	}
	tflog.Debug(ctx, "Resource read successfully", map[string]any{
		"response": res,
	})

	responseManifest := res.UnstructuredContent()

	responseMetadata := responseManifest["metadata"].(map[string]any)
	// we are expanding only for the sake of retrieving metadata
	manifest := ExpandModel(model)
	configMetadata := manifest["metadata"].(map[string]any)

	shimMetadata(responseMetadata, configMetadata, clientGetter.IgnoreLabels(), clientGetter.IgnoreAnnotations())

	FlattenManifest(responseManifest, model)
	setID(id, &model)
	return nil
}
