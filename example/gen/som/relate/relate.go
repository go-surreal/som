package relate

type Database interface {
	Query(statement string, vars any) (any, error)
}
