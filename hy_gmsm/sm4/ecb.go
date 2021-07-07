package sm4

import "github.com/tjfoc/gmsm/sm4"

func EncryptECB(key []byte, in []byte) ([]byte, error) {
	return sm4.Sm4Ecb(key, in, true)
}

func DecryptECB(key []byte, in []byte) ([]byte, error) {
	return sm4.Sm4Ecb(key, in, false)
}
