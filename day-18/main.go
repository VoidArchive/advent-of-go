package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	x, y int
}

type State struct {
	pos   Point
	steps int
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var bytePositions []Point

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			continue
		}
		x, _ := strconv.Atoi(parts[0])
		y, _ := strconv.Atoi(parts[1])

		bytePositions = append(bytePositions, Point{x, y})
	}

	const gridSize = 71
	const bytesToSimulate = 1024

	corrupted := make(map[Point]bool)
	for i := 0; i < bytesToSimulate && i < len(bytePositions); i++ {
		corrupted[bytePositions[i]] = true
	}

	start := Point{0, 0}
	end := Point{70, 70}

	steps := bfs(start, end, corrupted, gridSize)
	if steps == -1 {
		fmt.Println("No path found")
	} else {
		fmt.Printf("Part 1: Minimum steps needed: %d\n", steps)
	}

	blockingByte := findBlockingByte(bytePositions, start, end, gridSize)
	if blockingByte.x == -1 {
		fmt.Println("No blocking byte found")
	} else {
		fmt.Printf("Part 2: First blocking byte : %d, %d\n", blockingByte.x, blockingByte.y)
	}
}

func findBlockingByte(bytePositions []Point, start, end Point, gridSize int) Point {
	left, right := 0, len(bytePositions)-1
	result := Point{-1, -1}

	for left <= right {
		mid := (left + right) / 2

		corrupted := make(map[Point]bool)
		for i := 0; i <= mid; i++ {
			corrupted[bytePositions[i]] = true
		}

		if bfs(start, end, corrupted, gridSize) == -1 {
			result = bytePositions[mid]
			right = mid - 1
		} else {
			left = mid + 1
		}
	}
	return result
}

func bfs(start, end Point, corrupted map[Point]bool, gridSize int) int {
	if corrupted[start] || corrupted[end] {
		return -1
	}
	queue := []State{{start, 0}}
	visited := make(map[Point]bool)
	visited[start] = true

	directions := []Point{{0, 1}, {1, 0}, {0, -1}, {-1, 0}} // down, right, up, left

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current.pos == end {
			return current.steps
		}

		for _, dir := range directions {
			newPos := Point{
				x: current.pos.x + dir.x,
				y: current.pos.y + dir.y,
			}

			if newPos.x < 0 || newPos.x >= gridSize || newPos.y < 0 || newPos.y >= gridSize {
				continue
			}

			if corrupted[newPos] || visited[newPos] {
				continue
			}

			visited[newPos] = true
			queue = append(queue, State{newPos, current.steps + 1})
		}
	}

	return -1
}
