package main

import (
	"fmt"
	"os"
	"strings"
)

type Pattern struct {
	grid []string
	rows int
	cols int
}

func parseInput(input string) []Pattern {
	patterns := []Pattern{}
	blocks := strings.SplitSeq(strings.TrimSpace(input), "\n\n")

	for block := range blocks {
		lines := strings.Split(strings.TrimSpace(block), "\n")
		if len(lines) > 0 {
			patterns = append(patterns, Pattern{
				grid: lines,
				rows: len(lines),
				cols: len(lines[0]),
			})
		}
	}
	return patterns
}

// Check if there's a vertical reflection line between columns left and left+1
func (p *Pattern) hasVerticalReflection(left int) bool {
	return p.countVerticalMismatches(left) == 0
}

// Check if there's a horizontal reflection line between rows above and above+1
func (p *Pattern) hasHorizontalReflection(above int) bool {
	return p.countHorizontalMismatches(above) == 0
}

// Count mismatches for vertical reflection (for smudge detection)
func (p *Pattern) countVerticalMismatches(left int) int {
	mismatches := 0
	maxExtent := min(left+1, p.cols-left-1)

	for row := 0; row < p.rows; row++ {
		for offset := range maxExtent {
			leftCol := left - offset
			rightCol := left + 1 + offset

			if p.grid[row][leftCol] != p.grid[row][rightCol] {
				mismatches++
			}
		}
	}
	return mismatches
}

// Count mismatches for horizontal reflection (for smudge detection)
func (p *Pattern) countHorizontalMismatches(above int) int {
	mismatches := 0
	maxExtent := min(above+1, p.rows-above-1)

	for offset := range maxExtent {
		topRow := above - offset
		bottomRow := above + 1 + offset

		for col := 0; col < p.cols; col++ {
			if p.grid[topRow][col] != p.grid[bottomRow][col] {
				mismatches++
			}
		}
	}
	return mismatches
}

// Find reflection with exactly one smudge (one mismatch)
func (p *Pattern) findSmudgedReflection() int {
	// Try vertical reflections (between columns)
	for col := 0; col < p.cols-1; col++ {
		if p.countVerticalMismatches(col) == 1 {
			return col + 1 // columns to the left
		}
	}

	// Try horizontal reflections (between rows)
	for row := 0; row < p.rows-1; row++ {
		if p.countHorizontalMismatches(row) == 1 {
			return 100 * (row + 1) // 100 * rows above
		}
	}

	return 0 // No smudged reflection found
}

func (p *Pattern) findReflection() int {
	// Try vertical reflections (between columns)
	for col := 0; col < p.cols-1; col++ {
		if p.hasVerticalReflection(col) {
			return col + 1 // columns to the left
		}
	}

	// Try horizontal reflections (between rows)
	for row := 0; row < p.rows-1; row++ {
		if p.hasHorizontalReflection(row) {
			return 100 * (row + 1) // 100 * rows above
		}
	}

	return 0 // No reflection found
}

func solvePart1(input string) int {
	patterns := parseInput(input)
	total := 0

	for i, pattern := range patterns {
		score := pattern.findReflection()
		fmt.Printf("Pattern %d (Part 1): score = %d\n", i+1, score)
		total += score
	}

	return total
}

func solvePart2(input string) int {
	patterns := parseInput(input)
	total := 0

	for i, pattern := range patterns {
		score := pattern.findSmudgedReflection()
		fmt.Printf("Pattern %d (Part 2): score = %d\n", i+1, score)
		total += score
	}

	return total
}

func main() {
	input, err := os.ReadFile("input.txt")
	if err != nil {
		fmt.Printf("Error reading input.txt: %v\n", err)
		return
	}

	fmt.Println("=== Part 1 ===")
	result1 := solvePart1(string(input))
	fmt.Printf("Part 1 Result: %d\n\n", result1)

	fmt.Println("=== Part 2 ===")
	result2 := solvePart2(string(input))
	fmt.Printf("Part 2 Result: %d\n", result2)
}
