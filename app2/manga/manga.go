package manga

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"github.com/saintfish/chardet"
	"github.com/yoshiryo/linebot/app1/db"
	"github.com/yoshiryo/linebot/app2/model"
	"golang.org/x/net/html/charset"
)

func AlertMangeReleaseDay() string {
	if getMangeList() == "漫画が登録されてないよ！" {
		return "漫画が登録されてないよ！"
	} else {
		manga_list := strings.Split(getMangeList(), "　")
		manga_list = manga_list[:len(manga_list)-1]
		release_list := getMangeReleaseList()
		result := ""
		for i, day_release := range release_list {
			for _, manga := range manga_list {
				for _, release_manga := range day_release {
					if strings.Contains(release_manga, manga) {
						if i == 0 {
							result += "今日" + manga + "が発売だよ!" + "\n"
						} else if i == 1 {
							result += "明日" + manga + "が発売だよ!" + "\n"
						} else {
							result += "明後日" + manga + "が発売だよ!" + "\n"
						}
					}
				}
			}
		}
		return chop(result)
	}
}

func getMangeList() string {
	db, err := db.SqlConnect()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	db_result := []*model.Mangas{}
	//select * from stationsと同義
	error := db.Find(&db_result).Error
	if error != nil {
		fmt.Println(error)
	} else if len(db_result) == 0 {
		return "漫画が登録されてないよ！"
	}
	result := ""
	for _, user := range db_result {
		name := user.Name
		//返信のための文字列を作成
		result += name + "　"
	}
	return result
}

func getMangeReleaseList() [][]string {
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

func chop(s string) string {
	s = strings.TrimRight(s, "\n")
	if strings.HasSuffix(s, "\r") {
		s = strings.TrimRight(s, "\r")
	}
	return s
}
