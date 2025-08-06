package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	allWinning, allHave := parseInput("input.txt")
	part1 := solvePart1(allWinning, allHave)
	fmt.Printf("Part 1: %d\n", part1)
	part2 := solvePart2(allWinning, allHave)
	fmt.Printf("Part 2: %d\n", part2)
}

func parseInput(filepath string) ([][]string, [][]string) {
	var allWinning [][]string
	var allHave [][]string
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		numberParts := strings.TrimSpace(parts[1])
		cards := strings.Split(numberParts, "|")
		winNums := strings.Fields(cards[0])
		haveNums := strings.Fields(cards[1])
		allWinning = append(allWinning, winNums)
		allHave = append(allHave, haveNums)
	}
	return allWinning, allHave
}

func solvePart1(allWinning [][]string, allHave [][]string) int {
	totalPoints := 0
	for i := range allWinning {
		winSet := make(map[string]bool)
		matches := 0
		for _, num := range allWinning[i] {
			winSet[num] = true
		}
		for _, num := range allHave[i] {
			if winSet[num] {
				matches++
			}
		}
		if matches > 0 {
			totalPoints += 1 << (matches - 1)
		}
	}
	return totalPoints
}

func solvePart2(allWinning [][]string, allHave [][]string) int {
	totalCards := make([]int, len(allWinning))
	for i := range totalCards {
		totalCards[i] = 1
	}
	for i := range allWinning {
		winSet := make(map[string]bool)
		matches := 0
		for _, num := range allWinning[i] {
			winSet[num] = true
		}
		for _, num := range allHave[i] {
			if winSet[num] {
				matches++
			}
		}
		for j := 1; j <= matches && i+j < len(totalCards); j++ {
			totalCards[i+j] += totalCards[i]
		}
	}
	sum := 0
	for _, val := range totalCards {
		sum += val
	}
	return sum
}
