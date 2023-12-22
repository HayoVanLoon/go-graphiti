package graphiti

import (
	"strings"
	"fmt"
	"io"
	"errors"
)

type ErrNotFound[ID comparable] struct {
	Node ID
}

func (e *ErrNotFound[ID]) Error() string {
	return fmt.Sprintf("node not found: %v", e.Node)
}

func notFound[ID comparable](name ID) error {
	return &ErrNotFound[ID]{name}
}

type ErrBrokenEdge[ID comparable] struct {
	From, To ID
}

func (e *ErrBrokenEdge[ID]) Error() string {
	return fmt.Sprintf("broken edge: %v->%v", e.From, e.To)
}

func brokenEdge[ID comparable](from, to ID) error {
	return &ErrBrokenEdge[ID]{from, to}
}

type R struct {
	Steps int
	Dist  int
}

type Graph[ID comparable] struct {
	Nodes   []ID
	IndexOf map[ID]int
	Edges   [][3]int
	GroupOf []int
	M       [][]R
}

const unset = 0

func (g *Graph[ID]) Steps(from, to ID) (int, error) {
	x := g.getIdx(from)
	if x < 0 {
		return 0, notFound(from)
	}
	y := g.getIdx(to)
	if y < 0 {
		return 0, notFound(to)
	}

	if err := g.distance(x); err != nil {
		return unset, err
	}
	return g.M[x][y].Steps, nil
}

func (g *Graph[ID]) distance(from int) error {
	done := make([]bool, len(g.Nodes))
	q := g.getEdges(from)
	dist := 0
	for steps := 1; len(q) > 0; steps += 1 {
		var nextQ [][2]int
		var nextDist int
		for _, e := range q {
			c := dist + e[1]
			if shorter := g.set(from, e[0], c, steps); shorter {
				nextDist = c
			}
			for _, e2 := range g.getEdges(e[0]) {
				if !done[e2[0]] {
					done[e2[0]] = true
					nextQ = append(nextQ, e2)
				}
			}
		}
		q = nextQ
		dist = nextDist
	}
	return nil
}

func (g *Graph[ID]) AllDistances() error {
	for i := range g.Nodes {
		if err := g.distance(i); err != nil {
			return err
		}
	}
	return nil
}

func (g *Graph[ID]) String() string {
	addInt := func(w io.Writer, i int) {
		if i == unset {
			_, _ = fmt.Fprintf(w, "-")
		} else {
			_, _ = fmt.Fprintf(w, "%d", i)
		}
	}

	b := new(strings.Builder)
	for x := range g.M {
		for y := range g.M[x] {
			r := g.M[x][y]
			if r.Dist == unset && r.Steps == unset {
				b.WriteString("-\t")
				continue
			}
			addInt(b, r.Dist)
			b.WriteRune('/')
			addInt(b, r.Steps)
			b.WriteRune('\t')
		}
		_, _ = fmt.Fprintf(b, "%v (%d)\n", g.Nodes[x], g.GroupOf[x])
	}
	return b.String()
}

func (g *Graph[ID]) set(x, y int, dist, steps int) bool {
	r := g.M[x][y]
	shorter := r.Dist == unset || r.Dist > dist
	if shorter {
		g.M[x][y].Dist = dist
	}
	if r.Steps == 0 || r.Steps > steps {
		g.M[x][y].Steps = steps
	}
	return shorter
}

func (g *Graph[ID]) updateGroups(x, y int) {
	next := 0
	for _, i := range g.GroupOf {
		next = max(next, i)
	}
	next += 1

	if g.GroupOf[x] == unset {
		g.GroupOf[x] = next
	}
	if g.GroupOf[y] == unset {
		g.GroupOf[y] = next
	}
	if g.GroupOf[x] == g.GroupOf[y] {
		return
	}
	mn := min(g.GroupOf[x], g.GroupOf[y])
	mx := max(g.GroupOf[x], g.GroupOf[y])
	for i := range g.GroupOf {
		if g.GroupOf[i] == mx {
			g.GroupOf[i] = mn
		}
	}
}

func (g *Graph[ID]) HasCycles() bool {
	for xy := range g.M {
		if g.M[xy][xy] != (R{}) {
			return true
		}
	}
	return false
}

