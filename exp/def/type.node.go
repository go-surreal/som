package def

type Node struct {
	*Struct
}

func (n *Node) String() string {
	return n.describe("Node")
}
