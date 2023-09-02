// Code generated by github.com/marcbinz/som, DO NOT EDIT.
package relate

func NewGroup(db Database, unmarshal func(buf []byte, val any) error) *Group {
	return &Group{db: db, unmarshal: unmarshal}
}

type Group struct {
	db        Database
	unmarshal func(buf []byte, val any) error
}

func (n Group) Members() groupMember {
	return groupMember(n)
}
