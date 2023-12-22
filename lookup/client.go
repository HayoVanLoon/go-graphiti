// Package lookup implements a graph with a more exotic setup.
package lookup

import (
	"context"
	"errors"

	"github.com/HayoVanLoon/go-graphiti/basic"
)

// Client mimics a client to some external resource.
type Client struct {
	Nodes map[string]basic.Node
	Edges []basic.Edge
}

var ErrNotFound = errors.New("not found")

func (c Client) GetNode(_ context.Context, name string) (basic.Node, error) {
	n, ok := c.Nodes[name]
	if !ok {
		return basic.Node{}, ErrNotFound
	}
	return n, nil
}

func (c Client) GetEdges(_ context.Context, name string) ([]string, []int, error) {
	var ns []string
	var ds []int
	for _, e := range c.Edges {
		if e.From == name {
			if dst, ok := c.Nodes[e.To]; ok {
				ns = append(ns, dst.Name)
				ds = append(ds, e.Cost)
			}
		}
	}
	return ns, ds, nil
}
