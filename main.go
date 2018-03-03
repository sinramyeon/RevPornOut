package main

import (
	"fmt"
	"log"
	"runtime"
	"time"
)

func main() {

	log.Println("실행")
	runtime.GOMAXPROCS(4)
	client := ConnTwitter()

	// 리벤지 포르노 사이트 신고용
	tweetSearchWords := []string{"국산야동", "몰카", "노모", "야동", "도촬", "일본야동", "공짜야동", "포르노", "업스커트", "발정", "자위", "화장실", "고딩", "중딩", "16녀", "17녀", "18녀", "19녀", "변녀", "암캐", "보지", "자영", "자위영상"}

	// 합성 신고할 유저 찾기용
	userSearchWords := []string{""}

	/*
		RevPornOut(client, a)
		MakeExel()
	*/

	ticker := time.NewTicker(time.Minute * 15)
	ticker2 := time.NewTicker(time.Hour * 5)

	go func() {
		for t := range ticker.C {
			fmt.Println("트윗 서치", t)
			RevPornOut(client, tweetSearchWords)
			RevPornUserOut(client, userSearchWords)
		}

	}()

	go func() {
		for t := range ticker2.C {
			fmt.Println("트윗 저장", t)
			MakeExel()
		}

	}()

	time.Sleep(time.Hour * 24)
	ticker.Stop()
	fmt.Println("Ticker stopped")
	fmt.Println("Done")

}
