package parser

import "github.com/wzshiming/gotype"

type FeatureSet struct {
	Timestamps     bool
	OptimisticLock bool
	SoftDelete     bool
	Expiry         bool
	ExpiryDuration string
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
				fieldBase:   newBase("CreatedAt"),
				IsCreatedAt: true,
				IsUpdatedAt: false,
			},
			&FieldTime{
				fieldBase:   newBase("UpdatedAt"),
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
				fieldBase:   fieldBase{name: "DeletedAt", pointer: true},
				IsDeletedAt: true,
			},
		)
		return true, nil

	case "Expiry":
		expiry, err := parseExpiryTag(f.Tag().Get("som"))
		if err != nil {
			return true, err
		}
		features.Expiry = true
		features.ExpiryDuration = expiry
		*fields = append(*fields,
			&FieldTime{
				fieldBase:   newBase("ExpiresAt"),
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
		*fields = append(*fields, NewFieldVersion("Version"))
	}
}
