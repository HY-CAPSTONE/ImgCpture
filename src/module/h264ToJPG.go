package module

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func VerificaErro(e error) {
	if e != nil {
		panic(e)
	}
}

func GetFilePath() (string, string, string) {
	nowDir, err := os.Getwd()
	VerificaErro(err)
	if _, err3 := os.Stat("imgStuck"); os.IsNotExist(err3) {
		err = os.Mkdir("imgStuck", 0750)
		VerificaErro(err)
	}
	saveDir := filepath.Join(nowDir, "imgStuck") //저장 디렉토리
	timesheet := time.Now()
	dateName := []string{
		strconv.Itoa(timesheet.Year()),
		strconv.Itoa(int(timesheet.Month())),
		strconv.Itoa(timesheet.Day()),
		strconv.Itoa(timesheet.Hour()),
		strconv.Itoa(timesheet.Minute()),
		strconv.Itoa(timesheet.Second()),
	}
	newFileName := strings.Join(dateName, "-")
	newFileName = newFileName + ".jpg"
	fmt.Println(newFileName)

	FileNameWithDir := filepath.Join(saveDir, newFileName) //저장 디렉토리 + 파일명

	fileTarget := ""
	if runtime.GOOS == "linux" {
		fileTarget = "pipe:1"
	} else {
		fileTarget = newFileName //저장 대상 파일
	}

	return fileTarget, FileNameWithDir, newFileName
}

func MakeFile(fileOutput string) bool {
	var _, err3 = os.Stat(fileOutput)

	if os.IsNotExist(err3) {
		var file, err = os.Create(fileOutput)
		VerificaErro(err)
		defer file.Close()
		return false
	}
	return true
}

func SaveImg(fileTarget string, cmdName string, args []string) {
	f, err4 := os.OpenFile(fileTarget, os.O_RDWR|os.O_APPEND, 0666)
	VerificaErro(err4)

	cmd := exec.Command(cmdName, args...)
	stdout, err := cmd.StdoutPipe()
	VerificaErro(err)
	err2 := cmd.Start()
	VerificaErro(err2)

	chunk := make([]byte, 1024)
	for {
		nr, err5 := stdout.Read(chunk)

		if nr > 0 {
			validData := chunk[:nr]
			nw, err6 := f.Write(validData)
			fmt.Printf("Write %d bytes\n", nw)
			VerificaErro(err6)
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

func CapturePic(getImg chan bool, fileName *string) {
	fileTarget, FileNameWithDir, FileName := GetFilePath()
	if MakeFile(FileNameWithDir) {
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
		FileNameWithDir,
	}
	println(fileTarget)
	SaveImg(FileNameWithDir, cmdName, args)
	*fileName = "/imgStuck/" + FileName
	println("Downloading...Done")
	getImg <- true
}
