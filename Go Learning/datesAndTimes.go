package main

import (
	"fmt"
	"time"
)

func main() {
	consofTime := time.Now()
	fmt.Printf("When did Raj last said something that made someone feel unconfortable?\n")
	var stringy string
	fmt.Scanln(&stringy)
	fmt.Printf("Last time was %s\n", consofTime)
}
