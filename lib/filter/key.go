package filter

type Key struct {
	key   string
	open  bool
	close int
}

func NewKey() Key {
	return Key{}
}

func (k Key) Dot(field string) Key {
	if k.key == "" {
		k.key = field
	} else if k.open {
		k.key += " WHERE " + field
		k.close += 1
		k.open = false
	} else {
		k.key += "." + field
	}
	return k
}

func (k Key) In(elem string) Key {
	if k.open {
		k.key += ")"
		k.open = false
	}
	k.key += "->(" + elem
	k.open = true
	return k
}

func (k Key) Out(elem string) Key {
	if k.open {
		k.key += ")"
		k.open = false
	}
	k.key += "<-(" + elem
	k.open = true
	return k
}

func (k Key) Sub(field string) Key {
	k.key += " WHERE " + field
	return k
}
