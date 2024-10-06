package parser

import (
	"fmt"
	"github.com/go-surreal/som/exp/def"
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/packages"
	"log"
	"log/slog"
	"time"
)

const (
	mode = packages.NeedName | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo
)

type Parser struct {
	fileSet *token.FileSet

	imports []*def.Import
	nodes   []*def.Node
	edges   []*def.Edge
	structs []*def.Struct
}

func NewParser() *Parser {
	return &Parser{
		fileSet: token.NewFileSet(),
	}
}

func (p *Parser) Parse(path string) error {
	startTime := time.Now()
	defer func() {
		slog.Debug("Parsing completed.",
			"duration", time.Since(startTime).Round(time.Millisecond),
		)
	}()

	config := &packages.Config{
		Fset: p.fileSet,
		Mode: mode,
	}

	pkgs, err := packages.Load(config, path)
	if err != nil {
		return fmt.Errorf("could not load packages: %w", err)
	}

	for _, pkg := range pkgs {
		if err := p.parsePackage(pkg); err != nil {
			return fmt.Errorf("could not parse package: %w", err)
		}
	}

	//fileSet := token.NewFileSet()
	//
	//pkgs, err := parser.ParseDir(fileSet, path,
	//	func(info os.FileInfo) bool {
	//		return strings.HasSuffix(info.Name(), ".go")
	//	},
	//	parser.AllErrors,
	//)
	//if err != nil {
	//	return fmt.Errorf("could not parse code in source path: %w", err)
	//}
	//
	//for _, pkg := range pkgs {
	//	conf := types.Config{
	//		Importer: &customImporter{def: importer.Default()},
	//	}
	//
	//	// Create type-checking info
	//	info := &types.Info{
	//		Types:      make(map[ast.Expr]types.TypeAndValue),
	//		Defs:       make(map[*ast.Ident]types.Object),
	//		Uses:       make(map[*ast.Ident]types.Object),
	//		Implicits:  make(map[ast.Node]types.Object),
	//		Scopes:     make(map[ast.Node]*types.Scope),
	//		Selections: make(map[*ast.SelectorExpr]*types.Selection),
	//	}
	//
	//	// Create the package to type-check
	//	var files []*ast.File
	//	for _, file := range pkg.Files {
	//		files = append(files, file)
	//	}
	//
	//	// Type-check the package
	//	_, err := conf.Check(pkg.Name, fileSet, files, info)
	//	if err != nil {
	//		return fmt.Errorf("could not type-check package: %w", err)
	//	}
	//
	//	for a, b := range info.Types {
	//		fmt.Println(a, b)
	//	}
	//}

	//if len(pkgs) < 1 {
	//	return errors.New("no packages found in source path")
	//}
	//
	//if len(pkgs) > 1 {
	//	return errors.New("more than one package found in source path")
	//}
	//
	//for pkg := range maps.Values(pkgs) {
	//	if err := p.parse(pkg); err != nil {
	//		return fmt.Errorf("could not parse package: %w", err)
	//	}
	//}

	return nil
}

func (p *Parser) parsePackage(pkg *packages.Package) error {
	for _, fileAst := range pkg.Syntax {
		if err := p.parseFile(fileAst, pkg.TypesInfo); err != nil {
			return fmt.Errorf("could not parse file: %w", err)
		}
	}

	return nil
}

func (p *Parser) parseFile(fileAst *ast.File, info *types.Info) error {
	// TODO: use ast.Walk() instead?

	for expr, typ := range info.Types {
		fmt.Println(expr, typ)
	}

	for node := range ast.Preorder(fileAst) {
		switch mappedNode := node.(type) {

		case *ast.TypeSpec:
			if err := p.parseType(mappedNode, info); err != nil {
				return fmt.Errorf("could not parse type spec: %w", err)
			}
		}
	}

	return nil
}

