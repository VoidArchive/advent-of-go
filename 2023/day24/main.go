package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Vec3 struct {
	x, y, z float64
}

type HailStone struct {
	pos, vel Vec3
}

func parseInput(filename string) []HailStone {
	file, _ := os.Open(filename)
	defer file.Close()

	var hailstones []HailStone
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Split(line, "@")
		posStr := strings.Split(parts[0], ",")
		velStr := strings.Split(parts[1], ",")

		px, _ := strconv.ParseFloat(strings.TrimSpace(posStr[0]), 64)
		py, _ := strconv.ParseFloat(strings.TrimSpace(posStr[1]), 64)
		pz, _ := strconv.ParseFloat(strings.TrimSpace(posStr[2]), 64)
		vx, _ := strconv.ParseFloat(strings.TrimSpace(velStr[0]), 64)
		vy, _ := strconv.ParseFloat(strings.TrimSpace(velStr[1]), 64)
		vz, _ := strconv.ParseFloat(strings.TrimSpace(velStr[2]), 64)

		hailstones = append(hailstones, HailStone{
			pos: Vec3{px, py, pz},
			vel: Vec3{vx, vy, vz},
		})
	}
	return hailstones
}

func findIntersection2D(h1, h2 HailStone) (float64, float64, bool) {
	det := h1.vel.x*h2.vel.y - h1.vel.y*h2.vel.x
	if det == 0 {
		return 0, 0, false
	}

	dx := h2.pos.x - h1.pos.x
	dy := h2.pos.y - h1.pos.y

	t1 := (dx*h2.vel.y - dy*h2.vel.x) / det
	t2 := (dx*h1.vel.y - dy*h1.vel.x) / det

	if t1 < 0 || t2 < 0 {
		return 0, 0, false
	}

	x := h1.pos.x + t1*h1.vel.x
	y := h1.pos.y + t1*h1.vel.y
	return x, y, true
}

func part1(hailstones []HailStone, min, max float64) int {
	count := 0
	for i := range len(hailstones) {
		for j := i + 1; j < len(hailstones); j++ {
			x, y, valid := findIntersection2D(hailstones[i], hailstones[j])
			if valid && x >= min && x <= max && y >= min && y <= max {
				count++
			}
		}
	}
	return count
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func solveSystem(A [][]float64, b []float64) []float64 {
	n := len(b)
	for i := range n {
		A[i] = append(A[i], b[i])
	}

	for col := range n {
		pivot := col
		for row := col + 1; row < n; row++ {
			if abs(A[row][col]) > abs(A[pivot][col]) {
				pivot = row
			}
		}
		A[col], A[pivot] = A[pivot], A[col]

		for row := col + 1; row < n; row++ {
			factor := A[row][col] / A[col][col]
			for j := col; j <= n; j++ {
				A[row][j] -= factor * A[col][j]
			}
		}
	}

	x := make([]float64, n)
	for i := n - 1; i >= 0; i-- {
		x[i] = A[i][n]
		for j := i + 1; j < n; j++ {
			x[i] -= A[i][j] * x[j]
		}
		x[i] /= A[i][i]
	}
	return x
}

func part2(hailstones []HailStone) int64 {
	h0, h1, h2 := hailstones[0], hailstones[1], hailstones[2]
	A := [][]float64{
		{0, h0.vel.z - h1.vel.z, h1.vel.y - h0.vel.y, 0, h1.pos.z - h0.pos.z, h0.pos.y - h1.pos.y},
		{h1.vel.z - h0.vel.z, 0, h0.vel.x - h1.vel.x, h0.pos.z - h1.pos.z, 0, h1.pos.x - h0.pos.x},
		{h0.vel.y - h1.vel.y, h1.vel.x - h0.vel.x, 0, h1.pos.y - h0.pos.y, h0.pos.x - h1.pos.x, 0},
		{0, h0.vel.z - h2.vel.z, h2.vel.y - h0.vel.y, 0, h2.pos.z - h0.pos.z, h0.pos.y - h2.pos.y},
		{h2.vel.z - h0.vel.z, 0, h0.vel.x - h2.vel.x, h0.pos.z - h2.pos.z, 0, h2.pos.x - h0.pos.x},
		{h0.vel.y - h2.vel.y, h2.vel.x - h0.vel.x, 0, h2.pos.y - h0.pos.y, h0.pos.x - h2.pos.x, 0},
	}
	b := []float64{
		h0.pos.y*h0.vel.z - h0.vel.y*h0.pos.z - (h1.pos.y*h1.vel.z - h1.vel.y*h1.pos.z),
		h0.pos.z*h0.vel.x - h0.vel.z*h0.pos.x - (h1.pos.z*h1.vel.x - h1.vel.z*h1.pos.x),
		h0.pos.x*h0.vel.y - h0.vel.x*h0.pos.y - (h1.pos.x*h1.vel.y - h1.vel.x*h1.pos.y),
		h0.pos.y*h0.vel.z - h0.vel.y*h0.pos.z - (h2.pos.y*h2.vel.z - h2.vel.y*h2.pos.z),
		h0.pos.z*h0.vel.x - h0.vel.z*h0.pos.x - (h2.pos.z*h2.vel.x - h2.vel.z*h2.pos.x),
		h0.pos.x*h0.vel.y - h0.vel.x*h0.pos.y - (h2.pos.x*h2.vel.y - h2.vel.x*h2.pos.y),
	}
	solution := solveSystem(A, b)
	return int64(solution[0] + solution[1] + solution[2] + 0.5)
}

func main() {
	hailstones := parseInput("input.txt")
	fmt.Printf("Part 1: %d\n", part1(hailstones, 200000000000000, 400000000000000))
	fmt.Printf("Part 2: %d\n", part2(hailstones))
}
