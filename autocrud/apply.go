// Copyright IBM Corp. 2024
// SPDX-License-Identifier: MPL-2.0

package autocrud

import (
	"context"
	"reflect"

	"k8s.io/apimachinery/pkg/api/meta"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8sschema "k8s.io/apimachinery/pkg/runtime/schema"
	patchtypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func serverSideApply(ctx context.Context, clientGetter KubernetesClientGetter, apiVersion, kind string, model any) error {
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

	manifest := ExpandModel(model)

	var resourceInterface dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		metadata := manifest["metadata"].(map[string]interface{})
		namespace := "default"
		if v, ok := metadata["namespace"].(string); ok && v != "" {
			namespace = v
		}
		resourceInterface = client.Resource(mapping.Resource).Namespace(namespace)
	} else {
		resourceInterface = client.Resource(mapping.Resource)
	}

	obj := unstructured.Unstructured{}
	obj.SetUnstructuredContent(manifest)
	obj.SetKind(kind)
	obj.SetAPIVersion(apiVersion)
	payload, err := obj.MarshalJSON()
	if err != nil {
		return err
	}

	tflog.Debug(ctx, "Executing server-side apply operation", map[string]any{
		"manifest": obj,
	})
	res, err := resourceInterface.Patch(ctx, obj.GetName(), patchtypes.ApplyPatchType, payload,
		v1.PatchOptions{
			// FIXME this should be configurable
			FieldManager: "terraform",
		},
	)
	if err != nil {
		return err
	}

	tflog.Debug(ctx, "Server-side apply operation succeeded", map[string]any{
		"response": res,
	})

	responseManifest := res.UnstructuredContent()
	id := createID(responseManifest)

	responseMetadata := responseManifest["metadata"].(map[string]any)
	configMetadata := manifest["metadata"].(map[string]any)

	shimMetadata(responseMetadata, configMetadata, clientGetter.IgnoreLabels(), clientGetter.IgnoreAnnotations())

	err = FlattenManifest(responseManifest, model)
	if err != nil {
		return err
	}
	setID(id, model)

	return nil
}

func setID(ID string, model any) {
	// FIXME: we shouldnt need reflection here. We should make some sort
	//        of Model interface with SetID(), Expand(), Flatten()
	idval := reflect.ValueOf(model).Elem().FieldByName("ID")
	idval.Set(reflect.ValueOf(types.StringValue(ID)))
}
