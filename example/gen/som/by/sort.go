// Code generated by github.com/marcbinz/som, DO NOT EDIT.

package by

func keyed(base, key string) string {
	if base == "" {
		return key
	}
	return base + "." + key
}