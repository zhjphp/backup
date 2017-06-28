package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

//指纹集合文件存放目录 【暂时不需要，文件存储在与执行文件相同的目录】
var hashFilePath string

//指纹集合文件名称
var hashFileName string = "hash.zlbf"

//是否单独设置首次备份执行时间
var isDoFrist string

//首次备份运行的时间
var doFristTimeHour int = -1
var doFristTimeMin int = -1
var doFristTimeSec int = -1

//正常备份运行时间
var doTimeHour int = -1
var doTimeMin int = -1
var doTimeSec int = -1

var md5FilePath string

/*
备份文件目录存放 map
k:备份目录
v:存放目录
*/
var backupFilePathMap map[string]string

func main() {

	//读取输入
	reader := bufio.NewReader(os.Stdin)

	//暂时停用设置指纹文件目录，默认存储在当前程序执行目录下
	//setHashFilePath(reader)

	//设置备份时间
	setBackupTime(reader)

	//设置备份文件目录
	setBackupFilePathMap(reader)

	//创建指纹文件
	createHashFile()

	var hashMapContent *map[string]string

	for {
		now := time.Now()
		next := now.Add(time.Hour * 24)
		if isDoFrist == "y" || isDoFrist == "Y" {
			next = time.Date(now.Year(), now.Month(), now.Day(), doFristTimeHour, doFristTimeMin, doFristTimeSec, 0, now.Location())
		} else {
			next = time.Date(next.Year(), next.Month(), next.Day(), doTimeHour, doTimeMin, doTimeSec, 0, next.Location())
		}
		duration := next.Sub(now)
		t := time.NewTicker(duration)

		select {
		case <-t.C:
			//到时间后执行备份任务
			var targetPath string
			var partPath string

			//读取指纹文件到内存
			hashMapContent = readFileContent()

			//path:原始目录，targetPath：目标目录
			for key, value := range backupFilePathMap {
				targetPath = ""
				filepath.Walk(key, func(path string, f os.FileInfo, err error) error {
					partPath, _ = filepath.Rel(key, path)
					targetPath = filepath.Join(value, partPath)

					//path:原始文件地址，targetPath:备份文件地址
					//每个path都需要去比对md5文件中做比对，判断文件是否被修改过
					//如果文件是个目录则不写入指纹文件
					if f.IsDir() {
						copyFile(path, targetPath)
					} else {
						md5 := makeFileMd5(path) //获取文件md5
						isUpdate := comparedFileMd5(hashMapContent, md5, path)
						//如果修改过则复制文件，并更新md5文件
						if isUpdate {
							copyFile(path, targetPath)
						}
						//如果没有修改过则不执行任何操作
					}
					return nil
				})
			}
			//最终返回更新过的md5_jsong，写入json文件
			writeFileContent(hashMapContent, hashFilePath)
			//释放读取的指纹文件内存
			hashMapContent = nil
			isDoFrist = "n"
		}
	}
}

func checkDir(path string) {

}

//写入指纹文件
func writeFileContent(mapContent *map[string]string, path string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("writeFileContent模块出错")
		}
	}()

	jsonContent, err := json.Marshal(*mapContent)

	if err != nil {
		log.Println("指纹文件 json Marshal 失败： " + err.Error())
	}
	err = ioutil.WriteFile(path, jsonContent, 0777)
	if err != nil {
		log.Println("写入指纹文件失败： " + err.Error())
	}
}

func copyFile(basePath, targetPath string) {
	defer func() {
		if err_p := recover(); err_p != nil {
			fmt.Println("copyFile模块出错")
		}
	}()

	baseStat, err := os.Stat(basePath)
	if err != nil {
		log.Panicln("需要备份的文件检测失败，文件出现问题，无法复制")
		return
	}
	//targetStat, err := os.Stat(targetPath)
	_, err = os.Stat(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			//如果目标文件不存在
			if baseStat.IsDir() {
				//如果缺失的是一个空目录
				errMkDir := os.MkdirAll(targetPath, 0777)
				if errMkDir != nil {
					log.Println("创建目录 " + targetPath + " 失败")
				}
			} else {
				//如果缺失的是一个文件,则复制文件
				copyFileContent(basePath, targetPath)
			}
		} else {
			return
		}
	} else {
		//如果目标文件存在
		if baseStat.IsDir() {
			//如果是一个空目录
		} else {
			//如果是一个文件，则复制文件
			copyFileContent(basePath, targetPath)
		}
	}

}

