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
	"sync"
	"unsafe"

	intl_opts "github.com/vitalik-malkin/go-labs/randoom/internal/options"
)

type Seed struct {
	buf []byte

	bufReadOffset          int
	bufReadOffsetResetNeed bool

	m sync.Mutex

	resetCount int

	simpleResetMode bool
}

func Load(opts intl_opts.Options) (*Seed, error) {
	seedNums, err := loadSeedNums(opts)
	if err != nil {
		return nil, err
	}

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
		for i := 0; i < 4; i++ {
			if numBytes[i] == 0 {
				continue
			}
			buf = append(buf, numBytes[i])
		}
	}

	res := &Seed{buf: buf}
	res.Reset()

	return res, nil
}

func (s *Seed) Reset() {
	s.reset(true)
}

func (s *Seed) reset(zeroCounter bool) {
	s.m.Lock()

	bufLen := int32(len(s.buf))
	if s.simpleResetMode {
		s.bufReadOffset = int(cryptoRandomInt32(bufLen))
	} else {
		randomizationSeed := &Seed{
			buf:                    make([]byte, bufLen),
			bufReadOffset:          int(cryptoRandomInt32(bufLen)),
			bufReadOffsetResetNeed: false,
			resetCount:             0,
			simpleResetMode:        true,
		}
		copy(randomizationSeed.buf, s.buf)

		for i := int32(0); i < bufLen/2; i++ {
			n1, n2 := randomizationSeed.NextRandom(bufLen), randomizationSeed.NextRandom(bufLen)
			s.buf[n1], s.buf[n2] = s.buf[n2], s.buf[n1]
		}

		s.bufReadOffset = int(randomizationSeed.NextRandom(bufLen))
	}

	s.bufReadOffsetResetNeed = false
	if zeroCounter {
		s.resetCount = 0
	} else {
		s.resetCount++
	}

	s.m.Unlock()
}

func (s *Seed) ResetCount() int {
	s.m.Lock()
	x := s.resetCount
	s.m.Unlock()
	return x
}

func (s *Seed) Read(p []byte) (n int, err error) {
	s.m.Lock()

	if len(p) == 0 {
		s.m.Unlock()
		return 0, nil
	}
	if len(p) > len(s.buf) {
		s.m.Unlock()
		return 0, fmt.Errorf("len %d requested to read is too large; max is %d", len(p), len(s.buf))
	}

	if s.bufReadOffsetResetNeed {
		s.m.Unlock()
		s.reset(false)
		s.m.Lock()
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

	s.m.Unlock()
	return dstLen, nil
}

func cryptoRandomInt32(max int32) int32 {
	bigInt := big.NewInt(int64(max))
	res, err := cr.Int(cr.Reader, bigInt)
	if err != nil {
		log.Fatalf("error while generating number using cryptorand; err: %v", err)
		return 0
	}

	return int32(res.Int64())
}
