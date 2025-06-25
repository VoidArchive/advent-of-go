package main

import (
	"fmt"
	"sort"
)

type HorizontalSegment struct {
	Row int
	C1  int
	C2  int
}

type VerticalSegment struct {
	Col int
	R1  int
	R2  int
}

func calculateNumSides(region []Point, gardenMap [][]rune, plantType rune, numRows, numCols int) int {
	horizSegmentSet := make(map[HorizontalSegment]bool)
	vertSegmentSet := make(map[VerticalSegment]bool)

	for _, p := range region {
		nr, nc := p.R-1, p.C
		if nr < 0 || gardenMap[nr][nc] != plantType {
			horizSegmentSet[HorizontalSegment{p.R, p.C, p.C + 1}] = true
		}

		nr, nc = p.R+1, p.C
		if nr >= numRows || gardenMap[nr][nc] != plantType {
			horizSegmentSet[HorizontalSegment{p.R + 1, p.C, p.C + 1}] = true
		}

		nr, nc = p.R, p.C-1
		if nc < 0 || gardenMap[nr][nc] != plantType {
			vertSegmentSet[VerticalSegment{p.C, p.R, p.R + 1}] = true
		}
		nr, nc = p.R, p.C+1
		if nc >= numCols || gardenMap[nr][nc] != plantType {
			vertSegmentSet[VerticalSegment{p.C + 1, p.R, p.R + 1}] = true
		}
	}

	totalSides := 0
	groupedHoriz := make(map[int][]HorizontalSegment)
	for seg := range horizSegmentSet {
		groupedHoriz[seg.Row] = append(groupedHoriz[seg.Row], seg)
	}
	for _, segsInRow := range groupedHoriz {
		sort.Slice(segsInRow, func(i, j int) bool {
			return segsInRow[i].C1 < segsInRow[j].C1
		})

		if len(segsInRow) == 0 {
			continue
		}
		currentMergedSides := 0
		i := 0
		for i < len(segsInRow) {
			currentMergedSides++
			currentEndC2 := segsInRow[i].C2
			i++

			for i < len(segsInRow) && segsInRow[i].C1 == currentEndC2 {
				currentEndC2 = segsInRow[i].C2
				i++
			}
		}
		totalSides += currentMergedSides
	}

	groupedVert := make(map[int][]VerticalSegment)
	for seg := range vertSegmentSet {
		groupedVert[seg.Col] = append(groupedVert[seg.Col], seg)
	}
	for _, segsInCol := range groupedVert {
		sort.Slice(segsInCol, func(i, j int) bool {
			return segsInCol[i].R1 < segsInCol[j].R1
		})

		if len(segsInCol) == 0 {
			continue
		}
		currentMergedSides := 0
		i := 0
		for i < len(segsInCol) {
			currentMergedSides++
			currentEndR2 := segsInCol[i].R2
			i++

			for i < len(segsInCol) && segsInCol[i].R1 == currentEndR2 {
				currentEndR2 = segsInCol[i].R2
				i++
			}
		}
		totalSides += currentMergedSides
	}
	return totalSides
}

func solvePart2(gardenMap [][]rune) int {
	if len(gardenMap) == 0 {
		fmt.Println("Garden Map is empty")
		return 0
	}

	numRows := len(gardenMap)
	numCols := len(gardenMap[0])

	visited := make([][]bool, numRows)
	for i := range visited {
		visited[i] = make([]bool, numCols)
	}

	totalFencePrice := 0
	for r := range numRows {
		for c := range numCols {
			if !visited[r][c] {
				currentPlantType := gardenMap[r][c]

				pointsInRegion, currentRegArea, _ := bfs(gardenMap, visited, r, c, currentPlantType, numRows, numCols)

				numSides := calculateNumSidesCorners(pointsInRegion, numRows, numCols)

				regionPrice := currentRegArea * numSides
				totalFencePrice += regionPrice

			}
		}
	}
	return totalFencePrice
}

