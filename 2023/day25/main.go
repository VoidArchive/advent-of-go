package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strings"
)

func main() {
	names, W := parse("input.txt")
	minCut, partA, partB := stoerWagner(W)
	_ = minCut // always 3 for AoC Day 25 inputs

	fmt.Println(len(partA) * len(partB))

	fmt.Println("minCut:", minCut)
	fmt.Println("A:", partA)
	fmt.Println("B:", partB)
	_ = names
}

func parse(path string) ([]string, [][]int) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// First pass: collect all node names.
	nodes := make(map[string]struct{})
	type pair struct {
		u  string
		vs []string
	}
	var lines []pair

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			log.Fatalf("bad line: %q", line)
		}
		u := strings.TrimSpace(parts[0])
		vs := strings.Fields(strings.TrimSpace(parts[1]))
		lines = append(lines, pair{u, vs})
		nodes[u] = struct{}{}
		for _, v := range vs {
			nodes[v] = struct{}{}
		}
	}
	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}

	// Assign indices.
	names := make([]string, 0, len(nodes))
	for k := range nodes {
		names = append(names, k)
	}
	sort.Strings(names)
	id := make(map[string]int, len(names))
	for i, n := range names {
		id[n] = i
	}

	// Build symmetric weight matrix.
	n := len(names)
	W := make([][]int, n)
	for i := range W {
		W[i] = make([]int, n)
	}
	for _, ln := range lines {
		u := id[ln.u]
		for _, vname := range ln.vs {
			v := id[vname]
			if u == v {
				continue
			}
			W[u][v] += 1
			W[v][u] += 1
		}
	}
	return names, W
}

// stoerWagner computes the global minimum cut of an undirected weighted graph.
// W is an n×n symmetric matrix of nonnegative weights.
// Returns (minCutWeight, groupAOriginalIndices, groupBOriginalIndices).
func stoerWagner(W [][]int) (int, []int, []int) {
	n := len(W)
	if n == 0 {
		return 0, nil, nil
	}
	// verts holds the active supernodes by original index.
	verts := make([]int, n)
	for i := range verts {
		verts[i] = i
	}
	// comp maps original index -> set of original vertices in that supernode.
	comp := make([][]int, n)
	for i := range n {
		comp[i] = []int{i}
	}

	best := math.MaxInt
	var bestA []int

	m := n
	used := make([]bool, n)
	w := make([]int, n)

	for m > 1 {
		for i := range n {
			used[i] = false
			w[i] = 0
		}
		prev := -1
		order := make([]int, 0, m)

		// Maximum adjacency search over current m supernodes.
		for i := 0; i < m; i++ {
			sel := -1
			for _, v := range verts[:m] {
				if !used[v] && (sel == -1 || w[v] > w[sel]) {
					sel = v
				}
			}
			if sel == -1 {
				break
			}
			used[sel] = true
			order = append(order, sel)
			if i == m-1 {
				// Last added vertex is t; previous is s.
				t := sel
				s := prev

				// Current cut weight is w[t].
				if w[t] < best {
					best = w[t]
					// A is all vertices added before t in this phase (order[0..m-2]).
					inA := make(map[int]struct{}, 8)
					for _, u := range order[:len(order)-1] {
						inA[u] = struct{}{}
					}
					bestA = flattenGroup(comp, inA)
				}

				// Contract t into s.
				if s != -1 {
					for _, v := range verts[:m] {
						if v == s || v == t {
							continue
						}
						W[s][v] += W[t][v]
						W[v][s] = W[s][v]
					}
					// Merge component lists.
					comp[s] = append(comp[s], comp[t]...)
					// Remove t from verts by swapping with last active and shrinking m.
					for i2 := 0; i2 < m; i2++ {
						if verts[i2] == t {
							verts[i2], verts[m-1] = verts[m-1], verts[i2]
							break
						}
					}
					m--
				}
			} else {
				// Update weights to A.
				for _, v := range verts[:m] {
					if !used[v] {
						w[v] += W[sel][v]
					}
				}
				prev = sel
			}
		}
	}

	// Build complementary partition from bestA.
	inA := make(map[int]struct{}, len(bestA))
	for _, x := range bestA {
		inA[x] = struct{}{}
	}
	var bestB []int
	for i := range n {
		if _, ok := inA[i]; !ok {
			bestB = append(bestB, i)
		}
	}
	return best, bestA, bestB
}

func flattenGroup(comp [][]int, inA map[int]struct{}) []int {
	var out []int
	for super, members := range comp {
		if _, ok := inA[super]; ok {
			out = append(out, members...)
		}
	}
	return out
}