//复制文件内容
func copyFileContent(basePath, targetPath string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("copyFileContent模块出错")
		}
	}()

	baseFile, err := os.Open(basePath)
	if err != nil {
		log.Println("读取文件 " + basePath + " 失败")
		return
	}
	defer func() {
		err := baseFile.Close()
		if err != nil {
			log.Println("close文件 " + basePath + " 失败： " + err.Error())
		}
	}()
	targetFile, err := os.Create(targetPath)
	if err != nil {
		log.Println("创建文件 " + targetPath + " 失败： " + err.Error())
		return
	}
	defer func() {
		err := targetFile.Close()
		if err != nil {
			log.Println("close文件 " + targetPath + " 失败： " + err.Error())
		}
	}()
	copyData, err := io.Copy(targetFile, baseFile)
	if err != nil {
		log.Println("复制文件文件 " + basePath + " 失败： " + err.Error())
	}
	fmt.Println("正在复制文件： " + basePath + " 大小为： " + strconv.FormatInt(copyData, 10))
}

//比对指纹文件中的md5和新读取文件的md5
func comparedFileMd5(mapContent *map[string]string, md5 string, path string) bool {
	defer func() {
		if err := recover(); err != nil {
			log.Println("comparedFileMd5模块出错")
		}
	}()

	if contentMd5, ok := (*mapContent)[path]; ok {
		//如果md5存在，且不相同，则代表文件更新过，更新md5值，并且复制文件
		if md5 != contentMd5 {
			(*mapContent)[path] = md5
			return true
		} else {
			return false
		}
	} else {
		//如果md5不存在，则写入新的path
		(*mapContent)[path] = md5
		return true
	}
}

//生成文件md5
func makeFileMd5(filePath string) string {
	defer func() {
		if err_p := recover(); err_p != nil {
			fmt.Println("makeFileMd5模块出错")
		}
	}()

	hash := md5.New()
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("打开文件 " + filePath + " 准备验证MD5失败： " + err.Error())
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Println("close文件 " + filePath + " 失败： " + err.Error())
		}
	}()

	_, err = io.Copy(hash, file)
	if err != nil {
		log.Println("准备读取" + filePath + "文件内容验证md5失败： " + err.Error())
	}
	md5 := hex.EncodeToString(hash.Sum(nil))
	return md5
}

//读取整个指纹文件到内存
func readFileContent() *map[string]string {
	defer func() {
		if err := recover(); err != nil {
			log.Println("readFileContent模块出错")
		}
	}()

	file, err := os.Open(hashFilePath)
	defer func() {
		err := file.Close()
		if err != nil {
			log.Println("close指纹文件 " + hashFilePath + " 失败： " + err.Error())
		}
	}()

	jsonContent, err := ioutil.ReadAll(file)

	if err != nil {
		log.Println("读取指纹文件内容错误： " + err.Error())
	}
	content := make(map[string]string)
	err = json.Unmarshal(jsonContent, &content)
	if err != nil {
		log.Println("指纹文件 json Unmarshal 失败： " + err.Error())
	}
	return &content
}

//创建指纹文件
func createHashFile() {
	defer func() {
		if err_p := recover(); err_p != nil {
			fmt.Println("createHashFile模块出错")
		}
	}()

	hashFilePath = filepath.Join(getCurrentPath(), hashFileName)
	err := ioutil.WriteFile(hashFilePath, []byte("{}"), 0777)
	if err != nil {
		log.Println("创建指纹文件失败： " + err.Error())
	}
}

//获取文件当前执行路径
func getCurrentPath() string {
	defer func() {
		if err_p := recover(); err_p != nil {
			fmt.Println("getCurrentPath模块出错")
		}
	}()

	commandFile, err := exec.LookPath(os.Args[0])
	if err != nil {
		log.Println("获取当前执行文件命令行失败： " + err.Error())
	}
	absFile, err := filepath.Abs(commandFile)
	if err != nil {
		log.Println("获取当前执行文件abs路径失败： " + err.Error())
	}
	absPath := filepath.Dir(absFile)
	return absPath
}

