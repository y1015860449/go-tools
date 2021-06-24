package md5

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func StringMd5(txt string) (string, error) {
	m := md5.New()
	_, err := io.WriteString(m, txt)
	if err != nil {
		return "", err
	}
	arr := m.Sum(nil)
	return fmt.Sprintf("%x", arr), nil
}

func FileMd5(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	hash := md5.New()
	if _, err := io.Copy(hash, f); err != nil {
		fmt.Println("Copy", err)
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func Md5Sum32(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}