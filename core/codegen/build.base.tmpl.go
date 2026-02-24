package codegen

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/codegen/field"
	"github.com/go-surreal/som/core/parser"
)

func fieldValueStr(in *input, basePkg string, sf parser.ComplexIDField, accessor string) string {
	switch f := sf.Field.(type) {
	case *parser.FieldTime:
		return "&types.DateTime{Time: " + accessor + "}"
	case *parser.FieldDuration:
		return "&types.Duration{Duration: " + accessor + "}"
	case *parser.FieldNode:
		refNode := in.findNodeByName(f.Node)
		if refNode == nil {
			return accessor
		}
		tableName := refNode.NameDatabase()
		idVal := nodeRefValueStr(in, basePkg, refNode, accessor)
		return fmt.Sprintf("models.NewRecordID(%q, %s)", tableName, idVal)
	default:
		return accessor
	}
}

func nodeRefValueStr(in *input, basePkg string, refNode *field.NodeTable, accessor string) string {
	if !refNode.HasComplexID() {
		return "string(" + accessor + ".ID())"
	}
	cid := refNode.Source.ComplexID
	innerAccessor := accessor + ".ID()"
	if cid.Kind == parser.IDTypeArray {
		var elems []string
		for _, sf := range cid.Fields {
			elems = append(elems, fieldValueStr(in, basePkg, sf, innerAccessor+"."+sf.Name))
		}
		return "[]any{" + strings.Join(elems, ", ") + "}"
	}
	return mapValueStr(in, basePkg, cid.Fields, innerAccessor)
}

func recordIDValueStr(in *input, basePkg string, node *field.NodeTable, keyVar string) string {
	cid := node.Source.ComplexID
	if cid.Kind == parser.IDTypeArray {
		var elems []string
		for _, sf := range cid.Fields {
			elems = append(elems, fieldValueStr(in, basePkg, sf, keyVar+"."+sf.Name))
		}
		return "[]any{" + strings.Join(elems, ", ") + "}"
	}
	return mapValueStr(in, basePkg, cid.Fields, keyVar)
}

func mapValueStr(in *input, basePkg string, fields []parser.ComplexIDField, keyVar string) string {
	type entry struct {
		key string
		val string
	}
	var entries []entry
	for _, sf := range fields {
		entries = append(entries, entry{
			key: fmt.Sprintf("%q", sf.DBName),
			val: fieldValueStr(in, basePkg, sf, keyVar+"."+sf.Name),
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].key < entries[j].key
	})
	// Compute max key width for alignment
	maxKeyLen := 0
	for _, e := range entries {
		if len(e.key) > maxKeyLen {
			maxKeyLen = len(e.key)
		}
	}
	var lines []string
	for _, e := range entries {
		padding := strings.Repeat(" ", maxKeyLen-len(e.key))
		lines = append(lines, fmt.Sprintf("%s:%s %s", e.key, padding, e.val))
	}
	return "map[string]any{\n" + strings.Join(lines, ",\n") + ",\n}"
}

type hookEventData struct {
	Event       string
	EventLower  string
	TimingUpper string
	TimingLower string
	Comment     string
}

func (b *build) buildBaseFileFromTemplate(node *field.NodeTable) error {
	data := b.buildBaseTemplateData(node)
	return renderGoFile(
		b.fs.Writer(filepath.Join(def.PkgRepo, node.FileName())),
		def.PkgRepo,
		"baseFile",
		baseFileTmpl,
		data,
	)
}

