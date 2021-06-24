package sm4

import "github.com/tjfoc/gmsm/sm4"

func EncryptCFB(key []byte, in []byte) ([]byte, error) {
	return sm4.Sm4CFB(key, in, true)
}

func DecryptCFB(key []byte, in []byte) ([]byte, error) {
	return sm4.Sm4CFB(key, in, false)
}
