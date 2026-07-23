package parser

import "github.com/wzshiming/gotype"

type FeatureSet struct {
	Timestamps     bool
	OptimisticLock bool
	SoftDelete     bool
	TTL            bool
	TTLExpiry      string
}

// ParseFeature checks if an anonymous field is a known feature embed.
// Returns true if the field was matched as a feature (caller should continue to next field).
func ParseFeature(f gotype.Type, internalPkg string, features *FeatureSet, fields *[]Field) (bool, error) {
	if f.Elem().PkgPath() != internalPkg {
		return false, nil
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
		return true, nil

	case "OptimisticLock":
		features.OptimisticLock = true
		return true, nil

	case "SoftDelete":
		features.SoftDelete = true
		*fields = append(*fields,
			&FieldTime{
				fieldAtomic: &fieldAtomic{name: "DeletedAt", pointer: true},
				IsDeletedAt: true,
			},
		)
		return true, nil

	case "Expiry":
		expiry, err := parseTTLTag(f.Tag().Get("som"))
		if err != nil {
			return true, err
		}
		features.TTL = true
		features.TTLExpiry = expiry
		*fields = append(*fields,
			&FieldTime{
				fieldAtomic: &fieldAtomic{name: "ExpiresAt"},
				IsExpiresAt: true,
				ExpiresIn:   expiry,
			},
		)
		return true, nil
	}

	return false, nil
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
