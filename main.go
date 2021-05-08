package main

import "fmt"

func main() {
	fmt.Print("xx")
	a := [5]int{1, 2, 3, 4, 5}
	s := a[3:4:4]
	fmt.Println(s[0])
	fmt.Sprintf("abc%d", 123)
}
