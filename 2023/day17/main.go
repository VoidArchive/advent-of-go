package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"strconv"
)

type Point struct {
	x, y int
}

type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

func (d Direction) opposite() Direction {
	switch d {
	case Up:
		return Down
	case Down:
		return Up
	case Left:
		return Right
	case Right:
		return Left
	}
	return d
}

type State struct {
	pos   Point
	dir   Direction
	steps int
	heat  int
	index int
}

type PriorityQueue []*State

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].heat < pq[j].heat }
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	state := x.(*State)
	state.index = n
	*pq = append(*pq, state)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	state := old[n-1]
	old[n-1] = nil
	state.index = -1
	*pq = old[0 : n-1]
	return state
}

func readInput(filename string) [][]int {
	file, _ := os.Open(filename)
	defer file.Close()
	var grid [][]int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		row := make([]int, len(line))
		for i, ch := range line {
			digit, _ := strconv.Atoi(string(ch))
			row[i] = digit
		}
		grid = append(grid, row)
	}
	return grid
}

func solve(grid [][]int, minSteps, maxSteps int) int {
	rows, cols := len(grid), len(grid[0])
	target := Point{rows - 1, cols - 1}
	dx := []int{-1, 0, 1, 0}
	dy := []int{0, 1, 0, -1}

	visited := make(map[[4]int]int)
	pq := &PriorityQueue{}
	heap.Init(pq)

	heap.Push(pq, &State{Point{0, 0}, Right, 0, 0, -1})
	heap.Push(pq, &State{Point{0, 0}, Down, 0, 0, -1})

	for pq.Len() > 0 {
		current := heap.Pop(pq).(*State)

		if current.pos == target && current.steps >= minSteps {
			return current.heat
		}
		stateKey := [4]int{current.pos.x, current.pos.y, int(current.dir), current.steps}

		if prevHeat, exists := visited[stateKey]; exists && prevHeat <= current.heat {
			continue
		}
		visited[stateKey] = current.heat

		for newDir := Direction(0); newDir < 4; newDir++ {
			if newDir == current.dir.opposite() {
				continue
			}
			var newSteps int
			if newDir == current.dir {
				newSteps = current.steps + 1
				if newSteps > maxSteps {
					continue
				}
			} else {
				if current.steps < minSteps && current.steps > 0 {
					continue
				}
				newSteps = 1
			}
			newX := current.pos.x + dx[newDir]
			newY := current.pos.y + dy[newDir]

			if newX < 0 || newX >= rows || newY < 0 || newY >= cols {
				continue
			}
			newPos := Point{newX, newY}
			newHeat := current.heat + grid[newX][newY]

			newState := &State{
				newPos,
				newDir,
				newSteps,
				newHeat,
				-1,
			}
			heap.Push(pq, newState)
		}
	}
	return -1
}

func main() {
	grid := readInput("input.txt")
	fmt.Println("Part 1: ", solve(grid, 1, 3))
	fmt.Println("Part 2: ", solve(grid, 4, 10))
}
