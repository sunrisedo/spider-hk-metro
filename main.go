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

	"github.com/hu17889/go_spider/core/common/page"
	"github.com/hu17889/go_spider/core/pipeline"
	"github.com/hu17889/go_spider/core/spider"
	py "github.com/mozillazg/go-pinyin"
)

type MyPageProcesser struct {
}

func NewMyPageProcesser() *MyPageProcesser {
	return &MyPageProcesser{}
}

// Parse html dom here and record the parse result that we want to Page.
// Package goquery (http://godoc.org/github.com/PuerkitoBio/goquery) is used to parse html.
// func (this *MyPageProcesser) Process(p *page.Page) {
// 	if !p.IsSucc() {
// 		println(p.Errormsg())
// 		return
// 	}

// 	query := p.GetHtmlParser()
// 	var urls []string
// 	query.Find("h3[class='repo-list-name'] a").Each(func(i int, s *goquery.Selection) {
// 		href, _ := s.Attr("href")
// 		urls = append(urls, "http://github.com/"+href)
// 	})
// 	// these urls will be saved and crawed by other coroutines.
// 	p.AddTargetRequests(urls, "html")

// 	name := query.Find(".entry-title .author").Text()
// 	name = strings.Trim(name, " \t\n")
// 	repository := query.Find(".entry-title .js-current-repository").Text()
// 	repository = strings.Trim(repository, " \t\n")
// 	//readme, _ := query.Find("#readme").Html()
// 	if name == "" {
// 		p.SetSkip(true)
// 	}
// 	// the entity we want to save by Pipeline
// 	p.AddField("author", name)
// 	p.AddField("project", repository)
// 	//p.AddField("readme", readme)
// }

func (this *MyPageProcesser) Process(p *page.Page) {
	if !p.IsSucc() {
		println(p.Errormsg())
		return
	}

	query := p.GetHtmlParser()
	name := query.Find("id='RYGLineStatus'").Text()
	log.Println("name", name)
	name2 := query.Find("span").Text()
	log.Println("name2", name2, name)
	name3 := query.Find(".mtrimd_line_col .mtrimd_line_idicon").Text()
	log.Println("name3", name3)

	//p.AddField("readme", readme)
}

func (this *MyPageProcesser) Finish() {
	fmt.Printf("TODO:before end spider \r\n")
}

func main() {
	// GetLine()
	spider.NewSpider(NewMyPageProcesser(), "TaskName").
		AddUrl("http://www.mtr.com.hk/ch/customer/main/index.html", "html"). // Start url, html is the responce type ("html" or "json" or "jsonp" or "text")
		AddPipeline(pipeline.NewPipelineConsole()).                          // Print result on screen
		SetThreadnum(3).                                                     // Crawl request by three Coroutines
		Run()
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

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
func CreateFile2(fileName string, data []byte) {
	if checkFileIsExist(fileName) {
		if f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666); err != nil {
			log.Println("open file error:", err)
		} else if _, err := f.Write(data); err != nil {
			log.Println("write file error:", err)
		}
	} else if err := ioutil.WriteFile(fileName, data, 0666); err != nil {
		log.Println("create file error:", err)
	}
}

func CreateFile(fileName string, data []byte) {
	if f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm); err != nil {
		log.Println("open file error:", err)
	} else if _, err := f.Write(data); err != nil {
		log.Println("write file error:", err)
	}
}
