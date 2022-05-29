package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
)

func verificaErro(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	dir := "C:\\Users\\js_77\\Desktop\\"
	fileOutput := dir + "output.jpg"
	osEnv := string(runtime.GOOS)
	fileTarget := ""
	if osEnv == "linux" {
		fileTarget = "pipe:1"
	} else {
		fileTarget = fileOutput
	}
	cmdName := "ffmpeg"
	args := []string{
		"-y",
		"-i",
		"tcp://166.104.185.46:13072",
		"-s",
		"640x480",
		fileTarget,
	}
	println(fileTarget)

	var _, err3 = os.Stat(fileOutput)

	if os.IsNotExist(err3) {
		var file, err = os.Create(fileOutput)
		verificaErro(err)
		defer file.Close()
	}

	f, err4 := os.OpenFile(fileTarget, os.O_RDWR|os.O_APPEND, 0666)

	verificaErro(err4)

	println("Command start")
	cmd := exec.Command(cmdName, args...)
	stdout, err := cmd.StdoutPipe()
	verificaErro(err)
	err2 := cmd.Start()
	verificaErro(err2)

	chunk := make([]byte, 1024)
	for {
		println("Command end")
		nr, err5 := stdout.Read(chunk)
		fmt.Printf("Read %d bytes\n", nr)

		//do something with the data
		//e.g. write to file
		if nr > 0 {
			validData := chunk[:nr]
			nw, err6 := f.Write(validData)
			fmt.Printf("Write %d bytes\n", nw)
			verificaErro(err6)
		}

		if err5 != nil {
			//Reach end of file (stream), exit from loop
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
