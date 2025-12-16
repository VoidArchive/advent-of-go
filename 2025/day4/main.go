package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	start := time.Now()

	file, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var grid [][]byte
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		grid = append(grid, []byte(scanner.Text()))
	}

	fmt.Println("Part 1:", solvePart1(grid))
	fmt.Println("Part 2:", solvePart2(grid))

	fmt.Printf("Elapsed: %v\n", time.Since(start))
}

func solvePart1(grid [][]byte) int {
	// 8 directions: N, NE, E, SE, S, SW, W, NW
	dirs := [][2]int{
		{-1, 0},
		{-1, 1},
		{0, 1},
		{1, 1},
		{1, 0},
		{1, -1},
		{0, -1},
		{-1, -1},
	}

	count := 0
	rows := len(grid)
	cols := len(grid[0])

	for r := range rows {
		for c := range cols {
			if grid[r][c] != '@' {
				continue
			}
			adjacent := 0
			for _, d := range dirs {
				nr, nc := r+d[0], c+d[1]
				if 0 <= nr && nr < rows && 0 <= nc && nc < cols && grid[nr][nc] == '@' {
					adjacent++
				}
			}
			if adjacent < 4 {
				count++
			}
		}
	}
	return count
}

func solvePart2(grid [][]byte) int {
	// 8 directions: N, NE, E, SE, S, SW, W, NW
	dirs := [][2]int{
		{-1, 0},
		{-1, 1},
		{0, 1},
		{1, 1},
		{1, 0},
		{1, -1},
		{0, -1},
		{-1, -1},
	}

	count := 0
	rows := len(grid)
	cols := len(grid[0])

	for {
		removed := 0
		for r := range rows {
			for c := range cols {
				if grid[r][c] != '@' {
					continue
				}
				adjacent := 0
				for _, d := range dirs {
					nr, nc := r+d[0], c+d[1]
					if 0 <= nr && nr < rows && 0 <= nc && nc < cols && grid[nr][nc] == '@' {
						adjacent++
					}
				}
				if adjacent < 4 {
					removed++
					grid[r][c] = '.'
				}
			}
		}
		if removed == 0 {
			break
		}
		count += removed

	}
	return count
}
