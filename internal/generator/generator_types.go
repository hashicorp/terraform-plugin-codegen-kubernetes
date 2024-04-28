package generator

const (
	BoolAttributeType         = "BoolAttribute"
	StringAttributeType       = "StringAttribute"
	NumberAttributeType       = "NumberAttribute"
	Int64AttributeType        = "Int64Attribute"
	MapAttributeType          = "MapAttribute"
	ListAttributeType         = "ListAttribute"
	ObjectAttributeType       = "ObjectAttribute"
	SingleNestedAttributeType = "SingleNestedAttribute"
	ListNestedAttributeType   = "ListNestedAttribute"
)

// TODO: we need to expland these types to include float64, list, map, object

const (
	BoolElementType   = "BoolType"
	StringElementType = "StringType"
	NumberElementType = "NumberType"
	Int64ElementType  = "Int64Type"
)

const (
	BoolModelType   = "Bool"
	StringModelType = "String"
	NumberModelType = "Number"
	Int64ModelType  = "Int64"
)

const (
	BoolPlanModifierType   = "Bool"
	StringPlanModifierType = "String"
	NumberPlanModifierType = "Number"
	Int64PlanModifierType  = "Int64"
)

const (
	BoolPlanModifierPackage   = "boolplanmodifier"
	StringPlanModifierPackage = "stringplanmodifier"
	NumberPlanModifierPackage = "numberplanmodifier"
	Int64PlanModifierPackage  = "int64planmodifier"
)
