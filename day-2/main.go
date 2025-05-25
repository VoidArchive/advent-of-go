package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func isSafe(level []int) bool {
	var isIncreasing, isDecreasing bool

	for i := range level[:len(level)-1] {

		diff := level[i] - level[i+1]
		if diff == 0 || abs(diff) > 3 {
			return false
		}

		if diff > 0 {
			isIncreasing = true
		} else {
			isDecreasing = true
		}

		if isIncreasing && isDecreasing {
			return false
		}
	}
	return true
}

func isDampenedSafe(level []int) bool {
	if isSafe(level) {
		return true
	}

	for i := range level {
		modifiedLevel := make([]int, 0, len(level)-1)
		modifiedLevel = append(modifiedLevel, level[:i]...)
		modifiedLevel = append(modifiedLevel, level[i+1:]...)

		if isSafe(modifiedLevel) {
			return true
		}

	}
	return false
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	levels := [][]int{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		level := make([]int, 0, len(fields))
		for _, field := range fields {
			value, err := strconv.Atoi(field)
			if err != nil {
				fmt.Println("Error parsing value:", err)
				return
			}
			level = append(level, value)
		}
		levels = append(levels, level)
	}

	// INFO: Part 1
	safeCount := 0
	for _, level := range levels {
		if isSafe(level) {
			safeCount++
		}
	}
	modifiedSafeCount := 0
	for _, level := range levels {
		if isDampenedSafe(level) {
			modifiedSafeCount++
		}
	}
	fmt.Printf("Part 1: %d\n", safeCount)
	fmt.Printf("Part 2: %d\n", modifiedSafeCount)
}
