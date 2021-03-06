package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

// Graph represents a set of connected nodes with a root.
type Graph struct {
	Root  int     `json:"root"`
	Nodes [][]int `json:"nodes"`
}

// NewGraph creates a new Graph using an io.Reader.
// It expects the reader to contains a json in the form :
// { root: 0, nodes: [[1 3],[2],[3],[4]] }
func NewGraph(r io.Reader) (*Graph, error) {
	g := &Graph{}
	d := json.NewDecoder(r)
	err := d.Decode(g)
	if err != nil {
		return nil, err
	}
	err = g.validate()
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (g *Graph) validate() error {
	if g.Root < 0 || g.Root >= len(g.Nodes) {
		return errors.New("root must point to an existing node")
	}
	for _, node := range g.Nodes {
		for _, edge := range node {
			if edge < 0 || edge >= len(g.Nodes) {
				return errors.New("edges must point to an existing node")
			}
		}
	}
	return nil
}

// BFS implements the Breadth-first search algorithm.
// See https://en.wikipedia.org/wiki/Breadth-first_search
func (g *Graph) BFS() map[int]int {
	dist := make(map[int]int)
	nodes := make([]int, 0, len(g.Nodes))

	dist[g.Root] = 0
	nodes = append(nodes, g.Root)

	// We go through nodes while adding node.
	for i := 0; i < len(nodes); i++ {
		node := nodes[i]
		for _, neighbor := range g.Nodes[node] {
			_, ok := dist[neighbor]
			if !ok {
				dist[neighbor] = dist[node] + 1 // Can add edge weight if using a weighted graph
				nodes = append(nodes, neighbor)
			}
		}
	}

	return dist
}

func fatalErr(err error, format string, a ...interface{}) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}

func main() {
	f, err := os.Open("graph.json")
	fatalErr(err, "fail to open graph.json: %v\n", err)

	graph, err := NewGraph(f)
	fatalErr(err, "fail to read graph.json: %v\n", err)
	fmt.Printf("Graph: %+v\n", graph)

	dist := graph.BFS()
	fmt.Println("Distances: ", dist)
}
