package main

import (
	"fmt"
	"math/big"
	"os"
	"strings"
)

type Graph struct {
	nodes      map[string][2]string
	directions string
}

type TargetChecker interface {
	IsTarget(node string) bool
}

type ExactTarget string

func (e ExactTarget) IsTarget(node string) bool {
	return node == string(e)
}

type SuffixTarget string

func (s SuffixTarget) IsTarget(node string) bool {
	return strings.HasSuffix(node, string(s))
}

func ParseInput(input string) (*Graph, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("input too short")
	}

	var directions string
	var nodeStart int
	for i, line := range lines {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			directions = trimmed
			nodeStart = i + 1
			break
		}
	}

	if directions == "" {
		return nil, fmt.Errorf("no directions found")
	}

	for nodeStart < len(lines) && strings.TrimSpace(lines[nodeStart]) == "" {
		nodeStart++
	}

	nodes := make(map[string][2]string)

	for i := nodeStart; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		parts := strings.Split(line, " = ")
		if len(parts) != 2 {
			continue
		}
		node := strings.TrimSpace(parts[0])
		pairStr := strings.TrimSpace(parts[1])

		if !strings.HasPrefix(pairStr, "(") || !strings.HasSuffix(pairStr, ")") {
			continue
		}
		pairStr = pairStr[1 : len(pairStr)-1]

		lr := strings.Split(pairStr, ", ")
		if len(lr) != 2 {
			continue
		}

		left := strings.TrimSpace(lr[0])
		right := strings.TrimSpace(lr[1])
		nodes[node] = [2]string{left, right}
	}
	return &Graph{
		nodes:      nodes,
		directions: directions,
	}, nil
}

func (g *Graph) StepsToTarget(start string, checker TargetChecker) (uint64, error) {
	if checker.IsTarget(start) {
		return 0, nil
	}
	current := start
	var steps uint64
	dirLen := uint64(len(g.directions))

	for {
		if _, exists := g.nodes[current]; !exists {
			return 0, fmt.Errorf("nodes %s not found", current)
		}

		dir := g.directions[steps%dirLen]
		switch dir {
		case 'L':
			current = g.nodes[current][0]
		case 'R':
			current = g.nodes[current][1]
		default:
			return 0, fmt.Errorf("invalid direction: %c", dir)
		}

		steps++
		if checker.IsTarget(current) {
			return steps, nil
		}

		if steps > 1e9 {
			return 0, fmt.Errorf("too many steps, possible infinite loop")
		}
	}
}

func (g *Graph) FindStartingNodes(suffix string) []string {
	var starts []string
	for node := range g.nodes {
		if strings.HasSuffix(node, suffix) {
			starts = append(starts, node)
		}
	}
	return starts
}

func (g *Graph) SolvePart1() (uint64, error) {
	if _, hasAAA := g.nodes["AAA"]; !hasAAA {
		return 0, fmt.Errorf("AAA node not found")
	}
	if _, hasZZZ := g.nodes["ZZZ"]; !hasZZZ {
		return 0, fmt.Errorf("ZZZ node not found")
	}

	return g.StepsToTarget("AAA", ExactTarget("ZZZ"))
}

func (g *Graph) SolvePart2() (*big.Int, error) {
	starts := g.FindStartingNodes("A")
	if len(starts) == 0 {
		return big.NewInt(0), nil
	}
	periods := make([]*big.Int, 0, len(starts))
	checker := SuffixTarget("Z")

	for _, start := range starts {
		steps, err := g.StepsToTarget(start, checker)
		if err != nil {
			return nil, fmt.Errorf("no path from %s to Z node: %w", start, err)
		}
		if steps == 0 {
			return big.NewInt(0), nil
		}
		periods = append(periods, new(big.Int).SetUint64(steps))
	}
	return LCM(periods...), nil
}

func GCD(a, b *big.Int) *big.Int {
	return new(big.Int).GCD(nil, nil, a, b)
}

func LCM(nums ...*big.Int) *big.Int {
	if len(nums) == 0 {
		return big.NewInt(0)
	}
	if len(nums) == 1 {
		return new(big.Int).Set(nums[0])
	}
	result := new(big.Int).Set(nums[0])
	zero := big.NewInt(0)

	for i := 1; i < len(nums); i++ {
		if result.Cmp(zero) == 0 || nums[i].Cmp(zero) == 0 {
			return big.NewInt(0)
		}
		gcd := GCD(result, nums[i])
		result.Div(result, gcd)
		result.Mul(result, nums[i])
	}
	return result
}

func main() {
	input, _ := os.ReadFile("input.txt")
	graph, _ := ParseInput(string(input))

	part1, _ := graph.SolvePart1()
	fmt.Println(part1)

	part2, _ := graph.SolvePart2()
	fmt.Println(part2.String())
}
