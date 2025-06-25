package main

import (
	"bufio"
	"fmt"
	"os"
)

type Point struct {
	X, Y int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [1|2]")
		fmt.Println("  1 - Run Part 1 (normal warehouse)")
		fmt.Println("  2 - Run Part 2 (wide warehouse)")
		return
	}

	part := os.Args[1]
	grid, moves, robot := parseInputFile("input.txt")

	switch part {
	case "1":
		fmt.Println("=== PART 1: Normal Warehouse ===")
		solvePart1(grid, moves, robot)
	case "2":
		fmt.Println("=== PART 2: Wide Warehouse ===")
		solvePart2(grid, moves, robot)
	default:
		fmt.Println("Invalid part. Use 1 or 2")
	}
}

func solvePart1(grid [][]byte, moves string, robot Point) {
	fmt.Println("Initial grid:")
	printGridWithRobot(grid, robot)

	simulateNormal(grid, moves, &robot)

	fmt.Println("Final grid:")
	printGridWithRobot(grid, robot)

	sum := calculateNormalGPS(grid)
	fmt.Println("GPS sum:", sum)
}

func solvePart2(grid [][]byte, moves string, robot Point) {
	// Transform to wide warehouse
	wideGrid := transformToWide(grid)
	robot.X *= 2 // Robot's X position also doubles

	fmt.Println("Initial wide grid:")
	printGridWithRobot(wideGrid, robot)

	simulateWide(wideGrid, moves, &robot)

	fmt.Println("Final wide grid:")
	printGridWithRobot(wideGrid, robot)

	sum := calculateWideGPS(wideGrid)
	fmt.Println("GPS sum:", sum)
}

// ============ PART 1 FUNCTIONS ============

func simulateNormal(grid [][]byte, moves string, robot *Point) {
	for _, move := range moves {
		dx, dy := getDirection(byte(move))
		nx, ny := robot.X+dx, robot.Y+dy

		if grid[ny][nx] == '#' {
			continue
		} else if grid[ny][nx] == '.' {
			robot.X, robot.Y = nx, ny
		} else if grid[ny][nx] == 'O' {
			// Find end of box chain
			checkX, checkY := nx, ny
			for grid[checkY][checkX] == 'O' {
				checkX += dx
				checkY += dy
			}

			if grid[checkY][checkX] == '#' {
				continue
			} else if grid[checkY][checkX] == '.' {
				grid[checkY][checkX] = 'O'
				grid[ny][nx] = '.'
				robot.X, robot.Y = nx, ny
			}
		}
	}
}

func calculateNormalGPS(grid [][]byte) int {
	sum := 0
	for row := range grid {
		for col := range grid[row] {
			if grid[row][col] == 'O' {
				gps := 100*row + col
				sum += gps
			}
		}
	}
	return sum
}

// ============ PART 2 FUNCTIONS ============

func transformToWide(grid [][]byte) [][]byte {
	var wideGrid [][]byte

	for _, row := range grid {
		var wideRow []byte
		for _, cell := range row {
			switch cell {
			case '#':
				wideRow = append(wideRow, '#', '#')
			case 'O':
				wideRow = append(wideRow, '[', ']')
			case '.':
				wideRow = append(wideRow, '.', '.')
			case '@':
				wideRow = append(wideRow, '@', '.')
			}
		}
		wideGrid = append(wideGrid, wideRow)
	}

	return wideGrid
}

func simulateWide(grid [][]byte, moves string, robot *Point) {
	for _, move := range moves {
		dx, dy := getDirection(byte(move))
		nx, ny := robot.X+dx, robot.Y+dy

		if grid[ny][nx] == '#' {
			continue
		} else if grid[ny][nx] == '.' {
			robot.X, robot.Y = nx, ny
		} else if grid[ny][nx] == '[' || grid[ny][nx] == ']' {
			if canPushWide(grid, nx, ny, dx, dy) {
				pushWide(grid, nx, ny, dx, dy)
				robot.X, robot.Y = nx, ny
			}
		}
	}
}

func canPushWide(grid [][]byte, x, y, dx, dy int) bool {
	if dx != 0 { // Horizontal movement
		return canPushHorizontal(grid, x, y, dx)
	} else { // Vertical movement
		return canPushVertical(grid, x, y, dy)
	}
}

func canPushHorizontal(grid [][]byte, x, y, dx int) bool {
	// For horizontal, just check the line until we hit wall or empty space
	checkX := x
	for grid[y][checkX] == '[' || grid[y][checkX] == ']' {
		checkX += dx
	}
	return grid[y][checkX] == '.'
}

func canPushVertical(grid [][]byte, startX, startY, dy int) bool {
	// Get all boxes that would need to move
	boxesToMove := getAllBoxesToMove(grid, startX, startY, dy)

	// Check if all boxes can move
	for box := range boxesToMove {
		newY := box.Y + dy
		// Check both sides of the box
		if grid[newY][box.X] == '#' || grid[newY][box.X+1] == '#' {
			return false
		}
	}

	return true
}

