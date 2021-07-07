package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"log"
)

func Encrypt(text []byte, key []byte, iv []byte) (rs []byte, err error) {
	defer func() {
		if p := recover(); p != nil {
			str, ok := p.(string)
			if ok {
				rs = nil
				err = errors.New(str)
			} else {
				rs = nil
				err = errors.New("panic")
			}
		}
	}()
	//生成cipher.Block 数据块
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println("错误 -" + err.Error())
		return nil, err
	}
	//填充内容，如果不足16位字符
	blockSize := block.BlockSize()
	originData := pad(text, blockSize)
	//加密方式
	blockMode := cipher.NewCBCEncrypter(block, iv)
	//加密，输出到[]byte数组
	crypted := make([]byte, len(originData))
	blockMode.CryptBlocks(crypted, originData)
	return crypted, nil
}

func pad(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func Decrypt(data []byte, key []byte, iv []byte) (rs []byte, err error) {
	//生成密码数据块cipher.Block
	block, _ := aes.NewCipher(key)
	//解密模式
	blockMode := cipher.NewCBCDecrypter(block, iv)
	//输出到[]byte数组
	originData := make([]byte, len(data))
	blockMode.CryptBlocks(originData, data)
	//去除填充,并返回
	return unPad(originData), nil
}

func unPad(ciphertext []byte) []byte {
	length := len(ciphertext)
	//去掉最后一次的padding
	unPadding := int(ciphertext[length-1])
	return ciphertext[:(length - unPadding)]
}
