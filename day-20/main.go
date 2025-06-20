package main

import (
	"bufio"
	"fmt"
	"os"
)

type Point struct {
	x, y int
}

func main() {
	grid, start, end := parseInput()
	distances := findDistances(grid, start)
	normalTime := distances[end]
	fmt.Printf("Normal path time: %d picoseconds \n", normalTime)

	cheats := findCheats(distances, normalTime, 2)
	fmt.Printf("Cheats saving at least 100 picoseconds: %d\n", cheats)

	cheats2 := findCheats(distances, normalTime, 20)
	fmt.Printf("Cheats saving at least 100 picoseconds: %d\n", cheats2)
}

func parseInput() ([][]rune, Point, Point) {
	var grid [][]rune
	var start, end Point

	scanner := bufio.NewScanner(os.Stdin)
	y := 0

	for scanner.Scan() {
		line := scanner.Text()
		row := []rune(line)

		for x, ch := range row {
			if ch == 'S' {
				start = Point{x, y}
			} else if ch == 'E' {
				end = Point{x, y}
			}
		}
		grid = append(grid, row)
		y++
	}
	return grid, start, end
}

func findDistances(grid [][]rune, start Point) map[Point]int {
	distances := make(map[Point]int)
	queue := []Point{start}
	distances[start] = 0

	directions := []Point{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, dir := range directions {
			next := Point{current.x + dir.x, current.y + dir.y}

			if next.y >= 0 && next.y < len(grid) && next.x >= 0 && next.x < len(grid[0]) && grid[next.y][next.x] != '#' {
				if _, visited := distances[next]; !visited {
					distances[next] = distances[current] + 1
					queue = append(queue, next)
				}
			}
		}
	}
	return distances
}

func findCheats(distances map[Point]int, normalTime, maxCheatTime int) int {
	cheats := 0

	for startPos, distFromStart := range distances {
		for endPos, distToEnd := range distances {
			manhattanDist := abs(endPos.x-startPos.x) + abs(endPos.y-startPos.y)

			if manhattanDist > maxCheatTime {
				continue
			}

			if manhattanDist == 0 {
				continue
			}

			cheatTime := distFromStart + manhattanDist + (normalTime - distToEnd)
			timeSaved := normalTime - cheatTime

			if timeSaved >= 100 {
				cheats++
			}
		}
	}
	return cheats
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
