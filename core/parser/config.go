package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// AnalyzerDef represents a parsed analyzer definition from a //go:build som file.
type AnalyzerDef struct {
	Name       string
	VarName    string // The Go variable name (for reference in SearchDef)
	Tokenizers []string
	Filters    []FilterDef
}

// FilterDef represents a filter with optional parameters.
type FilterDef struct {
	Name   string
	Params []any // string or int/float parameters
}

// SearchDef represents a parsed search configuration from a //go:build som file.
type SearchDef struct {
	Name         string
	AnalyzerVar  string  // Variable name referencing an AnalyzerDef
	AnalyzerName string  // Resolved analyzer name
	BM25K1       float64 // BM25 k1 parameter (0 means not set)
	BM25B        float64 // BM25 b parameter (0 means not set)
	HasBM25      bool
	Highlights   bool
}

// ConfigOutput holds all parsed configuration from //go:build som files.
type ConfigOutput struct {
	Analyzers []AnalyzerDef
	Searches  []SearchDef
}

// ParseConfig parses all //go:build som files in the given directory
// and extracts analyzer and search configuration definitions.
func ParseConfig(dir string) (*ConfigOutput, error) {
	output := &ConfigOutput{}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("could not get absolute path: %w", err)
	}

	// Find all .go files with //go:build som tag
	entries, err := os.ReadDir(absDir)
	if err != nil {
		return nil, fmt.Errorf("could not read directory: %w", err)
	}

	fset := token.NewFileSet()

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}

		filePath := filepath.Join(absDir, entry.Name())

		// Check if file has //go:build som tag
		hasBuildTag, err := hasGoBuildSomTag(filePath)
		if err != nil {
			return nil, fmt.Errorf("could not check build tag for %s: %w", entry.Name(), err)
		}
		if !hasBuildTag {
			continue
		}

		// Parse the file
		f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("could not parse %s: %w", entry.Name(), err)
		}

		// Extract analyzer and search definitions
		analyzers, searches, err := extractDefinitions(f)
		if err != nil {
			return nil, fmt.Errorf("could not extract definitions from %s: %w", entry.Name(), err)
		}

		output.Analyzers = append(output.Analyzers, analyzers...)
		output.Searches = append(output.Searches, searches...)
	}

	// Resolve analyzer references in search definitions
	if err := resolveAnalyzerRefs(output); err != nil {
		return nil, err
	}

	return output, nil
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

// extractDefinitions extracts analyzer and search definitions from a parsed file.
func extractDefinitions(f *ast.File) ([]AnalyzerDef, []SearchDef, error) {
	var analyzers []AnalyzerDef
	var searches []SearchDef

	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			continue
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok || len(valueSpec.Values) == 0 {
				continue
			}

			varName := ""
			if len(valueSpec.Names) > 0 {
				varName = valueSpec.Names[0].Name
			}

			// Check if this is an Analyzer or Search call
			for _, value := range valueSpec.Values {
				if analyzer, ok := parseAnalyzerCall(value, varName); ok {
					analyzers = append(analyzers, analyzer)
				} else if search, ok := parseSearchCall(value); ok {
					searches = append(searches, search)
				}
			}
		}
	}

	return analyzers, searches, nil
}

// parseAnalyzerCall parses a define.Analyzer(...) or som.Analyzer(...) call chain.
func parseAnalyzerCall(expr ast.Expr, varName string) (AnalyzerDef, bool) {
	def := AnalyzerDef{VarName: varName}

	// Walk up the call chain (method calls are nested)
	current := expr
	for {
		call, ok := current.(*ast.CallExpr)
		if !ok {
			break
		}

		switch fn := call.Fun.(type) {
		case *ast.SelectorExpr:
			methodName := fn.Sel.Name

			switch methodName {
			case "Analyzer":
				// This is the root: define.Analyzer("name") or som.Analyzer("name")
				if len(call.Args) > 0 {
					if lit, ok := call.Args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
						def.Name, _ = strconv.Unquote(lit.Value)
					}
				}
				return def, def.Name != ""

			case "Tokenizers":
				// .Tokenizers(define.Blank, define.Punct, ...)
				for _, arg := range call.Args {
					if sel, ok := arg.(*ast.SelectorExpr); ok {
						def.Tokenizers = append(def.Tokenizers, strings.ToLower(sel.Sel.Name))
					}
				}
				current = fn.X

			case "Filters":
				// .Filters(define.Lowercase, define.Snowball("en"), ...)
				for _, arg := range call.Args {
					if filter, ok := parseFilterArg(arg); ok {
						def.Filters = append(def.Filters, filter)
					}
				}
				current = fn.X

			default:
				current = fn.X
			}

		case *ast.Ident:
			// Reached the end of the chain
			return def, false

		default:
			return def, false
		}
	}

	return def, false
}

