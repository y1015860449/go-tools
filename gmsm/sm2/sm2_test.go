package sm2

import (
	"encoding/hex"
	"log"
	"testing"
)


var priKey string
var pubKey string

func TestDecryptWithB64Key(t *testing.T) {
	text := "hello world! 测试 ***** ¥¥ 44 ^ ^ `` ) ( - + @"
	priKey, pubKey, _ := GenerateSm2KeyByB64()
	log.Printf("pub %s len %d \npri %s len %d", pubKey, len(pubKey), priKey, len(priKey))
	encryptData, _ := EncryptWithB64Key([]byte(text), pubKey)
	log.Printf("encrypt data %s", hex.EncodeToString(encryptData))
	decryptData, _ := DecryptWithB64Key(encryptData, priKey)
	log.Printf("decrypt data %s", string(decryptData))
	if text == string(decryptData) {
		log.Printf("success")
	}
}
