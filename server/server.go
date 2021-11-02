package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/yoshiryo/linebot/train"
)

const verifyToken = "00000000000000000000000000000000"

const helpMessage = `使い方
テキストメッセージ:
	"!manga"を送れば
	"!train [発車駅] [到着駅]"を送れば電車時間がわかるよ！
	"!addstation [発車駅] [到着駅] [ルート名]"を送れば電車のルート登録ができるよ！
	"!showroute"を送れば登録したルートがわかるよ！
	"!useroute [ルート名]"を送ればルート登録したもので電車時間がわかるよ！
	それ以外はやまびこを返すよ！
スタンプ:
	スタンプの情報を答えるよ！
それ以外:
	それ以外にはまだ対応してないよ！ごめんね...`

func Connect(cs, ct, port string) {
	bot, err := linebot.New(cs, ct)
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
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func getReplyMessage(event *linebot.Event, bot *linebot.Client) (replyMessage string) {
	// 来たメッセージの種類によって分岐する
	switch message := event.Message.(type) {
	// テキストメッセージが来たとき
	case *linebot.TextMessage:
		// !trainを含むとき
		if strings.Contains(message.Text, "!train") {
			words := strings.Fields(message.Text)
			if len(words) == 3 {
				return train.GetTrainTime(words[1], words[2])
			}
		} else if strings.Contains(message.Text, "!addstation") {
			words := strings.Fields(message.Text)
			if len(words) == 4 {
				return train.InsertStation(words[1], words[2], words[3])
			}
		} else if strings.Contains(message.Text, "!showroute") {
			return train.GetStation()
		} else if strings.Contains(message.Text, "!useroute") {
			words := strings.Fields(message.Text)
			if len(words) == 2 {
				return train.UseRoute(words[1])
			}
		} else if strings.Contains(message.Text, "!help") {
			return helpMessage
		}
		// !addmangaを含むとき
		/*
			else if strings.Contains(message.Text, "!addmanga") {
				words := strings.Fields(message.Text)
				if len(words) == 2 {
					return manga.Add_Manga(words[1])
				}
			}
		*/
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
