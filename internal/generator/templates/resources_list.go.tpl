// Code generated by hashicorp/terraform-plugin-codegen-kubernetes; DO NOT EDIT 
//
// This file contains the list of constructors for resources that have been autogenerated.
//
// This code was written by a robot on {{ .GeneratedTimestamp.Format "Jan 02, 2006 15:04:05 UTC" }}.

package provider

import (
    "github.com/hashicorp/terraform-plugin-framework/resource"

{{ range $val := .Packages }}
    "github.com/hashicorp/terraform-provider-kubernetes/internal/framework/provider/{{ $val }}"
{{- end }}
)

var generatedResources = []func() resource.Resource{
{{- range $val := .Resources }}
   {{ $val.Package }}.New{{ $val.Kind }},
{{- end }}
}
