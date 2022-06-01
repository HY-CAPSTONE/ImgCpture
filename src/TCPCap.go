package main

import (
	"PictureCap/src/module"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var nowDBInfo module.DBConn

func makeTerm(restart *chan struct{}) {
	println("Timer Start...")
	time.Sleep(15 * time.Minute)
	println("Restart...")
	*restart <- struct{}{}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	var targetDir string = module.FindRecentImg("99", nowDBInfo)[1:]
	buf, err := ioutil.ReadFile(targetDir)
	println(targetDir)

	if err != nil {

		log.Fatal(err)
	}
	log.Println(targetDir)
	if len(buf) < 100 {
		w.Write([]byte("This Photo is crash \n유감입니다.\n-윤지상-\n" + targetDir))
	} else {
		w.Header().Set("Content-Type", "image/jpeg")
		w.Write(buf)
	}
}

func srcMessage(DBInfo *module.DBConn) *module.DBConn {
	return DBInfo
}

func main() {
	//quit := make(chan struct{})
	getImg := make(chan bool) //0 작업전 1 이미지 저장 완료
	sendImg := make(chan bool)
	restart := make(chan struct{})
	fileName := ""
	dbInfo := module.DBConn{
		"mysql",
		"root",
		"wlgkcjf21gh",
		"112.170.208.72",
		"8080",
		"capstone",
	}
	nowDBInfo = dbInfo

	go func() {
		for true {
			println("Downloading...")
			go module.CapturePic(getImg, &fileName)
			<-getImg //이미지가 다운로드 받아질 때까지 대기
			println("SendStarting...")
			go module.Insert_to_db(&sendImg, dbInfo, fileName)
			<-sendImg
			go makeTerm(&restart)
			<-restart
		}
	}()
	//<-quit
	handler := http.HandlerFunc(handleRequest)

	http.Handle("/", handler)

	fmt.Println("Server started at port 8080")
	http.ListenAndServe(":8080", nil)
}
