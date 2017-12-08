package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/xuri/excelize"
)

// Global variables
var (
	Urls    []string
	TweetID []string
)

const Path = "/Saved/"

// ConnTwitter ...
// Connect With Twitter.
// I use env.go 's keys
func ConnTwitter() *twitter.Client {

	// 1. Get my auth keys
	var con TwitterConfig
	con = conf(con)

	// 2. you can use oauth1 to set config
	config := oauth1.NewConfig(con.ConfKey, con.ConfSecret)
	token := oauth1.NewToken(con.TokenKey, con.TokenSecret)

	// 3. connect with twitter.
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	return client

}

// RevPornOut...
// Search Tweets and Get site's url, and save as excel.
func RevPornOut(client *twitter.Client, keyword []string) {

	// search for keywords
	for _, v := range keyword {

		// 1. Search for keywords
		// if you're using standard twitter api, you couldn't search for more than 7 days...
		search, _, _ := client.Search.Tweets(&twitter.SearchTweetParams{
			Query: v,
			Count: 100,
		})

		// 2. get specific tweets
		for _, v := range search.Statuses {

			url := len(v.Entities.Urls)
			media := len(v.Entities.Media)

			// get tweets that has url
			if url != 0 {

				// get rid of spam and useless tweets
				TweetID = append(TweetID, v.User.ScreenName)
				for _, v := range v.Entities.Urls {

					site := v.ExpandedURL

					if !(strings.Contains(site, "nico") || strings.Contains(site, "kakao") || strings.Contains(site, "ask") || strings.Contains(site, "image") || strings.Contains(site, "video") || strings.Contains(site, "photo") || strings.Contains(site, "status") || strings.Contains(site, "twitter") || strings.Contains(site, "news") || strings.Contains(site, "tumblr") || strings.Contains(site, "postype") || strings.Contains(site, "ilbe") || strings.Contains(site, "naver") || strings.Contains(site, "file") || strings.Contains(site, "wordpress") || strings.Contains(site, "youtu") || strings.Contains(site, "media") || strings.Contains(site, "file") || strings.Contains(site, "daum") || strings.Contains(site, "tistory") || strings.Contains(site, "instiz") || strings.Contains(site, "instagram")) {

						Urls = append(Urls, v.ExpandedURL)

					}
				}

			}
			if media != 0 {

				TweetID = append(TweetID, v.User.ScreenName)
				for _, v := range v.Entities.Media {
					site := v.ExpandedURL

					if !(strings.Contains(site, "nico") || strings.Contains(site, "kakao") || strings.Contains(site, "ask") || strings.Contains(site, "image") || strings.Contains(site, "video") || strings.Contains(site, "photo") || strings.Contains(site, "status") || strings.Contains(site, "twitter") || strings.Contains(site, "news") || strings.Contains(site, "tumblr") || strings.Contains(site, "postype") || strings.Contains(site, "ilbe") || strings.Contains(site, "naver") || strings.Contains(site, "file") || strings.Contains(site, "wordpress") || strings.Contains(site, "youtu") || strings.Contains(site, "media") || strings.Contains(site, "file") || strings.Contains(site, "daum") || strings.Contains(site, "tistory") || strings.Contains(site, "instiz") || strings.Contains(site, "instagram")) {
						Urls = append(Urls, v.ExpandedURL)

					}

				}

			}

		}

	}

	// make excel
	header := map[string]string{"A1": "주소"}
	values := make(map[string]string)

	// excelize https://github.com/360EntSecGroup-Skylar/excelize
	for k, v := range Urls {

		values["A"+strconv.Itoa((k+2))] = v
	}

	style := `{"font":{"bold":true,"italic":true,"family":"Berlin Sans FB Demi","size":20,"color":"#777777"}}`

	//Down as Excel
	ExcelDown("SiteList.xlsx", style, header, values)

	//Down as File
	CreateFile(Urls)

}

// CreateFile ...
// make txt file
func CreateFile(url []string) error {

	filePath := MakeFolder() + ".txt"

	file, error := os.Create(filePath) // Truncates if file already exists, be careful!
	if error != nil {
		log.Fatalf("failed creating file: %s", error)
		return error
	}
	defer file.Close() // Make sure to close the file when you're done

	for _, v := range url {
		file.WriteString(v + `
			`)

	}

	if error != nil {
		log.Fatalf("failed writing to file: %s", error)
		return error
	}

	return error
}

// ExcelDown ...
// make excel file
func ExcelDown(fileNm, styleStr string, header, values map[string]string) error {
	xlsx := excelize.NewFile()
	for k, v := range header {
		xlsx.SetCellValue("Sheet1", k, v)
	}
	for k, v := range values {
		xlsx.SetCellValue("Sheet1", k, v)
	}

	styleID, err := xlsx.NewStyle(styleStr)
	if err != nil {
		log.Panic("[ERROR] xlsx.NewStyle() : ", err)
		return err
	}
	xlsx.SetCellStyle("Sheet1", "A1", "I1", styleID)

	filePath := MakeFolder() + ".xlsx"
	err = xlsx.SaveAs(filePath)
	if err != nil {
		log.Panic("[ERROR] xlsx.SaveAs() : ", err)
		return err
	}

	return nil
}

