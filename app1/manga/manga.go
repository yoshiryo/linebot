package manga

import (
	"fmt"
	"time"

	"github.com/yoshiryo/linebot/app1/db"
	"github.com/yoshiryo/linebot/app1/model"

	_ "github.com/go-sql-driver/mysql"
)

func AddManga(manga string) string {
	db, err := db.SqlConnect()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	error := db.Create(&model.Mangas{
		Name:     manga,
		UpdateAt: getDate(),
	}).Error
	if error != nil {
		fmt.Println(error)
	}
	return "漫画を追加したよ！"
}

func getDate() string {
	const layout = "2006-01-02 15:04:05"
	now := time.Now()
	return now.Format(layout)
}
