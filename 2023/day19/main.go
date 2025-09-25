package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Rule struct {
	Field     string
	Operator  string
	Value     int
	Target    string
	IsDefault bool
}

type Workflow struct {
	Name  string
	Rules []Rule
}

type Part struct {
	X, M, A, S int
}

type Range struct {
	Min, Max int
}

type RangeSet struct {
	X, M, A, S Range
}

func (r Range) Size() int64 {
	if r.Max < r.Min {
		return 0
	}
	return int64(r.Max - r.Min + 1)
}

func (rs RangeSet) Combinations() int64 {
	return rs.X.Size() * rs.M.Size() * rs.A.Size() * rs.S.Size()
}

func (rs RangeSet) IsEmpty() bool {
	return rs.X.Max < rs.X.Min || rs.M.Max < rs.M.Min || rs.A.Max < rs.A.Min || rs.S.Max < rs.S.Min
}

func parseWorkflows(lines []string) map[string]Workflow {
	workflows := make(map[string]Workflow)
	for _, line := range lines {
		if line == "" {
			break
		}
		parts := strings.Split(line, "{")
		name := parts[0]
		rulesStr := strings.TrimSuffix(parts[1], "}")

		var rules []Rule
		ruleParts := strings.Split(rulesStr, ",")

		for i, ruleStr := range ruleParts {
			if i == len(ruleParts)-1 {
				rules = append(rules, Rule{
					Target:    ruleStr,
					IsDefault: true,
				})
			} else {
				colonIdx := strings.Index(ruleStr, ":")
				condition := ruleStr[:colonIdx]
				target := ruleStr[colonIdx+1:]

				var field, operator string
				var value int

				if strings.Contains(condition, "<") {
					parts := strings.Split(condition, "<")
					field = parts[0]
					operator = "<"
					value, _ = strconv.Atoi(parts[1])
				} else if strings.Contains(condition, ">") {
					parts := strings.Split(condition, ">")
					field = parts[0]
					operator = ">"
					value, _ = strconv.Atoi(parts[1])
				}
				rules = append(rules, Rule{
					Field:    field,
					Operator: operator,
					Value:    value,
					Target:   target,
				})
			}
		}
		workflows[name] = Workflow{
			name,
			rules,
		}
	}
	return workflows
}

func parseParts(lines []string, startIdx int) []Part {
	var parts []Part

	for i := startIdx; i < len(lines); i++ {
		line := lines[i]
		if line == "" {
			continue
		}

		line = strings.TrimPrefix(line, "{")
		line = strings.TrimSuffix(line, "}")
		ratings := strings.Split(line, ",")
		part := Part{}

		for _, rating := range ratings {
			parts := strings.Split(rating, "=")
			field := parts[0]
			value, _ := strconv.Atoi(parts[1])

			switch field {
			case "x":
				part.X = value
			case "m":
				part.M = value
			case "a":
				part.A = value
			case "s":
				part.S = value
			}
		}
		parts = append(parts, part)
	}
	return parts
}

func evaluateRule(rule Rule, part Part) bool {
	if rule.IsDefault {
		return true
	}

	var fieldValue int
	switch rule.Field {
	case "x":
		fieldValue = part.X
	case "m":
		fieldValue = part.M
	case "a":
		fieldValue = part.A
	case "s":
		fieldValue = part.S
	}

	switch rule.Operator {
	case "<":
		return fieldValue < rule.Value
	case ">":
		return fieldValue > rule.Value
	}

	return false
}

func processPart(part Part, workflows map[string]Workflow) bool {
	currentWorkflow := "in"
	for {
		workflow := workflows[currentWorkflow]
		for _, rule := range workflow.Rules {
			if evaluateRule(rule, part) {
				switch rule.Target {
				case "A":
					return true
				case "R":
					return false
				default:
					currentWorkflow = rule.Target
					goto nextWorkflow
				}
			}
		}
	nextWorkflow:
	}
}

func splitRange(rs RangeSet, rule Rule) (matching, nonMatching RangeSet) {
	matching = rs
	nonMatching = rs

	if rule.IsDefault {
		nonMatching = RangeSet{
			Range{1, 0},
			Range{1, 0},
			Range{1, 0},
			Range{1, 0},
		}
		return matching, nonMatching
	}

	var currentRange *Range
	var matchingRange *Range
	var nonMatchingRange *Range

	switch rule.Field {
	case "x":
		currentRange = &rs.X
		matchingRange = &matching.X
		nonMatchingRange = &nonMatching.X
	case "m":
		currentRange = &rs.M
		matchingRange = &matching.M
		nonMatchingRange = &nonMatching.M
	case "a":
		currentRange = &rs.A
		matchingRange = &matching.A
		nonMatchingRange = &nonMatching.A
	case "s":
		currentRange = &rs.S
		matchingRange = &matching.S
		nonMatchingRange = &nonMatching.S
	}
	switch rule.Operator {
	case "<":
		if rule.Value <= currentRange.Min {
			*matchingRange = Range{1, 0}
		} else if rule.Value > currentRange.Max {
			*nonMatchingRange = Range{1, 0}
		} else {
			*matchingRange = Range{currentRange.Min, rule.Value - 1}
			*nonMatchingRange = Range{rule.Value, currentRange.Max}
		}
	case ">":
		if rule.Value >= currentRange.Max {
			*matchingRange = Range{1, 0}
		} else if rule.Value < currentRange.Min {
			*nonMatchingRange = Range{1, 0}
		} else {
			*matchingRange = Range{rule.Value + 1, currentRange.Max}
			*nonMatchingRange = Range{currentRange.Min, rule.Value}
		}
	}
	return matching, nonMatching
}

func countAcceptedCombinations(workflowName string, ranges RangeSet, workflows map[string]Workflow) int64 {
	if ranges.IsEmpty() {
		return 0
	}
	if workflowName == "A" {
		return ranges.Combinations()
	}
	if workflowName == "R" {
		return 0
	}
	workflow := workflows[workflowName]
	currentRanges := ranges
	totalAccepted := int64(0)

	for _, rule := range workflow.Rules {
		if currentRanges.IsEmpty() {
			break
		}
		matching, nonMatching := splitRange(currentRanges, rule)

		if !matching.IsEmpty() {
			totalAccepted += countAcceptedCombinations(rule.Target, matching, workflows)
		}
		currentRanges = nonMatching
	}
	return totalAccepted
}

func main() {
	file, _ := os.Open("input.txt")
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	blankLineIdx := -1
	for i, line := range lines {
		if line == "" {
			blankLineIdx = i
			break
		}
	}

	workflows := parseWorkflows(lines[:blankLineIdx])
	if blankLineIdx != -1 && blankLineIdx < len(lines)-1 {
		parts := parseParts(lines, blankLineIdx+1)
		totalRating := 0
		acceptedCount := 0

		for _, part := range parts {
			if processPart(part, workflows) {
				totalRating += part.X + part.M + part.A + part.S
				acceptedCount++
			}
		}
		fmt.Printf("Part 1 - Accepted Parts: %d\n", acceptedCount)
		fmt.Printf("Part 1 - Total rating sum: %d\n", totalRating)
	}

	initialRanges := RangeSet{
		Range{1, 4000},
		Range{1, 4000},
		Range{1, 4000},
		Range{1, 4000},
	}
	acceptedCombinations := countAcceptedCombinations("in", initialRanges, workflows)
	fmt.Printf("Part 2 - Accepted Combinations: %d\n", acceptedCombinations)
}
