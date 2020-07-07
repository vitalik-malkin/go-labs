package main

import (
	"fmt"
	"math"
	"math/rand"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	ch = make(chan string)
)

func main() {
	l1, err := selectTCPListener()
	l2, err := selectTCPListener()
	fmt.Printf("%v, %v, %v", l1, l2, err)
	// ...
	fmt.Println("done")
}

func donothing() string {
	return "nothing"
}

func ex1() {
	runtime.GOMAXPROCS(2)
	fmt.Printf("cpu: %v\n", runtime.NumCPU())
	fmt.Printf("threads: %v\n", runtime.GOMAXPROCS(0))

	go func() {
		for i := 0; i < 100; i++ {
			tstamp := strconv.FormatInt(time.Now().UnixNano(), 10)
			fmt.Printf("%v, %v\n", i, tstamp)
			time.Sleep(time.Millisecond * 10)
		}
	}()

	go func() {
		for i := 100; i < 200; i++ {
			tstamp := strconv.FormatInt(time.Now().UnixNano(), 10)
			fmt.Printf("%v, %v\n", i, tstamp)
			time.Sleep(time.Millisecond * 10)
		}
	}()

	runtime.Gosched()

	fmt.Println("done")
}

func selectTCPListener() (net.Listener, error) {
	const maxAttempts = 1
	const portNumRangeVolume = float64(45535)
	const minPortNum = 20000
	attemptNum := 0
	for {
		portNum := int(math.Ceil(rand.Float64()*portNumRangeVolume)) + minPortNum
		lis, err := net.Listen("tcp4", fmt.Sprintf(":%v", portNum))
		if err == nil {
			return lis, nil
		}
		if attemptNum++; isErrorAddressAlreadyInUse(err) && attemptNum < maxAttempts {
			continue
		}
		return nil, err
	}
}

func isErrorAddressAlreadyInUse(err error) bool {
	const windowsAddrInUserErrno = 10048
	opError, ok := err.(*net.OpError)
	if !ok {
		return false
	}
	syscallError, ok := opError.Err.(*os.SyscallError)
	if !ok {
		return false
	}
	syscallErrno, ok := syscallError.Err.(syscall.Errno)
	if !ok {
		return false
	}
	switch strings.ToUpper(runtime.GOOS) {
	case "WINDOWS":
		return syscallErrno == windowsAddrInUserErrno
	default:
		return syscallErrno == syscall.EADDRINUSE
	}
}
