package main

import "testing"

type testCase struct {
	target      int
	nums        []int
	expectPart1 bool
	expectPart2 bool
}

var testCases = []testCase{
	{190, []int{10, 19}, true, true},
	{3267, []int{81, 40, 27}, true, true},
	{83, []int{17, 5}, false, false},
	{156, []int{15, 6}, false, true},
	{7290, []int{6, 8, 6, 15}, false, true},
	{161011, []int{16, 10, 13}, false, false},
	{192, []int{17, 8, 14}, false, true},
	{21037, []int{9, 7, 18, 13}, false, false},
	{292, []int{11, 6, 16, 20}, true, true},
}

func TestValidCombinationExists(t *testing.T) {
	for _, tc := range testCases {
		got := validCombinationExists(tc.target, tc.nums)
		if got != tc.expectPart1 {
			t.Errorf("Part1: For target %d and nums %v: expected %v, got %v", tc.target, tc.nums, tc.expectPart1, got)
		}
	}
}

func TestValidCombinationExistsPart2(t *testing.T) {
	for _, tc := range testCases {
		got := validCombinationExistsPart2(tc.target, tc.nums)
		if got != tc.expectPart2 {
			t.Errorf("Part2: For target %d and nums %v: expected %v, got %v", tc.target, tc.nums, tc.expectPart2, got)
		}
	}
}

func TestParseLine(t *testing.T) {
	line := "190: 10 19"
	wantVal := 190
	wantNums := []int{10, 19}

	val, nums := parseLine(line)
	if val != wantVal {
		t.Errorf("expected %d, got %d", wantVal, val)
	}
	for i, v := range nums {
		if v != wantNums[i] {
			t.Errorf("expected %v, got %v", wantNums, nums)
			break
		}
	}
}
