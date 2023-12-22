package graphiti_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	g "github.com/HayoVanLoon/go-graphiti"
	"github.com/HayoVanLoon/go-graphiti/basic"
	"github.com/HayoVanLoon/go-graphiti/lookup"
)

func TestGraphiti_New(t *testing.T) {
	type args struct {
		graph basic.Graph
	}
	tests := []struct {
		name string
		args args
		want *g.Graph[string]
	}{
		{
			"simple DG",
			args{simpleDG()},
			&g.Graph[string]{
				Nodes:   []string{"A", "B", "C", "D"},
				IndexOf: map[string]int{"A": 0, "B": 1, "C": 2, "D": 3},
				Edges:   edgesDG(),
				GroupOf: []int{1, 1, 1, 1},
				M:       m4x4(),
			},
		},
		{
			"empty",
			args{basic.Graph{}},
			&g.Graph[string]{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gr := newGraphiti()
			actual, err := gr.New(tt.args.graph, nil)
			require.NoError(t, err)
			require.Equal(t, tt.want, actual)
		})
	}
}

func TestGraphiti_AllDistances(t *testing.T) {
	type args struct {
		g basic.Graph
	}
	tests := []struct {
		name string
		args args
		want *g.Graph[string]
	}{
		{
			"simple DAG",
			args{simpleDAG()},
			&g.Graph[string]{
				Nodes:   []string{"A", "B", "C", "D"},
				IndexOf: map[string]int{"A": 0, "B": 1, "C": 2, "D": 3},
				Edges:   edgesDAG(),
				GroupOf: []int{1, 1, 1, 1},
				M: [][]g.R{
					{r2(0), r2(1), r2(1), r2(2)},
					{r2(0), r2(0), r2(1), r2(2)},
					{r2(0), r2(0), r2(0), r2(1)},
					{r2(0), r2(0), r2(0), r2(0)},
				},
			},
		},
		{
			"simple DG",
			args{simpleDG()},
			&g.Graph[string]{
				Nodes:   []string{"A", "B", "C", "D"},
				IndexOf: map[string]int{"A": 0, "B": 1, "C": 2, "D": 3},
				Edges:   edgesDG(),
				GroupOf: []int{1, 1, 1, 1},
				M: [][]g.R{
					{r2(3), r2(1), r2(1), r2(2)},
					{r2(3), r2(4), r2(1), r2(2)},
					{r2(2), r2(3), r2(3), r2(1)},
					{r2(1), r2(2), r2(2), r2(3)},
				},
			},
		},
		{
			"simple disconnected DAG",
			args{simpleDisconnectedDAG()},
			&g.Graph[string]{
				Nodes:   []string{"A", "B", "C", "D"},
				IndexOf: map[string]int{"A": 0, "B": 1, "C": 2, "D": 3},
				Edges:   edgesDisconnectedDAG(),
				GroupOf: []int{1, 1, 2, 2},
				M: [][]g.R{
					{r2(0), r2(1), r2(0), r2(0)},
					{r2(0), r2(0), r2(0), r2(0)},
					{r2(0), r2(0), r2(0), r2(1)},
					{r2(0), r2(0), r2(0), r2(0)},
				},
			},
		},
		{
			"empty",
			args{basic.Graph{}},
			&g.Graph[string]{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gr := newGraphiti()
			actual, err := gr.New(tt.args.g, nil)
			require.NoError(t, err)
			_ = actual.AllDistances()
			require.NoError(t, err)
			require.Equalf(t, tt.want, actual, "expected:\n%s\ngot:\n%s", tt.want, actual)
		})
	}
}

func ExampleGraph_AllDistances() {
	initLookupExample()
	gr := g.Graphiti[struct{}, string, string]{
		GetName:  lookup.GetName,
		GetNode:  lookup.GetNode,
		GetNext:  lookup.GetNext,
		GetNodes: lookup.StartWith("A"),
	}
	graph, _ := gr.New(struct{}{}, nil)
	_ = graph.AllDistances()
	fmt.Println(graph)
	// Output:
	// -	1/1	1/1	2/2	A (1)
	// -	-	1/1	2/2	B (1)
	// -	-	-	1/1	C (1)
	// -	-	-	-	D (1)
}

func ExampleGraph_Groups() {
	gr := g.Graphiti[basic.Graph, basic.Node, string]{
		GetName:  basic.GetName,
		GetNode:  basic.GetNode,
		GetNext:  basic.GetNext,
		GetNodes: basic.GetNodes,
	}
	graph, _ := gr.New(simpleDisconnectedDAG(), nil)
	fmt.Println(graph.Groups())
	// Output:
	// [[A B] [C D]]
}