func (b *build) buildBaseTemplateData(node *field.NodeTable) map[string]any {
	pkgQuery := b.relativePkgPath(def.PkgQuery)
	pkgConv := b.relativePkgPath(def.PkgConv)
	somPkg := b.relativePkgPath()
	internalPkg := b.relativePkgPath("internal")
	relatePkg := b.relativePkgPath(def.PkgRelate)

	nameGo := node.NameGo()
	nameGoLower := node.NameGoLower()

	// Key type string
	keyType := "string"
	if node.HasComplexID() {
		keyType = "model." + node.Source.ComplexID.StructName
	}

	// RecordID func string
	var recordIDFunc string
	if node.HasComplexID() {
		ridValue := recordIDValueStr(b.input, b.basePkg(), node, "key")
		recordIDFunc = fmt.Sprintf("func(key %s) *models.RecordID {\n\t\t\t\trid := models.NewRecordID(%q, %s)\n\t\t\t\treturn &rid\n\t\t\t}",
			keyType, node.NameDatabase(), ridValue)
	} else {
		parseFn := "parseStringID"
		if node.Source.IDType == parser.IDTypeUUID {
			parseFn = "parseUUID"
		}
		recordIDFunc = fmt.Sprintf("func(id string) *models.RecordID {\n\t\t\t\trid := models.NewRecordID(%q, %s(id))\n\t\t\t\treturn &rid\n\t\t\t}",
			node.NameDatabase(), parseFn)
	}

	// RecordIDFromVar - expression for getting recordID from node variable
	var recordIDFromVar string
	if node.HasComplexID() {
		recordIDFromVar = "r.recordID(" + nameGoLower + ".ID())"
	} else {
		recordIDFromVar = "r.recordID(string(" + nameGoLower + ".ID()))"
	}

	// ID check function: takes (varName, errMsg) returns Go code
	idCheck := b.idCheckFunc(node)

	// Hook code functions
	beforeHooks := func(event string) string {
		return hooksStr(nameGoLower, "som", "Before", event)
	}
	afterHooks := func(event string) string {
		return hooksStr(nameGoLower, "som", "After", event)
	}

	// Imports
	imports := []goImport{
		{Path: "context"},
		{Path: "errors"},
		{Alias: "som", Path: somPkg},
		{Alias: "conv", Path: pkgConv},
		{Alias: "internal", Path: internalPkg},
		{Alias: "query", Path: pkgQuery},
		{Alias: "model", Path: b.input.sourcePkgPath},
		{Alias: "models", Path: def.PkgModels},
		{Path: "slices"},
		{Path: "sync"},
		{Path: "sync/atomic"},
	}

	if !node.HasComplexID() {
		imports = append(imports, goImport{Alias: "relate", Path: relatePkg})
	}

	if node.Source.SoftDelete {
		imports = append(imports, goImport{Path: "fmt"})
	}

	if node.Source.OptimisticLock {
		imports = append(imports, goImport{Path: "strings"})
	}

	if node.HasComplexID() && b.complexIDNeedsTypes(node) {
		imports = append(imports, goImport{Alias: "types", Path: b.relativePkgPath(def.PkgTypes)})
	}

	// Hook events for interface and method generation
	var hookEvents []hookEventData
	for _, event := range []string{"Create", "Update", "Delete"} {
		for _, timing := range []string{"Before", "After"} {
			methodName := "On" + timing + event
			var hookComment string
			switch timing {
			case "Before":
				hookComment = methodName + " registers a hook that runs before a record is " + strings.ToLower(event) + "d.\n" +
					"If the hook returns an error, the " + strings.ToLower(event) + " operation is aborted.\n" +
					"Returns a function that, when called, removes this hook.\n" +
					"\n" +
					"Note: Hooks are local to this application instance and are not\n" +
					"distributed across multiple instances of the application."
			case "After":
				hookComment = methodName + " registers a hook that runs after a record has been " + strings.ToLower(event) + "d.\n" +
					"If the hook returns an error, the error is returned to the caller.\n" +
					"Returns a function that, when called, removes this hook.\n" +
					"\n" +
					"Note: Hooks are local to this application instance and are not\n" +
					"distributed across multiple instances of the application."
			}
			hookEvents = append(hookEvents, hookEventData{
				Event:       event,
				EventLower:  strings.ToLower(event),
				TimingUpper: timing,
				TimingLower: strings.ToLower(timing),
				Comment:     formatGoComment(hookComment),
			})
		}
	}

	return map[string]any{
		"ImportBlock":    formatImportBlock(imports),
		"NameGo":        nameGo,
		"NameGoLower":   nameGoLower,
		"NameDB":        node.NameDatabase(),
		"HasComplexID":  node.HasComplexID(),
		"SoftDelete":    node.Source.SoftDelete,
		"OptimisticLock": node.Source.OptimisticLock,
		"KeyType":       keyType,
		"NewIDFunc":     idFuncName(node),
		"RecordIDFunc":  recordIDFunc,
		"RecordIDFromVar": recordIDFromVar,
		"IDCheck":       idCheck,
		"BeforeHooks":   beforeHooks,
		"AfterHooks":    afterHooks,
		"HookEvents":    hookEvents,
	}
}

