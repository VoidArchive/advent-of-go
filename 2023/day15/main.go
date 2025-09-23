package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Lens struct {
	Label       string
	FocalLength int
}

func hash(s string) int {
	current := 0
	for _, char := range s {
		ascii := int(char)
		current += ascii
		current *= 17
		current %= 256
	}
	return current
}

func part1(steps []string) int {
	sum := 0
	for _, step := range steps {
		hashValue := hash(step)
		sum += hashValue
	}
	return sum
}

func part2(steps []string) int {
	boxes := make([][]Lens, 256)
	for i := range boxes {
		boxes[i] = make([]Lens, 0)
	}

	for _, step := range steps {
		if strings.Contains(step, "=") {
			parts := strings.Split(step, "=")
			label := parts[0]
			focalLength, _ := strconv.Atoi(parts[1])
			boxNum := hash(label)

			found := false
			for i, lens := range boxes[boxNum] {
				if lens.Label == label {
					boxes[boxNum][i].FocalLength = focalLength
					found = true
					break
				}
			}
			if !found {
				boxes[boxNum] = append(boxes[boxNum], Lens{label, focalLength})
			}
		} else if strings.Contains(step, "-") {
			label := strings.TrimSuffix(step, "-")
			boxNum := hash(label)

			for i, lens := range boxes[boxNum] {
				if lens.Label == label {
					boxes[boxNum] = append(boxes[boxNum][:i], boxes[boxNum][i+1:]...)
					break
				}
			}
		}
	}
	totalPower := 0
	for boxNum, box := range boxes {
		for slotNum, lens := range box {
			power := (boxNum + 1) * (slotNum + 1) * lens.FocalLength
			totalPower += power
		}
	}
	return totalPower
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var input string

	for scanner.Scan() {
		input += scanner.Text()
	}
	steps := strings.Split(input, ",")
	fmt.Printf("Part 1 - Sum of hash values: %d\n", part1(steps))

	// Part 2
	fmt.Printf("Part 2 - Total focusing power: %d\n", part2(steps))
}
