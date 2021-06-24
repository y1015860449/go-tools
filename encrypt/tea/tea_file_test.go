package tea

import (
	"bufio"
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"go-tools/encrypt/md5"
	"math/rand"
	"os"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestTeaFile(t *testing.T) {
	testTeaFile(t, EncryptFileWithB64Key, DecryptFileWithB64Key)
	testTeaFile(t, SimpleEncryptFileWithB64Key, SimpleDecryptFileWithB64Key)
}

func testTeaFile(
	t *testing.T,
	encFun func(b64Key string, inFile, outFile string) error,
	decFun func(b64Key string, inFile, outFile string) error,
) {
	key := GenRandB64TeaKey()

	srcfile := "src.data"
	encfile := "encrypt.data"
	decfile := "decrypt.data"

	defer os.Remove(srcfile)
	defer os.Remove(encfile)
	defer os.Remove(decfile)

	// 测试1byte至数M大小的文件
	for i := 0; i < 24; i++ {
		if err := createRandFile(srcfile, 1<<i); err != nil {
			t.FailNow()
		}
		err := encFun(key, srcfile, encfile)
		assert.Nil(t, err)

		err = decFun(key, encfile, decfile)
		assert.Nil(t, err)

		assertFile(t, srcfile, decfile)
	}
}

func assertFile(t *testing.T, srcFile string, dstFile string) {
	srcMd5, err := md5.FileMd5(srcFile)
	assert.Nil(t, err)
	dstMd5, err := md5.FileMd5(dstFile)
	assert.Nil(t, err)
	assert.Equal(t, srcMd5, dstMd5)
}

// 创建指定长度的随机内容文件
func createRandFile(file string, len int) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	writer := bufio.NewWriter(f)
	defer writer.Flush()
	buf := make([]byte, 8)
	for i := 0; i < len; {
		binary.BigEndian.PutUint64(buf, rand.Uint64())
		if len < 8 {
			_, err := writer.Write(buf[:len])
			return err
		}
		n, err := writer.Write(buf)
		if err != nil {
			return err
		}
		i += n
	}
	return nil
}
