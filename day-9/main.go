package main

import (
	"bufio"
	"fmt"
	"os"
)

type Block byte

func parseDiskMap(input string) []Block {
	var layout []Block
	fileID := 0

	for i := range len(input) {
		size := int(input[i] - '0')

		if i%2 == 0 {
			char := Block('0' + fileID)
			for range size {
				layout = append(layout, char)
			}
			fileID++
		} else {
			for range size {
				layout = append(layout, '.')
			}
		}
	}
	return layout
}

func compact(layout []Block) {
	for {
		leftmostFreeIndex := -1
		for i := range len(layout) {
			if layout[i] == '.' {
				leftmostFreeIndex = i
				break
			}
		}
		if leftmostFreeIndex == -1 {
			break
		}
		rightmostFileIndex := -1
		for i := len(layout) - 1; i >= 0; i-- {
			if layout[i] != '.' {
				rightmostFileIndex = i
				break
			}
		}
		if rightmostFileIndex == -1 || rightmostFileIndex <= leftmostFreeIndex {
			break
		}
		layout[leftmostFreeIndex] = layout[rightmostFileIndex]
		layout[rightmostFileIndex] = '.'
	}
}

func checksum(layout []Block) int {
	sum := 0
	for i, block := range layout {
		if block != '.' {
			sum += i * int(block-'0')
		}
	}
	return sum
}

func solve(input string) int {
	layout := parseDiskMap(input)
	compact(layout)
	return checksum(layout)
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		input := scanner.Text()
		layout := parseDiskMap(input)
		fmt.Println("Before compact: ", string(layout))
		compact(layout)
		fmt.Println("After compact:  ", string(layout))
		fmt.Println("Part 1 Checksum:", solve(input))
	}

	exampleInput := "2333133121414131402"
	layout := parseDiskMap(exampleInput)
	fmt.Println("Before compact: ", string(layout))
	compact(layout)
	fmt.Println("After compact:  ", string(layout))
	fmt.Println("Checksum of example input:", checksum(layout))
}
