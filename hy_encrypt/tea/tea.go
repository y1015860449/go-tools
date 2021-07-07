package tea

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"math/rand"
	"time"
)

var (
	errIncorrectKeySize = errors.New("tea: incorrect key size")
	errDataTooShort     = errors.New("data too short")
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenRandB64TeaKey() string {
	return base64.StdEncoding.EncodeToString(genRandBytes(16))
}

// 生成指定位数的随机Bytes
func genRandBytes(l int) []byte {
	if l%8 != 0 {
		panic("a must be a multiple of 8")
	}
	var key []byte
	tmp := make([]byte, 8)
	for i := 0; i < l; i += 8 {
		binary.BigEndian.PutUint64(tmp, rand.Uint64())
		key = append(key, tmp...)
	}
	return key
}

func GenRandTeaKey() []byte {
	return genRandBytes(16)
}

func Encrypt(keyBytes []byte, plainBytes []byte) ([]byte, error) {
	if len(keyBytes) != 16 {
		return nil, errIncorrectKeySize
	}
	if len(plainBytes) < 8 {
		return nil, errDataTooShort
	}
	cipher, err := NewCipherWithRounds(keyBytes, 32)
	if err != nil {
		return nil, err
	}
	cipher.Encrypt(plainBytes, plainBytes)
	return plainBytes, nil
}

func Decrypt(keyBytes []byte, cipherBytes []byte) ([]byte, error) {
	if len(keyBytes) != 16 {
		return nil, errIncorrectKeySize
	}
	if len(cipherBytes) < 8 {
		return nil, errDataTooShort
	}
	cipher, err := NewCipherWithRounds(keyBytes, 32)
	if err != nil {
		return nil, err
	}
	cipher.Encrypt(cipherBytes, cipherBytes)
	return cipherBytes, nil
}
