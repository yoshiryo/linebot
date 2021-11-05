package main

import (
	"os"

	"github.com/yoshiryo/linebot/app1/server"
)

func main() {

	ChanellSecret := os.Getenv("ChanellSecret")
	ChanellToken := os.Getenv("ChanellToken")
	Port := "8080"
	server.Connect(ChanellSecret, ChanellToken, Port)
}
