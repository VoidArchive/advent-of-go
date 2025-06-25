package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func readLines(path string) [][]rune {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var grid [][]rune
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		grid = append(grid, []rune(line))
	}
	return grid
}

type Point struct {
	x, y int
}

func findAntinodes(grid [][]rune) map[Point]bool {
	width, height := len(grid[0]), len(grid)
	freqMap := make(map[rune][]Point)

	for y := range height {
		for x := range width {
			char := grid[y][x]
			if char != '.' {
				freqMap[char] = append(freqMap[char], Point{x, y})
			}
		}
	}
	result := make(map[Point]bool)

	for _, positions := range freqMap {
		for i := range len(positions) {
			for j := i + 1; j < len(positions); j++ {
				a := positions[i]
				b := positions[j]

				dx := b.x - a.x
				dy := b.y - a.y

				ax := a.x - dx
				ay := a.y - dy
				bx := b.x + dx
				by := b.y + dy

				if ax >= 0 && ax < width && ay >= 0 && ay < height {
					result[Point{ax, ay}] = true
				}

				if bx >= 0 && bx < width && by >= 0 && by < height {
					result[Point{bx, by}] = true
				}
			}
		}
	}
	return result
}

func printDebugGrid(grid [][]rune, antinodes map[Point]bool) {
	height := len(grid)
	width := len(grid[0])

	for y := range height {
		for x := range width {
			p := Point{x, y}
			if antinodes[p] {
				fmt.Print("#")
			} else if grid[y][x] != '.' {
				fmt.Print(string(grid[y][x]))
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func gcd(x, y int) int {
	for y != 0 {
		x, y = y, x%y
	}
	return x
}

func extendLine(start Point, dx, dy, width, height int) []Point {
	var points []Point
	x, y := start.x+dx, start.y+dy

	for 0 <= x && x < width && 0 <= y && y < height {
		points = append(points, Point{x, y})
		x += dx
		y += dy
	}

	return points
}

func findHarmonicAntinodes(grid [][]rune) map[Point]bool {
	width, height := len(grid[0]), len(grid)
	freqMap := make(map[rune][]Point)

	for y := range height {
		for x := range width {
			char := grid[y][x]
			if char != '.' {
				freqMap[char] = append(freqMap[char], Point{x, y})
			}
		}
	}
	result := make(map[Point]bool)

	for _, positions := range freqMap {
		n := len(positions)
		if n < 2 {
			continue
		}
		for i := range n {
			for j := i + 1; j < n; j++ {
				a := positions[i]
				b := positions[j]
				dx := b.x - a.x
				dy := b.y - a.y

				g := gcd(abs(dx), abs(dy))
				if g == 0 {
					continue
				}
				dx /= g
				dy /= g

				forward := extendLine(b, dx, dy, width, height)
				backward := extendLine(a, -dx, -dy, width, height)

				for _, p := range forward {
					result[p] = true
				}
				for _, p := range backward {
					result[p] = true
				}
				result[a] = true
				result[b] = true
			}
		}
	}

	return result
}

func main() {
	grid := readLines("input.txt")
	antinodes := findAntinodes(grid)
	harmonicAntinodes := findHarmonicAntinodes(grid)
	printDebugGrid(grid, harmonicAntinodes)
	fmt.Printf("Number of antinodes: %d\n", len(antinodes))
	fmt.Printf("Number of harmonic antinodes: %d\n", len(harmonicAntinodes))
}
