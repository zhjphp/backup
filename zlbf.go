package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

//指纹集合文件存放目录
var hashFilePath string

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

/*
备份文件目录存放 map
k:备份目录
v:存放目录
*/
var backupFilePathMap map[string]string

func main() {
	reader := bufio.NewReader(os.Stdin)
	//setHashFilePath(reader)
	//setBackupTime(reader)
}

func setBackupFilePathMap(reader *bufio.Reader) {
	backupFilePathMap = make(map[string]string)
	var inputKey []byte
	var inputValue []byte

	log.Println(`输入exit为结束当前命令行，并开始备份，windows目录中的“/”，请用“//”表示,如果备份目录的话目录的最后需要加上目录间隔符`)
}

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
		fmt.Println("首次：")
		fmt.Println(doFristTimeHour)
		fmt.Println(doFristTimeMin)
		fmt.Println(doFristTimeSec)
	}
}

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
	fmt.Println("日常：")
	fmt.Println(doTimeHour)
	fmt.Println(doTimeMin)
	fmt.Println(doTimeSec)

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
