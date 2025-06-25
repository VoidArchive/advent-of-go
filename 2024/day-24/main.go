package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Gate struct {
	input1, input2, output string
	op                     string
}

type Circuit struct {
	wires map[string]int
	gates []Gate
}

func parseInput(input string) *Circuit {
	lines := strings.Split(strings.TrimSpace(input), "\n")

	circuit := &Circuit{
		wires: make(map[string]int),
		gates: []Gate{},
	}

	// Parse initial wire values
	i := 0
	for i < len(lines) && strings.Contains(lines[i], ":") {
		parts := strings.Split(lines[i], ": ")
		wire := parts[0]
		value, _ := strconv.Atoi(parts[1])
		circuit.wires[wire] = value
		i++
	}

	// Skip empty line
	i++

	// Parse gates
	for i < len(lines) {
		parts := strings.Fields(lines[i])
		gate := Gate{
			input1: parts[0],
			op:     parts[1],
			input2: parts[2],
			output: parts[4],
		}
		circuit.gates = append(circuit.gates, gate)
		i++
	}

	return circuit
}

func (c *Circuit) simulate() {
	for {
		changed := false
		for _, gate := range c.gates {
			if _, exists := c.wires[gate.output]; exists {
				continue
			}

			val1, ok1 := c.wires[gate.input1]
			val2, ok2 := c.wires[gate.input2]

			if !ok1 || !ok2 {
				continue
			}

			var result int
			switch gate.op {
			case "AND":
				result = val1 & val2
			case "OR":
				result = val1 | val2
			case "XOR":
				result = val1 ^ val2
			}

			c.wires[gate.output] = result
			changed = true
		}

		if !changed {
			break
		}
	}
}

func (c *Circuit) getZValue() int64 {
	var zWires []string
	for wire := range c.wires {
		if strings.HasPrefix(wire, "z") {
			zWires = append(zWires, wire)
		}
	}
	sort.Strings(zWires)

	var result int64
	for i := len(zWires) - 1; i >= 0; i-- {
		result = result<<1 + int64(c.wires[zWires[i]])
	}

	return result
}

func (c *Circuit) findGateByOutput(output string) *Gate {
	for i, gate := range c.gates {
		if gate.output == output {
			return &c.gates[i]
		}
	}
	return nil
}

func (c *Circuit) findSwappedWires() []string {
	var swaps []string

	// Try 4 rounds of swaps (4 pairs = 8 wires)
	for round := range 4 {
		baseline := c.progress()
		fmt.Printf("Round %d: baseline progress = %d\n", round+1, baseline)

		found := false
		outputs := c.getAllOutputs()

		for x := range outputs {
			if found {
				break
			}
			for y := range outputs {
				if x == y {
					continue
				}

				// Try swapping x and y
				c.swapOutputs(x, y)
				newProgress := c.progress()

				if newProgress > baseline {
					fmt.Printf("Found beneficial swap: %s ↔ %s (progress: %d -> %d)\n", x, y, baseline, newProgress)
					swaps = append(swaps, x, y)
					found = true
					break
				}

				// Undo the swap if it didn't help
				c.swapOutputs(x, y)
			}
		}

		if !found {
			fmt.Printf("No beneficial swap found in round %d\n", round+1)
			break
		}
	}

	sort.Strings(swaps)
	return swaps
}

func (c *Circuit) getAllOutputs() map[string]bool {
	outputs := make(map[string]bool)
	for _, gate := range c.gates {
		outputs[gate.output] = true
	}
	return outputs
}

func (c *Circuit) swapOutputs(x, y string) {
	// Find gates with outputs x and y and swap their outputs
	for i := range c.gates {
		switch c.gates[i].output {
		case x:
			c.gates[i].output = y
		case y:
			c.gates[i].output = x
		}
	}
}

func (c *Circuit) progress() int {
	// Find how many bits are correctly computed
	i := 0
	for c.verify(i) {
		i++
	}
	return i
}

