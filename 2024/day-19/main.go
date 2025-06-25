package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	scanner.Scan()
	patternsLine := scanner.Text()
	patterns := make([]string, 0)
	for p := range strings.SplitSeq(patternsLine, ", ") {
		patterns = append(patterns, strings.TrimSpace(p))
	}

	scanner.Scan()
	designs := make([]string, 0)
	for scanner.Scan() {
		designs = append(designs, strings.TrimSpace(scanner.Text()))
	}

	possible := 0
	for _, design := range designs {
		if countWays(design, patterns) > 0 {
			possible++
		}
	}
	fmt.Printf("Possible designs: %d\n", possible)

	totalWays := 0
	for _, design := range designs {
		ways := countWays(design, patterns)
		totalWays += ways
	}
	fmt.Printf("Total Ways: %d\n", totalWays)
}

func countWays(design string, patterns []string) int {
	memo := make(map[string]int)
	return countWaysMemo(design, patterns, memo)
}

func countWaysMemo(design string, patterns []string, memo map[string]int) int {
	if design == "" {
		return 1
	}

	if result, exists := memo[design]; exists {
		return result
	}

	ways := 0
	for _, pattern := range patterns {
		if strings.HasPrefix(design, pattern) {
			remaining := design[len(pattern):]
			ways += countWaysMemo(remaining, patterns, memo)
		}
	}

	memo[design] = ways
	return ways
}
