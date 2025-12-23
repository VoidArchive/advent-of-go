package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Range struct {
	Start int
	End   int
}

func main() {
	file, _ := os.Open("input.txt")
	defer file.Close()
	scanner := bufio.NewScanner(file)

	var freshRanges []Range
	var ingredientIDs []int
	parsingRange := true

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			parsingRange = false
			continue
		}

		if parsingRange {

			parts := strings.Split(line, "-")
			if len(parts) != 2 {
				continue
			}
			x, _ := strconv.Atoi(parts[0])
			y, _ := strconv.Atoi(parts[1])
			freshRanges = append(freshRanges, Range{x, y})
		} else {

			ingredient, _ := strconv.Atoi(line)
			ingredientIDs = append(ingredientIDs, ingredient)
		}
	}

	part1 := 0
	for _, id := range ingredientIDs {
		for _, r := range freshRanges {
			if id >= r.Start && id <= r.End {
				part1++
				break
			}
		}
	}
	fmt.Printf("Part 1: %v\n", part1)

	sort.Slice(freshRanges, func(i, j int) bool {
		return freshRanges[i].Start < freshRanges[j].Start
	})

	merged := []Range{freshRanges[0]}
	for _, r := range freshRanges[1:] {
		last := &merged[len(merged)-1]
		if r.Start <= last.End+1 {
			last.End = max(last.End, r.End)
		} else {
			merged = append(merged, r)
		}
	}

	part2 := 0
	for _, r := range merged {
		part2 += r.End - r.Start + 1
	}

	fmt.Printf("Part 2: %v\n", part2)
}
