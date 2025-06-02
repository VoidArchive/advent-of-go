package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Block int

const EmptyBlock = Block('.')

func calculateChecksum(layout []Block) int64 {
	var sum int64 = 0
	for i, blockValue := range layout {
		if blockValue != EmptyBlock {
			numericFileID := int(blockValue - '0')
			sum += int64(i) * int64(numericFileID)
		}
	}
	return sum
}

func formatLayout(layout []Block) string {
	sb := strings.Builder{}
	for _, block := range layout {
		if block == EmptyBlock {
			sb.WriteRune('.')
		} else {
			sb.WriteString(strconv.Itoa(int(block - '0')))
		}
	}
	return sb.String()
}

// --- INFO: Part 1 ---

func parseDiskMapPart1(input string) []Block {
	var layout []Block
	numericFileID := 0

	for i := range len(input) {
		size := int(input[i] - '0')

		if i%2 == 0 {
			blockRepresentation := Block('0' + numericFileID)
			for range size {
				layout = append(layout, blockRepresentation)
			}
			numericFileID++
		} else {
			for range size {
				layout = append(layout, EmptyBlock)
			}
		}
	}
	return layout
}

func compactPart1(layout []Block) {
	for {
		leftmostFreeIndex := -1
		for i := range len(layout) {
			if layout[i] == EmptyBlock {
				leftmostFreeIndex = i
				break
			}
		}
		if leftmostFreeIndex == -1 {
			break
		}
		rightmostFileBlockIndex := -1
		for i := len(layout) - 1; i >= 0; i-- {
			if layout[i] != EmptyBlock {
				rightmostFileBlockIndex = i
				break
			}
		}
		if rightmostFileBlockIndex == -1 || rightmostFileBlockIndex <= leftmostFreeIndex {
			break
		}
		layout[leftmostFreeIndex] = layout[rightmostFileBlockIndex]
		layout[rightmostFileBlockIndex] = EmptyBlock
	}
}

func solvePart1(input string) int64 {
	layout := parseDiskMapPart1(input)
	compactPart1(layout)
	return calculateChecksum(layout)
}

// --- INFO: Part 2 ---

func parseDiskMapPart2(input string) (layout []Block, fileSizes map[int]int, maxFileID int) {
	layout = make([]Block, 0)
	fileSizes = make(map[int]int)
	maxFileID = -1
	currentFileID := 0

	for i := range len(input) {
		size := int(input[i] - '0')

		if i%2 == 0 {
			fileSizes[currentFileID] = size
			if currentFileID > maxFileID {
				maxFileID = currentFileID
			}
			blockRepresentation := Block('0' + currentFileID)
			if size > 0 {
				for range size {
					layout = append(layout, blockRepresentation)
				}
			}
			currentFileID++
		} else {
			if size > 0 {
				for range size {
					layout = append(layout, EmptyBlock)
				}
			}
		}
	}
	return layout, fileSizes, maxFileID
}

func compactPart2(layout []Block, fileSizes map[int]int, maxFileID int) {
	for fileIDToMove := maxFileID; fileIDToMove >= 0; fileIDToMove-- {
		sizeOfFileToMove, fileExistsInMap := fileSizes[fileIDToMove]

		if !fileExistsInMap || sizeOfFileToMove == 0 {
			continue
		}
		blockRepresentation := Block('0' + fileIDToMove)

		currentFileStartIndex := -1
		for i := 0; i <= len(layout)-sizeOfFileToMove; i++ {
			if layout[i] == blockRepresentation {
				isThisTheFile := true
				for k := 1; k < sizeOfFileToMove; k++ {
					if layout[i+k] != blockRepresentation {
						isThisTheFile = false
						break
					}
				}
				if isThisTheFile {
					currentFileStartIndex = i
					break
				}
			}
		}

		if currentFileStartIndex == -1 {
			continue
		}

		bestTargetSlotStartIndex := -1
		for potentialSlotStart := 0; potentialSlotStart <= currentFileStartIndex-sizeOfFileToMove; potentialSlotStart++ {
			isSpanFreeAndFits := true
			for k := range sizeOfFileToMove {
				if layout[potentialSlotStart+k] != EmptyBlock {
					isSpanFreeAndFits = false
					break
				}
			}
			if isSpanFreeAndFits {
				bestTargetSlotStartIndex = potentialSlotStart
				break
			}
		}

		if bestTargetSlotStartIndex != -1 {
			for k := range sizeOfFileToMove {
				layout[currentFileStartIndex+k] = EmptyBlock
				layout[bestTargetSlotStartIndex+k] = blockRepresentation
			}
		}
	}
}

func solvePart2(input string) int64 { // Return int64
	layout, fileSizes, maxFileID := parseDiskMapPart2(input)
	compactPart2(layout, fileSizes, maxFileID)
	return calculateChecksum(layout)
}

func main() {
	exampleInput := "2333133121414131402"

	fmt.Println("--- Example Part 1 ---")
	fmt.Println("Input:", exampleInput)
	layoutP1Ex := parseDiskMapPart1(exampleInput)
	fmt.Println("Initial Layout Part 1:", formatLayout(layoutP1Ex))
	compactPart1(layoutP1Ex)
	fmt.Println("Final Layout Part 1:  ", formatLayout(layoutP1Ex))
	fmt.Println("Checksum Example Part 1:", calculateChecksum(layoutP1Ex)) // Expected: 1928

	fmt.Println("\n--- Example Part 2 ---")
	fmt.Println("Input:", exampleInput)
	layoutP2Ex, fileSizesP2Ex, maxFileIDP2Ex := parseDiskMapPart2(exampleInput)
	fmt.Println("Initial Layout Part 2:", formatLayout(layoutP2Ex))
	compactPart2(layoutP2Ex, fileSizesP2Ex, maxFileIDP2Ex)
	fmt.Println("Final Layout Part 2:  ", formatLayout(layoutP2Ex))
	fmt.Println("Checksum Example Part 2:", calculateChecksum(layoutP2Ex)) // Expected: 2858

	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("\nError opening input.txt:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		puzzleInput := scanner.Text()
		fmt.Println("\n--- Puzzle Input Solutions ---")

		fmt.Println("Solving Part 1...")
		part1Solution := solvePart1(puzzleInput)
		fmt.Println("Part 1 Checksum:", part1Solution)

		fmt.Println("\nSolving Part 2...")
		part2Solution := solvePart2(puzzleInput)
		fmt.Println("Part 2 Checksum:", part2Solution)
	} else {
		fmt.Println("\nCould not read from input.txt or file is empty.")
	}
}
