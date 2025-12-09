// Copyright IBM Corp. 2024
// SPDX-License-Identifier: MPL-2.0

// metadata.go contains code that ports and shims backwards compatibility for
// the metadata of api responses
package autocrud

import (
	"net/url"
	"regexp"
	"strings"
)

func shimMetadata(responseMetadata map[string]any, configMetadata map[string]any, ignoreLabels, ignoreAnnotations []string) {
	ignoreKeys("labels", ignoreLabels, responseMetadata, configMetadata)
	ignoreKeys("annotations", ignoreAnnotations, responseMetadata, configMetadata)
	// the previous SDKv2 implementation assumed the zero-value even when
	// generation is not present in all API responses, this may have been in
	// error but this change maintains compatibility.
	if _, ok := responseMetadata["generation"]; !ok {
		responseMetadata["generation"] = int64(0)
	}
}

func ignoreKeys(typ string, ignore []string, responseMetadata, configMetadata map[string]any) {
	if _, ok := responseMetadata[typ]; ok {
		keys := responseMetadata[typ].(map[string]any)
		// remove internal use labels/annotations not set in config
		removeInternalKeys(keys, configMetadata[typ].(map[string]any))
		// remove regex matching labels/annotations specified by the user in the provider block
		removeKeys(keys, configMetadata[typ].(map[string]any), ignore)
		// if the remaining map is empty, set it to nil so the plan matches the config
		// fortunately this does not break the scenario where a user specifies an empty
		// map explicitly
		if len(keys) == 0 {
			delete(responseMetadata, typ)
		}
	}
}

func removeInternalKeys(m map[string]any, d map[string]any) {
	for k := range m {
		if isInternalKey(k) && !isKeyInMap(k, d) {
			delete(m, k)
		}
	}
}

func isKeyInMap(key string, d map[string]any) bool {
	_, ok := d[key]
	return ok
}

func isInternalKey(annotationKey string) bool {
	u, err := url.Parse("//" + annotationKey)
	if err != nil {
		return false
	}

	// allow user specified application specific keys
	if u.Hostname() == "app.kubernetes.io" {
		return false
	}

	// allow AWS load balancer configuration annotations
	if u.Hostname() == "service.beta.kubernetes.io" {
		return false
	}

	// internal *.kubernetes.io keys
	if strings.HasSuffix(u.Hostname(), "kubernetes.io") {
		return true
	}

	// Specific to DaemonSet annotations, generated & controlled by the server.
	if strings.Contains(annotationKey, "deprecated.daemonset.template.generation") {
		return true
	}
	return false
}

// removeKeys removes given Kubernetes metadata(annotations and labels) keys.
// In that case, they won't be available in the TF state file and will be ignored during apply/plan operations.
func removeKeys(m map[string]any, d map[string]any, ignoreKubernetesMetadataKeys []string) {
	for k := range m {
		if ignoreKey(k, ignoreKubernetesMetadataKeys) && !isKeyInMap(k, d) {
			delete(m, k)
		}
	}
}

// ignoreKey reports whether the Kubernetes metadata(annotations and labels) key contains
// any match of the regular expression pattern from the expressions slice.
func ignoreKey(key string, expressions []string) bool {
	for _, e := range expressions {
		if ok, _ := regexp.MatchString(e, key); ok {
			return true
		}
	}

	return false
}
