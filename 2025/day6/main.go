package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// 123 328  51 64
//  45 64  387 23
//   6 98  215 314
// *   +   *   +

// 123 * 45 * 6 = 33210
// 328 + 64 + 98 = 490
// 51 * 387 * 215 = 4243455
// 64 + 23 + 314 = 401

func readLines(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer func() { _ = file.Close() }()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func part1(lines []string) int {
	var grid [][]int
	total := 0
	for i := range 4 {
		fields := strings.Fields(lines[i])
		var row []int
		for _, f := range fields {
			n, _ := strconv.Atoi(f)
			row = append(row, n)
		}
		grid = append(grid, row)
	}

	var ops []rune
	for _, char := range lines[4] {
		if char == '*' || char == '+' {
			ops = append(ops, char)
		}
	}

	for col := range len(grid[0]) {
		op := ops[col]

		if op == '*' {
			product := 1
			for row := range len(grid) {
				product *= grid[row][col]
			}
			total += product

		} else {
			sum := 0
			for row := range len(grid) {
				sum += grid[row][col]
			}
			total += sum
		}
	}

	return total
}

func part2(lines []string) int {
	total := 0
	opLine := lines[4]

	isSeperator := func(col int) bool {
		for row := range 4 {
			if col < len(lines[row]) && lines[row][col] != ' ' {
				return false
			}
		}
		return true
	}

	col := 0
	for col < len(opLine) {
		if isSeperator(col) {
			col++
			continue
		}
		start := col
		for col < len(opLine) && !isSeperator(col) {
			col++
		}
		end := col - 1
		var op rune
		for c := start; c <= end; c++ {
			if c < len(opLine) && (opLine[c] == '*' || opLine[c] == '+') {
				op = rune(opLine[c])
				break
			}
		}
		if op == 0 {
			continue
		}

		var nums []int
		for c := start; c <= end; c++ {
			numStr := ""
			for row := range 4 {
				if c < len(lines[row]) && lines[row][c] >= '0' && lines[row][c] <= '9' {
					numStr += string(lines[row][c])
				}
			}
			if numStr != "" {
				n, _ := strconv.Atoi(numStr)
				nums = append(nums, n)
			}
		}

		if op == '*' {
			product := 1
			for _, n := range nums {
				product *= n
			}
			total += product
		} else {
			sum := 0
			for _, n := range nums {
				sum += n
			}
			total += sum
		}
	}
	return total
}

func main() {
	lines := readLines("input.txt")

	start := time.Now()
	fmt.Println("Part 1:", part1(lines))
	fmt.Println("Part 2:", part2(lines))
	fmt.Println("Time:", time.Since(start))
}
