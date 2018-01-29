package main

import (
	"fmt"
	"strings"
)

func main() {
	stringex := "parangaricutirimicuaro is a spanish tongue twister"
	fmt.Println("In go the function printf",
		"understands % verbs, as oposed to println")
	fmt.Printf("\n%T, %v\n", stringex, stringex)

	//playing with strings package

	fmt.Println(strings.Title(stringex))
	fmt.Println()

	// comparing strings

	comparingex1 := "ComSci 309"
	comparingex2 := "comsci 309"
	fmt.Println("Comparing in go uses simple == comparison\n",
		"The strings are equal:", comparingex1 == comparingex2)
	fmt.Println()

	fmt.Println("However using the \"strings\" method for non case sensitive")
	fmt.Println("The strings are equal:", strings.EqualFold(comparingex1, comparingex2))

}
