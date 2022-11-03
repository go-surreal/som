package by

func keyed(base, key string) string {
	if base == "" {
		return key
	}
	return base + "." + key
}
