//go:build embed

package repo

func newID(table string) string {
	return table + ":rand()"
}

func newULID(table string) string {
	return table + ":ulid()"
}

func newUUID(table string) string {
	return table + ":uuid()"
}
