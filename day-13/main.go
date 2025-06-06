package main

import (
	"errors"
	"fmt"
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
	if a < 0 || a > 100 {
		return 0, fmt.Errorf("a is out of bound: %d", a)
	}

	numerator := c1 - a1*a
	if numerator%b1 != 0 {
		return 0, fmt.Errorf("b is not an integer: %d / %d", numerator, b1)
	}

	b := numerator / b1
	if b < 0 || b > 100 {
		return 0, fmt.Errorf("b is out of bound %d", b)
	}

	cost := a*3 + b

	return cost, nil
}

func solvePart1(machines []Machine) int {
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

func main() {
	machines := []Machine{
		{Point{94, 34}, Point{22, 67}, Point{8400, 5400}},   // Should give 280
		{Point{17, 86}, Point{84, 37}, Point{7870, 6450}},   // Should give 200
		{Point{26, 66}, Point{67, 21}, Point{12748, 12176}}, // Should fail
	}
	result := solvePart1(machines)
	fmt.Println(result)
}
