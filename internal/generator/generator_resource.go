package generator

import (
	"log/slog"
	"path"
	"time"

	specresource "github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	specschema "github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

type ResourceGenerator struct {
	GeneratedTimestamp time.Time
	ResourceConfig     ResourceConfig
	Schema             SchemaGenerator
	ModelFields        ModelFieldsGenerator
}

func NewResourceGenerator(cfg ResourceConfig, spec specresource.Resource) ResourceGenerator {
	attributes := AttributesGenerator{{
		Name:          "id",
		AttributeType: StringAttributeType,
		Computed:      true,
		Description:   "The unique ID for this terraform resource",
	}}

	modelFields := ModelFieldsGenerator{{
		FieldName:     "ID",
		Type:          StringModelType,
		AttributeType: StringAttributeType,
		AttributeName: "id",
	}}

	return ResourceGenerator{
		GeneratedTimestamp: time.Now(),
		ResourceConfig:     cfg,
		ModelFields:        append(modelFields, GenerateModelFields(spec.Schema.Attributes, cfg.IgnoredAttributes, "")...),
		Schema: SchemaGenerator{
			Name:        cfg.Name,
			Description: cfg.Description,
			Attributes:  append(attributes, GenerateAttributes(spec.Schema.Attributes, cfg.IgnoredAttributes, cfg.ComputedAttributes, cfg.RequiredAttributes, cfg.SensitiveAttributes, cfg.ImmutableAttributes, cfg.DefaultValueAttributes, "")...),
		},
	}
}

func (g *ResourceGenerator) GenerateSchemaFunctionCode() string {
	imports := getPlanModifierImports(g.Schema.Attributes)

	if len(imports) > 0 {
		imports = append(imports, path.Join(schemaImportPath, "planmodifier"))
		g.Schema.Imports = imports
	}

	return renderTemplate(schemaFunctionTemplate, g)
}

func (g *ResourceGenerator) GenerateCRUDStubCode() string {
	return renderTemplate(crudStubsTemplate, g)
}

func (g *ResourceGenerator) GenerateResourceCode() string {
	return renderTemplate(resourceTemplate, g)
}

func (g *ResourceGenerator) GenerateModelCode() string {
	return renderTemplate(modelTemplate, g)
}

func (g *ResourceGenerator) GenerateAutoCRUDCode() string {
	return renderTemplate(autocrudTemplate, g)
}

func (g *ResourceGenerator) GenerateAutoCRUDHooksCode() string {
	return renderTemplate(autocrudHooksTemplate, g)
}

func (g *ResourceGenerator) GenerateDefaultValuesCode() string {
	return renderTemplate(defaultValuesTemplate, g)
}

// TODO create a walkAttributes function that abstracts the logic of traversing
// the spec for attributes

// FIXME this function has too many parameters now, should maybe be part of ResourceGenerator.
func GenerateAttributes(attrs specresource.Attributes, ignored, computed, required, sensitive, immutable, default_values []string, path string) AttributesGenerator {
	generatedAttrs := AttributesGenerator{}
	for _, attr := range attrs {
		attributePath := path + attr.Name

		if stringInSlice(attributePath, ignored) {
			continue
		}

		generatedAttr := AttributeGenerator{
			Name:         attr.Name,
			Required:     stringInSlice(attributePath, required),
			Computed:     stringInSlice(attributePath, computed),
			Sensitive:    stringInSlice(attributePath, sensitive),
      Immutable: stringInSlice(attributePath, immutable),
			DefaultValue: stringInSlice(attributePath, default_values),
		}
		switch {
		case attr.Bool != nil:
			if attr.Bool.Description != nil {
				generatedAttr.Description = *attr.Bool.Description
			}
			generatedAttr.AttributeType = BoolAttributeType
			generatedAttr.PlanModifierType = BoolPlanModifierType
			generatedAttr.PlanModifierPackage = BoolPlanModifierPackage
		case attr.String != nil:
			if attr.String.Description != nil {
				generatedAttr.Description = *attr.String.Description
			}
			generatedAttr.AttributeType = StringAttributeType
			generatedAttr.PlanModifierType = StringPlanModifierType
			generatedAttr.PlanModifierPackage = StringPlanModifierPackage
		case attr.Number != nil:
			if attr.Number.Description != nil {
				generatedAttr.Description = *attr.Number.Description
			}
			generatedAttr.AttributeType = NumberAttributeType
			generatedAttr.PlanModifierType = NumberPlanModifierType
			generatedAttr.PlanModifierPackage = NumberPlanModifierPackage
		case attr.Int64 != nil:
			if attr.Int64.Description != nil {
				generatedAttr.Description = *attr.Int64.Description
			}
			generatedAttr.AttributeType = Int64AttributeType
			generatedAttr.PlanModifierType = Int64PlanModifierType
			generatedAttr.PlanModifierPackage = Int64PlanModifierPackage
		case attr.Map != nil:
			if attr.Map.Description != nil {
				generatedAttr.Description = *attr.Map.Description
			}
			generatedAttr.AttributeType = MapAttributeType
			generatedAttr.ElementType = getElementType(attr.Map.ElementType)
		case attr.List != nil:
			if attr.List.Description != nil {
				generatedAttr.Description = *attr.List.Description
			}
			generatedAttr.AttributeType = ListAttributeType
			generatedAttr.ElementType = getElementType(attr.List.ElementType)
		case attr.SingleNested != nil:
			if attr.SingleNested.Description != nil {
				generatedAttr.Description = *attr.SingleNested.Description
			}
			generatedAttr.AttributeType = SingleNestedAttributeType
			generatedAttr.NestedAttributes = GenerateAttributes(attr.SingleNested.Attributes, ignored, computed, required, sensitive, immutable, default_values, attributePath+".")
		case attr.ListNested != nil:
			if attr.ListNested.Description != nil {
				generatedAttr.Description = *attr.ListNested.Description
			}
			generatedAttr.AttributeType = ListNestedAttributeType
			generatedAttr.NestedAttributes = GenerateAttributes(attr.ListNested.NestedObject.Attributes, ignored, computed, required, sensitive, immutable, default_values, attributePath+"[*].")
		}
		generatedAttrs = append(generatedAttrs, generatedAttr)
	}
	return generatedAttrs
}

func stringInSlice(str string, slice []string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func GenerateModelFields(attrs specresource.Attributes, ignored []string, path string) ModelFieldsGenerator {
	generatedModelFields := ModelFieldsGenerator{}
	for _, attr := range attrs {
		attributePath := path + attr.Name

		if stringInSlice(attributePath, ignored) {
			continue
		}

		generatedModelField := ModelFieldGenerator{
			FieldName:         MapTerraformAttributeToModel(attr.Name),
			ManifestFieldName: MapTerraformAttributeToKubernetes(attr.Name),
			AttributeName:     attr.Name,
		}
		switch {
		case attr.Bool != nil:
			generatedModelField.AttributeType = BoolAttributeType
			generatedModelField.Type = BoolModelType
		case attr.String != nil:
			generatedModelField.AttributeType = StringAttributeType
			generatedModelField.Type = StringModelType
		case attr.Number != nil:
			generatedModelField.AttributeType = NumberAttributeType
			generatedModelField.Type = NumberModelType
		case attr.Int64 != nil:
			generatedModelField.AttributeType = Int64AttributeType
			generatedModelField.Type = Int64ModelType
		case attr.Map != nil:
			generatedModelField.AttributeType = MapAttributeType
			generatedModelField.ElementType = getModelElementType(attr.Map.ElementType)
		case attr.List != nil:
			generatedModelField.AttributeType = ListAttributeType
			generatedModelField.ElementType = getModelElementType(attr.List.ElementType)
		case attr.SingleNested != nil:
			generatedModelField.AttributeType = SingleNestedAttributeType
			generatedModelField.NestedFields = GenerateModelFields(attr.SingleNested.Attributes, ignored, attributePath+".")
			if len(generatedModelField.NestedFields) == 0 {
				slog.Warn("Ignoring nested attribute with no schema", "name", attr.Name)
				continue
			}
		case attr.ListNested != nil:
			generatedModelField.AttributeType = ListNestedAttributeType
			generatedModelField.NestedFields = GenerateModelFields(attr.ListNested.NestedObject.Attributes, ignored, attributePath+"[*].")
			if len(generatedModelField.NestedFields) == 0 {
				slog.Warn("Ignoring nested attribute with no schema", "name", attr.Name)
				continue
			}
		}
		generatedModelFields = append(generatedModelFields, generatedModelField)
	}
	return generatedModelFields
}

func isComputed(c specschema.ComputedOptionalRequired) bool {
	return c == specschema.Computed || c == specschema.ComputedOptional
}

func isRequired(c specschema.ComputedOptionalRequired) bool {
	return c == specschema.Required
}

func isSensitive(s *bool) bool {
	return s != nil && *s
}

func getElementType(e specschema.ElementType) string {
	switch {
	case e.Bool != nil:
		return BoolElementType
	case e.String != nil:
		return StringElementType
	case e.Number != nil:
		return NumberElementType
	case e.Int64 != nil:
		return Int64ElementType
	}
	panic("unsupported element type")
}

func getModelElementType(e specschema.ElementType) string {
	switch {
	case e.Bool != nil:
		return BoolModelType
	case e.String != nil:
		return StringModelType
	case e.Number != nil:
		return NumberModelType
	case e.Int64 != nil:
		return Int64ModelType
	}
	panic("unsupported element type")
}
