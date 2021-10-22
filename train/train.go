package train

// 利用したい外部のコードを読み込む
import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"url/data"
	"url/db"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
)

func GetTrainTime(sta_station, des_station string) string {
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
	result := doc.Find(".time").Text()
	//取得した文字列を適切な文に変更
	result = result[strings.Index(result, "出発"):]
	result = strings.Replace(result, "出発", "", -1)

	//返信のための文字列を作成
	result = sta_station + "  " + des_station + "  所要時間" + "\n" +
		result[:13] + "　" + result[13:18] +
		result[18:31] + "　" + result[31:36] +
		result[36:49] + "　" + result[49:54]
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
	error := db.Create(&data.Stations{
		Name:           name,
		First_Station:  sta_station,
		Second_Station: des_station,
	}).Error
	if error != nil {
		fmt.Println(error)
	}
	return "追加しました！"
}

func GetStation() string {
	//mysqlとの接続開始
	db, err := db.SqlConnect()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	//
	db_result := []*data.Stations{}
	//select * from stationsと同義
	error := db.Find(&db_result).Error
	if error != nil {
		fmt.Println(error)
	} else if len(db_result) == 0 {
		return "登録されてないよ！"
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
	db_result := []*data.Stations{}
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
