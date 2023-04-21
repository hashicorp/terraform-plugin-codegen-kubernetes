package provider

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// FlattenValue takes a tftypes.Value and converts it to an unstructured Kubernetes manifest
func FlattenValue(in tftypes.Value) interface{} {
	return flattenValue(in, false)
}

func flattenValue(in tftypes.Value, unwrap bool) interface{} {
	switch {
	case in.Type().Is(tftypes.String):
		var s string
		in.As(&s)
		if s == "" {
			return nil
		}
		return s
	case in.Type().Is(tftypes.Bool):
		var b bool
		in.As(&b)
		return b
	case in.Type().Is(tftypes.Number):
		// TODO

	case in.Type().Is(tftypes.List{}) || in.Type().Is(tftypes.Tuple{}):
		var l []tftypes.Value
		in.As(&l)
		outl := []interface{}{}
		for _, v := range l {
			fv := flattenValue(v, false)
			outl = append(outl, fv)
		}
		if unwrap {
			return outl[0]
		}
		if len(outl) == 0 {
			return nil
		}
		return outl
	case in.Type().Is(tftypes.Map{}) || in.Type().Is(tftypes.Object{}):
		var m map[string]tftypes.Value
		in.As(&m)
		outm := map[string]interface{}{}
		for k, v := range m {
			unwrap := k == "metadata" // unwrap metadata from list
			if vv := flattenValue(v, unwrap); vv != nil {
				outm[Camelize(k)] = vv
			}
		}
		if len(outm) == 0 {
			return nil
		}
		return outm
	default:
		panic("unexpected type")
	}

	return nil
}

// Camelize converts a string containing snake_case into camelCase
// FIXME this wont work for variables containing ancronyms, e.g: pod_cidr, cluster_ip
// we should add a map of overrides so we can explicitly convert these
func Camelize(in string) string {
	out := ""
	cap := false
	for _, ch := range in {
		if ch == '_' {
			cap = true
			continue
		}
		if cap {
			out += strings.ToUpper(string(ch))
			cap = false
		} else {
			out += string(ch)
		}
	}
	return out
}
