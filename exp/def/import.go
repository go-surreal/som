package def

type Import struct {
	Name string
	Path string
}

func (i *Import) String() string {
	return i.Name + ": " + i.Path
}
