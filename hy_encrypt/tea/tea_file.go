package tea

import (
	"bufio"
	"encoding/base64"
	"errors"
	"io"
	"os"
)

const (
	SimpleEncryptLen = 32
)

// 简单加密
func SimpleEncryptFileWithB64Key(b64Key string, inFile, outFile string) error {
	return encryptFileWithB64Key(b64Key, inFile, outFile, SimpleEncryptLen)
}

func SimpleEncryptFileWithKey(key []byte, inFile, outFile string) error {
	return encryptFileWithKey(key, inFile, outFile, SimpleEncryptLen)
}

// 简单解密
func SimpleDecryptFileWithB64Key(b64Key string, inFile, outFile string) error {
	return decryptFileWithB64Key(b64Key, inFile, outFile, SimpleEncryptLen)
}

func SimpleDecryptFileWithKey(key []byte, inFile, outFile string) error {
	return decryptFileWithKey(key, inFile, outFile, SimpleEncryptLen)
}

// 完整加密(除去不够8byte的部分)
func EncryptFileWithB64Key(b64Key string, inFile, outFile string) error {
	return encryptFileWithB64Key(b64Key, inFile, outFile, 0)
}

func EncryptFileWithKey(key []byte, inFile, outFile string) error {
	return encryptFileWithKey(key, inFile, outFile, 0)
}

// 完整解密(除去不够8byte的部分)
func DecryptFileWithB64Key(b64Key string, inFile, outFile string) error {
	return decryptFileWithB64Key(b64Key, inFile, outFile, 0)
}

func DecryptFileWithKey(key []byte, inFile, outFile string) error {
	return decryptFileWithKey(key, inFile, outFile, 0)
}

func encryptFileWithB64Key(b64Key string, inFile, outFile string, maxEncLen uint32) error {

	if b64Key == "" || inFile == "" || outFile == "" {
		return errors.New("pram except")
	}

	keyBytes, err := base64.StdEncoding.DecodeString(b64Key)
	if err != nil {
		return err
	}
	return encryptFileWithKey(keyBytes, inFile, outFile, maxEncLen)
}

func encryptFileWithKey(key []byte, inFile, outFile string, maxEncLen uint32) error {
	if len(key) <= 0 || len(inFile) <= 0 || len(outFile) <= 0 {
		return errors.New("pram except")
	}
	cipher, err := NewCipherWithRounds(key, 32)
	if err != nil {
		return err
	}
	return encryptAndDecryptFileWithB64Key(cipher.Encrypt, inFile, outFile, maxEncLen)
}

func decryptFileWithB64Key(b64Key string, inFile, outFile string, maxEncLen uint32) error {

	if b64Key == "" || inFile == "" || outFile == "" {
		return errors.New("pram except")
	}

	keyBytes, err := base64.StdEncoding.DecodeString(b64Key)
	if err != nil {
		return err
	}
	return decryptFileWithKey(keyBytes, inFile, outFile, maxEncLen)
}

func decryptFileWithKey(key []byte, inFile, outFile string, maxEncLen uint32) error {
	if len(key) <= 0 || len(inFile) <= 0 || len(outFile) <= 0 {
		return errors.New("pram except")
	}
	cipher, err := NewCipherWithRounds(key, 32)
	if err != nil {
		return err
	}
	return encryptAndDecryptFileWithB64Key(cipher.Decrypt, inFile, outFile, maxEncLen)
}

func encryptAndDecryptFileWithB64Key(f func(dst, src []byte), inFile, outFile string, maxEncLen uint32) error {

	inf, err := os.Open(inFile)
	if err != nil {
		return err
	}
	defer inf.Close()

	_ = os.Remove(outFile)
	outf, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer outf.Close()

	reader := bufio.NewReader(inf)
	writer := bufio.NewWriter(outf)
	defer writer.Flush()

	buf := make([]byte, 8)
	var encLen uint32
	const maxBuf = 10240
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			needEnc := (n == 8) && (maxEncLen == 0 || encLen < maxEncLen)
			if needEnc { // 需要加密
				f(buf, buf[:n])
			}
			if _, err := writer.Write(buf[:n]); err != nil {
				return err
			}
			encLen += uint32(n)
			if !needEnc && len(buf) != maxBuf {
				buf = make([]byte, maxBuf) // 不需要分组加密是扩大buf
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}
	return nil
}
