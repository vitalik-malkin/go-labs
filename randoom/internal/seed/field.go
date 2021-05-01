package seed

import (
	"fmt"
	"log"
)

const (
	fieldMaxSize = 64
	oneUint64    = uint64(1)
)

type Field2 struct {
	F1 Field
	F2 Field
}

type Field struct {
	size int
	bits uint64
}

func (f *Field2) Eq(o Field2) bool {
	return f.F1.Eq(o.F1) && f.F2.Eq(o.F2)
}

func (f Field2) String() string {
	return fmt.Sprintf("%s; %s", f.F1, f.F2)
}

func (f *Field) Eq(o Field) bool {
	return f.size == o.size && f.bits == o.bits
}

func (f *Field) Set(num int) bool {
	if num < 1 || num > f.size {
		log.Fatalf("num should be in range [%d; %d]", 1, f.size)
		return false
	}

	n := (oneUint64 << (num - 1))
	if f.bits&n == n {
		return false
	}
	f.bits = f.bits | n
	return true
}

func (f *Field) Unset(num int) {
	if num < 1 || num > f.size {
		log.Fatalf("num should be in range [%d; %d]", 1, f.size)
		return
	}

	f.bits = f.bits ^ (oneUint64 << (num - 1))
}

func (f *Field) IsNeighbour(num int) bool {
	if num < 1 || num > f.size {
		log.Fatalf("num should be in range [%d; %d]", 1, f.size)
		return false
	}

	num--
	switch num {
	case 0:
		leftNum := oneUint64 << 1
		return f.bits&leftNum == leftNum
	case f.size - 1:
		rightNum := (oneUint64 << (f.size - 2))
		return f.bits&rightNum == rightNum
	default:
		rightNum := (oneUint64 << (num - 1))
		leftNum := (oneUint64 << (num + 1))
		return f.bits&rightNum == rightNum || f.bits&leftNum == leftNum
	}

}

func (f *Field) SetCount() int {
	n := f.bits
	c := 0
	for i := 0; i < f.size; i++ {
		if n&oneUint64 == oneUint64 {
			c++
		}
		n = n >> 1
	}

	return c
}

func (f *Field) Reset(size int) {
	if size > fieldMaxSize || size < 1 {
		log.Fatalf("field size %d is out of range; max size is %d", size, fieldMaxSize)
	}

	f.bits = 0
	f.size = size
}

func (f *Field) Nums() []int32 {
	setNums := make([]int32, 0, f.size)
	n := f.bits
	s := int32(f.size)
	for i := int32(0); i < s; i++ {
		if n&oneUint64 == oneUint64 {
			setNums = append(setNums, i+1)
		}
		n = n >> 1
	}

	return setNums
}

func (f Field) String() string {
	return fmt.Sprintf("%v", f.Nums())
}

func SimilarityDegree(f1, f2 Field) int {
	size := min(f1.size, f2.size)
	x := f1.bits & f2.bits
	c := 0
	for i := 0; i < size; i++ {
		if x&oneUint64 == oneUint64 {
			c++
		}
		x = x >> 1
	}

	return c
}

func SimilarityDegreeSlice(f []Field) int {
	res := 0

	for y := 0; y < len(f); y++ {
		for j := y + 1; j < len(f); j++ {
			d := SimilarityDegree(f[y], f[j])
			if d > res {
				res = d
			}
		}
	}

	return res
}

func MagicCheckF(f []Field) bool {
	maxOfNum := 0
	for i := 0; i < len(f); i++ {
		if i == 0 {
			maxOfNum = f[i].size
			continue
		}
		if maxOfNum != f[i].size {
			log.Fatal("MagicCheckF: all passed fields should be same size")
			return false
		}
	}

	numTotalUses := 0
	numUses := make([]int, maxOfNum)
	for num := 0; num < maxOfNum; num++ {
		numMask := oneUint64 << num
		for j := 0; j < len(f); j++ {
			if f[j].bits&numMask == numMask {
				numUses[num] = numUses[num] + 1
				numTotalUses = numTotalUses + 1
			}
		}
	}

	minUse := numTotalUses / maxOfNum
	maxOfMinUsePlusOne := numTotalUses % maxOfNum
	for i := 0; i < len(numUses); i++ {
		switch numUses[i] {
		case minUse:
		case minUse + 1:
			if maxOfMinUsePlusOne == 0 {
				return false
			}
			maxOfMinUsePlusOne--
		default:
			return false
		}
	}

	return true
}

func newField(size int) Field {
	if size > fieldMaxSize || size < 1 {
		log.Fatalf("field size %d is out of range; max size is %d", size, fieldMaxSize)
	}

	return Field{
		size: size,
	}
}
