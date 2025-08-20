package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
)

type p struct{ r, c int }

const (
	N = iota
	S
	W
	E
)

var (
	dirs = []p{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	opp  = []int{S, N, E, W}
)

func readInput() []string {
	path := "input.txt"
	if len(os.Args) > 1 && os.Args[1] != "" {
		path = os.Args[1]
	}

	if f, err := os.Open(path); err == nil {
		defer f.Close()
		return readAllLines(f)
	}
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		return readAllLines(os.Stdin)
	}
	return nil
}

func readAllLines(r io.Reader) []string {
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 1<<20), 1<<20)
	var out []string
	for sc.Scan() {
		out = append(out, sc.Text())
	}
	return out
}

func maxLen(xs []string) (m int) {
	for _, s := range xs {
		if len(s) > m {
			m = len(s)
		}
	}
	return
}

func in(r, c, h, w int) bool {
	return r >= 0 && r < h && c >= 0 && c < w
}

func make2D[T any](h, w int, init T) [][]T {
	m := make([][]T, h)
	for i := range m {
		row := make([]T, w)
		for j := range row {
			row[j] = init
		}
		m[i] = row
	}
	return m
}

func connDirs(ch rune) []int {
	switch ch {
	case '|':
		return []int{N, S}
	case '-':
		return []int{W, E}
	case 'L':
		return []int{N, E}
	case 'J':
		return []int{N, W}
	case '7':
		return []int{S, W}
	case 'F':
		return []int{S, E}
	default:
		return nil
	}
}

func connectsBack(ch rune, want int) bool {
	return slices.Contains(connDirs(ch), want)
}

func resolveS(sConn [4]bool) rune {
	switch {
	case sConn[N] && sConn[S]:
		return '|'
	case sConn[W] && sConn[E]:
		return '-'
	case sConn[N] && sConn[E]:
		return 'L'
	case sConn[N] && sConn[W]:
		return 'J'
	case sConn[S] && sConn[W]:
		return '7'
	case sConn[S] && sConn[E]:
		return 'F'
	default:
		return '.'
	}
}

func withSRune(ch rune, s rune) rune {
	if ch == 'S' {
		return s
	}
	return ch
}

func main() {
	lines := readInput()
	if len(lines) == 0 {
		return
	}

	h, w := len(lines), maxLen(lines)
	grid := make([][]rune, h)
	var start p
	for i, s := range lines {
		row := make([]rune, w)
		rs := []rune(s)
		for j := range row {
			if j < len(rs) {
				row[j] = rs[j]
				if row[j] == 'S' {
					start = p{i, j}
				}
			} else {
				row[j] = '.'
			}
		}
		grid[i] = row
	}

	sConn := [4]bool{}
	for d := range dirs {
		nr, nc := start.r+dirs[d].r, start.c+dirs[d].c
		if in(nr, nc, h, w) && connectsBack(grid[nr][nc], opp[d]) {
			sConn[d] = true
		}
	}

	sRune := resolveS(sConn)

	dist := make2D(h, w, -1)
	inLoop := make2D(h, w, 0)

	q := make([]p, 0, h*w)
	push := func(x p) { q = append(q, x) }
	pop := func() p { v := q[0]; q = q[1:]; return v }

	dist[start.r][start.c] = 0
	inLoop[start.r][start.c] = 1
	push(start)

	for len(q) > 0 {
		cur := pop()
		ch := grid[cur.r][cur.c]
		allowed := [4]bool{}
		if ch == 'S' {
			allowed = sConn
		} else {
			for _, d := range connDirs(ch) {
				allowed[d] = true
			}
		}

		for d := range dirs {
			if !allowed[d] {
				continue
			}
			nr, nc := cur.r+dirs[d].r, cur.c+dirs[d].c
			if !in(nr, nc, h, w) {
				continue
			}

			if !connectsBack(withSRune(grid[nr][nc], sRune), opp[d]) {
				continue
			}

			if dist[nr][nc] == -1 {
				dist[nr][nc] = dist[cur.r][cur.c] + 1
				inLoop[nr][nc] = 1
				push(p{nr, nc})
			}
		}
	}

	part1 := 0
	for r := range dist {
		for c := range dist[r] {
			if dist[r][c] > part1 {
				part1 = dist[r][c]
			}
		}
	}

	simple := make([][]rune, h)
	for r := range grid {
		row := make([]rune, w)
		for c := range grid[r] {
			if inLoop[r][c] == 1 {
				if grid[r][c] == 'S' {
					row[c] = sRune
				} else {
					row[c] = grid[r][c]
				}
			} else {
				row[c] = '.'
			}
		}
		simple[r] = row
	}

	part2 := 0
	for r := range simple {
		inside := false
		var pending rune
		for _, ch := range simple[r] {
			switch ch {
			case '.':
				if inside {
					part2++
				}
			case '|':
				inside = !inside
			case 'F', 'L':
				pending = ch
			case 'J':
				if pending == 'F' {
					inside = !inside
				}
			case '7':
				if pending == 'L' {
					inside = !inside
				}
				pending = 0
			default:

			}
		}
	}
	fmt.Println(part1)
	fmt.Println(part2)
}
