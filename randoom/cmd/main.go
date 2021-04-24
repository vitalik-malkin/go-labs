package main

import (
	"bufio"
	cr "crypto/rand"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"unsafe"
)

const (
	seedFile = "./seed.txt"
)

type seed []int

type Field struct {
}

func main() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Printf("error: %v", err)
		}
	}()

	s, err := loadSeed()
	if err != nil {
		log.Fatal(err)
	}

	fieldSet := s.nextRandom20FieldSet(4, 14)
	numStat := make([]int, 20)
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
	numStatMin, numStatMax := math.MaxInt32, math.MinInt32
	for _, v := range numStat {
		if numStatMin > v {
			numStatMin = v
		}
		if numStatMax < v {
			numStatMax = v
		}
	}
	fmt.Printf("\nStat: min=%d, max=%d", numStatMin, numStatMax)

	os.Exit(0)
}

func loadSeed() (seed, error) {
	file, err := os.Open(seedFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scan := bufio.NewScanner(file)
	scan.Split(bufio.ScanLines)

	r := []int{}
	for scan.Scan() {
		var n int
		n, err = strconv.Atoi(scan.Text())
		if err != nil {
			return nil, err
		}
		r = append(r, n)
	}

	return r, nil
}

func (s seed) Len() int {
	return len(s)
}

func (s seed) nextRandom() int {
	var buf [8]byte
	n, err := cr.Read(buf[:])
	if err != nil {
		log.Fatalf("error while populating buffer with random bytes: %v", err)
	} else if n != len(buf) {
		log.Fatal("invalid result of random function")
	}
	var num1, num2 int32

	bufPtr := unsafe.Pointer(&buf)

	num1 = (*(*int32)(bufPtr)) & int32(0x7fffffff)
	num2 = (*(*int32)(unsafe.Pointer(uintptr(bufPtr) + unsafe.Sizeof(num1)))) & 0x7fffffff
	var _ = num2

	return s[int(num1)%len(s)]
}

func (s seed) nextRandom20() int {
	return s.nextRandom()%20 + 1
}

func (s seed) nextRandom20FieldSet(fieldSize, setSize int) [][]int {

	numOverallRepeatLim := (fieldSize*setSize)/20 + 1
	numOverFieldRepeatLim :=
		func() int {
			const x = 1
			if x < numOverallRepeatLim {
				return x
			}
			return numOverallRepeatLim
		}()
	nextNumAttemptLim := (fieldSize * setSize)
	nextFieldAttemptLim := numOverallRepeatLim

	fieldSetAttempt := 0
nextFieldSet:
	fieldSetAttempt++
	fmt.Printf("===============\nNew field set (attempt %d)...\n", fieldSetAttempt)

	numStat := make([]int, 20)
	newField := make([]int, fieldSize)
	fieldSet := make([][]int, setSize)
	for j := 0; j < len(fieldSet); j++ {

		fieldAttempt := 0
	nextField:
		fieldAttempt++
		if fieldAttempt > nextFieldAttemptLim {
			goto nextFieldSet
		}

		for i := 0; i < len(newField); i++ {
			numAttempt := 0
		nextNum:
			numAttempt++
			if numAttempt > nextNumAttemptLim {
				if s.nextRandom()%2 == 0 {
					goto nextField
				}
				goto nextFieldSet
			}
			num := s.nextRandom20()
			if numStat[num-1] >= numOverallRepeatLim {
				goto nextNum
			}
			for y := 0; y < len(newField); y++ {
				sub := newField[y] - num
				if sub < 2 && sub > -2 {
					goto nextNum
				}
			}
			newField[i] = num
		}

		for y := 0; y < j; y++ {
			existingField := fieldSet[y]
			numMatchCount := 0
			for k := 0; k < len(newField); k++ {
				for l := 0; l < len(existingField); l++ {
					if newField[k] == existingField[l] {
						numMatchCount++
					}
				}
			}
			if numMatchCount > numOverFieldRepeatLim {
				goto nextField
			}
		}

		for _, num := range newField {
			numStat[num-1] = numStat[num-1] + 1
		}
		sort.Ints(newField)
		fieldSet[j] = newField
		newField = make([]int, fieldSize)
	}
	fmt.Print("done!\n\n")

	return fieldSet
}
