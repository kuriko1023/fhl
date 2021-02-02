package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, world!")
/*
	initDataset()
	for i := 0; i < 10; i++ {
		fmt.Println(generateA(5, 3))
	}
	for i := 0; i < 10; i++ {
		fmt.Println(generateB(i % 5 + 5))
	}
	for i := 0; i < 10; i++ {
		fmt.Println(generateC(3, 10))
	}
	furtherInit()
	sl1, sl2 := generateD(8)
	for i := 0; i < 8; i++ {
		fmt.Println(sl1[i], sl2[i])
	}
*/

	r := Room{Subject: &SubjectA{Word: "花"}}
	fmt.Println(r.Subject)

	s := SubjectB{}
	s.Parse("春花秋月何时了 3")
	fmt.Println(s)
	fmt.Println(s.Dump())
}
