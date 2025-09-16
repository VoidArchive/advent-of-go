package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type MemoKey struct {
	pos       int
	groupIdx  int
	groupSize int
}

func countArrangements(springs string, groups []int, memo map[MemoKey]int, pos, groupIdx, groupSize int) int {
	key := MemoKey{pos, groupIdx, groupSize}
	if val, exists := memo[key]; exists {
		return val
	}

	if pos == len(springs) {
		if groupIdx == len(groups) && groupSize == 0 {
			return 1
		}
		if groupIdx == len(groups)-1 && groupSize == groups[groupIdx] {
			return 1
		}
		return 0
	}
	result := 0
	current := springs[pos]

	possibilities := []rune{rune(current)}

	if current == '?' {
		possibilities = []rune{'.', '#'}
	}

	for _, char := range possibilities {
		switch char {
		case '.':
			if groupSize == 0 {
				result += countArrangements(springs, groups, memo, pos+1, groupIdx, 0)
			} else if groupIdx < len(groups) && groupSize == groups[groupIdx] {
				result += countArrangements(springs, groups, memo, pos+1, groupIdx+1, 0)
			}
		case '#':
			if groupIdx < len(groups) && groupSize < groups[groupIdx] {
				result += countArrangements(springs, groups, memo, pos+1, groupIdx, groupSize+1)
			}
		}
	}
	memo[key] = result
	return result
}

func solveLine(line string, unfold bool) int {
	parts := strings.Fields(line)
	if len(parts) != 2 {
		return 0
	}

	springs := parts[0]
	groupStrs := strings.Split(parts[1], ",")
	groups := make([]int, len(groupStrs))

	for i, s := range groupStrs {
		val, err := strconv.Atoi(s)
		if err != nil {
			return 0
		}
		groups[i] = val
	}
	if unfold {
		unfoldedSprings := make([]string, 5)
		for i := range 5 {
			unfoldedSprings[i] = springs
		}
		springs = strings.Join(unfoldedSprings, "?")
		originalGroups := make([]int, len(groups))
		copy(originalGroups, groups)
		groups = make([]int, 0, len(originalGroups)*5)
		for range 5 {
			groups = append(groups, originalGroups...)
		}
	}
	memo := make(map[MemoKey]int)
	return countArrangements(springs, groups, memo, 0, 0, 0)
}

func main() {
	var lines []string
	file, _ := os.Open("input.txt")
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	total := 0
	for _, line := range lines {
		arrangements := solveLine(line, false)
		fmt.Printf("%s -> %d arrangements\n", line, arrangements)
		total += arrangements
	}
	fmt.Printf("\nTotal arrangements: %d\n", total)

	totalPart2 := 0
	for _, line := range lines {
		arrangements := solveLine(line, true)
		totalPart2 += arrangements
	}
	fmt.Printf("\nPart2 Total = %d\n", totalPart2)
}
