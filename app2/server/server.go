package server

import (
	"log"

	"github.com/yoshiryo/linebot/app2/manga"

	"github.com/line/line-bot-sdk-go/linebot"
)

func Connect(cs, ct string) {
	// LINE Botクライアント生成する
	// BOT にはチャネルシークレットとチャネルトークンを環境変数から読み込み引数に渡す
	bot, err := linebot.New(
		cs,
		ct,
	)
	// エラーに値があればログに出力し終了する
	if err != nil {
		log.Fatal(err)
	}
	// テキストメッセージを生成する
	message := linebot.NewTextMessage(manga.AlertMangeReleaseDay())
	// テキストメッセージを友達登録しているユーザー全員に配信する
	if _, err := bot.BroadcastMessage(message).Do(); err != nil {
		log.Fatal(err)
	}
}
