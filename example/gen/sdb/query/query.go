package query

type Database interface {
	Query(statement string, vars map[string]any) (any, error)
}