func (b *build) idCheckFunc(node *field.NodeTable) func(string, string) string {
	if node.HasComplexID() {
		cid := node.Source.ComplexID
		if cid.HasNodeRef() {
			return func(varName, _ string) string {
				return b.nodeRefChecksStr(cid, varName)
			}
		}
		keyType := "model." + cid.StructName
		return func(varName, errMsg string) string {
			return fmt.Sprintf("var zeroKey %s\nif %s.ID() == zeroKey {\nreturn errors.New(%q)\n}", keyType, varName, errMsg)
		}
	}
	return func(varName, errMsg string) string {
		return fmt.Sprintf("if %s.ID() == \"\" {\nreturn errors.New(%q)\n}", varName, errMsg)
	}
}

func (b *build) nodeRefChecksStr(cid *parser.FieldComplexID, varName string) string {
	var parts []string
	for _, sf := range cid.Fields {
		fn, ok := sf.Field.(*parser.FieldNode)
		if !ok {
			continue
		}
		refNode := b.input.findNodeByName(fn.Node)
		if refNode == nil {
			continue
		}
		fieldErrMsg := sf.Name + ".ID must not be empty"
		accessor := varName + ".ID()." + sf.Name
		if !refNode.HasComplexID() {
			parts = append(parts, fmt.Sprintf("if %s.ID() == \"\" {\nreturn errors.New(%q)\n}", accessor, fieldErrMsg))
		} else if !refNode.Source.ComplexID.HasNodeRef() {
			zeroVar := "zero" + sf.Name + "Key"
			parts = append(parts, fmt.Sprintf("var %s model.%s\nif %s.ID() == %s {\nreturn errors.New(%q)\n}",
				zeroVar, refNode.Source.ComplexID.StructName, accessor, zeroVar, fieldErrMsg))
		}
	}
	return strings.Join(parts, "\n")
}

func hooksStr(nameGoLower, somAlias, timing, event string) string {
	hookIface := "On" + timing + event + "Hook"
	fieldName := strings.ToLower(timing) + event
	return fmt.Sprintf(`if h, ok := any(%s).(%s.%s); ok {
	if err := h.On%s%s(ctx); err != nil {
		return err
	}
}
r.mu.RLock()
%sHooks := make([]%sHook, len(r.%s))
copy(%sHooks, r.%s)
r.mu.RUnlock()
for _, h := range %sHooks {
	if err := h.fn(ctx, %s); err != nil {
		return err
	}
}`,
		nameGoLower, somAlias, hookIface,
		timing, event,
		fieldName, nameGoLower, fieldName,
		fieldName, fieldName,
		fieldName, nameGoLower)
}

func (b *build) complexIDNeedsTypes(node *field.NodeTable) bool {
	if !node.HasComplexID() {
		return false
	}
	return fieldsNeedTypes(b.input, node.Source.ComplexID.Fields)
}

func fieldsNeedTypes(in *input, fields []parser.ComplexIDField) bool {
	for _, sf := range fields {
		switch sf.Field.(type) {
		case *parser.FieldTime, *parser.FieldDuration:
			return true
		case *parser.FieldNode:
			fn := sf.Field.(*parser.FieldNode)
			refNode := in.findNodeByName(fn.Node)
			if refNode != nil && refNode.HasComplexID() {
				if fieldsNeedTypes(in, refNode.Source.ComplexID.Fields) {
					return true
				}
			}
		}
	}
	return false
}
