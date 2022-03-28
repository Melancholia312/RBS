package main

import (
	"flag"
	"fmt"
	"math/rand"
)

func main() {
	limit := flag.Int("limit", 0, "# of rand numbers")
	maxIn := flag.Int("max", 0, "the max value")
	flag.Parse()

	n := *limit
	max := *maxIn

	for i:= 0; i < n ; i++{
		fmt.Println(rand.Intn(max))
	}
}

