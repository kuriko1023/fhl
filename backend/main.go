package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, world!")
	initDataset()
	for i := 0; i < 10; i++ {
		fmt.Println(generateA(5, 3))
	}
	for i := 0; i < 100; i++ {
		generateB(3, 10)
	}
}
