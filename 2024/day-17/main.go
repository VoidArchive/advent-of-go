package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Computer struct {
	A, B, C int
	program []int
	ip      int
	output  []int
}

func NewComputer(a, b, c int, program []int) *Computer {
	return &Computer{
		a,
		b,
		c,
		program,
		0,
		make([]int, 0),
	}
}

func (c *Computer) getComboValue(operand int) int {
	switch operand {
	case 0, 1, 2, 3:
		return operand
	case 4:
		return c.A
	case 5:
		return c.B
	case 6:
		return c.C
	case 7:
		panic("Reserved combo operand 7 encountered")
	default:
		panic(fmt.Sprintf("Invalid combo operand: %d", operand))
	}
}

func (c *Computer) execute() {
	for c.ip < len(c.program) {
		if c.ip+1 >= len(c.program) {
			break
		}
		opcode := c.program[c.ip]
		operand := c.program[c.ip+1]

		switch opcode {
		// adv
		case 0:
			denominator := 1 << c.getComboValue(operand)
			c.A = c.A / denominator
		// bxl
		case 1:
			c.B = c.B ^ operand
		// bst
		case 2:
			c.B = c.getComboValue(operand) % 8
		// jnz
		case 3:
			if c.A != 0 {
				c.ip = operand
				continue
			}
		// bxc
		case 4:
			c.B = c.B ^ c.C

		// out
		case 5:
			c.output = append(c.output, c.getComboValue(operand)%8)
		// bdv
		case 6:
			denominator := 1 << c.getComboValue(operand)
			c.B = c.A / denominator
		// cdv
		case 7:
			denominator := 1 << c.getComboValue(operand)
			c.C = c.A / denominator
		default:
			panic(fmt.Sprintf("Unknown opcode : %d", opcode))
		}
		c.ip += 2
	}
}

func (c *Computer) getOutput() []int {
	return c.output
}

func (c *Computer) getOutputString() string {
	if len(c.output) == 0 {
		return ""
	}
	result := make([]string, len(c.output))
	for i, v := range c.output {
		result[i] = strconv.Itoa(v)
	}
	return strings.Join(result, ",")
}

func findQuineValue(program []int) int {
	candidates := []int{0}

	for pos := len(program) - 1; pos >= 0; pos-- {
		nextCandidates := []int{}

		for _, candidate := range candidates {
			for digit := range 8 {
				testA := candidate*8 + digit

				computer := NewComputer(testA, 0, 0, program)
				computer.execute()
				output := computer.getOutput()

				if len(output) >= len(program)-pos && output[0] == program[pos] {
					nextCandidates = append(nextCandidates, testA)
				}
			}
		}
		candidates = nextCandidates
	}
	minCandidate := -1
	for _, candidate := range candidates {
		if candidate > 0 {
			computer := NewComputer(candidate, 0, 0, program)
			computer.execute()
			output := computer.getOutput()

			if len(output) == len(program) {
				match := true
				for i := range len(program) {
					if output[i] != program[i] {
						match = false
						break
					}
				}
				if match && (minCandidate == -1 || candidate < minCandidate) {
					minCandidate = candidate
				}
			}
		}
	}
	return minCandidate
}

// Input
// Register A: 50230824
// Register B: 0
// Register C: 0
//
// Program: 2,4,1,3,7,5,0,3,1,4,4,7,5,5,3,0

func main() {
	program := []int{2, 4, 1, 3, 7, 5, 0, 3, 1, 4, 4, 7, 5, 5, 3, 0}
	computer := NewComputer(50230824, 0, 0, program)
	computer.execute()
	fmt.Printf("output: %s\n", computer.getOutputString())

	quineValue := findQuineValue(program)

	fmt.Println(quineValue)
}
