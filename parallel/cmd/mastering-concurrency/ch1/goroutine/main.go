package main

import (
	"fmt"
	"runtime"
	"strconv"
	"time"
)

var (
	ch = make(chan string)
)

func main() {
	ch <- donothing()
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
