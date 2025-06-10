package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type Point struct {
	x, y int
	dir  int
}

type State struct {
	P    Point
	cost int
}

type pqueue []*State

func (pq pqueue) Len() int           { return len(pq) }
func (pq pqueue) Less(i, j int) bool { return pq[i].cost < pq[j].cost }
func (pq pqueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}
func (pq *pqueue) Push(x any) { *pq = append(*pq, x.(*State)) }
func (pq *pqueue) Pop() any {
	n := len(*pq)
	it := (*pq)[n-1]
	*pq = (*pq)[:n-1]
	return it
}

func main() {
	grid, startPoint := parseGrid("input.txt")
	fmt.Printf("Start Point: %v\n", startPoint)

	// INFO: Part 1
	_, cost := dijkstraForward(grid, startPoint)
	fmt.Printf("Part 1 Solution: %d\n", cost)

	// INFO: Part 2
	optimalTiles := solvePart2(grid, startPoint)
	fmt.Printf("Part 2 Solution: %d\n", optimalTiles)
}

func parseGrid(input string) ([][]rune, Point) {
	file, err := os.Open(input)
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	grid := [][]rune{}
	var startPoint Point
	rowIndex := 0

	for scanner.Scan() {
		line := scanner.Text()
		grid = append(grid, []rune(line))

		for colIndex, char := range line {
			if char == 'S' {
				startPoint = Point{rowIndex, colIndex, 1}
			}
		}
		rowIndex++
	}
	return grid, startPoint
}

func dijkstraForward(grid [][]rune, start Point) (map[Point]int, int) {
	pq := pqueue{}
	heap.Push(&pq, &State{start, 0})
	dist := make(map[Point]int)
	dx := []int{-1, 0, 1, 0}
	dy := []int{0, 1, 0, -1}

	for pq.Len() > 0 {
		it := heap.Pop(&pq).(*State)
		cur, d := it.P, it.cost

		if cost, exists := dist[cur]; exists && d > cost {
			continue
		}
		dist[cur] = d

		if grid[cur.x][cur.y] == 'E' {
			return dist, d
		}

		// Go Forward
		newX := cur.x + dx[cur.dir]
		newY := cur.y + dy[cur.dir]
		if newX >= 0 && newX < len(grid) && newY >= 0 && newY < len(grid[0]) && grid[newX][newY] != '#' {
			newPoint := Point{newX, newY, cur.dir}
			if _, exists := dist[newPoint]; !exists {
				heap.Push(&pq, &State{newPoint, d + 1})
			}
		}

		// Turn Left
		leftDir := (cur.dir + 3) % 4
		newPoint := Point{cur.x, cur.y, leftDir}
		heap.Push(&pq, &State{newPoint, d + 1000})

		rightDir := (cur.dir + 1) % 4
		newPoint = Point{cur.x, cur.y, rightDir}
		heap.Push(&pq, &State{newPoint, d + 1000})
	}
	return dist, -1
}

func dijkstraBackward(grid [][]rune, endX, endY int) map[Point]int {
	pq := pqueue{}
	dist := make(map[Point]int)
	dx := []int{-1, 0, 1, 0}
	dy := []int{0, 1, 0, -1}

	for dir := range 4 {
		endState := Point{endX, endY, dir}
		heap.Push(&pq, &State{endState, 0})
		dist[endState] = 0
	}

	for pq.Len() > 0 {
		it := heap.Pop(&pq).(*State)
		cur, d := it.P, it.cost

		if cost, exists := dist[cur]; exists && d > cost {
			continue
		}

		// GO Backward
		backDir := (cur.dir + 2) % 4
		prevX := cur.x + dx[backDir]
		prevY := cur.y + dy[backDir]
		if prevX >= 0 && prevX < len(grid) && prevY >= 0 && prevY < len(grid[0]) && grid[prevX][prevY] != '#' {
			prevState := Point{prevX, prevY, cur.dir}
			if cost, exists := dist[prevState]; !exists || d+1 < cost {
				dist[prevState] = d + 1
				heap.Push(&pq, &State{prevState, d + 1})
			}
		}

		// Turn Left Back
		leftDir := (cur.dir + 3) % 4
		prevState := Point{cur.x, cur.y, leftDir}
		if cost, exists := dist[prevState]; !exists || d+1000 < cost {
			dist[prevState] = d + 1000
			heap.Push(&pq, &State{prevState, d + 1000})
		}
		rightDir := (cur.dir + 1) % 4
		prevState = Point{cur.x, cur.y, rightDir}
		if cost, exists := dist[prevState]; !exists || d+1000 < cost {
			dist[prevState] = d + 1000
			heap.Push(&pq, &State{prevState, d + 1000})
		}
	}
	return dist
}

func solvePart2(grid [][]rune, start Point) int {
	var endX, endY int
	for i, row := range grid {
		for j, cell := range row {
			if cell == 'E' {
				endX, endY = i, j
				break
			}
		}
	}
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
