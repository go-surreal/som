package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-surreal/som/core/util/gomod"
)

// AnalyzerDef represents an analyzer definition.
type AnalyzerDef struct {
	Name       string
	Tokenizers []string
	Filters    []FilterDef
}

// FilterDef represents a filter with optional parameters.
type FilterDef struct {
	Name   string
	Params []any // string or int/float parameters
}

// SearchDef represents a search configuration.
type SearchDef struct {
	Name         string
	AnalyzerName string
	BM25K1       float64
	BM25B        float64
	HasBM25      bool
	Highlights   bool
	Concurrently bool
}

// ViewDef represents a read-only view definition (DEFINE TABLE ... AS SELECT).
type ViewDef struct {
	// View is the Go name of the view model the definition applies to.
	View string
	// Source is the Go name of the source model the view selects from.
	Source string
	// Projections are the rendered "<expr> AS <column>" select entries.
	Projections []string
	// Where is the rendered (literal) filter condition, if any.
	Where string
	// GroupBy holds the rendered group-by field expressions, if any.
	GroupBy []string
}

// DefineOutput holds all parsed configuration from //go:build som files.
type DefineOutput struct {
	Analyzers []AnalyzerDef
	Searches  []SearchDef
	Views     []ViewDef
}

// defineOutputJSON matches the JSON structure from Definitions.ToJSON().
type defineOutputJSON struct {
	Analyzers []analyzerJSON `json:"analyzers"`
	Searches  []searchJSON   `json:"searches"`
	Views     []viewJSON     `json:"views"`
}

type viewJSON struct {
	View        string   `json:"view"`
	Source      string   `json:"source"`
	Projections []string `json:"projections"`
	Where       string   `json:"where,omitempty"`
	GroupBy     []string `json:"group_by,omitempty"`
}

type analyzerJSON struct {
	Name       string       `json:"name"`
	Tokenizers []string     `json:"tokenizers"`
	Filters    []filterJSON `json:"filters"`
}

type filterJSON struct {
	Name   string `json:"name"`
	Params []any  `json:"params,omitempty"`
}

type searchJSON struct {
	Name         string  `json:"name"`
	AnalyzerName string  `json:"analyzer_name"`
	BM25K1       float64 `json:"bm25_k1,omitempty"`
	BM25B        float64 `json:"bm25_b,omitempty"`
	HasBM25      bool    `json:"has_bm25"`
	Highlights   bool    `json:"highlights"`
	Concurrently bool    `json:"concurrently"`
}

// parseDefine parses the //go:build som definition file(s) by compiling and
// running the user's Definitions() function.
//
// Definitions live in the module-root package (next to go.mod), not in the
// model package: a view definition references the generated filter refs,
// which import the model package. Placing the definition in the model package
// would create an import cycle (model -> filter -> model). The module root can
// import both the model and the generated packages without a cycle.
func parseDefine(dir string) (*DefineOutput, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("could not get absolute path: %w", err)
	}

	// Get the module path; definitions are resolved against the module root.
	mod, err := gomod.FindGoMod(absDir)
	if err != nil {
		return nil, fmt.Errorf("could not find go.mod: %w", err)
	}

	rootDir := mod.Dir()
	modelPkg := mod.Module()

	// Check if the module root has any //go:build som file.
	hasDefine, err := hasDefineFiles(rootDir)
	if err != nil {
		return nil, err
	}
	if !hasDefine {
		return &DefineOutput{}, nil
	}

	// Create temp directory for main.go
	tempDir, err := os.MkdirTemp(rootDir, ".som_temp_")
	if err != nil {
		return nil, fmt.Errorf("could not create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Write temp main.go
	mainContent := fmt.Sprintf(`//go:build som

package main

import (
	"os"
	model "%s"
)

func main() {
	data, err := model.Definitions().ToJSON()
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
	os.Stdout.Write(data)
}
`, modelPkg)

	mainPath := filepath.Join(tempDir, "main.go")
	if err := os.WriteFile(mainPath, []byte(mainContent), 0644); err != nil {
		return nil, fmt.Errorf("could not write temp main.go: %w", err)
	}

	// Run go run -tags=som
	cmd := exec.Command("go", "run", "-tags=som", mainPath)
	cmd.Dir = mod.Dir()
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("failed to run Definitions(): %s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("failed to run Definitions(): %w", err)
	}

	// Parse JSON output
	var jsonOutput defineOutputJSON
	if err := json.Unmarshal(output, &jsonOutput); err != nil {
		return nil, fmt.Errorf("could not parse Definitions() output: %w", err)
	}

	// Convert to DefineOutput
	result := &DefineOutput{}

	for _, a := range jsonOutput.Analyzers {
		analyzer := AnalyzerDef{
			Name:       a.Name,
			Tokenizers: a.Tokenizers,
		}
		for _, f := range a.Filters {
			analyzer.Filters = append(analyzer.Filters, FilterDef(f))
		}
		result.Analyzers = append(result.Analyzers, analyzer)
	}

	for _, s := range jsonOutput.Searches {
		result.Searches = append(result.Searches, SearchDef(s))
	}

	for _, v := range jsonOutput.Views {
		result.Views = append(result.Views, ViewDef{
			View:        v.View,
			Source:      v.Source,
			Projections: v.Projections,
			Where:       v.Where,
			GroupBy:     v.GroupBy,
		})
	}

	return result, nil
}

// hasDefineFiles checks if the directory contains any //go:build som files.
func hasDefineFiles(dir string) (bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false, fmt.Errorf("could not read directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}

		filePath := filepath.Join(dir, entry.Name())
		hasBuildTag, err := hasGoBuildSomTag(filePath)
		if err != nil {
			return false, err
		}
		if hasBuildTag {
			return true, nil
		}
	}
	return false, nil
}

// hasGoBuildSomTag checks if a file has the //go:build som build constraint.
func hasGoBuildSomTag(filePath string) (bool, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false, err
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Check for new-style build tag
		if strings.HasPrefix(line, "//go:build") && strings.Contains(line, "som") {
			return true, nil
		}
		// Check for old-style build tag
		if strings.HasPrefix(line, "// +build") && strings.Contains(line, "som") {
			return true, nil
		}
		// Stop after package declaration
		if strings.HasPrefix(line, "package ") {
			break
		}
	}
	return false, nil
}
