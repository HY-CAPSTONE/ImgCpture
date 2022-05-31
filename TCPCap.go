package main

import (
	_ "github.com/go-sql-driver/mysql"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

func verificaErro(e error) {
	if e != nil {
		panic(e)
	}
}

type dbConn struct {
	driverName  string
	userName    string
	password    string
	ip          string
	port        string
	targetTable string
}

func (db dbConn) dbAddress() (string, string) {
	dataSourceName := db.userName + ":" + db.password + "@tcp(" + db.ip + ":" + db.port + ")/" + db.targetTable
	return db.driverName, dataSourceName
}

func getFilePath() (string, string, string) {
	nowDir, err := os.Getwd()
	verificaErro(err)
	err = os.Mkdir("imgStuck", 0750)
	verificaErro(err)
	saveDir := filepath.Join(nowDir, "imgStuck") //저장 디렉토리
	timesheet := time.Now()
	newFileName := string(timesheet.Year()) +
		string(int(timesheet.Month())) +
		string(timesheet.Day()) +
		string(timesheet.Hour()) +
		string(timesheet.Minute()) +
		string(timesheet.Second())

	FileNameWithDir := filepath.Join(saveDir, newFileName) //저장 디렉토리 + 파일명

	fileTarget := ""
	if runtime.GOOS == "linux" {
		fileTarget = "pipe:1"
	} else {
		fileTarget = newFileName //저장 대상 파일
	}
	return fileTarget, FileNameWithDir, newFileName
}

func makeFile(fileOutput string) bool {
	var _, err3 = os.Stat(fileOutput)

	if os.IsNotExist(err3) {
		var file, err = os.Create(fileOutput)
		verificaErro(err)
		defer file.Close()
		return false
	}
	return true
}

func main() {
	//ch := make(chan int, 1) //0 작업전 1 이미지 저장 완료
	dbInfo := dbConn{
		"mysql",
		"root",
		"wlgkcjf21gh",
		"112.170.208.72",
		"8080",
		"capstone",
	}

	go func() {
		for i := 0; i < 10; i++ {
			imgDir := capturePic()
			insert_to_db(dbInfo, imgDir)
			time.Sleep(10 * time.Minute)
		}
	}()

}