func (g *Graph[ID]) Groups() [][]ID {
	var gs [][]ID
	for idx, gr := range g.GroupOf {
		for i := len(gs); i < gr; i += 1 {
			gs = append(gs, []ID{})
		}
		gs[gr-1] = append(gs[gr-1], g.Nodes[idx])
	}
	return gs
}

func (g *Graph[ID]) getIdx(id ID) int {
	for i, id2 := range g.Nodes {
		if id == id2 {
			return i
		}
	}
	return -1
}

func (g *Graph[ID]) getEdges(id int) [][2]int {
	var q [][2]int
	for _, e := range g.Edges {
		if e[0] == id {
			q = append(q, [2]int{e[1], e[2]})
		}
	}
	return q
}

type Graphiti[G any, N any, ID comparable] struct {
	GetName  func(N) ID
	GetNode  func(g G, name ID) (N, bool)
	GetNext  func(g G, node N) ([]ID, []int)
	GetNodes func(G) []ID
}

func (gr Graphiti[G, N, ID]) New(graph G, nodes []N) (*Graph[ID], error) {
	g, errs := gr.newGraph(graph, nodes, false)
	return g, errors.Join(errs...)
}

func (gr Graphiti[G, N, ID]) Validate(graph G, nodes []N) (*Graph[ID], []error) {
	g, errs := gr.newGraph(graph, nodes, true)
	return g, errs
}

func (gr Graphiti[G, N, ID]) newGraph(graph G, nodes []N, persevere bool) (*Graph[ID], []error) {
	var errs []error
	g := &Graph[ID]{}
	if len(nodes) == 0 {
		g.Nodes = gr.GetNodes(graph)
	} else {
		for _, n := range nodes {
			g.Nodes = append(g.Nodes, gr.GetName(n))
		}
	}
	g.IndexOf = make(map[ID]int, len(g.Nodes))
	for i, id := range g.Nodes {
		g.IndexOf[id] = i
	}

	for i := 0; i < len(g.Nodes); i += 1 {
		g.M = append(g.M, make([]R, len(g.Nodes)))
	}
	g.GroupOf = make([]int, len(g.Nodes))

	for idx := 0; idx < len(g.Nodes); idx += 1 {
		ok, errs2 := gr.addEdges(graph, g.Nodes[idx], g, idx, persevere)
		errs = append(errs, errs2...)
		if !ok {
			return nil, errs
		}
	}

	if len(g.Nodes) == 0 {
		return &Graph[ID]{}, nil
	}
	return g, nil
}

func (gr Graphiti[G, N, ID]) addEdges(graph G, id ID, g *Graph[ID], idx int, persevere bool) (bool, []error) {
	var errs []error

	n, ok := gr.GetNode(graph, id)
	if !ok {
		errs = append(errs, notFound(id))
		if persevere {
			return true, errs
		}
		return false, errs
	}

	es, ds := gr.GetNext(graph, n)
	if len(es) != len(ds) {
		err := fmt.Errorf("GetNext slice lengths differ: %d != %d", len(es), len(ds))
		return false, []error{err}
	}
	for j := range es {
		if ds[j] < 0 {
			err := fmt.Errorf("edge with negative cost: (%v,%v)=%d", id, es[j], ds[j])
			errs = append(errs, err)
			if persevere {
				continue
			}
			return false, errs
		}
		idx2, ok := g.IndexOf[es[j]]
		if !ok {
			xn, ok := gr.GetNode(graph, es[j])
			if !ok {
				errs = append(errs, brokenEdge(id, es[j]))
				if persevere {
					continue
				}
				return false, errs
			}
			idx2 = gr.appendNewNode(g, xn)
		}
		g.updateGroups(idx, idx2)
		g.Edges = append(g.Edges, [3]int{idx, idx2, ds[j]})
	}
	return true, errs
}

func (gr Graphiti[G, N, ID]) appendNewNode(g *Graph[ID], newNode N) int {
	name := gr.GetName(newNode)

	g.Nodes = append(g.Nodes, name)
	g.GroupOf = append(g.GroupOf, 0)
	idx := len(g.Nodes) - 1
	g.IndexOf[name] = idx

	for x := 0; x < len(g.M); x += 1 {
		g.M[x] = append(g.M[x], R{})
	}
	g.M = append(g.M, make([]R, len(g.Nodes)))

	return idx
}
