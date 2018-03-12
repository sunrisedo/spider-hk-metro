package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	// "github.com/hu17889/go_spider/core/common/page"
	// "github.com/hu17889/go_spider/core/pipeline"
	// "github.com/hu17889/go_spider/core/spider"
	"github.com/gocolly/colly"
	py "github.com/mozillazg/go-pinyin"
)

func main() {
	// GetLine()
	dataCh := make(chan *LineData, 10)

	go DealData(dataCh)
	S8684Data(dataCh)

	var finish LineData
	finish.Key = "finish"
	dataCh <- &finish
	select {}
}

type LineInfo struct {
	Faresaver Faresaver
}

type Faresaver struct {
	Facilities []Facilities
}

type Facilities struct {
	LINE            string
	STATION_ID      string
	STATION_NAME_TC string
	STATION_NAME_EN string
}

type Metro struct {
	Versions   string      `json:"versions"`
	City       string      `json:"city"`
	Map        string      `json:"map"`
	CityKey    string      `json:"city_key"`
	LineData   []LineData  `json:"lineData"`
	EntrysData interface{} `json:"entrysData"`
}

type LineData struct {
	LineId       string   `json:"lineId"`
	Key          string   `json:"key"`
	LineName     string   `json:"lineName"`
	LineDiscribe string   `json:"lineDiscribe"`
	GisKey       string   `json:"gisKey"`
	Color        string   `json:"color"`
	Stages       []Stages `json:"stages"`
}

type Stages struct {
	Name     string   `json:"name"`
	EName    string   `json:"eName"`
	Index    int      `json:"index"`
	Code     string   `json:"code"`
	Position Position `json:"position"`
	RunTime  RunTime  `json:"runTime"`
}

type Position struct {
	Lat   float64 `json:"lat"`
	Lng   float64 `json:"lng"`
	Bdlat float64 `json:"bdlat"`
	Bdlng float64 `json:"bdlng"`
}

type RunTime struct {
	StartOne string `json:"startOne"`
	StartTow string `json:"startTow"`
	EndOne   string `json:"endOne"`
	EndTow   string `json:"endTow"`
}

func GetLine() {

	client := &http.Client{}
	reqest, _ := http.NewRequest("GET", "http://www.mtr.com.hk/st/data/fcdata_json.php", nil)
	reqest.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36")
	resp, errs := client.Do(reqest)
	if errs != nil {
		log.Println("get GetCoinInfo error:", errs)
		// SendError(&DingSend{Type: "Get coincola.com ETH price", Msg: errs.Error()})
		return
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("get body error:", err)
	}

	// log.Println("data", string(data))
	var ret LineInfo
	if err := json.Unmarshal(data, &ret); err != nil {
		// log.Println("data", string(data))
	}
	// log.Println("ret.Faresaver", ret.Faresaver.Facilities)
	// http://www.mtr.com.hk/alert/ryg_line_status.xml
	// .mtrimd_idicon.mtrimd_idicon--TWL {
	// 	background-color: #ff0000 !important; }
	//   .mtrimd_idicon.mtrimd_idicon--KTL {
	// 	background-color: #1a9431 !important; }
	//   .mtrimd_idicon.mtrimd_idicon--ISL {
	// 	background-color: #0860a8 !important; }
	//   .mtrimd_idicon.mtrimd_idicon--SIL {
	// 	background-color: #b5bd00 !important; }
	//   .mtrimd_idicon.mtrimd_idicon--TKL {
	// 	background-color: #6b208b !important; }
	//   .mtrimd_idicon.mtrimd_idicon--AEL {
	// 	background-color: #1c7670 !important; }
	//   .mtrimd_idicon.mtrimd_idicon--TCL {
	// 	background-color: #fe7f1d !important; }
	//   .mtrimd_idicon.mtrimd_idicon--DRL {
	// 	background-color: #f550a6 !important; }
	//   .mtrimd_idicon.mtrimd_idicon--EAL {
	// 	background-color: #5eb6e4 !important; }
	//   .mtrimd_idicon.mtrimd_idicon--MOL {
	// 	background-color: #9a3b26 !important; }
	//   .mtrimd_idicon.mtrimd_idicon--WRL {
	// 	background-color: #a40084 !important; }
	// http: //www.mtr.com.hk/ch/customer/css/servicestatus.css
	metro := Metro{Versions: "v0.0.1", City: "香港", CityKey: "hongkong", Map: "http://www.mtr.com.hk/ch/customer/st/index.php"}

	lineCnt := make(map[string]*LineData)
	a := py.Args{py.FIRST_LETTER, py.Heteronym, py.Separator}
	var lineIndex int
	for _, value := range ret.Faresaver.Facilities {
		lineInfo := strings.Split(value.LINE, "&")
		lineNames := py.LazyPinyin(lineInfo[0], a)
		// log.Println(lineNames, lineNamesss)
		key := fmt.Sprintf("%s%sL", strings.ToUpper(lineNames[0]), strings.ToUpper(lineNames[1]))
		if v, ok := lineCnt[key]; ok {
			var stages Stages
			stages.Index = len(v.Stages) + 1
			stages.Code = fmt.Sprintf("HK-%s-%d", v.LineId, stages.Index)
			stages.Name = value.STATION_NAME_TC
			stages.EName = value.STATION_NAME_EN
			v.Stages = append(v.Stages, stages)
			lineCnt[key] = v
			lineId, _ := strconv.Atoi(v.LineId)
			metro.LineData[lineId-1].Stages = v.Stages
		} else {
			log.Println(key)
			var lineData LineData
			var stages Stages
			lineIndex++
			lineData.LineId = strconv.Itoa(lineIndex)
			lineData.Key = key
			lineData.LineName = lineInfo[0]
			// lineData.Color = lineInfo[1][1 : len(lineInfo[1])-1]
			stages.Index = 1
			stages.Code = fmt.Sprintf("HK-%s-%d", lineData.LineId, stages.Index)
			stages.Name = value.STATION_NAME_TC
			stages.EName = value.STATION_NAME_EN
			lineData.Stages = append(lineData.Stages, stages)
			lineCnt[key] = &lineData
			metro.LineData = append(metro.LineData, lineData)
		}

	}

	obj, _ := json.Marshal(metro)
	CreateFile("./hongkong.json", obj)
	// log.Println("ret", metro)
}

