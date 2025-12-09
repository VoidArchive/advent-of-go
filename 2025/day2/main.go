package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func isInvalid(n int) bool {
	s := strconv.Itoa(n)
	if len(s)%2 != 0 {
		return false
	}
	half := len(s) / 2
	return s[:half] == s[half:]
}

func isInvalidPart2(n int) bool {
	s := strconv.Itoa(n)
	length := len(s)
	for patLen := 1; patLen <= length/2; patLen++ {
		if length%patLen != 0 {
			continue
		}
		pat := s[:patLen]
		valid := true

		for i := patLen; i < length; i += patLen {
			if s[i:i+patLen] != pat {
				valid = false
				break
			}
		}
		if valid {
			return true
		}
	}
	return false
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}

	defer file.Close()

	var ranges []string
	var sum1 int
	var sum2 int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		ranges = strings.Split(line, ",")
	}

	for _, r := range ranges {
		part := strings.Split(r, "-")
		r1, _ := strconv.Atoi(part[0])
		r2, _ := strconv.Atoi(part[1])

		for n := r1; n <= r2; n++ {
			if isInvalid(n) {
				sum1 += n
			}
			if isInvalidPart2(n) {
				sum2 += n
			}
		}
	}

	fmt.Printf("Part 1: %d\n", sum1)
	fmt.Printf("Part 2: %d\n", sum2)
}
