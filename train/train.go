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
	rslt := doc.Find(".time").Text()
	rslt = rslt[strings.Index(rslt, "出発"):]
	rslt = strings.Replace(rslt, "出発", "", -1)

	rslt = sta_station + "  " + des_station + "  所要時間" + "\n" +
		rslt[:13] + "　" + rslt[13:18] +
		rslt[18:31] + "　" + rslt[31:36] +
		rslt[36:49] + "　" + rslt[49:54]
	return rslt
}

func InsertStation(sta_station, des_station, name string) string {
	db, err := db.SqlConnect()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

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
	db, err := db.SqlConnect()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	result := []*data.Stations{}
	error := db.Find(&result).Error
	if error != nil {
		fmt.Println(error)
	} else if len(result) == 0 {
		return "登録されてないよ！"
	}

	kekka := ""
	for i, user := range result {
		name := user.Name
		first_station := user.First_Station
		second_station := user.Second_Station
		if i != len(result)-1 {
			kekka += "名前" + strconv.Itoa(i+1) + "：" + name + "\n" +
				"発車駅：" + first_station + "\n" +
				"到着駅：" + second_station + "\n" +
				"\n"
		} else {
			kekka += "名前" + strconv.Itoa(i+1) + "：" + name + "\n" +
				"発車駅：" + first_station + "\n" +
				"到着駅：" + second_station
		}

	}
	return kekka
}

func UseRoute(name string) string {
	db, err := db.SqlConnect()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	result := []*data.Stations{}
	error := db.Where("name = ?", name).First(&result).Error
	if error != nil {
		fmt.Println(error)
	} else if len(result) == 0 {
		return "登録されてないよ！"
	}
	first_station := result[0].First_Station
	second_station := result[0].Second_Station
	return GetTrainTime(first_station, second_station)
}
