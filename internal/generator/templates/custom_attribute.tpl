{{- if .WaitForRollout }}
			"wait_for_rollout": schema.BoolAttribute{
				Optional:    true,
				Description: "Terraform will wait for the default service account to be created.",
			},
{{- end }}

{{- if .WaitForDefaultServiceAccount }}
			"wait_for_default_service_account": schema.BoolAttribute{
				Optional:    true,
				Description: "Terraform will wait for the default service account to be created.",
			},
{{- end }}

{{- if .WaitForLoadBalancer }}
			"wait_for_load_balancer": schema.BoolAttribute{
				Optional:    true,
				Description: "Terraform will wait for the default service account to be created.",
			},
{{- end }}

{{- if .WaitForCompletion }}
			"wait_for_completion": schema.BoolAttribute{
				Optional:    true,
				Description: "Terraform will wait for the default service account to be created.",
			},
{{- end }}