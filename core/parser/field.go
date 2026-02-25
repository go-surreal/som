package parser

import (
	"fmt"
)

type Field interface {
	fmt.Stringer
	field()
	Name() string
	Pointer() bool
	setName(string)
	setPointer(bool)
	Index() *IndexInfo
	Search() *SearchInfo
	setIndex(*IndexInfo)
	setSearch(*SearchInfo)
	Validate() error
}

type fieldAtomic struct {
	name    string
	pointer bool
	index   *IndexInfo
	search  *SearchInfo
}

func (*fieldAtomic) field() {}

func (f *fieldAtomic) String() string {
	return f.Name()
}

func (f *fieldAtomic) Name() string {
	return f.name
}

func (f *fieldAtomic) setName(name string) {
	f.name = name
}

func (f *fieldAtomic) Pointer() bool {
	return f.pointer
}

func (f *fieldAtomic) setPointer(val bool) {
	f.pointer = val
}

func (f *fieldAtomic) Index() *IndexInfo {
	return f.index
}

func (f *fieldAtomic) setIndex(info *IndexInfo) {
	f.index = info
}

func (f *fieldAtomic) Search() *SearchInfo {
	return f.search
}

func (f *fieldAtomic) setSearch(info *SearchInfo) {
	f.search = info
}

func (f *fieldAtomic) Validate() error {
	if f.search != nil {
		return fmt.Errorf("field %s: fulltext index only supports string types (string, *string, []string, []*string, *[]string, *[]*string)", f.name)
	}
	return nil
}

type IDType string

const (
	IDTypeULID   IDType = "ULID"
	IDTypeUUID   IDType = "UUID"
	IDTypeRand   IDType = "Rand"
	IDTypeArray  IDType = "Array"
	IDTypeObject IDType = "Object"
)

type FieldID struct {
	*fieldAtomic
	Type IDType
}

func NewFieldID(name string, idType IDType) *FieldID {
	return &FieldID{fieldAtomic: &fieldAtomic{name: name}, Type: idType}
}

type FieldString struct {
	*fieldAtomic
}

func NewFieldString(name string) *FieldString {
	return &FieldString{fieldAtomic: &fieldAtomic{name: name}}
}

func (f *FieldString) Validate() error {
	return nil // string and *string support fulltext
}

type FieldNumeric struct {
	*fieldAtomic
	Type NumberType
}

func NewFieldNumeric(name string, numType NumberType) *FieldNumeric {
	return &FieldNumeric{fieldAtomic: &fieldAtomic{name: name}, Type: numType}
}

type NumberType int32

const (
	NumberInt NumberType = iota
	NumberInt8
	NumberInt16
	NumberInt32
	NumberInt64
	//NumberUint
	NumberUint8
	NumberUint16
	NumberUint32
	//NumberUint64
	//NumberUintptr
	NumberFloat32
	NumberFloat64
	NumberRune
)

type FieldBool struct {
	*fieldAtomic
}

func NewFieldBool(name string) *FieldBool {
	return &FieldBool{fieldAtomic: &fieldAtomic{name: name}}
}

type FieldByte struct {
	*fieldAtomic
}

func NewFieldByte(name string) *FieldByte {
	return &FieldByte{fieldAtomic: &fieldAtomic{name: name}}
}

type FieldDuration struct {
	*fieldAtomic
}

func NewFieldDuration(name string) *FieldDuration {
	return &FieldDuration{fieldAtomic: &fieldAtomic{name: name}}
}

type FieldMonth struct {
	*fieldAtomic
}

func NewFieldMonth(name string) *FieldMonth {
	return &FieldMonth{fieldAtomic: &fieldAtomic{name: name}}
}

type FieldWeekday struct {
	*fieldAtomic
}

func NewFieldWeekday(name string) *FieldWeekday {
	return &FieldWeekday{fieldAtomic: &fieldAtomic{name: name}}
}

type FieldTime struct {
	*fieldAtomic
	IsCreatedAt bool
	IsUpdatedAt bool
	IsDeletedAt bool
}

