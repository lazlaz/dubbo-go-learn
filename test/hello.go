package main

import (
	"fmt"
	"rsc.io/quote"
)

func main() {
	fmt.Println("go mod test start...")
	defer fmt.Println("go mod test end!")
	fmt.Println(quote.Hello())
}
