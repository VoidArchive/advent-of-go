package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	R, C int
}

type SearchContext struct {
	Grid [][]int
	Rows int
	Cols int
}

func readandParseInput(input string) ([][]int, int, int, error) {
	file, err := os.Open(input)
	if err != nil {
		return nil, 0, 0, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		return nil, 0, 0, err
	}
	if len(lines) == 0 {
		return nil, 0, 0, fmt.Errorf("no input data found")
	}
	rows := len(lines)
	cols := len(lines[0])
	parsedGrid := make([][]int, rows)
	for row := range rows {
		if len(lines[row]) != cols {
			return nil, 0, 0, fmt.Errorf("inconsistent row length at row %d", row)
		}
		parsedGrid[row] = make([]int, cols)
		for col, char := range lines[row] {
			height, err := strconv.Atoi(string(char))
			if err != nil {
				return nil, 0, 0, fmt.Errorf("invalid character '%s' at row %d, col %d", string(char), row, col)
			}
			parsedGrid[row][col] = height
		}
	}
	return parsedGrid, rows, cols, nil
}

func dfsFindNines(r, c, currentHeight int, ninesFound map[Point]bool, ctx *SearchContext) {
	if currentHeight == 9 {
		ninesFound[Point{R: r, C: c}] = true
		return
	}
	dr := []int{-1, 1, 0, 0}
	dc := []int{0, 0, -1, 1}
	for i := range 4 {
		nr, nc := r+dr[i], c+dc[i]
		if nr >= 0 && nr < ctx.Rows && nc >= 0 && nc < ctx.Cols {
			if ctx.Grid[nr][nc] == currentHeight+1 {
				dfsFindNines(nr, nc, currentHeight+1, ninesFound, ctx)
			}
		}
	}
}

func dfsCountTrailsPart2(r, c, currentHeight int, ctx *SearchContext, memo [][][]int) int {
	if currentHeight == 9 {
		return 1
	}
	if memo[r][c][currentHeight] != -1 {
		return memo[r][c][currentHeight]
	}
	numberOfTrails := 0
	dr := []int{-1, 1, 0, 0}
	dc := []int{0, 0, -1, 1}
	for i := range 4 {
		nr, nc := r+dr[i], c+dc[i]
		if nr >= 0 && nr < ctx.Rows && nc >= 0 && nc < ctx.Cols && ctx.Grid[nr][nc] == currentHeight+1 {
			numberOfTrails += dfsCountTrailsPart2(nr, nc, currentHeight+1, ctx, memo)
		}
	}
	memo[r][c][currentHeight] = numberOfTrails
	return numberOfTrails
}

func findTrailheads(grid [][]int, rows, cols int) []Point {
	var trailheads []Point
	for r := range rows {
		for c := range cols {
			if grid[r][c] == 0 {
				trailheads = append(trailheads, Point{R: r, C: c})
			}
		}
	}
	return trailheads
}

func solvePart1(grid [][]int, rows, cols int) int {
	ctx := &SearchContext{Grid: grid, Rows: rows, Cols: cols}
	trailheads := findTrailheads(grid, rows, cols)

	totalScore := 0
	for _, trailhead := range trailheads {
		ninesFound := make(map[Point]bool)
		dfsFindNines(trailhead.R, trailhead.C, 0, ninesFound, ctx)
		totalScore += len(ninesFound)
	}
	return totalScore
}

func solvePart2(grid [][]int, rows, cols int) int {
	ctx := &SearchContext{Grid: grid, Rows: rows, Cols: cols}
	trailheads := findTrailheads(grid, rows, cols)

	// Initialize memoization table
	memo := make([][][]int, rows)
	for r := range rows {
		memo[r] = make([][]int, cols)
		for c := range cols {
			memo[r][c] = make([]int, 10)
			for h := range 10 {
				memo[r][c][h] = -1
			}
		}
	}

	totalRating := 0
	for _, trailhead := range trailheads {
		rating := dfsCountTrailsPart2(trailhead.R, trailhead.C, 0, ctx, memo)
		totalRating += rating
	}
	return totalRating
}

func main() {
	grid, rows, cols, err := readandParseInput("input.txt")
	if err != nil {
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}

	totalScore := solvePart1(grid, rows, cols)
	fmt.Printf("Total score for Part 1: %d\n", totalScore)

	totalRating := solvePart2(grid, rows, cols)
	fmt.Printf("Total rating for Part 2: %d\n", totalRating)
}
