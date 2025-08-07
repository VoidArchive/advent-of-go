package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Range struct {
	src, dst, len uint64
}

type Stage []Range

type Interval struct {
	start, length uint64
}

func (s Stage) Apply(x uint64) uint64 {
	i := sort.Search(len(s), func(i int) bool { return x < s[i].src+s[i].len })
	if i < len(s) && x >= s[i].src {
		return s[i].dst + (x - s[i].src)
	}
	return x
}

func (s Stage) ApplyInterval(interval Interval) []Interval {
	var result []Interval
	cur := interval.start
	end := interval.start + interval.length

	for _, r := range s {
		// Skip ranges that don't overlap with our interval
		if r.src+r.len <= cur || r.src >= end {
			continue
		}

		// Gap before mapped segment
		if cur < r.src {
			gapLen := r.src - cur
			if cur+gapLen > end {
				gapLen = end - cur
			}
			result = append(result, Interval{cur, gapLen})
			cur += gapLen
		}

		if cur >= end {
			break
		}

		// Overlap segment
		overlapStart := cur
		overlapEnd := min(end, r.src+r.len)

		if overlapStart < overlapEnd {
			segLen := overlapEnd - overlapStart
			mappedStart := r.dst + (overlapStart - r.src)
			result = append(result, Interval{mappedStart, segLen})
			cur = overlapEnd
		}

		if cur >= end {
			break
		}
	}

	// Tail gap
	if cur < end {
		result = append(result, Interval{cur, end - cur})
	}

	return result
}

func parseSeeds(line string) []uint64 {
	parts := strings.Fields(line[strings.IndexByte(line, ':')+1:])
	out := make([]uint64, len(parts))
	for i, p := range parts {
		v, _ := strconv.ParseUint(p, 10, 64)
		out[i] = v
	}
	return out
}

func parseSeedRanges(line string) []Interval {
	parts := strings.Fields(line[strings.IndexByte(line, ':')+1:])
	var intervals []Interval
	for i := 0; i < len(parts); i += 2 {
		start, _ := strconv.ParseUint(parts[i], 10, 64)
		length, _ := strconv.ParseUint(parts[i+1], 10, 64)
		intervals = append(intervals, Interval{start, length})
	}
	return intervals
}

func parseInput(filename string) ([]uint64, []Interval, []Stage) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	if !sc.Scan() {
		log.Fatal("empty file")
	}
	seedLine := sc.Text()
	seeds := parseSeeds(seedLine)
	seedRanges := parseSeedRanges(seedLine)

	var stages []Stage
	for sc.Scan() {
		if strings.TrimSpace(sc.Text()) == "" {
			continue
		}
		var st Stage
		for sc.Scan() {
			line := strings.TrimSpace(sc.Text())
			if line == "" {
				break
			}
			if strings.Contains(line, ":") {
				continue
			}
			f := strings.Fields(line)
			dst, _ := strconv.ParseUint(f[0], 10, 64)
			src, _ := strconv.ParseUint(f[1], 10, 64)
			l, _ := strconv.ParseUint(f[2], 10, 64)
			st = append(st, Range{src, dst, l})
		}
		sort.Slice(st, func(i, j int) bool { return st[i].src < st[j].src })
		stages = append(stages, st)
	}

	return seeds, seedRanges, stages
}

func solvePart1(seeds []uint64, stages []Stage) uint64 {
	min := ^uint64(0)
	for _, s := range seeds {
		v := s
		for _, st := range stages {
			v = st.Apply(v)
		}
		if v < min {
			min = v
		}
	}
	return min
}

func solvePart2(seedRanges []Interval, stages []Stage) uint64 {
	intervals := make([]Interval, len(seedRanges))
	copy(intervals, seedRanges)

	for _, stage := range stages {
		var nextIntervals []Interval
		for _, interval := range intervals {
			mapped := stage.ApplyInterval(interval)
			nextIntervals = append(nextIntervals, mapped...)
		}
		intervals = nextIntervals
	}

	minLocation := ^uint64(0)
	for _, interval := range intervals {
		if interval.start < minLocation {
			minLocation = interval.start
		}
	}

	return minLocation
}

func main() {
	seeds, seedRanges, stages := parseInput("input.txt")

	fmt.Printf("Found %d seeds and %d seed ranges and %d mapping stages\n", len(seeds), len(seedRanges), len(stages))

	result1 := solvePart1(seeds, stages)
	fmt.Printf("Part 1: %d\n", result1)

	result2 := solvePart2(seedRanges, stages)
	fmt.Printf("Part 2: %d\n", result2)
}
