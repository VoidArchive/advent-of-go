package main

import (
	"fmt"
	"os"
	"strings"
)

var digitWords = map[string]int{
	"one": 1, "two": 2, "three": 3, "four": 4, "five": 5,
	"six": 6, "seven": 7, "eight": 8, "nine": 9,
}

func solveTrebuchet(input string, includeWords bool) int {
	lines := strings.Split(input, "\n")
	sum := 0

	for _, line := range lines {
		if line == "" {
			continue
		}
		var a, b int
		foundFirst := false
		for i := range line {
			if digit, found := findDigitAt(line, i, includeWords); found {
				if !foundFirst {
					a = digit
					foundFirst = true
				}
				b = digit
			}
		}
		if foundFirst {
			sum += a*10 + b
		}
	}
	return sum
}

func findDigitAt(line string, pos int, includeWords bool) (int, bool) {
	if pos >= len(line) {
		return 0, false
	}
	if isDigit(rune(line[pos])) {
		return digitValue(rune(line[pos])), true
	}

	if includeWords {
		remaining := line[pos:]
		for word, value := range digitWords {
			if strings.HasPrefix(remaining, word) {
				return value, true
			}
		}
	}
	return 0, false
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func digitValue(c rune) int {
	return int(c - '0')
}

func main() {
	content, _ := os.ReadFile("input.txt")

	result := solveTrebuchet(string(content), false)
	fmt.Printf("Part 1: %d\n", result)
	result2 := solveTrebuchet(string(content), true)
	fmt.Printf("Part 2: %d\n", result2)
}
