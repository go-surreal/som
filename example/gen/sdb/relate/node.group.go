package relate

func NewGroup(db Database) *Group {
	return &Group{db: db}
}

type Group struct {
	db Database
}

func (n Group) Members() memberOf {
	return memberOf{db: n.db}
}
