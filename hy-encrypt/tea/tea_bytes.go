// 以 bytes 的方式处理
// Less than 8 bytes will not be processed
package tea

import (
	"encoding/base64"
)

func BytesEncryptStr2StrWithB64Key(b64Key string, plainText string) (string, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(b64Key)
	if err != nil {
		return "", err
	}
	return StrEncryptStr2Str(keyBytes, plainText)
}

// 加密结果时经过base64的
func BytesEncryptStr2Str(keyBytes []byte, plainText string) (string, error) {
	encBytes, err := StrEncryptBytes2Bytes(keyBytes, []byte(plainText))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encBytes), nil
}

// 加密任意内容,key必须是16个字节的Base64串,加密内容以8byte为一组,不足的不处理
func BytesEncryptBytes2Bytes(keyBytes []byte, plainBytes []byte) ([]byte, error) {
	// 8个一组 分组加密
	srcLen := len(plainBytes)
	divisor := srcLen / 8       // 计算出整除的组
	buf := make([]byte, srcLen) // 调整后的数据buf
	copy(buf, plainBytes)

	cipher, err := NewCipherWithRounds(keyBytes, 32)
	for i := 0; i < divisor; i++ {
		if err != nil {
			return nil, err
		}
		start := i * 8
		end := start + 8
		cipher.Encrypt(buf[start:end], buf[start:end])
	}
	return buf[:], nil
}

func BytesDecryptStr2StrWithB64Key(b64Key string, cipherText string) (string, error) {

	keyBytes, err := base64.StdEncoding.DecodeString(b64Key)
	if err != nil {
		return "", err
	}
	return StrDecryptStr2Str(keyBytes, cipherText)

}

func BytesDecryptStr2Str(keyBytes []byte, cipherText string) (string, error) {
	cipherBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	bytes, err := StrDecryptBytes2Bytes(keyBytes, cipherBytes)
	return string(bytes), err
}

// 解密任意内容,key必须是16个字节
// 解密的内容是加密后的原始内容
// 解密返回的是原字符串
func BytesDecryptBytes2Bytes(keyBytes []byte, cipherBytes []byte) ([]byte, error) {

	// 8个一组 分组加密
	srcLen := len(cipherBytes)
	divisor := srcLen / 8       // 计算出整除的组
	buf := make([]byte, srcLen) // 调整后的数据buf
	copy(buf, cipherBytes)

	cipher, err := NewCipherWithRounds(keyBytes, 32)
	for i := 0; i < divisor; i++ {
		if err != nil {
			return nil, err
		}
		start := i * 8
		end := start + 8
		cipher.Decrypt(buf[start:end], buf[start:end])
	}
	return buf, nil
}
