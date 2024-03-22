package autocrud

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type DiagnosticsError struct {
	diag.Diagnostics
}

func (d DiagnosticsError) Error() string {
	if !d.HasError() {
		panic("No errors in diagnostics")
	}
	if d.ErrorsCount() == 1 {
		return d.Errors()[0].Detail()
	}
	message := "The following errors occured:"
	for _, e := range d.Errors() {
		message += "\n" + e.Detail()
	}
	return message
}
