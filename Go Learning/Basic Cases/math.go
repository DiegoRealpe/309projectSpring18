package main

import (
	"fmt"
	"math/big"
)

func main() {

	i1, i2, i3 := 12, 45, 68
	intSum := i1 + i2 + i3
	fmt.Println("Integer sum: ", intSum)

	var b1, b2, b3, bsum big.Float

	b1.SetFloat64(33.8)
	b1.
		b2.SetFloat64(19.1)
	b3.SetFloat64(3.56)
	bsum.Add(&b1, &b2).Add(&bsum, &b3)
	fmt.Printf("Float sum: %.20g", &bsum)
}
