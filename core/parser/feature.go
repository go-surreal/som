package parser

import "github.com/wzshiming/gotype"

type FeatureSet struct {
	Timestamps     bool
	OptimisticLock bool
	SoftDelete     bool
}

// ParseFeature checks if an anonymous field is a known feature embed.
// Returns true if the field was matched as a feature (caller should continue to next field).
func ParseFeature(f gotype.Type, internalPkg string, features *FeatureSet, fields *[]Field) bool {
	if f.Elem().PkgPath() != internalPkg {
		return false
	}

	switch f.Name() {
	case "Timestamps":
		features.Timestamps = true
		*fields = append(*fields,
			&FieldTime{
				fieldAtomic: &fieldAtomic{name: "CreatedAt"},
				IsCreatedAt: true,
				IsUpdatedAt: false,
			},
			&FieldTime{
				fieldAtomic: &fieldAtomic{name: "UpdatedAt"},
				IsCreatedAt: false,
				IsUpdatedAt: true,
			},
		)
		return true

	case "OptimisticLock":
		features.OptimisticLock = true
		return true

	case "SoftDelete":
		features.SoftDelete = true
		*fields = append(*fields,
			&FieldTime{
				fieldAtomic: &fieldAtomic{name: "DeletedAt", pointer: true},
				IsDeletedAt: true,
			},
		)
		return true
	}

	return false
}

// ApplyFeatures copies feature flags to the target booleans and appends
// the Version field if OptimisticLock is enabled.
func ApplyFeatures(features FeatureSet, timestamps, optimisticLock, softDelete *bool, fields *[]Field) {
	*timestamps = features.Timestamps
	*optimisticLock = features.OptimisticLock
	*softDelete = features.SoftDelete

	if features.OptimisticLock {
		*fields = append(*fields, &FieldVersion{&fieldAtomic{name: "Version", pointer: false}})
	}
}
