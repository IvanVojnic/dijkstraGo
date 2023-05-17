package main

import (
	"container/heap"
	"context"
	"fmt"
	"lab3/models"
	"lab3/repository"
	"os"
)

type PriorityQueue []models.Сrossroad

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].Dist < pq[j].Dist }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }

func (pq *PriorityQueue) Push(x interface{}) {
	node := x.(models.Сrossroad)
	*pq = append(*pq, node)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	*pq = old[0 : n-1]
	return node
}

func dijkstra(graph []*models.Сrossroad, start int, end int) int {
	n := len(graph)
	dist := make([]int, n)
	visited := make([]bool, n)

	for i := 0; i < n; i++ {
		dist[i] = 1e9
	}

	dist[start] = 0

	pq := make(PriorityQueue, 0)
	heap.Push(&pq, models.Сrossroad{CrossroadID: start, Dist: 0})

	for pq.Len() > 0 {
		node := heap.Pop(&pq).(models.Сrossroad)
		u := node.CrossroadID

		if u == end {
			return node.Dist
		}

		if visited[u] {
			continue
		}

		visited[u] = true

		for _, road := range graph[u].CrossroadRoads {
			v := road.EndCrossroad
			w := road.RoadTime
			if dist[u]+w < dist[v.CrossroadID] {
				dist[v.CrossroadID] = dist[u] + w
				heap.Push(&pq, models.Сrossroad{CrossroadID: v.CrossroadID, Dist: dist[v.CrossroadID]})
			}
		}
	}

	return -1
}

func main() {
	fileRead, err := os.Open("input.txt")
	if err != nil {
		fmt.Errorf("error while open file to read - %v", err)
	}
	defer fileRead.Close()
	fileWrite, err := os.Open("output.txt")
	if err != nil {
		fmt.Errorf("error while open file to read - %v", err)
	}
	defer fileWrite.Close()
	repos := repository.NewFileRepo(fileRead, fileWrite)
	graph, condition, err := repos.GetCrossroads(context.Background())
	if err != nil {
		fmt.Errorf("error while get crossroads - %v", err)
	}

	fmt.Printf("condition - %v, nodes - %v", condition, graph)
	fmt.Print("\n")
	// Print graph
	for _, vertex := range graph {
		fmt.Printf("Перекресток: %v\n", vertex.CrossroadID)
		for _, edge := range vertex.CrossroadRoads {
			fmt.Printf("  Дорога к перекрестоку %v, время прохождения: %d минут\n", edge.EndCrossroad.CrossroadID, edge.RoadTime+vertex.Time)
		}
	}
	result := dijkstra(graph, condition.S-1, condition.F-1)
	if result == -1 {
		fmt.Printf("No")
	} else {
		fmt.Printf("Yes - %v", 12)
	}
}
