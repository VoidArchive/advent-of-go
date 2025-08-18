package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var histories [][]int
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		history := make([]int, len(fields))

		for i, field := range fields {
			num, err := strconv.Atoi(field)
			if err != nil {
				fmt.Printf("Error parsing number %s: %v\n", field, err)
				os.Exit(1)
			}
			history[i] = num
		}

		histories = append(histories, history)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	sumNext := 0
	for _, history := range histories {
		next := extrapolateNext(history)
		sumNext += next
	}
	sumPrev := 0
	for _, history := range histories {
		prev := extrapolatePrevious(history)
		sumPrev += prev
	}

	fmt.Printf("Part 1 - Sum of next extrapolated values: %d\n", sumNext)
	fmt.Printf("Part 2 - Sum of previous extrapolated values: %d\n", sumPrev)
}

func buildDifferenceSequences(sequence []int) [][]int {
	sequences := [][]int{sequence}

	for {
		current := sequences[len(sequences)-1]
		if allZeros(current) {
			break
		}

		differences := make([]int, len(current)-1)
		for i := 0; i < len(current)-1; i++ {
			differences[i] = current[i+1] - current[i]
		}

		sequences = append(sequences, differences)
	}

	return sequences
}

func extrapolateNext(sequence []int) int {
	sequences := buildDifferenceSequences(sequence)
	nextValue := 0

	for i := len(sequences) - 2; i >= 0; i-- {
		lastValue := sequences[i][len(sequences[i])-1]
		nextValue = lastValue + nextValue
	}

	return nextValue
}

func extrapolatePrevious(sequence []int) int {
	sequences := buildDifferenceSequences(sequence)
	prevValue := 0

	for i := len(sequences) - 2; i >= 0; i-- {
		firstValue := sequences[i][0]
		prevValue = firstValue - prevValue
	}
	return prevValue
}

func allZeros(slice []int) bool {
	for _, v := range slice {
		if v != 0 {
			return false
		}
	}
	return true
}
