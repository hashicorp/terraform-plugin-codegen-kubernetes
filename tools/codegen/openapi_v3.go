package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/getkin/kin-openapi/openapi3"

	_ "embed"
)

// TODO add const for fields that should always be stripped, e.g. as fieldManager
// TODO replace camel case in descriptions with terraform snake case
// TODO singularize blocks that are arrays, e.g containers -> container
// TODO use enum field to add validators e.g ServiceSpec.type field
// TODO autogenerate schema for GetResources and GetDataSources

var (
	ref        = flag.String("ref", "", "reference to generate the Terraform schema for")
	pkg        = flag.String("pkg", "", "name of the Go package for this resource")
	outputFile = flag.String("o", "", "file to write the generated code to")
	jsonFile   = flag.String("json", "", "file to read the OpenAPI v3 schema information from")
	gofmt      = flag.Bool("fmt", true, "run the generated file through go fmt")
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

func getBlock(name string, schema *openapi3.Schema, requiredBlock bool) TerraformBlock {
	attributes := []TerraformAttribute{}
	blocks := []TerraformBlock{}

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
		log.Printf("WARN: schema for %q produced an empty block\n", name)
	}

	for name, prop := range properties {
		if name == "managedFields" {
			// skip managedFields for now
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
			blocks = append(blocks, getBlock(snakify(name), prop.Value, required))
			continue
		}
		attributes = append(attributes, TerraformAttribute{
			Name:          snakify(name),
			Description:   stripBackticks(prop.Value.Description),
			AttributeType: attributeType,
			ElementType:   elementType,
			Required:      required,
		})
	}

	block := TerraformBlock{
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

func main() {
	flag.Parse()

	doc, err := openapi3.NewLoader().LoadFromFile(*jsonFile)
	if err != nil {
		log.Fatalf("error loading OpenAPI specification: %v", err)
	}

	schema, ok := doc.Components.Schemas[*ref]
	if !ok {
		log.Fatalf("no schema for %q exists in OpenAPI document %q", *ref, *jsonFile)
	}

	parts := strings.Split(*ref, ".")
	kind := parts[len(parts)-1]
	resourceName := kind
	resource := TerraformResource{
		Package:               *pkg,
		ResourceName:          resourceName,
		TerraformResourceName: snakify(resourceName),
		ResourceBlock:         getBlock("root", schema.Value, false),
	}

	log.Printf("Generating Go code for terraform resource from OpenAPI ref %s", *ref)

	f, err := os.OpenFile(*outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("error opening file (%s): %s", *outputFile, err)
	}

	src := resource.Bytes()
	if *gofmt {
		src, err = format.Source(resource.Bytes())
		if err != nil {
			log.Fatalf("error formatting generated Go code: %s", err)
		}
	}

	if _, err := f.Write(src); err != nil {
		log.Fatalf("error writing to file %q: %s", *outputFile, err)
	}

	f.Close()
}

type TerraformAttribute struct {
	Name          string
	AttributeType string
	ElementType   string
	Required      bool
	Description   string
}

//go:embed templates/attribute.go.tpl
var attributeTemplate string

//go:embed templates/block.go.tpl
var blockTemplate string

//go:embed templates/resource.go.tpl
var resourceTemplate string

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

type TerraformBlock struct {
	Name        string
	Description string
	Attributes  []TerraformAttribute
	Blocks      []TerraformBlock
}

func (b TerraformBlock) String() string {
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
	// The Go package name the resource lives in
	Package string

	// The resource name in CamelCase format
	ResourceName string

	// The Terraform resource name in snake_case
	TerraformResourceName string

	ResourceBlock TerraformBlock
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
