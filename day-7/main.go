package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func parseLine(line string) (int, []int) {
	l := strings.Split(line, ": ")
	target, _ := strconv.Atoi(l[0])
	snums := strings.Split(l[1], " ")
	nums := make([]int, 0, len(snums))
	for _, n := range snums {
		num, err := strconv.Atoi(n)
		if err != nil {
			log.Fatal(err)
		}
		nums = append(nums, num)
	}
	return target, nums
}

func validCombinationExists(target int, nums []int) bool {
	n := len(nums) - 1
	max := 1 << n

	for i := range max {
		result := nums[0]
		for j := range n {
			if (i>>j)&1 == 0 {
				result += nums[j+1]
			} else {
				result *= nums[j+1]
			}
		}
		if result == target {
			return true
		}
	}
	return false
}

func concat(a, b int) int {
	multiplier := 1
	for temp := b; temp > 0; temp /= 10 {
		multiplier *= 10
	}
	return a*multiplier + b
}

func validCombinationExistsPart2(target int, nums []int) bool {
	n := len(nums) - 1
	max := int(math.Pow(3, float64(n)))

	for i := range max {
		ops := make([]int, n)
		x := i
		for j := range n {
			ops[j] = x % 3
			x /= 3
		}

		result := nums[0]
		for j := range n {
			switch ops[j] {
			case 0:
				result += nums[j+1]
			case 1:
				result *= nums[j+1]
			case 2:
				result = concat(result, nums[j+1])
			}
		}
		if result == target {
			return true
		}
	}
	return false
}

func readLines(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return lines
}

func main() {
	lines := readLines("input.txt")
	part1Result := 0
	part2Result := 0
	for _, line := range lines {
		target, nums := parseLine(line)
		if validCombinationExists(target, nums) {
			part1Result += target
		}

		if validCombinationExistsPart2(target, nums) {
			part2Result += target
		}
	}
	fmt.Println("Total valid combinations:", part1Result)
	fmt.Println("Total valid combinations for part 2:", part2Result)
}
