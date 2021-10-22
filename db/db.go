package db

import "github.com/jinzhu/gorm"

func SqlConnect() (database *gorm.DB, err error) {
	DBMS := "mysql"
	USER := "ueoai"
	PASS := "ueoai0622"
	PROTOCOL := "tcp(localhost:3306)"
	DBNAME := "linebot"

	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?charset=utf8&parseTime=true&loc=Asia%2FTokyo"
	return gorm.Open(DBMS, CONNECT)
}
