package sm4

import "github.com/tjfoc/gmsm/sm4"

func EncryptOFB(key []byte, in []byte) ([]byte, error) {
	return sm4.Sm4OFB(key, in, true)
}

func DecryptOFB(key []byte, in []byte) ([]byte, error) {
	return sm4.Sm4OFB(key, in, false)
}
