package sm4

import (
	"log"
	"testing"
)

func TestSm4(t *testing.T) {
	key := "11235fsdf54654f0"
	iv := "85rdsd=4rdeireuy"
	text := "hello world! 测试 ***** ¥¥ 44 ^ ^ `` ) ( - + @"
	encTxt, _ := EncryptECB([]byte(key), []byte(text))
	decTxt, _ := DecryptECB([]byte(key), encTxt)
	if text == string(decTxt) {
		log.Printf("sm4 ecb success")
	}

	encTxt, _ = EncryptCBC([]byte(key), []byte(text))
	decTxt, _ = DecryptCBC([]byte(key), encTxt)
	if text == string(decTxt) {
		log.Printf("sm4 cbc success")
	}

	encTxt, _ = EncryptCFB([]byte(key), []byte(text))
	decTxt, _ = DecryptCFB([]byte(key), encTxt)
	if text == string(decTxt) {
		log.Printf("sm4 cfb success")
	}

	encTxt, _ = EncryptOFB([]byte(key), []byte(text))
	decTxt, _ = DecryptOFB([]byte(key), encTxt)
	if text == string(decTxt) {
		log.Printf("sm4 ofb success")
	}

	encTxt, _ = EncryptSm4([]byte(key), []byte(iv), []byte(text))
	decTxt, _ = DecryptSm4([]byte(key), []byte(iv), encTxt)
	if text == string(decTxt) {
		log.Printf("sm4 success")
	}

}
