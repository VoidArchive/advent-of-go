package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Robot struct {
	x, y   int
	dx, dy int
}

func parseRobot(line string) Robot {
	parts := strings.Split(line, " ")

	posStr := strings.TrimPrefix(parts[0], "p=")
	posCoords := strings.Split(posStr, ",")
	x, _ := strconv.Atoi(posCoords[0])
	y, _ := strconv.Atoi(posCoords[1])

	velStr := strings.TrimPrefix(parts[1], "v=")
	velCoords := strings.Split(velStr, ",")
	dx, _ := strconv.Atoi(velCoords[0])
	dy, _ := strconv.Atoi(velCoords[1])

	return Robot{x, y, dx, dy}
}

func simulatePosition(robot Robot, seconds, width, height int) (int, int) {
	finalX := (robot.x + robot.dx*seconds) % width
	finalY := (robot.y + robot.dy*seconds) % height

	if finalX < 0 {
		finalX += width
	}
	if finalY < 0 {
		finalY += height
	}
	return finalX, finalY
}

func calculateSafetyFactor(robots []Robot, seconds, width, height int) int {
	midX := width / 2
	midY := height / 2

	quadrants := [4]int{0, 0, 0, 0}

	for _, robot := range robots {
		x, y := simulatePosition(robot, seconds, width, height)

		if x == midX || y == midY {
			continue
		}

		var quadrant int
		if x < midX && y < midY {
			quadrant = 0
		} else if x > midX && y < midY {
			quadrant = 1
		} else if x < midX && y > midY {
			quadrant = 2
		} else {
			quadrant = 3
		}
		quadrants[quadrant]++
	}

	safetyFactor := 1
	for _, count := range quadrants {
		safetyFactor *= count
	}

	return safetyFactor
}

func calculateClustering(robots []Robot, seconds, width, height int) int {
	positions := make(map[[2]int]bool)
	for _, robot := range robots {
		x, y := simulatePosition(robot, seconds, width, height)
		positions[[2]int{x, y}] = true
	}

	clustered := 0
	for pos := range positions {
		x, y := pos[0], pos[1]

		neighbors := 0
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				if dx == 0 && dy == 0 {
					continue
				}
				nx, ny := x+dx, y+dy
				if positions[[2]int{nx, ny}] {
					neighbors++
				}
			}
		}
		if neighbors > 0 {
			clustered++
		}
	}
	return clustered
}

func visualizeRobots(robots []Robot, seconds, width, height int) {
	grid := make([][]int, height)
	for i := range grid {
		grid[i] = make([]int, width)
	}

	for _, robot := range robots {
		x, y := simulatePosition(robot, seconds, width, height)
		grid[y][x]++
	}

	fmt.Printf("\n Second %d:\n", seconds)
	for _, row := range grid {
		for _, count := range row {
			if count == 0 {
				fmt.Print(".")
			} else {
				fmt.Print("#")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func findChristmasTree(robots []Robot, width, height int) int {
	maxClustering := 0
	bestSecond := 0

	for seconds := 1; seconds <= 10000; seconds++ {
		clustering := calculateClustering(robots, seconds, width, height)

		if clustering > maxClustering {
			maxClustering = clustering
			bestSecond = seconds

			if clustering > len(robots)/2 {
				fmt.Printf("High clustering at second %d: %d robots clustered\n", seconds, clustering)
				// visualizeRobots(robots, seconds, width, height)
			}
		}

		if seconds%1000 == 0 {
			fmt.Printf("Checked %d seconds, max clustering: %d at second %d\n", seconds, maxClustering, bestSecond)
		}
	}
	return bestSecond
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var robots []Robot
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			robots = append(robots, parseRobot(line))
		}

	}
	width := 101
	height := 103
	seconds := 100

	// INFO: PART 1
	safetyFactor := calculateSafetyFactor(robots, seconds, width, height)
	fmt.Println(safetyFactor)

	// INFO: PART 2
	// Part 2
	fmt.Println("\nSearching for Christmas tree pattern...")
	treeSecond := findChristmasTree(robots, width, height)
	fmt.Printf("Part 2 - Christmas tree appears at second: %d\n", treeSecond)

	fmt.Printf("\nChecking second 6752 specifically:")
	clustering6752 := calculateClustering(robots, 6752, width, height)
	fmt.Printf("Second 6752 clustering: %d\n", clustering6752)

	fmt.Println("\nVisualization at detected best second:")
	visualizeRobots(robots, treeSecond, width, height)

	if treeSecond != 6752 {
		fmt.Println("\nVisualization at second 6752:")
		visualizeRobots(robots, 6752, width, height)
	}
}
