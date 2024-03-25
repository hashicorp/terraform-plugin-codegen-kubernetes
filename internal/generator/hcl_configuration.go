package generator

import "github.com/hashicorp/hcl/v2/hclsimple"

// GeneratorConfig is the top level code generator configuration
type GeneratorConfig struct {
	Resources  []ResourceConfig   `hcl:"resource,block"`
	DataSource []DataSourceConfig `hcl:"data,block"`
}

// ResourceConfig configures code generation for a Terraform resource
type ResourceConfig struct {
	// Name is the terraform name of this resource
	Name string `hcl:"name,label"`

	// Package is the name of the Go package for the source files for this resource
	Package string `hcl:"package"`

	// OutputFilenamePrefix is a prefix to be added to all source files generated
	// for this resource
	OutputFilenamePrefix string `hcl:"output_filename_prefix"`

	// APIVersion is the Kubernetes API version of the resource
	APIVersion string `hcl:"api_version"`

	// Kind is the Kubernetes kind of the resource
	Kind string `hcl:"kind"`

	// Description is a Markdown description for the resource
	Description string `hcl:"description"`

	// IgnoredAttributes is a list of attribute paths to omit from the resource
	IgnoredAttributes []string `hcl:"ignored_attributes,optional"`

	// FIXME not implemented yet
	// RequiredAttributes is a list of attribute paths to mark as required in the schema
	RequiredAttributes []string `hcl:"required_attributes,optional"`

	// FIXME not implemented yet
	// ComputedAttributes is a list of attribute paths to mark as computed in the schema
	ComputedAttributes []string `hcl:"computed_attributes,optional"`

	// FIXME not implemented yet
	// SensitiveAttributes is a list of attribute paths to mark as sensitive in the schema
	SensitiveAttributes []string `hcl:"sensitive_attributes,optional"`

	// FIXME not implemented yet
	// ImmutableAttributes is a list of attribute paths to mark as requiring a forced
	// replacement if changed in the schema
	ImmutableAttributes []string `hcl:"immutable_attributes,optional"`

	// Generate controls generator specific options
	Generate GenerateConfig `hcl:"generate,block"`

	// OpenAPIConfig configures options for the OpenAPI to Framework IR generator
	OpenAPIConfig TerraformPluginGenOpenAPIConfig `hcl:"openapi,block"`

	// Disabled tells the generator to skip this configuration
	Disabled bool `hcl:"disabled,optional"`
}

// DataSourceConfig configures code generation for a Terraform data source
type DataSourceConfig struct {
	// TODO implement data source generation
}

// TerraformPluginGenOpenAPIConfig supplies configuration to tfplugingen-openapi
// See: https://github.com/hashicorp/terraform-plugin-codegen-openapi
type TerraformPluginGenOpenAPIConfig struct {
	// Filename is the filename for the OpenAPI JSON specification
	Filename string `hcl:"filename"`

	// CreatePath is the POST path for the resource in the OpenAPI spec, e.g. /api/v1/namespaces/{namespace}/configmaps
	CreatePath string `hcl:"create_path"`

	// ReadPath is the GET path for the resource in the OpenAPI spec, e.g. /api/v1/namespaces/{namespace}/configmaps/{name}
	ReadPath string `hcl:"read_path"`
}

// CRUDAutoOptions configures options for the autocrud template
type CRUDAutoOptions struct {
	WaitForDeletion bool   `hcl:"wait_for_deletion,optional"`
	Hooks           *Hooks `hcl:"hooks,block"`
}

// Hooks configures which hooks to include for autocrud template if necessary
type Hooks struct {
	BeforeCreate bool `hcl:"before_create,optional"`
	AfterCreate  bool `hcl:"after_create,optional"`
	BeforeRead   bool `hcl:"before_read,optional"`
	AfterRead    bool `hcl:"after_read,optional"`
	BeforeUpdate bool `hcl:"before_update,optional"`
	AfterUpdate  bool `hcl:"after_update,optional"`
	BeforeDelete bool `hcl:"before_delete,optional"`
	AfterDelete  bool `hcl:"after_delete,optional"`
}

// GenerateConfig configures the options for what we should generate
type GenerateConfig struct {
	Schema          bool             `hcl:"schema,optional"`
	Overrides       bool             `hcl:"overrides,optional"`
	Model           bool             `hcl:"model,optional"`
	CRUDAuto        bool             `hcl:"autocrud,optional"`
	CRUDAutoOptions *CRUDAutoOptions `hcl:"autocrud_options,block"`
	CRUDStubs       bool             `hcl:"crud_stubs,optional"`
}

// ParseHCLConfig parses the .hcl configuraiton file and
// produces a GeneratorConfig
func ParseHCLConfig(filename string) (GeneratorConfig, error) {
	config := GeneratorConfig{}
	err := hclsimple.DecodeFile(filename, nil, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

// Checks whether hooks are used to prevent file from being generated if block is empty or all set to false.
func (h *Hooks) IsEmpty() bool {
	return !(h.AfterCreate || h.BeforeCreate || h.AfterRead || h.BeforeRead || h.AfterUpdate || h.BeforeUpdate || h.AfterDelete || h.BeforeDelete)
}
