package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/getkin/kin-openapi/openapi3"

	_ "embed"
)

// TODO add const for fields that should always be stripped, e.g. as fieldManager
// TODO replace camel case in descriptions with terraform snake case
// TODO singularize blocks that are arrays, e.g containers -> container
// TODO use enum field to add validators e.g ServiceSpec.type field
// TODO autogenerate schema for GetResources and GetDataSources

var (
	gofmt = flag.Bool("fmt", true, "run the generated files through go fmt")

	cfgfile = flag.String("cfg", "gen.json", "path to the JSON configuration file for the code generator")
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func snakify(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func stripBackticks(str string) string {
	return strings.ReplaceAll(str, "`", "")
}

func skipField(f string) bool {
	for _, v := range []string{"managedFields", "ownerReferences"} {
		if f == v {
			return true
		}
	}
	return false
}

func stringInSlice(s string, ss []string) bool {
	for _, sss := range ss {
		if sss == s {
			return true
		}
	}
	return false
}

func getSchema(name string, schema *openapi3.Schema, requiredBlock bool, ignoredFields, computedFields []string) TerraformSchema {
	attributes := []TerraformAttribute{}
	blocks := []TerraformSchema{}

	// NOTE some schemas are just a reference to another schema
	if len(schema.AllOf) > 0 {
		schema = schema.AllOf[0].Value
	}

	properties := schema.Properties
	if schema.Type == "array" {
		properties = schema.Items.Value.Properties
		if len(properties) == 0 && len(schema.Items.Value.AllOf) > 0 {
			properties = schema.Items.Value.AllOf[0].Value.Properties
		}
	}
	if len(properties) == 0 {
		log.Printf("warning: schema for %q produced an empty block\n", name)
	}

	for name, prop := range properties {
		// FIXME make this configurable
		if stringInSlice(name, ignoredFields) {
			continue
		}

		required := false
		for _, r := range schema.Required {
			if r == name {
				required = true
			}
		}
		attributeType, elementType, err := getAttributeType(prop.Value)
		if err != nil {
			// this means the attribute needs to be a block
			blocks = append(blocks, getSchema(snakify(name), prop.Value, required, ignoredFields, computedFields))
			continue
		}
		attributes = append(attributes, TerraformAttribute{
			Name:          snakify(name),
			Description:   stripBackticks(prop.Value.Description),
			AttributeType: attributeType,
			ElementType:   elementType,
			Required:      required,
			Computed:      stringInSlice(name, computedFields),
		})
	}

	block := TerraformSchema{
		Name:        name,
		Description: schema.Description,
		Attributes:  attributes,
		Blocks:      blocks,
	}

	return block
}

func getElementType(value *openapi3.Schema) (string, error) {
	switch value.Type {
	case "boolean":
		return "BoolType", nil
	case "string":
		return "StringType", nil
	case "number":
		return "NumberType", nil
	case "integer":
		return "Int64Type", nil
	}
	return "", fmt.Errorf("cannot use complex type as element type")
}

// getAttributeType will return the string representation of the framework type
// needed for the supplied OpenAPI schema. If the schema cannot be represented
// as a framework type the function will return an error
func getAttributeType(schema *openapi3.Schema) (string, string, error) {
	switch schema.Type {
	case "boolean":
		return "BoolAttribute", "", nil
	case "string":
		return "StringAttribute", "", nil
	case "number":
		return "NumberAttribute", "", nil
	case "integer":
		return "Int64Attribute", "", nil
	case "array":
		elementType, err := getElementType(schema.Items.Value)
		if err != nil {
			return "", "", err
		}
		return "ListAttribute", elementType, nil
	case "object":
		// NOTE objects where the schema's only additional property is any string
		// should be a map of strings
		if addProps := schema.AdditionalProperties; addProps != nil {
			if addProps.Value.Type == "string" {
				elementType, err := getElementType(schema.AdditionalProperties.Value)
				if err != nil {
					return "", "", err
				}
				return "MapAttribute", elementType, nil
			}
		}
	}

	// NOTE fields that can be one of int or string should use only string
	if len(schema.OneOf) == 2 {
		return "StringAttribute", "", nil
	}

	// NOTE some fields have a schema that is just a reference to another schema
	if len(schema.AllOf) > 0 {
		return getAttributeType(schema.AllOf[0].Value)
	}

	return "", "", fmt.Errorf("complex schema cannot be an attribute and must be a block")
}

type OpenAPIv3Config struct {
	Source string `json:"source"`
	Ref    string `json:"ref"`
}

type ResourceConfig struct {
	Package         string          `json:"package"`
	ResourceName    string          `json:"resource_name"`
	Kind            string          `json:"kind"`
	APIVersion      string          `json:"apiVersion"`
	Filename        string          `json:"filename"`
	OpenAPIv3Config OpenAPIv3Config `json:"openapi_v3"`
	IgnoreFields    []string        `json:"ignore_fields"`
	ComputedFields  []string        `json:"computed_fields"`
}

type Config struct {
	ResourcesConfig []ResourceConfig `json:"resources"`
}

func main() {
	flag.Parse()

	f, err := os.Open(*cfgfile)
	if err != nil {
		log.Fatalf("could not open JSON configuration file %q: %v", *cfgfile, err)
	}

	var config Config
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("could not decode JSON configuration file: %v", err)
	}

	resources := []string{}
	for _, r := range config.ResourcesConfig {
		doc, err := openapi3.NewLoader().LoadFromFile("../../tools/codegen/data/" + r.OpenAPIv3Config.Source)
		if err != nil {
			log.Fatalf("error loading OpenAPI specification: %v", err)
		}

		schema, ok := doc.Components.Schemas[r.OpenAPIv3Config.Ref]
		if !ok {
			log.Fatalf("no schema for %q exists in OpenAPI document %q", r.OpenAPIv3Config.Ref, r.OpenAPIv3Config.Source)
		}

		parts := strings.Split(r.OpenAPIv3Config.Ref, ".")
		kind := parts[len(parts)-1]
		resource := TerraformResource{
			Package:            r.Package,
			Kind:               kind,
			APIVersion:         r.APIVersion,
			ResourceName:       r.ResourceName,
			Root:               getSchema("root", schema.Value, false, r.IgnoreFields, r.ComputedFields),
			GeneratedTimestamp: time.Now(),
		}

		log.Printf("Generating Go code for terraform resource %q from OpenAPI ref %s",
			resource.ResourceName, r.OpenAPIv3Config.Ref)

		f, err := os.OpenFile(r.Filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatalf("error opening file (%s): %s", r.Filename, err)
		}

		src := resource.Bytes()
		if *gofmt {
			src, err = format.Source(resource.Bytes())
			if err != nil {
				log.Fatalf("error formatting generated Go code: %s", err)
			}
		}

		if _, err := f.Write(src); err != nil {
			log.Fatalf("error writing to file %q: %s", r.Filename, err)
		}

		f.Close()

		resources = append(resources, kind)
	}

	log.Printf("Generating list of resources in resources_list.go")
	resourceListFilename := "resources_list.go"
	resourcesList := ResourcesList{
		Package:   "provider",
		Resources: resources,
	}
	f, err = os.OpenFile(resourceListFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("error opening file (%s): %s", resourceListFilename, err)
	}
	src, err := format.Source(resourcesList.Bytes())
	if err != nil {
		log.Fatalf("error formatting generated Go code: %s", err)
	}
	if _, err := f.Write(src); err != nil {
		log.Fatalf("error writing to file %q: %s", resourceListFilename, err)
	}

	f.Close()
}

//go:embed templates/attribute.go.tpl
var attributeTemplate string

//go:embed templates/schema.go.tpl
var blockTemplate string

//go:embed templates/resource.go.tpl
var resourceTemplate string

//go:embed templates/resources_list.go.tpl
var resourcesListTemplate string

type TerraformAttribute struct {
	Name          string
	AttributeType string
	ElementType   string
	Required      bool
	Description   string
	Computed      bool
}

func (a TerraformAttribute) String() string {
	tpl, err := template.New("").Parse(attributeTemplate)
	if err != nil {
		panic(fmt.Sprintf("could not parse template: %v", err))
	}
	var buf bytes.Buffer
	err = tpl.Execute(&buf, a)
	if err != nil {
		panic(fmt.Sprintf("error executing template: %v", err))
	}
	return buf.String()
}

type TerraformSchema struct {
	Name        string
	Description string
	Attributes  []TerraformAttribute
	Blocks      []TerraformSchema
}

func (b TerraformSchema) String() string {
	tpl, err := template.New("").Parse(blockTemplate)
	if err != nil {
		panic(fmt.Sprintf("could not parse template: %v", err))
	}
	var buf bytes.Buffer
	err = tpl.Execute(&buf, b)
	if err != nil {
		panic(fmt.Sprintf("error executing template: %v", err))
	}
	return buf.String()
}

type TerraformResource struct {
	// GeneratedTimestamp
	GeneratedTimestamp time.Time

	// Package is the Go package name the resource lives in
	Package string

	// Kind is the Kubernetes resource kind
	Kind string

	// APIVersion is the Kubernetes resource apiVersion
	APIVersion string

	// ResourceName is the Terraform resource name in snake_case
	ResourceName string

	Root TerraformSchema
}

func (r TerraformResource) String() string {
	tpl, err := template.New("").Parse(resourceTemplate)
	if err != nil {
		panic(fmt.Sprintf("could not parse template: %v", err))
	}
	var buf bytes.Buffer
	err = tpl.Execute(&buf, r)
	if err != nil {
		panic(fmt.Sprintf("error executing template: %v", err))
	}
	return buf.String()
}

func (r TerraformResource) Bytes() []byte {
	return []byte(r.String())
}

type ResourcesList struct {
	// The Go package name the resource lives in
	Package string

	// The list of resource names
	Resources []string
}

func (r ResourcesList) String() string {
	tpl, err := template.New("").Parse(resourcesListTemplate)
	if err != nil {
		panic(fmt.Sprintf("could not parse template: %v", err))
	}
	var buf bytes.Buffer
	err = tpl.Execute(&buf, r)
	if err != nil {
		panic(fmt.Sprintf("error executing template: %v", err))
	}
	return buf.String()
}

func (r ResourcesList) Bytes() []byte {
	return []byte(r.String())
}
