package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"strings"
)

func processOneBlink(currentStones []*big.Int) []*big.Int {
	nextStones := make([]*big.Int, 0, len(currentStones))

	bigZero := big.NewInt(0)
	bigOne := big.NewInt(1)
	multipiler := big.NewInt(2024)

	for _, stoneValue := range currentStones {
		if stoneValue.Cmp(bigZero) == 0 {
			nextStones = append(nextStones, new(big.Int).Set(bigOne))
			continue
		}
		stoneStr := stoneValue.String()
		if len(stoneStr)%2 == 0 {
			midpoint := len(stoneStr) / 2
			leftStr := stoneStr[:midpoint]
			rightStr := stoneStr[midpoint:]

			leftNum, _ := new(big.Int).SetString(leftStr, 10)
			rightNum, _ := new(big.Int).SetString(rightStr, 10)

			nextStones = append(nextStones, leftNum)
			nextStones = append(nextStones, rightNum)
			continue
		}
		newValue := new(big.Int)
		newValue.Mul(stoneValue, multipiler)
		nextStones = append(nextStones, newValue)
	}
	return nextStones
}

func readInput(input string) []*big.Int {
	file, err := os.Open(input)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	stones := []*big.Int{}
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		for _, part := range parts {
			num, ok := new(big.Int).SetString(part, 10)
			if !ok {
				fmt.Printf("Error converting %s to a number %v\n ", part, err)
				os.Exit(1)
			}
			stones = append(stones, num)
		}

	}
	return stones
}

func main() {
	stones := readInput("input.txt")
	numberOfBlinks := 75

	for blink := 1; blink <= numberOfBlinks; blink++ {
		stones = processOneBlink(stones)
		if blink%5 == 0 || blink == 1 || blink == numberOfBlinks {
			fmt.Printf("After %d blink(s): %d stones\n", blink, len(stones))
		}
	}
	fmt.Println("Part 1:", len(stones))
}
