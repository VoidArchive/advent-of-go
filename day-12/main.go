package main

import (
	"bufio"
	"fmt"
	"os"
)

type Point struct {
	R int
	C int
}

func readInput(filePath string) ([][]rune, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %w", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	gardenMap := [][]rune{}

	for scanner.Scan() {
		line := scanner.Text()
		gardenMap = append(gardenMap, []rune(line))
	}
	return gardenMap, nil
}

func bfs(gardenMap [][]rune, visited [][]bool, startR, startC int, plantType rune, numRows, numCols int) (area, perimeter int) {
	dr := []int{-1, 1, 0, 0}
	dc := []int{0, 0, -1, 1}

	queue := []Point{{startR, startC}}

	for len(queue) > 0 {
		currPoint := queue[0]
		queue = queue[1:]
		area++

		for i := range 4 {
			nextR, nextC := currPoint.R+dr[i], currPoint.C+dc[i]
			if nextR < 0 || nextR >= numRows || nextC < 0 || nextC >= numCols || gardenMap[nextR][nextC] != plantType {
				perimeter++
			} else {
				if !visited[nextR][nextC] {
					visited[nextR][nextC] = true
					queue = append(queue, Point{nextR, nextC})
				}
			}
		}
	}
	return area, perimeter
}

func solvePart1(gardenMap [][]rune) int {
	if len(gardenMap) == 0 {
		fmt.Println("Garden Map is empty")
		return 0
	}

	numRows := len(gardenMap)
	numCols := len(gardenMap[0])

	visited := make([][]bool, numRows)
	for i := range visited {
		visited[i] = make([]bool, numCols)
	}

	totalFencePrice := 0
	for r := range numRows {
		for c := range numCols {
			if !visited[r][c] {
				currentPlantType := gardenMap[r][c]
				visited[r][c] = true

				// INFO: Calling BFS
				area, perimeter := bfs(gardenMap, visited, r, c, currentPlantType, numRows, numCols)

				regionPrice := area * perimeter
				totalFencePrice += regionPrice
			}
		}
	}
	return totalFencePrice
}

func main() {
	gardenMap, err := readInput("input.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	resultPart1 := solvePart1(gardenMap)
	fmt.Println(resultPart1)
}
