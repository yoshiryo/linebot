package main

import (
	"os"

	"github.com/yoshiryo/linebot/app2/server"
)

func main() {
	ChanellSecret := os.Getenv("ChanellSecret")
	ChanellToken := os.Getenv("ChanellToken")
	server.Connect(ChanellSecret, ChanellToken)
}
