package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"url/manga"
	"url/train"

	_ "github.com/go-sql-driver/mysql"
	"github.com/line/line-bot-sdk-go/linebot"
)

const verifyToken = "00000000000000000000000000000000"

const helpMessage = `使い方
テキストメッセージ:
	"おみくじ"がメッセージに入ってれば今日の運勢を占うよ！
	"!manga"を送れば
	"!train"を送れば電車時間がわかるよ！
	それ以外はやまびこを返すよ！
スタンプ:
	スタンプの情報を答えるよ！
それ以外:
	それ以外にはまだ対応してないよ！ごめんね...`
	
func main() {

	ChanellSecret := "10db005b09514726670d87966b5a443c"
	ChanellToken := "Jc5PdOo/OR+MwtnjzgIrslHYbCxU9TzNocjLcwEaIHS07EGIJQHTJi7Z0ppNAcY6Yh2hvRH+BsMxH8CSZwnVPSkHgUhOEk1RgdH9VEk6ksvfWviGyUzOCStt7FoDSt0dOoqwMYNlPwaY89e4vGy67QdB04t89/1O/w1cDnyilFU="
	Port := "8080"
	bot, err := linebot.New(ChanellSecret, ChanellToken)
	if err != nil {
		log.Fatal(err)
	}

	// LINEサーバからのリクエストを受け取ったときの処理
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Print("Accessed\n")

		// リクエストを扱いやすい形に変換する
		events, err := bot.ParseRequest(req)
		// 変換に失敗したとき
		if err != nil {
			fmt.Println("ParseRequest error:", err)
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}

		// LINEサーバから来たメッセージによってやる処理を変える
		for _, event := range events {
			// LINEサーバのverify時は何もしない
			if event.ReplyToken == verifyToken {
				return
			}

			// メッセージが来たとき
			if event.Type == linebot.EventTypeMessage {
				// 返信を生成する
				replyMessage := getReplyMessage(event, bot)
				// 生成した返信を送信する
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	})

	// LINEサーバからのリクエストを受け取る
	if err := http.ListenAndServe(":"+Port, nil); err != nil {
		log.Fatal(err)
	}
}

// 返信を生成する
func getReplyMessage(event *linebot.Event, bot *linebot.Client) (replyMessage string) {
	// 来たメッセージの種類によって分岐する
	switch message := event.Message.(type) {
	// テキストメッセージが来たとき
	case *linebot.TextMessage:
		// さらに「おみくじ」という文字列が含まれているとき
		if message.Text == "おみくじ" {
			// おみくじ結果を取得する
			return getFortune()
		} 
		else if strings.Contains(message.Text, "!train") {
			words := strings.Fields(message.Text)
			if len(words) == 3 {
				return train.GetTrainTime(words[1], words[2])
			}
		} 
		else if strings.Contains(message.Text, "!addmanga") {
			words := strings.Fields(message.Text)
			if len(words) == 2 {
				return manga.Add_Manga(words[1])
			}
		}
		// そうじゃないときはオウム返しする
		return message.Text
	// スタンプが来たとき
	case *linebot.StickerMessage:
		replyMessage := fmt.Sprintf("sticker id is %s, stickerResourceType is %s", message.StickerID, message.StickerResourceType)
		return replyMessage
	// どっちでもないとき
	default:
		return helpMessage
	}
}

