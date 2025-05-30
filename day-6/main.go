package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type State struct {
	x, y, dir int
}

var directions = []State{
	{-1, 0, 0}, // up
	{0, 1, 0},  // right
	{1, 0, 0},  // down
	{0, -1, 0}, // left
}

var symbols = map[byte]int{
	'^': 0,
	'>': 1,
	'v': 2,
	'<': 3,
}

func readLines(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return lines
}

func copyGrid(lines []string) [][]byte {
	grid := make([][]byte, len(lines))
	for i := range lines {
		grid[i] = []byte(lines[i])
	}
	return grid
}

func findGuard(grid [][]byte) (int, int, int) {
	for i := range grid {
		for j, ch := range grid[i] {
			if d, ok := symbols[ch]; ok {
				grid[i][j] = '.'
				return i, j, d
			}
		}
	}
	panic("Guard not found in the grid")
}

func simulateGuardPath(grid [][]byte, sx, sy, dir int) int {
	visited := map[[2]int]bool{
		{sx, sy}: true,
	}
	x, y := sx, sy

	for {
		dx, dy := directions[dir].x, directions[dir].y
		nx, ny := x+dx, y+dy

		if nx < 0 || nx >= len(grid) || ny < 0 || ny >= len(grid[0]) {
			break
		}

		if grid[nx][ny] == '#' {
			dir = (dir + 1) % 4
		} else {
			x, y = nx, ny
			visited[[2]int{x, y}] = true
		}
	}
	return len(visited)
}

func causesLoop(grid [][]byte, sx, sy, dir int) bool {
	seen := make(map[State]bool)
	x, y := sx, sy

	for {
		state := State{x, y, dir}
		if seen[state] {
			return true
		}
		seen[state] = true

		dx, dy := directions[dir].x, directions[dir].y
		nx, ny := x+dx, y+dy

		if nx < 0 || nx >= len(grid) || ny < 0 || ny >= len(grid[0]) {
			return false
		}

		if grid[nx][ny] == '#' {
			dir = (dir + 1) % 4
		} else {
			x, y = nx, ny
		}
	}
}

func countLoopObstacles(lines []string, sx, sy, dir int) int {
	count := 0

	for i, row := range lines {
		for j := range row {
			if row[j] != '.' || (i == sx && j == sy) {
				continue
			}

			grid := copyGrid(lines)
			grid[i][j] = '#'

			if causesLoop(grid, sx, sy, dir) {
				count++
			}
		}
	}
	return count
}

func main() {
	lines := readLines("input.txt")
	grid := copyGrid(lines)
	sx, sy, dir := findGuard(grid)
	// INFO: Part 1
	fmt.Println("Part 1:", simulateGuardPath(grid, sx, sy, dir))
	fmt.Println("Part 2:", countLoopObstacles(lines, sx, sy, dir))
}
