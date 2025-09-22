package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Part 1
	grid1 := parseInput("input.txt")
	tiltNorth(grid1)
	load1 := calculateLoad(grid1)
	fmt.Printf("Part 1 - Total load after tilting north: %d\n", load1)

	// Part 2
	grid2 := parseInput("input.txt")
	load2 := spinCycles(grid2, 1000000000)
	fmt.Printf("Part 2 - Total load after 1000000000 cycles: %d\n", load2)
}

func parseInput(filename string) [][]rune {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var grid [][]rune
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			grid = append(grid, []rune(line))
		}
	}

	return grid
}

func tiltNorth(grid [][]rune) {
	rows := len(grid)
	cols := len(grid[0])

	for col := range cols {
		for row := range rows {
			if grid[row][col] == 'O' {
				newRow := row
				for newRow > 0 && grid[newRow-1][col] == '.' {
					newRow--
				}
				if newRow != row {
					grid[row][col] = '.'
					grid[newRow][col] = 'O'
				}
			}
		}
	}
}

func tiltWest(grid [][]rune) {
	rows := len(grid)
	cols := len(grid[0])

	for row := range rows {
		for col := range cols {
			if grid[row][col] == 'O' {
				newCol := col
				for newCol > 0 && grid[row][newCol-1] == '.' {
					newCol--
				}
				if newCol != col {
					grid[row][col] = '.'
					grid[row][newCol] = 'O'
				}
			}
		}
	}
}

func tiltSouth(grid [][]rune) {
	rows := len(grid)
	cols := len(grid[0])

	for col := range cols {
		for row := rows - 1; row >= 0; row-- {
			if grid[row][col] == 'O' {
				newRow := row
				for newRow < rows-1 && grid[newRow+1][col] == '.' {
					newRow++
				}
				if newRow != row {
					grid[row][col] = '.'
					grid[newRow][col] = 'O'
				}
			}
		}
	}
}

func tiltEast(grid [][]rune) {
	rows := len(grid)
	cols := len(grid[0])

	for row := range rows {
		for col := cols - 1; col >= 0; col-- {
			if grid[row][col] == 'O' {
				newCol := col
				for newCol < cols-1 && grid[row][newCol+1] == '.' {
					newCol++
				}
				if newCol != col {
					grid[row][col] = '.'
					grid[row][newCol] = 'O'
				}
			}
		}
	}
}

func spinCycle(grid [][]rune) {
	tiltNorth(grid)
	tiltWest(grid)
	tiltSouth(grid)
	tiltEast(grid)
}

func gridToString(grid [][]rune) string {
	var sb strings.Builder
	for _, row := range grid {
		sb.WriteString(string(row))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func spinCycles(grid [][]rune, cycles int) int {
	seen := make(map[string]int)
	loads := make([]int, 0)

	for i := range cycles {
		spinCycle(grid)

		gridState := gridToString(grid)
		if firstSeen, exists := seen[gridState]; exists {
			// Found a cycle
			cycleLength := i - firstSeen
			remaining := (cycles - 1 - i) % cycleLength
			targetIndex := firstSeen + remaining
			return loads[targetIndex]
		}

		seen[gridState] = i
		loads = append(loads, calculateLoad(grid))
	}

	return calculateLoad(grid)
}

func calculateLoad(grid [][]rune) int {
	rows := len(grid)
	totalLoad := 0

	for row := range rows {
		for col := 0; col < len(grid[row]); col++ {
			if grid[row][col] == 'O' {
				load := rows - row
				totalLoad += load
			}
		}
	}

	return totalLoad
}

func testExample() {
	example := []string{
		"O....#....",
		"O.OO#....#",
		".....##...",
		"OO.#O....O",
		".O.....O#.",
		"O.#..O.#.#",
		"..O..#O..O",
		".......O..",
		"#....###..",
		"#OO..#....",
	}

	grid := make([][]rune, len(example))
	for i, line := range example {
		grid[i] = []rune(line)
	}

	fmt.Println("Testing Part 1:")
	grid1 := make([][]rune, len(grid))
	for i := range grid {
		grid1[i] = make([]rune, len(grid[i]))
		copy(grid1[i], grid[i])
	}

	tiltNorth(grid1)
	load1 := calculateLoad(grid1)
	fmt.Printf("Part 1 load: %d (expected: 136)\n", load1)

	fmt.Println("\nTesting Part 2:")
	load2 := spinCycles(grid, 1000000000)
	fmt.Printf("Part 2 load: %d (expected: 64)\n", load2)
}

func printGrid(grid [][]rune) {
	for _, row := range grid {
		fmt.Println(string(row))
	}
}
