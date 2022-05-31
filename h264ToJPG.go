package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

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

func capturePic() string {
	fileTarget, FileNameWithDir, FileName := getFilePath()
	if makeFile(FileNameWithDir) {
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
	return "/imgStuck/" + FileName
}
