package seed

import (
	"bufio"
	cr "crypto/rand"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"strconv"
	"unsafe"

	intl_opts "github.com/vitalik-malkin/go-labs/randoom/internal/options"
)

type Seed struct {
	buf []byte

	bufReadOffset          int
	bufReadOffsetResetNeed bool
}

func Load(opts intl_opts.Options) (*Seed, error) {
	seedNums, err := loadSeedNums(opts)
	if err != nil {
		return nil, err
	}

	tmpSeed, err := load(seedNums)
	if err != nil {
		return nil, err
	}

	randomizeSeedNums(seedNums, tmpSeed)

	return load(seedNums)
}

func loadSeedNums(opts intl_opts.Options) ([]int32, error) {
	file, err := os.Open(opts.SeedFile())
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scan := bufio.NewScanner(file)
	scan.Split(bufio.ScanLines)

	seedNums := []int32{}
	for scan.Scan() {
		var text = scan.Text()
		if len(text) == 0 {
			continue
		}

		var seedNum int
		seedNum, err = strconv.Atoi(text)
		if err != nil {
			return nil, err
		}
		if seedNum < 0 || seedNum > math.MaxInt32 {
			return nil, fmt.Errorf("seed num %d is out of valid range;", seedNum)
		}
		seedNums = append(seedNums, int32(seedNum))
	}

	return seedNums, nil
}

func load(seedNums []int32) (*Seed, error) {
	buf := []byte{}
	for _, seedNum := range seedNums {
		num := int32(seedNum)
		numPtr := unsafe.Pointer(&num)
		numBytes := (*[4]byte)(numPtr)
		pZero := false
		for i := 0; i < 4; i++ {
			if numBytes[i] == 0 && pZero {
				continue
			}
			buf = append(buf, numBytes[i])
			pZero = numBytes[i] == 0
		}
	}

	res := &Seed{buf: buf}
	res.resetOffset()

	return res, nil
}

func randomizeSeedNums(seedNums []int32, s *Seed) {
	l := int32(len(seedNums))
	for i := int32(0); i < l/2; i++ {
		n1, n2 := s.NextRandom(l), s.NextRandom(l)
		seedNums[n1], seedNums[n2] = seedNums[n2], seedNums[n1]
	}
}

func (s *Seed) resetOffset() {
	if len(s.buf) < 2 {
		s.bufReadOffset = 0
		return
	}

	bufLen := big.NewInt(int64(len(s.buf)))
	newOffsetBig, err := cr.Int(cr.Reader, bufLen)
	if err != nil {
		log.Fatalf("error while generating random offset; err: %v", err)
		return
	}

	s.bufReadOffset = int(newOffsetBig.Int64())
	s.bufReadOffsetResetNeed = false
}

func (s *Seed) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	if len(p) > len(s.buf) {
		return 0, fmt.Errorf("len %d requested to read is too large; max is %d", len(p), len(s.buf))
	}

	if s.bufReadOffsetResetNeed {
		s.resetOffset()
	}

	var (
		srcPos = s.bufReadOffset
		dstPos = 0
		dstLen = len(p)
	)
	for {
		lenCopied := copy(p[dstPos:], s.buf[srcPos:])
		dstPos = dstPos + lenCopied
		if dstLen == dstPos {
			break
		}
		srcPos = 0
	}
	if len(s.buf)-s.bufReadOffset > dstLen {
		s.bufReadOffset = s.bufReadOffset + dstLen
	} else {
		s.bufReadOffsetResetNeed = true
	}

	return dstLen, nil
}
