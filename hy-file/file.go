package hy_file

import (
	"errors"
	"os"
)

// 判断文件是否存在
func GetFileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

// 获取文件名称
func GetFileName(path string) (string, error) {
	if len(path) <= 0 {
		return "", errors.New("file path empty")
	}
	if fi, err := os.Stat(path); err != nil {
		return "", err
	} else {
		return fi.Name(), nil
	}
}

// 获取文件大小
func GetFileLength(path string) (int64, error) {
	if len(path) <= 0 {
		return 0, errors.New("file path empty")
	}
	if fi, err := os.Stat(path); err != nil {
		return 0, err
	} else {
		return fi.Size(), nil
	}
}