func CreateFile(fileName string, data []byte) {
	if f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm); err != nil {
		log.Println("open file error:", err)
	} else if _, err := f.Write(data); err != nil {
		log.Println("write file error:", err)
	}
}

func ToPy(name string) string {
	lineNames := py.LazyPinyin(name, py.Args{py.FIRST_LETTER, py.Heteronym, py.Separator})
	// log.Println(lineNames, lineNamesss)
	var pys string
	for _, value := range lineNames {
		pys = fmt.Sprintf("%s%s", pys, strings.ToUpper(value))
	}
	return pys
}
func S8684Data(dataCh chan *LineData) {
	c := colly.NewCollector()
	lineIndex := 0

	c.OnHTML("div .ib-box ul li", func(e *colly.HTMLElement) {
		// log.Printf("Contains %s", e.Name)
		exit := false
		e.ForEach("font", func(_ int, el *colly.HTMLElement) {
			// log.Printf("ForEach %s", el.Text)
			// log.Printf("Contains %s", el.Text)
			if strings.Contains(el.Text, "未开通") {
				exit = true
				return
			}
		})
		if exit {
			// log.Printf("exit %v", exit)
			return
		}
		e.ForEach("a", func(_ int, el *colly.HTMLElement) {
			// log.Printf("ForEach %s", el.Text)
			// log.Printf("Contains %s", el.Text)
			if strings.Contains(el.Text, "轻铁") {
				return
			}

			style := el.Attr("style")
			styles := strings.Split(style, ":")
			if styles[0] != "background-color" {
				return
			}
			color := styles[1]
			link := el.Attr("href")

			key := ""
			if strings.Contains(el.Text, "-") {
				key = strings.Split(el.Text, "(")[0]
			} else {
				key = el.Text
			}
			lineKey := ToPy(key)

			lineIndex++
			var lineData LineData
			lineData.LineId = strconv.Itoa(lineIndex)
			lineData.Key = lineKey
			lineData.LineName = el.Text
			lineData.Color = color[:len(color)-1]
			// stages.Index = 1
			// stages.Code = fmt.Sprintf("HK-%s-%d", lineData.LineId, stages.Index)
			// stages.Name = value.STATION_NAME_TC
			// stages.EName = value.STATION_NAME_EN
			dataCh <- &lineData
			log.Printf("text:%s, key:%s, style:%s, link:%s", el.Text, lineData.Key, lineData.Color, link)
			e.Request.Visit(link)
		})
	})

	c.OnHTML("table tbody tr", func(e *colly.HTMLElement) {
		log.Println("table", e.Text)
		texts := strings.Split(e.Text, " ")
		log.Println("texts", texts, e.Attr("href"))
		// var stage Stages
		// type Stages struct {
		// 	Name     string   `json:"name"`
		// 	EName    string   `json:"eName"`
		// 	Index    int      `json:"index"`
		// 	Code     string   `json:"code"`
		// 	Position Position `json:"position"`
		// 	RunTime  RunTime  `json:"runTime"`
		// }
		// e.Request.Visit(link)
	})

	c.OnHTML("table tbody tr", func(e *colly.HTMLElement) {
		texts := strings.Split(e.Text, " ")
		log.Println("texts", texts)
		e.ForEach("td a", func(_ int, el *colly.HTMLElement) {
			link := el.Attr("href")
			log.Println("href", link)
			e.Request.Visit(link)
		})
		// var stage Stages
		// type Stages struct {
		// 	Name     string   `json:"name"`
		// 	EName    string   `json:"eName"`
		// 	Index    int      `json:"index"`
		// 	Code     string   `json:"code"`
		// 	Position Position `json:"position"`
		// 	RunTime  RunTime  `json:"runTime"`
		// }
	})

	// c.OnHTML("table tbody tr", func(e *colly.HTMLElement) {
	// 	texts := strings.Split(e.Text, " ")
	// 	log.Println("texts", texts)
	// 	e.ForEach("td a", func(_ int, el *colly.HTMLElement) {
	// 		link := el.Attr("href")
	// 		log.Println("href", link)
	// 		e.Request.Visit(link)
	// 	})
	// 	// var stage Stages
	// 	// type Stages struct {
	// 	// 	Name     string   `json:"name"`
	// 	// 	EName    string   `json:"eName"`
	// 	// 	Index    int      `json:"index"`
	// 	// 	Code     string   `json:"code"`
	// 	// 	Position Position `json:"position"`
	// 	// 	RunTime  RunTime  `json:"runTime"`
	// 	// }
	// })
	c.Visit("http://hkdt.8684.cn/")
}

