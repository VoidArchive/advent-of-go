package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	scanner.Scan()
	patternsLine := scanner.Text()
	patterns := make([]string, 0)
	for p := range strings.SplitSeq(patternsLine, ", ") {
		patterns = append(patterns, strings.TrimSpace(p))
	}

	fmt.Println(patterns)
}