// MakeFolder ...
// make folder for save files
func MakeFolder() string {

	t := time.Now()
	sYear := strconv.Itoa(t.Year())
	sMonth := strconv.Itoa(int(t.Month()))
	day := t.Day()
	hour := strconv.Itoa(t.Hour())

	filename := fmt.Sprintf("%s.log", t.Format("2006010215"))
	st := sYear + sMonth
	filePath := Path + st + "/" + strconv.Itoa(day) + "/" + hour + "/" + filename

	os.MkdirAll(Path, os.ModePerm)

	return filePath
}

/*
자동신고용
agouti + selenium으로 트위터 자동 신고기능 제작
package crawling

import (
	"encoding/xml"
	"log"
	"net/http"
	env "slackbot/envsetting"
	"strings"

	"github.com/parnurzeal/gorequest"
	"github.com/sclevine/agouti"
)

// 이벤트 게시글 json parsing	용
type Write struct {
	Day        string
	AuthorName string
	Text       string
}

type EntryData struct {
	Key   string `xml:"name,attr"`
	Value string `xml:"text"`
}

type ViewEntry struct {
	Key   string      `xml:"unid,attr"`
	Value []EntryData `xml:"entrydata"`
}
type ViewEntries struct {
	XMLName     xml.Name    `xml:viewentries`
	ViewEntries []ViewEntry `xml:"viewentry"`
}

// 이벤트 얻어오기
func GetEvent() map[string]string {

	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	// 인터파크 사내접속을 위한 token 생성
	token := MakeToken()

	// 사내 이벤트 게시판 xml로 들어가서 파징
	parsed := new(ViewEntries)
	_, body, _ := gorequest.New().Get(
		"http://ione.interpark.com/gw/app/bult/bbs00000.nsf/wviwnotice?ReadViewEntries&start=1&count=14&restricttocategory=03&page=1||_=1504081645868",
	).Type("xml").AddCookie(
		&http.Cookie{Name: "LtpaToken", Value: token},
	).End()

	_ = xml.Unmarshal([]byte(body), &parsed)

	// 결과 정리
	var event Write
	var eventlist []Write

	for _, v := range parsed.ViewEntries {
		var entrydata []EntryData
		entrydata = v.Value

		for key, val := range entrydata {

			if event.AuthorName != "" || event.Day != "" || event.Text != "" {
				eventlist = append(eventlist, event)
				event.AuthorName = ""
				event.Day = ""
				event.Text = ""
			}

			switch key {
			case 1:
				event.Day = val.Value
			case 2:
				event.AuthorName = val.Value
			case 3:
				event.Text = val.Value
			}

		}
	}

	// 그 중 최신 3개만 가져옴
	returnlist := make(map[string]string)
	var loop = 0

	for _, v := range eventlist {
		if loop < 3 {
			returnlist[v.Text] = v.Day + " " + v.AuthorName
			loop++
		}
	}

	return returnlist
}

//ltpa 토큰 만들기
func MakeToken() string {

	//agouti 이용. chromedriver, phantomjs가 %PATH%에 있거나
	//mac인경우에는 brew로 설치 필요!

	driver := agouti.ChromeDriver()
	if err := driver.Start(); err != nil {
		log.Fatalln(err)
	}

	defer recover()
	defer driver.Stop()

	page, err := driver.NewPage(agouti.Browser("phantomjs"))
	if err != nil {
		log.Fatalln(err)
	}

	// 접속 (진짜 크롬 창이 뜸)
	if err := page.Navigate("http://ione.interpark.com/"); err != nil {
		log.Fatalln(err)
	}

	// 로그인
	var interenv env.Interpark
	interenv = env.InterparkLogin(interenv)
	ID := interenv.ID
	PW := interenv.PW
	page.FindByID("Username").SendKeys(ID)
	page.FindByID("Password").SendKeys(PW)

	page.FindByClass("loginSubmit").Click()

	// 이벤트 주소 접속
	if err := page.Navigate("http://ione.interpark.com/gw/app/bult/bbs00000.nsf/wviwnotice?ReadViewEntries&start=1&count=14&restricttocategory=03&page=1||_=1504081645868"); err != nil {
		log.Fatalln(err)
	}

	// 쿠키 얻기
	cookie, err := page.GetCookies()
	if err != nil {
		log.Fatalln(err)
	}

	// 토큰 추출
	for _, v := range cookie {
		if strings.Contains(v.Name, "LtpaToken") {
			return v.Value
		}
	}

	return ""

}

*/

/*

agouti in chrome

$ brew install chromedriver
$ go get github.com/sclevine/agouti
$ go run main.go
// main.go
package main

import (
    "github.com/sclevine/agouti"
    "log"
)

func main() {
    driver := agouti.ChromeDriver()
    if err := driver.Start(); err != nil {
        log.Fatalf("Failed to start driver:%v", err)
    }
    defer driver.Stop()

    page, err := driver.NewPage(agouti.Browser("chrome"))
    if err != nil {
        log.Fatalf("Failed to open page:%v", err)
    }

    if err := page.Navigate("http://qiita.com/"); err != nil {
        log.Fatalf("Failed to navigate:%v", err)
    }
    page.Screenshot("/tmp/chrome_qiita.jpg")
}
*/
