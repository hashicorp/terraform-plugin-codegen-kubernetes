package main

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/lmittmann/tint"

	"github.com/hashicorp/terraform-plugin-codegen-kubernetes/internal/generator"
)

var configFilePattern = regexp.MustCompile(`generate(_.+)?\.hcl`)

func main() {
	// setup slog with colour to make it easier to read
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
		}),
	))

	generateFiles := []string{}
	filepath.Walk("./", func(path string, info fs.FileInfo, err error) error {
		filename := filepath.Base(path)
		if configFilePattern.MatchString(filename) {
			slog.Info("Found generator config file", "filename", filename)
			generateFiles = append(generateFiles, path)
		}
		return nil
	})

	generatedResources := []generator.ResourceConfig{}
	for _, f := range generateFiles {
		config, err := generator.ParseHCLConfig(f)
		if err != nil {
			slog.Error("Error parsing configuration", "filename", f, "err", err)
			os.Exit(1)
		}
		resources, err := generateFrameworkCode(f, config)
		if err != nil {
			slog.Error("Error generating framework code", "err", err)
			os.Exit(1)
		}
		generatedResources = append(generatedResources, resources...)
	}

	// generate resources list file
	resourcesList := generator.ResourcesListGenerator{
		time.Now(),
		generatedResources,
		generatePackageList(generatedResources),
	}
	outputFilename := "resources_list_gen.go"
	generator.WriteFormattedSourceFile("./provider", outputFilename, resourcesList.String())
	slog.Info("Generated resources list source file", "filename", outputFilename)
}

func generatePackageList(resources []generator.ResourceConfig) []string {
	packages := []string{}
	packageMap := map[string]struct{}{}
	for _, r := range resources {
		packageMap[r.Package] = struct{}{}
	}
	for k := range packageMap {
		packages = append(packages, k)
	}
	return packages
}

func generateFrameworkCode(path string, config generator.GeneratorConfig) ([]generator.ResourceConfig, error) {
	wd := filepath.Dir(path)

	generatedResources := []generator.ResourceConfig{}
	for _, r := range config.Resources {
		if r.Disabled {
			slog.Warn("Code generation is disabled, skipping", "resource", r.Name)
			continue
		}
		slog.Info("Generating framework code", "resource", r.Name)
		spec, err := generator.GenerateResourceSpec(r)
		if err != nil {
			return nil, fmt.Errorf("error generating provider spec: %v", err)
		}

		gen := generator.NewResourceGenerator(r, spec)

		// generate resource
		resourceCode := gen.GenerateResourceCode()
		outputFilename := fmt.Sprintf("%s_gen.go", r.OutputFilenamePrefix)
		generator.WriteFormattedSourceFile(wd, outputFilename, resourceCode)
		slog.Info("Generated resource source file", "filename", outputFilename)

		// generate schema
		if r.Generate.Schema {
			schemaCode := gen.GenerateSchemaFunctionCode()
			outputFilename = fmt.Sprintf("%s_schema_gen.go", r.OutputFilenamePrefix)
			generator.WriteFormattedSourceFile(wd, outputFilename, schemaCode)
			slog.Info("Generated schema source file", "filename", outputFilename)
		}

		// generate CRUD stubs
		if r.Generate.CRUDStubs {
			crudStubCode := gen.GenerateCRUDStubCode()
			outputFilename = fmt.Sprintf("%s_crud.go", r.OutputFilenamePrefix)
			generator.WriteFormattedSourceFile(wd, outputFilename, crudStubCode)
			slog.Info("Generated CRUD stub source file", "filename", outputFilename)
		}

		// generate auto CRUD functions
		if r.Generate.CRUDAuto {
			crudStubCode := gen.GenerateAutoCRUDCode()
			outputFilename = fmt.Sprintf("%s_crud_gen.go", r.OutputFilenamePrefix)
			generator.WriteFormattedSourceFile(wd, outputFilename, crudStubCode)
			slog.Info("Generated autocrud source file", "filename", outputFilename)
		}

		// generate model
		if r.Generate.Model {
			crudStubCode := gen.GenerateModelCode()
			outputFilename = fmt.Sprintf("%s_model_gen.go", r.OutputFilenamePrefix)
			generator.WriteFormattedSourceFile(wd, outputFilename, crudStubCode)
			slog.Info("Generated model source file", "filename", outputFilename)
		}

		generatedResources = append(generatedResources, r)
	}
	return generatedResources, nil
}
