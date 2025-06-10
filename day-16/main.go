package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"os"
)

const (
	North = 0
	East  = 1
	South = 2
	West  = 3
)

type State struct {
	x, y int
	dir  int
}

type QueueItem struct {
	state State
	cost  int
}

type PriorityQueue []*QueueItem

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].cost < pq[j].cost }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x any)        { *pq = append(*pq, x.(*QueueItem)) }
func (pq *PriorityQueue) Pop() any          { item := (*pq)[len(*pq)-1]; *pq = (*pq)[:len(*pq)-1]; return item }

// Directions vectors: North, East, South, West
var directions = [][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

func main() {
	grid, start := parseGrid("example-input.txt")
	fmt.Printf("Start: (%d,%d) facing %s\n", start.x, start.y, directionName(start.dir))

	// INFO: Part 1: Find minimum cost
	_, minCost := dijkstraForward(grid, start)
	fmt.Printf("Part 1: %d\n", minCost)

	// INFO: Part 2: Count optimal path time
	optimalTiles := countOptimalTiles(grid, start)
	fmt.Printf("Part 2: %d\n", optimalTiles)
}

func parseGrid(filename string) ([][]rune, State) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	var grid [][]rune
	var start State
	scanner := bufio.NewScanner(file)

	for row := 0; scanner.Scan(); row++ {
		line := scanner.Text()
		grid = append(grid, []rune(line))

		for col, cell := range line {
			if cell == 'S' {
				start = State{row, col, East}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	return grid, start
}

func dijkstraForward(grid [][]rune, start State) (map[State]int, int) {
	pq := &PriorityQueue{}
	heap.Push(pq, &QueueItem{start, 0})
	distances := make(map[State]int)

	for pq.Len() > 0 {
		current := heap.Pop(pq).(*QueueItem)
		state, cost := current.state, current.cost

		if prevCost, exists := distances[state]; exists && cost > prevCost {
			continue
		}
		distances[state] = cost

		if grid[state.x][state.y] == 'E' {
			return distances, cost
		}

		// Go Forward
		dx, dy := directions[state.dir][0], directions[state.dir][1]
		newX := state.x + dx
		newY := state.y + dy

		if isValid(grid, newX, newY) {
			newState := State{newX, newY, state.dir}
			if _, visited := distances[newState]; !visited {
				heap.Push(pq, &QueueItem{newState, cost + 1})
			}
		}

		// Turn Left and Left
		leftDir := (state.dir + 3) % 4
		rightDir := (state.dir + 1) % 4

		for _, newDir := range []int{leftDir, rightDir} {
			newState := State{state.x, state.y, newDir}
			if _, visited := distances[newState]; !visited {
				heap.Push(pq, &QueueItem{newState, cost + 1000})
			}
		}
	}
	return distances, -1
}

func dijkstraBackward(grid [][]rune, endX, endY int) map[State]int {
	pq := &PriorityQueue{}
	distances := make(map[State]int)

	for dir := range 4 {
		endState := State{endX, endY, dir}
		heap.Push(pq, &QueueItem{endState, 0})
		distances[endState] = 0
	}

	for pq.Len() > 0 {
		current := heap.Pop(pq).(*QueueItem)
		state, cost := current.state, current.cost

		if prevCost, exists := distances[state]; exists && cost > prevCost {
			continue
		}

		// GO Backward
		backDir := (state.dir + 2) % 4
		dx, dy := directions[backDir][0], directions[backDir][1]
		prevX := state.x + dx
		prevY := state.y + dy
		if isValid(grid, prevX, prevY) {
			prevState := State{prevX, prevY, state.dir}
			newCost := cost + 1
			if prevCost, exists := distances[prevState]; !exists || newCost < prevCost {
				distances[prevState] = newCost
				heap.Push(pq, &QueueItem{prevState, newCost})
			}
		}

		// Reverse Turn left/right
		leftDir := (state.dir + 3) % 4
		rightDir := (state.dir + 1) % 4

		for _, prevDir := range []int{leftDir, rightDir} {
			prevState := State{state.x, state.y, prevDir}
			newCost := cost + 1000
			if prevCost, exists := distances[prevState]; !exists || newCost < prevCost {
				distances[prevState] = newCost
				heap.Push(pq, &QueueItem{prevState, newCost})
			}
		}

	}
	return distances
}

func countOptimalTiles(grid [][]rune, start State) int {
	var endX, endY int
	for i, row := range grid {
		for j, cell := range row {
			if cell == 'E' {
				endX, endY = i, j
				goto found
			}
		}
	}
found:

	distFromStart, minCost := dijkstraForward(grid, start)
	distToEnd := dijkstraBackward(grid, endX, endY)

	optimalTiles := make(map[[2]int]bool)
	for state, costFromStart := range distFromStart {
		if costToEnd, exists := distToEnd[state]; exists {
			if costFromStart+costToEnd == minCost {
				optimalTiles[[2]int{state.x, state.y}] = true
			}
		}
	}
	return len(optimalTiles)
}

func isValid(grid [][]rune, x, y int) bool {
	return x >= 0 && x < len(grid) && y >= 0 && y < len(grid[0]) && grid[x][y] != '#'
}

func directionName(dir int) string {
	names := []string{"North", "East", "South", "West"}
	return names[dir]
}
