package basic

import "sort"

type Node struct {
	Name string
}

type Edge struct {
	From string
	To   string
	Cost int
}

type Graph struct {
	Nodes map[string]Node
	Edges []Edge
}

func GetName(n Node) string {
	return n.Name
}

func GetNode(g Graph, name string) (Node, bool) {
	n, ok := g.Nodes[name]
	return n, ok
}

func GetNext(g Graph, n Node) ([]string, []int) {
	var ns []string
	var ds []int
	for _, e := range g.Edges {
		if e.From == n.Name {
			dst, ok := GetNode(g, e.To)
			if ok {
				ns = append(ns, dst.Name)
				ds = append(ds, e.Cost)
			}
		}
	}
	return ns, ds
}

func GetNodes(g Graph) []string {
	var ns []string
	for _, v := range g.Nodes {
		ns = append(ns, v.Name)
	}
	sort.Strings(ns)
	return ns
}
