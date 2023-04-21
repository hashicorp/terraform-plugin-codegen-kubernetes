package provider

import (
	"strings"
	"unicode"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func UniversalExpand(in interface{}, ignoredFields []string) tftypes.Value {
	return expand(in, ignoredFields)
}

func expand(in interface{}, ignoredFields []string) tftypes.Value {
	switch v := in.(type) {
	case string:
		return tftypes.NewValue(tftypes.String, v)
	case bool:
		return tftypes.NewValue(tftypes.Bool, v)
	// TODO handle numbers
	case map[string]interface{}:
		outm := map[string]tftypes.Value{}
		attrtypes := map[string]tftypes.Type{}
		for k, vv := range v {
			if stringInSlice(k, ignoredFields) {
				continue
			}
			vvv := expand(vv, ignoredFields)
			if k == "metadata" { // FIXME we need to make this configurable
				tt := tftypes.List{ElementType: vvv.Type()}
				attrtypes[k] = tt
				outm[k] = tftypes.NewValue(tt, []tftypes.Value{vvv})
			} else {
				outm[k] = vvv
				attrtypes[k] = vvv.Type()
			}
		}
		stringMap := true
		for _, t := range attrtypes {
			if !t.Is(tftypes.String) {
				stringMap = false
			}
		}
		if stringMap {
			return tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, outm)
		}
		for k, vv := range outm {
			delete(outm, k)
			delete(attrtypes, k)

			outm[Snakify(k)] = vv
			attrtypes[Snakify(k)] = vv.Type()
		}
		return tftypes.NewValue(tftypes.Object{AttributeTypes: attrtypes}, outm)
	case []interface{}:
		outl := []tftypes.Value{}
		for _, vv := range v {
			outl = append(outl, expand(vv, ignoredFields))
		}
		// FIXME need to figure out if this is list or tuple
		return tftypes.NewValue(tftypes.List{ElementType: outl[0].Type()}, outl)
	}

	// FIXME return an error here
	return tftypes.Value{}
}

// Snakify converts a camelCase string into snake_case
func Snakify(in string) string {
	out := ""
	prevcap := false
	for _, ch := range in {
		cap := unicode.IsUpper(ch)
		if cap && !prevcap {
			out += "_"
		}
		out += strings.ToLower(string(ch))
		prevcap = cap
	}
	return out
}
