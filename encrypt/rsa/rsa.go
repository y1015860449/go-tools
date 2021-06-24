package rsa

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
)

func GenRsaKey(bits int) ([]byte, []byte, error) {
	key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	priBytes := x509.MarshalPKCS1PrivateKey(key)
	pubBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		return nil, nil, err
	}
	return priBytes, pubBytes, nil
}

// 获取公私钥对
func GenRsaKeyB64(bits int) (string, string, error) {
	pub, pri, err := GenRsaKey(bits)
	if err != nil {
		return "", "", err
	}
	return base64.StdEncoding.EncodeToString(pri), base64.StdEncoding.EncodeToString(pub), nil
}

// 公钥加密
func EncryptWithBase64Key(origData []byte, pubKey string) ([]byte, error) {
	block, err := base64.StdEncoding.DecodeString(pubKey)
	if block == nil || err != nil {
		return nil, errors.New("public key error")
	}
	return Encrypt(origData, block)
}

func Encrypt(origData []byte, pubKey []byte) ([]byte, error) {
	pubInterface, err := x509.ParsePKIXPublicKey(pubKey)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)

	partLen := pub.N.BitLen()/8 - 11
	chunks := split(origData, partLen)
	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		bytes, err := rsa.EncryptPKCS1v15(rand.Reader, pub, chunk)
		if err != nil {
			return nil, err
		}
		buffer.Write(bytes)
	}
	return buffer.Bytes(), nil
}

// 私钥解密
func DecryptWithBase64Key(cipherText []byte, priKey string) ([]byte, error) {
	block, err := base64.StdEncoding.DecodeString(priKey)
	if block == nil || err != nil {
		return nil, errors.New("private key error")
	}
	pri, err := x509.ParsePKCS1PrivateKey(block)
	if err != nil {
		return nil, err
	}
	partLen := pri.PublicKey.N.BitLen() / 8
	chunks := split(cipherText, partLen)
	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, pri, chunk)
		if err != nil {
			return nil, err
		}
		buffer.Write(decrypted)
	}
	return buffer.Bytes(), nil
}

func Decrypt(cipherText []byte, priKey []byte) ([]byte, error) {
	pri, err := x509.ParsePKCS1PrivateKey(priKey)
	if err != nil {
		return nil, err
	}
	partLen := pri.PublicKey.N.BitLen() / 8
	chunks := split(cipherText, partLen)
	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, pri, chunk)
		if err != nil {
			return nil, err
		}
		buffer.Write(decrypted)
	}
	return buffer.Bytes(), nil
}

func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}
