package main

import (
	"context"
	"fmt"
	"log"
	"os"

	intl_opts "github.com/vitalik-malkin/go-labs/randoom/internal/options"
	intl_seed "github.com/vitalik-malkin/go-labs/randoom/internal/seed"
)

type findMatchR struct {
	SessionNum           int
	TNum                 int
	SeedOffsetResetCount int
}

func main() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Printf("error: %v", err)
		}
	}()

	// m, err := findMatch(40000)
	// if err != nil {
	// 	panic(err)
	// }
	// for _, i := range m {
	// 	fmt.Printf("%+v\n", i)
	// }
	// os.Exit(0)

	opts := intl_opts.Load()

	opts1 := opts
	opts1.SetFieldSetSize(opts.FieldSetSize() / 2)
	opts1.SetGenFieldSetAttemptLimit(opts.GenFieldSetAttemptLimit() / 2)
	opts2 := opts
	opts2.SetFieldSetSize(opts.FieldSetSize() - opts1.FieldSetSize())
	opts2.SetGenFieldSetAttemptLimit(opts.GenFieldSetAttemptLimit() - opts1.GenFieldSetAttemptLimit())

	seed1, err := intl_seed.Load(opts1)
	if err != nil {
		log.Fatal(err)
	}
	seed2, err := intl_seed.Load(opts2)
	if err != nil {
		log.Fatal(err)
	}

	var fieldSetGen1, fieldSetGen2 func(intl_opts.Options) [][]int32

	switch opts.GeneratorVersion() {
	case 2:
		fieldSetGen1, fieldSetGen2 = seed1.NextRandomFieldSetV2, seed2.NextRandomFieldSetV2
	default:
		fieldSetGen1, fieldSetGen2 = seed1.NextRandomFieldSet, seed2.NextRandomFieldSet
	}

	var fieldSet [][]int32

	fieldSet1 := fieldSetGen1(opts1)
	fieldSet2 := fieldSetGen2(opts2)

	if fieldSet1 != nil && fieldSet2 != nil {
		fieldSet = make([][]int32, opts.FieldSetSize())

		c := 0
		for _, f := range fieldSet1 {
			fieldSet[c] = f
			c = c + 2
		}

		c = 1
		for _, f := range fieldSet2 {
			fieldSet[c] = f
			c = c + 2
		}
	}

	if fieldSet != nil {
		numStat := make([]int, opts.MaxOfNum())
		fmt.Print("Sets:\n")
		for i := 0; i < len(fieldSet); i++ {
			if i == 0 {
				fmt.Print("1) ")
			} else {
				if i%2 == 0 {
					fmt.Print("\n")
					fmt.Printf("%d) ", i/2+1)
				} else {
					fmt.Print("; ")
				}
			}
			fmt.Printf("%v", fieldSet[i])

			for u := 0; u < len(fieldSet[i]); u++ {
				numStat[fieldSet[i][u]-1] = numStat[fieldSet[i][u]-1] + 1
			}

		}

		numStatByCount := make(map[int][]int)
		for num, count := range numStat {
			numStatByCount[count] = append(numStatByCount[count], num+1)
		}
		fmt.Printf("\nStat (num repeats):\n%+v", numStatByCount)

		os.Exit(0)
	} else {
		fmt.Print("No sets.\n")
		os.Exit(1)
	}
}

func findMatch(maxTPerSession int) ([]findMatchR, error) {
	opts := intl_opts.Load()

	targetSeed, err := intl_seed.Load(opts)
	if err != nil {
		return nil, err
	}

	tSeed, err := intl_seed.Load(opts)
	if err != nil {
		return nil, err
	}

	sessionNum := 0
	resCount := 100
	res := make([]findMatchR, 0)
	for len(res) != resCount {
		sessionNum++
		var target intl_seed.Field2
		target = targetSeed.GenerateField2(opts)
		fmt.Printf("\nSESS %d, %s. Begin\n", sessionNum, target)

		sessionW := false
		sessionWF := intl_seed.Field2{}
		tNum := 0
		ctx, cancel := context.WithCancel(context.Background())
		tChan := tSeed.Field2RandomStream(ctx, opts)
		for t := range tChan {
			tNum++
			fmt.Printf("SESS %d. %d: %s\n", sessionNum, tNum, t)

			if target.Eq(t) {
				sessionW = true
				sessionWF = t
				cancel()
				break
			}

			if tNum == maxTPerSession {
				cancel()
				break
			}
		}
		cancel()

		if sessionW {
			res = append(res, findMatchR{SessionNum: sessionNum, TNum: tNum, SeedOffsetResetCount: tSeed.ResetCount()})
			fmt.Printf("SESS %d, %s. WIN, %s, %d. End\n", sessionNum, target, sessionWF, tNum)
			sessionNum = 0
			tSeed.Reset()
		} else {
			fmt.Printf("SESS %d, %s. End\n", sessionNum, target)
		}

	}

	return res, nil
}
