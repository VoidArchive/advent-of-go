package main

import (
	"bufio"
	"fmt"
	"os"
)

func maxJoltagePart1(bank string) int {
	n := len(bank)
	maxVal := 0

	for i := range n {
		for j := i + 1; j < n; j++ {
			val := int(bank[i]-'0')*10 + int(bank[j]-'0')
			if val > maxVal {
				maxVal = val
			}
		}
	}

	return maxVal
}

func maxJoltagePart2(bank string) uint64 {
	n := len(bank)
	k := 12

	// dp[i][j] stores the maximum number formed by picking j digits from first i characters
	// We store as string to handle large numbers, then convert at the end
	dp := make([][]string, n+1)
	for i := range dp {
		dp[i] = make([]string, k+1)
		for j := range dp[i] {
			dp[i][j] = ""
		}
	}

	dp[0][0] = "0"

	for i := 1; i <= n; i++ {
		digit := string(bank[i-1])
		for j := 0; j <= k && j <= i; j++ {
			// Option 1: don't pick current digit
			if dp[i-1][j] != "" {
				if dp[i][j] == "" || compareNumStr(dp[i-1][j], dp[i][j]) > 0 {
					dp[i][j] = dp[i-1][j]
				}
			}

			// Option 2: pick current digit
			if j > 0 && dp[i-1][j-1] != "" {
				var newVal string
				if dp[i-1][j-1] == "0" {
					newVal = digit
				} else {
					newVal = dp[i-1][j-1] + digit
				}
				if dp[i][j] == "" || compareNumStr(newVal, dp[i][j]) > 0 {
					dp[i][j] = newVal
				}
			}
		}
	}

	result := dp[n][k]
	if result == "" {
		return 0
	}

	var val uint64
	for _, c := range result {
		val = val*10 + uint64(c-'0')
	}
	return val
}

func compareNumStr(a, b string) int {
	if len(a) != len(b) {
		if len(a) > len(b) {
			return 1
		}
		return -1
	}
	if a > b {
		return 1
	}
	if a < b {
		return -1
	}
	return 0
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			lines = append(lines, line)
		}
	}

	totalPart1 := 0
	for _, line := range lines {
		totalPart1 += maxJoltagePart1(line)
	}
	fmt.Println("Part 1:", totalPart1)

	var totalPart2 uint64
	for _, line := range lines {
		totalPart2 += maxJoltagePart2(line)
	}
	fmt.Println("Part 2:", totalPart2)
}
