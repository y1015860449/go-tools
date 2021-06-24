package tea

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

var b64TeaKey = GenRandB64TeaKey()

func TestTea(t *testing.T) {

	for i := 0; i < 100000; i++ {

		plainText := fmt.Sprintf("%d", i)
		cipherText, err := StrEncryptStr2StrWithB64Key(b64TeaKey, plainText)
		assert.Nil(t, err)

		decryptText, err := StrDecryptStr2StrWithB64Key(b64TeaKey, cipherText)
		if err != nil || plainText != decryptText {
			t.FailNow()
		}
	}
}

// 测试对比app目前的加解密
func TestTeaWithApp(t *testing.T) {
	teaKey := "MVvnfm6bFtrgfzYBlPA8Gg=="
	plainText := "一名魔兽"
	encryptText, err := StrEncryptStr2StrWithB64Key(teaKey, plainText)
	assert.Nil(t, err)
	assert.Equal(t, "ONs5KWUUhu2GKQyFD2H9MA==", encryptText)

	decryptText, err := StrDecryptStr2StrWithB64Key(teaKey, encryptText)
	assert.Nil(t, err)

	assert.EqualValues(t, plainText, decryptText)
}

func TestGenRandKey(t *testing.T) {
	record := make(map[string]struct{})
	// 生成一批 随机bytes 不能重复
	for i := 0; i < 10000; i++ {
		key := GenRandB64TeaKey()
		if _, ok := record[key]; ok || len(key) == 0 {
			t.FailNow()
			return
		}
		record[key] = struct{}{}
	}
}

func Test_genRandBytes(t *testing.T) {
	for i := 8; i <= math.MaxInt8; i += 8 {
		keyBytes := genRandBytes(i)
		assert.EqualValues(t, i, len(keyBytes))
	}
}
