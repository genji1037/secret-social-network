package memgraph

// Node represent node.
type Node struct {
	Name      string
	Consensus []*Node
	Values    []float64
}
