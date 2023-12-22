// Package lookup provides an example implementation of Graphiti functions.
//
// This implementation does not specify graph type. All graph information is
// retrieved from some (mock) upstream service.
package lookup

import (
	"context"
	"log"
)

// GetName implements the Graphiti function of the same name. Trivial, as the
// node type is already string.
func GetName(name string) string {
	return name
}

// GetNode implement the Graphiti function of the same name. It looks up a node
// from an 'external' source. There is no graph here, so the first argument can
// be ignored.
func GetNode(_ struct{}, name string) (string, bool) {
	n, err := FetchNode(context.Background(), name)
	if err != nil {
		log.Println("error", err)
		return "", false
	}
	return n, true
}

// GetNext implement the Graphiti function of the same name. It looks up edge
// information from an 'external' source. There is no graph here, so the first
// argument is ignored.
func GetNext(_ struct{}, n string) ([]string, []int) {
	ns, ds, err := FetchEdges(context.Background(), n)
	if err != nil {
		log.Println("error", err)
	}
	return ns, ds
}

// StartWith returns a function returning the starting node.
func StartWith(start string) func(struct{}) []string {
	return func(struct{}) []string {
		return []string{start}
	}
}

func FetchNode(ctx context.Context, name string) (string, error) {
	client := NewClient()

	n, err := client.GetNode(ctx, name)
	if err != nil {
		return "", err
	}
	return n.Name, nil
}

func FetchEdges(ctx context.Context, node string) ([]string, []int, error) {
	client := NewClient()

	ns, ds, err := client.GetEdges(ctx, node)
	if err != nil {
		return nil, nil, err
	}
	return ns, ds, nil
}
