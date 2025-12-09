// Copyright IBM Corp. 2024
// SPDX-License-Identifier: MPL-2.0

package generator

import (
	"fmt"
	"go/format"
	"os"
	"path/filepath"
)

// WriteFormattedSourceFile runs Go code through format before writing to a file
func WriteFormattedSourceFile(wd, path string, contents string) error {
	src, err := format.Source([]byte(contents))
	outputPath := filepath.Join(wd, path)
	if err != nil {
		// if there's an error, write the unformattable Go code to the file so we cans see what broke
		writeErr := os.WriteFile(outputPath, []byte(contents), os.ModePerm)
		if writeErr != nil {
			return writeErr
		}
		return fmt.Errorf("failed to format Go file %q", outputPath)
	}
	return os.WriteFile(outputPath, src, os.ModePerm)
}
