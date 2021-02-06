package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var Config struct {
	Port      int    `json:"port"`
	AppID     string `json:"appid"`
	AppSecret string `json:"appsecret"`
	Debug     bool   `json:"debug"`
}

var db *sql.DB

func main() {
	fmt.Println("Hello, world!")
/*
	initDataset()
	initSymDel()
	var a, b int
	a, b = checkSentenceInDataset([]string{"悠哉悠哉", "辗转反侧"})
	fmt.Println(a, b)
	a, b = checkSentenceInDataset([]string{"悠哉游哉", "辗转反侧"})
	fmt.Println(a, b)
	a, b = checkSentenceInDataset([]string{"辗转反侧", "呜呜呜呜"})
	fmt.Println(a, b)
	a, b = checkSentenceInDataset([]string{"梳洗罢", "独倚望江楼"})
	fmt.Println(a, b)
	a, b = checkSentenceInDataset([]string{"梳洗黑", "独倚望江楼"})
	fmt.Println(a, b)
	a, b = checkSentenceInDataset([]string{"江南好"})
	fmt.Println(a, b)
	return

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
	configPath := os.Getenv("CONFIG")
	if configPath == "" {
		configPath = "config.json"
	}
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(content, &Config); err != nil {
		panic(err)
	}

	if db, err = SetUpDatabase(); err != nil {
		panic(err)
	}
	defer db.Close()

	SetUpHttp()
}
