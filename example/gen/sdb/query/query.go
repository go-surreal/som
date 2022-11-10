package query

type Database interface {
	Query(statement string, vars map[string]any) (any, error)
}
	
type idNode struct {
	ID string
}
	
type countResult struct {
	Count int	
}
