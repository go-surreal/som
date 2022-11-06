package depth

func keyed(base, key string) string {
	if base == "" {
		return key
	}
	return string(base) + "." + key
}

func User(key string, depth int) []string {
	if depth == 0 {
		return []string{key}
	}

	var fields []string

	for _, field := range Group("main_group", depth-1) {
		fields = append(fields, field)
	}

	return fields
}

func Group(key string, depth int) []string {
	if depth == 0 {
		return []string{key}
	}

	return []string{key}
}
