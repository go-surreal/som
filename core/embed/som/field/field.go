//go:build embed

package field

func keyed(base, key string) string {
	if base == "" {
		return key
	}
	return base + "." + key
}
