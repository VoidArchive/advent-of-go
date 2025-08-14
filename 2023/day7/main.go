package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Hand struct {
	cards string
	bid   int
	kind  int
	ranks []int
}

const (
	HighCard = iota + 1
	OnePair
	TwoPair
	ThreeKind
	FullHouse
	FourKind
	FiveKind
)

var cardRankPart1 = map[byte]int{
	'2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9,
	'T': 10, 'J': 11, 'Q': 12, 'K': 13, 'A': 14,
}

var cardRankPart2 = map[byte]int{
	'J': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9,
	'T': 10, 'Q': 11, 'K': 12, 'A': 13,
}

func classifyHandPart1(cards string) int {
	freq := make(map[byte]int)
	for i := range len(cards) {
		freq[cards[i]]++
	}
	maxFreq, pairs := 0, 0
	for _, count := range freq {
		if count > maxFreq {
			maxFreq = count
		}
		if count == 2 {
			pairs++
		}
	}
	unique := len(freq)
	switch {
	case maxFreq == 5:
		return FiveKind
	case maxFreq == 4:
		return FourKind
	case maxFreq == 3 && unique == 2:
		return FullHouse
	case maxFreq == 3:
		return ThreeKind
	case pairs == 2:
		return TwoPair
	case maxFreq == 2:
		return OnePair
	default:
		return HighCard
	}
}

func classifyHandPart2(cards string) int {
	freq := make(map[byte]int)
	jokers := 0

	for i := range len(cards) {
		if cards[i] == 'J' {
			jokers++
		} else {
			freq[cards[i]]++
		}
	}
	if len(freq) == 0 {
		return FiveKind
	}
	maxFreq, pairs := 0, 0
	for _, count := range freq {
		if count > maxFreq {
			maxFreq = count
		}
		if count == 2 {
			pairs++
		}
	}
	maxFreq += jokers
	if jokers > 0 {
		pairs = 0
		for _, count := range freq {
			if count == 2 {
				pairs++
			}
		}
	}
	unique := len(freq)
	switch {
	case maxFreq == 5:
		return FiveKind
	case maxFreq == 4:
		return FourKind
	case maxFreq == 3 && unique == 2:
		return FullHouse
	case maxFreq == 3:
		return ThreeKind
	case pairs == 2:
		return TwoPair
	case maxFreq == 2:
		return OnePair
	default:
		return HighCard
	}
}

func parseHand(line string, part2 bool) (Hand, error) {
	parts := strings.Fields(line)
	bid, _ := strconv.Atoi(parts[1])
	cards := parts[0]
	ranks := make([]int, 5)

	var cardRank map[byte]int
	var classifyFunc func(string) int
	if part2 {
		cardRank = cardRankPart2
		classifyFunc = classifyHandPart2
	} else {

		cardRank = cardRankPart1
		classifyFunc = classifyHandPart1
	}
	for i, card := range []byte(cards) {
		rank := cardRank[card]
		ranks[i] = rank
	}
	return Hand{
		cards,
		bid,
		classifyFunc(cards),
		ranks,
	}, nil
}

func compareHands(a, b Hand) int {
	if a.kind != b.kind {
		return a.kind - b.kind
	}
	for i := range 5 {
		if a.ranks[i] != b.ranks[i] {
			return a.ranks[i] - b.ranks[i]
		}
	}
	return 0
}

func solve(filePath string, part2 bool) int {
	var hands []Hand
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		hand, _ := parseHand(line, part2)
		hands = append(hands, hand)
	}
	slices.SortFunc(hands, compareHands)
	total := 0
	for i, hand := range hands {
		rank := i + 1
		total += rank * hand.bid
	}

	return total
}

func main() {
	resultPart1 := solve("input.txt", false)
	fmt.Println(resultPart1)

	resultPart2 := solve("input.txt", true)
	fmt.Println(resultPart2)
}
