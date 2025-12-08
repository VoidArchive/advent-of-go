package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func countZeroCrossings(start, distance int, left bool) (newPos, zeros int) {
	if distance == 0 {
		return start, 0
	}

	if left {
		newPos = ((start-distance)%100 + 100) % 100
		if start == 0 {
			zeros = distance / 100
		} else if distance >= start {
			zeros = 1 + (distance-start)/100
		}
	} else {
		newPos = (start + distance) % 100
		stepsToZero := (100 - start) % 100
		if stepsToZero == 0 {
			stepsToZero = 100
		}
		if distance >= stepsToZero {
			zeros = 1 + (distance-stepsToZero)/100
		}
	}
	return newPos, zeros
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	pos := 50
	part1 := 0
	part2 := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		direction := line[0]
		distance, _ := strconv.Atoi(line[1:])

		newPos, zeros := countZeroCrossings(pos, distance, direction == 'L')
		pos = newPos
		part2 += zeros

		if pos == 0 {
			part1++
		}
	}

	fmt.Printf("Part 1: %d\n", part1)
	fmt.Printf("Part 2: %d\n", part2)
}
