package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	input := strings.TrimSpace(string(data))
	schematics := strings.Split(input, "\n\n")

	var locks [][]int
	var keys [][]int

	for _, schematic := range schematics {
		lines := strings.Split(schematic, "\n")

		// Check if it's a lock (top row all #, bottom row all .)
		if strings.Count(lines[0], "#") == 5 && strings.Count(lines[len(lines)-1], ".") == 5 {
			// It's a lock - convert to heights
			heights := make([]int, 5)
			for col := range 5 {
				for row := 1; row < len(lines); row++ {
					if lines[row][col] == '#' {
						heights[col]++
					} else {
						break
					}
				}
			}
			locks = append(locks, heights)
		} else if strings.Count(lines[0], ".") == 5 && strings.Count(lines[len(lines)-1], "#") == 5 {
			// It's a key - convert to heights
			heights := make([]int, 5)
			for col := range 5 {
				for row := len(lines) - 2; row >= 0; row-- {
					if lines[row][col] == '#' {
						heights[col]++
					} else {
						break
					}
				}
			}
			keys = append(keys, heights)
		}
	}

	// Count valid lock/key pairs
	validPairs := 0
	maxHeight := 5 // Available space is 5 (7 rows - 2 for top/bottom)

	for _, lock := range locks {
		for _, key := range keys {
			valid := true
			for i := range 5 {
				if lock[i]+key[i] > maxHeight {
					valid = false
					break
				}
			}
			if valid {
				validPairs++
			}
		}
	}

	fmt.Println(validPairs)
}

