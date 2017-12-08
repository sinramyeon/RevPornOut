package main

import (
	"log"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(4)
	client := ConnTwitter()
	start := time.Now()

	a := []string{"국산야동", "몰카", "노모", "야동", "도촬", "일본야동", "공짜야동", "포르노", "업스커트", "발정", "자위", "화장실", "고딩", "중딩", "16녀", "17녀", "18녀", "19녀", "변녀", "암캐", "보지", "자영", "자위영상"}
	RevPornOut(client, a)

	elapsed := time.Since(start)
	log.Printf("Binomial took %s", elapsed)

}
