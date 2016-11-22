package main

import "fmt"

func main() {

	graph, root := constructGraph()
	fmt.Println("Graph: ", graph)

	dist := bfs(graph, root)

	fmt.Println("Distances: ", dist)
}

func constructGraph() ([][]int, int) {
	var nodes, edges int
	fmt.Printf("Nb nodes: ")
	fmt.Scanln(&nodes)

	graph := make([][]int, nodes)
	for i := 0; i < nodes; i++ {
		graph[i] = make([]int, 0)
	}

	fmt.Printf("Nb edges: ")
	fmt.Scanln(&edges)

	for i := 0; i < edges; i++ {
		var node1, node2 int
		fmt.Print("Edge ", i+1, ": ")
		fmt.Scanf("%d %d\n", &node1, &node2)
		graph[node1] = append(graph[node1], node2)
		graph[node2] = append(graph[node2], node1)
	}

	var root int
	fmt.Printf("Root node: ")
	fmt.Scanln(&root)

	return graph, root
}

func bfs(graph [][]int, root int) map[int]int {
	distance := make(map[int]int)

	done := make(chan bool)
	queue := make(chan int, len(graph))
	distance[root] = 0
	go func() {
		for {
			select {
			case node := <-queue:
				for _, neighbor := range graph[node] {
					_, prs := distance[neighbor]
					if !prs {
						distance[neighbor] = distance[node] + 1 // Can add edge weight if using a weighted graph
						queue <- neighbor
					}
				}
			default:
				done <- true
			}
		}
	}()

	queue <- root

	<-done

	return distance
}
