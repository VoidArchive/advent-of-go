package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Point struct {
	x, y, z int
}

type Brick struct {
	start, end Point
	id         int
}

func parseBrick(line string, id int) Brick {
	parts := strings.Split(line, "~")
	start := parsePoint(parts[0])
	end := parsePoint(parts[1])
	return Brick{start, end, id}
}

func parsePoint(s string) Point {
	coords := strings.Split(s, ",")
	x, _ := strconv.Atoi(coords[0])
	y, _ := strconv.Atoi(coords[1])
	z, _ := strconv.Atoi(coords[2])
	return Point{x, y, z}
}

func (b Brick) getBlocks() []Point {
	blocks := []Point{}
	for x := min(b.start.x, b.end.x); x <= max(b.start.x, b.end.x); x++ {
		for y := min(b.start.y, b.end.y); y <= max(b.start.y, b.end.y); y++ {
			for z := min(b.start.z, b.end.z); z <= max(b.start.z, b.end.z); z++ {
				blocks = append(blocks, Point{x, y, z})
			}
		}
	}
	return blocks
}

func (b Brick) minZ() int {
	return min(b.start.z, b.end.z)
}

func settleBricks(bricks []Brick) ([]Brick, map[Point]int) {
	sort.Slice(bricks, func(i, j int) bool {
		return bricks[i].minZ() < bricks[j].minZ()
	})

	occupied := make(map[Point]int)
	settled := make([]Brick, len(bricks))

	for i, brick := range bricks {
		blocks := brick.getBlocks()

		maxFallZ := 1
		for _, block := range blocks {
			for z := block.z - 1; z >= 1; z-- {
				checkPoint := Point{block.x, block.y, z}
				if _, exists := occupied[checkPoint]; exists {
					maxFallZ = max(maxFallZ, z+1)
					break
				}
			}
		}

		drop := brick.minZ() - maxFallZ
		newBrick := Brick{
			start: Point{brick.start.x, brick.start.y, brick.start.z - drop},
			end:   Point{brick.end.x, brick.end.y, brick.end.z - drop},
			id:    brick.id,
		}

		for _, block := range newBrick.getBlocks() {
			occupied[block] = newBrick.id
		}

		settled[i] = newBrick
	}

	return settled, occupied
}

func buildSupportGraph(bricks []Brick, occupied map[Point]int) (map[int][]int, map[int][]int) {
	supports := make(map[int][]int)
	supportedBy := make(map[int][]int)

	for _, brick := range bricks {
		supportedBySet := make(map[int]bool)

		for _, block := range brick.getBlocks() {
			if block.z == brick.minZ() {
				below := Point{block.x, block.y, block.z - 1}
				if belowID, exists := occupied[below]; exists && belowID != brick.id {
					supportedBySet[belowID] = true
				}
			}
		}

		for id := range supportedBySet {
			supportedBy[brick.id] = append(supportedBy[brick.id], id)
			supports[id] = append(supports[id], brick.id)
		}
	}

	return supports, supportedBy
}

func canDisintegrate(brickID int, supports map[int][]int, supportedBy map[int][]int) bool {
	for _, supportedID := range supports[brickID] {
		if len(supportedBy[supportedID]) == 1 {
			return false
		}
	}
	return true
}

func countChainReaction(brickID int, supports map[int][]int, supportedBy map[int][]int) int {
	fallen := make(map[int]bool)
	fallen[brickID] = true

	queue := []int{brickID}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, supportedID := range supports[current] {
			if fallen[supportedID] {
				continue
			}

			allSupportsFallen := true
			for _, supportID := range supportedBy[supportedID] {
				if !fallen[supportID] {
					allSupportsFallen = false
					break
				}
			}

			if allSupportsFallen {
				fallen[supportedID] = true
				queue = append(queue, supportedID)
			}
		}
	}

	return len(fallen) - 1
}

func main() {
	file, _ := os.Open("input.txt")
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var bricks []Brick
	id := 0

	for scanner.Scan() {
		bricks = append(bricks, parseBrick(scanner.Text(), id))
		id++
	}

	settled, occupied := settleBricks(bricks)
	supports, supportedBy := buildSupportGraph(settled, occupied)

	safeCount := 0
	totalFallen := 0

	for _, brick := range settled {
		if canDisintegrate(brick.id, supports, supportedBy) {
			safeCount++
		}
		totalFallen += countChainReaction(brick.id, supports, supportedBy)
	}

	fmt.Println("Part 1:", safeCount)
	fmt.Println("Part 2:", totalFallen)
}