func DealData(dataCh chan *LineData) {
	log.Println("Deal data start.")
	line := make(map[string]LineData)
	for {
		select {
		case data := <-dataCh:
			if data.Key == "finish" {
				// obj, _ := json.Marshal(metro)
				// CreateFile("./hongkong.json", obj)
				log.Println("Finish data start.", line)
				os.Exit(0)
				return
			}
			if _, ok := line[data.Key]; ok {
				line[data.Key+"A"] = *data
			} else {
				line[data.Key] = *data
			}
		}
	}
}

// func ToMetro() {
// 	metro := Metro{Versions: "v0.0.1", City: "香港", CityKey: "hongkong", Map: "http://www.mtr.com.hk/ch/customer/st/index.php"}

// 	lineCnt := make(map[string]*LineData)
// 	a := py.Args{py.FIRST_LETTER, py.Heteronym, py.Separator}
// 	var lineIndex int
// 	for _, value := range ret.Faresaver.Facilities {
// 		lineInfo := strings.Split(value.LINE, "&")
// 		lineNames := py.LazyPinyin(lineInfo[0], a)
// 		// log.Println(lineNames, lineNamesss)
// 		key := fmt.Sprintf("%s%sL", strings.ToUpper(lineNames[0]), strings.ToUpper(lineNames[1]))
// 		if v, ok := lineCnt[key]; ok {
// 			var stages Stages
// 			stages.Index = len(v.Stages) + 1
// 			stages.Code = fmt.Sprintf("HK-%s-%d", v.LineId, stages.Index)
// 			stages.Name = value.STATION_NAME_TC
// 			stages.EName = value.STATION_NAME_EN
// 			v.Stages = append(v.Stages, stages)
// 			lineCnt[key] = v
// 			lineId, _ := strconv.Atoi(v.LineId)
// 			metro.LineData[lineId-1].Stages = v.Stages
// 		} else {
// 			log.Println(key)
// 			var lineData LineData
// 			var stages Stages
// 			lineIndex++
// 			lineData.LineId = strconv.Itoa(lineIndex)
// 			lineData.Key = key
// 			lineData.LineName = lineInfo[0]
// 			// lineData.Color = lineInfo[1][1 : len(lineInfo[1])-1]
// 			stages.Index = 1
// 			stages.Code = fmt.Sprintf("HK-%s-%d", lineData.LineId, stages.Index)
// 			stages.Name = value.STATION_NAME_TC
// 			stages.EName = value.STATION_NAME_EN
// 			lineData.Stages = append(lineData.Stages, stages)
// 			lineCnt[key] = &lineData
// 			metro.LineData = append(metro.LineData, lineData)
// 		}

// 	}

// 	obj, _ := json.Marshal(metro)
// 	CreateFile("./hongkong.json", obj)
// }
