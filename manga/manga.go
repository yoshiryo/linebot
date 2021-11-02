package manga

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"url/data"
	"url/db"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
)

func AddManga(manga string) string {
	release_time := getReleaseMangeTime(manga)
	db, err := db.SqlConnect()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	error := db.Create(&data.Manga{
		Name:     manga,
		Release:  release_time,
		UpdateAt: getDate(),
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
	return doc.Find(".time").Text()
}

func getDate() string {
	const layout = "2006-01-02 15:04:05"
	now := time.Now()
	return now.Format(layout)
}

func InfoManga() string {

}