func TestGraphiti_Steps(t *testing.T) {
	type args struct {
		g    basic.Graph
		from string
		to   string
	}
	type want struct {
		value int
		err   require.ErrorAssertionFunc
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"simple DAG",
			args{simpleDAG(), "C", "B"},
			want{0, require.NoError},
		},
		{
			"simple DG",
			args{simpleDG(), "C", "B"},
			want{3, require.NoError},
		},
		{
			"empty",
			args{basic.Graph{}, "C", "B"},
			want{0, require.Error},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gr := newGraphiti()
			graph, _ := gr.New(tt.args.g, nil)
			actual, err := graph.Steps(tt.args.from, tt.args.to)
			tt.want.err(t, err)
			require.Equal(t, tt.want.value, actual)
		})
	}
}

func TestM_HasCycles(t *testing.T) {
	type args struct {
		g basic.Graph
	}
	type want struct {
		value bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"simple DAG",
			args{simpleDAG()},
			want{false},
		},
		{
			"simple DG",
			args{simpleDG()},
			want{true},
		},
		{
			"empty",
			args{basic.Graph{}},
			want{false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gr := newGraphiti()
			graph, _ := gr.New(tt.args.g, nil)
			_ = graph.AllDistances()
			require.Equal(t, tt.want.value, graph.HasCycles())
		})
	}
}

func TestGraph_Groups(t *testing.T) {
	type args struct {
		g basic.Graph
	}
	type want struct {
		value [][]string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"simple DAG",
			args{simpleDAG()},
			want{[][]string{{"A", "B", "C", "D"}}},
		},
		{
			"simple disconnected DAG",
			args{simpleDisconnectedDAG()},
			want{[][]string{{"A", "B"}, {"C", "D"}}},
		},
		{
			"empty",
			args{basic.Graph{}},
			want{nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gr := newGraphiti()
			graph, _ := gr.New(tt.args.g, nil)
			require.Equal(t, tt.want.value, graph.Groups())
		})
	}
}

func newGraphiti() g.Graphiti[basic.Graph, basic.Node, string] {
	return g.Graphiti[basic.Graph, basic.Node, string]{
		GetName:  basic.GetName,
		GetNode:  basic.GetNode,
		GetNext:  basic.GetNext,
		GetNodes: basic.GetNodes,
	}
}

func m4x4() [][]g.R {
	return [][]g.R{
		make([]g.R, 4),
		make([]g.R, 4),
		make([]g.R, 4),
		make([]g.R, 4),
	}
}

func r2(i int) g.R {
	return g.R{Dist: i, Steps: i}
}

func simpleDAG() basic.Graph {
	return basic.Graph{
		Nodes: map[string]basic.Node{
			"A": {Name: "A"},
			"B": {Name: "B"},
			"C": {Name: "C"},
			"D": {Name: "D"},
		},
		Edges: []basic.Edge{
			{From: "A", To: "B", Cost: 1},
			{From: "A", To: "C", Cost: 1},
			{From: "B", To: "C", Cost: 1},
			{From: "C", To: "D", Cost: 1},
		},
	}
}

func edgesDAG() [][3]int {
	return [][3]int{
		{0, 1, 1},
		{0, 2, 1},
		{1, 2, 1},
		{2, 3, 1},
	}
}

func simpleDG() basic.Graph {
	return basic.Graph{
		Nodes: map[string]basic.Node{
			"A": {Name: "A"},
			"B": {Name: "B"},
			"C": {Name: "C"},
			"D": {Name: "D"},
		},
		Edges: []basic.Edge{
			{From: "A", To: "B", Cost: 1},
			{From: "A", To: "C", Cost: 1},
			{From: "B", To: "C", Cost: 1},
			{From: "C", To: "D", Cost: 1},
			{From: "D", To: "A", Cost: 1},
		},
	}
}

func edgesDG() [][3]int {
	return [][3]int{
		{0, 1, 1},
		{0, 2, 1},
		{1, 2, 1},
		{2, 3, 1},
		{3, 0, 1},
	}
}

func simpleDisconnectedDAG() basic.Graph {
	return basic.Graph{
		Nodes: map[string]basic.Node{
			"A": {Name: "A"},
			"B": {Name: "B"},
			"C": {Name: "C"},
			"D": {Name: "D"},
		},
		Edges: []basic.Edge{
			{From: "A", To: "B", Cost: 1},
			{From: "C", To: "D", Cost: 1},
		},
	}
}

func edgesDisconnectedDAG() [][3]int {
	return [][3]int{
		{0, 1, 1},
		{2, 3, 1},
	}
}

func initLookupExample() {
	data := simpleDAG()
	lookup.InitDatabase(data.Nodes, data.Edges)
}
