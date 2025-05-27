package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	grid, err := readGrid("input.txt")
	if err != nil {
		log.Fatal("reading grid: ", err)
	}

	fmt.Println("Part 1:", countXMAS(grid))
	fmt.Println("Part 2:", countXMASX(grid))
}

func readGrid(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	grid := make([]string, 0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		grid = append(grid, line)
	}

	return grid, scanner.Err()
}

func inbounds(x, y, rows, cols int) bool {
	return x >= 0 && x < rows && y >= 0 && y < cols
}

func countXMAS(g []string) int {
	dirs := [][2]int{
		{-1, -1},
		{-1, 0},
		{-1, 1},
		{0, -1},
		{0, 1},
		{1, -1},
		{1, 0},
		{1, 1},
	}

	target := "XMAS"
	count := 0
	rows, cols := len(g), len(g[0])
	for r, row := range g {
		for c, char := range row {
			if char != 'X' {
				continue
			}

			for _, dir := range dirs {
				match := true

				for i := 1; i < len(target); i++ {
					x, y := r+dir[0]*i, c+dir[1]*i
					if !inbounds(x, y, rows, cols) || g[x][y] != target[i] {
						match = false
						break
					}
				}

				if match {
					count++
				}
			}
		}
	}
	return count
}

func countXMASX(g []string) int {
	rows, cols := len(g), len(g[0])
	count := 0

	for r := 1; r < rows-1; r++ {
		for c := 1; c < cols-1; c++ {
			if g[r][c] != 'A' {
				continue
			}

			tl, tr := g[r-1][c-1], g[r-1][c+1]
			bl, br := g[r+1][c-1], g[r+1][c+1]

			diag1 := (tl == 'M' && br == 'S') || (tl == 'S' && br == 'M')
			diag2 := (tr == 'M' && bl == 'S') || (tr == 'S' && bl == 'M')

			if diag1 && diag2 {
				count++
			}
		}
	}
	return count
}