func NewFieldTime(name string) *FieldTime {
	return &FieldTime{fieldAtomic: &fieldAtomic{name: name}}
}

type UUIDPackage string

const (
	UUIDPackageGoogle UUIDPackage = "github.com/google/uuid"
	UUIDPackageGofrs  UUIDPackage = "github.com/gofrs/uuid"
)

type FieldUUID struct {
	*fieldAtomic
	Package UUIDPackage
}

func NewFieldUUID(name string, pkg UUIDPackage) *FieldUUID {
	return &FieldUUID{
		fieldAtomic: &fieldAtomic{
			name: name,
		},
		Package: pkg,
	}
}

type FieldURL struct {
	*fieldAtomic
}

func NewFieldURL(name string) *FieldURL {
	return &FieldURL{fieldAtomic: &fieldAtomic{name: name}}
}

type FieldPassword struct {
	*fieldAtomic
	Algorithm PasswordAlgorithm
}

func NewFieldPassword(name string, algo PasswordAlgorithm) *FieldPassword {
	return &FieldPassword{fieldAtomic: &fieldAtomic{name: name}, Algorithm: algo}
}

type PasswordAlgorithm string

const (
	PasswordBcrypt PasswordAlgorithm = "Bcrypt"
	PasswordArgon2 PasswordAlgorithm = "Argon2"
	PasswordPbkdf2 PasswordAlgorithm = "Pbkdf2"
	PasswordScrypt PasswordAlgorithm = "Scrypt"
)

type FieldEmail struct {
	*fieldAtomic
}

func NewFieldEmail(name string) *FieldEmail {
	return &FieldEmail{fieldAtomic: &fieldAtomic{name: name}}
}

type FieldNode struct {
	*fieldAtomic
	Node string
}

func NewFieldNode(name string, node string) *FieldNode {
	return &FieldNode{fieldAtomic: &fieldAtomic{name: name}, Node: node}
}

type FieldEdge struct {
	*fieldAtomic
	Edge string
}

func NewFieldEdge(name string, edge string) *FieldEdge {
	return &FieldEdge{fieldAtomic: &fieldAtomic{name: name}, Edge: edge}
}

type FieldEnum struct {
	*fieldAtomic
	Typ string
}

func NewFieldEnum(name string, typ string) *FieldEnum {
	return &FieldEnum{fieldAtomic: &fieldAtomic{name: name}, Typ: typ}
}

type FieldStruct struct {
	*fieldAtomic
	Struct string
}

func NewFieldStruct(name string, structName string) *FieldStruct {
	return &FieldStruct{fieldAtomic: &fieldAtomic{name: name}, Struct: structName}
}

type FieldSlice struct {
	*fieldAtomic
	Field Field
}

func NewFieldSlice(name string, inner Field) *FieldSlice {
	return &FieldSlice{fieldAtomic: &fieldAtomic{name: name}, Field: inner}
}

func (f *FieldSlice) Validate() error {
	if err := f.Field.Validate(); err != nil {
		return err
	}

	if f.search != nil {
		if _, ok := f.Field.(*FieldString); !ok {
			return fmt.Errorf("field %s: fulltext index only supports string slice types, got slice of %T", f.name, f.Field)
		}
	}
	return nil
}

type FieldVersion struct {
	*fieldAtomic
}

type ComplexIDField struct {
	Name   string
	DBName string
	Field  Field
}

type FieldComplexID struct {
	*fieldAtomic
	Kind       IDType
	StructName string
	Fields     []ComplexIDField
}

func NewFieldComplexID(name string, kind IDType, structName string, fields []ComplexIDField) *FieldComplexID {
	return &FieldComplexID{
		fieldAtomic: &fieldAtomic{name: name},
		Kind:        kind,
		StructName:  structName,
		Fields:      fields,
	}
}

func (f *FieldComplexID) HasNodeRef() bool {
	for _, sf := range f.Fields {
		if _, ok := sf.Field.(*FieldNode); ok {
			return true
		}
	}
	return false
}
