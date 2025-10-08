package main

import (
	"bufio"
	"fmt"
	"os"
)

type Point struct {
	x, y int
}

type Edge struct {
	to   Point
	dist int
}

func main() {
	file, _ := os.Open("input.txt")
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var grid []string

	for scanner.Scan() {
		grid = append(grid, scanner.Text())
	}

	rows := len(grid)
	cols := len(grid[0])

	var start, end Point
	for x := range cols {
		if grid[0][x] == '.' {
			start = Point{x, 0}
		}
		if grid[rows-1][x] == '.' {
			end = Point{x, rows - 1}
		}
	}

	maxSteps1 := dfs(grid, start, end, make(map[Point]bool), 0, true)
	fmt.Println("Part 1:", maxSteps1)

	graph := buildGraph(grid, start, end)
	maxSteps2 := dfsGraph(graph, start, end, make(map[Point]bool), 0)
	fmt.Println("Part 2:", maxSteps2)
}

func buildGraph(grid []string, start, end Point) map[Point][]Edge {
	rows := len(grid)
	cols := len(grid[0])

	junctions := make(map[Point]bool)
	junctions[start] = true
	junctions[end] = true

	for y := range rows {
		for x := range cols {
			if grid[y][x] == '#' {
				continue
			}
			neighbors := 0
			for _, d := range []Point{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} {
				nx, ny := x+d.x, y+d.y
				if nx >= 0 && nx < cols && ny >= 0 && ny < rows && grid[ny][nx] != '#' {
					neighbors++
				}
			}
			if neighbors > 2 {
				junctions[Point{x, y}] = true
			}
		}
	}

	graph := make(map[Point][]Edge)
	for junction := range junctions {
		edges := exploreFromJunction(grid, junction, junctions, rows, cols)
		graph[junction] = edges
	}
	return graph
}

func exploreFromJunction(grid []string, start Point, junctions map[Point]bool, rows, cols int) []Edge {
	var edges []Edge
	visited := make(map[Point]bool)
	visited[start] = true

	for _, d := range []Point{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} {
		next := Point{start.x + d.x, start.y + d.y}
		if next.x < 0 || next.x >= cols || next.y < 0 || next.y >= rows {
			continue
		}
		if grid[next.y][next.x] == '#' {
			continue
		}
		edge := walkPath(grid, start, next, junctions, rows, cols)
		if edge.dist > 0 {
			edges = append(edges, edge)
		}
	}
	return edges
}

func walkPath(grid []string, from, start Point, junctions map[Point]bool, rows, cols int) Edge {
	visited := make(map[Point]bool)
	visited[from] = true
	current := start
	steps := 1

	for {
		visited[current] = true
		if junctions[current] {
			return Edge{current, steps}
		}

		found := false
		for _, d := range []Point{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} {
			next := Point{current.x + d.x, current.y + d.y}
			if next.x < 0 || next.x >= cols || next.y < 0 || next.y >= rows {
				continue
			}
			if grid[next.y][next.x] == '#' {
				continue
			}
			if visited[next] {
				continue
			}
			current = next
			steps++
			found = true
			break
		}
		if !found {
			return Edge{Point{-1, -1}, 0}
		}
	}
}

func dfsGraph(graph map[Point][]Edge, pos, end Point, visited map[Point]bool, dist int) int {
	if pos == end {
		return dist
	}
	visited[pos] = true
	defer delete(visited, pos)
	maxPath := -1
	for _, edge := range graph[pos] {
		if visited[edge.to] {
			continue
		}
		result := dfsGraph(graph, edge.to, end, visited, dist+edge.dist)
		if result > maxPath {
			maxPath = result
		}
	}
	return maxPath
}

func dfs(grid []string, pos, end Point, visited map[Point]bool, steps int, useSlopes bool) int {
	if pos == end {
		return steps
	}

	rows := len(grid)
	cols := len(grid[0])

	visited[pos] = true
	defer delete(visited, pos)

	maxPath := -1
	directions := getDirections(grid, pos, useSlopes)
	for _, dir := range directions {
		next := Point{pos.x + dir.x, pos.y + dir.y}
		if next.x < 0 || next.x >= cols || next.y < 0 || next.y >= rows {
			continue
		}
		if grid[next.y][next.x] == '#' {
			continue
		}
		if visited[next] {
			continue
		}
		result := dfs(grid, next, end, visited, steps+1, useSlopes)
		if result > maxPath {
			maxPath = result
		}
	}
	return maxPath
}

func getDirections(grid []string, pos Point, useSlopes bool) []Point {
	if !useSlopes {
		return []Point{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
	}
	char := grid[pos.y][pos.x]

	switch char {
	case '^':
		return []Point{{0, -1}}
	case 'v':
		return []Point{{0, 1}}
	case '<':
		return []Point{{-1, 0}}
	case '>':
		return []Point{{1, 0}}
	default:
		return []Point{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}

	}
}
