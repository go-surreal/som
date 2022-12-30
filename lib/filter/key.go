package filter

type KeyPart interface{}

type KeyField struct {
	name string
}

type KeyNode struct {
	name    string
	filters []Of[any]
}

type KeyEdge struct {
	name      string
	direction string
	filters   []Of[any]
}

type Key []KeyPart

func NewKey() Key {
	return Key{}
}

func (k Key) Field(name string) Key {
	return append(k, KeyField{
		name: name,
	})
}

func (k Key) Node(name string, filters []Of[any]) Key {
	return append(k, KeyNode{
		name:    name,
		filters: filters,
	})
}

func (k Key) EdgeIn(name string, filters []Of[any]) Key {
	return append(k, KeyEdge{
		name:      name,
		direction: "->",
		filters:   filters,
	})
}

func (k Key) EdgeOut(name string, filters []Of[any]) Key {
	return append(k, KeyEdge{
		name:      name,
		direction: "<-",
		filters:   filters,
	})
}
