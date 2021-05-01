package seed

import (
	"context"
	cr "crypto/rand"
	"fmt"
	"log"
	"math"
	"math/big"
	"sort"

	gt "github.com/buger/goterm"

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

func (s *Seed) NextRandomFieldSet(opts intl_opts.Options) [][]int32 {
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

func (s *Seed) NextRandomFieldSetV2(opts intl_opts.Options) [][]int32 {
	var (
		fieldSize                      = opts.FieldSize()
		setSize                        = opts.FieldSetSize()
		maxOfNum                       = opts.MaxOfNum()
		attemptLimit                   = opts.GenFieldSetAttemptLimit()
		numTotalUses                   = setSize * fieldSize
		numUses                        = make([]int, maxOfNum)
		numMinUse                      = numTotalUses / int(maxOfNum)
		numMaxOfMinUsePlusOne          = numTotalUses % int(maxOfNum)
		genNumAttemptLimit             = opts.GenNumAttemptLimit()
		genFieldAttemptLimit           = opts.GenFieldAttemptLimit()
		maxOfFieldsWithNeighboringNums = opts.MaxOfFieldsWithNeighboringNums()
	)

	gt.Clear()
	gt.MoveCursor(1, 1)
	gt.Print(gt.Color("======================================================", gt.GREEN))
	gt.Flush()

	successAttempt := 0
	fields := make([]Field, setSize)
genFieldSet:
	for attempt := 0; successAttempt == 0 && attempt < attemptLimit; attempt++ {
		if attempt < 9999 {
			gt.MoveCursor(1, 2)
			gt.Printf("Attempt %d from %dma...", (attempt + 1), attemptLimit/10000)
			gt.Flush()
		} else if (attempt+1)%10000 == 0 {
			gt.MoveCursor(1, 2)
			gt.Printf("Attempt %dma from %dma...", (attempt+1)/10000, attemptLimit/10000)
			gt.Flush()
		}

		// reset all
		var (
			numMaxOfMinUsePlusOneVal          int
			maxOfFieldsWithNeighboringNumsVal int
		)
		{
			for fIdx := 0; fIdx < setSize; fIdx++ {
				fields[fIdx].Reset(int(maxOfNum))
			}
			for i := 0; i < len(numUses); i++ {
				numUses[i] = 0
			}
			numMaxOfMinUsePlusOneVal = numMaxOfMinUsePlusOne
			maxOfFieldsWithNeighboringNumsVal = maxOfFieldsWithNeighboringNums
		}

		// fill all fields
		for fIdx := 0; fIdx < setSize; fIdx++ {
			genFieldAttempt := 0
		genField:
			genFieldAttempt++
			if genFieldAttempt > genFieldAttemptLimit {
				continue genFieldSet
			}

			// fill field.
			genNumAttempt := 0
			withNeighboringNums := false
			countOfMinUsePlusOne := 0
			for c := 0; c != fieldSize; {
				genNumAttempt++
				if genNumAttempt > genNumAttemptLimit {
					continue genFieldSet
				}
				num := int(s.NextRandom(maxOfNum))
				isNeighbourNum := fields[fIdx].IsNeighbour(num + 1)
				if isNeighbourNum && maxOfFieldsWithNeighboringNumsVal == 0 {
					continue
				}
				if fields[fIdx].Set(num + 1) {
					l := numUses[num] + 1
					switch {
					case l <= numMinUse:
						numUses[num] = l
						if isNeighbourNum {
							withNeighboringNums = true
						}
					case l == numMinUse+1 && numMaxOfMinUsePlusOneVal > 0:
						numUses[num] = l
						numMaxOfMinUsePlusOneVal--
						countOfMinUsePlusOne++
						if isNeighbourNum {
							withNeighboringNums = true
						}
					default:
						fields[fIdx].Unset(num + 1)
						continue
					}
					c++
				}
			}

			if SimilarityDegreeSlice(fields[:fIdx+1]) > opts.FieldSimilarityDegree() {
				for _, num := range fields[fIdx].Nums() {
					numUses[num-1] = numUses[num-1] - 1
				}
				fields[fIdx].Reset(int(maxOfNum))
				numMaxOfMinUsePlusOneVal = numMaxOfMinUsePlusOneVal + countOfMinUsePlusOne
				goto genField
			}

			if withNeighboringNums {
				maxOfFieldsWithNeighboringNumsVal--
			}

		}

		// verify constraints
		if !MagicCheckF(fields) {
			continue
		}

		successAttempt = attempt + 1
	}

	if successAttempt != 0 {
		fmt.Printf("\ndone at %d attempt!\n\n", successAttempt)
		result := make([][]int32, setSize)
		for t := 0; t < setSize; t++ {
			result[t] = fields[t].Nums()
		}
		return result
	}

	fmt.Printf("\n===============\nAttempts limit reached. No results.\n")
	return nil
}

func (s *Seed) GenerateField2(opts intl_opts.Options) Field2 {
	f2 := []Field{{}, {}}
	f2[0].Reset(int(opts.MaxOfNum()))
	f2[1].Reset(int(opts.MaxOfNum()))

	for i := 0; i < len(f2); i++ {
		c := 0
		for c != opts.FieldSize() {
			n := s.NextRandom(opts.MaxOfNum()) + 1
			if f2[i].Set(int(n)) {
				c++
			}
		}
	}

	return Field2{F1: f2[0], F2: f2[1]}
}

func (s *Seed) Field2RandomStream(ctx context.Context, opts intl_opts.Options) <-chan Field2 {
	ch := make(chan Field2)

	go func() {
		f2 := []Field{{}, {}}
		maxOfNum := int(opts.MaxOfNum())
		for {
			select {
			case <-ctx.Done():
				close(ch)
				return
			default:
				f2[0].Reset(maxOfNum)
				f2[1].Reset(maxOfNum)
				for i := 0; i < len(f2); i++ {
					c := 0
					for c != opts.FieldSize() {
						n := s.NextRandom(opts.MaxOfNum()) + 1
						if f2[i].Set(int(n)) {
							c++
						}
					}
				}

				f := Field2{
					F1: f2[0],
					F2: f2[1],
				}

				select {
				case ch <- f:
				case <-ctx.Done():
					close(ch)
					return
				}
			}
		}
	}()

	return ch
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
