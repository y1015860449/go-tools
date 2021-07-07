package sm2

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
)

func GenerateSm2Key() ([]byte, []byte, error) {
	key, err := sm2.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	priBytes, err := x509.MarshalSm2PrivateKey(key, nil)
	if err != nil {
		return nil, nil, err
	}
	pubBytes, err := x509.MarshalSm2PublicKey(&key.PublicKey)
	if err != nil {
		return nil, nil, err
	}
	return priBytes, pubBytes, nil
}

func GenerateSm2KeyByB64() (string, string, error) {
	pri, pub, err := GenerateSm2Key()
	if err != nil {
		return "", "", err
	}
	return base64.StdEncoding.EncodeToString(pri), base64.StdEncoding.EncodeToString(pub), nil
}

func Encrypt(origData, pubKey []byte) ([]byte, error) {
	pub, err := x509.ParseSm2PublicKey(pubKey)
	if err != nil {
		return nil, err
	}
	return sm2.Encrypt(pub, origData, rand.Reader)
}

func Decrypt(cipherText, priKey []byte) ([]byte, error) {
	pri, err := x509.ParsePKCS8PrivateKey(priKey, nil)
	if err != nil {
		return nil, err
	}
	return sm2.Decrypt(pri, cipherText)
}

func EncryptWithB64Key(origData []byte, pubKey string) ([]byte, error) {
	pub, err := base64.StdEncoding.DecodeString(pubKey)
	if err != nil {
		return nil, err
	}
	return Encrypt(origData, pub)
}

func DecryptWithB64Key(cipherText []byte, priKey string) ([]byte, error) {
	pri, err := base64.StdEncoding.DecodeString(priKey)
	if err != nil {
		return nil, err
	}
	return Decrypt(cipherText, pri)
}
