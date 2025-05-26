package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func sumValidMul(input string) int {
	re := regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)`)
	matches := re.FindAllStringSubmatch(input, -1)

	sum := 0
	for _, match := range matches {
		x, _ := strconv.Atoi(match[1])
		y, _ := strconv.Atoi(match[2])
		sum += x * y
	}

	return sum
}

func sumConditionalValidMul(input string) int {
	re := regexp.MustCompile(`do\(\)|don't\(\)|mul\((\d{1,3}),(\d{1,3})\)`)

	enabled := true
	sum := 0

	matches := re.FindAllStringSubmatch(input, -1)
	for _, match := range matches {

		switch match[0] {
		case "do()":
			enabled = true
		case "don't()":
			enabled = false
		default:
			if !enabled {
				continue
			}
			x, _ := strconv.Atoi(match[1])
			y, _ := strconv.Atoi(match[2])

			sum += x * y

		}
	}
	return sum
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var builder strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		builder.WriteString(scanner.Text())
	}

	// INFO: Part 1
	part1Result := sumValidMul(builder.String())
	fmt.Printf("Part 1: Sum of valid mul operations: %d\n", part1Result)

	// INFO: Part 2

	part2Result := sumConditionalValidMul(builder.String())
	fmt.Printf("Part 2: Sum of conditional mul operations: %d\n", part2Result)
}
