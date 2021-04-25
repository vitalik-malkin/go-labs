package seed

import (
	cr "crypto/rand"
	"fmt"
	"log"
	"math"
	"math/big"
	"sort"

	intl_opts "github.com/vitalik-malkin/go-labs/randoom/internal/options"
)

func (s *Seed) NextRandom(max int32) int32 {
	maxBig := big.NewInt(int64(math.MaxInt32))
	rndBig, err := cr.Int(s, maxBig)
	if err != nil {
		log.Fatalf("error while generating next random int; err: %v", err)
		return 0
	}
	rnd := int32(rndBig.Int64())
	factor := float64(rnd) / float64(math.MaxInt32)
	res := int32(math.Floor(factor * float64(max)))
	return res
}

func (s *Seed) NextRandom20() int32 {
	return s.NextRandom(int32(20)) + 1
}

func (s *Seed) NextRandom20FieldSet(opts intl_opts.Options) [][]int32 {
	var (
		fieldSize = int32(opts.FieldSize())
		setSize   = int32(opts.FieldSetSize())
	)

	fieldSetAttempt := 0
nextFieldSet:
	fieldSetAttempt++
	if fieldSetAttempt > opts.GenFieldSetAttemptLimit() {
		fmt.Printf("===============\nAttempts limit reached. No results.\n")
		return nil
	}
	fmt.Printf("===============\nNew field set (attempt %d)...\n", fieldSetAttempt)

	numStat := make([]int, int(opts.MaxOfNum()))
	fieldSet := make([][]int32, setSize)
	for j := 0; j < len(fieldSet); j++ {
		newField := make([]int32, fieldSize)

		fieldAttempt := 0
	nextField:
		fieldAttempt++
		if fieldAttempt > opts.GenFieldAttemptLimit() {
			goto nextFieldSet
		}

		for i := 0; i < len(newField); i++ {
			newField[i] = 0
		}
		for i := 0; i < len(newField); i++ {
			numAttempt := 0
		nextNum:
			numAttempt++
			if numAttempt > opts.GenNumAttemptLimit() {
				if s.NextRandom(opts.MaxOfNum())%2 == 0 {
					goto nextField
				}
				goto nextFieldSet
			}
			num := s.NextRandom(opts.MaxOfNum()) + 1
			newField[i] = num
			if numStat[num-1] > opts.MaxOfRepeatOfNumPerFieldSet() || isViolatingNeighboringCond(newField, opts) || hasSimilarity(fieldSet, newField, opts) {
				goto nextNum
			}
		}

		sort.Slice(newField, func(i, y int) bool { return newField[i] < newField[y] })
		for _, num := range newField {
			numStat[num-1] = numStat[num-1] + 1
		}
		fieldSet[j] = newField
	}
	fmt.Print("done!\n\n")

	return fieldSet
}

func isViolatingNeighboringCond(field []int32, opts intl_opts.Options) bool {
	c := 0
	for k := 0; k < len(field); k++ {
		if field[k] == 0 {
			continue
		}

		for j := k + 1; j < len(field); j++ {
			if field[j] == 0 {
				continue
			}

			sub := field[k] - field[j]
			if -2 < sub && sub < 2 {
				c++
				if sub == 0 || c > opts.MaxOfNeighboringNumsInField() {
					return true
				}
			}
		}
	}

	return false
}

func hasSimilarity(fieldSet [][]int32, field []int32, opts intl_opts.Options) bool {
	for y := 0; y < len(fieldSet); y++ {
		existingField := fieldSet[y]
		if existingField == nil {
			continue
		}

		similarNumCount := 0
		for f := 0; f < len(field); f++ {
			for e := 0; e < len(existingField); e++ {
				if field[f] == existingField[e] {
					similarNumCount++
				}
			}
		}

		if similarNumCount > opts.FieldSimilarityDegree() {
			return true
		}
	}

	return false
}

func min(j, k int) int {
	if k < j {
		return k
	}
	return j
}
