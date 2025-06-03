package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	R, C int
}

type SearchContext struct {
	Grid [][]int
	Rows int
	Cols int
}

func readandParseInput(input string) ([][]int, int, int, error) {
	file, err := os.Open(input)
	if err != nil {
		return nil, 0, 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		return nil, 0, 0, err
	}
	if len(lines) == 0 {
		return nil, 0, 0, fmt.Errorf("no input data found")
	}
	rows := len(lines)
	cols := len(lines[0])

	parsedGrid := make([][]int, rows)
	for row := range rows {
		if len(lines[row]) != cols {
			return nil, 0, 0, fmt.Errorf("inconsistent row length at row %d", row)
		}
		parsedGrid[row] = make([]int, cols)
		for col, char := range lines[row] {
			height, err := strconv.Atoi(string(char))
			if err != nil {
				return nil, 0, 0, fmt.Errorf("invalid character '%s' at row %d, col %d", string(char), row, col)
			}
			parsedGrid[row][col] = height
		}
	}
	return parsedGrid, rows, cols, nil
}

func dfsFindNines(r, c, currentHeight int, ninesFound map[Point]bool, ctx *SearchContext) {
	if currentHeight == 9 {
		ninesFound[Point{R: r, C: c}] = true
		return
	}
	dr := []int{-1, 1, 0, 0}
	dc := []int{0, 0, -1, 1}

	for i := range 4 {
		nr, nc := r+dr[i], c+dc[i]

		if nr >= 0 && nr < ctx.Rows && nc >= 0 && nc < ctx.Cols {
			if ctx.Grid[nr][nc] == currentHeight+1 {
				dfsFindNines(nr, nc, currentHeight+1, ninesFound, ctx)
			}
		}
	}
}

func dfsCountTrailsPart2(r, c, currentHeight int, ctx *SearchContext, memo [][][]int) int {
	if currentHeight == 9 {
		return 1
	}
	if memo[r][c][currentHeight] != -1 {
		return memo[r][c][currentHeight]
	}
	numberOfTrails := 0
	dr := []int{-1, 1, 0, 0}
	dc := []int{0, 0, -1, 1}

	for i := range 4 {
		nr, nc := r+dr[i], c+dc[i]
		if nr >= 0 && nr < ctx.Rows && nc >= 0 && nc < ctx.Cols && ctx.Grid[nr][nc] == currentHeight+1 {
			numberOfTrails += dfsCountTrailsPart2(nr, nc, currentHeight+1, ctx, memo)
		}
	}
	memo[r][c][currentHeight] = numberOfTrails
	return numberOfTrails
}

func solvePart1(trailheads []Point, ctx *SearchContext) (int, []int) {
	totalScore := 0
	individualScores := make([]int, 0, len(trailheads))

	for _, th := range trailheads {
		reachableNines := make(map[Point]bool)
		dfsFindNines(th.R, th.C, 0, reachableNines, ctx)
		score := len(reachableNines)
		totalScore += score
		individualScores = append(individualScores, score)
	}
	return totalScore, individualScores
}

func solvePart2(trailheads []Point, ctx *SearchContext) (int, []int) {
	totalRaiting := 0
	individualRatings := make([]int, 0, len(trailheads))

	for _, th := range trailheads {
		memo := make([][][]int, ctx.Rows)
		for rMemo := range ctx.Rows {
			memo[rMemo] = make([][]int, ctx.Cols)
			for cMemo := range ctx.Cols {
				memo[rMemo][cMemo] = make([]int, 10)
				for hMemo := range 10 {
					memo[rMemo][cMemo][hMemo] = -1
				}
			}
		}
		raiting := dfsCountTrailsPart2(th.R, th.C, 0, ctx, memo)
		totalRaiting += raiting
		individualRatings = append(individualRatings, raiting)
	}
	return totalRaiting, individualRatings
}

func runAndPrintSolutions(grid [][]int, rows int, cols int) {
	fmt.Println("Parsed Grid (" + strconv.Itoa(rows) + "x" + strconv.Itoa(cols) + ") - first few lines if large:")
	for r := 0; r < rows && r < 5; r++ {
		fmt.Println(grid[r])
	}
	if rows > 5 {
		fmt.Println("...")
	}

	var trailheads []Point
	for r := range rows {
		for c := range cols {
			if grid[r][c] == 0 {
				trailheads = append(trailheads, Point{r, c})
			}
		}
	}

	if len(trailheads) == 0 {
		fmt.Println("No trailheads found.")
		return
	}
	fmt.Printf("Found %d trailheads.\n", len(trailheads))

	searchCtx := &SearchContext{Grid: grid, Rows: rows, Cols: cols}

	// ---INFO: Solve Part 1 ---
	totalScoreP1, individualScoresP1 := solvePart1(trailheads, searchCtx)
	fmt.Println("\n--- Part One ---")
	fmt.Println("Individual trailhead scores:", individualScoresP1)
	fmt.Println("Sum of the scores of all trailheads (Part 1 Answer):", totalScoreP1)

	// ---INFO: Solve Part 2 ---
	totalRatingP2, individualRatingsP2 := solvePart2(trailheads, searchCtx)
	fmt.Println("\n--- Part Two ---")
	fmt.Println("Individual trailhead ratings:", individualRatingsP2)
	fmt.Println("Sum of the ratings of all trailheads (Part 2 Answer):", totalRatingP2)
}

func main() {
	grid, rows, cols, err := readandParseInput("input.txt")
	if err != nil {
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}
	runAndPrintSolutions(grid, rows, cols)
}
