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
	for i := 0; i < 10; i++ {
		generateB(3, 10)
	}
	/*sl1, sl2 := generateC(8)
	for i := 0; i < 8; i++ {
		fmt.Println(sl1[i], sl2[i])
	}*/
}