func (c *Circuit) verify(num int) bool {
	wire := fmt.Sprintf("z%02d", num)
	return c.verifyZ(wire, num)
}

func (c *Circuit) verifyZ(wire string, num int) bool {
	gate := c.findGateByOutput(wire)
	if gate == nil {
		return false
	}

	if gate.op != "XOR" {
		return false
	}

	if num == 0 {
		inputs := []string{gate.input1, gate.input2}
		sort.Strings(inputs)
		return inputs[0] == "x00" && inputs[1] == "y00"
	}

	// For higher bits: one input should be intermediate XOR, other should be carry
	return (c.verifyIntermediateXOR(gate.input1, num) && c.verifyCarryBit(gate.input2, num)) ||
		(c.verifyIntermediateXOR(gate.input2, num) && c.verifyCarryBit(gate.input1, num))
}

func (c *Circuit) verifyIntermediateXOR(wire string, num int) bool {
	gate := c.findGateByOutput(wire)
	if gate == nil {
		return false
	}

	if gate.op != "XOR" {
		return false
	}

	expected := []string{fmt.Sprintf("x%02d", num), fmt.Sprintf("y%02d", num)}
	actual := []string{gate.input1, gate.input2}
	sort.Strings(expected)
	sort.Strings(actual)

	return expected[0] == actual[0] && expected[1] == actual[1]
}

func (c *Circuit) verifyCarryBit(wire string, num int) bool {
	gate := c.findGateByOutput(wire)
	if gate == nil {
		return false
	}

	if num == 1 {
		if gate.op != "AND" {
			return false
		}
		inputs := []string{gate.input1, gate.input2}
		sort.Strings(inputs)
		return inputs[0] == "x00" && inputs[1] == "y00"
	}

	if gate.op != "OR" {
		return false
	}

	return (c.verifyDirectCarry(gate.input1, num-1) && c.verifyRecarry(gate.input2, num-1)) ||
		(c.verifyDirectCarry(gate.input2, num-1) && c.verifyRecarry(gate.input1, num-1))
}

func (c *Circuit) verifyDirectCarry(wire string, num int) bool {
	gate := c.findGateByOutput(wire)
	if gate == nil {
		return false
	}

	if gate.op != "AND" {
		return false
	}

	expected := []string{fmt.Sprintf("x%02d", num), fmt.Sprintf("y%02d", num)}
	actual := []string{gate.input1, gate.input2}
	sort.Strings(expected)
	sort.Strings(actual)

	return expected[0] == actual[0] && expected[1] == actual[1]
}

func (c *Circuit) verifyRecarry(wire string, num int) bool {
	gate := c.findGateByOutput(wire)
	if gate == nil {
		return false
	}

	if gate.op != "AND" {
		return false
	}

	return (c.verifyIntermediateXOR(gate.input1, num) && c.verifyCarryBit(gate.input2, num)) ||
		(c.verifyIntermediateXOR(gate.input2, num) && c.verifyCarryBit(gate.input1, num))
}

func solvePart1(input string) int64 {
	circuit := parseInput(input)
	circuit.simulate()
	return circuit.getZValue()
}

func solvePart2(input string) string {
	circuit := parseInput(input)

	fmt.Println("=== FINDING SWAPPED WIRES ===")
	swappedWires := circuit.findSwappedWires()

	fmt.Printf("\nFound %d swapped wires: %v\n", len(swappedWires), swappedWires)

	if len(swappedWires) != 8 {
		fmt.Printf("Warning: Expected 8 wires, found %d\n", len(swappedWires))
	}

	return strings.Join(swappedWires, ",")
}

func main() {
	content, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatal("Error reading input.txt:", err)
	}

	input := string(content)

	fmt.Printf("Part 1: %d\n", solvePart1(input))
	fmt.Printf("Part 2: %s\n", solvePart2(input))
}