// parseFilterArg parses a filter argument (either a selector like define.Lowercase
// or a call like define.Snowball("en")).
func parseFilterArg(arg ast.Expr) (FilterDef, bool) {
	switch v := arg.(type) {
	case *ast.SelectorExpr:
		// Simple filter like define.Lowercase
		return FilterDef{Name: strings.ToLower(v.Sel.Name)}, true

	case *ast.CallExpr:
		// Filter with parameters like define.Snowball("en")
		if sel, ok := v.Fun.(*ast.SelectorExpr); ok {
			filter := FilterDef{Name: strings.ToLower(sel.Sel.Name)}
			for _, arg := range v.Args {
				if lit, ok := arg.(*ast.BasicLit); ok {
					switch lit.Kind {
					case token.STRING:
						val, _ := strconv.Unquote(lit.Value)
						filter.Params = append(filter.Params, val)
					case token.INT:
						val, _ := strconv.Atoi(lit.Value)
						filter.Params = append(filter.Params, val)
					case token.FLOAT:
						val, _ := strconv.ParseFloat(lit.Value, 64)
						filter.Params = append(filter.Params, val)
					}
				}
			}
			return filter, true
		}
	}
	return FilterDef{}, false
}

// parseSearchCall parses a define.Search(...) or som.Search(...) call chain.
func parseSearchCall(expr ast.Expr) (SearchDef, bool) {
	def := SearchDef{}

	current := expr
	for {
		call, ok := current.(*ast.CallExpr)
		if !ok {
			break
		}

		switch fn := call.Fun.(type) {
		case *ast.SelectorExpr:
			methodName := fn.Sel.Name

			switch methodName {
			case "Search":
				// Root: define.Search("name")
				if len(call.Args) > 0 {
					if lit, ok := call.Args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
						def.Name, _ = strconv.Unquote(lit.Value)
					}
				}
				return def, def.Name != ""

			case "Analyzer":
				// .Analyzer(varName) - reference to an analyzer variable
				if len(call.Args) > 0 {
					if ident, ok := call.Args[0].(*ast.Ident); ok {
						def.AnalyzerVar = ident.Name
					}
				}
				current = fn.X

			case "BM25":
				// .BM25(k1, b)
				def.HasBM25 = true
				if len(call.Args) >= 2 {
					if lit, ok := call.Args[0].(*ast.BasicLit); ok {
						def.BM25K1, _ = strconv.ParseFloat(lit.Value, 64)
					}
					if lit, ok := call.Args[1].(*ast.BasicLit); ok {
						def.BM25B, _ = strconv.ParseFloat(lit.Value, 64)
					}
				}
				current = fn.X

			case "Highlights":
				// .Highlights()
				def.Highlights = true
				current = fn.X

			default:
				current = fn.X
			}

		case *ast.Ident:
			return def, false

		default:
			return def, false
		}
	}

	return def, false
}

// resolveAnalyzerRefs resolves analyzer variable references in search definitions.
func resolveAnalyzerRefs(output *ConfigOutput) error {
	// Build a map of variable name -> analyzer name
	varToName := make(map[string]string)
	for _, a := range output.Analyzers {
		if a.VarName != "" && a.VarName != "_" {
			varToName[a.VarName] = a.Name
		}
	}

	// Resolve references
	for i := range output.Searches {
		if output.Searches[i].AnalyzerVar != "" {
			name, ok := varToName[output.Searches[i].AnalyzerVar]
			if !ok {
				return fmt.Errorf("search config %q references unknown analyzer variable %q",
					output.Searches[i].Name, output.Searches[i].AnalyzerVar)
			}
			output.Searches[i].AnalyzerName = name
		}
	}

	return nil
}
