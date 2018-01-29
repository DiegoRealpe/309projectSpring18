package main

import "fmt"

func main() {
	fmt.Printf("Computer Science 309 is ")
	var stringy string
	fmt.Scanln(&stringy)
	fmt.Println("Yes, Master. Computer Science is", stringy, "indeed!")
}