//设置备份文件
func setBackupFilePathMap(reader *bufio.Reader) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("setBackupFilePathMap模块出错")
		}
	}()

	backupFilePathMap = make(map[string]string)
	var inputKey []byte
	var inputValue []byte
	var key string
	var value string
	var errInput error
	var count int64 = 0
	log.Println(`输入exit为结束当前行，并开始备份`)

	for {
		inputKey = []byte{}
		inputValue = []byte{}
		key = ""
		value = ""
		count++
		log.Println("请输第【" + strconv.FormatInt(count, 10) + "】组数据：")

		for len(inputKey) == 0 {
			log.Println("请输入要备份的目录,回车结束：")
			inputKey, _, errInput = reader.ReadLine()
			if len(inputKey) == 0 {
				log.Println("您未输入任何内容")
				inputKey = []byte{}
			}
			if errInput != nil {
				log.Println("读取输入错误：" + errInput.Error())
				inputKey = []byte{}
			}
		}
		key = string(inputKey)
		if (key == "exit") || (key == "EXIT") {
			break
		}

		for len(inputValue) == 0 {
			log.Println("请输入备份文件存放目录,回车结束：")
			inputValue, _, errInput = reader.ReadLine()
			if len(inputValue) == 0 {
				log.Println("您未输入任何内容")
				inputValue = []byte{}
			}
			if errInput != nil {
				log.Println("读取输入错误：" + errInput.Error())
				inputValue = []byte{}
			}
		}
		value = string(inputValue)
		if (string(inputValue) == "exit") || (string(inputValue) == "EXIT") {
			break
		}

		backupFilePathMap[key] = value
	}
}

//设置首次备份时间
func setFristBackupTime(reader *bufio.Reader) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("setFristBackupTime模块出错")
		}
	}()

	var doFristTimeInput []byte
	var errDoFristTimeInput error
	var errToInt error

	//如果设置需要设置首次备份时间
	if (isDoFrist == "y") || (isDoFrist == "Y") {
		for (doFristTimeHour) < 0 || (doFristTimeHour > 24) {
			log.Println("请输入首次备份开始的小时数(只支持0-24的整数，对应一天中的0至24点)：")
			doFristTimeInput, _, errDoFristTimeInput = reader.ReadLine()
			if errDoFristTimeInput != nil {
				log.Println("读取输入错误：" + errDoFristTimeInput.Error())
			} else if len(doFristTimeInput) == 0 {
				log.Println("您未输入任何内容")
			}
			doFristTimeHour, errToInt = strconv.Atoi(string(doFristTimeInput))
			if errToInt != nil {
				log.Println("小时数输入错误")
				doFristTimeHour = -1
			} else if (doFristTimeHour) < 0 || (doFristTimeHour > 24) {
				log.Println("请输入0-24之间的整数")
			}
		}

		for (doFristTimeMin) < 0 || (doFristTimeMin > 60) {
			log.Println("请输入首次备份开始的分钟数(只支持0-60的整数")
			doFristTimeInput, _, errDoFristTimeInput = reader.ReadLine()
			if errDoFristTimeInput != nil {
				log.Println("读取输入错误：" + errDoFristTimeInput.Error())
			} else if len(doFristTimeInput) == 0 {
				log.Println("您未输入任何内容")
			}
			doFristTimeMin, errToInt = strconv.Atoi(string(doFristTimeInput))
			if errToInt != nil {
				log.Println("分钟数输入错误")
				doFristTimeMin = -1
			} else if (doFristTimeMin) < 0 || (doFristTimeMin > 60) {
				log.Println("请输入0-60之间的整数")
			}
		}

		for (doFristTimeSec) < 0 || (doFristTimeSec > 60) {
			log.Println("请输入首次备份开始的秒数(只支持0-60的整数")
			doFristTimeInput, _, errDoFristTimeInput = reader.ReadLine()
			if errDoFristTimeInput != nil {
				log.Println("读取输入错误：" + errDoFristTimeInput.Error())
			} else if len(doFristTimeInput) == 0 {
				log.Println("您未输入任何内容")
			}
			doFristTimeSec, errToInt = strconv.Atoi(string(doFristTimeInput))
			if errToInt != nil {
				log.Println("秒数输入错误")
				doFristTimeSec = -1
			} else if (doFristTimeSec) < 0 || (doFristTimeSec > 60) {
				log.Println("请输入0-60之间的整数")
			}
		}
	}
}

