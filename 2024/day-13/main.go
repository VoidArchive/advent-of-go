package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Point struct {
	X int
	Y int
}
type Machine struct {
	buttonA Point
	buttonB Point
	Prize   Point
}

func solveLinearSystem(buttonA, buttonB, Prize Point) (int, error) {
	// INFO: Solve using elimination method

	a1, b1, c1 := buttonA.X, buttonB.X, Prize.X
	a2, b2, c2 := buttonA.Y, buttonB.Y, Prize.Y

	m1, m2 := b2, b1

	e1a, e1c := a1*m1, c1*m1
	e2a, e2c := a2*m2, c2*m2

	finalA, finalC := e1a-e2a, e1c-e2c

	if finalA == 0 {
		return 0, errors.New("division by zero: no unique solution")
	}

	if finalC%finalA != 0 {
		return 0, fmt.Errorf("a is not an integer: %d / %d", finalC, finalA)
	}

	a := finalC / finalA

	// INFO: PART 1: constraint
	// if a < 0 || a > 100 {
	// 	return 0, fmt.Errorf("a is out of bound: %d", a)
	// }

	numerator := c1 - a1*a
	if numerator%b1 != 0 {
		return 0, fmt.Errorf("b is not an integer: %d / %d", numerator, b1)
	}

	b := numerator / b1

	// INFO: PART1: constraint
	// if b < 0 || b > 100 {
	// 	return 0, fmt.Errorf("b is out of bound %d", b)
	// }

	if a < 0 || b < 0 {
		return 0, fmt.Errorf("negative presses not allowed")
	}

	cost := a*3 + b

	return cost, nil
}

func solve(machines []Machine) int {
	totalCost := 0
	for _, machine := range machines {
		cost, err := solveLinearSystem(machine.buttonA, machine.buttonB, machine.Prize)
		if err != nil {
			fmt.Printf("Machine is not solvable: %v\n", err)
		} else {
			totalCost += cost
		}
	}
	return totalCost
}

func parseInputFile(filename string) ([]Machine, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var machines []Machine

	buttonRegex := regexp.MustCompile(`Button [AB]: X\+(\d+), Y\+(\d+)`)
	prizeRegex := regexp.MustCompile(`Prize: X=(\d+), Y=(\d+)`)

	var ax, ay, bx, by int

	// INFO: Part2: Offset
	offset := 10000000000000

	for scanner.Scan() {
		line := scanner.Text()

		if matches := buttonRegex.FindStringSubmatch(line); matches != nil {
			x, _ := strconv.Atoi(matches[1])
			y, _ := strconv.Atoi(matches[2])

			if strings.Contains(line, "Button A") {
				ax, ay = x, y
			} else {
				bx, by = x, y
			}
		} else if matches := prizeRegex.FindStringSubmatch(line); matches != nil {
			px, _ := strconv.Atoi(matches[1])
			py, _ := strconv.Atoi(matches[2])

			machines = append(machines, Machine{
				Point{ax, ay},
				Point{bx, by},
				Point{px + offset, py + offset},
			})
		}
	}
	return machines, scanner.Err()
}

func main() {
	machines, _ := parseInputFile("input.txt")
	result := solve(machines)
	fmt.Println(result)
}
