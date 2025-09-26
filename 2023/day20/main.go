package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Pulse struct {
	from, to string
	high     bool
}

type Module struct {
	name  string
	typ   byte // '%', '&', 'b'
	dests []string
	state bool
	mem   map[string]bool
}

func (m *Module) process(from string, high bool) []Pulse {
	switch m.typ {
	case '%': // flip-flop
		if high {
			return nil
		}
		m.state = !m.state
		return m.send(m.state)
	case '&': // conjunction
		m.mem[from] = high
		allHigh := true
		for _, v := range m.mem {
			if !v {
				allHigh = false
				break
			}
		}
		return m.send(!allHigh)
	default: // broadcast
		return m.send(high)
	}
}

func (m *Module) send(high bool) []Pulse {
	pulses := make([]Pulse, len(m.dests))
	for i, dest := range m.dests {
		pulses[i] = Pulse{m.name, dest, high}
	}
	return pulses
}

func parse(lines []string) map[string]*Module {
	modules := make(map[string]*Module)

	for _, line := range lines {
		parts := strings.Split(line, " -> ")
		spec, dests := parts[0], strings.Split(parts[1], ", ")

		var name string
		var typ byte
		if spec == "broadcaster" {
			name, typ = "broadcaster", 'b'
		} else {
			name, typ = spec[1:], spec[0]
		}

		modules[name] = &Module{
			name:  name,
			typ:   typ,
			dests: dests,
			mem:   make(map[string]bool),
		}
	}

	// Wire conjunction inputs
	for name, mod := range modules {
		for _, dest := range mod.dests {
			if destMod := modules[dest]; destMod != nil && destMod.typ == '&' {
				destMod.mem[name] = false
			}
		}
	}

	return modules
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int64) int64 { return a * b / gcd(a, b) }

func simulate(modules map[string]*Module, presses int64, findRx bool) (int, int, int64) {
	low, high := 0, 0
	cycles := make(map[string]int64)
	var rxInputs []string

	if findRx {
		for name, mod := range modules {
			for _, dest := range mod.dests {
				if dest == "rx" {
					for iName, iMod := range modules {
						for _, iDest := range iMod.dests {
							if iDest == name {
								rxInputs = append(rxInputs, iName)
							}
						}
					}
				}
			}
		}
	}

	for press := int64(1); press <= presses; press++ {
		queue := []Pulse{{"button", "broadcaster", false}}

		for len(queue) > 0 {
			p := queue[0]
			queue = queue[1:]

			if p.high {
				high++
			} else {
				low++
			}

			if findRx && p.high {
				for _, input := range rxInputs {
					if p.from == input && cycles[input] == 0 {
						cycles[input] = press
						if len(cycles) == len(rxInputs) {
							result := cycles[rxInputs[0]]
							for i := 1; i < len(rxInputs); i++ {
								result = lcm(result, cycles[rxInputs[i]])
							}
							return low, high, result
						}
					}
				}
			}

			if mod := modules[p.to]; mod != nil {
				queue = append(queue, mod.process(p.from, p.high)...)
			}
		}
	}

	return low, high, -1
}

func main() {
	file, _ := os.Open("input.txt")
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if line := strings.TrimSpace(scanner.Text()); line != "" {
			lines = append(lines, line)
		}
	}

	modules := parse(lines)

	// Part 1
	low, high, _ := simulate(modules, 1000, false)
	fmt.Printf("Part 1: %d\n", low*high)

	// Reset
	for _, mod := range modules {
		mod.state = false
		for k := range mod.mem {
			mod.mem[k] = false
		}
	}

	// Part 2
	_, _, result := simulate(modules, 10000000, true)
	fmt.Printf("Part 2: %d\n", result)
}
