package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Point struct {
	x, y int
}

var numericKeypad = map[Point]rune{
	{0, 0}: '7', {0, 1}: '8', {0, 2}: '9',
	{1, 0}: '4', {1, 1}: '5', {1, 2}: '6',
	{2, 0}: '1', {2, 1}: '2', {2, 2}: '3',
	{3, 1}: '0', {3, 2}: 'A',
}

var directionalKeypad = map[Point]rune{
	{0, 1}: '^', {0, 2}: 'A',
	{1, 0}: '<', {1, 1}: 'v', {1, 2}: '>',
}

var memo = make(map[string]int)

func getPosition(keypad map[Point]rune, key rune) Point {
	for pos, k := range keypad {
		if k == key {
			return pos
		}
	}
	panic(fmt.Sprintf("Key %c not found", key))
}

// Find the shortest path between two keys, returning the sequence
func findShortestPath(keypad map[Point]rune, from, to rune) []string {
	if from == to {
		return []string{"A"}
	}

	start := getPosition(keypad, from)
	target := getPosition(keypad, to)

	type state struct {
		pos  Point
		path string
	}
	queue := []state{{start, ""}}
	visited := make(map[Point]int)
	var allPaths []string
	minLength := -1
	directions := map[rune]Point{
		'^': {-1, 0}, 'v': {1, 0},
		'<': {0, -1}, '>': {0, 1},
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if minLength != -1 && len(current.path) > minLength {
			break
		}

		if current.pos == target {
			if minLength == -1 {
				minLength = len(current.path)
			}
			allPaths = append(allPaths, current.path+"A")
			continue
		}
		if prev, exists := visited[current.pos]; exists && prev < len(current.path) {
			continue
		}
		visited[current.pos] = len(current.path)

		for dirChar, delta := range directions {
			newPos := Point{current.pos.x + delta.x, current.pos.y + delta.y}
			if _, exists := keypad[newPos]; exists {
				queue = append(queue, state{newPos, current.path + string(dirChar)})
			}
		}
	}
	return allPaths
}

func isValidPath(keypad map[Point]rune, start Point, path string) bool {
	pos := start
	directions := map[rune]Point{
		'^': {-1, 0}, 'v': {1, 0},
		'<': {0, -1}, '>': {0, 1},
	}

	for _, dir := range path {
		if delta, ok := directions[dir]; ok {
			pos = Point{pos.x + delta.x, pos.y + delta.y}
			if _, exists := keypad[pos]; !exists {
				return false
			}
		}
	}
	return true
}

func findMinLength(sequence string, depth int) int {
	if depth == 0 {
		return len(sequence)
	}

	key := fmt.Sprintf("%s-%d", sequence, depth)
	if result, exists := memo[key]; exists {
		return result
	}

	total := 0
	current := 'A'

	for _, char := range sequence {
		allPaths := findShortestPath(directionalKeypad, current, char)
		minCost := int(^uint(0) >> 1)

		for _, path := range allPaths {
			cost := findMinLength(path, depth-1)
			if cost < minCost {
				minCost = cost
			}
		}
		total += minCost
		current = char
	}

	memo[key] = total
	return total
}

func solveCode(code string, directionalLevels int) int {
	total := 0
	current := 'A'

	for _, char := range code {
		allPaths := findShortestPath(numericKeypad, current, char)
		minCost := int(^uint(0) >> 1)

		for _, path := range allPaths {
			cost := findMinLength(path, directionalLevels)
			if cost < minCost {
				minCost = cost
			}
		}
		total += minCost
		current = char
	}
	return total
}

func main() {
	codes := []string{"805A", "964A", "459A", "968A", "671A"}
	totalComplexity := 0

	for _, code := range codes {
		// INFO: For part 2 change the directionalLevels to 25, without memo this would
		// have been impossible
		length := solveCode(code, 25)
		numericPart := strings.TrimSuffix(code, "A")
		numeric, _ := strconv.Atoi(numericPart)
		complexity := length * numeric

		fmt.Printf("Code %s: length=%d, numeric=%d, complexity=%d\n",
			code, length, numeric, complexity)
		totalComplexity += complexity
	}

	fmt.Printf("Total complexity: %d\n", totalComplexity)
}
