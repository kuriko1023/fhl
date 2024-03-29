package main

import (
	"encoding/json"
	"os"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var Config struct {
	Port      int    `json:"port"`
	AppID     string `json:"appid"`
	AppSecret string `json:"appsecret"`
	Debug     bool   `json:"debug"`

	AllowOrigin string `json:"allow_origin"`

	ArticleCache int `json:"article_cache"`
}

var db *sql.DB

func main() {
	// 读取配置
	configPath := os.Getenv("CONFIG")
	if configPath == "" {
		configPath = "config.json"
	}
	content, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(content, &Config); err != nil {
		panic(err)
	}

	initDataset()
	/*var y bool
	var a, b int
	y, a, b = lookupText([]string{"悠哉悠哉", "辗转反侧"})
	fmt.Println(y, getArticle(a).Content[b:])
	y, a, b = lookupText([]string{"悠哉游哉", "辗转反侧"})
	fmt.Println(y, getArticle(a).Content[b:])
	y, a, b = lookupText([]string{"辗转反侧", "呜呜"})
	fmt.Println(y, a, b)
	y, a, b = lookupText([]string{"梳洗罢", "独倚望江楼"})
	fmt.Println(y, getArticle(a).Content[b:])
	y, a, b = lookupText([]string{"梳洗黑", "独倚望江楼"})
	fmt.Println(y, getArticle(a).Content[b:])
	y, a, b = lookupText([]string{"江南好"})
	fmt.Println(y, getArticle(a).Content[b:])
	y, a, b = lookupText([]string{"赤壁沉埋水流", "徒留名载空舟"})
	fmt.Println(y, getArticle(a).Content[b:])
	y, a, b = lookupText([]string{"赤壁沉埋水流", "徒留名空舟"})
	fmt.Println(y, a, b)

	for i := 0; i < 10; i++ {
		fmt.Println(generateA(5, 3))
	}
	for i := 0; i < 10; i++ {
		fmt.Println(generateB(i%5 + 5))
	}
	for i := 0; i < 10; i++ {
		fmt.Println(generateC(3, 10))
	}
	sl1, sl2 := generateD(8)
	for i := 0; i < 8; i++ {
		fmt.Println(sl1[i], sl2[i])
	}
	os.Stdin.Read(make([]byte, 1))
	return

	r := Room{Subject: &SubjectA{Word: "花"}}
	fmt.Println(r.Subject)

	var s Subject
	s = &SubjectB{}
	s.Parse("春花秋月何时了/3")
	fmt.Println(s)
	fmt.Println(s.Dump())

	s = &SubjectC{}
	s.Parse("古 梦 雁/长 舟 送 寄 事 神 不 生 西风 多少/1000010011")
	c, d := s.Answer("千古兴亡多少事", SideHost)
	fmt.Println(c, d)
	c, d = s.Answer("千古兴亡多少事", SideHost)
	fmt.Println(c, d)
	fmt.Println(s.Dump())

	s = &SubjectD{}
	s.Parse("万 书 今 凉 得 来 柳 欲/一片 丝 如此 孤 庭 细 舟 觉/00000000/00000000")
	c, d = s.Answer("孤帆一片日边来", SideHost)
	fmt.Println(c, d)
	c, d = s.Answer("孤蓬万里征", SideHost)
	fmt.Println(c, d)
	fmt.Println(s.Dump())
	os.Stdin.Read(make([]byte, 1))*/

	if db, err = SetUpDatabase(); err != nil {
		panic(err)
	}
	defer db.Close()

	SetUpHttp()
}
