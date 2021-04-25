package main

import (
	"fmt"
	"log"
	"math"
	"os"

	intl_opts "github.com/vitalik-malkin/go-labs/randoom/internal/options"
	intl_seed "github.com/vitalik-malkin/go-labs/randoom/internal/seed"
)

func main() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Printf("error: %v", err)
		}
	}()

	opts := intl_opts.Load()
	s, err := intl_seed.Load(opts)
	if err != nil {
		log.Fatal(err)
	}

	fieldSet := s.NextRandom20FieldSet(opts)
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
		numStatMinNum, numStatMaxNum := 0, 0
		numStatMin, numStatMax := math.MaxInt32, math.MinInt32
		for num, v := range numStat {
			if numStatMin > v {
				numStatMin = v
				numStatMinNum = num + 1
			}
			if numStatMax < v {
				numStatMax = v
				numStatMaxNum = num + 1
			}
		}
		fmt.Printf("\nStat (num repeats): min=%d (num %d), max=%d (num %d)", numStatMin, numStatMinNum, numStatMax, numStatMaxNum)

		os.Exit(0)
	} else {
		fmt.Print("No sets.\n")
		os.Exit(1)
	}
}
