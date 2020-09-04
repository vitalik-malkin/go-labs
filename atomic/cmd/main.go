package main

import (
	"fmt"
	"runtime"
)

func main() {

	state := uint64(1)

	heldID1 := (state >> 1) + 1
	state = (heldID1 << 1) | (state & 1)

	heldID2 := (state >> 1) + 1
	state = (heldID2 << 1) | (state & 1)

	heldID := (state >> 1)

	//	state |= 7

	f1 := state&1 == 1

	fmt.Printf("%v %v, %v %v", f1, heldID1, heldID2, heldID)

	runtime.Breakpoint()

	state = state >> 63

	runtime.Breakpoint()

}