//设置日常执行备份时间
func setNormalBackupTime(reader *bufio.Reader) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("setBackupTime模块出错")
		}
	}()

	var doTimeInput []byte
	var errDoTimeInput error
	var errToInt error

	for (doTimeHour) < 0 || (doTimeHour > 24) {
		log.Println("请输入日常备份开始的小时数(只支持0-24的整数，对应一天中的0至24点)：")
		doTimeInput, _, errDoTimeInput = reader.ReadLine()
		if errDoTimeInput != nil {
			log.Println("读取输入错误：" + errDoTimeInput.Error())
		} else if len(doTimeInput) == 0 {
			log.Println("您未输入任何内容")
		}
		doTimeHour, errToInt = strconv.Atoi(string(doTimeInput))
		if errToInt != nil {
			log.Println("小时数输入错误")
			doTimeHour = -1
		} else if (doTimeHour) < 0 || (doTimeHour > 24) {
			log.Println("请输入0-24之间的整数")
		}
	}

	for (doTimeMin) < 0 || (doTimeMin > 60) {
		log.Println("请输入日常备份开始的分钟数(只支持0-60的整数")
		doTimeInput, _, errDoTimeInput = reader.ReadLine()
		if errDoTimeInput != nil {
			log.Println("读取输入错误：" + errDoTimeInput.Error())
		} else if len(doTimeInput) == 0 {
			log.Println("您未输入任何内容")
		}
		doTimeMin, errToInt = strconv.Atoi(string(doTimeInput))
		if errToInt != nil {
			log.Println("分钟数输入错误")
			doTimeMin = -1
		} else if (doTimeMin) < 0 || (doTimeMin > 60) {
			log.Println("请输入0-60之间的整数")
		}
	}

	for (doTimeSec) < 0 || (doTimeSec > 60) {
		log.Println("请输入日常备份开始的秒数(只支持0-60的整数：")
		doTimeInput, _, errDoTimeInput = reader.ReadLine()
		if errDoTimeInput != nil {
			log.Println("读取输入错误：" + errDoTimeInput.Error())
		} else if len(doTimeInput) == 0 {
			log.Println("您未输入任何内容")
		}
		doTimeSec, errToInt = strconv.Atoi(string(doTimeInput))
		if errToInt != nil {
			log.Println("秒数输入错误")
			doTimeSec = -1
		} else if (doTimeSec) < 0 || (doTimeSec > 60) {
			log.Println("请输入0-60之间的整数")
		}
	}
}

//设置备份启动时间
func setBackupTime(reader *bufio.Reader) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("setBackupTime模块出错")
		}
	}()

	var isDoFristInput []byte
	var errIsDoFristInput error

	for (isDoFrist != "y") && (isDoFrist != "Y") && (isDoFrist != "n") && (isDoFrist != "N") {
		log.Println("是否需要单独设置首次备份时间？(y/n)")
		isDoFristInput, _, errIsDoFristInput = reader.ReadLine()
		if errIsDoFristInput != nil {
			log.Println("读取输入错误：" + errIsDoFristInput.Error())
		} else {
			isDoFrist = string(isDoFristInput)
			if len(isDoFrist) == 0 {
				log.Println("您未输入任何内容")
			} else if (isDoFrist != "y") && (isDoFrist != "Y") && (isDoFrist != "n") && (isDoFrist != "N") {
				log.Println("请输入y/n")
			}
		}
	}

	setFristBackupTime(reader)
	setNormalBackupTime(reader)

}

//设置指纹集合文件存放目录
func setHashFilePath(reader *bufio.Reader) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("setHashFilePath模块出错")
		}
	}()

	var hashFilePathInput []byte
	var errHashFilePathInput error

	for len(hashFilePathInput) == 0 {
		log.Println("请输入文件指纹存放目录，回车结束输入：")
		hashFilePathInput, _, errHashFilePathInput = reader.ReadLine()
		if errHashFilePathInput != nil {
			log.Println("读取输入错误：" + errHashFilePathInput.Error())
		} else {
			hashFilePath = string(hashFilePathInput)
			if len(hashFilePathInput) == 0 {
				log.Println("您未输入任何内容")
			}
		}
	}
}
