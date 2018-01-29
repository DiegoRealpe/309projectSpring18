package main

import "fmt"

func main() {

	str1 := "The quick red fox"
	str2 := "jumped over"
	str3 := "the lazy brown dog."
	//aNumber := 42
	//isTrue := isTrue

	strlen, err := fmt.Println(str1, str2, str3)

	if err == nil {
		fmt.Println("String length:", strlen)
		fmt.Println("Programming is fun!")
	}                                                             \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
}
