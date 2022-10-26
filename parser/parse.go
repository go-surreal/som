package parser

import (
	"errors"
	"fmt"
	"github.com/wzshiming/gotype"
	"go/ast"
	"golang.org/x/mod/modfile"
	"os"
	"path"
)

const fileGoMod = "go.mod"

const packagePath = "github.com/marcbinz/sdb"

func Parse(dir string) (*Result, error) {
	res := &Result{}

	imp := gotype.NewImporter()

	n, err := imp.Import(dir, "")
	if err != nil {
		return nil, err
	}

	pkgPath, err := parseMod(dir)
	if err != nil {
		return nil, err
	}

	res.PkgPath = path.Join(pkgPath, "example", "model") // TODO

	nc := n.NumChild()
	for i := 0; i < nc; i++ {
		v := n.Child(i)

		if !ast.IsExported(v.Name()) {
			continue
		}

		if v.Kind() == gotype.Struct {
			done := false
			nf := v.NumField()
			for j := 0; j < nf; j++ {
				f := v.Field(j)

				if f.Name() == "Node" && f.Elem().PkgPath() == packagePath {
					node := Node{Name: v.Name()}

					for k := 0; k < nf; k++ {
						f2 := v.Field(k)

						if k == j {
							continue
						}

						// TODO: ignore unexported fields?!

						var field Field

						// prevent ID from not being a string type
						if f2.Name() == "ID" && f2.Elem().Kind() != gotype.String {
							return nil, fmt.Errorf("field ID of model %s must be a string", f.Name())
						}

						switch f2.Elem().Kind() {
						case gotype.String:
							if f2.Name() == "ID" {
								field = FieldID{fieldAtomic{Name: f2.Name()}}
							} else if f2.Elem().String() != "string" && f2.Elem().PkgPath() == "github.com/marcbinz/sdb" { // TODO: might not be an enum..?!
								field = FieldEnum{fieldAtomic{Name: f2.Name()}, f2.Elem().String()}
							} else {
								field = FieldString{fieldAtomic{Name: f2.Name()}}
							}
						case gotype.Int:
							field = FieldInt{fieldAtomic{Name: f2.Name()}}
						case gotype.Int32:
							field = FieldInt32{fieldAtomic{Name: f2.Name()}}
						case gotype.Int64:
							field = FieldInt64{fieldAtomic{Name: f2.Name()}}
						case gotype.Float32:
							field = FieldFloat32{fieldAtomic{Name: f2.Name()}}
						case gotype.Float64:
							field = FieldFloat64{fieldAtomic{Name: f2.Name()}}
						case gotype.Bool:
							field = FieldBool{fieldAtomic{Name: f2.Name()}}
						case gotype.Struct:
							// TODO: prevent structs (or general types) from another package (except time and uuid)!
							if f2.Elem().PkgPath() == "time" {
								field = FieldTime{fieldAtomic{Name: f2.Name()}}
							} else if f2.Elem().PkgPath() == "github.com/google/uuid" {
								field = FieldUUID{fieldAtomic{Name: f2.Name()}}
							} else {
								field = FieldStruct{fieldAtomic{Name: f2.Name()}, false} // TODO: handle pointers
							}
						case gotype.Slice:
							field = FieldSlice{fieldAtomic{Name: f2.Name()}, f2.Elem().Elem().Name()}
						case gotype.Map:
							field = FieldMap{fieldAtomic{Name: f2.Name()}, f2.Elem().Key().Name(), f2.Elem().Elem().Name()}
						default:
							return nil, fmt.Errorf("field %s has unsupported type %s", f2.Name(), f2.Elem().Kind())
						}

						node.Fields = append(node.Fields, field)
					}

					res.Nodes = append(res.Nodes, node)
					done = true
					continue
				}
			}

			if done {
				continue
			}
		}

		if v.Kind() == gotype.String {
			if v.PkgPath() == packagePath {
				res.Enums = append(res.Enums, Enum{
					Name: v.Name(),
				})
				continue
			}
		}

		if v.Kind() == gotype.Declaration {
			res.EnumValues = append(res.EnumValues, EnumValue{
				Enum:     v.Declaration().Name(),
				Variable: v.Name(),
				Value:    v.Value(),
			})
		}
	}

	return res, nil
}

type Result struct {
	PkgPath    string
	Nodes      []Node
	Structs    []Struct
	Enums      []Enum
	EnumValues []EnumValue
}

//
// -- TOP LEVEL
//

type Node struct {
	Name   string
	Fields []Field
}

type Struct struct {
	Name string
}

type Enum struct {
	Name string
}

type EnumValue struct {
	Enum     string
	Variable string
	Value    string
}

//
// -- FIELDS
//

type Field interface {
	field()
	GetName() string
}

type isField struct{}

func (isField) field() {}

type fieldAtomic struct {
	isField
	Name string
}

func (f fieldAtomic) GetName() string {
	return f.Name
}

type FieldID struct {
	fieldAtomic
}

type FieldString struct {
	fieldAtomic
}

type FieldInt struct {
	fieldAtomic
}

type FieldInt32 struct {
	fieldAtomic
}

type FieldInt64 struct {
	fieldAtomic
}

type FieldFloat32 struct {
	fieldAtomic
}

type FieldFloat64 struct {
	fieldAtomic
}

type FieldBool struct {
	fieldAtomic
}

type FieldTime struct {
	fieldAtomic
}

type FieldUUID struct {
	fieldAtomic
}

type FieldEnum struct {
	fieldAtomic
	Typ string
}

type FieldStruct struct {
	fieldAtomic
	Pointer bool
}

type FieldSlice struct {
	fieldAtomic
	Value string
}

type FieldMap struct {
	fieldAtomic
	Key   string
	Value string
}

//
// -- HELPER
//

func findAndReadModFile(dir string) ([]byte, string, error) {
	for dir != "" {
		data, err := os.ReadFile(path.Join(dir, fileGoMod))

		if err == nil {
			return data, dir, nil
		}

		if !errors.Is(err, os.ErrNotExist) {
			return nil, "", err
		}

		dir = path.Dir(dir)
	}

	return nil, "", errors.New("could not find go.mod in worktree")
}

func parseMod(dir string) (string, error) {
	data, _, err := findAndReadModFile(dir)
	if err != nil {
		return "", err
	}

	f, err := modfile.Parse(fileGoMod, data, nil)
	if err != nil {
		return "", err
	}

	return f.Module.Mod.Path, nil
}
