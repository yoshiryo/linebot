package db

import "github.com/jinzhu/gorm"

func SqlConnect() (database *gorm.DB, err error) {
	DBMS := "mysql"
	USER := "go_example"
	PASS := "12345!"
	PROTOCOL := "tcp(localhost:3306)"
	DBNAME := "go_example"

	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?charset=utf8&parseTime=true&loc=Asia%2FTokyo"
	return gorm.Open(DBMS, CONNECT)
}
