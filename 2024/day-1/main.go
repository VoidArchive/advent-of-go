package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var columnOne []int
	var columnTwo []int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		a, _ := strconv.Atoi(fields[0])
		b, _ := strconv.Atoi(fields[1])

		columnOne = append(columnOne, a)
		columnTwo = append(columnTwo, b)
	}

	sort.Ints(columnOne)
	sort.Ints(columnTwo)

	// INFO: Part 1
	totalDistance := 0
	for i := range columnOne {
		totalDistance += abs(columnOne[i] - columnTwo[i])
	}

	fmt.Printf("Total distance: %d\n", totalDistance)

	// INFO: Part 2
	freq := make(map[int]int)
	for _, num := range columnTwo {
		freq[num]++
	}

	similarityScore := 0
	for _, num := range columnOne {
		similarityScore += num * freq[num]
	}
	fmt.Printf("Similarity score: %d\n", similarityScore)
}
