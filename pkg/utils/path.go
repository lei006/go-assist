package utils

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}

func PathMustExist(fileName string) error {

	_, err := os.Stat(fileName)
	if os.IsNotExist(err) == false {
		return nil
	}

	return os.MkdirAll(fileName, os.ModePerm) //创建多级目录
}

//使用 go run . 运行..
func Is_Go_Run_Mod() bool {
	return GetBinPath() != GetWorkPath()
}

// 运行在 vs中...
func Is_RunAtVs() bool {
	bin_path := GetBinPath()
	return strings.Contains(bin_path, "Temp/go-build")
}

func GetBinPath() string {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		return ""
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}

func GetWorkPath() string {
	str, _ := os.Getwd()

	return strings.Replace(str, "\\", "/", -1) //将\替换成/
}

func GetExePath() string {

	workPath := ""

	if Is_RunAtVs() {
		//运行在vs 中
		workPath = GetWorkPath() + "/"
	} else {
		workPath = GetBinPath() + "/"
	}

	return workPath
}
