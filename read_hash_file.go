package main

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func writeFile(file *os.File, filePath string, md5 string) {
	defer func() {
		if err_p := recover(); err_p != nil {
			fmt.Println("writeFile模块出错")
		}
	}()
	content := filePath + "," + md5 + "\n"
	size, err := io.WriteString(file, content)
	if err != nil {
		log.Println("写入文件指纹失败： " + err.Error())
	}
}

func readLine(file *os.File, filePath string, md5 string) bool {
	defer func() {
		if err_p := recover(); err_p != nil {
			fmt.Println("readLine模块出错")
		}
	}()

	lineReader := bufio.NewReader(file)
	for {
		line, err := lineReader.ReadString("\n")
		dataSlice := strings.Split(line, ",")
		if err == io.EOF {
			//如果不存在filePath，则写入不存在的文件

		}
		if dataSlice[0] == filePath {
			if dataSlice[1] == md5 {
				return true
			}
		}
	}
}
