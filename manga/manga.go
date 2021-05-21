package manga

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
)

func add_Manga(manga string) string {
	t := getReleaseMangeTime(manga)
	db, err := sqlConnect()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	error := db.Create(&Users{
		Name:         manga,
		release_time: t,
		UpdateAt:     getDate(),
	}).Error
	if error != nil {
		fmt.Println(error)
	}
	return ("追加しました")
}

func getReleaseMangeTime(manga string) string {
	url := ""

	// Getリクエスト
	res, _ := http.Get(url)
	defer res.Body.Close()

	// 読み取り
	buf, _ := ioutil.ReadAll(res.Body)

	// 文字コード判定
	det := chardet.NewTextDetector()
	detRslt, _ := det.DetectBest(buf)
	// => EUC-JP

	// 文字コード変換
	bReader := bytes.NewReader(buf)
	reader, _ := charset.NewReaderLabel(detRslt.Charset, bReader)

	// HTMLパース
	doc, _ := goquery.NewDocumentFromReader(reader)

	// titleを抜き出し
	t := doc.Find(".time").Text()
	return t
}

func getDate() string {
	const layout = "2006-01-02 15:04:05"
	now := time.Now()
	return now.Format(layout)
}

// SQLConnect DB接続
func sqlConnect() (database *gorm.DB, err error) {
	DBMS := "mysql"
	USER := "go_example"
	PASS := "12345!"
	PROTOCOL := "tcp(localhost:3306)"
	DBNAME := "go_example"

	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?charset=utf8&parseTime=true&loc=Asia%2FTokyo"
	return gorm.Open(DBMS, CONNECT)
}

type Users struct {
	ID           int
	Name         string `json:"name"`
	release_time int
	UpdateAt     string `json:"updateAt" sql:"not null;type:date"`
}
