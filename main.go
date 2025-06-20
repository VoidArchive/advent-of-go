package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func mix(secret, value int) int {
	return secret ^ value
}

func prune(secret int) int {
	return secret % 16777216
}

func nextSecret(secret int) int {
	result := secret * 64
	secret = mix(secret, result)
	secret = prune(secret)

	result = secret / 32
	secret = mix(secret, result)
	secret = prune(secret)

	result = secret * 2048
	secret = mix(secret, result)
	secret = prune(secret)

	return secret
}

func generateNthSecret(initial, n int) int {
	current := initial
	for range n {
		current = nextSecret(current)
	}
	return current
}

func solvePart1(initials []int) int {
	totalSum := 0
	for _, initial := range initials {
		secret2000 := generateNthSecret(initial, 2000)
		fmt.Printf("Initial: %d, 2000th: %d\n", initial, secret2000)
		totalSum += secret2000
	}
	return totalSum
}

func parseInputFile(filename string) ([]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var numbers []int
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		num, err := strconv.Atoi(line)
		if err != nil {
			return nil, err
		}
		numbers = append(numbers, num)
	}
	return numbers, nil
}

func generatePricesAndChanges(initial, count int) ([]int, []int) {
	prices := make([]int, count+1)
	prices[0] = initial % 10

	current := initial
	for i := 1; i <= count; i++ {
		current = nextSecret(current)
		prices[i] = current % 10
	}
	changes := make([]int, count)
	for i := range count {
		changes[i] = prices[i+1] - prices[i]
	}
	return prices, changes
}

func findFirstOccurrence(changes, prices []int, sequence [4]int) int {
	if len(changes) < 4 {
		return 0
	}
	for i := 0; i <= len(changes)-4; i++ {
		if changes[i] == sequence[0] &&
			changes[i+1] == sequence[1] &&
			changes[i+2] == sequence[2] &&
			changes[i+3] == sequence[3] {
			return prices[i+4]
		}
	}
	return 0
}

func solvePart2(initials []int) int {
	sequenceTotals := make(map[[4]int]int)
	for buyerIdx, initial := range initials {
		if buyerIdx%100 == 0 {
			fmt.Printf("Processing buyer %d/%d...\n", buyerIdx+1, len(initials))
		}
		prices, changes := generatePricesAndChanges(initial, 2000)
		foundSequences := make(map[[4]int]bool)

		for i := 0; i <= len(changes)-4; i++ {
			sequence := [4]int{changes[i], changes[i+1], changes[i+2], changes[i+3]}

			if !foundSequences[sequence] {
				foundSequences[sequence] = true
				price := prices[i+4]
				sequenceTotals[sequence] += price
			}
		}
	}

	maxBananas := 0
	var bestSequence [4]int

	for sequence, total := range sequenceTotals {
		if total > maxBananas {
			maxBananas = total
			bestSequence = sequence
		}
	}
	fmt.Printf("Best sequence: [%d,%d,%d,%d] gives %d bananas\n",
		bestSequence[0], bestSequence[1], bestSequence[2], bestSequence[3], maxBananas)

	return maxBananas
}

func main() {
	puzzleInput, err := parseInputFile("input.txt")
	if err != nil {
		return
	}

	result1 := solvePart1(puzzleInput)
	fmt.Printf("Part 1 answer :%d\n", result1)

	result2 := solvePart2(puzzleInput)
	fmt.Printf("Part 2 answer :%d\n", result2)
}
