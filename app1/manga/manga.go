package manga

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/yoshiryo/linebot/db"
	"github.com/yoshiryo/linebot/model"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
)

func AddManga(manga string) string {
	db, err := db.SqlConnect()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	error := db.Create(&model.Manga{
		Name:     manga,
		UpdateAt: getDate(),
	}).Error
	if error != nil {
		fmt.Println(error)
	}
	return ("追加しました")
}

func getMangeList() [][]string {
	d := getDates()
	url_list := []string{
		"https://sinkan.net/?group=Comic&action_top=true&start=" + d[0][0:10],
		"https://sinkan.net/?group=Comic&action_top=true&start=" + d[1][0:10],
		"https://sinkan.net/?group=Comic&action_top=true&start=" + d[2][0:10]}
	var days_manga_list [][]string
	for _, url := range url_list {
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
		var manga_list []string
		for i := 1; i < 201; i++ {
			selector := "#today > div.autopagerize_page_element > table:nth-child(" + strconv.Itoa(i) + ") > tbody > tr > td.i_info > div.i_title > a"
			t := doc.Find(selector).Text()
			if len(t) > 0 {
				manga_list = append(manga_list, t)
			} else {
				break
			}
		}
		days_manga_list = append(days_manga_list, manga_list)
	}
	return days_manga_list
}

func getDate() string {
	const layout = "2006-01-02 15:04:05"
	now := time.Now()
	return now.Format(layout)
}

func getDates() []string {
	const layout = "2006-01-02 15:04:05"
	now := time.Now()
	tommorow := now.Add(24 * time.Hour)
	day_after_tomorrow := now.Add(48 * time.Hour)
	dates := []string{
		now.Format(layout),
		tommorow.Format(layout),
		day_after_tomorrow.Format(layout)}
	return dates
}
