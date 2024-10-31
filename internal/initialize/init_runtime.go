package initialize

import (
	"fmt"
	"os"
	"task/internal/global"
)

func InitRuntime() {
	// 检查本地是否存在runtime文件夹
	// 获取当前程序的路径
	path := GetWorkDir()
	rtPath := isExistsCreatePath(path, "task_runtime")
	global.RuntimePath = rtPath
	isExistsCreatePath(rtPath, "data")
	isExistsCreatePath(rtPath, "config")
	isExistsCreatePath(rtPath, "scripts")
	isExistsCreatePath(rtPath, "logs")
}

// GetWorkDir
// 获取程序运行目录
func GetWorkDir() string {
	pwd, _ := os.Getwd()
	return pwd
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		// 创建文件夹
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return false, err
		} else {
			return true, nil
		}
	}
	return false, err
}

func isExistsCreatePath(path, name string) string {
	p := fmt.Sprintf("%s/%s", path, name)
	exists, err := PathExists(p)
	if err != nil {
		panic(err.Error())
	}

	if !exists {
		panic("创建【config】文件夹失败")
	}

	return p
}