func calculateNumSidesDebug(region []Point, gardenMap [][]rune, plantType rune, numRows, numCols int) int {
	fmt.Printf("Calculating sides for region with %d points of type %c\n", len(region), plantType)

	horizSegmentSet := make(map[HorizontalSegment]bool)
	vertSegmentSet := make(map[VerticalSegment]bool)

	// Collect all segments
	for _, p := range region {
		// Top edge
		nr, nc := p.R-1, p.C
		if nr < 0 || gardenMap[nr][nc] != plantType {
			seg := HorizontalSegment{p.R, p.C, p.C + 1}
			horizSegmentSet[seg] = true
			fmt.Printf("Top edge: row %d, cols %d-%d\n", seg.Row, seg.C1, seg.C2)
		}

		// Bottom edge
		nr, nc = p.R+1, p.C
		if nr >= numRows || gardenMap[nr][nc] != plantType {
			seg := HorizontalSegment{p.R + 1, p.C, p.C + 1}
			horizSegmentSet[seg] = true
			fmt.Printf("Bottom edge: row %d, cols %d-%d\n", seg.Row, seg.C1, seg.C2)
		}

		// Left edge
		nr, nc = p.R, p.C-1
		if nc < 0 || gardenMap[nr][nc] != plantType {
			seg := VerticalSegment{p.C, p.R, p.R + 1}
			vertSegmentSet[seg] = true
			fmt.Printf("Left edge: col %d, rows %d-%d\n", seg.Col, seg.R1, seg.R2)
		}

		// Right edge
		nr, nc = p.R, p.C+1
		if nc >= numCols || gardenMap[nr][nc] != plantType {
			seg := VerticalSegment{p.C + 1, p.R, p.R + 1}
			vertSegmentSet[seg] = true
			fmt.Printf("Right edge: col %d, rows %d-%d\n", seg.Col, seg.R1, seg.R2)
		}
	}

	totalSides := 0

	// Process horizontal segments
	fmt.Printf("\nProcessing horizontal segments:\n")
	groupedHoriz := make(map[int][]HorizontalSegment)
	for seg := range horizSegmentSet {
		groupedHoriz[seg.Row] = append(groupedHoriz[seg.Row], seg)
	}

	for row, segsInRow := range groupedHoriz {
		fmt.Printf("Row %d has %d segments: ", row, len(segsInRow))
		for _, seg := range segsInRow {
			fmt.Printf("[%d-%d] ", seg.C1, seg.C2)
		}
		fmt.Println()

		sort.Slice(segsInRow, func(i, j int) bool {
			return segsInRow[i].C1 < segsInRow[j].C1
		})

		if len(segsInRow) == 0 {
			continue
		}

		currentMergedSides := 0
		i := 0
		for i < len(segsInRow) {
			currentMergedSides++
			currentEndC2 := segsInRow[i].C2
			startC1 := segsInRow[i].C1
			fmt.Printf("  Starting new side from %d to %d", startC1, currentEndC2)
			i++

			for i < len(segsInRow) && segsInRow[i].C1 == currentEndC2 {
				fmt.Printf(" -> extending to %d", segsInRow[i].C2)
				currentEndC2 = segsInRow[i].C2
				i++
			}
			fmt.Printf(" (final side: %d to %d)\n", startC1, currentEndC2)
		}
		fmt.Printf("  Total merged sides for row %d: %d\n", row, currentMergedSides)
		totalSides += currentMergedSides
	}

	// Process vertical segments
	fmt.Printf("\nProcessing vertical segments:\n")
	groupedVert := make(map[int][]VerticalSegment)
	for seg := range vertSegmentSet {
		groupedVert[seg.Col] = append(groupedVert[seg.Col], seg)
	}

	for col, segsInCol := range groupedVert {
		fmt.Printf("Col %d has %d segments: ", col, len(segsInCol))
		for _, seg := range segsInCol {
			fmt.Printf("[%d-%d] ", seg.R1, seg.R2)
		}
		fmt.Println()

		sort.Slice(segsInCol, func(i, j int) bool {
			return segsInCol[i].R1 < segsInCol[j].R1
		})

		if len(segsInCol) == 0 {
			continue
		}

		currentMergedSides := 0
		i := 0
		for i < len(segsInCol) {
			currentMergedSides++
			currentEndR2 := segsInCol[i].R2
			startR1 := segsInCol[i].R1
			fmt.Printf("  Starting new side from %d to %d", startR1, currentEndR2)
			i++

			for i < len(segsInCol) && segsInCol[i].R1 == currentEndR2 {
				fmt.Printf(" -> extending to %d", segsInCol[i].R2)
				currentEndR2 = segsInCol[i].R2
				i++
			}
			fmt.Printf(" (final side: %d to %d)\n", startR1, currentEndR2)
		}
		fmt.Printf("  Total merged sides for col %d: %d\n", col, currentMergedSides)
		totalSides += currentMergedSides
	}

	fmt.Printf("Total sides: %d\n\n", totalSides)
	return totalSides
}

func calculateNumSidesCorners(region []Point, numRows, numCols int) int {
	// Create a set of region points for fast lookup
	regionSet := make(map[Point]bool)
	for _, p := range region {
		regionSet[p] = true
	}

	// Count corners - number of corners equals number of sides
	corners := 0

	for _, p := range region {
		r, c := p.R, p.C

		// Check all 4 corners of this cell
		// For each corner, check the 2 adjacent cells

		// Top-left corner
		topExists := r > 0 && regionSet[Point{r - 1, c}]
		leftExists := c > 0 && regionSet[Point{r, c - 1}]
		topLeftExists := r > 0 && c > 0 && regionSet[Point{r - 1, c - 1}]

		// Convex corner (external corner)
		if !topExists && !leftExists {
			corners++
		}
		// Concave corner (internal corner)
		if topExists && leftExists && !topLeftExists {
			corners++
		}

		// Top-right corner
		topExists = r > 0 && regionSet[Point{r - 1, c}]
		rightExists := c < numCols-1 && regionSet[Point{r, c + 1}]
		topRightExists := r > 0 && c < numCols-1 && regionSet[Point{r - 1, c + 1}]

		// Convex corner
		if !topExists && !rightExists {
			corners++
		}
		// Concave corner
		if topExists && rightExists && !topRightExists {
			corners++
		}

		// Bottom-left corner
		bottomExists := r < numRows-1 && regionSet[Point{r + 1, c}]
		leftExists = c > 0 && regionSet[Point{r, c - 1}]
		bottomLeftExists := r < numRows-1 && c > 0 && regionSet[Point{r + 1, c - 1}]

		// Convex corner
		if !bottomExists && !leftExists {
			corners++
		}
		// Concave corner
		if bottomExists && leftExists && !bottomLeftExists {
			corners++
		}

		// Bottom-right corner
		bottomExists = r < numRows-1 && regionSet[Point{r + 1, c}]
		rightExists = c < numCols-1 && regionSet[Point{r, c + 1}]
		bottomRightExists := r < numRows-1 && c < numCols-1 && regionSet[Point{r + 1, c + 1}]

		// Convex corner
		if !bottomExists && !rightExists {
			corners++
		}
		// Concave corner
		if bottomExists && rightExists && !bottomRightExists {
			corners++
		}
	}

	return corners
}
