package lookup

import (
	"context"
	"errors"

	"github.com/HayoVanLoon/go-graphiti/basic"
)

type Database struct {
	Graph basic.Graph
}

// InitDatabase initialises the 'database'. It should be called before starting
// Graphiti graph building.
func InitDatabase(nodes map[string]basic.Node, edges []basic.Edge) {
	DB.Graph = basic.Graph{Nodes: nodes, Edges: edges}
}

var DB = new(Database)

type Client struct {
	db *Database
}

var ErrNotFound = errors.New("not found")

func (c Client) GetNode(_ context.Context, name string) (basic.Node, error) {
	n, ok := c.db.Graph.Nodes[name]
	if !ok {
		return basic.Node{}, ErrNotFound
	}
	return n, nil
}

func (c Client) GetEdges(_ context.Context, name string) ([]string, []int, error) {
	var ns []string
	var ds []int
	for _, e := range c.db.Graph.Edges {
		if e.From == name {
			if dst, ok := c.db.Graph.Nodes[e.To]; ok {
				ns = append(ns, dst.Name)
				ds = append(ds, e.Cost)
			}
		}
	}
	return ns, ds, nil
}

func NewClient() Client {
	return Client{DB}
}
