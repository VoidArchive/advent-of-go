package main

import (
	"bufio"
	"fmt"
	"os"
)

type Point struct {
	row, col int
}

func countReachablePlots(grid [][]rune, start Point, steps int) int {
	rows := len(grid)
	cols := len(grid[0])

	current := make(map[Point]bool)
	current[start] = true

	dirs := []Point{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	for range steps {
		next := make(map[Point]bool)
		for pos := range current {
			for _, dir := range dirs {
				newPos := Point{pos.row + dir.row, pos.col + dir.col}
				if newPos.row >= 0 && newPos.row < rows &&
					newPos.col >= 0 && newPos.col < cols &&
					grid[newPos.row][newPos.col] != '#' {
					next[newPos] = true
				}
			}
		}
		current = next
	}
	return len(current)
}

func countReachablePlotsInfinite(grid [][]rune, start Point, steps int) int {
	rows := len(grid)
	cols := len(grid[0])

	current := make(map[Point]bool)
	current[start] = true

	dirs := []Point{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	halfSize := rows / 2
	samples := []int{}
	sampleSteps := []int{halfSize, halfSize + rows, halfSize + 2*rows}
	sampleIdx := 0

	for step := range steps {
		next := make(map[Point]bool)
		for pos := range current {
			for _, dir := range dirs {
				newPos := Point{pos.row + dir.row, pos.col + dir.col}
				gridRow := ((newPos.row % rows) + rows) % rows
				gridCol := ((newPos.col % cols) + cols) % cols
				if grid[gridRow][gridCol] != '#' {
					next[newPos] = true
				}
			}
		}
		current = next

		if sampleIdx < len(sampleSteps) && step+1 == sampleSteps[sampleIdx] {
			samples = append(samples, len(current))
			sampleIdx++
			if len(samples) == 3 {
				return extrapolateQuadratic(samples, steps, rows)
			}
		}
	}
	return len(current)
}

func extrapolateQuadratic(samples []int, target, gridSize int) int {
	y0, y1, y2 := samples[0], samples[1], samples[2]
	a := (y2 - 2*y1 + y0) / 2
	b := y1 - y0 - a
	c := y0
	halfSize := gridSize / 2
	x := (target - halfSize) / gridSize
	return a*x*x + b*x + c
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var grid [][]rune
	var start Point
	scanner := bufio.NewScanner(file)
	row := 0

	for scanner.Scan() {
		line := scanner.Text()
		gridRow := []rune(line)
		for col, ch := range gridRow {
			if ch == 'S' {
				start = Point{row, col}
			}
		}
		grid = append(grid, gridRow)
		row++
	}

	result1 := countReachablePlots(grid, start, 64)
	fmt.Println("Part 1:", result1)

	result2 := countReachablePlotsInfinite(grid, start, 26501365)

	fmt.Println("Part 2:", result2)
}
