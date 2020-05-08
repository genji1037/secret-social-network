package memgraph

type Node struct {
	Name      string
	Consensus []*Node
	Values    []float64
}