func findInPackage(pkg *packages.Package, fset *token.FileSet) {
	for _, fileAst := range pkg.Syntax {
		ast.Inspect(fileAst, func(n ast.Node) bool {
			if structTy, ok := n.(*ast.StructType); ok {
				findInFields(structTy.Fields, n, pkg.TypesInfo, fset)
			} else if interfaceTy, ok := n.(*ast.InterfaceType); ok {
				findInFields(interfaceTy.Methods, n, pkg.TypesInfo, fset)
			}

			return true
		})
	}
}

func findInFields(fl *ast.FieldList, n ast.Node, tinfo *types.Info, fset *token.FileSet) {
	type FieldReport struct {
		Name string
		Kind string
		Type types.Type
	}
	var reps []FieldReport

	for _, field := range fl.List {
		if field.Names == nil {
			tv, ok := tinfo.Types[field.Type]
			if !ok {
				log.Fatal("not found", field.Type)
			}

			embName := fmt.Sprintf("%v", field.Type)

			_, hostIsStruct := n.(*ast.StructType)
			var kind string

			switch typ := tv.Type.Underlying().(type) {
			case *types.Struct:
				if hostIsStruct {
					kind = "struct (s@s)"
				} else {
					kind = "struct (s@i)"
				}
				reps = append(reps, FieldReport{embName, kind, typ})
			case *types.Interface:
				if hostIsStruct {
					kind = "interface (i@s)"
				} else {
					kind = "interface (i@i)"
				}
				reps = append(reps, FieldReport{embName, kind, typ})
			default:
			}
		}
	}

	if len(reps) > 0 {
		fmt.Printf("Found at %v\n", fset.Position(n.Pos()))

		for _, report := range reps {
			fmt.Printf("--> field '%s' is embedded %s: %s\n", report.Name, report.Kind, report.Type)
		}
		fmt.Println("")
	}
}

//type customImporter struct {
//	def types.Importer
//}
//
//func (i *customImporter) Import(path string) (*types.Package, error) {
//	if strings.HasPrefix(path, "github.com/go-surreal/som/exp/") {
//		list := filepath.SplitList(path)
//		return types.NewPackage(strings.TrimPrefix("github.com/go-surreal/som/exp/", path), list[len(list)-1]), nil
//	}
//
//	if path == "github.com/go-surreal/som" {
//		abs, _ := filepath.Abs("./..")
//		return types.NewPackage(abs, "som"), nil
//	}
//
//	return i.def.Import(path)
//}
//
//func isBasic(t types.Type) bool {
//	switch x := t.(type) {
//	case *types.Basic:
//		return true
//	case *types.Slice:
//		return true
//	case *types.Map:
//		return true
//	case *types.Pointer:
//		return isBasic(x.Elem())
//	default:
//		return false
//	}
//}

//func (p *Parser) parse(pkg ast.Node) error {
//	// TODO: use ast.Walk() instead?
//
//	for node := range ast.Preorder(pkg) {
//		switch mappedNode := node.(type) {
//
//		case *ast.TypeSpec:
//			if err := p.parseTypeSpec(mappedNode); err != nil {
//				return fmt.Errorf("could not parse type spec: %w", err)
//			}
//
//		case *ast.ImportSpec:
//			p.parseImportSpec(mappedNode)
//		}
//	}
//
//	return nil
//}

func (p *Parser) String() string {
	out := "\n"

	out += "---- Imports ----\n\n"

	for _, imp := range p.imports {
		out += fmt.Sprintf("%s\n", imp)
	}

	out += "\n---- Structs ----\n\n"

	for _, str := range p.structs {
		out += fmt.Sprintf("%s\n", str)
	}

	out += "\n---- Nodes ----\n\n"

	for _, node := range p.nodes {
		out += fmt.Sprintf("%s\n", node)
	}

	out += "\n---- Edges ----\n\n"

	for _, edge := range p.edges {
		out += fmt.Sprintf("%s\n", edge)
	}

	return out
}
