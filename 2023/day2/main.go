package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type GameReveal struct {
	Red, Green, Blue int
}

func parseGame(line string) (int, []GameReveal) {
	parts := strings.Split(line, ": ")
	gameID, _ := strconv.Atoi(strings.TrimPrefix(parts[0], "Game "))

	var reveals []GameReveal
	revealStrings := strings.Split(parts[1], "; ")

	for _, revealStr := range revealStrings {
		reveal := GameReveal{}
		cubes := strings.Split(revealStr, ", ")

		for _, cube := range cubes {
			cubeParts := strings.Split(strings.TrimSpace(cube), " ")
			count, _ := strconv.Atoi(cubeParts[0])
			color := cubeParts[1]

			switch color {
			case "red":
				reveal.Red = count
			case "green":
				reveal.Green = count
			case "blue":
				reveal.Blue = count
			}
		}
		reveals = append(reveals, reveal)
	}
	return gameID, reveals
}

func isGamePossible(reveals []GameReveal, maxRed, maxGreen, maxBlue int) bool {
	for _, reveal := range reveals {
		if reveal.Red > maxRed || reveal.Green > maxGreen || reveal.Blue > maxBlue {
			return false
		}
	}
	return true
}

func findMinimumCubes(reveals []GameReveal) (int, int, int) {
	maxRed, maxGreen, maxBlue := 0, 0, 0

	for _, reveal := range reveals {
		if reveal.Red > maxRed {
			maxRed = reveal.Red
		}
		if reveal.Green > maxGreen {
			maxGreen = reveal.Green
		}
		if reveal.Blue > maxBlue {
			maxBlue = reveal.Blue
		}
	}
	return maxRed, maxGreen, maxBlue
}

func main() {
	file, _ := os.Open("input.txt")
	defer file.Close()

	scanner := bufio.NewScanner(file)
	sum := 0
	totalPower := 0

	for scanner.Scan() {
		gameID, reveals := parseGame(scanner.Text())

		if isGamePossible(reveals, 12, 13, 14) {
			sum += gameID
		}
		minRed, minGreen, minBlue := findMinimumCubes(reveals)
		power := minRed * minGreen * minBlue
		totalPower += power

	}
	fmt.Printf("Part 1: %d\n", sum)
	fmt.Printf("Part 2: %d\n", totalPower)
}
