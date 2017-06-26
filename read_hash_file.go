package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

//读取整个指纹文件到内存
func readFileContent(file *os.File) *map[string]string {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("readFileContent模块出错")
		}
	}()
	jsonContent, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("读取指纹文件内容错误： " + err.Error())
	}
	content := make(map[string]string)
	err := json.Unmarshal(jsonContent, &content)
	if err != nil {
		log.Println("指纹文件 json decode 失败： " + err.Error())
	}
	return &content
}

//比对指纹文件中的md5和新读取文件的md5
func comparedFileMd5(mapContent *map[string]string, filePath string, md5 string) *map[string]string {
	if contentMd5, ok := (*mapContent)[filePath]; ok {
		//如果md5存在，且不相同，则代表文件更新过，更新md5值
		if md5 != contentMd5 {
			(*mapContent)[filePath] = md5
		}
	} else {
		//如果md5不存在，则写入新的path
		(*mapContent)[filePath] = md5
	}
	return mapContent
}

func writeFileContent(mapContent *map[string]string) {
	jsonContent, err := json.Marshal(*mapContent)
}
