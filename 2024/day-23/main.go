package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	adj := parseInput("input.txt")

	part1 := findTrianglesWithT(adj)
	fmt.Printf("Part 1 - Triangles with 't': %d\n", part1)

	part2 := findMaxClique(adj)
	fmt.Printf("Part 2 - Password: %s\n", part2)
}

func parseInput(filename string) map[string]map[string]bool {
	file, _ := os.Open(filename)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	adj := make(map[string]map[string]bool)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "-")
		a, b := parts[0], parts[1]

		if adj[a] == nil {
			adj[a] = make(map[string]bool)
		}
		if adj[b] == nil {
			adj[b] = make(map[string]bool)
		}

		adj[a][b] = true
		adj[b][a] = true
	}
	return adj
}

func getAllNodes(adj map[string]map[string]bool) []string {
	nodes := make([]string, 0, len(adj))
	for node := range adj {
		nodes = append(nodes, node)
	}
	sort.Strings(nodes)
	return nodes
}

// Bron-Kerbosch algorithm for maximum clique
func findMaxClique(adj map[string]map[string]bool) string {
	nodes := getAllNodes(adj)
	nodeSet := make(map[string]bool)
	for _, node := range nodes {
		nodeSet[node] = true
	}

	var maxClique []string

	var bronKerbosch func(r, p, x map[string]bool)
	bronKerbosch = func(r, p, x map[string]bool) {
		if len(p) == 0 && len(x) == 0 {
			// Found maximal clique
			if len(r) > len(maxClique) {
				maxClique = make([]string, 0, len(r))
				for node := range r {
					maxClique = append(maxClique, node)
				}
			}
			return
		}

		// Choose pivot to minimize branching
		pivot := ""
		maxDegree := -1
		for node := range union(p, x) {
			degree := len(intersection(adj[node], p))
			if degree > maxDegree {
				maxDegree = degree
				pivot = node
			}
		}

		// Process nodes not connected to pivot
		for node := range difference(p, adj[pivot]) {
			newR := union(r, map[string]bool{node: true})
			newP := intersection(p, adj[node])
			newX := intersection(x, adj[node])

			bronKerbosch(newR, newP, newX)

			delete(p, node)
			x[node] = true
		}
	}

	bronKerbosch(make(map[string]bool), nodeSet, make(map[string]bool))

	sort.Strings(maxClique)
	return strings.Join(maxClique, ",")
}

// Set operations
func union(a, b map[string]bool) map[string]bool {
	result := make(map[string]bool)
	for k := range a {
		result[k] = true
	}
	for k := range b {
		result[k] = true
	}
	return result
}

func intersection(a, b map[string]bool) map[string]bool {
	result := make(map[string]bool)
	for k := range a {
		if b[k] {
			result[k] = true
		}
	}
	return result
}

func difference(a, b map[string]bool) map[string]bool {
	result := make(map[string]bool)
	for k := range a {
		if !b[k] {
			result[k] = true
		}
	}
	return result
}

func findTrianglesWithT(adj map[string]map[string]bool) int {
	nodes := getAllNodes(adj)
	count := 0
	seen := make(map[string]bool)

	for i := range nodes {
		for j := i + 1; j < len(nodes); j++ {
			for k := j + 1; k < len(nodes); k++ {
				a, b, c := nodes[i], nodes[j], nodes[k]

				if adj[a][b] && adj[b][c] && adj[a][c] {
					key := a + "," + b + "," + c

					if !seen[key] {
						seen[key] = true

						if strings.HasPrefix(a, "t") ||
							strings.HasPrefix(b, "t") ||
							strings.HasPrefix(c, "t") {
							count++
						}
					}
				}
			}
		}
	}

	return count
}
