package train

// 利用したい外部のコードを読み込む
import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/yoshiryo/linebot/app1/db"
	"github.com/yoshiryo/linebot/app1/model"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
)

func GetTrainTime(sta_station, des_station string) string {
	encode_sta := url.QueryEscape(sta_station) //URL encode
	encode_des := url.QueryEscape(des_station) //URL encode
	date := getDate()                          //2006-01-02 15:04:05
	year := date[:4]
	month := date[5:7]
	day := date[8:10]
	hour := date[11:13]
	min1 := date[14:15]
	min2 := date[15:16]
	url := "https://transit.yahoo.co.jp/search/result?from=" + encode_sta + "&to=" + encode_des + "&fromgid=&togid=&flatlon=&tlatlon=&via=&viacode=&y=" + year + "&m=" + month + "&d=" + day + "&hh=" + hour + "&m1=" + min1 + "&m2=" + min2 + "&type=1&ticket=ic&expkind=1&userpass=1&ws=3&s=0&al=1&shin=1&ex=1&hb=1&lb=1&sr=1&"
	fmt.Println(date)
	fmt.Println(url)
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
	result := doc.Find(".time").Text()
	//取得した文字列を適切な文に変更
	result = result[strings.Index(result, "出発"):]
	result = strings.Replace(result, "出発", "", -1)

	//返信のための文字列を作成
	result = sta_station + "  " + des_station + "  所要時間" + "\n" +
		result[:13] + "　" + result[13:18] +
		result[18:31] + "　" + result[31:36] +
		result[36:49] + "　" + result[49:54] + "\n" +
		url
	return result
}

func InsertStation(sta_station, des_station, name string) string {
	//mysqlとの接続開始
	db, err := db.SqlConnect()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	//mysqlにinsertする
	error := db.Create(&model.Stations{
		Name:           name,
		First_Station:  sta_station,
		Second_Station: des_station,
	}).Error
	if error != nil {
		fmt.Println(error)
	}
	return "駅を追加したよ！"
}

func GetStation() string {
	//mysqlとの接続開始
	db, err := db.SqlConnect()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	//
	db_result := []*model.Stations{}
	//select * from stationsと同義
	error := db.Find(&db_result).Error
	if error != nil {
		fmt.Println(error)
	} else if len(db_result) == 0 {
		return "駅が登録されてないよ！"
	}

	result := ""
	//selectで取得したものを一つずつ適切な形に処理
	for i, user := range db_result {
		name := user.Name
		first_station := user.First_Station
		second_station := user.Second_Station
		//返信のための文字列を作成
		if i != len(result)-1 {
			result += "名前" + strconv.Itoa(i+1) + "：" + name + "\n" +
				"発車駅：" + first_station + "\n" +
				"到着駅：" + second_station + "\n" +
				"\n"
		} else {
			result += "名前" + strconv.Itoa(i+1) + "：" + name + "\n" +
				"発車駅：" + first_station + "\n" +
				"到着駅：" + second_station
		}

	}
	return result
}

func UseRoute(name string) string {
	//mysqlとの接続開始
	db, err := db.SqlConnect()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	//
	db_result := []*model.Stations{}
	//select * from stations where name = "?"と同義
	error := db.Where("name = ?", name).First(&db_result).Error
	if error != nil {
		fmt.Println(error)
	} else if len(db_result) == 0 {
		return "登録されてないよ！"
	}
	first_station := db_result[0].First_Station
	second_station := db_result[0].Second_Station
	//GetTrainTime関数を使って時刻を割り出してreturn
	return GetTrainTime(first_station, second_station)
}

func getDate() string {
	const layout = "2006-01-02 15:04:05"
	now := time.Now()
	nowUTC := now.UTC()
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	nowJST := nowUTC.In(jst)
	return nowJST.Format(layout)
}
