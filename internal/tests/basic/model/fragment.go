package model

import (
	"time"

	"som.test/gen/som"
)

// AllTypesCore projects a handful of plain fields of AllTypes.
type AllTypesCore struct {
	som.Fragment[AllTypes]

	FieldString string
	FieldInt    int
	FieldBool   bool
	FieldTime   time.Time
	FieldEnum   Role
}

// AllTypesRelations projects fields of AllTypes that need conversion:
// a renamed field, a struct, a record link and a slice of links.
type AllTypesRelations struct {
	som.Fragment[AllTypes]

	FieldRenamed     string
	FieldCredentials Credentials
	FieldNodePtr     *SpecialTypes
	FieldNodeSlice   []SpecialTypes
}

// SpecialTypesName projects only the name of a SpecialTypes record, whose node
// has soft delete enabled.
type SpecialTypesName struct {
	som.Fragment[SpecialTypes]

	Name string
}

// WeatherReading projects a field of a node with a complex (array) id.
type WeatherReading struct {
	som.Fragment[Weather]

	Temperature float64
}
