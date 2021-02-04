package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var Config struct {
	Port int `json:"port"`
	AppID string `json:"appid"`
	AppSecret string `json:"appsecret"`
}

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

	r := Room{Subject: &SubjectA{Word: "花"}}
	fmt.Println(r.Subject)

	var s Subject
	s = &SubjectB{}
	s.Parse("春花秋月何时了/3")
	fmt.Println(s)
	fmt.Println(s.Dump())

	s = &SubjectC{}
	s.Parse("古 梦 雁/长 舟 送 寄 事 神 不 生 西风 多少/1000010011")
	a, b := s.Answer("千古兴亡多少事", SideHost)
	fmt.Println(a, b)
	a, b = s.Answer("千古兴亡多少事", SideHost)
	fmt.Println(a, b)
	fmt.Println(s.Dump())

	s = &SubjectD{}
	s.Parse("万 书 今 凉 得 来 柳 欲/一片 丝 如此 孤 庭 细 舟 觉/00000000/00000000")
	a, b = s.Answer("孤帆一片日边来", SideHost)
	fmt.Println(a, b)
	a, b = s.Answer("孤蓬万里征", SideHost)
	fmt.Println(a, b)
	fmt.Println(s.Dump())
*/

	// 读取配置
	content, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(content, &Config)
	if err != nil {
		panic(err)
	}
	fmt.Println(Config)

	SetUpHttp()
}
