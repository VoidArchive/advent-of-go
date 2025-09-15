package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Galaxy struct {
	row, col int
}

func parseInput(filename string) ([]Galaxy, int, int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, 0, 0, err
	}
	defer file.Close()

	var galaxies []Galaxy
	var grid []string

	scanner := bufio.NewScanner(file)
	row := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		grid = append(grid, line)

		for col, char := range line {
			if char == '#' {
				galaxies = append(galaxies, Galaxy{row: row, col: col})
			}
		}
		row++
	}

	if err := scanner.Err(); err != nil {
		return nil, 0, 0, err
	}

	rows := len(grid)
	cols := 0
	if rows > 0 {
		cols = len(grid[0])
	}

	return galaxies, rows, cols, nil
}

func findEmptyRowsCols(galaxies []Galaxy, totalRows, totalCols int) ([]int, []int) {
	// Find which rows and columns have galaxies
	hasGalaxyInRow := make([]bool, totalRows)
	hasGalaxyInCol := make([]bool, totalCols)

	for _, galaxy := range galaxies {
		hasGalaxyInRow[galaxy.row] = true
		hasGalaxyInCol[galaxy.col] = true
	}

	// Collect empty rows and columns
	var emptyRows, emptyCols []int
	for i := 0; i < totalRows; i++ {
		if !hasGalaxyInRow[i] {
			emptyRows = append(emptyRows, i)
		}
	}
	for i := 0; i < totalCols; i++ {
		if !hasGalaxyInCol[i] {
			emptyCols = append(emptyCols, i)
		}
	}

	return emptyRows, emptyCols
}

func manhattanDistance(g1, g2 Galaxy, emptyRows, emptyCols []int, expansionFactor int) int {
	// Basic Manhattan distance
	distance := abs(g1.row-g2.row) + abs(g1.col-g2.col)

	// Add extra distance for crossing empty rows
	minRow, maxRow := min(g1.row, g2.row), max(g1.row, g2.row)
	for _, emptyRow := range emptyRows {
		if emptyRow > minRow && emptyRow < maxRow {
			distance += expansionFactor - 1
		}
	}

	// Add extra distance for crossing empty columns
	minCol, maxCol := min(g1.col, g2.col), max(g1.col, g2.col)
	for _, emptyCol := range emptyCols {
		if emptyCol > minCol && emptyCol < maxCol {
			distance += expansionFactor - 1
		}
	}

	return distance
}

func solvePart(galaxies []Galaxy, emptyRows, emptyCols []int, expansionFactor int) int {
	totalDistance := 0

	// Calculate distance between every pair of galaxies
	for i := range galaxies {
		for j := i + 1; j < len(galaxies); j++ {
			distance := manhattanDistance(galaxies[i], galaxies[j], emptyRows, emptyCols, expansionFactor)
			totalDistance += distance
		}
	}

	return totalDistance
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	galaxies, totalRows, totalCols, err := parseInput("input.txt")
	if err != nil {
		fmt.Printf("Error parsing input: %v\n", err)
		return
	}

	emptyRows, emptyCols := findEmptyRowsCols(galaxies, totalRows, totalCols)

	part1Result := solvePart(galaxies, emptyRows, emptyCols, 2)
	part2Result := solvePart(galaxies, emptyRows, emptyCols, 1000000)

	fmt.Printf("Part 1: %d\n", part1Result)
	fmt.Printf("Part 2: %d\n", part2Result)
}
