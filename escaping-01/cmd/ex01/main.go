package main

import (
	"fmt"
)

func main() {
	x := "25"
	testA(&x)

}

func testA(a *string) {
	fmt.Printf("%v", a)
}

func f1() (string, bool) {
	var m map[string]string
	return m[""]

}
