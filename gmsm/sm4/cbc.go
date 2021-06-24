package sm4

import (
	"bytes"
	"crypto/cipher"
	"github.com/tjfoc/gmsm/sm4"
)

func EncryptCBC(key []byte, in []byte) ([]byte, error) {
	return sm4.Sm4Cbc(key, in, true)
}

func DecryptCBC(key []byte, in []byte) ([]byte, error) {
	return sm4.Sm4Cbc(key, in, false)
}

func EncryptSm4(key, iv []byte, plainText []byte) ([]byte, error) {
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	paddData := pading(plainText, block.BlockSize())
	blokMode := cipher.NewCBCEncrypter(block, iv[:block.BlockSize()])
	cipherText := make([]byte, len(paddData))
	blokMode.CryptBlocks(cipherText, paddData)
	return cipherText, nil
}

func DecryptSm4(key, iv []byte, cipherText []byte) ([]byte, error) {
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv[:block.BlockSize()])
	blockMode.CryptBlocks(cipherText, cipherText)
	plainText := unPading(cipherText)
	return plainText, nil
}

func pading(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func unPading(ciphertext []byte) []byte {
	length := len(ciphertext)
	unPadding := int(ciphertext[length-1])
	return ciphertext[:(length - unPadding)]
}
