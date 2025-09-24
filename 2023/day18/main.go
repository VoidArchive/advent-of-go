package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	x, y int64
}

type Instruction struct {
	dir   rune
	steps int64
	color string
}

func parseInput(filename string, usePart2 bool) []Instruction {
	file, _ := os.Open(filename)
	defer file.Close()
	var instructions []Instruction
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 3 {
			continue
		}
		if !usePart2 {
			dir := rune(parts[0][0])
			steps, _ := strconv.ParseInt(parts[1], 10, 64)
			color := strings.Trim(parts[2], "()")

			instructions = append(instructions, Instruction{
				dir, steps, color,
			})
		} else {
			color := strings.Trim(parts[2], "(#)")
			distanceHex := color[:5]
			steps, _ := strconv.ParseInt(distanceHex, 16, 64)
			dirCode := color[5]
			var dir rune
			switch dirCode {
			case '0':
				dir = 'R'
			case '1':
				dir = 'D'
			case '2':
				dir = 'L'
			case '3':
				dir = 'U'
			default:
				log.Fatal("Invalid direction code")
			}
			instructions = append(instructions, Instruction{
				dir, steps, color,
			})
		}
	}
	return instructions
}

func calculateArea(instructions []Instruction) int64 {
	vertices := getVertices(instructions)
	area := shoelaceArea(vertices)
	perimeter := getPerimeter(instructions)

	return area + perimeter/2 + 1
}

func getVertices(instructions []Instruction) []Point {
	vertices := []Point{{0, 0}}
	current := Point{0, 0}

	for _, inst := range instructions {
		switch inst.dir {
		case 'R':
			current.x += inst.steps
		case 'L':
			current.x -= inst.steps
		case 'U':
			current.y -= inst.steps
		case 'D':
			current.y += inst.steps
		}
		vertices = append(vertices, current)
	}
	return vertices
}

func shoelaceArea(vertices []Point) int64 {
	n := len(vertices)
	area := int64(0)

	for i := range n - 1 {
		area += vertices[i].x*vertices[i+1].y - vertices[i+1].x*vertices[i].y
	}
	if area < 0 {
		area = -area
	}
	return area / 2
}

func getPerimeter(instructions []Instruction) int64 {
	perimeter := int64(0)
	for _, inst := range instructions {
		perimeter += inst.steps
	}
	return perimeter
}

func main() {
	instructions1 := parseInput("input.txt", false)
	area1 := calculateArea(instructions1)
	fmt.Printf("Part 1: %d\n", area1)

	instructions2 := parseInput("input.txt", true)
	area2 := calculateArea(instructions2)
	fmt.Printf("Part 2: %d\n", area2)
}
