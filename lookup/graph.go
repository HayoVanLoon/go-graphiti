// Package lookup implements Graphiti functions using lookups.
package lookup

import (
	"context"
	"log"
)

type NotGraph struct {
	// Start is the identifier of the starting 'node'.
	Start  string
	Client Client
}

func (ng NotGraph) LookupNode(ctx context.Context, name string) (string, error) {
	n, err := ng.Client.GetNode(ctx, name)
	if err != nil {
		return "", err
	}
	return n.Name, nil
}

func (ng NotGraph) LookupEdges(ctx context.Context, n string) ([]string, []int, error) {
	ns, ds, err := ng.Client.GetEdges(ctx, n)
	if err != nil {
		return nil, nil, err
	}
	return ns, ds, nil
}

func (NotGraph) GetName(name string) string {
	return name
}

func (ng NotGraph) GetNode(_ struct{}, name string) (string, bool) {
	n, err := ng.LookupNode(context.Background(), name)
	if err != nil {
		log.Println("error", err)
		return "", false
	}
	return n, true
}

func (ng NotGraph) GetNext(_ struct{}, n string) ([]string, []int) {
	ns, ds, err := ng.LookupEdges(context.Background(), n)
	if err != nil {
		log.Println("error", err)
	}
	return ns, ds
}

func (ng NotGraph) GetNodes(struct{}) []string {
	return []string{ng.Start}
}
