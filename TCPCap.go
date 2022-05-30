package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

func verificaErro(e error) {
	if e != nil {
		panic(e)
	}
}

func getFilePath() (string, string) {
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

	newFileName = filepath.Join(saveDir, newFileName) //저장 디렉토리 + 파일명

	fileTarget := ""
	if runtime.GOOS == "linux" {
		fileTarget = "pipe:1"
	} else {
		fileTarget = newFileName //저장 대상 파일
	}
	return fileTarget, newFileName
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

func saveImg(fileTarget string, cmdName string, args []string) {
	f, err4 := os.OpenFile(fileTarget, os.O_RDWR|os.O_APPEND, 0666)
	verificaErro(err4)

	cmd := exec.Command(cmdName, args...)
	stdout, err := cmd.StdoutPipe()
	verificaErro(err)
	err2 := cmd.Start()
	verificaErro(err2)

	chunk := make([]byte, 1024)
	for {
		nr, err5 := stdout.Read(chunk)

		if nr > 0 {
			validData := chunk[:nr]
			nw, err6 := f.Write(validData)
			fmt.Printf("Write %d bytes\n", nw)
			verificaErro(err6)
		}
		if err5 != nil {
			if err5 == io.EOF {
				break
			}
			fmt.Printf("Error = %v\n", err5)
			continue
		}
	}
	if err := cmd.Wait(); err != nil {
		fmt.Printf("Wait command error: %v\n", err)
	}
}

func capturePic() {
	fileTarget, newFileName := getFilePath()
	if makeFile(newFileName) {
		println("Make File Confirm!")
	}

	IP := "tcp://166.104.185.46:13072"

	cmdName := "ffmpeg"
	args := []string{
		"-y",
		"-i",
		IP,
		"-s",
		"640x480",
		fileTarget,
	}

	saveImg(fileTarget, cmdName, args)
}

func main() {
	ch := make(chan int, 1) //0 작업전 1 이미지 저장 완료

	go func() {
		for i := 0; i < 10; i++ {
			capturePic()
			ch <- 1
			time.Sleep(10 * time.Minute)
		}
	}()
}
