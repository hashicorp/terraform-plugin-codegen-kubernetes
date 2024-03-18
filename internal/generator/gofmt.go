package generator

import (
	"go/format"
	"os"
	"path/filepath"
)

// WriteFormattedSourceFile runs Go code through format before writing to a file
func WriteFormattedSourceFile(wd, path string, contents string) error {
	src, err := format.Source([]byte(contents))
	if err != nil {
		return err
	}
	outputPath := filepath.Join(wd, path)
	return os.WriteFile(outputPath, src, os.ModePerm)
}