func pushWide(grid [][]byte, x, y, dx, dy int) {
	if dx != 0 { // Horizontal movement
		pushHorizontal(grid, x, y, dx)
	} else { // Vertical movement
		pushVertical(grid, x, y, dy)
	}
}

func pushHorizontal(grid [][]byte, x, y, dx int) {
	// Find the end of the box chain
	endX := x
	for grid[y][endX] == '[' || grid[y][endX] == ']' {
		endX += dx
	}

	// Move boxes from end to start
	if dx > 0 { // Moving right
		for i := endX; i > x; i-- {
			grid[y][i] = grid[y][i-1]
		}
	} else { // Moving left
		for i := endX; i < x; i++ {
			grid[y][i] = grid[y][i+1]
		}
	}

	grid[y][x] = '.'
}

func pushVertical(grid [][]byte, startX, startY, dy int) {
	// Get all boxes that need to move
	allBoxes := getAllBoxesToMove(grid, startX, startY, dy)

	// Convert to slice and sort by Y coordinate
	var sortedBoxes []Point
	for box := range allBoxes {
		sortedBoxes = append(sortedBoxes, box)
	}

	// Sort: if pushing up (dy=-1), move top boxes first; if pushing down (dy=1), move bottom boxes first
	for i := range len(sortedBoxes) - 1 {
		for j := i + 1; j < len(sortedBoxes); j++ {
			if (dy < 0 && sortedBoxes[i].Y > sortedBoxes[j].Y) ||
				(dy > 0 && sortedBoxes[i].Y < sortedBoxes[j].Y) {
				sortedBoxes[i], sortedBoxes[j] = sortedBoxes[j], sortedBoxes[i]
			}
		}
	}

	// Move each box
	for _, box := range sortedBoxes {
		// Clear old position
		grid[box.Y][box.X] = '.'
		grid[box.Y][box.X+1] = '.'

		// Set new position
		newY := box.Y + dy
		grid[newY][box.X] = '['
		grid[newY][box.X+1] = ']'
	}
}

func getAllBoxesToMove(grid [][]byte, startX, startY, dy int) map[Point]bool {
	boxes := make(map[Point]bool)
	toProcess := []Point{}

	// Start with the initial box (always store box position as its left '[' coordinate)
	if grid[startY][startX] == '[' {
		toProcess = append(toProcess, Point{startX, startY})
	} else { // grid[startY][startX] == ']'
		toProcess = append(toProcess, Point{startX - 1, startY})
	}

	// Process all connected boxes using BFS
	for len(toProcess) > 0 {
		box := toProcess[0]
		toProcess = toProcess[1:]

		// Skip if already processed
		if boxes[box] {
			continue
		}
		boxes[box] = true

		// Check what's above/below this box
		newY := box.Y + dy
		leftX, rightX := box.X, box.X+1

		// Check left side
		switch grid[newY][leftX] {
		case '[':
			toProcess = append(toProcess, Point{leftX, newY})
		case ']':
			toProcess = append(toProcess, Point{leftX - 1, newY})
		}

		// Check right side
		switch grid[newY][rightX] {
		case '[':
			toProcess = append(toProcess, Point{rightX, newY})
		case ']':
			toProcess = append(toProcess, Point{rightX - 1, newY})
		}
	}

	return boxes
}

func calculateWideGPS(grid [][]byte) int {
	sum := 0
	for row := range grid {
		for col := range grid[row] {
			if grid[row][col] == '[' { // Only count left part of box
				gps := 100*row + col
				sum += gps
			}
		}
	}
	return sum
}

// ============ SHARED FUNCTIONS ============

func printGridWithRobot(grid [][]byte, robot Point) {
	for y, row := range grid {
		for x, cell := range row {
			if x == robot.X && y == robot.Y {
				fmt.Print("@")
			} else {
				fmt.Print(string(cell))
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func getDirection(move byte) (int, int) {
	switch move {
	case '^':
		return 0, -1
	case 'v':
		return 0, 1
	case '<':
		return -1, 0
	case '>':
		return 1, 0
	}
	return 0, 0
}

func parseInputFile(input string) ([][]byte, string, Point) {
	file, err := os.Open(input)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var grid [][]byte
	var moves string
	var robot Point

	// Parse grid
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		row := []byte(line)
		for x, cell := range row {
			if cell == '@' {
				robot.X, robot.Y = x, len(grid)
				row[x] = '.' // Replace robot with empty space
			}
		}
		grid = append(grid, row)
	}

	// Parse moves
	for scanner.Scan() {
		line := scanner.Text()
		moves += line
	}

	return grid, moves, robot
}
