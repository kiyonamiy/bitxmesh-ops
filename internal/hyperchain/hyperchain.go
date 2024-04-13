package hyperchain

type Node struct {
	host string
	port int
}

func NewNode(host string, port int) *Node {
	return &Node{
		host,
		port,
	}
}

func (n *Node) Host() string {
	return n.host
}

func (n *Node) Port() int {
	return n.port
}
