package generator

type AttributeGenerator struct {
	Name string

	AttributeType string
	ElementType   string

	PlanModifierType    string
	PlanModifierPackage string

	Required    bool
	Description string
	Computed    bool
	Sensitive   bool
	Immutable   bool

	NestedAttributes AttributesGenerator
}

func (g AttributeGenerator) String() string {
	return renderTemplate(attributeTemplate, g)
}

type AttributesGenerator []AttributeGenerator

func (g AttributesGenerator) String() string {
	return renderTemplate(attributesTemplate, g)
}

const sampleValidatorFile = `package validators

import (
    "context"
    "fmt"
    "unicode"

    "github.com/hashicorp/terraform-plugin-framework/diag"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/validator"
    "github.com/hashicorp/terraform-plugin-framework/validator/stringvalidator"
)

type nameStartsUppercaseValidator struct{}

func (v nameStartsUppercaseValidator) Description(ctx context.Context) string {
    return "Ensures the string starts with an uppercase letter."
}

func (v nameStartsUppercaseValidator) MarkdownDescription(ctx context.Context) string {
    return "Ensures the string starts with an uppercase letter."
}

func (v nameStartsUppercaseValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
    if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
        return
    }

    value := req.ConfigValue.ValueString()
    if len(value) == 0 || !unicode.IsUpper(rune(value[0])) {
        resp.Diagnostics.AddError(
            "Invalid Name",
            fmt.Sprintf("The value for %q must start with an uppercase letter.", req.Path),
        )
    }
}

func NameStartsUppercase() validator.String {
    return nameStartsUppercaseValidator{}
}`
