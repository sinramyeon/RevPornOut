package main

import (
	"crypto/rand"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	//_ "github.com/sclevine/agouti"
	"github.com/xuri/excelize"
)

// 트위터 데이터 파징용 스트럭트
type Twitter struct {
	Day        string
	AuthorName string
	Text       string
}

// Global variables
var (
	Urls    []string
	TweetID []string
	TweetUser []string
)

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

	//var Urls []string
	//var TweetID []string

	// 키워드 별로 검색 실행 - > go 루틴으로 나눠 병렬처리
	for _, v := range keyword {

		// 1. 키워드를 검색
		// 스탠다드에서는 7일 이내 것만 검색 가능
		search, _, _ := client.Search.Tweets(&twitter.SearchTweetParams{
			Query: v,
			Count: 100,
		})

		// 2. 검색한 키워드 내에서 나눔
		for _, v := range search.Statuses {

			url := len(v.Entities.Urls)
			media := len(v.Entities.Media)

			// 주소를 갖고 있는 트윗만 꺼내옴
			if url != 0 {

				// 내부에서 사진 등의 링크는 거르고, 실제 링크만 저장
				TweetID = append(TweetID, v.User.ScreenName)
				for _, v := range v.Entities.Urls {

					site := v.ExpandedURL

					if !(strings.Contains(site, "nico") || strings.Contains(site, "kakao") || strings.Contains(site, "ask") || strings.Contains(site, "image") || strings.Contains(site, "video") || strings.Contains(site, "photo") || strings.Contains(site, "status") || strings.Contains(site, "twitter") || strings.Contains(site, "news") || strings.Contains(site, "tumblr") || strings.Contains(site, "postype") || strings.Contains(site, "ilbe") || strings.Contains(site, "naver") || strings.Contains(site, "file") || strings.Contains(site, "wordpress") || strings.Contains(site, "youtu") || strings.Contains(site, "media") || strings.Contains(site, "file") || strings.Contains(site, "daum") || strings.Contains(site, "tistory") || strings.Contains(site, "instiz") || strings.Contains(site, "instagram")) {

						Urls = append(Urls, v.ExpandedURL)

					}
				}

			}
			if media != 0 {

				// 내부에서 사진 등의 링크는 거르고, 실제 링크만 저장
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

	Urls = removeDuplicatesUnordered(Urls)
}

// MakeExcel ...
// Make Excel File
func MakeExcel() {

	// 엑셀 저장
	header := map[string]string{"A1": "주소"}
	values := make(map[string]string)

	// 해당 엑셀 라이브러리 https://github.com/360EntSecGroup-Skylar/excelize
	for k, v := range Urls {

		values["A"+strconv.Itoa((k+2))] = v
	}

	style := `{"font":{"bold":true,"italic":true,"family":"Berlin Sans FB Demi","size":20,"color":"#777777"}}`

	if len(Urls) != 0 {

		ExcelDown("SiteList.xlsx", style, header, values)

		// 텍스트파일 저장

		CreateFile(Urls)
	}

}

// CreateFile ...
// Create Excel File
func CreateFile(url []string) error {

	uuid := CreateUUID()

	file, error := os.Create(uuid + ".txt") // Truncates if file already exists, be careful!
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
// Download Excel File
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

	uuid := CreateUUID()
	filepath := "/public/temp/" + uuid + ".xlsx"

	os.MkdirAll("/public/temp/", os.ModePerm)

	err = xlsx.SaveAs(filepath)
	if err != nil {
		log.Panic("[ERROR] xlsx.SaveAs() : ", err)
		return err
	}

	return nil
}

// CreateUUID ...
// create a random UUID with from RFC 4122
// adapted from http://github.com/nu7hatch/gouuid
func CreateUUID() (uuid string) {
	u := new([16]byte)
	_, err := rand.Read(u[:])
	if err != nil {
		log.Fatalln("Cannot generate UUID", err)
	}

	// 0x40 is reserved variant from RFC 4122
	u[8] = (u[8] | 0x40) & 0x7F
	// Set the four most significant bits (bits 12 through 15) of the
	// time_hi_and_version field to the 4-bit version number.
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	uuid = fmt.Sprintf("%x%x%x%x%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return
}

// keepDoingSomething will keep trying to doSomething() until either
// we get a result from doSomething() or the timeout expires
func removeDuplicatesUnordered(elements []string) []string {
	encountered := map[string]bool{}

	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v]] = true
	}

	// Place all keys from the map into a slice.
	result := []string{}
	for key, _ := range encountered {
		result = append(result, key)
	}
	return result
}


// RevPornUserOut ...
// Search Tweets and Get user ID
func RevPornUserOut(client *twitter.Client, keyword []string) {

	// 키워드 별로 검색 실행 - > go 루틴으로 나눠 병렬처리
	for _, v := range keyword {

		// 1. 키워드를 검색
		// 스탠다드에서는 7일 이내 것만 검색 가능
		search, _, _ := client.Search.Tweets(&twitter.SearchTweetParams{
			Query: v,
			Count: 100,
		})

		// 2. 검색한 키워드 내에서 나눔
		for _, v := range search.Statuses {

			id := v.IDStr

			// id가 없지 않을 때
			if id != nil {

				TweetUser = append(TweetUser, id)

			}

	TweetUser = removeDuplicatesUnordered(TweetUser)

}

// TweetBlockUser ...
// make user block tweet string
func TweetBlockUser(client *twitter.Client){

	string := `<신고를 부르는 계정 타래>

	신고에 동참해주시는 분들과 제보해주시는 분들, 모두 고맙습니다😊
	
	❌ 미디어주의 ❌
	
	특히 아래의 아이디를 집중적으로 신고해주세요.
	다른 아이디도 po신고wer 부탁드립니다`
	blank := ""
	atMark := "@"
	blockUser = ""

	for _, i := range TweetUser {

		// "make tweets like @id @id "

		blockUser+atMark+blank
		 

	}

}


// 트윗 쓰기
func SendTweet(client *twitter.Client, str string) {

	client.Statuses.Update(str, nil)


}



/*

	
	/*
		{
		"tweet": {
		"created_at": "Thu Apr 06 15:24:15 +0000 2017",
		"id_str": "850006245121695744",
		"text": "1\/ Today we\u2019re sharing our vision for the future of the Twitter API platform!\nhttps:\/\/t.co\/XweGngmxlP",
		"user": {
		"id": 2244994945,
		"name": "Twitter Dev",
		"screen_name": "TwitterDev",
		"location": "Internet",
		"url": "https:\/\/dev.twitter.com\/",
		"description": "Your official source for Twitter Platform news, updates & events. Need technical help? Visit https:\/\/twittercommunity.com\/ \u2328\ufe0f #TapIntoTwitter"
		},
		"place": {
		
		},
		"entities": {
		"hashtags": [
			
		],
		"urls": [
			{
			"url": "https:\/\/t.co\/XweGngmxlP",
			"unwound": {
				"url": "https:\/\/cards.twitter.com\/cards\/18ce53wgo4h\/3xo1c",
				"title": "Building the Future of the Twitter API Platform"
			}
			}
		],
		"user_mentions": [
			
		]
		}
	}
	}


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
