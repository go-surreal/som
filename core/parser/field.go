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

type FieldString struct {
	*fieldAtomic
}

func (f *FieldString) Validate() error {
	return nil // string and *string support fulltext
}

type FieldNumeric struct {
	*fieldAtomic
	Type NumberType
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

type FieldByte struct {
	*fieldAtomic
}

type FieldDuration struct {
	*fieldAtomic
}

type FieldTime struct {
	*fieldAtomic
	IsCreatedAt bool
	IsUpdatedAt bool
	IsDeletedAt bool
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

func NewFieldUUID(name string, pointer bool, pkg UUIDPackage) *FieldUUID {
	return &FieldUUID{
		fieldAtomic: &fieldAtomic{
			name:    name,
			pointer: pointer,
		},
		Package: pkg,
	}
}

type FieldURL struct {
	*fieldAtomic
}

type FieldPassword struct {
	*fieldAtomic
	Algorithm PasswordAlgorithm
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

type FieldNode struct {
	*fieldAtomic
	Node string
}

type FieldEdge struct {
	*fieldAtomic
	Edge string
}

type FieldEnum struct {
	*fieldAtomic
	Typ string
}

type FieldStruct struct {
	*fieldAtomic
	Struct string
}

type FieldSlice struct {
	*fieldAtomic
	// Value  string
	Field Field
	// IsNode bool
	// IsEdge bool
	// IsEnum bool
}

func (f *FieldSlice) Validate() error {
	// First validate the element
	if err := f.Field.Validate(); err != nil {
		return err
	}

	// Fulltext is only valid for string slices
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

func (f *FieldComplexID) HasNodeRef() bool {
	for _, sf := range f.Fields {
		if _, ok := sf.Field.(*FieldNode); ok {
			return true
		}
	}
	return false
}

// type FieldMap struct {
// 	fieldAtomic
// 	Key   string
// 	Value string
// }
