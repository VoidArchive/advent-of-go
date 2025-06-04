package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"strings"
)

var (
	memo       map[string]map[int]*big.Int
	bigZero    *big.Int
	bigOne     *big.Int
	multiplier *big.Int
)

func initGlobal() {
	memo = make(map[string]map[int]*big.Int)
	bigZero = big.NewInt(0)
	bigOne = big.NewInt(1)
	multiplier = big.NewInt(2024)
}

func countDescendantStones(currentStoneValue *big.Int, blinksRemaining int) *big.Int {
	if blinksRemaining == 0 {
		return new(big.Int).Set(bigOne)
	}
	currentStoneStr := currentStoneValue.String()
	if kMap, ok := memo[currentStoneStr]; ok {
		if count, ok2 := kMap[blinksRemaining]; ok2 {
			return new(big.Int).Set(count)
		}
	}
	calculateCount := big.NewInt(0)
	if currentStoneValue.Cmp(bigZero) == 0 {
		calculateCount.Add(calculateCount, countDescendantStones(
			bigOne,
			blinksRemaining-1,
		))
	} else {
		if len(currentStoneStr)%2 == 0 {
			midpoint := len(currentStoneStr) / 2
			leftStr := currentStoneStr[:midpoint]
			rightStr := currentStoneStr[midpoint:]

			leftNum, _ := new(big.Int).SetString(leftStr, 10)
			rightNum, _ := new(big.Int).SetString(rightStr, 10)

			calculateCount.Add(calculateCount, countDescendantStones(leftNum, blinksRemaining-1))
			calculateCount.Add(calculateCount, countDescendantStones(rightNum, blinksRemaining-1))
		} else {
			nextValue := new(big.Int).Mul(currentStoneValue, multiplier)
			calculateCount.Add(calculateCount, countDescendantStones(nextValue, blinksRemaining-1))
		}
	}
	if _, ok := memo[currentStoneStr]; !ok {
		memo[currentStoneStr] = make(map[int]*big.Int)
	}
	memo[currentStoneStr][blinksRemaining] = new(big.Int).Set(calculateCount)
	return calculateCount
}

// func processOneBlink(currentStones []*big.Int) []*big.Int {
// 	nextStones := make([]*big.Int, 0, len(currentStones))
//
// 	for _, stoneValue := range currentStones {
// 		if stoneValue.Cmp(bigZero) == 0 {
// 			nextStones = append(nextStones, new(big.Int).Set(bigOne))
// 			continue
// 		}
// 		stoneStr := stoneValue.String()
// 		if len(stoneStr)%2 == 0 {
// 			midpoint := len(stoneStr) / 2
// 			leftStr := stoneStr[:midpoint]
// 			rightStr := stoneStr[midpoint:]
//
// 			leftNum, _ := new(big.Int).SetString(leftStr, 10)
// 			rightNum, _ := new(big.Int).SetString(rightStr, 10)
//
// 			nextStones = append(nextStones, leftNum)
// 			nextStones = append(nextStones, rightNum)
// 			continue
// 		}
// 		newValue := new(big.Int)
// 		newValue.Mul(stoneValue, multipiler)
// 		nextStones = append(nextStones, newValue)
// 	}
// 	return nextStones
// }

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
	initGlobal()
	stones := readInput("input.txt")
	numberOfBlinks := 75
	totalStonesCount := big.NewInt(0)

	for i, stone := range stones {
		fmt.Printf("Processing initial stone %d/%d: %s\n",
			i+1, len(stones), stone.String())
		countStone := countDescendantStones(stone, numberOfBlinks)
		totalStonesCount.Add(totalStonesCount, countStone)
		fmt.Printf("  Stone %s will become %s stones.\n", stone.String(), countStone.String())
	}
	fmt.Printf("After %d blinks, the total number of stones will be: %s\n",
		numberOfBlinks, totalStonesCount.String())
}
