package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Time:        61     67     75     71
// Distance:   430   1036   1307   1150
//

func parseLine(line string) []int {
	parts := strings.Fields(line)[1:]
	var nums []int
	for _, part := range parts {
		n, _ := strconv.Atoi(part)
		nums = append(nums, n)
	}
	return nums
}

func parseLinetoSingleInt(line string) int {
	parts := strings.Fields(line)[1:]
	combined := strings.Join(parts, "")
	val, _ := strconv.Atoi(combined)
	return val
}

func waysToWin(time, record int) int {
	count := 0
	for hold := range time {
		speed := hold
		moveTime := time - hold
		distance := speed * moveTime
		if distance > record {
			count++
		}
	}
	return count
}

func main() {
	input := `Time:        61     67     75     71
Distance:   430   1036   1307   1150`
	lines := strings.Split(input, "\n")
	times := parseLine(lines[0])
	distances := parseLine(lines[1])

	part1Result := 1
	for i := range len(times) {
		ways := waysToWin(times[i], distances[i])
		part1Result *= ways
	}
	fmt.Println(part1Result)

	part2Time := parseLinetoSingleInt(lines[0])
	part2Dist := parseLinetoSingleInt(lines[1])
	part2Result := waysToWin(part2Time, part2Dist)
	fmt.Println(part2Result)
}
