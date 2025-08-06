package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Point struct {
	row, col int
}

func main() {
	exampleGrid := parseInput("input.txt")
	eans := solvePart1(exampleGrid)
	ans := solvePart2(exampleGrid)
	fmt.Println(eans)
	fmt.Println(ans)
}

func parseInput(filename string) [][]rune {
	file, _ := os.Open(filename)
	defer file.Close()

	grid := [][]rune{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		row := []rune(line)
		grid = append(grid, row)
	}
	return grid
}

func solvePart1(grid [][]rune) int {
	sum := 0
	for row := range len(grid) {
		var currentNumber string
		var numberStartCol int
		var inNumber bool
		for col := range len(grid[row]) {
			char := grid[row][col]

			if isDigit(char) {
				if !inNumber {
					inNumber = true
					numberStartCol = col
					currentNumber = ""
				}
				currentNumber += string(char)
			} else {
				if inNumber {
					if isAdjacentToSymbol(grid, row, numberStartCol, col-1) {
						num, _ := strconv.Atoi(currentNumber)
						sum += num
					}
					inNumber = false
				}
			}
		}
		if inNumber {
			if isAdjacentToSymbol(grid, row, numberStartCol, len(grid[row])-1) {
				num, _ := strconv.Atoi(currentNumber)
				sum += num
			}
		}
	}
	return sum
}

func isDigit(char rune) bool {
	return char >= '0' && char <= '9'
}

func isAdjacentToSymbol(grid [][]rune, row, startCol, endCol int) bool {
	minRow, maxRow := row-1, row+1
	minCol, maxCol := startCol-1, endCol+1

	for r := minRow; r <= maxRow; r++ {
		for c := minCol; c <= maxCol; c++ {
			if r >= 0 && r < len(grid) && c >= 0 && c < len(grid[r]) {
				if r == row && c >= startCol && c <= endCol {
					continue
				}
				char := grid[r][c]
				if char != '.' && !isDigit(char) {
					return true
				}
			}
		}
	}
	return false
}

func solvePart2(grid [][]rune) int {
	sum := 0
	for row := range len(grid) {
		for col := range len(grid[row]) {
			if grid[row][col] == '*' {
				numbers := findAdjacentNumbers(grid, row, col)

				if len(numbers) == 2 {
					gearRatio := numbers[0] * numbers[1]
					sum += gearRatio
				}
			}
		}
	}
	return sum
}

func findAdjacentNumbers(grid [][]rune, gearRow, gearCol int) []int {
	numbers := []int{}
	visited := make(map[Point]bool)

	for dr := -1; dr <= 1; dr++ {
		for dc := -1; dc <= 1; dc++ {
			if dr == 0 && dc == 0 {
				continue
			}
			r := gearRow + dr
			c := gearCol + dc

			if r >= 0 && r < len(grid) && c >= 0 && c < len(grid[r]) {
				if isDigit(grid[r][c]) {
					point := Point{r, c}
					if !visited[point] {
						number, startCol, endCol := extractNumberAt(grid, r, c)
						numbers = append(numbers, number)

						for col := startCol; col <= endCol; col++ {
							visited[Point{r, col}] = true
						}
					}
				}
			}
		}
	}
	return numbers
}

func extractNumberAt(grid [][]rune, row, col int) (int, int, int) {
	start := col
	for start > 0 && isDigit(grid[row][start-1]) {
		start--
	}

	end := col
	for end < len(grid[row])-1 && isDigit(grid[row][end+1]) {
		end++
	}

	numberStr := ""
	for c := start; c <= end; c++ {
		numberStr += string(grid[row][c])
	}
	number, _ := strconv.Atoi(numberStr)
	return number, start, end
}
