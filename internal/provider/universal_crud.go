package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"k8s.io/apimachinery/pkg/api/meta"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8sschema "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
)

func UniversalCreate(ctx context.Context, clientGetter *KubernetesClientGetter, kind, apiVersion string, ignoredFields []string, req resource.CreateRequest) (tftypes.Value, error) {
	client, err := clientGetter.DynamicClient()
	if err != nil {
		return tftypes.Value{}, err
	}
	discoveryClient, err := clientGetter.DiscoveryClient()
	if err != nil {
		return tftypes.Value{}, err
	}
	agr, err := restmapper.GetAPIGroupResources(discoveryClient)
	if err != nil {
		return tftypes.Value{}, err
	}
	gvk := k8sschema.FromAPIVersionAndKind(apiVersion, kind)
	restMapper := restmapper.NewDiscoveryRESTMapper(agr)
	mapping, err := restMapper.RESTMapping(gvk.GroupKind(), apiVersion)
	if err != nil {
		return tftypes.Value{}, err
	}

	flattenedManifest := UniversalFlatten(req.Plan.Raw, ignoredFields).(map[string]interface{})

	var resourceInterface dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		metadata := flattenedManifest["metadata"].(map[string]interface{})
		namespace := "default"
		if v, ok := metadata["namespace"].(string); ok && v != "" {
			namespace = v
		}
		resourceInterface = client.Resource(mapping.Resource).Namespace(namespace)
	} else {
		resourceInterface = client.Resource(mapping.Resource)
	}

	data := unstructured.Unstructured{}
	data.Object = flattenedManifest
	res, err := resourceInterface.Create(ctx,
		&data,
		v1.CreateOptions{
			// FIXME this should be configurable
			FieldManager: "terraform",
		},
	)
	if err != nil {
		return tftypes.Value{}, err
	}

	responseManifest := res.Object
	responseManifest["id"] = createID(responseManifest)
	state := UniversalExpand(responseManifest, ignoredFields)
	return state, nil
}

func UniversalRead(ctx context.Context, clientGetter *KubernetesClientGetter, kind, apiVersion string, ignoredFields []string, req resource.ReadRequest) (tftypes.Value, error) {
	client, err := clientGetter.DynamicClient()
	if err != nil {
		return tftypes.Value{}, err
	}
	discoveryClient, err := clientGetter.DiscoveryClient()
	if err != nil {
		return tftypes.Value{}, err
	}
	agr, err := restmapper.GetAPIGroupResources(discoveryClient)
	if err != nil {
		return tftypes.Value{}, err
	}
	gvk := k8sschema.FromAPIVersionAndKind(apiVersion, kind)
	restMapper := restmapper.NewDiscoveryRESTMapper(agr)
	mapping, err := restMapper.RESTMapping(gvk.GroupKind(), apiVersion)
	if err != nil {
		return tftypes.Value{}, err
	}

	var id string
	req.State.GetAttribute(ctx, path.Root("id"), &id)
	name, namespace := parseID(id)

	var resourceInterface dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		if namespace != "" {
			namespace = "default"
		}
		resourceInterface = client.Resource(mapping.Resource).Namespace(namespace)
	} else {
		resourceInterface = client.Resource(mapping.Resource)
	}

	res, err := resourceInterface.Get(ctx, name, v1.GetOptions{})
	if err != nil {
		return tftypes.Value{}, err
	}

	responseManifest := res.Object
	responseManifest["id"] = id
	state := UniversalExpand(responseManifest, ignoredFields)
	return state, nil
}

func UniversalUpdate(ctx context.Context, clientGetter *KubernetesClientGetter, kind, apiVersion string, ignoredFields []string, req resource.UpdateRequest) (tftypes.Value, error) {
	client, err := clientGetter.DynamicClient()
	if err != nil {
		return tftypes.Value{}, err
	}
	discoveryClient, err := clientGetter.DiscoveryClient()
	if err != nil {
		return tftypes.Value{}, err
	}
	agr, err := restmapper.GetAPIGroupResources(discoveryClient)
	if err != nil {
		return tftypes.Value{}, err
	}
	gvk := k8sschema.FromAPIVersionAndKind(apiVersion, kind)
	restMapper := restmapper.NewDiscoveryRESTMapper(agr)
	mapping, err := restMapper.RESTMapping(gvk.GroupKind(), apiVersion)
	if err != nil {
		return tftypes.Value{}, err
	}

	flattenedManifest := UniversalFlatten(req.Plan.Raw, ignoredFields).(map[string]interface{})

	var resourceInterface dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		metadata := flattenedManifest["metadata"].(map[string]interface{})
		namespace := "default"
		if v, ok := metadata["namespace"].(string); ok && v != "" {
			namespace = v
		}
		resourceInterface = client.Resource(mapping.Resource).Namespace(namespace)
	} else {
		resourceInterface = client.Resource(mapping.Resource)
	}

	data := unstructured.Unstructured{}
	data.Object = flattenedManifest
	res, err := resourceInterface.Update(ctx,
		&data,
		v1.UpdateOptions{
			// FIXME this should be configurable
			FieldManager: "terraform",
		},
	)
	if err != nil {
		return tftypes.Value{}, err
	}

	responseManifest := res.Object
	responseManifest["id"] = createID(responseManifest)
	state := UniversalExpand(responseManifest, ignoredFields)
	return state, nil
}

func UniversalDelete(ctx context.Context, clientGetter *KubernetesClientGetter, kind, apiVersion string, req resource.DeleteRequest) error {
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
		if namespace != "" {
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

	return nil
}

func createID(manifest map[string]interface{}) string {
	if metadata, ok := manifest["metadata"].(map[string]interface{}); ok {
		id := metadata["name"].(string)
		if ns, ok := metadata["namespace"].(string); ok && ns != "" {
			id = id + "/" + ns
		}
		return id
	}
	return ""
}

func parseID(id string) (string, string) {
	parts := strings.Split(id, "/")
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], parts[1]
}
