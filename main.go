package main

// 利用したい外部のコードを読み込む
import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
	_ "github.com/go-sql-driver/mysql"
    "github.com/jinzhu/gorm"
)

const verifyToken = "00000000000000000000000000000000"

// main関数は最初に呼び出されることが決まっている
func main() {
	// ランダムな数値を生成する際のシード値の設定
	rand.Seed(time.Now().UnixNano())

	ChanellSecret := "ef6c38e18536dcfe37ea3f7103dfbf3d"
	ChanellToken := "2/nj45Eap/xlWsQOBMuWQ52n14ZttH1jRa/+1b078ELyvT8Nris0Xjd82vfPr1ZbaVGE8aU0Qn2b/3gWNqsEPivCoavF1BcekmuovjK5rxlmNnF/ffzxLtbejPgBDGzwk8aPGqB1o7ttXlRnpxz1EQdB04t89/1O/w1cDnyilFU="
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

const helpMessage = `使い方
テキストメッセージ:
	"おみくじ"がメッセージに入ってれば今日の運勢を占うよ！
	それ以外はやまびこを返すよ！
スタンプ:
	スタンプの情報を答えるよ！
それ以外:
	それ以外にはまだ対応してないよ！ごめんね...`

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
		} else if strings.Contains(message.Text, "電車") {
			words := strings.Fields(message.Text)
			if len(words) == 3 {
				return getTrainTime(words[1], words[2])
			}
		} else if strings.Contains(message.Text, "addmanga") {
			words := strings.Fields(message.Text)
			if len(words) == 2 {
				return add_Manga(words[1])
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

// おみくじ結果の生成
func getFortune() string {
	oracles := map[int]string{
		0: "大吉",
		1: "中吉",
		2: "小吉",
		3: "末吉",
		4: "吉",
		5: "凶",
		6: "末凶",
		7: "小凶",
		8: "中凶",
		9: "大凶",
	}
	// rand.Intn(10)は1～10のランダムな整数を返す
	return oracles[rand.Intn(10)]
}

func getTrainTime(sta_station, des_station string) string {
	url := "https://transit.yahoo.co.jp/search/result?flatlon=&fromgid=&from=" + sta_station + "&tlatlon=&togid=&to=" + des_station + "&viacode=&via=&viacode=&via=&viacode=&via=&y=&m=&d=&hh=&m2=&m1=&type=1&ticket=ic&expkind=1&ws=3&s=0&al=1&shin=1&ex=1&hb=1&lb=1&sr=1&kw=" + des_station

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

	// titleを抜き出し
	rslt := doc.Find(".time").Text()
	return rslt
}
add_Manga(manga string) string {
	db, err := sqlConnect()
	if err != nil {
		panic(err.Error())
	} 
	defer db.Close()
 
    error := db.Create(&Users{
        Name:          manga,
		release_time:  t,
        UpdateAt: getDate(),
    }).Error
    if error != nil {
        fmt.Println(error)
    }
	return ("追加しました") 
}


// SQLConnect DB接続
func sqlConnect() (database *gorm.DB, err error) {
	DBMS := "mysql"
	USER := "go_example"
	PASS := "12345!"
	PROTOCOL := "tcp(localhost:3306)"
	DBNAME := "go_example"

	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?charset=utf8&parseTime=true&loc=Asia%2FTokyo"
	return gorm.Open(DBMS, CONNECT)
}