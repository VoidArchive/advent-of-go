package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Rule struct {
	before,
	after int
}

func readInput(path string) ([]Rule, [][]int) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)

	var rules []Rule
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			break
		}

		a, _ := strconv.Atoi(strings.Split(line, "|")[0])
		b, _ := strconv.Atoi(strings.Split(line, "|")[1])
		rules = append(rules, Rule{before: a, after: b})
	}

	var updates [][]int
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		raw := strings.Split(line, ",")
		up := make([]int, len(raw))
		for i, t := range raw {
			v, _ := strconv.Atoi(strings.TrimSpace(t))
			up[i] = v
		}
		updates = append(updates, up)
	}

	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}

	return rules, updates
}

func valid(update []int, rules []Rule) bool {
	pos := make(map[int]int, len(update))
	for i, p := range update {
		pos[p] = i
	}
	for _, r := range rules {
		i, okA := pos[r.before]
		j, okB := pos[r.after]
		if okA && okB && i > j {
			return false
		}
	}
	return true
}

func reorder(update []int, rules []Rule) []int {
	in := make(map[int]int)
	adj := make(map[int][]int)
	present := make(map[int]bool, len(update))

	for _, p := range update {
		present[p] = true
	}

	for _, r := range rules {
		if present[r.before] && present[r.after] {
			adj[r.before] = append(adj[r.before], r.after)
			in[r.after]++
		}
	}

	// Kahn's algorithm for deterministic tie-break
	queue := make([]int, 0)
	for _, p := range update {
		if in[p] == 0 {
			queue = append(queue, p)
		}
	}

	sort.Ints(queue)

	out := make([]int, 0, len(update))
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		out = append(out, n)

		for _, m := range adj[n] {
			in[m]--
			if in[m] == 0 {
				queue = append(queue, m)
			}
		}
		sort.Ints(queue)
	}
	return out
}

func part1(rules []Rule, updates [][]int) int {
	sum := 0
	for _, up := range updates {
		if valid(up, rules) {
			sum += up[len(up)/2]
		}
	}
	return sum
}

func part2(rules []Rule, updates [][]int) int {
	sum := 0
	for _, up := range updates {
		if !valid(up, rules) {
			re := reorder(up, rules)
			sum += re[len(re)/2]
		}
	}
	return sum
}

func main() {
	rules, updates := readInput("input.txt")
	fmt.Println(rules)
	fmt.Println(updates)
	fmt.Println("Part 1:", part1(rules, updates))
	fmt.Println("Part 2:", part2(rules, updates))
}
