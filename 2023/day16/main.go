package main

import (
	"fmt"
	"os"
	"strings"
)

type Position struct {
	row, col int
}

type Direction struct {
	dr, dc int
}

type Beam struct {
	pos Position
	dir Direction
}

var (
	Right = Direction{0, 1}
	Left  = Direction{0, -1}
	Up    = Direction{-1, 0}
	Down  = Direction{1, 0}
)

func simulateBeam(grid []string, startPos Position, startDir Direction) int {
	rows, cols := len(grid), len(grid[0])

	// Track visited states to avoid infinite loops
	visited := make(map[string]bool)
	// Track energized tiles
	energized := make(map[Position]bool)

	// BFS/DFS to simulate beam propagation
	var beams []Beam
	beams = append(beams, Beam{startPos, startDir})

	for len(beams) > 0 {
		beam := beams[0]
		beams = beams[1:]

		// Check bounds
		if beam.pos.row < 0 || beam.pos.row >= rows ||
			beam.pos.col < 0 || beam.pos.col >= cols {
			continue
		}

		// Create state key to detect cycles
		stateKey := fmt.Sprintf("%d,%d,%d,%d",
			beam.pos.row, beam.pos.col, beam.dir.dr, beam.dir.dc)

		if visited[stateKey] {
			continue
		}
		visited[stateKey] = true

		// Mark tile as energized
		energized[beam.pos] = true

		// Get current tile
		tile := grid[beam.pos.row][beam.pos.col]

		switch tile {
		case '.':
			// Continue in same direction
			newPos := Position{
				beam.pos.row + beam.dir.dr,
				beam.pos.col + beam.dir.dc,
			}
			beams = append(beams, Beam{newPos, beam.dir})

		case '/':
			// Reflect: right->up, left->down, up->right, down->left
			var newDir Direction
			if beam.dir == Right {
				newDir = Up
			} else if beam.dir == Left {
				newDir = Down
			} else if beam.dir == Up {
				newDir = Right
			} else { // Down
				newDir = Left
			}
			newPos := Position{
				beam.pos.row + newDir.dr,
				beam.pos.col + newDir.dc,
			}
			beams = append(beams, Beam{newPos, newDir})

		case '\\':
			// Reflect: right->down, left->up, up->left, down->right
			var newDir Direction
			if beam.dir == Right {
				newDir = Down
			} else if beam.dir == Left {
				newDir = Up
			} else if beam.dir == Up {
				newDir = Left
			} else { // Down
				newDir = Right
			}
			newPos := Position{
				beam.pos.row + newDir.dr,
				beam.pos.col + newDir.dc,
			}
			beams = append(beams, Beam{newPos, newDir})

		case '|':
			// Vertical splitter
			if beam.dir == Up || beam.dir == Down {
				// Pointy end - pass through
				newPos := Position{
					beam.pos.row + beam.dir.dr,
					beam.pos.col + beam.dir.dc,
				}
				beams = append(beams, Beam{newPos, beam.dir})
			} else {
				// Flat side - split into up and down
				upPos := Position{
					beam.pos.row + Up.dr,
					beam.pos.col + Up.dc,
				}
				downPos := Position{
					beam.pos.row + Down.dr,
					beam.pos.col + Down.dc,
				}
				beams = append(beams, Beam{upPos, Up})
				beams = append(beams, Beam{downPos, Down})
			}

		case '-':
			// Horizontal splitter
			if beam.dir == Left || beam.dir == Right {
				// Pointy end - pass through
				newPos := Position{
					beam.pos.row + beam.dir.dr,
					beam.pos.col + beam.dir.dc,
				}
				beams = append(beams, Beam{newPos, beam.dir})
			} else {
				// Flat side - split into left and right
				leftPos := Position{
					beam.pos.row + Left.dr,
					beam.pos.col + Left.dc,
				}
				rightPos := Position{
					beam.pos.row + Right.dr,
					beam.pos.col + Right.dc,
				}
				beams = append(beams, Beam{leftPos, Left})
				beams = append(beams, Beam{rightPos, Right})
			}
		}
	}

	return len(energized)
}

func solvePart1(grid []string) int {
	return simulateBeam(grid, Position{0, 0}, Right)
}

func solvePart2(grid []string) int {
	rows, cols := len(grid), len(grid[0])
	maxEnergized := 0

	// Test all starting positions along edges

	// Top row (heading down)
	for col := 0; col < cols; col++ {
		energized := simulateBeam(grid, Position{0, col}, Down)
		if energized > maxEnergized {
			maxEnergized = energized
		}
	}

	// Bottom row (heading up)
	for col := 0; col < cols; col++ {
		energized := simulateBeam(grid, Position{rows - 1, col}, Up)
		if energized > maxEnergized {
			maxEnergized = energized
		}
	}

	// Left column (heading right)
	for row := 0; row < rows; row++ {
		energized := simulateBeam(grid, Position{row, 0}, Right)
		if energized > maxEnergized {
			maxEnergized = energized
		}
	}

	// Right column (heading left)
	for row := 0; row < rows; row++ {
		energized := simulateBeam(grid, Position{row, cols - 1}, Left)
		if energized > maxEnergized {
			maxEnergized = energized
		}
	}

	return maxEnergized
}

func main() {
	// Read input from file
	content, err := os.ReadFile("input.txt")
	if err != nil {
		fmt.Printf("Error reading input.txt: %v\n", err)
		return
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")

	// Solve both parts
	part1 := solvePart1(lines)
	fmt.Printf("Part 1: %d\n", part1)

	part2 := solvePart2(lines)
	fmt.Printf("Part 2: %d\n", part2)
}
